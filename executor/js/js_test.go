package js_test

import (
	"io"
	"reflect"
	"serverless/config"
	"strings"
	"testing"
	"time"

	"serverless/executor"
	"serverless/executor/js"
)

type mockQuery struct {
	readFile func(id string) (content io.Reader, checksum string, err error)
	checksum func(id string) (string, error)
}

func (m *mockQuery) ReadFile(id string) (content io.Reader, checksum string, err error) {
	return m.readFile(id)
}

func (m *mockQuery) Checksum(id string) (string, error) {
	return m.checksum(id)
}

func TestExecutor_Execute(t *testing.T) {
	type args struct {
		id  string
		cfg map[string]string
		req executor.Request
	}
	tests := []struct {
		name    string
		source  mockQuery
		args    args
		want    executor.Response
		wantErr bool
	}{
		{
			name: "Demo",
			source: mockQuery{
				readFile: func(id string) (content io.Reader, checksum string, err error) {
					const source = `
						function handle(req, cfg) {
							return {
								status_code: 200,
								headers: {
									"Content-Type": "application/json"
								},
								body: JSON.stringify({
									message: "Hello, " + cfg["adjective"] + " " + req.body["name"] + "!",
								}),
							}
						}
					`
					return strings.NewReader(source), "123", nil
				},
				checksum: func(id string) (string, error) {
					return "123", nil
				},
			},
			args: args{
				id: "test/echo", // assumed to adhere format <username>/<http_route>
				cfg: map[string]string{
					"adjective": "wonderful",
				},
				req: executor.Request{
					Method: "GET",
					URL:    "https://example.com",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: map[string]string{"name": "world"},
				},
			},
			want: executor.Response{
				StatusCode: 200,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: `{"message":"Hello, wonderful world!"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := js.NewExecutor(config.Executor{HotDuration: 30 * time.Minute}, &tt.source)
			got, err := exec.Execute(tt.args.id, tt.args.cfg, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute() got = %v, want %v", got, tt.want)
			}
		})
	}
}
