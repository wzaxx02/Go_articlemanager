<template>
  <div class="attachment-manager">
    <h3>文件附件</h3>
    
    <!-- 上传区域 -->
    <div class="upload-section">
      <a-upload
        :multiple="true"
        :action="uploadUrl"
        :headers="headers"
        :data="{ article_id: articleId }"
        :fileList="fileList"
        @change="handleChange"
        @remove="handleRemove"
        @progress="handleProgress"
      >
        <a-button type="primary" :disabled="!articleId || articleId <= 0">
          <a-icon type="upload" /> 上传附件
        </a-button>
        <span class="upload-hint" v-if="!articleId || articleId <= 0">请先保存文章后再上传附件</span>
      </a-upload>
    </div>
    
    <!-- 上传进度 -->
    <div v-if="uploadingFiles.length > 0" class="progress-section">
      <div v-for="(file, index) in uploadingFiles" :key="index" class="progress-item">
        <div class="progress-info">
          <span class="filename">{{ file.name }}</span>
          <span class="percent">{{ file.percent }}%</span>
        </div>
        <a-progress :percent="file.percent" :status="file.status" size="small" />
      </div>
    </div>
    
    <!-- 附件列表 -->
    <div class="attachment-list" v-if="attachments.length > 0">
      <h4>已上传附件列表</h4>
      <a-table :dataSource="attachments" :columns="columns" rowKey="ID" size="small">
        <template slot="action" slot-scope="text, record">
          <a-button type="link" @click="downloadAttachment(record)">
            <a-icon type="download" /> 下载
          </a-button>
          <a-button type="link" @click="copyAttachmentLink(record.ID)" style="margin: 0 8px;">
            <a-icon type="link" /> 复制链接
          </a-button>
          <a-button type="link" @click="deleteAttachment(record.ID)" class="delete-btn">
            <a-icon type="delete" /> 删除
          </a-button>
        </template>
        <template slot="fileSize" slot-scope="text">
          {{ formatFileSize(text) }}
        </template>
      </a-table>
    </div>
    
    <!-- 下载进度模态框 -->
    <a-modal
      title="下载进度"
      :visible="downloadModalVisible"
      :footer="null"
      @cancel="downloadModalVisible = false"
    >
      <div class="download-progress">
        <p>正在下载: {{ currentDownload.fileName }}</p>
        <a-progress :percent="currentDownload.percent" status="active" />
      </div>
    </a-modal>
  </div>
</template>

<script>
import { Url } from '../../plugin/http'

