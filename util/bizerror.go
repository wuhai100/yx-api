package util

import (
  "net/http"
)

type BizError struct {
  httpStatus int
  code       string
  message    string
}

func (e *BizError) HttpStatus() int {
  return e.httpStatus
}

func (e *BizError) Code() string {
  return e.code
}

func (e *BizError) Error() string {
  return e.message
}

func CustomBizError(message string) *BizError {
  return &BizError{http.StatusOK, "500", message}
}

func CodeBizError(code, message string) *BizError {
  return &BizError{http.StatusOK, code, message}
}
