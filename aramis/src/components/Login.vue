<template>
  <div class="login-pf">
    <div v-if="sessionExpired" class="alert alert-warning alert-dismissable absolute-center-message extra-padding">
      <span class="pficon pficon-warning-triangle-o"></span>
      <strong>{{$t('login.session_expired')}}</strong>. {{$t('login.session_expired_desc')}}.
    </div>
    <button @click="doLogin()" class="btn btn-primary btn-lg login-big absolute-center-top" type="button">{{ $t("login.login") }}</button>
    <iframe class="iframe-container" :src="iframeURL"></iframe>
  </div>
</template>

<script>
import LoginService from "./../services/login";
import StorageService from "./../services/storage";
import UtilService from "./../services/util";
import { setTimeout } from "timers";

export default {
  name: "login",
  mixins: [LoginService, StorageService, UtilService],
  data() {
    var sessionExpired = false;
    var isSessionExpired = this.get("sessionExpired") || false;
    if (isSessionExpired) {
      sessionExpired = true;
      this.delete("sessionExpired");
    }
    var isAutoLogin = this.get("autoLogin") || false;
    if (isAutoLogin) {
      this.doLogin();
      this.delete("autoLogin");
    }
    return {
      sessionExpired: sessionExpired,
      appName: CONFIG.APP_NAME,
      iframeURL: CONFIG.FRAME_URL
    };
  },
  methods: {
    doLogin() {
      this.auth0Login();
    },
    getImage(plan) {
      if (plan == "fiorentina") {
        return require("./../assets/fiorentina-white.svg");
      }
      if (plan == "pizza") {
        return require("./../assets/pizza-white.svg");
      }
      if (plan == "crostino") {
        return require("./../assets/crostino-white.svg");
      }
      if (plan == "lasagna") {
        return require("./../assets/lasagna-white.svg");
      }
    }
  }
};
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
  width: 180px !important;
  height: 40px !important;
  font-size: 20px !important;
}
</style>
