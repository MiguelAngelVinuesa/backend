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

// PutPlayerGlobalStateResponse PutPlayerGlobalStateResponse
//
// swagger:model PutPlayerGlobalStateResponse
type PutPlayerGlobalStateResponse struct {

	// Boolean determining whether the request was successful
	// Example: true
	// Required: true
	Success bool `json:"success"`
}

// Validate validates this put player global state response
func (m *PutPlayerGlobalStateResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSuccess(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PutPlayerGlobalStateResponse) validateSuccess(formats strfmt.Registry) error {

	if err := validate.Required("success", "body", bool(m.Success)); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this put player global state response based on context it is used
func (m *PutPlayerGlobalStateResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PutPlayerGlobalStateResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PutPlayerGlobalStateResponse) UnmarshalBinary(b []byte) error {
	var res PutPlayerGlobalStateResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
