<template>
  <div>
    <h2>Server
      <strong class="soft">{{server.inventory && server.inventory.networking.fqdn || '-'}}</strong>
    </h2>
    <div class="container-fluid container-cards-pf">
      <div class="row row-cards-pf no-padding-top row-divider">

        <div class="col-xs-12 col-sm-6 col-md-6 col-lg-6">
          <div class="card-pf card-pf-accented">
            <h2 class="card-pf-title">
              {{$t('servers.info')}}
            </h2>
            <div class="card-pf-body">
              <div v-if="!isLoadingInventory">
                <div>
                  <span>{{$t('servers.version')}}</span>
                  <span class="right">
                    <strong class="soft">{{server.inventory && server.inventory.os.release.full || '-'}}</strong>
                  </span>
                </div>
                <div>
                  <span>{{$t('servers.public_ip')}}</span>
                  <span class="right">
                    <strong class="soft">{{server.info.public_ip || '-'}}</strong>
                  </span>
                </div>
                <div>
                  <span>{{$t('servers.ns_lookup')}}</span>
                  <span class="right">
                    <strong class="soft">{{server.info.public_ip || '-'}}</strong>
                  </span>
                </div>
                <div>
                  <span>{{$t('servers.uptime')}}</span>
                  <span class="right">
                    <strong class="soft">{{server.inventory && server.inventory.system_uptime.seconds | secondsInReadable}}</strong>
                  </span>
                </div>
              </div>
              <div v-if="isLoadingInventory" class="spinner spinner-sm"></div>
            </div>
            <div class="card-pf-footer">
              <div class="dropdown card-pf-time-frame-filter">
                <delete-server :obj="server.info" :redir="'servers'"></delete-server>
              </div>
              <p>
                <a class="card-pf-link-with-icon">
                </a>
              </p>
            </div>
          </div>
        </div>

        <div class="col-xs-12 col-sm-6 col-md-6 col-lg-6">
          <div class="card-pf card-pf-accented">
            <h2 class="card-pf-title">
              {{$t('servers.plan')}}
            </h2>
            <div class="card-pf-body">
              <div v-if="!isLoadingInfo">
                <div>
                  <span>{{$t('servers.plan_type')}}</span>
                  <span class="right">
                    <strong class="soft">{{server.info.subscription.subscription_plan && server.info.subscription.subscription_plan.name || '-'}}</strong>
                  </span>
                </div>
                <div>
                  <span>{{$t('servers.description')}}</span>
                  <span class="right">
                    <strong class="soft">{{server.info.subscription.subscription_plan && server.info.subscription.subscription_plan.description
                      || '-'}}</strong>
                  </span>
                </div>
                <div>
                  <span>{{$t('servers.until')}}</span>
                  <span class="right">
                    <strong class="soft">{{server.info.subscription && server.info.subscription.valid_until || '-' | formatDate}}</strong>
                  </span>
                </div>
                <div>
                  <span>{{$t('servers.status')}}</span>
                  <span class="right">
                    <span data-toggle="tooltip" data-placement="left" :title="$t('servers.license_valid')" v-if="!isExpired(server.info.subscription.valid_until)"
                      class="pficon pficon-ok"></span>
                    <span data-toggle="tooltip" data-placement="left" :title="$t('servers.expired')" v-if="isExpired(server.info.subscription.valid_until)"
                      class="pficon pficon-warning-triangle-o"></span>
                  </span>
                </div>
              </div>
              <div v-if="isLoadingInfo" class="spinner spinner-sm"></div>
            </div>
            <div class="card-pf-footer">
              <div class="dropdown card-pf-time-frame-filter">
                <renew-button v-if="server.info.id" v-bind:obj="server.info" :update="getServerInventory"></renew-button>
              </div>
              <p>
                <a class="card-pf-link-with-icon">
                </a>
              </p>
            </div>
          </div>
        </div>
      </div>

      <div class="row row-cards-pf no-padding-top row-divider">
        <div class="col-xs-12 col-sm-6 col-md-4 col-lg-4">
          <div class="card-pf card-pf-accented">
            <h2 class="card-pf-title">
              {{$t('servers.status')}}
              <span data-toggle="tooltip" data-placement="left" :title="$t('servers.active')" v-if="server.info.status == 'active'" class="pficon pficon-ok right"></span>
              <span data-toggle="tooltip" data-placement="left" :title="$t('servers.no_active')" v-if="server.info.status == 'no_active'"
                class="pficon pficon-error-circle-o right"></span>
              <span data-toggle="tooltip" data-placement="left" :title="$t('servers.no_comm')" v-if="server.info.status == 'no_comm'" class="pficon pficon-help right"></span>
            </h2>
            <div class="card-pf-body">
              <span>{{$t('servers.last_check')}}</span>
              <span class="right">
                <strong v-show="!isLoadingHeartbeat" class="soft">{{server.heartbeat | dateFromNow}}</strong>
                <div v-show="isLoadingHeartbeat" class="spinner spinner-sm"></div>
              </span>
            </div>
          </div>
        </div>

        <div class="col-xs-12 col-sm-6 col-md-4 col-lg-4">
          <div class="card-pf card-pf-accented">
            <h2 class="card-pf-title">
              {{$t('servers.inventory')}}
              <span v-if="!server.inventory" data-toggle="tooltip" data-placement="left" :title="$t('servers.inventory_not_available')"
                :class="['pficon', 'pficon-help', 'right']"></span>
              <span v-if="server.inventory" data-toggle="tooltip" data-placement="left" :title="$t('servers.inventory_available')" :class="['pficon', 'pficon-ok', 'right']"></span>
            </h2>
            <div class="card-pf-body">
              <span>{{$t('servers.last_update')}}</span>
              <span class="right">
                <strong v-show="!isLoadingInventory" class="soft">{{server.inventory && server.inventory.timestamp | dateFromNow}}</strong>
                <div v-show="isLoadingInventory" class="spinner spinner-sm"></div>
              </span>
            </div>
          </div>
        </div>

        <div class="col-xs-12 col-sm-12 col-md-4 col-lg-4">
          <div class="card-pf card-pf-accented">
            <h2 class="card-pf-title">
              {{$t('servers.alerts')}}
              <span :class="['pficon', server.alerts.length == 0 ? 'pficon-ok' : 'pficon-error-circle-o', 'right']"></span>
            </h2>
            <div class="card-pf-body">
              <p v-if="!isLoadingAlerts">{{$t('servers.no_alerts_found')}}</p>
              <div v-if="isLoadingAlerts" class="spinner spinner-sm"></div>
            </div>
          </div>
        </div>
      </div>

      <div class="row row-cards-pf no-padding-top">
        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-6">
          <div class="card-pf card-pf-accented">
            <h2 class="card-pf-title">
              {{$t('servers.resources')}}
            </h2>
            <div class="card-pf-body">
              <p v-if="!isLoadingInventory">[card contents]</p>
              <div v-if="isLoadingInventory" class="spinner spinner-sm"></div>
            </div>
            <div class="card-pf-footer">
              <div class="dropdown card-pf-time-frame-filter">
                <button type="button" class="btn btn-default dropdown-toggle" data-toggle="dropdown">
                  Last 30 Days
                  <span class="caret"></span>
                </button>
                <ul class="dropdown-menu dropdown-menu-right" role="menu">
                  <li class="selected">
                    <a href="#">Last 30 Days</a>
                  </li>
                  <li>
                    <a href="#">Last 60 Days</a>
                  </li>
                  <li>
                    <a href="#">Last 90 Days</a>
                  </li>
                </ul>
              </div>
              <p>
                <a href="#" class="card-pf-link-with-icon">
                  <span class="pficon pficon-add-circle-o"></span>Add New Cluster
                </a>
              </p>
            </div>
          </div>
        </div>

        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-6">
          <div class="card-pf card-pf-accented">
            <h2 class="card-pf-title">
              {{$t('servers.alerts_details')}}
            </h2>
            <div class="card-pf-body">
              <p v-if="!isLoadingAlerts">[card contents]</p>
              <div v-if="isLoadingAlerts" class="spinner spinner-sm"></div>
            </div>
            <div class="card-pf-footer">
              <div class="dropdown card-pf-time-frame-filter">
                <button type="button" class="btn btn-default dropdown-toggle" data-toggle="dropdown">
                  Last 30 Days
                  <span class="caret"></span>
                </button>
                <ul class="dropdown-menu dropdown-menu-right" role="menu">
                  <li class="selected">
                    <a href="#">Last 30 Days</a>
                  </li>
                  <li>
                    <a href="#">Last 60 Days</a>
                  </li>
                  <li>
                    <a href="#">Last 90 Days</a>
                  </li>
                </ul>
              </div>
              <p>
                <a href="#" class="card-pf-link-with-icon">
                  <span class="pficon pficon-add-circle-o"></span>Add New Cluster
                </a>
              </p>
            </div>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script>
  import LoginService from './../services/login';
  import StorageService from './../services/storage';
  import UtilService from './../services/util';
  import _ from 'lodash'

  import RenewButton from './directives/RenewButton.vue';
  import DeleteServer from './directives/DeleteServer.vue';
  import {
    setTimeout
  } from 'timers';

  export default {
    name: 'server',
    mixins: [LoginService, StorageService, UtilService],
    components: {
      renewButton: RenewButton,
      deleteServer: DeleteServer
    },
    created() {
      // get server info
      this.getServerInfo()

      // get server inventory
      this.getServerInventory()

      // get server heartbeats
      this.getServerHeartbeats()

      // get server alerts
      this.getServerAlerts()
    },
    data() {
      setTimeout(function () {
        $('[data-toggle="tooltip"]').tooltip()
      }, 500)
      return {
        server: {
          info: {},
          inventory: null,
          heartbeat: {},
          alerts: [],
          isLoadingInfo: true,
          isLoadingInventory: true,
          isLoadingHeartbeat: true,
          isLoadingAlerts: true
        }
      }
    },
    methods: {
      isExpired(date) {
        return new Date().toISOString() > date
      },
      getServerInfo() {
        this.isLoadingInfo = true
        this.$http.get('http://' + this.$root.$options.api_host + '/api/ui/systems/' + this.$route.params.id, {
          headers: {
            'Authorization': 'Bearer ' + this.get('access_token', false) || ''
          }
        }).then(function (success) {
          this.server.info = success.body
          this.isLoadingInfo = false
        }, function (error) {
          console.error(error)
          this.isLoadingInfo = false
        });
      },
      getServerInventory() {
        this.isLoadingInventory = true
        this.$http.get('http://' + this.$root.$options.api_host + '/api/ui/inventories/' + this.$route.params.id, {
          headers: {
            'Authorization': 'Bearer ' + this.get('access_token', false) || ''
          }
        }).then(function (success) {
          this.server.inventory = success.body.data
          this.server.inventory.timestamp = success.body.timestamp
          this.isLoadingInventory = false
        }, function (error) {
          console.error(error)
          this.isLoadingInventory = false
        });
      },
      getServerAlerts() {
        this.isLoadingAlerts = true
        this.$http.get('http://' + this.$root.$options.api_host + '/api/ui/alerts/' + this.$route.params.id, {
          headers: {
            'Authorization': 'Bearer ' + this.get('access_token', false) || ''
          }
        }).then(function (success) {
          this.server.alerts = success.body
          this.isLoadingAlerts = false
        }, function (error) {
          console.error(error)
          this.isLoadingAlerts = false
        });
      },
      getServerHeartbeats() {
        this.isLoadingHeartbeat = true
        this.$http.get('http://' + this.$root.$options.api_host + '/api/ui/heartbeats/' + this.$route.params.id, {
          headers: {
            'Authorization': 'Bearer ' + this.get('access_token', false) || ''
          }
        }).then(function (success) {
          this.server.heartbeat = success.body.timestamp
          this.isLoadingHeartbeat = false
        }, function (error) {
          console.error(error)
          this.isLoadingHeartbeat = false
        });
      }
    }
  }

</script>

<style scoped>


</style>
