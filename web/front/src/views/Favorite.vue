<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card
          class="ma-3"
          v-for="item in artList"
          :key="item.id"
          link
          @click="$router.push(`/article/detail/${item.ID}`)"
        >
          <v-row no-gutters class="d-flex align-center">
            <v-avatar class="ma-3 hidden-sm-and-down" size="125" tile>
              <v-img :src="item.img"></v-img>
            </v-avatar>
            <v-col>
              <v-card-title>
                <v-chip color="purple" outlined label class="mr-3 white--text">{{
                  item.category_name
                }}</v-chip>
                <div>{{ item.title }}</div>
              </v-card-title>
              <v-card-subtitle class="mt-1" v-text="item.desc"></v-card-subtitle>
              <v-divider class="mx-4"></v-divider>
              <v-card-text class="d-flex align-center">
                <div class="d-flex align-center">
                  <v-icon class="mr-1" small>{{ 'mdi-calendar-month' }}</v-icon>
                  <span>{{
                    item.CreatedAt | dateformat('YYYY-MM-DD HH:MM')
                  }}</span>
                </div>
                <div class="mx-4 d-flex align-center">
                  <v-icon class="mr-1" small>{{ 'mdi-comment' }}</v-icon>
                  <span>{{ item.comment_count }}</span>
                </div>
                <div class="mx-1 d-flex align-center">
                  <v-icon class="mr-1" small>{{ 'mdi-eye' }}</v-icon>
                  <span>{{ item.read_count }}</span>
                </div>
              </v-card-text>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>
    
    <!-- 分页 -->
    <v-pagination
      color="indigo"
      total-visible="7"
      v-model="queryParam.pagenum"
      :length="Math.ceil(total / queryParam.pagesize)"
      @input="getBookMarkArticleList()"
    ></v-pagination>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      artList: [], // 文章列表
      queryParam: {
        pagesize: 5,
        pagenum: 1
      },
      total: 0,
      headers: {
        username: '',
        user_id: ''
      }
    }
  },
  mounted() {
    // 检查是否需要从cookie中恢复会话
    this.checkSessionStatus()
  },
  created() {
    // 先初始化headers，稍后在mounted中可能会更新
    this.headers = {
      username: window.sessionStorage.getItem('username'),
      user_id: window.sessionStorage.getItem('user_id')
    }
  },
  methods: {
    // 检查会话状态并从cookie恢复
    checkSessionStatus() {
      // 检查sessionStorage中是否有用户ID
      let userId = window.sessionStorage.getItem('user_id')
      
      // 如果sessionStorage中没有，尝试从cookie获取
      if (!userId) {
        userId = this.getCookie('user_id')
        const username = this.getCookie('username')
        const token = this.getCookie('token')
        
        // 如果cookie中有数据，则恢复会话
        if (userId && username && token) {
          window.sessionStorage.setItem('user_id', userId)
          window.sessionStorage.setItem('username', username)
          window.sessionStorage.setItem('token', token)
          
          // 更新headers
          this.headers = {
            username: username,
            user_id: userId
          }
        }
      }
      
      // 获取收藏文章列表
      if (this.headers.user_id) {
        this.getBookMarkArticleList()
      } else {
        // 如果未登录，可以显示提示或重定向
        this.$message.info('请先登录后查看收藏')
        this.$router.push('/')
      }
    },
    
    // 获取cookie的方法
    getCookie(name) {
      const nameEQ = name + '='
      const ca = document.cookie.split(';')
      for (let i = 0; i < ca.length; i++) {
        let c = ca[i]
        while (c.charAt(0) === ' ') c = c.substring(1, c.length)
        if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length)
      }
      return null
    },
    
    async getBookMarkArticleList() {
      try {
        // 确保有用户ID
        if (!this.headers.user_id) {
          return
        }
        
        const { data: res } = await this.$http.get('bookmarkedarticles', {
          params: {
            user_id: this.headers.user_id,
            pagesize: this.queryParam.pagesize,
            pagenum: this.queryParam.pagenum
          }
        })
        this.artList = res.data
        this.total = res.total
      } catch (error) {
        console.error('获取文章失败:', error)
      }
    }
  }
}
</script>