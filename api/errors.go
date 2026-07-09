package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("API error %d: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

func ParseErrorFromBody(statusCode int, body []byte) error {
	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return &APIError{
			Code:    statusCode,
			Message: fmt.Sprintf("HTTP %d", statusCode),
		}
	}
	if apiErr.Code == 0 {
		apiErr.Code = statusCode
	}
	return &apiErr
}

func parseError(resp *resty.Response) error {
	if resp == nil {
		return fmt.Errorf("nil response")
	}
	var apiErr APIError
	if err := json.Unmarshal(resp.Body(), &apiErr); err != nil {
		return &APIError{
			Code:    resp.StatusCode(),
			Message: fmt.Sprintf("HTTP %d: %s", resp.StatusCode(), resp.Status()),
		}
	}
	if apiErr.Code == 0 {
		apiErr.Code = resp.StatusCode()
	}
	if apiErr.Message == "" {
		apiErr.Message = resp.Status()
	}
	return &apiErr
}

func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == 404
	}
	return false
}

func IsUnauthorized(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == 401
	}
	return false
}

func IsConflict(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == 409
	}
	return false
}
