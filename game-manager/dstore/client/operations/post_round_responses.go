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

// PostRoundReader is a Reader for the PostRound structure.
type PostRoundReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostRoundReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPostRoundOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPostRoundBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewPostRoundInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewPostRoundOK creates a PostRoundOK with default headers values
func NewPostRoundOK() *PostRoundOK {
	return &PostRoundOK{}
}

/* PostRoundOK describes a response with status code 200, with default header values.

200
*/
type PostRoundOK struct {
	Payload *models.RoundResponse
}

func (o *PostRoundOK) Error() string {
	return fmt.Sprintf("[POST /round][%d] postRoundOK  %+v", 200, o.Payload)
}
func (o *PostRoundOK) GetPayload() *models.RoundResponse {
	return o.Payload
}

func (o *PostRoundOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RoundResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostRoundBadRequest creates a PostRoundBadRequest with default headers values
func NewPostRoundBadRequest() *PostRoundBadRequest {
	return &PostRoundBadRequest{}
}

/* PostRoundBadRequest describes a response with status code 400, with default header values.

400
*/
type PostRoundBadRequest struct {
	Payload *models.ErrorResponse
}

func (o *PostRoundBadRequest) Error() string {
	return fmt.Sprintf("[POST /round][%d] postRoundBadRequest  %+v", 400, o.Payload)
}
func (o *PostRoundBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *PostRoundBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostRoundInternalServerError creates a PostRoundInternalServerError with default headers values
func NewPostRoundInternalServerError() *PostRoundInternalServerError {
	return &PostRoundInternalServerError{}
}

/* PostRoundInternalServerError describes a response with status code 500, with default header values.

500
*/
type PostRoundInternalServerError struct {
	Payload *models.ErrorResponse
}

func (o *PostRoundInternalServerError) Error() string {
	return fmt.Sprintf("[POST /round][%d] postRoundInternalServerError  %+v", 500, o.Payload)
}
func (o *PostRoundInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *PostRoundInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
