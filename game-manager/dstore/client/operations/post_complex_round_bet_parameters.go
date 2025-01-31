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

// NewPostComplexRoundBetParams creates a new PostComplexRoundBetParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPostComplexRoundBetParams() *PostComplexRoundBetParams {
	return &PostComplexRoundBetParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPostComplexRoundBetParamsWithTimeout creates a new PostComplexRoundBetParams object
// with the ability to set a timeout on a request.
func NewPostComplexRoundBetParamsWithTimeout(timeout time.Duration) *PostComplexRoundBetParams {
	return &PostComplexRoundBetParams{
		timeout: timeout,
	}
}

// NewPostComplexRoundBetParamsWithContext creates a new PostComplexRoundBetParams object
// with the ability to set a context for a request.
func NewPostComplexRoundBetParamsWithContext(ctx context.Context) *PostComplexRoundBetParams {
	return &PostComplexRoundBetParams{
		Context: ctx,
	}
}

// NewPostComplexRoundBetParamsWithHTTPClient creates a new PostComplexRoundBetParams object
// with the ability to set a custom HTTPClient for a request.
func NewPostComplexRoundBetParamsWithHTTPClient(client *http.Client) *PostComplexRoundBetParams {
	return &PostComplexRoundBetParams{
		HTTPClient: client,
	}
}

/* PostComplexRoundBetParams contains all the parameters to send to the API endpoint
   for the post complex round bet operation.

   Typically these are written to a http.Request.
*/
type PostComplexRoundBetParams struct {

	// XAPIKey.
	XAPIKey string

	// Payload.
	Payload *models.ComplexRoundBetRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the post complex round bet params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostComplexRoundBetParams) WithDefaults() *PostComplexRoundBetParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the post complex round bet params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostComplexRoundBetParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the post complex round bet params
func (o *PostComplexRoundBetParams) WithTimeout(timeout time.Duration) *PostComplexRoundBetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post complex round bet params
func (o *PostComplexRoundBetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post complex round bet params
func (o *PostComplexRoundBetParams) WithContext(ctx context.Context) *PostComplexRoundBetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post complex round bet params
func (o *PostComplexRoundBetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post complex round bet params
func (o *PostComplexRoundBetParams) WithHTTPClient(client *http.Client) *PostComplexRoundBetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post complex round bet params
func (o *PostComplexRoundBetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithXAPIKey adds the xAPIKey to the post complex round bet params
func (o *PostComplexRoundBetParams) WithXAPIKey(xAPIKey string) *PostComplexRoundBetParams {
	o.SetXAPIKey(xAPIKey)
	return o
}

// SetXAPIKey adds the xApiKey to the post complex round bet params
func (o *PostComplexRoundBetParams) SetXAPIKey(xAPIKey string) {
	o.XAPIKey = xAPIKey
}

// WithPayload adds the payload to the post complex round bet params
func (o *PostComplexRoundBetParams) WithPayload(payload *models.ComplexRoundBetRequest) *PostComplexRoundBetParams {
	o.SetPayload(payload)
	return o
}

// SetPayload adds the payload to the post complex round bet params
func (o *PostComplexRoundBetParams) SetPayload(payload *models.ComplexRoundBetRequest) {
	o.Payload = payload
}

// WriteToRequest writes these params to a swagger request
func (o *PostComplexRoundBetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
