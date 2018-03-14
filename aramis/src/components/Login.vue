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
                  <strong>800 €</strong> / year
                </h1>
              </p>
              <p>Stable Updates repository</p>
              <p>Community support</p>
              <p>Professional Support - 100€ per incident</p>
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
      <div class="right">
        <a class="unsplash" href="https://unsplash.com/@jplenio?utm_medium=referral&amp;utm_campaign=photographer-credit&amp;utm_content=creditBadge"
          target="_blank" rel="noopener noreferrer" title="Download free do whatever you want high-resolution photos from Johannes Plenio">
          <span style="display:inline-block;padding:2px 3px;">
            <svg xmlns="http://www.w3.org/2000/svg" style="height:12px;width:auto;position:relative;vertical-align:middle;top:-1px;fill:white;"
              viewBox="0 0 32 32">
              <title>unsplash-logo</title>
              <path d="M20.8 18.1c0 2.7-2.2 4.8-4.8 4.8s-4.8-2.1-4.8-4.8c0-2.7 2.2-4.8 4.8-4.8 2.7.1 4.8 2.2 4.8 4.8zm11.2-7.4v14.9c0 2.3-1.9 4.3-4.3 4.3h-23.4c-2.4 0-4.3-1.9-4.3-4.3v-15c0-2.3 1.9-4.3 4.3-4.3h3.7l.8-2.3c.4-1.1 1.7-2 2.9-2h8.6c1.2 0 2.5.9 2.9 2l.8 2.4h3.7c2.4 0 4.3 1.9 4.3 4.3zm-8.6 7.5c0-4.1-3.3-7.5-7.5-7.5-4.1 0-7.5 3.4-7.5 7.5s3.3 7.5 7.5 7.5c4.2-.1 7.5-3.4 7.5-7.5z"></path>
            </svg>
          </span>
          <span style="display:inline-block;padding:2px 3px;">Johannes Plenio</span>
        </a>
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
      position: absolute;
      right: 0;
      top: -10px;
    }
    #plans-table {
      margin-top: 70px;
      padding: 40px;
      margin-bottom: 25px;
    }
  }

  @media (max-width: 768px) {
    #plans-table {
      margin-top: 0px;
      padding: 40px;
      margin-bottom: 35px;
    }
    #badge {
      margin-bottom: 0px;

    }
  }

  .sub-desc {
    margin-top: -15px;
    font-size: 14px;
    color: #fff;
  }

  .logo-container-1 {
    margin: -20px;
    background: #82c1e8;
    margin-bottom: 10px;
    padding-bottom: 25px;
    color: white;
  }

  .logo-container-2 {
    margin: -20px;
    background: #212121;
    margin-bottom: 10px;
    padding-bottom: 25px;
    color: white;
  }

  .logo-container-3 {
    margin: -20px;
    background: #b71c1c;
    margin-bottom: 10px;
    padding-bottom: 25px;
    color: white;
  }

  .logo-container-4 {
    margin: -20px;
    background: #19425b;
    margin-bottom: 10px;
    padding-bottom: 25px;
    color: white;
  }

  .logo-plan-img {
    width: 120px;
    margin-bottom: -15px;
  }

  .disabled {
    color: #9c9c9c;
  }

  .login-cont {
    position: inherit;
  }

  .card-plan {
    margin-bottom: 10px;
  }

  .card-pf-body {
    min-height: 250px !important;
  }

  .login-big {
    width: 150px;
    height: 40px;
    font-size: 18px;
  }

</style>
