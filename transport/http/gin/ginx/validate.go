/**
 * @Time: 2021/2/26 2:41 下午
 * @Author: varluffy
 */

package ginx

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/varluffy/rich/errcode"
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

func ShouldBind(c *gin.Context, v interface{}) error {
	if err := c.ShouldBind(v); err != nil {
		return warpError(c, err)
	}
	return nil
}

func ShouldBindUri(c *gin.Context, v interface{}) error {
	if err := c.ShouldBindUri(v); err != nil {
		return warpError(c, err)
	}
	return nil
}

func warpError(c *gin.Context, err error) error {
	switch err.(type) {
	case validator.ValidationErrors:
		trans := c.Value("trans")
		utrans, _ := trans.(ut.Translator)
		errs := err.(validator.ValidationErrors)
		return errcode.New(400, translateErrors(utrans, errs))
	case *json.UnmarshalTypeError:
		unmarshalTypeError := err.(*json.UnmarshalTypeError)
		message := fmt.Errorf("%s 类型错误，期望类型 %s", unmarshalTypeError.Field, unmarshalTypeError.Type.String()).Error()
		return errcode.New(400, message)
	default:
		return errcode.New(400, "参数解析失败，未知类型")
	}
}

func translateErrors(trans ut.Translator, errs validator.ValidationErrors) string {
	errList := make([]string, 0)
	for _, e := range errs {
		msg := e.Translate(trans)
		errList = append(errList, msg)
	}
	return strings.Join(errList, ",")
}
