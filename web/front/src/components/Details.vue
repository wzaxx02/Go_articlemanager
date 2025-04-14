<template>
  <div>
    <div class="d-flex justify-center pa-3 ma-1 text-h4 font-weight-bold">{{ artInfo.title }}</div>
    <div class="d-flex justify-center align-center">
      <div class="d-flex mx-10 justify-center">
        <v-icon class="mr-1" color="indigo" small>
          mdi-calendar-month
        </v-icon>
        <span>{{ artInfo.CreatedAt | dateformat('YYYY-MM-DD') }}</span>
      </div>
      <div class="d-flex mr-10 justify-center">
        <v-icon class="mr-1" color="pink" small>mdi-comment</v-icon>
        <span>{{ total }}</span>
      </div>
      <div class="d-flex mr-10 justify-center">
        <v-icon class="mr-1" color="green" small>mdi-eye</v-icon>
        <span>{{ artInfo.read_count }}</span>
      </div>
      <div class="d-flex mr-10 justify-center" @click="collectBtnClick">
        <v-icon class="mr-1" color="green" small>{{ bookmarkIcon }}</v-icon>
        <span>{{ collectText }}</span>
      </div>
    </div>
    <v-divider class="pa-3 ma-3"></v-divider>
    <v-alert class="ma-4" elevation="1" color="indigo" dark border="left" outlined>
      {{ artInfo.desc }}
    </v-alert>

    <div v-html="artInfo.content" class="content ma-5 pa-3 text-justify"></div>

    <!-- 附件区域 -->
    <v-sheet v-if="attachments.length > 0" class="ma-5 pa-3">
      <v-divider class="mb-3"></v-divider>
      <h3 class="headline mb-3">文章附件</h3>
      <v-alert v-if="!headers.user_id" dense type="info" class="mb-3">
        请登录后查看和下载附件
      </v-alert>
      <v-list dense>
        <v-list-item
          v-for="attachment in attachments"
          :key="attachment.ID"
          @click="downloadAttachment(attachment.ID, attachment.fileName)"
          link
        >
          <v-list-item-icon>
            <v-icon :color="getFileIconColor(attachment.fileType)">{{ getFileIcon(attachment.fileType) }}</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>{{ attachment.fileName }}</v-list-item-title>
            <v-list-item-subtitle>{{ formatFileSize(attachment.fileSize) }}</v-list-item-subtitle>
          </v-list-item-content>
          <v-list-item-action>
            <v-btn icon small>
              <v-icon color="indigo">mdi-download</v-icon>
            </v-btn>
          </v-list-item-action>
        </v-list-item>
      </v-list>
    </v-sheet>

    <v-divider class="ma-5"></v-divider>
    <v-sheet class="ma-3 pa-3">
      <div>
        <v-list
          outlined
          class="ma-3 pa-3"
          v-for="item in commentList"
          :key="item.ID"
          v-show="item.status === 1"
        >
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title>
                {{ item.username }}
                {{ item.CreatedAt | dateformat('YYYY-MM-DD') }}
              </v-list-item-title>
              <v-list-item-subtitle class="mr-3">
                {{ item.content }}
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </div>
      <div class="text-center" v-if="commentList">
        <v-pagination
          class="my-2"
          total-visible="7"
          v-model="queryParam.pagenum"
          :length="Math.ceil(total / queryParam.pagesize)"
          @input="getCommentList()"
        ></v-pagination>
      </div>
      <div>
        <v-card flat>
          <v-alert v-if="!headers.username" class="ma-3" dense outlined type="error">
            你还未登录，请登录后留言
          </v-alert>
          <div v-if="headers.username">
            <v-textarea class="mx-3" outlined v-model="comment.content"></v-textarea>
            <v-btn class="ml-3 mb-1" dark color="indigo" small @click="pushComment()">确定</v-btn>
          </div>
        </v-card>
      </div>
    </v-sheet>
  </div>
</template>

