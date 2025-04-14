<template>
  <div>
    <a-card>
      <h3>{{ id ? '编辑文章' : '新增文章' }}</h3>

      <a-form-model :model="artInfo" ref="artInfoRef" :rules="artInfoRules"
        :hideRequiredMark="true">
        <a-row :gutter="24">
          <a-col :span="16">
            <a-form-model-item label="文章标题" prop="title">
              <a-input style="width: 300px" v-model="artInfo.title"></a-input>
            </a-form-model-item>
            <a-form-model-item label="文章描述" prop="desc">
              <a-input type="textarea" v-model="artInfo.desc"></a-input>
            </a-form-model-item>
          </a-col>
          <a-col :span="8">
            <a-form-model-item label="文章分类" prop="cid">
              <a-select style="width: 200px" v-model="artInfo.cid" placeholder="请选择分类"
                @change="cateChange">
                <a-select-option v-for="item in Catelist" :key="item.id" :value="item.id">
                  {{
                      item.name
                  }}
                </a-select-option>
              </a-select>
            </a-form-model-item>

            <a-form-model-item label="文章缩略图" prop="img">
              <a-upload listType="picture" name="file" :action="upUrl" :headers="headers"
                @change="upChange">
                <a-button>
                  <a-icon type="upload" />点击上传
                </a-button>

                <template v-if="id">
                  <img :src="artInfo.img" style="width: 120px; height: 100px; margin-left: 15px" />
                </template>
              </a-upload>
            </a-form-model-item>
          </a-col>
        </a-row>
        <a-form-model-item label="文章内容" prop="content">
          <Editor v-model="artInfo.content"></Editor>
        </a-form-model-item>

        <!-- 新增附件管理器组件 -->
        <AttachmentManager :articleId="artInfo.id" ref="attachmentManager" />

        <a-form-model-item>
          <a-button type="danger" style="margin-right: 15px" @click="artOk(artInfo.id)">
            {{
                artInfo.id ? '更新' : '提交'
            }}
          </a-button>
          <a-button type="primary" @click="addCancel">取消</a-button>
        </a-form-model-item>
      </a-form-model>
    </a-card>
  </div>
</template>

<script>
import { Url } from '../../plugin/http'
import Editor from '../editor/index'
// 导入附件管理器组件
import AttachmentManager from './AttachmentManager.vue'

export default {
  components: { 
    Editor,
    AttachmentManager // 注册附件管理器组件
  },
  props: ['id'],
  data() {
    return {
      artInfo: {
        id: 0,
        title: '',
        cid: undefined,
        desc: '',
        content: '',
        img: '',
      },
      Catelist: [],
      upUrl: Url + 'upload',
      headers: {},
      fileList: [],
      artInfoRules: {
        title: [{ required: true, message: '请输入文章标题', trigger: 'change' }],
        cid: [{ required: true, message: '请选择文章分类', trigger: 'change' }],
        desc: [
          { required: true, message: '请输入文章描述', trigger: 'change' },
          { max: 120, message: '描述最多可写120个字符', trigger: 'change' },
        ],
        // img: [{ required: true, message: '请选择文章缩略图', trigger: 'change' }],
        content: [{ required: true, message: '请输入文章内容', trigger: 'change' }],
      },
    }
  },
  mounted() {
    this.getCateList()
    this.headers = { Authorization: `Bearer ${window.sessionStorage.getItem('token')}` }
    if (this.id) {
      this.getArtInfo(this.id)
    }
  },
  methods: {
    // 查询文章信息
    async getArtInfo(id) {
      try {
        // 使用正确的API路径 - 后端提供的路径是 admin/article/info/:id
        const { data: res } = await this.$http.get(`admin/article/info/${id}`)
        
        if (res.status !== 200) {
          // 处理特定错误码，但不要立即清除会话
          if ([1004, 1005, 1006, 1007].includes(res.status)) {
            this.$message.error(res.message || '登录状态失效，请重新登录')
            // 仅当确认是授权问题时再清除会话并跳转
            window.sessionStorage.clear()
            this.$router.push('/login')
            return
          }
          this.$message.error(res.message || '获取文章信息失败')
          return
        }
        
        // 正确设置文章信息
        this.artInfo = res.data
        
        // 确保ID字段正确赋值
        // 后端返回的ID字段可能是大写的"ID"
        if (res.data.ID !== undefined) {
          this.artInfo.id = res.data.ID
        }
      } catch (error) {
        console.error('获取文章信息出错:', error)
        this.$message.error('获取文章信息失败，可能是网络问题或登录已过期')
      }
    },
    // 获取分类列表
    async getCateList() {
      const { data: res } = await this.$http.get('admin/category')
      if (res.status !== 200) return this.$message.error(res.message)
      this.Catelist = res.data
    },
    // 选择分类
    cateChange(value) {
      this.artInfo.cid = value
    },
    // 上传图片
    upChange(info) {
      if (info.file.status !== 'uploading') {
      }
      if (info.file.status === 'done') {
        this.$message.success(`图片上传成功`)
        this.artInfo.img = info.file.response.url
      } else if (info.file.status === 'error') {
        this.$message.error(`图片上传失败`)
      }
    },
    // 提交&&更新文章
    artOk(id) {
      this.$refs.artInfoRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数验证未通过，请按要求录入文章内容')
        if (id === 0) {
          const { data: res } = await this.$http.post('article/add', this.artInfo)
          if (res.status !== 200) return this.$message.error(res.message)
          
          // 更新文章ID，使附件管理器可以使用
          if (res.data && res.data.ID) {
            this.artInfo.id = res.data.ID
            // 提示用户现在可以上传附件
            this.$message.info('文章保存成功，现在可以上传附件了')
            // 等待DOM更新后刷新附件管理器
            this.$nextTick(() => {
              if (this.$refs.attachmentManager) {
                this.$refs.attachmentManager.getAttachments()
              }
            })
          } else {
            this.$router.push('/artlist')
            this.$message.success('添加文章成功')
          }
        } else {
          const { data: res } = await this.$http.put(`article/${id}`, this.artInfo)
          if (res.status !== 200) return this.$message.error(res.message)

          this.$router.push('/artlist')
          this.$message.success('更新文章成功')
        }
      })
    },

    addCancel() {
      this.$refs.artInfoRef.resetFields()
    },
  },
}
</script>