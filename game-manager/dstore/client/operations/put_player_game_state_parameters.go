// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/dstore/models"
)

// NewPutPlayerGameStateParams creates a new PutPlayerGameStateParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPutPlayerGameStateParams() *PutPlayerGameStateParams {
	return &PutPlayerGameStateParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPutPlayerGameStateParamsWithTimeout creates a new PutPlayerGameStateParams object
// with the ability to set a timeout on a request.
func NewPutPlayerGameStateParamsWithTimeout(timeout time.Duration) *PutPlayerGameStateParams {
	return &PutPlayerGameStateParams{
		timeout: timeout,
	}
}

// NewPutPlayerGameStateParamsWithContext creates a new PutPlayerGameStateParams object
// with the ability to set a context for a request.
func NewPutPlayerGameStateParamsWithContext(ctx context.Context) *PutPlayerGameStateParams {
	return &PutPlayerGameStateParams{
		Context: ctx,
	}
}

// NewPutPlayerGameStateParamsWithHTTPClient creates a new PutPlayerGameStateParams object
// with the ability to set a custom HTTPClient for a request.
func NewPutPlayerGameStateParamsWithHTTPClient(client *http.Client) *PutPlayerGameStateParams {
	return &PutPlayerGameStateParams{
		HTTPClient: client,
	}
}

/* PutPlayerGameStateParams contains all the parameters to send to the API endpoint
   for the put player game state operation.

   Typically these are written to a http.Request.
*/
type PutPlayerGameStateParams struct {

	// XAPIKey.
	XAPIKey string

	// Payload.
	Payload *models.PutPlayerGameStateRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the put player game state params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutPlayerGameStateParams) WithDefaults() *PutPlayerGameStateParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the put player game state params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutPlayerGameStateParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the put player game state params
func (o *PutPlayerGameStateParams) WithTimeout(timeout time.Duration) *PutPlayerGameStateParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the put player game state params
func (o *PutPlayerGameStateParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the put player game state params
func (o *PutPlayerGameStateParams) WithContext(ctx context.Context) *PutPlayerGameStateParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the put player game state params
func (o *PutPlayerGameStateParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the put player game state params
func (o *PutPlayerGameStateParams) WithHTTPClient(client *http.Client) *PutPlayerGameStateParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the put player game state params
func (o *PutPlayerGameStateParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithXAPIKey adds the xAPIKey to the put player game state params
func (o *PutPlayerGameStateParams) WithXAPIKey(xAPIKey string) *PutPlayerGameStateParams {
	o.SetXAPIKey(xAPIKey)
	return o
}

// SetXAPIKey adds the xApiKey to the put player game state params
func (o *PutPlayerGameStateParams) SetXAPIKey(xAPIKey string) {
	o.XAPIKey = xAPIKey
}

// WithPayload adds the payload to the put player game state params
func (o *PutPlayerGameStateParams) WithPayload(payload *models.PutPlayerGameStateRequest) *PutPlayerGameStateParams {
	o.SetPayload(payload)
	return o
}

// SetPayload adds the payload to the put player game state params
func (o *PutPlayerGameStateParams) SetPayload(payload *models.PutPlayerGameStateRequest) {
	o.Payload = payload
}

// WriteToRequest writes these params to a swagger request
func (o *PutPlayerGameStateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// header param X-API-Key
	if err := r.SetHeaderParam("X-API-Key", o.XAPIKey); err != nil {
		return err
	}
	if o.Payload != nil {
		if err := r.SetBodyParam(o.Payload); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
