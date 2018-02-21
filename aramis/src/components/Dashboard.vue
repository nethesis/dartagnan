<template>
  <div>

    <h2>Dashboard</h2>

    <!-- modals -->
    <div class="modal fade" id="about-modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content about-modal-pf">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
          </div>
          <div class="modal-body">
            <h1>Aramis</h1>
            <div class="product-versions-pf">
              <ul class="list-unstyled">
                <li>
                  <strong>Athos</strong>
                  <a target="blank" href="https://github.com/nethesis/dartagnan/tree/master/athos">GitHub</a>
                </li>
                <li>
                  <strong>Aramis</strong>
                  <a target="blank" href="https://github.com/nethesis/dartagnan/tree/master/aramis">GitHub</a>
                </li>
              </ul>
            </div>
          </div>
          <div class="modal-footer">
            <img class="about-logo" src="./../assets/logo.png" alt="Patternfly Symbol">
          </div>
        </div>
      </div>
    </div>
    <!-- end modals -->
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
