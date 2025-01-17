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

// NewPostComplexRoundWinParams creates a new PostComplexRoundWinParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPostComplexRoundWinParams() *PostComplexRoundWinParams {
	return &PostComplexRoundWinParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPostComplexRoundWinParamsWithTimeout creates a new PostComplexRoundWinParams object
// with the ability to set a timeout on a request.
func NewPostComplexRoundWinParamsWithTimeout(timeout time.Duration) *PostComplexRoundWinParams {
	return &PostComplexRoundWinParams{
		timeout: timeout,
	}
}

// NewPostComplexRoundWinParamsWithContext creates a new PostComplexRoundWinParams object
// with the ability to set a context for a request.
func NewPostComplexRoundWinParamsWithContext(ctx context.Context) *PostComplexRoundWinParams {
	return &PostComplexRoundWinParams{
		Context: ctx,
	}
}

// NewPostComplexRoundWinParamsWithHTTPClient creates a new PostComplexRoundWinParams object
// with the ability to set a custom HTTPClient for a request.
func NewPostComplexRoundWinParamsWithHTTPClient(client *http.Client) *PostComplexRoundWinParams {
	return &PostComplexRoundWinParams{
		HTTPClient: client,
	}
}

/* PostComplexRoundWinParams contains all the parameters to send to the API endpoint
   for the post complex round win operation.

   Typically these are written to a http.Request.
*/
type PostComplexRoundWinParams struct {

	// XAPIKey.
	XAPIKey string

	// Payload.
	Payload *models.ComplexRoundWinRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the post complex round win params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostComplexRoundWinParams) WithDefaults() *PostComplexRoundWinParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the post complex round win params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostComplexRoundWinParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the post complex round win params
func (o *PostComplexRoundWinParams) WithTimeout(timeout time.Duration) *PostComplexRoundWinParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post complex round win params
func (o *PostComplexRoundWinParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post complex round win params
func (o *PostComplexRoundWinParams) WithContext(ctx context.Context) *PostComplexRoundWinParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post complex round win params
func (o *PostComplexRoundWinParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post complex round win params
func (o *PostComplexRoundWinParams) WithHTTPClient(client *http.Client) *PostComplexRoundWinParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post complex round win params
func (o *PostComplexRoundWinParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithXAPIKey adds the xAPIKey to the post complex round win params
func (o *PostComplexRoundWinParams) WithXAPIKey(xAPIKey string) *PostComplexRoundWinParams {
	o.SetXAPIKey(xAPIKey)
	return o
}

// SetXAPIKey adds the xApiKey to the post complex round win params
func (o *PostComplexRoundWinParams) SetXAPIKey(xAPIKey string) {
	o.XAPIKey = xAPIKey
}

// WithPayload adds the payload to the post complex round win params
func (o *PostComplexRoundWinParams) WithPayload(payload *models.ComplexRoundWinRequest) *PostComplexRoundWinParams {
	o.SetPayload(payload)
	return o
}

// SetPayload adds the payload to the post complex round win params
func (o *PostComplexRoundWinParams) SetPayload(payload *models.ComplexRoundWinRequest) {
	o.Payload = payload
}

// WriteToRequest writes these params to a swagger request
func (o *PostComplexRoundWinParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
