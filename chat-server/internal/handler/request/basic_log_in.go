package request

import "github.com/go-playground/validator/v10"

type BasicAuthLogInRequest struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

func (r *BasicAuthLogInRequest) Validate(v *validator.Validate) error {
	err := v.Struct(r)
	if err != nil {
		return err
	}

	return nil
}
