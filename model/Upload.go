package model

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"web_blog/utils"
	"web_blog/utils/errmsg"
)

// UpLoadFile 上传文件到本地
func UpLoadFile(file multipart.File, fileSize int64, fileHeader *multipart.FileHeader) (string, int) {
	// 确保目标目录存在
	err := os.MkdirAll(utils.UploadPath, os.ModePerm)
	if err != nil {
		return "", errmsg.ERROR
	}

	// 确定保存的文件路径
	fileName := fileHeader.Filename
	savePath := filepath.Join(utils.UploadPath, fileName)

	// 创建目标文件
	dstFile, err := os.Create(savePath)
	if err != nil {
		return "", errmsg.ERROR
	}
	defer dstFile.Close()

	// 将上传的文件内容写入目标文件
	_, err = io.Copy(dstFile, file)
	if err != nil {
		return "", errmsg.ERROR
	}

	// 返回文件的访问路径
	url := "http://localhost:3000/uploads/" + fileName
	return url, errmsg.SUCCESS
}
