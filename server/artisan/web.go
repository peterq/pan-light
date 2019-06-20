package artisan

import "github.com/kataras/iris/context"

func ApiHandler(handler func(ctx context.Context, param map[string]interface{}) (result interface{}, err error)) func(ctx context.Context) {

	return func(ctx context.Context) {
		var param map[string]interface{}
		var result interface{}

		err := ctx.ReadJSON(&param)
		if err == nil {
			defer func() {
				if err := recover(); err != nil {
					if e1, ok := err.(error); !ok {
						panic(e1)
					}
					ctx.JSON(map[string]interface{}{
						"success": false,
						"message": err.(error).Error(),
						"code":    ErrorFrom(err.(error)).Code(),
					})
				}
			}()
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
