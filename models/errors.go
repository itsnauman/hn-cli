package models

import (
	"fmt"
	"strings"
)

// ErrorOutput is the structured error type rendered to stdout.
type ErrorOutput struct {
	Error string `toon:"error" json:"error"`
	Code  int    `toon:"code" json:"code"`
	Hint  string `toon:"hint" json:"hint"`
}

func NewNotFoundError(resource, id string) *ErrorOutput {
	return &ErrorOutput{
		Error: fmt.Sprintf("%s not found: %s", resource, id),
		Code:  404,
		Hint:  fmt.Sprintf("check that the %s ID is correct", resource),
	}
}

func NewAPIError(err error) *ErrorOutput {
	return &ErrorOutput{
		Error: fmt.Sprintf("api error: %s", err),
		Code:  500,
		Hint:  "the Hacker News API may be temporarily unavailable — retry in a moment",
	}
}

// NewErrorFromFetch returns a not-found error if the error message contains "not found",
// otherwise returns a generic API error.
func NewErrorFromFetch(resource, id string, err error) *ErrorOutput {
	if strings.Contains(err.Error(), "not found") {
		return NewNotFoundError(resource, id)
	}
	return NewAPIError(err)
}

func NewValidationError(msg, hint string) *ErrorOutput {
	return &ErrorOutput{
		Error: msg,
		Code:  400,
		Hint:  hint,
	}
}
