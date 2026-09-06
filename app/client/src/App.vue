<template>
  <AppConsole />
</template>

<script>
import AppConsole from './components/AppConsole.vue'

export default {
  name: 'App',
  components: {
    AppConsole
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>