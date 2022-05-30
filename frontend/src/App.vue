<template>
  <img alt="Vue logo" src="./assets/logo.png">
  <EngagementComponent />
</template>

<script>
import EngagementComponent from './components/EngagementComponent.vue'

export default {
  name: 'App',
  components: {
    EngagementComponent
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
