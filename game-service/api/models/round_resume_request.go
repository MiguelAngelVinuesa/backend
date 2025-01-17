// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// RoundResumeRequest Resume round after user choice
//
// Endpoint to resume the spins after the player made a selection.
//
// swagger:model RoundResumeRequest
type RoundResumeRequest struct {

	// i18n
	I18n *PrefetchI18n `json:"i18n,omitempty"`

	// Optional object with player choice(s).
	// Example: {"stickySymbol":5,"wing":"north"}
	PlayerChoice interface{} `json:"playerChoice,omitempty"`

	// First round identification.
	// Example: 5418324dc7884ad7b7d6e0fff31e4d1a
	// Required: true
	RoundID string `json:"roundId"`

	// Player session ID.
	// Example: bot9897cc03f5d7b43923a73bfaffc2d7dd43
	// Required: true
	SessionID string `json:"sessionId"`
}
