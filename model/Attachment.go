package model

import (
	"path/filepath"
	"web_blog/utils/errmsg"

	"gorm.io/gorm"
)

// Attachment 文章附件模型
type Attachment struct {
	gorm.Model
	ArticleID uint   `gorm:"type:int;not null" json:"article_id"`
	FileName  string `gorm:"type:varchar(255);not null" json:"fileName"`
	FilePath  string `gorm:"type:varchar(255);not null" json:"filePath"`
	FileSize  int64  `gorm:"type:bigint;not null" json:"fileSize"`
	FileType  string `gorm:"type:varchar(100);not null" json:"fileType"`
}

// CreateAttachment 新增附件
func CreateAttachment(data *Attachment) int {
	// 确保文件名不为空
	if data.FileName == "" && data.FilePath != "" {
		data.FileName = filepath.Base(data.FilePath)
	}

	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// GetAttachments 获取文章的所有附件
func GetAttachments(articleID int) ([]Attachment, int) {
	var attachments []Attachment
	err := db.Where("article_id = ?", articleID).Find(&attachments).Error
	if err != nil {
		return nil, errmsg.ERROR
	}

	// 处理每个附件，确保信息完整
	for i := range attachments {
		if attachments[i].FileName == "" && attachments[i].FilePath != "" {
			attachments[i].FileName = filepath.Base(attachments[i].FilePath)
		}
	}

	return attachments, errmsg.SUCCESS
}

// DeleteAttachment 删除附件
func DeleteAttachment(id int) int {
	var attachment Attachment
	err := db.Where("id = ?", id).Delete(&attachment).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// GetAttachmentInfo 获取单个附件信息
func GetAttachmentInfo(id int) (Attachment, int) {
	var attachment Attachment
	err := db.Where("id = ?", id).First(&attachment).Error
	if err != nil {
		return attachment, errmsg.ERROR
	}

	// 确保文件名不为空
	if attachment.FileName == "" && attachment.FilePath != "" {
		attachment.FileName = filepath.Base(attachment.FilePath)
	}

	return attachment, errmsg.SUCCESS
}

// BatchCreateAttachments 批量创建附件
func BatchCreateAttachments(attachments *[]Attachment) int {
	// 确保所有附件都有完整信息
	for i := range *attachments {
		if (*attachments)[i].FileName == "" && (*attachments)[i].FilePath != "" {
			(*attachments)[i].FileName = filepath.Base((*attachments)[i].FilePath)
		}
	}

	err := db.Create(&attachments).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
