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
                  {{ user.nickname }}
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
              <a href="#" class="dropdown-toggle nav-item-iconic" id="dropdownMenu1" data-toggle="dropdown" aria-haspopup="true" aria-expanded="true">
                <span :title="$t('menu.help')" class="fa pficon-help"></span>
                <span class="caret"></span>
              </a>
              <ul class="dropdown-menu" aria-labelledby="dropdownMenu1">
                <li>
                  <a target="blank" href="https://github.com/nethesis/dartagnan">{{ $t("menu.help") }}</a>
                </li>
                <li>
                  <a href="#" data-toggle="modal" data-target="#about-modal">{{ $t("menu.about") }}</a>
                </li>
              </ul>
            </li>
          </ul>
        </nav>

      </nav>
      <!-- end top bar -->

      <!-- left menu -->
      <div class="nav-pf-vertical nav-pf-vertical-with-sub-menus nav-pf-persistent-secondary">
        <ul class="list-group">

          <li v-bind:class="[getCurrentPath('dashboard') ? 'active' : '', 'list-group-item']">
            <router-link to="/dashboard">
              <span class="fa fa-dashboard"></span>
              <span class="list-group-item-value">{{ $t('menu.dashboard') }}</span>
            </router-link>
          </li>

          <li v-bind:class="[getCurrentPath('alerts') ? 'active' : '', 'list-group-item']">
            <router-link to="/alerts">
              <span class="pficon pficon-warning-triangle-o"></span>
              <span class="list-group-item-value">{{ $t("menu.alerts") }}</span>

            </router-link>
          </li>

          <li v-bind:class="[getCurrentPath('servers') ? 'active' : '', 'list-group-item']">
            <router-link to="/servers">
              <span class="pficon pficon-server"></span>
              <span class="list-group-item-value">{{ $t("menu.servers") }}</span>

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
          <li class="list-group-item secondary-nav-item-pf mobile-nav-item-pf visible-xs-block" data-target="#help-secondary">
            <a href="#">
              <span class="pficon pficon-help" data-toggle="tooltip" title="" data-original-title="Help"></span>
              <span class="list-group-item-value">{{ $t("menu.help") }}</span>
            </a>
            <div id="help-secondary" class="nav-pf-secondary-nav">
              <div class="nav-item-pf-header">
                <a href="#" class="secondary-collapse-toggle-pf" data-toggle="collapse-secondary-nav"></a>
                <span>{{ $t("menu.help") }}</span>
              </div>
              <ul class="list-group">
                <li class="list-group-item">
                  <a target="blank" href="https://github.com/nethesis/dartagnan">
                    <span class="list-group-item-value">{{ $t("menu.help") }}</span>
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
  </div>
</template>

<script>
  import LoginService from './services/login';
  import StorageService from './services/storage';
  import UtilService from './services/util';
  import {
    setTimeout
  } from 'timers';

  export default {
    name: 'app',
    mixins: [LoginService, StorageService, UtilService],
    data() {
      // is logged
      var isLogged = this.auth0CheckAuth()
      var user = this.get('logged_user') || null

      if (this.$route.path == '/callback') {

      } else {
        if (isLogged) {
          this.initGraphics()
          this.$router.push({
            path: this.$route.path
          })
        } else {
          this.showBody()
          this.$router.push({
            path: 'login'
          })
        }
      }

      return {
        user: user,
        isLogged: isLogged,
      }
    },
    methods: {
      getCurrentPath(route) {
        return this.$route.path.split('/')[1] === route
      },
      doLogout() {
        var context = this
        this.auth0Logout(function () {
          context.isLogged = false
          context.resetGraphics()
        })
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
      }
    }
  }

</script>

<style src="./styles/main.css">


</style>
