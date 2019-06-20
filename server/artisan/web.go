package artisan

import (
	"errors"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/peterq/pan-light/server/pc-api/middleware"
	"strings"
	"time"
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

type ThrottleOption struct {
	Duration time.Duration // 时间窗口
	Number   int           // 允许操作次数
	GetKey   func(ctx context.Context) string
}

type throttleState struct {
	Time  int64   `json:"time"`
	Water float64 `json:"water"`
}

func (o ThrottleOption) hit(ctx context.Context) time.Duration {
	key := o.GetKey(ctx)
	stateKey := "throttle-" + key

	var state throttleState
	RedisGet(stateKey, &state)
	// 计算现有数量
	speed := float64(o.Number) / float64(o.Duration/time.Second) // 每秒流失的水
	du := time.Duration(time.Now().UnixNano() - state.Time)
	water := state.Water - speed*float64(du/time.Second)
	if water < 0 {
		water = 0
	}
	water++
	if water > float64(o.Number) {
		return time.Duration((water-float64(o.Number))/speed) * time.Second
	}
	state.Time = time.Now().UnixNano()
	state.Water = water
	err := RedisSet(stateKey, state, o.Duration)
	if err != nil {
		panic(err)
	}
	return 0
}

// 频率限制器, 漏斗算法
func Throttle(options ...ThrottleOption) func(ctx context.Context) {
	for i, option := range options {
		if option.GetKey == nil {
			options[i].GetKey = func(ctx context.Context) string {
				return strings.Join([]string{
					ctx.GetCurrentRoute().Name(),
					middleware.CotextLoginInfo(ctx).Uk(),
					fmt.Sprint(option.Duration),
					fmt.Sprint(option.Number),
				}, ".")
			}
		}
	}
	return func(ctx context.Context) {
		for _, option := range options {
			ban := option.hit(ctx)
			if ban > 0 {
				panic(NewError(fmt.Sprintf("run out of %d call in %s, try after %s",
					option.Number, option.Duration, ban), iris.StatusTooManyRequests, nil))
			}
		}
		ctx.Next()
	}
}
