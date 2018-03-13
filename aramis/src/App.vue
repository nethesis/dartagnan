<template>
  <div>
    <div v-if="isLogged" class="cards-pf">
      <!-- top bar -->
      <nav class="navbar navbar-pf-vertical">
        <div class="navbar-header">
          <button type="button" class="navbar-toggle">
            <span class="sr-only"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          <router-link to="/dashboard" class="navbar-brand">
            <img class="navbar-brand-icon" src="./assets/logo.png" alt="" />
            <p class="navbar-brand-name">Aramis</p>
          </router-link>
        </div>
        <nav class="collapse navbar-collapse">
          <ul class="nav navbar-nav navbar-right navbar-iconic navbar-utility">

            <li class="dropdown">
              <a href="#" class="dropdown-toggle nav-item-iconic" id="dropdownMenu2" data-toggle="dropdown" aria-haspopup="true" aria-expanded="true">
                <img class="profile-picture-mini" :src="user.picture" />
                <p class="login-main-name">{{ user.name }}</p>
                <p class="login-main-type">
                  {{ user.email }}
                  <span class="caret"></span>
                </p>
              </a>
              <ul class="dropdown-menu" aria-labelledby="dropdownMenu2">
                <li>
                  <router-link to="/profile">{{ $t("menu.profile") }}</router-link>
                </li>
                <li>
                  <a v-on:click="doLogout()" href="#">Logout</a>
                </li>
              </ul>
            </li>
            <li class="dropdown">
              <a href="#" class="dropdown-toggle nav-item-iconic" id="dropdownMenu1">
                <span :title="$t('menu.help')" class="fa pficon-help" data-toggle="modal" data-target="#about-modal"></span>
              </a>
            </li>
          </ul>
        </nav>

      </nav>
      <!-- end top bar -->

      <!-- left menu -->
      <div class="nav-pf-vertical nav-pf-vertical-with-sub-menus nav-pf-persistent-secondary">
        <ul class="list-group">

          <!-- <li v-bind:class="[getCurrentPath('dashboard') ? 'active' : '', 'list-group-item']">
            <router-link to="/dashboard">
              <span class="fa fa-dashboard"></span>
              <span class="list-group-item-value">{{ $t('menu.dashboard') }}</span>
            </router-link>
          </li> -->

          <li v-bind:class="[getCurrentPath('servers') ? 'active' : '', 'list-group-item']">
            <router-link to="/servers">
              <span class="pficon pficon-server"></span>
              <span class="list-group-item-value">{{ $t("menu.servers") }}</span>

            </router-link>
          </li>

          <li v-bind:class="[getCurrentPath('alerts') ? 'active' : '', 'list-group-item']">
            <router-link to="/alerts">
              <span class="pficon pficon-warning-triangle-o"></span>
              <span class="list-group-item-value">{{ $t("menu.alerts") }}</span>

            </router-link>
          </li>

          <li class="list-group-item secondary-nav-item-pf mobile-nav-item-pf visible-xs-block" data-target="#user-secondary">
            <a href="#">
              <span class="pficon pficon-user" data-toggle="tooltip" title="" data-original-title="User"></span>
              <span class="list-group-item-value">{{ user.name }}</span>
            </a>
            <div id="user-secondary" class="nav-pf-secondary-nav">
              <div class="nav-item-pf-header">
                <a href="#" class="secondary-collapse-toggle-pf" data-toggle="collapse-secondary-nav"></a>
                <span>{{ user.name }}</span>
              </div>

              <ul class="list-group">
                <li class="list-group-item">
                  <router-link to="/profile">
                    <span class="list-group-item-value">{{ $t("menu.profile") }}</span>
                  </router-link>
                </li>

                <li class="list-group-item">
                  <a v-on:click="doLogout()" href="#">
                    <span class="list-group-item-value">{{ $t("menu.logout") }}</span>
                  </a>
                </li>
              </ul>
            </div>
          </li>

        </ul>

      </div>
      <!-- end left menu -->
      <!-- main view -->
      <div class="container-fluid container-cards-pf container-pf-nav-pf-vertical nav-pf-persistent-secondary">
        <router-view></router-view>
      </div>
      <!-- end main view -->
    </div>
    <div v-if="!isLogged">
      <router-view></router-view>
    </div>

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
            <h1>Dartagnan</h1>
            <div class="product-versions-pf">
              <ul class="list-unstyled">
                <li>
                  <strong>Dartagnan on</strong>
                  <a target="blank" href="https://github.com/nethesis/dartagna">GitHub</a>
                </li>
                <li>
                  <strong>Docs</strong>
                  <a target="blank" href="https://nethesis.github.io/dartagnan/">Link</a>
                </li>
              </ul>
            </div>
          </div>
          <div class="modal-footer">
            <img class="about-logo" src="./assets/logo.png" alt="Patternfly Symbol">
          </div>
        </div>
      </div>
    </div>
    <!-- end modals -->
  </div>
</template>

<script>
  import LoginService from './services/login';
  import StorageService from './services/storage';
  import UtilService from './services/util';
  import {
    setTimeout
  } from 'timers';
  import {
    error
  } from 'util';

  export default {
    name: 'app',
    mixins: [LoginService, StorageService, UtilService],
    data() {
      // is logged
      var isLogged = this.auth0CheckAuth()
      var user = this.get('logged_user') || null

      // save route query params
      if (Object.keys(this.get('query_params') || {}).length == 0) {
        if (Object.keys(this.$route.query).length > 0) {
          this.set('query_params', this.$route.query)
        }
      }

      // get action if exists
      var action = this.get('query_params') && this.get('query_params').action || null

      if (this.$route.path !== '/callback') {
        if (isLogged) {
          this.initGraphics()
          // handle action
          this.$router.push({
            path: this.handleAction(action)
          })
        } else {
          this.showBody()
          this.$router.push({
            path: '/login'
          })
        }
      }

      return {
        user: user,
        isLogged: isLogged,
        action: action
      }
    },
    methods: {
      handleAction(action, path) {
        switch (action) {
          case 'newServer':
            this.isLogged = this.auth0CheckAuth()
            this.initGraphics()
            return 'servers'

          default:
            return this.$route.path == '/' ? 'dashboard' : path || this.$route.path
        }
      },
      getCurrentPath(route) {
        return this.$route.path.split('/')[1] === route
      },
      doLogout() {
        this.$router.push('/login')
        this.auth0Logout()
        this.isLogged = false
        this.resetGraphics()
      },
      initGraphics() {
        $('body').addClass('logged')
        $('body').removeClass('not-logged')
        this.showBody()
        setTimeout(function () {
          $().setupVerticalNavigation(true);
        }, 0);
      },
      resetGraphics() {
        $('body').addClass('not-logged')
        $('body').removeClass('logged')
        window.location.reload(true)
      },
      showBody() {
        $('body').show()
        $('body').addClass('not-logged')
      },
      routeTo(route) {
        this.$router.push('/' + route)
      }
    }
  }

</script>

<style src="./styles/main.css">


</style>
