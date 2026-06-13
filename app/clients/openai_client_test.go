package clients

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChatCompleteSuccess(t *testing.T) {
	var captured openAIRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization header = %q, want Bearer test-key", got)
		}
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &captured)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"choices":[{"message":{"role":"assistant","content":"{\"ok\":true}"}}]}`)
	}))
	defer server.Close()

	client := newTestClient(server.URL, "test-key", "gpt-4o-mini")
	got, err := client.ChatComplete(context.Background(), "sys prompt", "user prompt")
	if err != nil {
		t.Fatalf("ChatComplete() error = %v", err)
	}
	if got != `{"ok":true}` {
		t.Fatalf("ChatComplete() = %q, want JSON content", got)
	}

	// Request must force JSON output and carry both messages.
	if captured.ResponseFormat.Type != "json_object" {
		t.Errorf("response_format.type = %q, want json_object", captured.ResponseFormat.Type)
	}
	if len(captured.Messages) != 2 || captured.Messages[0].Role != "system" || captured.Messages[1].Role != "user" {
		t.Errorf("messages = %#v, want system+user", captured.Messages)
	}
	if captured.Model != "gpt-4o-mini" {
		t.Errorf("model = %q, want gpt-4o-mini", captured.Model)
	}
}

func TestChatCompleteAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"error":{"message":"Invalid API key"}}`)
	}))
	defer server.Close()

	client := newTestClient(server.URL, "bad-key", "gpt-4o-mini")
	_, err := client.ChatComplete(context.Background(), "sys", "user")
	if err == nil {
		t.Fatal("ChatComplete() should return error on API error payload")
	}
	if !strings.Contains(err.Error(), "Invalid API key") {
		t.Fatalf("ChatComplete() error = %v, want it to surface the API message", err)
	}
}

func TestChatCompleteNoChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"choices":[]}`)
	}))
	defer server.Close()

	client := newTestClient(server.URL, "k", "gpt-4o-mini")
	_, err := client.ChatComplete(context.Background(), "sys", "user")
	if err == nil {
		t.Fatal("ChatComplete() should error when no choices are returned")
	}
}

func TestChatCompleteMalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `not json at all`)
	}))
	defer server.Close()

	client := newTestClient(server.URL, "k", "gpt-4o-mini")
	_, err := client.ChatComplete(context.Background(), "sys", "user")
	if err == nil {
		t.Fatal("ChatComplete() should error on malformed response body")
	}
}

// newTestClient builds a client pointed at a test server.
func newTestClient(url, apiKey, model string) *OpenAIClient {
	c := NewOpenAIClient(apiKey, model, 5)
	c.endpoint = url
	return c
}
