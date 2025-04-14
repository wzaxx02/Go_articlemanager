<template>
    <div>
        <div style="height: 9vh;">
            <TopBar></TopBar>
        </div>
        <div style="margin-top: 1.5vh;">
          <router-view></router-view>
        </div>
        <div>
        </div>
    </div>
</template>

<script>
import TopBar from '../../components/TopBar.vue'
// import { RouterView } from 'vue-router'

export default {
  components: { 
    TopBar,
  }
}
</script>

<style></style>