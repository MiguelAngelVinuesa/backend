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

// NewPostRoundParams creates a new PostRoundParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPostRoundParams() *PostRoundParams {
	return &PostRoundParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPostRoundParamsWithTimeout creates a new PostRoundParams object
// with the ability to set a timeout on a request.
func NewPostRoundParamsWithTimeout(timeout time.Duration) *PostRoundParams {
	return &PostRoundParams{
		timeout: timeout,
	}
}

// NewPostRoundParamsWithContext creates a new PostRoundParams object
// with the ability to set a context for a request.
func NewPostRoundParamsWithContext(ctx context.Context) *PostRoundParams {
	return &PostRoundParams{
		Context: ctx,
	}
}

// NewPostRoundParamsWithHTTPClient creates a new PostRoundParams object
// with the ability to set a custom HTTPClient for a request.
func NewPostRoundParamsWithHTTPClient(client *http.Client) *PostRoundParams {
	return &PostRoundParams{
		HTTPClient: client,
	}
}

/* PostRoundParams contains all the parameters to send to the API endpoint
   for the post round operation.

   Typically these are written to a http.Request.
*/
type PostRoundParams struct {

	// XAPIKey.
	XAPIKey string

	// Payload.
	Payload *models.RoundRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the post round params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostRoundParams) WithDefaults() *PostRoundParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the post round params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostRoundParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the post round params
func (o *PostRoundParams) WithTimeout(timeout time.Duration) *PostRoundParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post round params
func (o *PostRoundParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post round params
func (o *PostRoundParams) WithContext(ctx context.Context) *PostRoundParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post round params
func (o *PostRoundParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post round params
func (o *PostRoundParams) WithHTTPClient(client *http.Client) *PostRoundParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post round params
func (o *PostRoundParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithXAPIKey adds the xAPIKey to the post round params
func (o *PostRoundParams) WithXAPIKey(xAPIKey string) *PostRoundParams {
	o.SetXAPIKey(xAPIKey)
	return o
}

// SetXAPIKey adds the xApiKey to the post round params
func (o *PostRoundParams) SetXAPIKey(xAPIKey string) {
	o.XAPIKey = xAPIKey
}

// WithPayload adds the payload to the post round params
func (o *PostRoundParams) WithPayload(payload *models.RoundRequest) *PostRoundParams {
	o.SetPayload(payload)
	return o
}

// SetPayload adds the payload to the post round params
func (o *PostRoundParams) SetPayload(payload *models.RoundRequest) {
	o.Payload = payload
}

// WriteToRequest writes these params to a swagger request
func (o *PostRoundParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
