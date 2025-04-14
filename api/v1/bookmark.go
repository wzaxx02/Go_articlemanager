package v1

import (
	"net/http"
	"strconv"
	"web_blog/model"
	"web_blog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

func Addbookmark(c *gin.Context) {
	// 定义结构体用于绑定请求体
	var request struct {
		UserID    uint `json:"user_id"`
		ArticleID uint `json:"article_id"`
	}

	// 从请求体中解析 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数解析失败",
		})
		return
	}

	if model.CheckBookMark(request.UserID, request.ArticleID) {
		model.DeleteBookmark(request.UserID, request.ArticleID)
	}

	// 创建书签
	data := model.Bookmark{
		UserId:    request.UserID,
		ArticleId: request.ArticleID,
	}

	code := model.AddBookmark(&data)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

func DeleteBookmark(c *gin.Context) {
	// 定义结构体用于绑定请求体
	var request struct {
		UserID    uint `json:"user_id"`
		ArticleID uint `json:"article_id"`
	}

	// 从请求体中解析 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数解析失败",
		})
		return
	}

	// 删除书签
	code := model.DeleteBookmark(request.UserID, request.ArticleID)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

func GetBookmarkListFront(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	data, total, code := model.GetBookmarkListFront(id, pageSize, pageNum)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})

}

func CheckBookMark(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))
	articleId, _ := strconv.Atoi(c.Query("article_id"))

	isBookmarked := model.CheckBookMark(uint(userId), uint(articleId))

	m := make(map[string]bool)
	m["isBookmarked"] = isBookmarked

	c.JSON(http.StatusOK, gin.H{
		"status": "200",
		"data":   m,
	})
}
