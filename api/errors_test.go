package api

import (
	"net/http"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	e := &APIError{Code: 404, Message: "Not Found"}
	if e.Error() != "API error 404: Not Found" {
		t.Errorf("expected 'API error 404: Not Found', got '%s'", e.Error())
	}

	e2 := &APIError{Code: 500, Message: "Internal Error", Details: "something broke"}
	if e2.Error() != "API error 500: Internal Error (something broke)" {
		t.Errorf("expected 'API error 500: Internal Error (something broke)', got '%s'", e2.Error())
	}
}

func TestParseError(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 404, APIError{Code: 404, Message: "not found"})
	}))
	defer server.Close()

	resp, err := client.R().Get("/")
	if err != nil {
		t.Fatal(err)
	}

	err = parseError(resp)
	if !IsNotFound(err) {
		t.Errorf("expected IsNotFound to be true, got false: %v", err)
	}
}

func TestParseErrorNonJSON(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("plain text error"))
	}))
	defer server.Close()

	resp, err := client.R().Get("/")
	if err != nil {
		t.Fatal(err)
	}

	err = parseError(resp)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Code != 500 {
		t.Errorf("expected code 500, got %d", apiErr.Code)
	}
}

func TestParseErrorEmptyMessage(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: ""})
	}))
	defer server.Close()

	resp, err := client.R().Get("/")
	if err != nil {
		t.Fatal(err)
	}

	err = parseError(resp)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Message == "" {
		t.Error("expected non-empty message from Status() fallback")
	}
}

func TestParseErrorCodeZero(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 503, map[string]string{"message": "service unavailable"})
	}))
	defer server.Close()

	resp, err := client.R().Get("/")
	if err != nil {
		t.Fatal(err)
	}

	err = parseError(resp)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Code != 503 {
		t.Errorf("expected code 503, got %d", apiErr.Code)
	}
}

func TestParseErrorNilResponse(t *testing.T) {
	err := parseError(nil)
	if err == nil {
		t.Error("expected error for nil response")
	}
}

func TestIsNotFound(t *testing.T) {
	if IsNotFound(nil) {
		t.Error("nil should not be IsNotFound")
	}
	if IsNotFound(&APIError{Code: 404}) != true {
		t.Error("404 should be IsNotFound")
	}
	if IsNotFound(&APIError{Code: 500}) != false {
		t.Error("500 should not be IsNotFound")
	}
}

func TestIsUnauthorized(t *testing.T) {
	if IsUnauthorized(nil) {
		t.Error("nil should not be IsUnauthorized")
	}
	if IsUnauthorized(&APIError{Code: 401}) != true {
		t.Error("401 should be IsUnauthorized")
	}
	if IsUnauthorized(&APIError{Code: 403}) != false {
		t.Error("403 should not be IsUnauthorized")
	}
}

func TestIsConflict(t *testing.T) {
	if IsConflict(nil) {
		t.Error("nil should not be IsConflict")
	}
	if IsConflict(&APIError{Code: 409}) != true {
		t.Error("409 should be IsConflict")
	}
	if IsConflict(&APIError{Code: 400}) != false {
		t.Error("400 should not be IsConflict")
	}
}

func TestParseErrorFromBody(t *testing.T) {
	err := ParseErrorFromBody(500, []byte(`{"code":500,"message":"boom"}`))
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Code != 500 || apiErr.Message != "boom" {
		t.Errorf("unexpected error: %+v", apiErr)
	}
}

func TestParseErrorFromBodyNonJSON(t *testing.T) {
	err := ParseErrorFromBody(503, []byte("not json"))
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Code != 503 {
		t.Errorf("expected code 503, got %d", apiErr.Code)
	}
}

func TestParseErrorFromBodyDefaults(t *testing.T) {
	err := ParseErrorFromBody(503, []byte(`{}`))
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Code != 503 {
		t.Errorf("expected code 503, got %d", apiErr.Code)
	}
}
