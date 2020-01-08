package controller

import (
  "github.com/labstack/echo/v4"
  "net/http"
  "strconv"
  "yx-api/service"
  "yx-api/util"
)

var (
  log               = util.AppLog.With("file", "error.go")
  preventFrequent1s = service.PreventFrequentOperation("1s") // 防止重复提交设置1s
  preventFrequent2s = service.PreventFrequentOperation("2s") // 防止重复提交设置5s
)

type JsonError struct {
  httpStatus int
  json       []byte
  message    string
}

func (e *JsonError) Error() string {
  return e.message
}

func (e *JsonError) Json() []byte {
  return e.json
}

func (e *JsonError) HttpStatus() int {
  return e.httpStatus
}

// BindRoutes func
func init() {
  util.EchoInst.HTTPErrorHandler = func(err error, c echo.Context) {
    log.With("error_handler", strconv.Itoa(http.StatusInternalServerError)).Errorf("%v \n %v \n %v", c.Request().Method, c.Request().URL, err)

		if c.Response().Committed {
			return
		}

    if bizError, ok := err.(*util.BizError); ok {
      c.JSON(bizError.HttpStatus(), map[string]interface{} {
        "errno": bizError.Code(),
        "data":  bizError.Error(),
      })
      return
    }

		c.JSON(200, map[string]interface{} {
			"errno": 500,
			"data": err.Error(),
		})
	}
}
