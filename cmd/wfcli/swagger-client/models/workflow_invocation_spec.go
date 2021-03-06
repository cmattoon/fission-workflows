package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// WorkflowInvocationSpec Workflow Invocation Model
// swagger:model WorkflowInvocationSpec
type WorkflowInvocationSpec struct {

	// inputs
	Inputs map[string]TypedValue `json:"inputs,omitempty"`

	// workflow Id
	WorkflowID string `json:"workflowId,omitempty"`
}

// Validate validates this workflow invocation spec
func (m *WorkflowInvocationSpec) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateInputs(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *WorkflowInvocationSpec) validateInputs(formats strfmt.Registry) error {

	if swag.IsZero(m.Inputs) { // not required
		return nil
	}

	if err := validate.Required("inputs", "body", m.Inputs); err != nil {
		return err
	}

	return nil
}
