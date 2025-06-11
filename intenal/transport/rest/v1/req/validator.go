package req

import "github.com/go-playground/validator"

var validate = validator.New()

func (r RegisterRequest) Validate() error {
	return validate.Struct(r)
}

func (r LoginRequest) Validate() error {
	return validate.Struct(r)
}
