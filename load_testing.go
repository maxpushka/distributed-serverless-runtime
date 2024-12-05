package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"serverless/router/schema"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func main() {
	// CSV file setup
	csvFile, err := os.Create("load_test_results.csv")
	if err != nil {
		log.Fatal("Could not create CSV file:", err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write CSV headers
	err = writer.Write([]string{"Endpoint", "Total Requests", "Successful Requests", "Failed Requests", "Error Rate (%)", "Median Response Time (ms)", "95th Percentile Response Time (ms)"})
	if err != nil {
		log.Fatal("Could not write CSV headers:", err)
	}

	// Base URL of your server
	baseURL := "http://localhost:8080" // Replace with your server's URL

	// Step 1: Create an admin user with password 'admin'
	fmt.Println("Creating admin user...")
	adminUser := schema.Credentials{Username: "admin", Password: "admin"}
	err = createUser(baseURL, adminUser)
	if err != nil {
		log.Fatal("Failed to create admin user:", err)
	}

	// Step 2: Stress test of creating users with random names and passwords for 5 minutes
	fmt.Println("Starting stress test: Creating users...")
	//err = stressTestCreateUsers(baseURL, 5*time.Minute, writer)
	err = stressTestCreateUsers(baseURL, 5*time.Second, writer)
	if err != nil {
		log.Fatal("Stress test failed:", err)
	}

	// Step 3: Stress test of logging in admin user for 5 minutes
	fmt.Println("Starting stress test: Logging in admin user...")
	//err = stressTestLoginUser(baseURL, adminUser, 5*time.Minute, writer)
	err = stressTestLoginUser(baseURL, adminUser, 5*time.Second, writer)
	if err != nil {
		log.Fatal("Stress test failed:", err)
	}

	// Step 4: Create a route for admin user with the name 'test_route'
	fmt.Println("Creating route 'test_route' for admin user...")
	routeID, err := createRoute(baseURL, adminUser, "test_route")
	if err != nil {
		log.Fatal("Failed to create route:", err)
	}

	// Step 5: Stress test for creating routes and getting routes data for admin user for 5 minutes
	fmt.Println("Starting stress test: Creating routes...")
	//err = stressTestCreateRoutes(baseURL, adminUser, 5*time.Minute, writer)
	err = stressTestCreateRoutes(baseURL, adminUser, 5*time.Second, writer)
	if err != nil {
		log.Fatal("Stress test failed:", err)
	}

	// Step 6: Stress test for creating routes and getting routes data for admin user for 5 minutes
	fmt.Println("Starting stress test: Getting routes...")
	//err = stressTestGetRoutes(baseURL, adminUser, 5*time.Minute, writer)
	err = stressTestGetRoutes(baseURL, adminUser, 5*time.Second, writer)
	if err != nil {
		log.Fatal("Stress test failed:", err)
	}

	// Step 7: Upload config and executable into 'test_route'
	fmt.Println("Uploading config and executable to 'test_route'...")
	err = uploadExecutable(baseURL, adminUser, routeID)
	if err != nil {
		log.Fatal("Failed to upload executable:", err)
	}

	// Step 8: Run execute stress test for 10 minutes
	fmt.Println("Starting stress test: Executing route...")
	//err = stressTestExecuteRoute(baseURL, adminUser, routeID, 10*time.Minute, writer)
	err = stressTestExecuteRoute(baseURL, adminUser, routeID, 5*time.Second, writer)
	if err != nil {
		log.Fatal("Stress test failed:", err)
	}

	fmt.Println("Load testing completed. Results are saved in 'load_test_results.csv'.")
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Function to create a user
func createUser(baseURL string, user schema.Credentials) error {
	url := baseURL + "/register"
	body, _ := json.Marshal(user)
	resp, err := httpPost(url, body, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		return fmt.Errorf("Failed to create user: %s", resp.Status)
	}
	return nil
}

// Function to perform stress test of creating users
func stressTestCreateUsers(baseURL string, duration time.Duration, writer *csv.Writer) error {
	targeter := NewCreateUserTargeter(baseURL)
	return runStressTest(targeter, duration, "Create User", writer)
}

func NewCreateUserTargeter(baseURL string) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		username := randString(16) + "_user_" + randString(16)
		password := randString(16) + "_pass_" + randString(16)
		user := schema.Credentials{Username: username, Password: password}
		body, _ := json.Marshal(user)
		tgt.Method = "POST"
		tgt.URL = baseURL + "/register"
		tgt.Body = body
		tgt.Header = map[string][]string{"Content-Type": {"application/json"}}
		return nil
	}
}

// Function to perform stress test of logging in user
func stressTestLoginUser(baseURL string, user schema.Credentials, duration time.Duration, writer *csv.Writer) error {
	body, _ := json.Marshal(user)
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    baseURL + "/login",
		Body:   body,
		Header: map[string][]string{"Content-Type": {"application/json"}},
	})
	return runStressTest(targeter, duration, "Login Admin", writer)
}

