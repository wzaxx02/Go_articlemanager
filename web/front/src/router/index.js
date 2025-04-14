import Vue from 'vue'
import VueRouter from 'vue-router'
import Router from 'vue-router'
const ArticleList = () =>
  import(/* webpackChunkName: "group-index" */ '../components/ArticleList.vue')
const Detail = () =>
  import(/* webpackChunkName: "group-detail" */ '../components/Details.vue')
const Category = () =>
  import(/* webpackChunkName: "group-category" */ '../components/CateList.vue')
const Search = () =>
  import(/* webpackChunkName: "group-search" */ '../components/Search.vue')
const Home = () => import('../views/Home.vue')
// const TopBar = () => import('../components/TopBar.vue')
const Layout = () => import('../views/Layout/LayoutView.vue')
const Favorite = () => import('../views/Favorite.vue')
Vue.use(VueRouter)
Vue.use(Router)

// //获取原型对象上的push函数
const originalPush = VueRouter.prototype.push
// //修改原型对象中的push方法
VueRouter.prototype.push = function push(location) {
  return originalPush.call(this, location).catch(err => err)
}

// const originalPush = VueRouter.prototype.push
// VueRouter.prototype.push = function push(location) {
//   return originalPush.call(this, location).catch((err) => err)
// }

const routes = [
  
  
  {
    path: '/',
    component: Layout,
    name: 'layout',
    redirect: '/home',
    children: [
      {
        path: 'home',
        name: 'home',
        component: Home,
      },
      {
        path: 'favorite',
        name: 'favorite',
        component: Favorite
      },
      {
        path: '/article/detail/:id',
        component: Detail,
        meta: { title: window.sessionStorage.getItem('title') },
        props: true
      },
      {
        path: 'category/:cid',
        // path: '/category',
        component: Category,
        meta: { title: '分类信息' },
        props: true
      },
      {
        path: 'search/:title',
        component: Search,
        meta: { title: '搜索结果' },
        props: true
      },
      { 
        path: 'articleList',
         component: ArticleList,
          meta: { title: '大作业' } 
      },
    ]
  },
  // {
  //   path: '/home',
  //   component: Home,
  // },
  // {
  //   path: '/login',
  //   component: TopBar,
  // },

]

// 创建路由实例
const router = new VueRouter({
  mode: 'history', // 使用 HTML5 History 模式
  base: process.env.BASE_URL,
  routes // 路由配置
})

router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = to.meta.title ? to.meta.title : '加载中'
  }
  next()
})

export default router
