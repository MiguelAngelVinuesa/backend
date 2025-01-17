// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// RoundNextRequest RoundNextRequest
//
// swagger:model RoundNextRequest
type RoundNextRequest struct {

	// Player game state JSON
	// Example: {}
	PlayerGameState *string `json:"playerGameState,omitempty"`

	// Round ID
	// Example: 7sKCVTqlnw9sRV858GFbINXotSoSCZrK
	// Required: true
	RoundID string `json:"roundId"`

	// Round state JSON
	// Example: {}
	RoundState *string `json:"roundState,omitempty"`

	// Player session ID
	// Example: bot9897cc03f5d7b43923a73bfaffc2d7dd43
	// Required: true
	SessionID string `json:"sessionId"`

	// Game session state JSON
	// Example: {}
	SessionState *string `json:"sessionState,omitempty"`

	// Spin number to return (ignored if 0)
	// Example: 1
	Spin *int32 `json:"spin,omitempty"`
}

// Validate validates this round next request
func (m *RoundNextRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRoundID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSessionID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RoundNextRequest) validateRoundID(formats strfmt.Registry) error {

	if err := validate.RequiredString("roundId", "body", m.RoundID); err != nil {
		return err
	}

	return nil
}

func (m *RoundNextRequest) validateSessionID(formats strfmt.Registry) error {

	if err := validate.RequiredString("sessionId", "body", m.SessionID); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this round next request based on context it is used
func (m *RoundNextRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RoundNextRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RoundNextRequest) UnmarshalBinary(b []byte) error {
	var res RoundNextRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