// Function to create a route
func createRoute(baseURL string, user schema.Credentials, routeName string) (string, error) {
	// Login to get the token
	token, err := loginUser(baseURL, user)
	if err != nil {
		return "", err
	}

	url := baseURL + "/api/routes"
	routeData := map[string]string{"name": routeName}
	body, _ := json.Marshal(routeData)
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}
	resp, err := httpPost(url, body, headers)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != 201 {
		return "", fmt.Errorf("Failed to create route: %s", resp.Status)
	}

	// Decode the JSON response
	var responseData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return "", fmt.Errorf("Failed to decode response: %v", err)
	}

	// Access the "data" field
	data, ok := responseData["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Data field not found in response")
	}

	// Extract the "id" from "data"
	routeID, ok := data["id"].(string)
	if !ok {
		// If "id" is a number, handle accordingly
		idFloat, ok := data["id"].(float64)
		if !ok {
			return "", fmt.Errorf("Route ID not found or invalid in response data")
		}
		routeID = fmt.Sprintf("%.0f", idFloat)
	}

	return routeID, nil
}

// Function to perform stress test of creating routes
func stressTestCreateRoutes(baseURL string, user schema.Credentials, duration time.Duration, writer *csv.Writer) error {
	// Login to get the token
	token, err := loginUser(baseURL, user)
	if err != nil {
		return err
	}

	headers := map[string][]string{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + token},
	}

	// Create route targeter
	createRouteTargeter := NewCreateRouteTargeter(baseURL, headers)

	// Run stress test
	err = runStressTest(createRouteTargeter, duration, "Create Routes", writer)
	if err != nil {
		return err
	}

	return nil
}

// Function to perform stress test of getting routes
func stressTestGetRoutes(baseURL string, user schema.Credentials, duration time.Duration, writer *csv.Writer) error {
	// Login to get the token
	token, err := loginUser(baseURL, user)
	if err != nil {
		return err
	}

	headers := map[string][]string{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + token},
	}

	// Get routes targeter
	getRoutesTargeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    baseURL + "/api/routes",
		Header: headers,
	})

	// Run stress test
	err = runStressTest(getRoutesTargeter, duration, "Get Routes", writer)
	if err != nil {
		return err
	}

	return nil
}

func NewCreateRouteTargeter(baseURL string, headers map[string][]string) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		routeName := "route_" + randString(8)
		routeData := map[string]string{"name": routeName}
		body, _ := json.Marshal(routeData)
		tgt.Method = "POST"
		tgt.URL = baseURL + "/api/routes"
		tgt.Body = body
		tgt.Header = headers
		return nil
	}
}

