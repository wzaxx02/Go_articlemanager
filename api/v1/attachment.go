package v1

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
	"web_blog/model"
	"web_blog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// 上传路径 - 确保这与实际存储文件的位置匹配
const uploadPath = "static/uploads/attachments/"

// UploadAttachment 上传单个附件
func UploadAttachment(c *gin.Context) {
	// 获取文章ID
	articleID, _ := strconv.Atoi(c.PostForm("article_id"))
	if articleID <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
			"error":   "无效的文章ID",
		})
		return
	}

	// 验证文章是否存在
	_, code := model.GetArtInfo(articleID)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}

	// 获取上传文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
			"error":   err.Error(),
		})
		return
	}

	// 确保上传目录存在
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
			"error":   "创建上传目录失败: " + err.Error(),
		})
		return
	}

	// 生成唯一文件名，保留原始扩展名
	fileExt := filepath.Ext(file.Filename)
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + fileExt
	filePath := path.Join(uploadPath, fileName)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
			"error":   "保存文件失败: " + err.Error(),
		})
		return
	}

	// 创建附件记录
	attachment := model.Attachment{
		ArticleID: uint(articleID),
		FileName:  file.Filename, // 保存原始文件名
		FilePath:  filePath,      // 保存文件实际路径
		FileSize:  file.Size,
		FileType:  fileExt[1:], // 去掉点号
	}

	code = model.CreateAttachment(&attachment)
	if code != errmsg.SUCCESS {
		// 删除已上传的文件
		os.Remove(filePath)
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"data":    attachment,
		"url":     filePath,
	})
}

// BatchUploadAttachments 批量上传附件
func BatchUploadAttachments(c *gin.Context) {
	// 获取文章ID
	articleID, _ := strconv.Atoi(c.PostForm("article_id"))
	if articleID <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
			"error":   "无效的文章ID",
		})
		return
	}

	// 验证文章是否存在
	_, code := model.GetArtInfo(articleID)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}

	// 获取上传文件
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
			"error":   err.Error(),
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
			"error":   "未上传任何文件",
		})
		return
	}

	// 确保上传目录存在
	err = os.MkdirAll(uploadPath, 0755)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
			"error":   "创建上传目录失败: " + err.Error(),
		})
		return
	}

	var successFiles []model.Attachment
	var failedFiles []string
	var attachments []model.Attachment

	// 处理每个文件
	for _, file := range files {
		// 生成唯一文件名
		fileExt := path.Ext(file.Filename)
		fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + fileExt
		filePath := path.Join(uploadPath, fileName)

		// 保存文件
		err = c.SaveUploadedFile(file, filePath)
		if err != nil {
			failedFiles = append(failedFiles, file.Filename)
			continue
		}

		// 创建附件记录
		attachment := model.Attachment{
			ArticleID: uint(articleID),
			FileName:  file.Filename,
			FilePath:  filePath,
			FileSize:  file.Size,
			FileType:  fileExt[1:], // 去掉点号
		}

		attachments = append(attachments, attachment)
		successFiles = append(successFiles, attachment)
	}

	// 批量保存到数据库
	if len(attachments) > 0 {
		code = model.BatchCreateAttachments(&attachments)
		if code != errmsg.SUCCESS {
			// 删除已上传的文件
			for _, attachment := range attachments {
				os.Remove(attachment.FilePath)
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        code,
		"message":       errmsg.GetErrMsg(code),
		"success_count": len(successFiles),
		"failed_count":  len(failedFiles),
		"success_files": successFiles,
		"failed_files":  failedFiles,
	})
}

// GetAttachments 获取文章的所有附件
func GetAttachments(c *gin.Context) {
	articleID, _ := strconv.Atoi(c.Param("id"))
	attachments, code := model.GetAttachments(articleID)

	// 确保返回完整的附件信息
	if code == errmsg.SUCCESS && attachments != nil {
		// 这里可以添加额外处理，确保每个附件的信息完整
		for i := range attachments {
			if attachments[i].FileName == "" {
				// 如果文件名为空，尝试从文件路径中提取
				if attachments[i].FilePath != "" {
					attachments[i].FileName = filepath.Base(attachments[i].FilePath)
				} else {
					attachments[i].FileName = "未命名文件." + attachments[i].FileType
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"data":    attachments,
	})
}

// DeleteAttachment 删除附件
func DeleteAttachment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// 先获取附件信息
	attachment, code := model.GetAttachmentInfo(id)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}

	// 删除物理文件
	if err := os.Remove(attachment.FilePath); err != nil {
		// 记录错误但继续删除数据库记录
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "删除物理文件失败: " + err.Error(),
			"path":    attachment.FilePath,
		})
		return
	}

	// 删除数据库记录
	code = model.DeleteAttachment(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// DownloadAttachment 下载附件 - 改进版本
func DownloadAttachment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// 检查token (允许通过查询参数传递token)
	token := c.Query("token")
	if token == "" {
		// 如果查询参数中没有token，尝试从header中获取
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}
	}

	// 获取附件信息
	attachment, code := model.GetAttachmentInfo(id)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}

	// 检查文件是否存在
	fileInfo, err := os.Stat(attachment.FilePath)
	if os.IsNotExist(err) {
		// 尝试修正路径 - 可能文件保存在uploads而不是static/uploads
		altPath := filepath.Join("uploads", "attachments", filepath.Base(attachment.FilePath))
		fileInfo, err = os.Stat(altPath)

		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR,
				"message": errmsg.GetErrMsg(errmsg.ERROR),
				"error":   "文件不存在，路径: " + attachment.FilePath + ", 备选路径: " + altPath,
			})
			return
		}

		// 更新路径
		attachment.FilePath = altPath
	}

	// 设置正确的MIME类型
	mimeType := "application/octet-stream" // 默认值
	ext := filepath.Ext(attachment.FileName)
	if ext != "" {
		switch filepath.Ext(attachment.FileName)[1:] {
		case "mp3":
			mimeType = "audio/mpeg"
		case "wav":
			mimeType = "audio/wav"
		case "pdf":
			mimeType = "application/pdf"
		case "jpg", "jpeg":
			mimeType = "image/jpeg"
		case "png":
			mimeType = "image/png"
		case "mp4":
			mimeType = "video/mp4"
		case "doc", "docx":
			mimeType = "application/msword"
		case "xls", "xlsx":
			mimeType = "application/vnd.ms-excel"
		}
	}

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename="+attachment.FileName)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", mimeType)
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// 发送文件
	c.File(attachment.FilePath)
}
