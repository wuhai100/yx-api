package controller

import (
  "yx-api/service"
  "yx-api/util"
  "github.com/yb7/echoswg"
  "github.com/labstack/echo/v4"
)

// 上传文件接口
func uploadFile(c echo.Context) *util.ResponseData {
  rd := &util.ResponseData{}
  file, err := c.FormFile("file")
  if err != nil {
    rd.Errno = 501
    rd.Data = err
    return rd
  }

  return util.ResultData(service.UploadFile(file))
}

// 上传文件
func init() {
  g := echoswg.NewApiGroup(util.EchoInst, "上传文件", "/admin")
	g.SetDescription("Admin版本")
  g.POST("/upload", uploadFile, "上传文件接口")
}
