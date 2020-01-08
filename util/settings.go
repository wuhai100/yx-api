package util

import (
  "os"
)

type Settings struct {
  DBUrl        string // 只读数据库
  WriterDBUrl  string // 写操作数据库
  RedisAddr    string // Redis服务器
  RedisPwd     string // Redis密码
  AppPort      string // 服务端口号
  ApiEnv       string // 服务环境dev|test|prod
  DriverName   string // 数据库驱动
  Ak           string // 华为云
  Sk           string // 华为云
  Endpoint     string // 华为云
  WxAppID      string // 微信AppID 用于用户认证and支付
  WxSecret     string // 微信Secret 用于用户认证
  WxKey        string // 微信Key 用于支付
  WxMchId      string // 微信MchId 用于支付
  WxNotifyHost string // 微信支付回调地址
}

var settings = Settings{}

func Get() Settings {
  settings.DriverName = "mysql"
  settings.DBUrl = os.Getenv("DB_URL")
  settings.RedisAddr = os.Getenv("REDIS_ADDR")
  settings.RedisPwd = os.Getenv("REDIS_PWD")
  settings.AppPort = os.Getenv("APP_PORT")
  settings.ApiEnv = os.Getenv("API_ENV")
  settings.Ak = os.Getenv("AK")
  settings.Sk = os.Getenv("SK")
  settings.Endpoint = os.Getenv("ENDPOINT")
  settings.WxAppID = os.Getenv("WX_APPID")
  settings.WxSecret = os.Getenv("WX_SECRET")
  settings.WxKey = os.Getenv("WX_KEY")
  settings.WxMchId = os.Getenv("WX_MCHID")
  settings.WxNotifyHost = os.Getenv("WX_NOTIFY_HOST")
  return settings
}
