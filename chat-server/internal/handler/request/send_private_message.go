package request

import "github.com/go-playground/validator/v10"

type SendPrivateMessageRequest struct {
	Sender    string `validate:"required"`
	Recipient string `validate:"required"`
	Content   string `json:"content" validate:"required"`
}

func (r *SendPrivateMessageRequest) Validate(v *validator.Validate) error {
	err := v.Struct(r)
	if err != nil {
		return err
	}

	return nil
}