export default {
  props: {
    articleId: {
      type: [Number, String],
      default: 0
    }
  },
  
  data() {
    return {
      uploadUrl: Url + 'attachment/upload',
      headers: {},
      fileList: [],
      attachments: [],
      uploadingFiles: [], // 正在上传的文件列表及进度
      downloadModalVisible: false, // 下载进度模态框可见性
      currentDownload: { // 当前下载信息
        fileName: '',
        percent: 0,
        id: null
      },
      columns: [
        {
          title: '文件名',
          dataIndex: 'fileName',
          key: 'fileName',
          render: (text) => text || '未知文件名'
        },
        {
          title: '类型',
          dataIndex: 'fileType',
          key: 'fileType',
          render: (text) => text || '未知类型'
        },
        {
          title: '大小',
          dataIndex: 'fileSize',
          key: 'fileSize',
          scopedSlots: { customRender: 'fileSize' },
        },
        {
          title: '操作',
          key: 'action',
          scopedSlots: { customRender: 'action' },
        }
      ]
    }
  },
  
  mounted() {
    this.headers = { Authorization: `Bearer ${window.sessionStorage.getItem('token')}` }
    
    if (this.articleId && this.articleId > 0) {
      this.getAttachments()
    }
  },
  
  watch: {
    articleId(newVal) {
      if (newVal && newVal > 0) {
        this.getAttachments()
      } else {
        this.attachments = []
      }
    }
  },
  
  methods: {
    // 获取文章附件列表
    async getAttachments() {
      if (!this.articleId || this.articleId <= 0) return
      
      try {
        const { data: res } = await this.$http.get(`attachment/list/${this.articleId}`)
        if (res.status !== 200) {
          this.$message.error(res.message || '获取附件列表失败')
          return
        }
        
        // 确保每个附件的数据都是完整的
        this.attachments = res.data.map(attachment => {
          // 处理可能的字段名大小写问题 (ID vs id, FileName vs fileName 等)
          return {
            ID: attachment.ID || attachment.id || 0,
            fileName: attachment.fileName || attachment.FileName || '未知文件名',
            fileType: attachment.fileType || attachment.FileType || '未知类型',
            fileSize: attachment.fileSize || attachment.FileSize || 0,
            filePath: attachment.filePath || attachment.FilePath || '',
            // 其他可能的字段...
          }
        })
        
        console.log('附件数据:', this.attachments)
      } catch (error) {
        console.error('获取附件列表失败:', error)
        this.$message.error('获取附件列表失败，请检查网络连接')
      }
    },
    
    // 处理上传进度
    handleProgress(step, file) {
      // 查找是否已存在该文件
      const index = this.uploadingFiles.findIndex(item => item.uid === file.uid)
      
      if (index !== -1) {
        // 更新已存在文件的进度
        this.uploadingFiles[index].percent = Math.floor(step.percent)
        this.uploadingFiles[index].status = step.percent >= 100 ? 'success' : 'active'
      } else {
        // 添加新文件
        this.uploadingFiles.push({
          uid: file.uid,
          name: file.name,
          percent: Math.floor(step.percent),
          status: 'active'
        })
      }
      
      // 100%完成后，延迟移除上传进度条
      if (step.percent >= 100) {
        setTimeout(() => {
          this.uploadingFiles = this.uploadingFiles.filter(item => item.uid !== file.uid)
        }, 3000)
      }
    },
    
    // 处理上传状态变化
    handleChange(info) {
      let fileList = [...info.fileList]
      
      // 更新上传进度
      if (info.event) {
        this.handleProgress(info.event, info.file)
      }
      
      // 只显示最近上传的文件
      fileList = fileList.slice(-5)
      
      // 更新状态
      this.fileList = fileList
      
      // 处理上传状态
      if (info.file.status === 'done') {
        if (info.file.response && info.file.response.status === 200) {
          this.$message.success(`${info.file.name} 上传成功`)
          // 刷新附件列表
          this.getAttachments()
        } else {
          const errorMsg = info.file.response ? info.file.response.message : '未知错误'
          this.$message.error(`${info.file.name} 上传失败: ${errorMsg}`)
        }
      } else if (info.file.status === 'error') {
        this.$message.error(`${info.file.name} 上传失败`)
      }
    },
    
    // 下载附件并显示进度
    downloadAttachment(attachment) {
      if (!attachment || !attachment.ID) return;
      
      // 确保使用正确的文件名
      const fileName = attachment.fileName || `file.${attachment.fileType}`;
      
      // 设置下载信息
      this.currentDownload = {
        id: attachment.ID,
        fileName: fileName,
        percent: 0
      };
      
      // 显示下载进度模态框
      this.downloadModalVisible = true;
      
      // 创建下载链接
      const downloadUrl = `${Url}attachment/download/${attachment.ID}`;
      
      // 创建XMLHttpRequest对象来跟踪下载进度
      const xhr = new XMLHttpRequest();
      xhr.open('GET', downloadUrl, true);
      xhr.responseType = 'blob';
      
      // 设置认证头
      xhr.setRequestHeader('Authorization', this.headers.Authorization);
      
      // 监听进度事件
      xhr.addEventListener('progress', (event) => {
        if (event.lengthComputable) {
          const percentComplete = Math.round((event.loaded / event.total) * 100);
          this.currentDownload.percent = percentComplete;
          
          // 如果完成，延迟关闭模态框
          if (percentComplete >= 100) {
            setTimeout(() => {
              this.downloadModalVisible = false;
            }, 1000);
          }
        }
      });
      
      // 下载完成时处理
      xhr.addEventListener('load', () => {
        if (xhr.status === 200) {
          // 获取正确的MIME类型
          let mimeType = 'application/octet-stream'; // 默认值
          
          // 尝试根据文件类型设置正确的MIME类型
          if (attachment.fileType) {
            const fileType = attachment.fileType.toLowerCase();
            if (fileType === 'mp3') {
              mimeType = 'audio/mpeg';
            } else if (fileType === 'mp4') {
              mimeType = 'video/mp4';
            } else if (fileType === 'pdf') {
              mimeType = 'application/pdf';
            } else if (['jpg', 'jpeg'].includes(fileType)) {
              mimeType = 'image/jpeg';
            } else if (fileType === 'png') {
              mimeType = 'image/png';
            } else if (fileType === 'doc' || fileType === 'docx') {
              mimeType = 'application/msword';
            } else if (fileType === 'xls' || fileType === 'xlsx') {
              mimeType = 'application/vnd.ms-excel';
            }
          }
          
          // 创建URL
          const blob = new Blob([xhr.response], { type: mimeType });
          const url = window.URL.createObjectURL(blob);
          
          // 创建临时链接并触发下载
          const a = document.createElement('a');
          a.style.display = 'none';
          a.href = url;
          a.download = fileName; // 使用原始文件名
          document.body.appendChild(a);
          a.click();
          
          // 清理
          window.URL.revokeObjectURL(url);
          document.body.removeChild(a);
          
          // 完成进度条
          this.currentDownload.percent = 100;
          
          // 延迟关闭模态框
          setTimeout(() => {
            this.downloadModalVisible = false;
          }, 1000);
        } else {
          this.$message.error('下载文件失败');
          this.downloadModalVisible = false;
        }
      });
      
      // 错误处理
      xhr.addEventListener('error', () => {
        console.error('下载错误');
        this.$message.error('下载文件时发生错误');
        this.downloadModalVisible = false;
      });
      
      // 开始下载
      xhr.send();
    },
    
    // 复制附件链接
    copyAttachmentLink(id) {
      const link = `${Url}attachment/download/${id}`;
      
      // 使用Clipboard API
      navigator.clipboard.writeText(link).then(() => {
        this.$message.success('链接已复制到剪贴板');
      }).catch(err => {
        console.error('复制失败:', err);
        this.$message.error('复制链接失败');
      });
    },
    
    // 删除附件
    async deleteAttachment(id) {
      this.$confirm({
        title: '确认删除',
        content: '确定要删除这个附件吗？删除后无法恢复。',
        okText: '确定',
        cancelText: '取消',
        onOk: async () => {
          try {
            const { data: res } = await this.$http.delete(`attachment/${id}`)
            if (res.status !== 200) {
              this.$message.error(res.message || '删除附件失败')
              return
            }
            
            this.$message.success('附件删除成功')
            this.getAttachments()
          } catch (error) {
            console.error('删除附件失败:', error)
            this.$message.error('删除附件失败，请检查网络连接')
          }
        }
      })
    },
    
    // 移除上传列表中的文件
    handleRemove(file) {
      const index = this.fileList.indexOf(file)
      const newFileList = this.fileList.slice()
      newFileList.splice(index, 1)
      this.fileList = newFileList
      
      // 同时移除进度条
      this.uploadingFiles = this.uploadingFiles.filter(item => item.uid !== file.uid);
    },
    
    // 格式化文件大小
    formatFileSize(size) {
      // 转换为数字以确保正确处理
      const numSize = Number(size);
      
      // 检查是否是有效数字
      if (isNaN(numSize) || numSize === undefined || numSize === null) {
        console.log('无效的文件大小值:', size);
        return '未知';
      }
      
      if (numSize < 1024) {
        return numSize + ' B';
      } else if (numSize < 1024 * 1024) {
        return (numSize / 1024).toFixed(2) + ' KB';
      } else if (numSize < 1024 * 1024 * 1024) {
        return (numSize / (1024 * 1024)).toFixed(2) + ' MB';
      } else {
        return (numSize / (1024 * 1024 * 1024)).toFixed(2) + ' GB';
      }
    }
  }
}
</script>

<style scoped>
.attachment-manager {
  margin-top: 20px;
  padding: 20px;
  background-color: #fafafa;
  border-radius: 4px;
  border: 1px solid #e8e8e8;
}

.upload-section {
  margin-bottom: 20px;
}

.upload-hint {
  margin-left: 10px;
  color: #ff4d4f;
}

.attachment-list {
  margin-top: 20px;
}

.delete-btn {
  color: #ff4d4f;
}

.progress-section {
  margin: 15px 0;
  padding: 10px;
  background-color: #fff;
  border: 1px dashed #d9d9d9;
  border-radius: 4px;
}

.progress-item {
  margin: 8px 0;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}

.filename {
  font-size: 14px;
  color: rgba(0, 0, 0, 0.65);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 80%;
}

.percent {
  font-size: 14px;
  color: rgba(0, 0, 0, 0.45);
}

.download-progress {
  padding: 10px;
}
</style>