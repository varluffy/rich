/**
 * @Time: 2021/2/26 2:41 下午
 * @Author: varluffy
 * @Description: validate
 */

package ginwrap

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/varluffy/ginx/errcode"
	"strings"
)

var _ error = (*ValidateError)(nil)

type ValidateError struct {
	Key     string
	Message string
}

type ValidateErrors []*ValidateError

func (v *ValidateError) Error() string {
	return v.Message
}

func (v ValidateErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidateErrors) Errors() []string {
	errs := make([]string, len(v))
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

func BindAndValid(c *gin.Context, v interface{}) error {
	if err := c.ShouldBind(v); err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			trans := c.Value("trans")
			utrans, _ := trans.(ut.Translator)
			errs := err.(validator.ValidationErrors)
			return errcode.BadRequest(400, translateErrors(utrans, errs))
		case *json.UnmarshalTypeError:
			unmarshalTypeError := err.(*json.UnmarshalTypeError)
			message := fmt.Errorf("%s 类型错误，期望类型 %s", unmarshalTypeError.Field, unmarshalTypeError.Type.String()).Error()
			return errcode.BadRequest(400, message)
		default:
			return errcode.InternalServer(500, "unknown, invalid params")
		}
	}
	return nil
}

func translateErrors(trans ut.Translator, errs validator.ValidationErrors) string {
	errList := make([]string, 0)
	for _, e := range errs {
		msg := e.Translate(trans)
		errList = append(errList, msg)
	}
	return strings.Join(errList, ",")
}
