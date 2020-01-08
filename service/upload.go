package service

import (
  "bytes"
  "yx-api/util"
  "yx-api/util/obs"
  "fmt"
  "github.com/skip2/go-qrcode"
  "mime/multipart"
  "path"
  "strings"
)

// Obs跟文件夹
const obsRootDir = "yx"
// Obs访问路径前缀
const ObsURLPrefix = "https://obs-product-img.obs.cn-east-2.myhuaweicloud.com/"
// 创建ObsClient结构体
var obsClient *obs.ObsClient

type Image struct {
  URL  string `json:"url" desc:"图片URL"`
  Name string `json:"name" desc:"上传时图片名称"`
  ETag string `json:"eTag" desc:"上传后的文件标示，如果同一个文件这个tag是一样的"`
}

// 华为的云存储服务，该方法在main.go中启动时进行初始化
func InitObsClient() {
  obsClient, _ = obs.New(conf.Ak, conf.Sk, conf.Endpoint)
}

// 上传文件方法
func UploadFile(img *multipart.FileHeader) (Image, error) {
  suffix := path.Ext(img.Filename)
  fileName := fmt.Sprintf("%s/images/%s%s", obsRootDir, util.GetRandomString(10, 3), suffix)
  input := &obs.PutObjectInput{}
  input.Bucket = "obs-product-img"
  input.Key = fileName
  fd, err := img.Open()
  input.Body = fd
  output, err := obsClient.PutObject(input)
  if err == nil {
    fmt.Printf("RequestId:%s\n", output.RequestId)
    fmt.Printf("ETag:%s\n", output.ETag)
  } else if obsError, ok := err.(obs.ObsError); ok {
    fmt.Printf("Code:%s\n", obsError.Code)
    fmt.Printf("Message:%s\n", obsError.Message)
  }
  return Image{URL:ObsURLPrefix + fileName, Name:img.Filename, ETag:output.ETag}, err
}

// 订单中生成的提货二维码
func UploadQrCode(code string) (string, error) {
  var png []byte
  png, err := qrcode.Encode(code, qrcode.Medium, 256)
  fileName := fmt.Sprintf("qrcode/%s%s", util.GetRandomString(10, 3), code + ".png")
  input := &obs.PutObjectInput{}
  input.Bucket = "obs-product-img"
  input.Key = fileName
  input.Body = bytes.NewReader(png)
  output, err := obsClient.PutObject(input)
  if err == nil {
    fmt.Printf("RequestId:%s\n", output.RequestId)
    fmt.Printf("ETag:%s\n", output.ETag)
  } else if obsError, ok := err.(obs.ObsError); ok {
    fmt.Printf("Code:%s\n", obsError.Code)
    fmt.Printf("Message:%s\n", obsError.Message)
  }
  return fileName, err
}

func appendPrefix(fileName string) string {
  if fileName == "" {
    return fileName
  }
  return ObsURLPrefix + fileName
}

func clearPrefix(fileName string) string {
  return strings.ReplaceAll(fileName, ObsURLPrefix, "")
}

func obsPrefix(fileName string) string {
  if fileName == "" {
    return fileName
  } else if strings.HasPrefix(fileName, "https://") {
    return strings.ReplaceAll(fileName, ObsURLPrefix, "")
  }
  return ObsURLPrefix + fileName
}
