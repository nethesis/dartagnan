<template>
  <div>
    <h2>Alerts</h2>

    <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
      <div class="card-pf card-pf-view">
        <div class="card-pf-body">
          <div class="card-pf-top-element">
            <span class="fa fa-wifi blue card-pf-icon-circle"></span>
          </div>
          <h2 class="card-pf-title text-center">NethSpot</h2>
          <div class="card-pf-items text-center">
            <div class="card-pf-item">
              <span class="fa fa-link"></span>
              <span class="card-pf-item-text"
                ><a href="https://my.nethspot.com" target="_blank"
                  >https://my.nethspot.com</a
                ></span
              >
            </div>
            <div class="card-pf-item">
              <button @click="getCredentials()" class="btn btn-primary">
                {{ $t("integrations.nethspot_get_credentials") }}
              </button>
              <div
                v-if="isLoading"
                class="spinner spinner-sm loader-cred"
              ></div>
            </div>
            <div
              v-if="isSuccess"
              class="alert alert-success alert-dismissable adjust-size"
            >
              <span class="pficon pficon-ok"></span>
              <ul class="list-integrations">
                <li>
                  <strong>Username: </strong>
                  <code>{{ nethspot.username }}</code>
                </li>
                <li>
                  <strong>Password: </strong>
                  <code>{{ nethspot.password }}</code>
                </li>
              </ul>
            </div>
            <div
              v-if="isSuccess"
              class="alert alert-warning alert-dismissable adjust-size"
            >
              <span class="pficon pficon-warning-triangle-o"></span>
              <strong>{{ $t("integrations.warning") }}!</strong>
              {{ $t("integrations.credentials_alert") }}
            </div>
            <div
              v-if="isError"
              class="alert alert-danger alert-dismissable adjust-size"
            >
              <span class="pficon pficon-error-circle-o"></span>
              <ul class="no-style">
                <li>
                  <strong>Error: </strong>
                  {{ nethspot.error }}.
                </li>
              </ul>
            </div>
          </div>
          <div class="card-pf-info text-center">
            <div class="text-center">
              {{ $t("integrations.nethspot_text") }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- /row -->
  </div>
</template>

<script>
import LoginService from "./../services/login";
import StorageService from "./../services/storage";
import UtilService from "./../services/util";

export default {
  name: "alerts",
  mixins: [LoginService, StorageService, UtilService],
  updated: function () {
    $('[data-toggle="tooltip"]').tooltip();
  },
  data() {
    return {
      isLoading: false,
      isError: false,
      isSuccess: false,
      nethspot: {
        username: "",
        password: "",
        error: "",
      },
    };
  },
  methods: {
    getCredentials() {
      this.isLoading = true;
      this.isError = false;
      this.isSuccess = false;
      this.$http
        .post(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/ui/integrations/nethspot",
          {},
          {
            headers: {
              Authorization: "Bearer " + this.get("access_token", false) || "",
            },
          }
        )
        .then(
          function (success) {
            this.nethspot.username = success.body.data.username;
            this.nethspot.password = success.body.data.password;
            this.nethspot.error = null;
            this.isLoading = false;
            this.isSuccess = true;
            this.isError = false;
          },
          function (error) {
            console.error(error);
            this.nethspot.username = null;
            this.nethspot.password = null;
            this.nethspot.error = error.body.data.error;
            this.isLoading = false;
            this.isSuccess = false;
            this.isError = true;
          }
        );
    },
  },
};
</script>

<style scoped>
.adjust-size {
  margin-top: 15px;
}

.list-integrations {
  list-style: none;
  text-align: left;
}

.no-style {
  list-style: none;
}

.loader-cred {
  float: right;
  margin-top: 4px;
  margin-left: 10px;
}
</style>
