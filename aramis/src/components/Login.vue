<template>
  <div class="login-pf">
    <span id="badge">
      <img class="brand-logo" src="/static/logo.png" alt=" logo" />
    </span>
    <div v-if="sessionExpired" class="alert alert-warning alert-dismissable absolute-center-message">
      <span class="pficon pficon-warning-triangle-o"></span>
      <strong>{{$t('login.session_expired')}}</strong>. {{$t('login.session_expired_desc')}}.
    </div>
    <div id="plans-table">
      <div class="row row-cards-pf">

        <div class="col-xs-12 col-sm-6 col-md-3 col-lg-3 card-plan">
          <div class="card-pf text-center">
            <div class="logo-container-1">
              <img :src="getImage('crostino')" class="logo-plan-img">
              <h2 class="card-pf-title">
                <strong>Crostino</strong>
              </h2>
              <h4 class="sub-desc">Good Starter</h4>
            </div>
            <div class="card-pf-body">
              <p>
                <h1>
                  <strong>48 €</strong> / year
                </h1>
              </p>
              <p>Stable Updates repository</p>
              <p>Community support</p>
              <p class="disabled">Asset Portal</p>
              <p class="disabled">Monitoring Portal</p>
              <p class="disabled">Phone support</p>
            </div>
            <div class="card-pf-footer">
              <p>
                <a :href="infoUrl" target="_blank" class="btn btn-secondary btn-lg " type="button">{{$t('login.more_info')}}</a>
              </p>
            </div>
          </div>
        </div>
        <div class="col-xs-12 col-sm-6 col-md-3 col-lg-3 card-plan">
          <div class="card-pf text-center">
            <div class="logo-container-2">
              <img :src="getImage('lasagna')" class="logo-plan-img">
              <h2 class="card-pf-title">
                <strong>Lasagna</strong>
              </h2>
              <h4 class="sub-desc">Homemade first plate</h4>
            </div>
            <div class="card-pf-body">
              <p>
                <h1>
                  <strong>250 €</strong> / year
                </h1>
              </p>
              <p>Stable Updates repository</p>
              <p>Professional support via Customer Portal + SSH</p>
              <p>
                <strong class="soft">3</strong> support tickets/year</p>
              <p>Asset Portal</p>
              <p class="disabled">Monitoring Portal</p>
              <p class="disabled">Phone support</p>
            </div>
            <div class="card-pf-footer">
              <p>
                <a :href="infoUrl" target="_blank" class="btn btn-secondary btn-lg " type="button">{{$t('login.more_info')}}</a>
              </p>
            </div>
          </div>
        </div>
        <div class="col-xs-12 col-sm-6 col-md-3 col-lg-3 card-plan">
          <div class="card-pf text-center">
            <div class="logo-container-3">
              <img :src="getImage('fiorentina')" class="logo-plan-img">
              <h2 class="card-pf-title">
                <strong>Fiorentina</strong>
              </h2>
              <h4 class="sub-desc">The main Course</h4>
            </div>
            <div class="card-pf-body">
              <p>
                <h1>
                  <strong>450 €</strong> / year
                </h1>
              </p>
              <p>Stable Updates repository</p>
              <p>Professional support via Customer Portal + SSH</p>
              <p>
                <strong class="soft">6</strong> support tickets/year</p>
              <p>Asset Portal</p>
              <p>Monitoring Portal</p>
              <p class="disabled">Phone support</p>
            </div>
            <div class="card-pf-footer">
              <p>
                <a :href="infoUrl" target="_blank" class="btn btn-secondary btn-lg " type="button">{{$t('login.more_info')}}</a>
              </p>
            </div>
          </div>
        </div>
        <div class="col-xs-12 col-sm-6 col-md-3 col-lg-3 card-plan">
          <div class="card-pf text-center">
            <div class="logo-container-4">
              <img :src="getImage('pizza')" class="logo-plan-img">
              <h2 class="card-pf-title">
                <strong>Pizza</strong>
              </h2>
              <h4 class="sub-desc">What else?</h4>
            </div>
            <div class="card-pf-body">
              <p>
                <h1>
                  <strong>800 €</strong> / year
                </h1>
              </p>
              <p>Stable Updates repository</p>
              <p>Professional support via Customer Portal + SSH</p>
              <p>
                <strong class="soft">12</strong> support tickets/year</p>
              <p>Asset Portal</p>
              <p>Monitoring Portal</p>
              <p>Phone support</p>
            </div>
            <div class="card-pf-footer">
              <p>
                <a :href="infoUrl" target="_blank" class="btn btn-secondary btn-lg " type="button">{{$t('login.more_info')}}</a>
              </p>
            </div>
          </div>
        </div>
      </div>
      <!-- /row -->
    </div>
    <div class="container login-cont">
      <div class="row">
        <div class="col-sm-12 col-brand">
          <div id="brand">
            <h1>{{appName}}</h1>
          </div>
        </div>
        <div class="col-sm-12 col-md-12 col-lg-12 login">
          <p>
            <strong>{{ $t("login.welcome_title_pre") }} {{appName}} {{ $t("login.welcome_title_suf") }}.</strong>
            {{ $t("login.welcome_subtitle") }}
          </p>
          <button @click="doLogin()" class="btn btn-primary btn-lg login-big" type="button">{{ $t("login.login") }}</button>
        </div>
      </div>
    </div>
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
    name: 'login',
    mixins: [LoginService, StorageService, UtilService],
    data() {
      var sessionExpired = false
      if (this.$parent.action == 'sessionExpired') {
        sessionExpired = true
        this.delete('query_params')
      }
      return {
        sessionExpired: sessionExpired,
        appName: CONFIG.APP_NAME,
        infoUrl: CONFIG.INFO_URL
      }
    },
    methods: {
      doLogin() {
        this.auth0Login()
      },
      getImage(plan) {
        if (plan == 'fiorentina') {
          return require('./../assets/fiorentina-white.svg')
        }
        if (plan == 'pizza') {
          return require('./../assets/pizza-white.svg')
        }
        if (plan == 'crostino') {
          return require('./../assets/crostino-white.svg')
        }
        if (plan == 'lasagna') {
          return require('./../assets/lasagna-white.svg')
        }
      },
    }
  }

