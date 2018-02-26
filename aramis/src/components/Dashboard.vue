<template>
  <div>
    <h2>Dashboard</h2>
  </div>
</template>

<script>
  import LoginService from './../services/login';
  import StorageService from './../services/storage';
  import UtilService from './../services/util';
  import {
    setTimeout
  } from 'timers';

  export default {
    name: 'app',
    mixins: [LoginService, StorageService, UtilService],
    data() {
      // is logged
      this.$parent.isLogged = this.auth0CheckAuth()
      this.$parent.user = this.get('logged_user') || null

      this.initGraphics()

      return {
        user: this.$parent.user,
      }
    },
    methods: {
      getCurrentPath(route) {
        return this.$route.path.split('/')[1] === route
      },
      doLogin() {
        this.auth0Login()
      },
      initGraphics() {
        $('body').addClass('logged')
        $('body').removeClass('not-logged')
        setTimeout(function () {
          $().setupVerticalNavigation(true);
        }, 1000);
        this.showBody()
      },
      resetGraphics() {
        $('body').addClass('not-logged')
        $('body').removeClass('logged')
        window.location.reload()
      },
      showBody() {
        $('body').show()
        $('body').addClass('not-logged')
      }
    }
  }

</script>

<style scoped>

</style>
