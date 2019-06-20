package artisan

import (
	"errors"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

func ApiHandler(handler func(ctx context.Context, param map[string]interface{}) (result interface{}, err error)) func(ctx context.Context) {

	return func(ctx context.Context) {
		var param map[string]interface{}
		var result interface{}

		err := ctx.ReadJSON(&param)
		if err == nil {
			result, err = handler(ctx, param)
		} else {
			err = NewError("decode input json error", -1, err)
		}

		if err != nil {
			ctx.JSON(map[string]interface{}{
				"success": false,
				"message": err.Error(),
				"code":    ErrorFrom(err.(error)).Code(),
			})
		} else {
			ctx.JSON(map[string]interface{}{
				"success": true,
				"result":  result,
			})
		}
	}
}

func ApiRecover(ctx context.Context) {
	defer func() {
		e := recover()
		if e != nil {
			// 转换为error
			var err error
			var appErr AppError
			var ok bool
			if err, ok = e.(error); !ok {
				err = errors.New(fmt.Sprint(e))
			}

			// 转换为 app error
			if appErr, ok = e.(AppError); !ok {
				ctx.StatusCode(iris.StatusInternalServerError)
				appErr = NewError("internal server error", 500, err)
			}

			ctx.JSON(map[string]interface{}{
				"success": false,
				"message": appErr.Error(),
				"code":    appErr.Code(),
			})
			ctx.StopExecution()
		}
	}()
	ctx.Next()
}