</script>

<style scoped>
  @media (min-width: 768px) {
    #badge {
      position: absolute !important;
      right: 0 !important;
      top: -10px !important;
    }
    #plans-table {
      margin-top: 70px !important;
      padding: 40px !important;
      margin-bottom: 25px !important;
    }
  }

  @media (max-width: 768px) {
    #plans-table {
      margin-top: 0px !important;
      padding: 40px !important;
      margin-bottom: 35px !important;
    }
    #badge {
      margin-bottom: 0px !important;

    }
  }

  .sub-desc {
    margin-top: -15px !important;
    font-size: 14px !important;
    color: #fff !important;
  }

  .logo-container-1 {
    margin: -20px !important;
    background: #82c1e8 !important;
    margin-bottom: 10px !important;
    padding-bottom: 25px !important;
    color: white !important;
  }

  .logo-container-2 {
    margin: -20px !important;
    background: #212121 !important;
    margin-bottom: 10px !important;
    padding-bottom: 25px !important;
    color: white !important;
  }

  .logo-container-3 {
    margin: -20px !important;
    background: #b71c1c !important;
    margin-bottom: 10px !important;
    padding-bottom: 25px !important;
    color: white !important;
  }

  .logo-container-4 {
    margin: -20px !important;
    background: #19425b !important;
    margin-bottom: 10px !important;
    padding-bottom: 25px !important;
    color: white !important;
  }

  .logo-plan-img {
    width: 120px;
    margin-bottom: -15px !important;
  }

  .disabled {
    display: none !important;
  }

  .login-cont {
    position: inherit !important;
  }

  .card-plan {
    margin-bottom: 10px !important;
  }

  .card-pf-body {
    min-height: 250px !important;
  }

  .login-big {
    width: 150px !important;
    height: 40px !important;
    font-size: 18px !important;
  }

</style>
