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
              <v-carousel-item v-for="(image, i) in images" :key="i" :src="image"></v-carousel-item>
            </v-carousel>

            <!-- 动态分类 -->
            <v-row>
              <v-col v-for="(category, index) in cateList" :key="index" cols="12" md="3">
                <v-card class="ma-2" outlined>
                  <v-card-title class="headline" @click="toPage(category.id)">
                    <router-link :to="`/category/${category.id}`">{{ category.name }}</router-link>
                  </v-card-title>
                </v-card>
              </v-col>
            </v-row>

            <v-row>
              <v-col>
                <v-tabs v-model="tab" align-tabs="center" color="deep-purple-accent-4">
                  <v-tab text="Option 1" value="1">足球</v-tab>
                  <v-tab value="2">篮球</v-tab>
                  <v-tab value="3">乒乓球</v-tab>
                  <v-tab value="4">羽毛球</v-tab>
                </v-tabs>

                <v-tabs-items v-model="tab">
                  <v-tab-item key="1">
                    <v-list lines="two">
                      <v-list-item v-for="item in categoriesList" :key="item.ID">
                        <router-link :to="`/article/detail/${item.ID}`" class="article-link">
                          <div style="width: 100%;height: 12vh;display: flex;justify-content: space-around;margin: 1vh 0;">
                            <div
                              :style="{ width: '10%', height: '100%', backgroundSize: 'contain', backgroundImage: `url(${item.img})` }"
                              class="article-image"
                            >
                            </div>
                            <div style="width: 88%;height: 100%;">
                              <span style="font-size: 1.2vw;font-weight: 600;">{{ item.title }}</span>
                              <div v-html="item.content"
                                style="height: 3vh;margin: 0;color: rgb(100, 100, 100);font-size: 1.1vw;display: -webkit-box;-webkit-box-orient: vertical;-webkit-line-clamp: 1;overflow: hidden;">
                              </div>
                              <div style="font-size: 1.1vw;display: flex;justify-content: flex-start;">
                                <span>发布时间：{{ item.CreatedAt }}</span>
                                <span style="width: 1vw;"></span>
                                <span>阅读数量：{{ item.read_count }}</span>
                              </div>
                            </div>
                          </div>
                        </router-link>
                      </v-list-item>
                    </v-list>
                  </v-tab-item>
                  
                  <v-tab-item key="2">
                    <v-list lines="two">
                      <v-list-item v-for="item in categoriesList7" :key="item.ID">
                        <router-link :to="`/article/detail/${item.ID}`" class="article-link">
                          <div style="width: 100%;height: 10vh;display: flex;justify-content: space-around;margin: 1vh 0;">
                            <div
                              :style="{ width: '10%', height: '100%', backgroundSize: 'contain', backgroundImage: `url(${item.img})` }"
                              class="article-image"
                            >
                            </div>
                            <div style="width: 88%;height: 100%;">
                              <span style="font-size: 1.2vw;font-weight: 600;">{{ item.title }}</span>
                              <div v-html="item.content"
                                style="height: 3vh;margin: 0;color: rgb(100, 100, 100);font-size: 1.1vw;display: -webkit-box;-webkit-box-orient: vertical;-webkit-line-clamp: 1;overflow: hidden;">
                              </div>
                              <div style="font-size: 1.1vw;display: flex;justify-content: flex-start;">
                                <span>发布时间：{{ item.CreatedAt }}</span>
                                <span style="width: 1vw;"></span>
                                <span>阅读数量：{{ item.read_count }}</span>
                              </div>
                            </div>
                          </div>
                        </router-link>
                      </v-list-item>
                    </v-list>
                  </v-tab-item>
                  
                  <v-tab-item key="3">
                    <v-list lines="two">
                      <v-list-item v-for="item in categoriesList8" :key="item.ID">
                        <router-link :to="`/article/detail/${item.ID}`" class="article-link">
                          <div style="width: 100%;height: 10vh;display: flex;justify-content: space-around;margin: 1vh 0;">
                            <div
                              :style="{ width: '10%', height: '100%', backgroundSize: 'contain', backgroundImage: `url(${item.img})` }"
                              class="article-image"
                            >
                            </div>
                            <div style="width: 88%;height: 100%;">
                              <span style="font-size: 1.2vw;font-weight: 600;">{{ item.title }}</span>
                              <div v-html="item.content"
                                style="height: 3vh;margin: 0;color: rgb(100, 100, 100);font-size: 1.1vw;display: -webkit-box;-webkit-box-orient: vertical;-webkit-line-clamp: 1;overflow: hidden;">
                              </div>
                              <div style="font-size: 1.1vw;display: flex;justify-content: flex-start;">
                                <span>发布时间：{{ item.CreatedAt }}</span>
                                <span style="width: 1vw;"></span>
                                <span>阅读数量：{{ item.read_count }}</span>
                              </div>
                            </div>
                          </div>
                        </router-link>
                      </v-list-item>
                    </v-list>
                  </v-tab-item>
                  
                  <v-tab-item key="4">
                    <v-list lines="two">
                      <v-list-item v-for="item in categoriesList9" :key="item.ID">
                        <router-link :to="`/article/detail/${item.ID}`" class="article-link">
                          <div style="width: 100%;height: 10vh;display: flex;justify-content: space-around;margin: 1vh 0;">
                            <div
                              :style="{ width: '10%', height: '100%', backgroundSize: 'contain', backgroundImage: `url(${item.img})` }"
                              class="article-image"
                            >
                            </div>
                            <div style="width: 88%;height: 100%;">
                              <span style="font-size: 1.2vw;font-weight: 600;">{{ item.title }}</span>
                              <div v-html="item.content"
                                style="height: 3vh;margin: 0;color: rgb(100, 100, 100);font-size: 1.1vw;display: -webkit-box;-webkit-box-orient: vertical;-webkit-line-clamp: 1;overflow: hidden;">
                              </div>
                              <div style="font-size: 1.1vw;display: flex;justify-content: flex-start;">
                                <span>发布时间：{{ item.CreatedAt }}</span>
                                <span style="width: 1vw;"></span>
                                <span>阅读数量：{{ item.read_count }}</span>
                              </div>
                            </div>
                          </div>
                        </router-link>
                      </v-list-item>
                    </v-list>
                  </v-tab-item>
                </v-tabs-items>
              </v-col>
            </v-row>
          </v-col>
        </v-row>
      </v-container>
    </v-main>

    <Footer></Footer>
  </v-app>
