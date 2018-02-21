<template>
  <div class="spinner spinner-lg absolute-center"></div>
</template>

<script>
  import LoginService from './../services/login';
  import StorageService from './../services/storage';

  export default {
    name: 'callback',
    mixins: [LoginService, StorageService],
    data() {
      var context = this
      this.auth0Handler(function (path) {
        context.$router.push({
          path: path
        })
      })
      return {}
    }
  }

</script>

<style scoped>

</style>
