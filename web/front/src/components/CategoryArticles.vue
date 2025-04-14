<template>
    <v-app app>
      <TopBar></TopBar>
  
      <v-main class="grey lighten-3">
        <v-container>
          <v-row>
            <v-col cols="12" md="3">
              <Nav></Nav>
            </v-col>
            <v-col cols="12" md="9">
              <v-carousel>
                <v-carousel-item
                  v-for="(image, i) in images"
                  :key="i"
                  :src="image"
                ></v-carousel-item>
              </v-carousel>
  
              <!-- 动态分类 -->
              <v-row>
                <v-col v-for="(category, index) in cateList" :key="index" cols="12" md="3">
                  <v-card class="ma-2" outlined>
                    <v-card-title class="headline">
                      <router-link :to="`/category/${category.id}`">{{ category.name }}</router-link>
                    </v-card-title>
                  </v-card>
                </v-col>
              </v-row>
            </v-col>
            
          </v-row>
        </v-container>
      </v-main>
  
      <Footer></Footer>
    </v-app>
  </template>
  
  <script>W
  import TopBar from '../components/TopBar'
  import Footer from '../components/Footer'
  import Nav from '../components/Nav'
  export default {
    components: { TopBar, Footer, Nav },
    data() {
      return {
        images: [
          'https://news.scnu.edu.cn/media/image/2023/07/20230705b15134.jpg.v',
          'https://th.bing.com/th?id=OIP.e4ejkxOmaWxlnmR5-eh_2wHaFj&w=288&h=216&c=8&rs=1&qlt=90&o=6&dpr=2&pid=3.1&rm=2',
          'https://ts1.cn.mm.bing.net/th/id/R-C.efa96ab8ba5ecfb3fa5182819bc26aaf?rik=9P1VC6o4ZzB%2fWw&riu=http%3a%2f%2fimg.fotomen.cn%2f2013%2f05%2f2315.jpg&ehk=7nGOx8JBiYgRZjhzNJRF%2bA%2bJbsoRIk44plxoQU97wbI%3d&risl=&pid=ImgRaw&r=0',
        ],
        cateList: [], // 动态分类列表
      }
    },
    methods: {
      // 获取分类
      async GetCateList() {
        try {
          const { data: res } = await this.$http.get('category')
          this.cateList = res.data
        } catch (error) {
          console.error('获取分类失败:', error)
        }
      },
    },
    created() {
      this.GetCateList() // 页面加载时获取分类数据
    },
  }
  </script>