module.exports = {
  // transpileDependencies: ['vuetify'],
  assetsDir: 'static',
  chainWebpack: config => {
    config.plugin('html').tap(args => {
      args[0].title = '大作业'
      return args
    })
  },
  productionSourceMap: false
}
