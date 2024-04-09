package request

import "github.com/go-playground/validator/v10"

type ViewUserListRequest struct {
	Username string `validate:"required"`
}

func (r *ViewUserListRequest) Validate(v *validator.Validate) error {
	err := v.Struct(r)
	if err != nil {
		return err
	}

	return nil
}
