// Copyright © 2019 @ken-aio <suguru.akiho@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"
)

var validate = validator.New()

func validateParams(p interface{}) []validator.FieldError {
	errs := validate.Struct(p)
	return extractValidationErrors(errs)
}

func extractValidationErrors(err error) []validator.FieldError {
	fieldErrors := make([]validator.FieldError, 0)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldErrors = append(fieldErrors, err.(validator.FieldError))
		}
		return fieldErrors
	}

	return nil
}

func validationErrorToText(e validator.FieldError, text string) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", text)
	case "max":
		return fmt.Sprintf("%s cannot be greater than %s", text, e.Param())
	case "min":
		return fmt.Sprintf("%s must be greater than %s", text, e.Param())
	}
	return fmt.Sprintf("%s is not valid %s", text, e.Value())
}