<script>
export default {
  props: ['id'],
  data() {
    return {
      artInfo: {},
      commentList: [],
      comment: {
        content: ''
      },
      total: 0,
      headers: {
        username: '',
        user_id: 0
      },
      queryParam: {
        pagesize: 5,
        pagenum: 1
      },
      collect: false,
      attachments: [] // 文章附件列表
    }
  },
  computed: {
    bookmarkIcon() {
      return this.collect ? 'mdi-star' : 'mdi-star-outline';
    },
    collectText() {
      return this.collect ? '已收藏' : '收藏';
    }
  },
  created() {
    this.headers = {
      username: window.sessionStorage.getItem('username'),
      user_id: window.sessionStorage.getItem('user_id')
    }
  },
  mounted() {
    this.checkBookMarked();
    this.getArtInfo();
    this.getCommentList();
    this.getAttachments(); // 获取文章附件
  },
  methods: {
    handlePrint() { 
      alert(this.collect)
    },
    async collectBtnClick() {
      try {
        let res;
        if (this.collect) {
          res = await this.$http.delete(`delbookmark`, {
            data: { // Axios中delete请求的body需要使用data字段
              article_id: parseInt(this.id),
              user_id: parseInt(this.headers.user_id),
            }
          });
        } else {
          res = await this.$http.post(`addbookmark`, {
            article_id: parseInt(this.id),
            user_id: parseInt(this.headers.user_id),
          });
        }

        if (res.data.status !== 200) throw new Error(res.data.message);
        
        this.collect = !this.collect; // 切换状态
        this.$message.success(this.collect ? '收藏成功' : '取消收藏');
      } catch (error) {
        this.$message.error(error.message || '操作失败');
      }
    },
    // 查询文章收藏状态，并将结果转换为 Boolean 类型
    async checkBookMarked() {
      const { data: res } = await this.$http.get(`checkbookmarked`, {
        params: {
          article_id: this.id,
          user_id: parseInt(this.headers.user_id),
        }
      });
      // 确保返回 true 时，collect 为 true（兼容可能返回字符串的情况）
      this.collect = res.data.isBookmarked
    },
    // 查询文章信息
async getArtInfo() {
  try {
    const articleId = this.id;
    
    const { data: res } = await this.$http.get(`article/info/${articleId}`);
    
    if (res.status !== 200) {
      this.$message.error(res.message || '获取文章信息失败');
      return;
    }
    
    this.artInfo = res.data;
    
    if (res.data.ID !== undefined) {
      this.artInfo.id = res.data.ID;
    }
    
  } catch (error) {
    console.error('获取文章信息出错:', error);
    this.$message.error('获取文章信息失败，可能是网络问题');
  }
},
    // Add this method to increment the article read count
    async incrementReadCount(articleId) {
      try {
        // Make a separate request to increment the read count
        await this.$http.put(`article/read/${articleId}`);
      } catch (error) {
        console.error('增加阅读计数失败:', error);
        // Silently fail - not critical to user experience
      }
    },
    async getCommentList() {
      const { data: res } = await this.$http.get(`commentfront/${this.id}`, {
        params: {
          pagesize: this.queryParam.pagesize,
          pagenum: this.queryParam.pagenum
        }
      });
      this.commentList = res.data;
      this.total = res.total;
    },
    async pushComment() {
      const { data: res } = await this.$http.post('addcomment', {
        article_id: parseInt(this.id),
        content: this.comment.content,
        user_id: parseInt(this.headers.user_id),
        username: this.headers.username
      });
      if (res.status !== 200) return this.$message.error(res.message);
      this.$message.success('评论成功，待审核后显示');
      this.$router.go(0);
    },
    
    // 获取文章附件列表 - 修改后的方法
    async getAttachments() {
      try {
        if (!this.id) return;
        
        // 准备请求配置，如果用户已登录则添加认证头
        const config = {};
        if (this.headers.user_id) {
          config.headers = {
            'Authorization': `Bearer ${window.sessionStorage.getItem('token')}`
          };
        }
        
        const { data: res } = await this.$http.get(`attachment/list/${this.id}`, config);
        if (res.status !== 200) {
          // 仅在用户已登录时显示错误
          if (this.headers.user_id) {
            console.error('获取附件列表失败:', res.message);
          }
          return;
        }
        
        this.attachments = res.data;
      } catch (error) {
        console.error('获取附件列表出错:', error);
        // 仅在用户已登录时显示错误消息
        if (this.headers.user_id) {
          this.$message.error('获取附件列表失败，请重新登录');
        }
      }
    },
    
    // 下载附件 - 修改后的方法
    downloadAttachment(id, fileName) {
      if (!id) return;
      
      // 检查用户是否已登录
      if (!this.headers.user_id) {
        this.$message.warning('请先登录后再下载附件');
        return;
      }
      
      // 创建一个包含token的URL，确保URL格式正确
      const token = window.sessionStorage.getItem('token');
      let baseUrl = this.$http.defaults.baseURL;
      
      // 确保baseURL以斜杠结尾
      if (!baseUrl.endsWith('/')) {
        baseUrl += '/';
      }
      
      const downloadUrl = `${baseUrl}attachment/download/${id}?token=${token}`;
      console.log("下载URL:", downloadUrl); // 调试用
      
      // 创建一个隐形的a标签
      const a = document.createElement('a');
      a.style.display = 'none';
      a.href = downloadUrl;
      a.download = fileName || `attachment-${id}`;
      
      // 添加到文档，触发点击，然后移除
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      
      this.$message.success('开始下载，请稍等...');
    },
    
    // 根据文件类型获取图标
    getFileIcon(fileType) {
      if (!fileType) return 'mdi-file-outline';
      
      const type = fileType.toLowerCase();
      
      if (type.includes('pdf')) return 'mdi-file-pdf-box';
      if (type.includes('word') || type.includes('doc')) return 'mdi-file-word-box';
      if (type.includes('excel') || type.includes('xls')) return 'mdi-file-excel-box';
      if (type.includes('ppt') || type.includes('powerpoint')) return 'mdi-file-powerpoint-box';
      if (type.includes('image') || type.includes('jpg') || type.includes('png')) return 'mdi-file-image-box';
      if (type.includes('audio') || type.includes('mp3')) return 'mdi-file-music-box';
      if (type.includes('video') || type.includes('mp4')) return 'mdi-file-video-box';
      if (type.includes('zip') || type.includes('rar')) return 'mdi-zip-box';
      
      return 'mdi-file-outline';
    },
    
    // 根据文件类型获取图标颜色
    getFileIconColor(fileType) {
      if (!fileType) return 'grey';
      
      const type = fileType.toLowerCase();
      
      if (type.includes('pdf')) return 'red';
      if (type.includes('word') || type.includes('doc')) return 'blue';
      if (type.includes('excel') || type.includes('xls')) return 'green';
      if (type.includes('ppt') || type.includes('powerpoint')) return 'orange';
      if (type.includes('image')) return 'purple';
      if (type.includes('audio')) return 'teal';
      if (type.includes('video')) return 'pink';
      if (type.includes('zip') || type.includes('rar')) return 'amber';
      
      return 'grey';
    },
    
    // 格式化文件大小
    formatFileSize(size) {
      if (!size && size !== 0) return '未知大小';
      
      if (size < 1024) {
        return size + ' B';
      } else if (size < 1024 * 1024) {
        return (size / 1024).toFixed(2) + ' KB';
      } else if (size < 1024 * 1024 * 1024) {
        return (size / (1024 * 1024)).toFixed(2) + ' MB';
      } else {
        return (size / (1024 * 1024 * 1024)).toFixed(2) + ' GB';
      }
    }
  }
}
</script>

<style scoped>
.content >>> div,
img,
span {
  width: auto;
  max-width: 100%;
}

.content >>> pre,
code {
  margin: 10px;
  padding: 14px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: rgba(27, 31, 35, 0.05);
  border-left-width: 0.5rem;
  border-left-style: solid;
  border-color: #fdfdfd;
}
</style>