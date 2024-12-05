package schema

type ExecuteParams struct {
	RequestBody any               `json:"requestBody"`
	RouteConfig map[string]string `json:"routeConfig"`
}