// Function to upload config and executable
func uploadExecutable(baseURL string, user schema.Credentials, routeID string) error {
	token, err := loginUser(baseURL, user)
	if err != nil {
		return err
	}
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	execURL := fmt.Sprintf("%s/api/routes/%s/executable", baseURL, routeID)
	execFilePath := "test_files/executable.js" // Update this path with the actual executable file path
	execFile, err := os.Open(execFilePath)
	if err != nil {
		return fmt.Errorf("Failed to open executable file: %v", err)
	}
	defer execFile.Close()

	// Upload the file using httpPostFile
	resp, err := httpPostFile(execURL, nil, "file", "executable_filename", execFile, headers)
	if err != nil {
		return fmt.Errorf("Failed to upload executable: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to upload executable: %s", resp.Status)
	}

	return nil
}

// Function to perform stress test of executing route
func stressTestExecuteRoute(baseURL string, user schema.Credentials, routeID string, duration time.Duration, writer *csv.Writer) error {
	token, err := loginUser(baseURL, user)
	if err != nil {
		return err
	}
	headers := map[string][]string{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + token},
	}
	executeURL := fmt.Sprintf("%s/api/routes/%s/execute", baseURL, routeID)
	routeData := schema.ExecuteParams{
		RouteConfig: map[string]string{"adjective": "wonderful"},
		RequestBody: map[string]string{"name": "world"},
	}
	body, _ := json.Marshal(routeData)
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    executeURL,
		Header: headers,
		Body:   body,
	})
	return runStressTest(targeter, duration, "Execute Route", writer)
}

// Function to run a stress test using Vegeta
func runStressTest(targeter vegeta.Targeter, duration time.Duration, endpoint string, writer *csv.Writer) error {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second} // 100 requests per second
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics

	for res := range attacker.Attack(targeter, rate, duration, endpoint) {
		metrics.Add(res)
	}
	metrics.Close()

	// Calculate error rate
	errorRate := (1 - metrics.Success) * 100

	// Write results to CSV
	record := []string{
		endpoint,
		fmt.Sprintf("%d", metrics.Requests),
		fmt.Sprintf("%d", int(float64(metrics.Requests)*metrics.Success)),
		fmt.Sprintf("%d", int(float64(metrics.Requests)*(1-metrics.Success))),
		fmt.Sprintf("%.2f", errorRate),
		fmt.Sprintf("%.2f", float64(metrics.Latencies.P50)/float64(time.Millisecond)),
		fmt.Sprintf("%.2f", float64(metrics.Latencies.P95)/float64(time.Millisecond)),
	}
	err := writer.Write(record)
	if err != nil {
		return fmt.Errorf("Could not write to CSV: %v", err)
	}

	return nil
}

// Helper function to login user and get token
func loginUser(baseURL string, user schema.Credentials) (string, error) {
	url := baseURL + "/login"
	body, _ := json.Marshal(user)
	resp, err := httpPost(url, body, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Failed to login: %s", resp.Status)
	}
	var responseData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return "", fmt.Errorf("Failed to decode response: %v", err)
	}
	data, ok := responseData["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Data field not found in response")
	}
	token, ok := data["token"].(string)
	if !ok {
		return "", fmt.Errorf("Token not found in response data")
	}
	return token, nil
}

// Helper function to perform HTTP POST requests
func httpPost(url string, body []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	return client.Do(req)
}

// Helper function to perform HTTP POST requests with file upload
func httpPostFile(url string, params map[string]string, fileParamName, fileName string, file io.Reader, headers map[string]string) (*http.Response, error) {
	// Create a buffer to hold the form data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add the file to the form data
	fw, err := w.CreateFormFile(fileParamName, fileName)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return nil, err
	}

	// Add other form fields if necessary
	for key, val := range params {
		if err = w.WriteField(key, val); err != nil {
			return nil, err
		}
	}

	// Close the writer to finalize the form data
	w.Close()

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return nil, err
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", w.FormDataContentType())
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	client := &http.Client{}
	return client.Do(req)
}
