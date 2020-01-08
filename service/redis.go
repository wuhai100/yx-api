package service

import (
  "github.com/labstack/echo/v4"
  "time"
  "yx-api/util"
  "fmt"
  "gopkg.in/redis.v5"
)

const (
  keyUserInfo      = "yx_user_info"      // 用户信息
)

var redisLog = util.AppLog.With("file", "service.redis.go")
var redisCache *redis.Client

func OpenRedis() *redis.Client {
  log := redisLog.With("func", "OpenRedis")
  fmt.Println("redis info:", conf.RedisAddr, conf.RedisPwd)
  redisCache = redis.NewClient(&redis.Options{Addr: conf.RedisAddr, Password: conf.RedisPwd})
  sc := redisCache.Ping()
  log.Infof("Connectiong to %v ", sc.Val())
  if sc.Err() != nil {
    log.Errorf("OpenRedis conection set up failed, %s\n", sc.Err())
    panic(sc.Err())
  }
  log.Infof(redisCache.String() + "\n")
  log.Infof("Redis conection set up successfully \n")
  return redisCache
}

func CloseRedis() {
  log := redisLog.With("func", "CloseRedis")
  err := redisCache.Close()
  if err != nil {
    log.Errorf("close error %v ", err)
  }
}

func PreventFrequentOperation(intervalStr string) func(echo.Context) error {
  interval, err := time.ParseDuration(intervalStr)
  if err != nil {
    fmt.Printf("error parse [%s] to time duration\n", intervalStr)
    return nil
  }

  return func(ctx echo.Context) error {
    c, _ := ctx.FormParams()
    key := fmt.Sprintf("%v,%v,%v,%v", ctx.Request().Method, ctx.Request().URL.Path, c.Get("sign"), c.Get("token"))
    pre, _ := redisCache.Get(key).Int64()
    previous := time.Unix(pre, 0)
    now := time.Now()
    if previous.Add(interval).After(now) {
      return util.CodeBizError("700", "您操作过快哦！")
    } else {
      redisCache.Set(key, now.Unix(), 60*time.Second)
    }
    return nil
  }
}
