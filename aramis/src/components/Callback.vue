<template>
  <div class="spinner spinner-lg spinner-inverse absolute-center"></div>
</template>

<script>
  import LoginService from './../services/login';
  import StorageService from './../services/storage';

  export default {
    name: 'callback',
    mixins: [LoginService, StorageService],
    data() {
      var context = this
      this.$parent.showBody()
      this.auth0Handler(function (path) {
        context.$router.push({
          path: context.$parent.handleAction(context.$parent.action, path)
        })
      })
      return {}
    }
  }

</script>

<style scoped>


</style>
