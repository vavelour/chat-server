package request

import "github.com/go-playground/validator/v10"

type ShowPublicMessageRequest struct {
	Limit  int `json:"limit" validate:"min=1"`
	Offset int `json:"offset" validate:"min=0"`
}

func (r *ShowPublicMessageRequest) Validate(v *validator.Validate) error {
	err := v.Struct(r)
	if err != nil {
		return err
	}

	return nil
}
