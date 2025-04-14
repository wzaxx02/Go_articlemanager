package model

import (
	"web_blog/utils/errmsg"

	"gorm.io/gorm"
)

type Bookmark struct {
	gorm.Model
	UserId    uint `json:"user_id"`
	ArticleId uint `json:"article_id"`
}

func AddBookmark(data *Bookmark) int {
	err = db.Save(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

func DeleteBookmark(userId, articleId uint) int {
	result := db.Where("user_id = ? AND article_id = ?", userId, articleId).Delete(&Bookmark{})
	if result.Error != nil {
		return errmsg.ERROR // 返回错误码
	}
	return errmsg.SUCCESS
}

func CheckBookMark(userId, articleId uint) bool {
	var bookmark Bookmark
	err = db.Where("article_id = ? AND user_id = ? AND deleted_at IS NULL", articleId, userId).First(&bookmark).Error

	return err == nil
}

func GetBookmarkListFront(id int, pageSize int, pageNum int) (Bookmark, int64, int) {
	var bookmarklist Bookmark
	var total int64
	db.Model(&Bookmark{}).Where("article_id = ?", id).Count(&total)
	err = db.Model(&Bookmark{}).Limit(pageSize).Offset((pageNum-1)*pageSize).Order("bookmark.created_at DESC").Select("bookmark.id, article.title, user_id, article_id, user.username, bookmark.created_at").Joins("LEFT JOIN article ON bookmark.article_id = article.id").Joins("LEFT JOIN user ON bookmark.user_id = user.id").Where("bookmark.user_id = ?", id).Scan(&bookmarklist).Error
	if err != nil {
		return bookmarklist, 0, errmsg.ERROR
	}
	return bookmarklist, total, errmsg.SUCCESS
}
