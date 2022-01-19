/**
 * @Time: 2021/2/26 2:41 下午
 * @Author: varluffy
 */

package ginx

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/varluffy/rich/errcode"
	"reflect"
	"strings"
)

var _ error = (*ValidateError)(nil)
var trans ut.Translator

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
		return warpError(err)
	}
	return nil
}

func ShouldBindUri(c *gin.Context, v interface{}) error {
	if err := c.ShouldBindUri(v); err != nil {
		return warpError(err)
	}
	return nil
}

func warpError(err error) error {
	switch err.(type) {
	case validator.ValidationErrors:
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

func TransInit(local string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New() //chinese
		enT := en.New() //english
		uni := ut.New(enT, zhT, enT)

		var o bool
		trans, o = uni.GetTranslator(local)
		if !o {
			return fmt.Errorf("uni.GetTranslator(%s) failed", local)
		}
		//register translate
		// 注册翻译器
		switch local {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = chTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}
