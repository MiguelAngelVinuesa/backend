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

// NewPostComplexRoundCompleteParams creates a new PostComplexRoundCompleteParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPostComplexRoundCompleteParams() *PostComplexRoundCompleteParams {
	return &PostComplexRoundCompleteParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPostComplexRoundCompleteParamsWithTimeout creates a new PostComplexRoundCompleteParams object
// with the ability to set a timeout on a request.
func NewPostComplexRoundCompleteParamsWithTimeout(timeout time.Duration) *PostComplexRoundCompleteParams {
	return &PostComplexRoundCompleteParams{
		timeout: timeout,
	}
}

// NewPostComplexRoundCompleteParamsWithContext creates a new PostComplexRoundCompleteParams object
// with the ability to set a context for a request.
func NewPostComplexRoundCompleteParamsWithContext(ctx context.Context) *PostComplexRoundCompleteParams {
	return &PostComplexRoundCompleteParams{
		Context: ctx,
	}
}

// NewPostComplexRoundCompleteParamsWithHTTPClient creates a new PostComplexRoundCompleteParams object
// with the ability to set a custom HTTPClient for a request.
func NewPostComplexRoundCompleteParamsWithHTTPClient(client *http.Client) *PostComplexRoundCompleteParams {
	return &PostComplexRoundCompleteParams{
		HTTPClient: client,
	}
}

/* PostComplexRoundCompleteParams contains all the parameters to send to the API endpoint
   for the post complex round complete operation.

   Typically these are written to a http.Request.
*/
type PostComplexRoundCompleteParams struct {

	// XAPIKey.
	XAPIKey string

	// Payload.
	Payload *models.ComplexRoundCompleteRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the post complex round complete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostComplexRoundCompleteParams) WithDefaults() *PostComplexRoundCompleteParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the post complex round complete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostComplexRoundCompleteParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the post complex round complete params
func (o *PostComplexRoundCompleteParams) WithTimeout(timeout time.Duration) *PostComplexRoundCompleteParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post complex round complete params
func (o *PostComplexRoundCompleteParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post complex round complete params
func (o *PostComplexRoundCompleteParams) WithContext(ctx context.Context) *PostComplexRoundCompleteParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post complex round complete params
func (o *PostComplexRoundCompleteParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post complex round complete params
func (o *PostComplexRoundCompleteParams) WithHTTPClient(client *http.Client) *PostComplexRoundCompleteParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post complex round complete params
func (o *PostComplexRoundCompleteParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithXAPIKey adds the xAPIKey to the post complex round complete params
func (o *PostComplexRoundCompleteParams) WithXAPIKey(xAPIKey string) *PostComplexRoundCompleteParams {
	o.SetXAPIKey(xAPIKey)
	return o
}

// SetXAPIKey adds the xApiKey to the post complex round complete params
func (o *PostComplexRoundCompleteParams) SetXAPIKey(xAPIKey string) {
	o.XAPIKey = xAPIKey
}

// WithPayload adds the payload to the post complex round complete params
func (o *PostComplexRoundCompleteParams) WithPayload(payload *models.ComplexRoundCompleteRequest) *PostComplexRoundCompleteParams {
	o.SetPayload(payload)
	return o
}

// SetPayload adds the payload to the post complex round complete params
func (o *PostComplexRoundCompleteParams) SetPayload(payload *models.ComplexRoundCompleteRequest) {
	o.Payload = payload
}

// WriteToRequest writes these params to a swagger request
func (o *PostComplexRoundCompleteParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
