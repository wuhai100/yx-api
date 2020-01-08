package controller

import (
  "github.com/yb7/echoswg"
  "yx-api/model"
  "yx-api/service"
  "yx-api/util"
)


func getBannerPagination(req *struct{ model.BannerDto }) *util.ResponseData {
  return util.ResultPageData(service.GetBannerPagination(req.BannerDto))
}

func addBanner(user *service.AdminUser, req *struct{ Body service.BannerDto }) *util.ResponseData {
  req.Body.Uid = service.ExtractAdminUid(user)
  return util.ResultData(nil, service.AddBanner(req.Body))
}

func editBanner(user *service.AdminUser, req *struct{ Body service.BannerDto }) *util.ResponseData {
  req.Body.Uid = service.ExtractAdminUid(user)
  return util.ResultData(nil, service.EditBanner(req.Body))
}

func editBannerStatus(req *struct{ Body struct{ ID, Status int } }) *util.ResponseData {
  return util.ResultData(nil, service.EditBannerStatus(req.Body.ID, req.Body.Status))
}

func editBannerSort(req *struct { Body service.BannerSortDto }) *util.ResponseData {
  return util.ResultData(nil, service.EditBannerSort(req.Body))
}

func init() {
  f := echoswg.NewApiGroup(util.EchoInst, "轮播图相关接口", "/admin")
  f.SetDescription("Admin版本")
  f.GET("/banner/list", verifyAdminUser, getBannerPagination, "轮播图列表")
  f.POST("/banner", preventFrequent2s, verifyAdminUser, addBanner, "轮播图话题")
  f.PUT("/banner", preventFrequent2s, verifyAdminUser, editBanner, "轮播图话题")
  f.PUT("/banner/status", preventFrequent2s, verifyAdminUser, editBannerStatus, "编辑轮播图状态")
  f.PUT("/banner/sort", preventFrequent2s, verifyAdminUser, editBannerSort, "调整轮播图顺序")

  p := echoswg.NewApiGroup(util.EchoInst, "APP端轮播图信息", "/v1")
  p.SetDescription("APP版本")
}
