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

// NewPostBoResetGameStateParams creates a new PostBoResetGameStateParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPostBoResetGameStateParams() *PostBoResetGameStateParams {
	return &PostBoResetGameStateParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPostBoResetGameStateParamsWithTimeout creates a new PostBoResetGameStateParams object
// with the ability to set a timeout on a request.
func NewPostBoResetGameStateParamsWithTimeout(timeout time.Duration) *PostBoResetGameStateParams {
	return &PostBoResetGameStateParams{
		timeout: timeout,
	}
}

// NewPostBoResetGameStateParamsWithContext creates a new PostBoResetGameStateParams object
// with the ability to set a context for a request.
func NewPostBoResetGameStateParamsWithContext(ctx context.Context) *PostBoResetGameStateParams {
	return &PostBoResetGameStateParams{
		Context: ctx,
	}
}

// NewPostBoResetGameStateParamsWithHTTPClient creates a new PostBoResetGameStateParams object
// with the ability to set a custom HTTPClient for a request.
func NewPostBoResetGameStateParamsWithHTTPClient(client *http.Client) *PostBoResetGameStateParams {
	return &PostBoResetGameStateParams{
		HTTPClient: client,
	}
}

/* PostBoResetGameStateParams contains all the parameters to send to the API endpoint
   for the post bo reset game state operation.

   Typically these are written to a http.Request.
*/
type PostBoResetGameStateParams struct {

	// XAPIKey.
	XAPIKey string

	// CasinoID.
	CasinoID string

	// PlayerID.
	PlayerID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the post bo reset game state params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostBoResetGameStateParams) WithDefaults() *PostBoResetGameStateParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the post bo reset game state params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostBoResetGameStateParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the post bo reset game state params
func (o *PostBoResetGameStateParams) WithTimeout(timeout time.Duration) *PostBoResetGameStateParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post bo reset game state params
func (o *PostBoResetGameStateParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post bo reset game state params
func (o *PostBoResetGameStateParams) WithContext(ctx context.Context) *PostBoResetGameStateParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post bo reset game state params
func (o *PostBoResetGameStateParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post bo reset game state params
func (o *PostBoResetGameStateParams) WithHTTPClient(client *http.Client) *PostBoResetGameStateParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post bo reset game state params
func (o *PostBoResetGameStateParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithXAPIKey adds the xAPIKey to the post bo reset game state params
func (o *PostBoResetGameStateParams) WithXAPIKey(xAPIKey string) *PostBoResetGameStateParams {
	o.SetXAPIKey(xAPIKey)
	return o
}

// SetXAPIKey adds the xApiKey to the post bo reset game state params
func (o *PostBoResetGameStateParams) SetXAPIKey(xAPIKey string) {
	o.XAPIKey = xAPIKey
}

// WithCasinoID adds the casinoID to the post bo reset game state params
func (o *PostBoResetGameStateParams) WithCasinoID(casinoID string) *PostBoResetGameStateParams {
	o.SetCasinoID(casinoID)
	return o
}

// SetCasinoID adds the casinoId to the post bo reset game state params
func (o *PostBoResetGameStateParams) SetCasinoID(casinoID string) {
	o.CasinoID = casinoID
}

// WithPlayerID adds the playerID to the post bo reset game state params
func (o *PostBoResetGameStateParams) WithPlayerID(playerID string) *PostBoResetGameStateParams {
	o.SetPlayerID(playerID)
	return o
}

// SetPlayerID adds the playerId to the post bo reset game state params
func (o *PostBoResetGameStateParams) SetPlayerID(playerID string) {
	o.PlayerID = playerID
}

// WriteToRequest writes these params to a swagger request
func (o *PostBoResetGameStateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// header param X-API-Key
	if err := r.SetHeaderParam("X-API-Key", o.XAPIKey); err != nil {
		return err
	}

	// query param casinoId
	qrCasinoID := o.CasinoID
	qCasinoID := qrCasinoID
	if qCasinoID != "" {

		if err := r.SetQueryParam("casinoId", qCasinoID); err != nil {
			return err
		}
	}

	// query param playerId
	qrPlayerID := o.PlayerID
	qPlayerID := qrPlayerID
	if qPlayerID != "" {

		if err := r.SetQueryParam("playerId", qPlayerID); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
