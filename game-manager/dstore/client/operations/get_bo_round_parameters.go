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
)

// NewGetBoRoundParams creates a new GetBoRoundParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetBoRoundParams() *GetBoRoundParams {
	return &GetBoRoundParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetBoRoundParamsWithTimeout creates a new GetBoRoundParams object
// with the ability to set a timeout on a request.
func NewGetBoRoundParamsWithTimeout(timeout time.Duration) *GetBoRoundParams {
	return &GetBoRoundParams{
		timeout: timeout,
	}
}

// NewGetBoRoundParamsWithContext creates a new GetBoRoundParams object
// with the ability to set a context for a request.
func NewGetBoRoundParamsWithContext(ctx context.Context) *GetBoRoundParams {
	return &GetBoRoundParams{
		Context: ctx,
	}
}

// NewGetBoRoundParamsWithHTTPClient creates a new GetBoRoundParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetBoRoundParamsWithHTTPClient(client *http.Client) *GetBoRoundParams {
	return &GetBoRoundParams{
		HTTPClient: client,
	}
}

/* GetBoRoundParams contains all the parameters to send to the API endpoint
   for the get bo round operation.

   Typically these are written to a http.Request.
*/
type GetBoRoundParams struct {

	// XAPIKey.
	XAPIKey string

	// ID.
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get bo round params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetBoRoundParams) WithDefaults() *GetBoRoundParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get bo round params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetBoRoundParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get bo round params
func (o *GetBoRoundParams) WithTimeout(timeout time.Duration) *GetBoRoundParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get bo round params
func (o *GetBoRoundParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get bo round params
func (o *GetBoRoundParams) WithContext(ctx context.Context) *GetBoRoundParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get bo round params
func (o *GetBoRoundParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get bo round params
func (o *GetBoRoundParams) WithHTTPClient(client *http.Client) *GetBoRoundParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get bo round params
func (o *GetBoRoundParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithXAPIKey adds the xAPIKey to the get bo round params
func (o *GetBoRoundParams) WithXAPIKey(xAPIKey string) *GetBoRoundParams {
	o.SetXAPIKey(xAPIKey)
	return o
}

// SetXAPIKey adds the xApiKey to the get bo round params
func (o *GetBoRoundParams) SetXAPIKey(xAPIKey string) {
	o.XAPIKey = xAPIKey
}

// WithID adds the id to the get bo round params
func (o *GetBoRoundParams) WithID(id string) *GetBoRoundParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get bo round params
func (o *GetBoRoundParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetBoRoundParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// header param X-API-Key
	if err := r.SetHeaderParam("X-API-Key", o.XAPIKey); err != nil {
		return err
	}

	// query param id
	qrID := o.ID
	qID := qrID
	if qID != "" {

		if err := r.SetQueryParam("id", qID); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
