package main

import (
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

type JSONResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Image string `json:"image"`
	} `json:"data"`
}

func main() {
	g := gin.Default()

	g.LoadHTMLFiles("index.html")
	g.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	g.GET("/qr", func(c *gin.Context) {
		responseMsg := ""
		data := c.Query("data")
		size := c.Query("size")
		returnType := c.DefaultQuery("type", "image")

		if len(data) == 0 || len(size) == 0 {
			responseMsg = "缺少必要参数"
			if returnType == "json" {
				c.JSON(http.StatusOK, JSONResponse{
					Code: 500,
					Msg:  responseMsg,
				})
			} else {
				c.String(http.StatusOK, responseMsg)
			}
			return
		}

		intSize, err := strconv.Atoi(size)
		if err != nil {
			responseMsg = "size参数不合法"
			if returnType == "json" {
				c.JSON(http.StatusOK, JSONResponse{
					Code: 500,
					Msg:  responseMsg,
				})
			} else {
				c.String(http.StatusOK, responseMsg)
			}
			return
		}

		if intSize >= 1000 {
			responseMsg = "size参数不合法，size必须小于等于1000"
			if returnType == "json" {
				c.JSON(http.StatusOK, JSONResponse{
					Code: 500,
					Msg:  responseMsg,
				})
			} else {
				c.String(http.StatusOK, responseMsg)
			}
			return
		}

		var png []byte
		png, err = qrcode.Encode(data, qrcode.Medium, intSize)
		if err != nil {
			responseMsg = "图片生成失败"
			if returnType == "json" {
				c.JSON(http.StatusOK, JSONResponse{
					Code: 500,
					Msg:  responseMsg,
				})
			} else {
				c.String(http.StatusOK, responseMsg)
			}
			return
		}

		if returnType == "image" {
			c.Data(http.StatusOK, "image/png", png)
		} else if returnType == "json" {
			response := JSONResponse{
				Code: http.StatusOK,
				Msg:  "Success",
				Data: struct {
					Image string "json:\"image\""
				}{
					Image: base64.StdEncoding.EncodeToString(png),
				},
			}
			c.JSON(http.StatusOK, response)
		} else {
			c.String(http.StatusOK, "无效返回参数")
		}
	})

	g.Run()
}
