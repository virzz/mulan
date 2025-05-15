package req

import (
	"reflect"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func init() {
	binding.Validator = nil
}

var Validator binding.StructValidator = &defaultValidator{}

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *defaultValidator) ValidateStruct(obj any) error {
	if obj == nil {
		return nil
	}
	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(binding.SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

// validateStruct receives struct type
func (v *defaultValidator) validateStruct(obj any) error {
	v.lazyinit()
	return v.validate.Struct(obj)
}

func (v *defaultValidator) Engine() any {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
	})
}

func ShouldBind(c *gin.Context, obj any) error {
	err := binding.Default(c.Request.Method, c.ContentType()).
		Bind(c.Request, obj)
	if err != nil {
		zap.L().Warn("Failed to bind", zap.Error(err))
	}
	if len(c.Params) > 0 {
		if err := c.ShouldBindUri(obj); err != nil {
			zap.L().Warn("Failed to bind uri", zap.Error(err))
		}
	}
	return Validator.ValidateStruct(obj)
}

func Bind(c *gin.Context, obj any) error {
	binding.Default(c.Request.Method, c.ContentType()).Bind(c.Request, obj)
	if len(c.Params) > 0 {
		c.ShouldBindUri(obj)
	}
	return Validator.ValidateStruct(obj)
}
