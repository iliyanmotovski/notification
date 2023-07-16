package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func ShouldBindJSON(ctx *gin.Context, form interface{}) error {
	if ctx.Request.Body == nil {
		return fmt.Errorf("body cannot be nil")
	}

	if err := ctx.ShouldBindBodyWith(form, binding.JSON); err != nil {
		res := handleErrors(err)
		if res != nil {
			return res
		}

		return err
	}

	return nil
}

type Error struct {
	GlobalError string            `json:"global_error,omitempty"`
	FieldsError map[string]string `json:"fields_error,omitempty"`
}

type FieldErrors map[string]string

func (fe FieldErrors) Error() string {
	var result string
	for _, val := range fe {
		result += val + "\n\r"
	}

	return result
}

func HandleError(ctx *gin.Context, err error) bool {
	errType, ok := err.(FieldErrors)
	if ok && errType != nil {
		errorResponseFields(ctx, errType)
		return true
	}

	if err != nil {
		errorResponseGlobal(ctx, err)
		return true
	}

	return false
}

func errorResponseFields(ctx *gin.Context, fieldsError FieldErrors) {
	ctx.JSON(http.StatusBadRequest, &Error{
		FieldsError: fieldsError,
	})
}

func errorResponseGlobal(ctx *gin.Context, globalError interface{}) {
	result := &Error{}

	err, ok := globalError.(error)
	if ok {
		result.GlobalError = err.Error()
	} else {
		result.GlobalError = globalError.(string)
	}

	ctx.JSON(http.StatusBadRequest, result)
}

func handleErrors(formErrors interface{}) error {
	errors, ok := formErrors.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	var fieldErrors FieldErrors = make(map[string]string)
	for _, err := range errors {
		fieldErrors[err.Field()] = err.Error()
	}

	return fieldErrors
}
