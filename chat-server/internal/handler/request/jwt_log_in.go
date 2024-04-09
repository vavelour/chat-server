package request

import "github.com/go-playground/validator/v10"

type JWTLogInRequest struct {
	Token string `validate:"required"`
}

func (r *JWTLogInRequest) Validate(v *validator.Validate) error {
	err := v.Struct(r)
	if err != nil {
		return err
	}

	return nil
}
