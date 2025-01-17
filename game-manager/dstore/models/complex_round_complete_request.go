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

// ComplexRoundCompleteRequest ComplexRoundCompleteRequest
//
// swagger:model ComplexRoundCompleteRequest
type ComplexRoundCompleteRequest struct {

	// Player game state JSON
	// Example: {}
	PlayerGameState *string `json:"playerGameState,omitempty"`

	// Initial round result JSON
	// Example: []
	Result *string `json:"result,omitempty"`

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
}

// Validate validates this complex round complete request
func (m *ComplexRoundCompleteRequest) Validate(formats strfmt.Registry) error {
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

func (m *ComplexRoundCompleteRequest) validateRoundID(formats strfmt.Registry) error {

	if err := validate.RequiredString("roundId", "body", m.RoundID); err != nil {
		return err
	}

	return nil
}

func (m *ComplexRoundCompleteRequest) validateSessionID(formats strfmt.Registry) error {

	if err := validate.RequiredString("sessionId", "body", m.SessionID); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this complex round complete request based on context it is used
func (m *ComplexRoundCompleteRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ComplexRoundCompleteRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComplexRoundCompleteRequest) UnmarshalBinary(b []byte) error {
	var res ComplexRoundCompleteRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
