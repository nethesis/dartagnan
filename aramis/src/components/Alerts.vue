<template>
  <div>
    <h2>Alerts</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <!-- filters -->
    <div v-if="!isLoading" class="row toolbar-pf filters-container">
      <div class="col-sm-12">
        <div class="toolbar-pf-actions">
          <div class="form-group toolbar-pf-filter">
            <label class="sr-only" for="filter">Name</label>
            <div class="input-group">
              <input v-model="filters.search" type="text" class="form-control" id="filter" :placeholder="$t('alerts.filter_by')">
              <div class="input-group-btn">
                <button class="btn btn-default" type="button">
                  <span class="fa fa-search"></span>
                </button>
              </div>
            </div>
          </div>
          <div class="form-group">
            <label>{{$t('alerts.priority')}}</label>
            <div class="dropdown btn-group renew-button">
              <button type="button" class="btn btn-default dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                {{$t('alerts.'+filters.currentPriority)}}
                <span class="caret"></span>
              </button>
              <ul class="dropdown-menu">
                <li @click="setFilterPriority(s)" v-for="(s, i) in filters.priorities" v-bind:key="i" :class="s == filters.currentPriority ? 'selected' : ''">
                  <a class="filter-option">{{$t('alerts.'+s)}}</a>
                </li>
              </ul>
            </div>
          </div>
        </div>
        <div class="row toolbar-pf-results">
          <div class="col-sm-12">
            <h5>{{filteredAlerts().length}} {{$t('alerts.results')}}</h5>
            <span v-if="filters.search.length > 0 || filters.currentPriority !== 'all'">
              <p>{{$t('alerts.active_filters')}}:</p>
              <ul class="list-inline">
                <li v-if="filters.search.length > 0">
                  <span class="label label-info">
                    {{$t('alerts.search')}}: {{filters.search}}
                    <a class="filter-option" @click="clearFilters('search')">
                      <span class="pficon pficon-close"></span>
                    </a>
                  </span>
                </li>
                <li v-if="filters.currentPriority !== 'all'">
                  <span class="label label-info">
                    {{$t('alerts.priority')}}: {{$t('alerts.'+filters.currentPriority)}}
                    <a class="filter-option" @click="clearFilters('priority')">
                      <span class="pficon pficon-close"></span>
                    </a>
                  </span>
                </li>
                <!-- <li v-if="filters.currentPlan.code !== 'all'">
                  <span class="label label-info">
                    {{$t('alerts.plan_type')}}: {{filters.currentPlan && filters.currentPlan.name || '-'}}
                    <a class="filter-option" @click="clearFilters('plan')">
                      <span class="pficon pficon-close"></span>
                    </a>
                  </span>
                </li> -->
              </ul>
              <p>
                <a class="filter-option" @click="clearFilters('all')">{{$t('alerts.clear_filters')}}</a>
              </p>
            </span>
          </div>
        </div>
      </div>
    </div>
    <!-- end filters -->
    <div v-if="!isLoading && alerts.length == 0" class="blank-slate-pf">
      <div class="blank-slate-pf-icon">
        <span class="pficon pficon-ok"></span>
      </div>
      <h1>
        {{ $t('alerts.no_alerts_title')}}
      </h1>
      <p>
        {{ $t('alerts.no_alerts_subtitle')}}
      </p>
    </div>

    <div v-if="!isLoading && alerts.length > 0" class="row row-cards-pf servers-container">
      <vue-good-table v-if="!isLoading" :perPage="5" :columns="columns" :rows="filteredAlerts()" :lineNumbers="false" :defaultSortBy="{field: 'priority', type: 'asc'}"
        :globalSearch="false" :paginate="true" styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText"
        :rowsPerPageText="tableLangsTexts.rowsPerPageText" :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder"
        :ofText="tableLangsTexts.ofText" class="container-fluid">
        <template slot="table-row" slot-scope="props">
          <td>
            <span :class="['fa fa-exclamation-triangle details-info', props.row.priority == 'HIGH'? 'red' : props.row.priority == 'AVERAGE' ? 'orange' : 'yellow' ]"
              data-toggle="tooltip" data-placement="left" :title="$t('alerts.'+props.row.priority)"></span>
            <strong>{{ props.row.namei18n }}</strong>
          </td>
          <td class="fancy">
            <strong @click="$parent.routeTo('servers/'+props.row.system.id)" class="soft click-hover">{{ props.row.system.hostname || '-'}}</strong>
          </td>
          <td class="fancy">{{ props.row.timestamp | formatDate}}</td>
          <td>
            <strong>{{ props.row.status }}</strong>
          </td>
          <td class="fancy">
            <span class="system-note">{{ props.row.note || '-' }}</span>
            <edit-note :obj="props.row"></edit-note>
          </td>
          <td>
            <delete-alert :obj="props.row" :update="listAlerts"></delete-alert>
          </td>
        </template>
      </vue-good-table>
    </div>
    <!-- /row -->
  </div>
