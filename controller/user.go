package controller

import (
  "yx-api/model"
  "yx-api/service"
  "yx-api/util"
  "github.com/labstack/echo/v4"
  "github.com/yb7/echoswg"
)

func adminLogin(req *struct { Body model.AdminUserDto }) *util.ResponseData {
  return util.ResultData(service.AdminLogin(req.Body))
}

func adminLogout(req *struct { Token string }) *util.ResponseData {
  return util.ResultData(nil, service.AdminLogout(req.Token))
}

func adminUserList(req *struct { model.AdminUserDto }) *util.ResponseData {
  return util.ResultData(service.GetAdminUserList(req.AdminUserDto))
}

func addAdminUser(req *struct { Body model.AdminUserDto }) *util.ResponseData {
  return util.ResultData(nil, service.AddAdminUser(req.Body))
}

func editAdminUser(req *struct { Body model.AdminUserDto }) *util.ResponseData {
  return util.ResultData(nil, service.EditAdminUser(req.Body))
}

func verifyAdminUser(req *struct{ Token string }) (*service.AdminUser, error) {
  return service.CheckAdminLogin(req.Token)
}

func login(req *struct { Body model.UserDto }) *util.ResponseData {
  return util.ResultData(service.Login(req.Body))
}

func register(req *struct { Body model.UserDto }) *util.ResponseData {
  return util.ResultData(nil, service.Register(req.Body))
}

func getCountry() *util.ResponseData {
  return util.ResultData(service.GetCountry())
}

func setUserFollow(user *service.User, req *struct { Body model.UserFollowDto }) *util.ResponseData {
  req.Body.Uid = service.ExtractUid(user)
  return util.ResultData(nil, service.SetUserFollow(req.Body))
}

func checkSign(ctx echo.Context) error {
  return util.CheckSign(ctx)
}

func verifyUserToken(ctx echo.Context) (*service.User, error) {
  return service.VerifyUserToken(ctx)
}

// 校验签名并且获取用户信息,如果token校验失败返回没有登陆状态
func verifyUserNotLogin(ctx echo.Context) (*service.User, error) {
  user, err := service.VerifyUserToken(ctx)
  if err == nil {
    return user, nil
  }
  if err.(*util.BizError).Code() == "603" {
    return nil, nil
  }
  return user, err
}

func init() {
  f := echoswg.NewApiGroup(util.EchoInst, "管理后端用户登录和信息", "/admin")
  f.SetDescription("Admin版本")
  f.POST("/login", preventFrequent2s, adminLogin, "用户登录")
  f.POST("/logout", verifyAdminUser, adminLogout, "用户登出")
  f.GET("/user/list", verifyAdminUser, adminUserList, "用户列表")
  f.POST("/user", preventFrequent2s, verifyAdminUser, addAdminUser, "新增用户")
  f.PUT("/user", preventFrequent2s, verifyAdminUser, editAdminUser, "修改用户")

  p := echoswg.NewApiGroup(util.EchoInst, "APP端用户信息", "/v1")
  p.SetDescription("APP版本")
  p.POST("/login", preventFrequent1s, checkSign, login, "用户登录")
  p.POST("/register", preventFrequent2s, checkSign, register, "用户注册")
  p.GET("/country", checkSign, getCountry, "国家列表")
  p.PUT("/user/follow", verifyUserToken, setUserFollow, "用户关注和取消关注")
}
