// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// RoundNextRequest Next result request.
//
// Retrieve the 'next' game result from a stored slot-machine round.
//
// swagger:model RoundNextRequest
type RoundNextRequest struct {

	// i18n
	I18n *PrefetchI18n `json:"i18n,omitempty"`

	// Round identification.
	// Example: 5418324dc7884ad7b7d6e0fff31e4d1a
	// Required: true
	RoundID string `json:"roundId"`

	// Player session ID
	// Example: bot9697cc03f5d7b43923a73bfaffc2d7dd43
	// Required: true
	SessionID string `json:"sessionId"`

	// Spin result sequence number.
	// Example: 2
	// Required: true
	SpinSeq int64 `json:"spinSeq"`
}
