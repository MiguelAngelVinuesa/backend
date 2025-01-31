// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/dstore/models"
)

// GetPlayerGameStateReader is a Reader for the GetPlayerGameState structure.
type GetPlayerGameStateReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPlayerGameStateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPlayerGameStateOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetPlayerGameStateBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetPlayerGameStateInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetPlayerGameStateOK creates a GetPlayerGameStateOK with default headers values
func NewGetPlayerGameStateOK() *GetPlayerGameStateOK {
	return &GetPlayerGameStateOK{}
}

/* GetPlayerGameStateOK describes a response with status code 200, with default header values.

200
*/
type GetPlayerGameStateOK struct {
	Payload *models.PlayerGameStateResponse
}

func (o *GetPlayerGameStateOK) Error() string {
	return fmt.Sprintf("[GET /player-game-state][%d] getPlayerGameStateOK  %+v", 200, o.Payload)
}
func (o *GetPlayerGameStateOK) GetPayload() *models.PlayerGameStateResponse {
	return o.Payload
}

func (o *GetPlayerGameStateOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PlayerGameStateResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPlayerGameStateBadRequest creates a GetPlayerGameStateBadRequest with default headers values
func NewGetPlayerGameStateBadRequest() *GetPlayerGameStateBadRequest {
	return &GetPlayerGameStateBadRequest{}
}

/* GetPlayerGameStateBadRequest describes a response with status code 400, with default header values.

400
*/
type GetPlayerGameStateBadRequest struct {
	Payload *models.ErrorResponse
}

func (o *GetPlayerGameStateBadRequest) Error() string {
	return fmt.Sprintf("[GET /player-game-state][%d] getPlayerGameStateBadRequest  %+v", 400, o.Payload)
}
func (o *GetPlayerGameStateBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetPlayerGameStateBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPlayerGameStateInternalServerError creates a GetPlayerGameStateInternalServerError with default headers values
func NewGetPlayerGameStateInternalServerError() *GetPlayerGameStateInternalServerError {
	return &GetPlayerGameStateInternalServerError{}
}

/* GetPlayerGameStateInternalServerError describes a response with status code 500, with default header values.

500
*/
type GetPlayerGameStateInternalServerError struct {
	Payload *models.ErrorResponse
}

func (o *GetPlayerGameStateInternalServerError) Error() string {
	return fmt.Sprintf("[GET /player-game-state][%d] getPlayerGameStateInternalServerError  %+v", 500, o.Payload)
}
func (o *GetPlayerGameStateInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetPlayerGameStateInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
