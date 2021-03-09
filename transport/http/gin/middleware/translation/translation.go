/**
 * @Time: 2021/2/25 6:11 下午
 * @Author: varluffy
 */

package translation

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

// Translation 参数绑定翻译
func Translation() gin.HandlerFunc {
	return func(c *gin.Context) {
		uni := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
		locale := c.GetHeader("locale")
		trans, _ := uni.GetTranslator(locale)
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			v.RegisterTagNameFunc(func(fld reflect.StructField) string {
				name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
				if name == "-" {
					return ""
				}
				return name
			})

			switch locale {
			case "zh":
				_ = zhtranslations.RegisterDefaultTranslations(v, trans)
				break
			case "en":
				_ = entranslations.RegisterDefaultTranslations(v, trans)
				break
			default:
				_ = zhtranslations.RegisterDefaultTranslations(v, trans)
				break
			}
			c.Set("trans", trans)
		}

		c.Next()
	}
}
