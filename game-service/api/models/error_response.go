// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// ErrorResponse Error response.
//
// swagger:model ErrorResponse
type ErrorResponse struct {

	// An error code indicating why the request failed.
	// Example: 101
	// Required: true
	ErrorCode int64 `json:"errorCode"`

	// The error level for the error code (R=retry; F=fatal).
	// Example: F
	// Required: true
	ErrorLevel string `json:"errorLevel"`

	// An error message describing what went wrong.
	// Example: Bad request
	// Required: true
	Message string `json:"message"`

	// Indicates if the request was successful.
	// Example: false
	// Required: true
	Success bool `json:"success"`
}
