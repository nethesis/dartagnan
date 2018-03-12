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
                <div class="details-info">
                  <span>{{$t('servers.version')}}</span>
                  <span class="right">
                    <strong class="soft">{{server.inventory && server.inventory.os.release.full || '-'}}</strong>
                  </span>
                </div>
                <div class="details-info">
                  <span>{{$t('servers.public_ip')}}</span>
                  <span class="right">
                    <strong class="soft">{{server.info.public_ip || '-'}}</strong>
                  </span>
                </div>
                <div class="details-info">
                  <span>{{$t('servers.ns_lookup')}}</span>
                  <span class="right">
                    <strong class="soft">{{server.info.public_ip || '-'}}</strong>
                  </span>
                </div>
                <div class="details-info">
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
              {{$t('servers.subscription')}}
            </h2>
            <div class="card-pf-body">
              <div v-if="!isLoadingInfo">
                <div class="details-info">
                  <span>{{$t('servers.plan_type')}}</span>
                  <span class="right">
                    <strong class="soft">{{server.info.subscription.subscription_plan && server.info.subscription.subscription_plan.name || '-'}}</strong>
                  </span>
                </div>
                <div class="details-info">
                  <span>{{$t('servers.description')}}</span>
                  <span class="right">
                    <strong class="soft">{{server.info.subscription.subscription_plan && server.info.subscription.subscription_plan.description
                      || '-'}}</strong>
                  </span>
                </div>
                <div class="details-info">
                  <span>{{$t('servers.until')}}</span>
                  <span class="right">
                    <strong class="soft">{{server.info.subscription && server.info.subscription.valid_until || '-' | formatDate(false)}}</strong>
                  </span>
                </div>
                <div class="details-info">
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
                <renew-button v-if="server.info.id" v-bind:obj="server.info" :update="getServerInfo"></renew-button>
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
              <p v-if="!isLoadingAlerts && server.alerts.length == 0">{{$t('servers.no_alerts_found')}}</p>
              <span v-if="!isLoadingAlerts && server.alerts.length > 0">{{$t('servers.total')}}</span>
              <span class="right">
                <strong v-show="!isLoadingAlerts && server.alerts.length > 0" class="soft">{{server.alerts.length}}</strong>
              </span>
              <div v-if="isLoadingAlerts" class="spinner spinner-sm"></div>
            </div>
          </div>
        </div>
      </div>

      <div class="row row-cards-pf no-padding-top">
        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-6">
          <div class="card-pf card-pf-accented">
            <h2 class="card-pf-title">
              {{$t('servers.alerts_details')}}
            </h2>
            <div class="card-pf-body">
              <vue-good-table v-if="!isLoadingAlerts" :perPage="5" :columns="columns" :rows="server.alerts" :lineNumbers="false" :defaultSortBy="{field: 'priority', type: 'asc'}"
                :globalSearch="true" :paginate="true" styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText"
                :rowsPerPageText="tableLangsTexts.rowsPerPageText" :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder"
                :ofText="tableLangsTexts.ofText" class="container-fluid">
                <template slot="table-row" slot-scope="props">
                  <td>
                    <strong>{{ props.row.alert_id }}</strong>
                  </td>
                  <td class="fancy">{{ props.row.timestamp | formatDate}}</td>
                  <td>
                    <strong>{{ props.row.status }}</strong>
                  </td>
                  <td class="fancy">
                    <p>{{ props.row.note || '-' }}</p>
                  </td>
                  <td class="fancy">
                    <span :class="['fa fa-exclamation-triangle details-info', props.row.priority == 'HIGH'? 'red' : props.row.priority == 'AVERAGE' ? 'orange' : 'yellow' ]"
                      data-toggle="tooltip" data-placement="left" :title="$t('alert.'+props.row.priority)"></span>
                  </td>
                </template>
              </vue-good-table>
              <div v-if="isLoadingAlerts" class="spinner spinner-sm"></div>
            </div>
          </div>
        </div>

        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-6">
          <div class="card-pf card-pf-accented">
            <h2 class="card-pf-title">
              {{$t('servers.resources')}}
            </h2>
            <div class="card-pf-body resources-container">
              <div v-if="!isLoadingInventory && !server.inventory" class="alert alert-warning alert-dismissable">
                <span class="pficon pficon-warning-triangle-o"></span>
                <strong>{{$t('servers.inventory_not_available')}}</strong>. {{$t('servers.inventory_not_available_desc')}}.
              </div>
              <ul v-if="!isLoadingInventory && server.inventory" class="nav nav-tabs nav-tabs-pf" id="myTab" role="tablist">
                <li class="nav-item">
                  <a class="nav-link active" id="system-tab-parent" data-toggle="tab" href="#system-tab" role="tab" aria-controls="system"
                    aria-selected="true">{{$t('servers.system')}}</a>
                </li>
                <li @click="initMemoryCharts()" class="nav-item">
                  <a class="nav-link" id="memory-tab-parent" data-toggle="tab" href="#memory-tab" role="tab" aria-controls="memory" aria-selected="false">{{$t('servers.memory')}}</a>
                </li>
                <li @click="initStorageCharts()" class="nav-item">
                  <a class="nav-link" id="storage-tab-parent" data-toggle="tab" href="#storage-tab" role="tab" aria-controls="storage" aria-selected="false">{{$t('servers.storage')}}</a>
                </li>
                <li class="nav-item">
                  <a class="nav-link" id="network-tab-parent" data-toggle="tab" href="#network-tab" role="tab" aria-controls="network" aria-selected="false">{{$t('servers.network')}}</a>
                </li>
              </ul>
              <div v-if="!isLoadingInventory && server.inventory" class="tab-content" id="myServerResourceContent">
                <div class="tab-pane fade active" id="system-tab" role="tabpanel" aria-labelledby="system-tab">
                  <div class="container-fluid container-cards-pf">
                    <div class="row row-cards-pf">
                      <!-- OS -->
                      <div class="col-xs-12 col-sm-6 col-md-4 col-lg-6 resources-panel">
                        <div class="panel panel-default">
                          <div class="panel-heading">
                            <h3 class="panel-title">
                              <a class="icon-header-panel">
                                <span class="pficon pficon-running"></span>
                              </a>{{$t('servers.os')}}</h3>
                          </div>
                          <div class="panel-body">
                            <span class="details-info " ng-class="">{{server.inventory.os.name}}</span>
                            <div class="text-right">{{$t('servers.vendor')}}:
                              <b>
                                <span class="">{{server.inventory.os.family}}</span>
                              </b>
                            </div>
                            <div class="text-right">{{$t('servers.release')}}:
                              <b class="">{{server.inventory.os.release.full}}</b>
                            </div>
                          </div>
                        </div>
                      </div>
                      <!-- Kernel -->
                      <div class="col-xs-12 col-sm-6 col-md-4 col-lg-6 resources-panel">
                        <div class="panel panel-default">
                          <div class="panel-heading">
                            <h3 class="panel-title">
                              <a class="icon-header-panel">
                                <span class="fa fa-cogs"></span>
                              </a>{{$t('servers.kernel')}}</h3>
                          </div>
                          <div class="panel-body">
                            <span class="details-info " ng-class="">{{server.inventory.kernel}}</span>
                            <div class="text-right">{{$t('servers.release')}}:
                              <b>
                                <span class="">{{server.inventory.kernelrelease}}</span>
                              </b>
                            </div>
                            <div class="text-right">{{$t('servers.architecture')}}:
                              <b class="">{{server.inventory.os.architecture}}</b>
                            </div>
                          </div>
                        </div>
                      </div>
                      <!-- Machine -->
                      <div class="col-xs-12 col-sm-6 col-md-4 col-lg-6 resources-panel">
                        <div class="panel panel-default">
                          <div class="panel-heading">
                            <h3 class="panel-title">
                              <a class="icon-header-panel">
                                <span class="pficon pficon-screen"></span>
                              </a>{{$t('servers.machine')}}</h3>
                          </div>
                          <div class="panel-body">
                            <span class="details-info " ng-class="">{{server.inventory.virtual}}</span>
                            <div class="text-right ">{{$t('servers.uuid')}}:
                              <b>
                                <span class="">{{server.inventory.dmi.product.uuid}}</span>
                              </b>
                            </div>
                            <div class="text-right resource-info-details"></div>
                          </div>
                        </div>
                      </div>
                      <!-- CPU -->
                      <div class="col-xs-12 col-sm-6 col-md-4 col-lg-6 resources-panel">
                        <div class="panel panel-default">
                          <div class="panel-heading">
                            <h3 class="panel-title">
                              <a class="icon-header-panel">
                                <span class="pficon pficon-cpu"></span>
                              </a>{{$t('servers.cpu')}}</h3>
                          </div>
                          <div class="panel-body">
                            <span class="details-info" ng-class="">{{server.inventory.processors.count}} Core</span>
                            <div class="text-right">{{$t('servers.model')}}:
                              <b>
                                <span class="">{{server.inventory.processors.models[0]}}</span>
                              </b>
                            </div>
                            <div class="text-right resource-info-details"></div>
                          </div>
                        </div>
                      </div>
                      <!-- BIOS -->
                      <div class="col-xs-12 col-sm-6 col-md-4 col-lg-6 resources-panel">
                        <div class="panel panel-default">
                          <div class="panel-heading">
                            <h3 class="panel-title">
                              <a class="icon-header-panel">
                                <span class="pficon pficon-enterprise"></span>
                              </a>{{$t('servers.bios')}}</h3>
                          </div>
                          <div class="panel-body">
                            <span class="details-info " ng-class="">{{server.inventory.dmi.bios.version}}</span>
                            <div class="text-right">{{$t('servers.vendor')}}:
                              <b>
                                <span class="">{{server.inventory.dmi.bios.vendor}}</span>
                              </b>
                            </div>
                            <div class="text-right resource-info-details"></div>
                          </div>
                        </div>
                      </div>

                    </div>
                  </div>
                </div>
                <div class="tab-pane fade" id="memory-tab" role="tabpanel" aria-labelledby="memory-tab">
                  <div class="row row-cards-pf">
                    <div class="col-xs-12 col-sm-6 col-md-6 col-lg-6 resources-panel">
                      <div class="panel panel-default">
                        <div class="panel-heading">
                          <h3 class="panel-title">
                            <a class="icon-header-panel">
                              <span class="pficon pficon-memory"></span>
                            </a>{{$t('servers.ram')}}</h3>
                        </div>
                        <div class="panel-body">
                          <div id="ram-chart" class="text-center"></div>
                        </div>
                      </div>
                    </div>
                    <div class="col-xs-12 col-sm-6 col-md-6 col-lg-6 resources-panel">
                      <div class="panel panel-default">
                        <div class="panel-heading">
                          <h3 class="panel-title">
                            <a class="icon-header-panel">
                              <span class="fa fa-exchange"></span>
                            </a>{{$t('servers.swap')}}</h3>
                        </div>
                        <div class="panel-body">
                          <div id="swap-chart" class="text-center"></div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="tab-pane fade" id="storage-tab" role="tabpanel" aria-labelledby="storage-tab">
                  <div class="row row-cards-pf">
                    <div v-for="(k, v) in server.inventory.mountpoints" v-bind:key="v" class="col-xs-12 col-sm-6 col-md-6 col-lg-6 resources-panel">
                      <div class="panel panel-default">
                        <div class="panel-heading">
                          <h3 class="panel-title">
                            <a class="icon-header-panel">
                              <span class="fa fa-hdd-o"></span>
                            </a>{{v}}</h3>
                        </div>
                        <div class="panel-body">
                          <div :id="'mount-chart-'+parseMount(v)" class="text-center"></div>
                          <div class="text-right ">{{$t('servers.total')}}:
                            <b>
                              <span class="">{{server.inventory.mountpoints[v].size_bytes | byteFormat}}</span>
                            </b>
                          </div>
                          <div class="text-right ">{{$t('servers.device')}}:
                            <b>
                              <span class="">{{server.inventory.mountpoints[v].device}}</span>
                            </b>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="tab-pane fade" id="network-tab" role="tabpanel" aria-labelledby="network-tab">
                  <div class="row row-cards-pf">
                    <div v-if="e.type != 'xdsl-disabled'" v-for="(e, k) in server.inventory.esmithdb.networks" v-bind:key="k" class="col-xs-12 col-sm-6 col-md-4 col-lg-6 resources-panel">
                      <div class="panel panel-default">
                        <div class="panel-heading">
                          <h3 class="panel-title">
                            <a class="icon-header-panel">
                              <span class="pficon pficon-network"></span>
                            </a>{{$t('servers.interface')}} {{e.name}}</h3>
                        </div>
                        <div class="panel-body">
                          <div>
                            <span>{{$t('servers.ip')}}</span>
                            <span class="details-info" ng-class="">{{e.props.ipaddr || '-'}}</span>
                          </div>
                          <div>
                            <span>{{$t('servers.netmask')}}</span>
                            <span class="details-info" ng-class="">{{e.props.netmask || '-'}}</span>
                          </div>
                          <div>
                            <span>{{$t('servers.gateway')}}</span>
                            <span class="details-info" ng-class="">{{e.props.gateway || '-'}}</span>
                          </div>
                          <div class="text-right">{{$t('servers.type')}}:
                            <b>
                              <span class="">{{e.type}}</span>
                            </b>
                          </div>
                          <div class="text-right">{{$t('servers.role')}}:
                            <b>
                              <span :class="e.props.role">{{e.props.role}}</span>
                            </b>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
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
  import c3 from 'patternfly/node_modules/c3'

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
          isLoadingAlerts: true,
        },
        columns: [{
            label: this.$i18n.t('alert.alert_id'),
            field: 'alert_id',
            filterable: true,
          }, {
            label: this.$i18n.t('alert.timestamp'),
            field: 'timestamp',
            filterable: true,
          },
          {
            label: this.$i18n.t('alert.status'),
            field: 'status',
            filterable: true,
          },
          {
            label: this.$i18n.t('alert.note'),
            field: 'note',
            filterable: true,
          },
          {
            label: this.$i18n.t('alert.priority'),
            field: 'priority',
            filterable: true,
          },
        ],
        rows: [],
        tableLangsTexts: this.tableLangs(),
      }
    },
    methods: {
      isExpired(date) {
        return new Date().toISOString() > date
      },
      getServerInfo() {
        this.isLoadingInfo = true
        this.$http.get('https://' + this.$root.$options.api_host + '/api/ui/systems/' + this.$route.params.id, {
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
        this.$http.get('https://' + this.$root.$options.api_host + '/api/ui/inventories/' + this.$route.params.id, {
          headers: {
            'Authorization': 'Bearer ' + this.get('access_token', false) || ''
          }
        }).then(function (success) {
          this.server.inventory = success.body.data
          this.server.inventory.timestamp = success.body.timestamp
          this.isLoadingInventory = false
          // init tab
          setTimeout(function () {
            $('#system-tab-parent').click()
          }, 500)
        }, function (error) {
          console.error(error)
          this.isLoadingInventory = false
        });
      },
      getServerAlerts() {
        this.isLoadingAlerts = true
        this.$http.get('https://' + this.$root.$options.api_host + '/api/ui/alerts/' + this.$route.params.id, {
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
        this.$http.get('https://' + this.$root.$options.api_host + '/api/ui/heartbeats/' + this.$route.params.id, {
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
      },
      initMemoryCharts() {
        var c3ChartDefaults = $().c3ChartDefaults();
        var ramConfig = c3ChartDefaults.getDefaultDonutConfig('A');
        var swapConfig = c3ChartDefaults.getDefaultDonutConfig('A');
        ramConfig.bindto = '#ram-chart';
        swapConfig.bindto = '#swap-chart';
        ramConfig.data = {
          type: "donut",
          columns: [
            ["Used", this.server.inventory.memory.system.used_bytes],
            ["Available", this.server.inventory.memory.system.available_bytes],
          ],
          groups: [
            ["used", "available"]
          ],
          order: null
        };
        swapConfig.data = {
          type: "donut",
          columns: [
            ["Used", this.server.inventory.memory.swap.used_bytes],
            ["Available", this.server.inventory.memory.swap.available_bytes],
          ],
          groups: [
            ["used", "available"]
          ],
          order: null
        };
        ramConfig.size = {
          width: 180,
          height: 180
        };
        swapConfig.size = {
          width: 180,
          height: 180
        };

        ramConfig.tooltip = {
          contents: $().pfGetUtilizationDonutTooltipContentsFn('GB')
        };
        swapConfig.tooltip = {
          contents: $().pfGetUtilizationDonutTooltipContentsFn('GB')
        };

        c3.generate(ramConfig);
        c3.generate(swapConfig);
        $().pfSetDonutChartTitle("#ram-chart", this.$options.filters.byteFormat(this.server.inventory.memory.system.used_bytes),
          " Used");
        $().pfSetDonutChartTitle("#swap-chart", this.$options.filters.byteFormat(this.server.inventory.memory.swap.used_bytes),
          " Used");
      },
      parseMount(value) {
        return 'm' + value.substring(1);
      },
      initStorageCharts() {
        var c3ChartDefaults = $().c3ChartDefaults();
        var mountConfig = c3ChartDefaults.getDefaultDonutConfig('A');

        for (var m in this.server.inventory.mountpoints) {
          var mount = this.server.inventory.mountpoints[m]

          mountConfig.bindto = '#mount-chart-' + this.parseMount(m);
          mountConfig.data = {
            type: "donut",
            columns: [
              ["Used", mount.used_bytes],
              ["Available", mount.available_bytes],
            ],
            groups: [
              ["used", "available"]
            ],
            order: null
          };

          mountConfig.size = {
            width: 180,
            height: 180
          };


          mountConfig.tooltip = {
            contents: $().pfGetUtilizationDonutTooltipContentsFn('GB')
          };

          c3.generate(mountConfig);
          $().pfSetDonutChartTitle("#mount-chart-" + this.parseMount(m), this.$options.filters.byteFormat(mount
              .used_bytes),
            " Used");
        }

      }
    }
  }

</script>

<style scoped>


</style>