</template>

<script>
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
      tab: null,
      categoriesList: [], //分类详细数据
      categoriesList7: [], //分类详细数据
      categoriesList8: [], //分类详细数据
      categoriesList9: [], //分类详细数据
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
    toPage(id) {
      this.$router.push(`/category/${id}`)
    },
    //获取列表数据
    async GetList() {
      try {
        const resp = await this.$http.get('/categories-with-articles/6')
        this.categoriesList = resp.data.data
      } catch (error) {
        console.error(error);
        this.$message.error('获取列表数据失败')
      }
    },
    async GetList7() {
      try {
        const resp = await this.$http.get('/categories-with-articles/7')
        this.categoriesList7 = resp.data.data
      } catch (error) {
        console.error(error);
        this.$message.error('获取列表数据失败')
      }
    },
    async GetList8() {
      try {
        const resp = await this.$http.get('/categories-with-articles/8')
        this.categoriesList8 = resp.data.data
      } catch (error) {
        console.error(error);
        this.$message.error('获取列表数据失败')
      }
    },
    async GetList9() {
      try {
        const resp = await this.$http.get('/categories-with-articles/9')
        this.categoriesList9 = resp.data.data
      } catch (error) {
        console.error(error);
        this.$message.error('获取列表数据失败')
      }
    }
  },
  created() {
    this.GetCateList() // 页面加载时获取分类数据
    this.GetList()
    this.GetList7()
    this.GetList8()
    this.GetList9()
  },
}
</script>

<style scoped>
.article-link {
  display: block;
  text-decoration: none;
  color: inherit;
  width: 100%;
}

.article-image {
  transition: transform 0.2s ease;
  cursor: pointer;
}

.article-link:hover .article-image {
  transform: scale(1.05);
}

.article-link:hover {
  background-color: rgba(0, 0, 0, 0.03);
}
</style>