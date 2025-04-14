package v1

import (
	"net/http"
	"web_blog/model"
	"web_blog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// UpLoad 上传图片接口
func UpLoad(c *gin.Context) {
	// 获取上传的文件
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "无法获取上传文件",
		})
		return
	}

	fileSize := fileHeader.Size

	// 调用模型层的文件上传方法
	url, code := model.UpLoadFile(file, fileSize, fileHeader)

	// 返回 JSON 响应
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"url":     url,
	})
}
