package controller

import (
  "github.com/yb7/echoswg"
  "yx-api/model"
  "yx-api/service"
  "yx-api/util"
)


func getTagPagination(req *struct{ model.TagDto }) *util.ResponseData {
  return util.ResultPageData(service.GetTagPagination(req.TagDto))
}

func getEffectiveTag() *util.ResponseData {
  return util.ResultData(service.GetEffectiveTag())
}

func addTag(req *struct{ Body service.TagDto }) *util.ResponseData {
  return util.ResultData(nil, service.AddTag(req.Body))
}

func editTag(req *struct{ Body service.TagDto }) *util.ResponseData {
  return util.ResultData(nil, service.EditTag(req.Body))
}

func editTagStatus(req *struct{ Body struct{ ID, Status int } }) *util.ResponseData {
  return util.ResultData(nil, service.EditTagStatus(req.Body.ID, req.Body.Status))
}

func editTagSort(req *struct { Body service.TagSortDto }) *util.ResponseData {
  return util.ResultData(nil, service.EditTagSort(req.Body))
}

func init() {
  f := echoswg.NewApiGroup(util.EchoInst, "标签相关接口", "/admin")
  f.SetDescription("Admin版本")
  f.GET("/tag/list", verifyAdminUser, getTagPagination, "标签列表")
  f.GET("/tag/map", verifyAdminUser, getEffectiveTag, "标签列表Map")
  f.POST("/tag", preventFrequent2s, verifyAdminUser, addTag, "标签话题")
  f.PUT("/tag", preventFrequent2s, verifyAdminUser, editTag, "标签话题")
  f.PUT("/tag/status", preventFrequent2s, verifyAdminUser, editTagStatus, "编辑标签状态")
  f.PUT("/tag/sort", preventFrequent2s, verifyAdminUser, editTagSort, "调整标签顺序")

  p := echoswg.NewApiGroup(util.EchoInst, "APP端标签信息", "/v1")
  p.SetDescription("APP版本")
}
