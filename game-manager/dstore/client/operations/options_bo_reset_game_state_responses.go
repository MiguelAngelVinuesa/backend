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

// OptionsBoResetGameStateReader is a Reader for the OptionsBoResetGameState structure.
type OptionsBoResetGameStateReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *OptionsBoResetGameStateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewOptionsBoResetGameStateOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewOptionsBoResetGameStateOK creates a OptionsBoResetGameStateOK with default headers values
func NewOptionsBoResetGameStateOK() *OptionsBoResetGameStateOK {
	return &OptionsBoResetGameStateOK{}
}

/* OptionsBoResetGameStateOK describes a response with status code 200, with default header values.

200
*/
type OptionsBoResetGameStateOK struct {
	Payload *models.CORSResponse
}

func (o *OptionsBoResetGameStateOK) Error() string {
	return fmt.Sprintf("[OPTIONS /bo/reset-game-state][%d] optionsBoResetGameStateOK  %+v", 200, o.Payload)
}
func (o *OptionsBoResetGameStateOK) GetPayload() *models.CORSResponse {
	return o.Payload
}

func (o *OptionsBoResetGameStateOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.CORSResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
