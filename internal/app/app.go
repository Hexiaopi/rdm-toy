package app

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hexiaopi/rdm-toy/internal/retcode"
)

type CommResponse struct {
	*retcode.RetErr
	Data interface{} `json:"data,omitempty"`
}

type ErrResponse struct {
	*retcode.RetErr
}

type ListResponse struct {
	*retcode.RetErr
	Total int64       `json:"total"`
	Data  interface{} `json:"data,omitempty"`
}

func ToResponseCode(writer http.ResponseWriter, err error) {
	var response ErrResponse
	var code *retcode.RetErr
	if errors.As(err, &code) {
		response = ErrResponse{
			code,
		}
	}
	result, _ := json.Marshal(response)
	writer.Write(result)
}

func ToResponseData(writer http.ResponseWriter, data interface{}) {
	response := CommResponse{
		RetErr: retcode.Success,
		Data:   data,
	}
	result, _ := json.Marshal(response)
	writer.Write(result)
}

func ToResponseList(writer http.ResponseWriter, total int64, data interface{}) {
	response := ListResponse{
		RetErr: retcode.Success,
		Total:  total,
		Data:   data,
	}
	result, _ := json.Marshal(response)
	writer.Write(result)
}

type HandlerFunc func(ctx *gin.Context) error

func Wrap(handler HandlerFunc) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		method := ctx.Request.Method
		response := CommResponse{RetErr: retcode.Success}
		if err := handler(ctx); err != nil {
			var code *retcode.RetErr
			if errors.As(err, &code) {
				response.RetErr = code
				ctx.Error(code)
			} else {
				response.RetErr = retcode.UnknownError
				ctx.Error(retcode.UnknownError)
			}
		}
		retcode.RequestRetcodeCounter.WithLabelValues(path, method, response.Code, response.Desc).Inc()
		ctx.JSON(http.StatusOK, response)
	}
}

type HandlerDataFunc func(ctx *gin.Context) (res interface{}, err error)

func WrapData(handler HandlerDataFunc) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		method := ctx.Request.Method
		response := CommResponse{RetErr: retcode.Success}
		res, err := handler(ctx)
		if err != nil {
			var code *retcode.RetErr
			if errors.As(err, &code) {
				response.RetErr = code
				ctx.Error(code)
			} else {
				response.RetErr = retcode.UnknownError
				ctx.Error(retcode.UnknownError)
			}
		} else {
			response.Data = res
		}
		retcode.RequestRetcodeCounter.WithLabelValues(path, method, response.Code, response.Desc).Inc()
		ctx.JSON(http.StatusOK, response)
	}
}

type HandlerListFunc func(ctx *gin.Context) (res interface{}, total int64, err error)

func WrapList(handler HandlerListFunc) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		method := ctx.Request.Method
		response := ListResponse{RetErr: retcode.Success}
		res, total, err := handler(ctx)
		if err != nil {
			var code *retcode.RetErr
			if errors.As(err, &code) {
				response.RetErr = code
				ctx.Error(code)
			} else {
				response.RetErr = retcode.UnknownError
				ctx.Error(retcode.UnknownError)
			}
		} else {
			response.Total = total
			response.Data = res
		}
		retcode.RequestRetcodeCounter.WithLabelValues(path, method, response.Code, response.Desc).Inc()
		ctx.JSON(http.StatusOK, response)
	}
}
