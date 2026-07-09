package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"signalstack-ai/backend/config"
	"signalstack-ai/backend/routes"
)

func TestIngestionChatAndEvaluationEndpoints(t *testing.T) {
	router := NewRouter(config.Config{}, routes.Dependencies{})

	cases := []struct {
		name   string
		method string
		path   string
	}{
		{name: "ingest all", method: http.MethodPost, path: "/api/ingest/all"},
		{name: "chat ask", method: http.MethodPost, path: "/api/chat/ask"},
		{name: "evaluation run", method: http.MethodPost, path: "/api/evaluate/run"},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(testCase.method, testCase.path, nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d", recorder.Code)
			}
		})
	}
}