</template>

<script>
  import LoginService from './../services/login';
  import StorageService from './../services/storage';
  import UtilService from './../services/util';
  import {
    setTimeout
  } from 'timers';

  import EditNote from './directives/EditNote.vue';
  import DeleteAlert from './directives/DeleteAlert.vue';

  export default {
    name: 'alerts',
    mixins: [LoginService, StorageService, UtilService],
    components: {
      editNote: EditNote,
      deleteAlert: DeleteAlert
    },
    updated: function () {
      $('[data-toggle="tooltip"]').tooltip()
    },
    data() {
      // get alerts list
      this.listAlerts()

      return {
        isLoading: true,
        filters: {
          search: '',
          currentPriority: this.get('alerts_filter_priority') || 'all',
          priorities: ['all', 'high', 'average', 'warning'],
        },
        alerts: [],
        columns: [{
            label: this.$i18n.t('alerts.alert_id'),
            field: 'namei18n',
            filterable: false,
          }, {
            label: this.$i18n.t('alerts.server'),
            field: 'namei18n',
            filterable: false,
          }, {
            label: this.$i18n.t('alerts.timestamp'),
            field: 'timestamp',
            filterable: false,
          },
          {
            label: this.$i18n.t('alerts.status'),
            field: 'status',
            filterable: false,
          },
          {
            label: this.$i18n.t('alerts.note'),
            field: 'note',
            filterable: false,
            sortable: false
          },
          {
            label: this.$i18n.t('alerts.action'),
            filterable: false,
            sortable: false
          }
        ],
        tableLangsTexts: this.tableLangs(),
      }
    },
    methods: {
      listAlerts() {
        this.isLoading = true
        this.$http.get(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/ui/alerts', {
          headers: {
            'Authorization': 'Bearer ' + this.get('access_token', false) || ''
          }
        }).then(function (success) {
          this.alerts = success.body
          this.isLoading = false
        }, function (error) {
          console.error(error)
          this.alerts = []
          this.isLoading = false
        });
      },
      clearFilters(filter) {
        if (filter == 'all') {
          this.filters.search = ''
          this.filters.currentPriority = 'all'
          this.set('alerts_filter_priority', 'all')
        }

        if (filter == 'search') {
          this.filters.search = ''
        }

        if (filter == 'priority') {
          this.filters.currentPriority = 'all'
          this.set('alerts_filter_priority', 'all')
        }
      },
      setFilterPriority(priority) {
        this.filters.currentPriority = priority
        this.set('alerts_filter_priority', priority)
        this.filteredAlerts()
      },
      filteredAlerts() {
        var filter = {}
        if (this.filters.currentPriority != 'all') {
          filter['priority'] = this.filters.currentPriority.toUpperCase()
        }
        var filtered = _.filter(this.alerts, filter);
        if (this.filters.search.length > 0) {
          var context = this
          filtered = _.filter(filtered, function (item) {
            return item.namei18n.toLowerCase().startsWith(context.filters.search.toLowerCase()) || item.alert_id.toLowerCase()
              .startsWith(context.filters.search.toLowerCase())
          })
        }
        return _.orderBy(filtered, ['priority'], 'asc')
      }
    }
  }

</script>

<style scoped>


</style>
