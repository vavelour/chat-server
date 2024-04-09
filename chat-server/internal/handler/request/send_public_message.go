package request

import "github.com/go-playground/validator/v10"

type SendPublicMessageRequest struct {
	Sender    string `validate:"required"`
	Recipient string
	Content   string `json:"content" validate:"required"`
}

func (r *SendPublicMessageRequest) Validate(v *validator.Validate) error {
	err := v.Struct(r)
	if err != nil {
		return err
	}

	return nil
}
