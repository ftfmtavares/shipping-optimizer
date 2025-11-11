package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/ftfmtavares/shipping-optimizer/internal/instrumentation"
	"github.com/stretchr/testify/assert"
)

func testRequest(t *testing.T, method, target string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, target, body)
	assert.NoError(t, err)

	req.Header.Add("Accept", `application/json`)
	return req
}

func TestHTTPServer(t *testing.T) {
	logger := instrumentation.NewLogger()

	testCases := []struct {
		desc                string
		serverConfiguration func(*HTTPServer)
		request             *http.Request
		expectedResponse    func(assert.TestingT, *http.Response)
	}{
		{
			desc:                "server with no routes",
			serverConfiguration: func(server *HTTPServer) {},
			request:             testRequest(t, http.MethodGet, "http://localhost:8000/health", nil),
			expectedResponse: func(tt assert.TestingT, res *http.Response) {
				assert.Equal(tt, http.StatusNotFound, res.StatusCode)
			},
		},
		{
			desc: "health check",
			serverConfiguration: func(server *HTTPServer) {
				server.WithHealthCheck()
			},
			request: testRequest(t, http.MethodGet, "http://localhost:8000/health", nil),
			expectedResponse: func(tt assert.TestingT, res *http.Response) {
				assert.Equal(tt, http.StatusOK, res.StatusCode)

				body, err := io.ReadAll(res.Body)
				assert.NoError(tt, err)

				var responseMap map[string]string
				err = json.Unmarshal(body, &responseMap)
				assert.NoError(tt, err)
				assert.Equal(tt, "healthy", responseMap["status"])
				assert.NotEmpty(tt, responseMap["time"])
			},
		},
		{
			desc: "routed service",
			serverConfiguration: func(server *HTTPServer) {
				server.WithServiceHandler("/test", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte(`{"ok":"true"}`))
				}, http.MethodGet)
			},
			request: testRequest(t, http.MethodGet, "http://localhost:8000/test", nil),
			expectedResponse: func(tt assert.TestingT, res *http.Response) {
				assert.Equal(tt, http.StatusAccepted, res.StatusCode)

				assert.Equal(tt, "application/json", res.Header.Get("Content-Type"))

				body, err := io.ReadAll(res.Body)
				assert.NoError(tt, err)

				var responseMap map[string]string
				err = json.Unmarshal(body, &responseMap)
				assert.NoError(tt, err)
				assert.Equal(tt, "true", responseMap["ok"])
			},
		},
		{
			desc: "static web page",
			serverConfiguration: func(server *HTTPServer) {
				server.WithStatic("/", "./testdata")
			},
			request: testRequest(t, http.MethodGet, "http://localhost:8000/test.html", nil),
			expectedResponse: func(tt assert.TestingT, res *http.Response) {
				assert.Equal(tt, http.StatusOK, res.StatusCode)

				assert.Equal(tt, "text/html; charset=utf-8", res.Header.Get("Content-Type"))

				body, err := io.ReadAll(res.Body)
				assert.NoError(tt, err)
				assert.Equal(tt, "<!DOCTYPE html><html lang=\"en\"></html>", string(body))
			},
		},
		{
			desc: "service panic",
			serverConfiguration: func(server *HTTPServer) {
				server.WithServiceHandler("/test", func(w http.ResponseWriter, r *http.Request) {
					panic("something bad")
				}, http.MethodGet)
			},
			request: testRequest(t, http.MethodGet, "http://localhost:8000/test", nil),
			expectedResponse: func(tt assert.TestingT, res *http.Response) {
				assert.Equal(tt, http.StatusInternalServerError, res.StatusCode)

				assert.Equal(tt, "application/json", res.Header.Get("Content-Type"))
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			s := NewHTTPServer(HTTPServerConfig{
				Address:      "localhost",
				Port:         8000,
				ReadTimeout:  time.Second * 30,
				WriteTimeout: time.Second * 30,
				IdleTimeout:  time.Second * 30,
				Logger:       logger,
			})

			tC.serverConfiguration(&s)

			s.StartHTTPServerAsync()
			time.Sleep(200 * time.Millisecond)

			client := &http.Client{Timeout: 2 * time.Second}
			resp, err := client.Do(tC.request)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.NotNil(t, resp)
			tC.expectedResponse(t, resp)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err = s.server.Shutdown(ctx)
			assert.NoError(t, err)
		})
	}
}
