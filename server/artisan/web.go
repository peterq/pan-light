package artisan

import (
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/peterq/pan-light/server/artisan/cache"
	"github.com/peterq/pan-light/server/conf"
	"github.com/peterq/pan-light/server/pc-api/middleware"
	"runtime/debug"
	"strings"
	"time"
)

func ApiHandler(handler func(ctx iris.Context, param JsonMap) (result interface{}, err error)) func(ctx iris.Context) {

	return func(ctx iris.Context) {
		var param JsonMap
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

func ApiRecover(ctx iris.Context) {
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
			if appErr, ok = err.(AppError); !ok {
				ctx.StatusCode(iris.StatusInternalServerError)
				appErr = NewError("internal server error", 500, err)
				if conf.Conf.Debug {
					ctx.Application().Logger().Error(err, string(debug.Stack()))
				} else {
					ctx.Application().Logger().Error(err)
				}
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
	GetKey   func(ctx iris.Context) string
}

type throttleState struct {
	Time  int64   `json:"time"`
	Water float64 `json:"water"`
}

func (o ThrottleOption) hit(ctx iris.Context) time.Duration {
	key := o.GetKey(ctx)
	stateKey := "throttle-" + key

	var state throttleState
	cache.RedisGet(stateKey, &state)
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
	err := cache.RedisSet(stateKey, state, o.Duration)
	if err != nil {
		panic(err)
	}
	return 0
}

// 频率限制器, 漏斗算法
func Throttle(options ...ThrottleOption) func(ctx iris.Context) {
	for i, option := range options {
		if option.GetKey == nil {
			options[i].GetKey = func(ctx iris.Context) string {
				return strings.Join([]string{
					ctx.GetCurrentRoute().Name(),
					middleware.ContextLoginInfo(ctx).Uk(),
					fmt.Sprint(option.Duration),
					fmt.Sprint(option.Number),
				}, ".")
			}
		}
	}
	return func(ctx iris.Context) {
		if conf.Conf.Debug {
			ctx.Next()
			return
		}
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

type JsonMap map[string]interface{}
type JsonValue struct {
	name string
	data interface{}
}

func (m JsonMap) Get(keys ...string) JsonValue {
	if len(keys) == 1 {
		keys = strings.Split(keys[0], ".")
	}
	if len(keys) == 1 {
		return JsonValue{
			name: keys[0],
			data: m[keys[0]],
		}
	}
	parent := m.Get(keys[:len(keys)-1]...).Map()
	return JsonValue{
		name: strings.Join(keys, "."),
		data: parent[keys[len(keys)-1]],
	}
}

func (v JsonValue) Map() JsonMap {
	if m, ok := v.data.(map[string]interface{}); ok {
		return m
	}
	panic(NewError(fmt.Sprintf("%s needs to be map, %T given", v.name, v.data), -1, nil))
}

func (v JsonValue) String(ignore ...interface{}) string {
	if m, ok := v.data.(string); ok {
		return m
	}
	panic(NewError(fmt.Sprintf("%s needs to be string, %T given", v.name, v.data), -1, nil))
}

func (v JsonValue) Int() int {
	if m, ok := v.data.(float64); ok {
		return int(m)
	}
	panic(NewError(fmt.Sprintf("%s needs to be int, %T given", v.name, v.data), -1, nil))
}

func (v JsonValue) Int64() int64 {
	if m, ok := v.data.(float64); ok {
		return int64(m)
	}
	panic(NewError(fmt.Sprintf("%s needs to be int, %T given", v.name, v.data), -1, nil))
}

func (v JsonValue) Float() float64 {
	if m, ok := v.data.(float64); ok {
		return m
	}
	panic(NewError(fmt.Sprintf("%s needs to be float, %T given", v.name, v.data), -1, nil))
}

func (v JsonValue) Array() []JsonValue {
	if m, ok := v.data.([]interface{}); ok {
		arr := make([]JsonValue, len(m))
		for idx, value := range m {
			arr[idx] = JsonValue{
				name: v.name + " index of " + fmt.Sprint(idx),
				data: value,
			}
		}
		return arr
	}
	panic(NewError(fmt.Sprintf("%s needs to be array, %T given", v.name, v.data), -1, nil))
}

var App *iris.Application

func NewApp() *iris.Application {
	App = iris.New()
	return App
}
