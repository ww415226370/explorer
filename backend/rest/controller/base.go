package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/types"
	"github.com/irisnet/explorer/backend/utils"
	"github.com/irisnet/irishub-sync/module/logger"
	"net/http"
	"strconv"
	"time"
)

func GetString(request *http.Request, key string) (result string) {
	apiName := request.RequestURI

	request.ParseForm()
	if len(request.Form[key]) > 0 {
		result = request.Form[key][0]
		logger.Info("Api Param", logger.String("Api", apiName), logger.String(key, result))
	}
	return
}

func GetInt(request *http.Request, key string) (result int) {
	value := GetString(request, key)
	if len(value) == 0 {
		return
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		logger.Error("param is not int type", logger.String("param", key))
	}
	return
}

func Var(request *http.Request, key string) (result string) {
	args := mux.Vars(request)
	result = args[key]
	return
}

func GetPage(r *http.Request) (int, int) {
	page := Var(r, "page")
	size := Var(r, "size")
	iPage := 0
	iSize := 20
	if p, ok := utils.ParseInt(page); ok {
		iPage = int(p)
	}
	if s, ok := utils.ParseInt(size); ok {
		iSize = int(s)
	}
	return iPage, iSize
}

func buildResponse(data ...interface{}) model.Response {
	var resp = model.Response{
		Code: types.ErrorCodeSuccess.Code,
	}

	if len(data) == 2 {
		if succ, ok := data[1].(bool); ok && !succ {
			err := data[0].(types.Error)
			resp.Code = err.Code
			resp.Msg = err.Msg
		}
	}
	resp.Data = data[0]
	return resp
}
func writeResponse(writer http.ResponseWriter, data ...interface{}) {
	//resp := buildResponse(data)
	resp := data[0]
	resultByte, err := json.Marshal(resp)
	if err != nil {
		logger.Error("json.Marshal failed")
	}
	writer.Write(resultByte)
}

// 用户处理逻辑函数
type Action func(request *http.Request) interface{}

//封装用户处理逻辑函数，捕获panic异常，统一处理
func wrap(action Action) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		apiName := request.RequestURI
		defer func() {
			if r := recover(); r != nil {
				switch r.(type) {
				case types.Error:
					writeResponse(writer, r, false)
					break
				case error:
					err := r.(error)
					e := types.Error{
						Code: types.ErrorCodeUnKnown.Code,
						Msg:  err.Error(),
					}
					writeResponse(writer, e, false)
					break
				}
				logger.Error("API execute failed", logger.Any("errMsg", r))

			}
		}()

		start := time.Now()
		result := action(request)
		end := time.Now()

		bizTime := end.Unix() - start.Unix()

		logger.Info("process information", logger.String("Api", apiName), logger.Int64("coast(s)", bizTime), logger.Any("result", result))

		if bizTime >= 3 {
			logger.Warn("api coast most time", logger.String("Api", apiName))
		}

		writeResponse(writer, result)
	}
}

//处理api接口
// url : api路径
// method : api请求方法
// action : 用户处理逻辑
func doApi(r *mux.Router, url, method string, action Action) {
	wrapperAction := wrap(action)
	r.HandleFunc(url, wrapperAction).Methods(method)
}
