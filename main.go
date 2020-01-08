package main

import (
  "github.com/yb7/asd"
  _ "yx-api/controller"
  "yx-api/service"
  "yx-api/util"
  "github.com/labstack/echo/v4/middleware"
  "github.com/yb7/echoswg"
)

func main() {
  conf := util.Get()

  service.OpenDB()
  defer service.CloseDB()

  service.OpenRedis()
  defer service.CloseRedis()

  // 华为的OBS桶配置初始
  service.InitObsClient()

  // asd缓存redis服务
  asd.InitRedisPool(asd.RedisOptions{})

  e := util.EchoInst
  e.Static("/swagger", "swagger")
  e.Static("/", "public")
  e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format: `${time_rfc3339} ${method} ${uri} ${status} cost:${latency_human} bytes:${bytes_in}->${bytes_out}}` + "\n",
  }))
  e.Use(middleware.Recover())

  e.GET("/api-docs", echoswg.GenApiDoc("API", "服务API"))

  e.Logger.Fatal(e.Start(":" + conf.AppPort))
}
