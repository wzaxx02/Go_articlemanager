<template>
  <a-card>
    <a-form-model labelAlign="left" :label-col="{ span: 2 }" :wrapper-col="{ span: 12 }">


      <a-form-model-item label="网站备案号">
        <a-input style="width: 300px" v-model="profileInfo.icp_record"></a-input>
      </a-form-model-item>

      <a-form-model-item label="体育快讯">
        <a-input style="width: 300px" v-model="profileInfo.qq_chat"></a-input>
      </a-form-model-item>

      <a-form-model-item label="体育比分">
        <a-input style="width: 300px" v-model="profileInfo.wechat"></a-input>
      </a-form-model-item>

      <a-form-model-item label="youtube地址">
        <a-input style="width: 300px" v-model="profileInfo.weibo"></a-input>
      </a-form-model-item>

      <a-form-model-item label="Twitter地址">
        <a-input style="width: 300px" v-model="profileInfo.bili"></a-input>
      </a-form-model-item>

      <a-form-model-item label="facebook地址">
        <a-input style="width: 300px" v-model="profileInfo.email"></a-input>
      </a-form-model-item>



      <a-form-model-item>
        <a-button type="danger" style="margin-right: 15px" @click="updateProfile">更新</a-button>
      </a-form-model-item>
    </a-form-model>
  </a-card>
</template>
<script>
import { Url } from '../../plugin/http'

export default {
  data() {
    return {
      profileInfo: {
        id: 1,
        name: '',
        desc: '',
        qq_chat: '',
        wechat: '',
        weibo: '',
        bili: '',
        email: '',
        img: '',
        avatar: '',
        icp_record: '',
      },
      upUrl: Url + 'upload',
      headers: {},
    }
  },
  created() {
    this.getProfileInfo()
    this.headers = { Authorization: `Bearer ${window.sessionStorage.getItem('token')}` }
  },
  methods: {
    // 获取个人设置
    async getProfileInfo() {
      const { data: res } = await this.$http.get(`admin/profile/${this.profileInfo.id}`)
      if (res.status !== 200) {
        if (res.status === 1004 || 1005 || 1006 || 1007) {
          window.sessionStorage.clear()
          this.$router.push('/login')
        }
        this.$message.error(res.message)
      }
      this.profileInfo = res.data
    },

    // 上传头像
    avatarChange(info) {
      if (info.file.status !== 'uploading') {
      }
      if (info.file.status === 'done') {
        this.$message.success(`图片上传成功`)
        const imgUrl = info.file.response.url
        this.profileInfo.avatar = imgUrl
      } else if (info.file.status === 'error') {
        this.$message.error(`图片上传失败`)
      }
    },

    // 上传头像背景图
    imgChange(info) {
      if (info.file.status !== 'uploading') {
      }
      if (info.file.status === 'done') {
        this.$message.success(`图片上传成功`)
        const imgUrl = info.file.response.url
        this.profileInfo.img = imgUrl
      } else if (info.file.status === 'error') {
        this.$message.error(`图片上传失败`)
      }
    },

    // 更新
    async updateProfile() {
      const { data: res } = await this.$http.put(`profile/${this.profileInfo.id}`, this.profileInfo)
      if (res.status !== 200) return this.$message.error(res.message)
      this.$message.success(`个人信息更新成功`)
      this.$router.push('/index')
    },
  },
}
</script>

<style scoped>
.upBtn {
  margin-right: 10px;
}
</style>