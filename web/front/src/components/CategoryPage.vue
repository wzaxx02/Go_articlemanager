<template>
    <v-container>
      <v-row>
        <v-col v-for="item in artList" :key="item.id" cols="12" md="4">
          <v-card class="ma-3" outlined @click="$router.push(`/article/detail/${item.ID}`)">
            <v-img :src="item.img" height="200px"></v-img>
            <v-card-title>{{ item.title }}</v-card-title>
            <v-card-subtitle>{{ item.desc }}</v-card-subtitle>
          </v-card>
        </v-col>
      </v-row>
  
      <!-- 分页 -->
      <v-pagination
        color="indigo"
        total-visible="7"
        v-model="queryParam.pagenum"
        :length="Math.ceil(total / queryParam.pagesize)"
        @input="getArtList"
      ></v-pagination>
    </v-container>
  </template>
  
  <script>
  export default {
    props: ['id'],  // 接收路由传递的分类 ID
    data() {
      return {
        artList: [],  // 文章列表
        queryParam: {
          pagesize: 5,
          pagenum: 1
        },
        total: 0
      }
    },
    mounted() {
      this.getArtList()  // 获取分类下的文章
    },
    methods: {
      async getArtList() {
        try {
          const { data: res } = await this.$http.get('article/category', {
            params: {
              categoryId: this.id,  // 使用传入的分类 ID
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