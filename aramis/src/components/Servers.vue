<template>
  <div>
    <h2>Servers</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <button
      v-if="!isLoading && servers.length > 0"
      @click="counters.ns8 >= counters.limit ? undefined : addServer('ns8')"
      :class="[
        'btn btn-primary btn-lg create-server-ns8',
        counters.ns8 >= counters.limit ? 'disabled' : '',
      ]"
    >
      {{ $t("servers.create_server_ns8") }}
    </button>
    <button
      v-if="!isLoading && servers.length > 0"
      @click="counters.nsec >= counters.limit ? undefined : addServer('nsec')"
      :class="[
        'btn btn-primary btn-lg create-server-nsec',
        counters.nsec >= counters.limit ? 'disabled' : '',
      ]"
    >
      {{ $t("servers.create_server_nsec") }}
    </button>

    <!-- filters -->
    <div v-if="!isLoading" class="row toolbar-pf filters-container">
      <div class="col-sm-12">
        <div class="toolbar-pf-actions">
          <div class="form-group toolbar-pf-filter">
            <label class="sr-only" for="filter">Name</label>
            <div class="input-group">
              <input
                v-model="filters.search"
                type="text"
                class="form-control"
                id="filter"
                :placeholder="$t('servers.filter_by')"
              />
              <div class="input-group-btn">
                <button class="btn btn-default" type="button">
                  <span class="fa fa-search"></span>
                </button>
              </div>
            </div>
          </div>
          <div class="form-group">
            <label>{{ $t("servers.status") }}</label>
            <div class="dropdown btn-group renew-button">
              <button
                type="button"
                class="btn btn-default dropdown-toggle"
                data-toggle="dropdown"
                aria-haspopup="true"
                aria-expanded="false"
              >
                {{ $t("servers." + filters.currentStatus) }}
                <span class="caret"></span>
              </button>
              <ul class="dropdown-menu">
                <li
                  @click="setFilterStatus(s)"
                  v-for="(s, i) in filters.statuses"
                  v-bind:key="i"
                  :class="s == filters.currentStatus ? 'selected' : ''"
                >
                  <a class="filter-option">{{ $t("servers." + s) }}</a>
                </li>
              </ul>
            </div>
          </div>
          <div class="form-group">
            <label>{{ $t("servers.plan_type") }}</label>
            <div class="dropdown btn-group renew-button">
              <button
                type="button"
                class="btn btn-default dropdown-toggle"
                data-toggle="dropdown"
                aria-haspopup="true"
                aria-expanded="false"
              >
                {{ (filters.currentPlan && filters.currentPlan.name) || "-" }}
                <span class="caret"></span>
              </button>
              <ul class="dropdown-menu">
                <li
                  @click="setFilterPlan(l)"
                  v-for="(l, i) in filters.plans"
                  v-bind:key="i"
                  :class="
                    l && l.code == filters.currentPlan.code ? 'selected' : ''
                  "
                >
                  <a class="filter-option">{{
                    l && l.code == "all"
                      ? $t("servers.all")
                      : (l && l.name) || "-"
                  }}</a>
                </li>
              </ul>
            </div>
          </div>
        </div>
        <div class="row toolbar-pf-results">
          <div class="col-sm-12">
            <h5>{{ filteredServers().length }} {{ $t("servers.results") }}</h5>
            <span
              v-if="
                filters.search.length > 0 ||
                filters.currentStatus !== 'all' ||
                filters.currentPlan.code !== 'all'
              "
            >
              <p>{{ $t("servers.active_filters") }}:</p>
              <ul class="list-inline">
                <li v-if="filters.search.length > 0">
                  <span class="label label-info">
                    {{ $t("servers.search") }}: {{ filters.search }}
                    <a class="filter-option" @click="clearFilters('search')">
                      <span class="pficon pficon-close"></span>
                    </a>
                  </span>
                </li>
                <li v-if="filters.currentStatus !== 'all'">
                  <span class="label label-info">
                    {{ $t("servers.status") }}:
                    {{ $t("servers." + filters.currentStatus) }}
                    <a class="filter-option" @click="clearFilters('status')">
                      <span class="pficon pficon-close"></span>
                    </a>
                  </span>
                </li>
                <li v-if="filters.currentPlan.code !== 'all'">
                  <span class="label label-info">
                    {{ $t("servers.plan_type") }}:
                    {{
                      (filters.currentPlan && filters.currentPlan.name) || "-"
                    }}
                    <a class="filter-option" @click="clearFilters('plan')">
                      <span class="pficon pficon-close"></span>
                    </a>
                  </span>
                </li>
              </ul>
              <p>
                <a class="filter-option" @click="clearFilters('all')">{{
                  $t("servers.clear_filters")
                }}</a>
              </p>
            </span>
          </div>
        </div>
      </div>
    </div>
    <!-- end filters -->

    <div
      v-if="servers.length > 0"
      class="alert alert-warning alert-dismissable info-trial"
    >
      <span class="pficon pficon-warning-triangle-o"></span>
      <strong>{{ $t("servers.trials_limit") }}:</strong>
      {{
        $t("servers.trials_limit_text", {
          ns8: counters.limit - counters.ns8,
          nsec: counters.limit - counters.nsec,
        })
      }}.
    </div>

    <div v-if="!isLoading && servers.length == 0" class="blank-slate-pf">
      <div class="blank-slate-pf-icon">
        <span class="pficon pficon pficon-server"></span>
      </div>
      <h1>
        {{ $t("servers.no_servers_title") }}
      </h1>
      <p>
        {{ $t("servers.no_servers_subtitle") }}
      </p>
      <div class="blank-slate-pf-main-action">
        <button @click="addServer('ns8')" class="btn btn-primary btn-lg">
          {{ $t("servers.create_server_ns8") }}
        </button>
        <button @click="addServer('nsec')" class="btn btn-primary btn-lg">
          {{ $t("servers.create_server_nsec") }}
        </button>
      </div>
    </div>

    <div
      v-if="!isLoading && servers.length > 0"
      class="row row-cards-pf servers-container"
    >
      <div
        v-for="s in filteredServers()"
        v-bind:key="s.id"
        class="col-xs-12 col-sm-12 col-md-6 col-lg-3"
      >
        <div
          :class="[
            isExpired(s.subscription.valid_until) ? 'disabled-top' : '',
            'card-pf card-pf-view card-pf-accented',
          ]"
        >
          <div class="card-pf-body adjust-height">
            <span
              v-if="s.alerts > 0"
              class="fa fa-exclamation-triangle fa-big orange pull-right fa-2x"
              data-toggle="tooltip"
              data-placement="left"
              :title="$t('servers.alerts') + ': ' + s.alerts"
            ></span>
            <div
              @click="$parent.routeTo('servers/' + s.id)"
              class="card-pf-top-element click-hover"
            >
              <img
                data-toggle="tooltip"
                data-placement="top"
                :title="
                  isUnsupported(s.subscription.subscription_plan.base_code)
                    ? s.subscription.subscription_plan.name + ' (EOL)'
                    : s.subscription.subscription_plan.name
                "
                :src="getImage(s)"
                :class="[
                  isExpired(s.subscription.valid_until) ? 'expired-svg' : '',
                  'pficon pficon-server card-pf-icon-circle adjust-icon-size',
                ]"
              />
            </div>
            <h2
              @click="$parent.routeTo('servers/' + s.id)"
              class="card-pf-title text-center click-hover"
            >
              {{ s.hostname || "-" }}
            </h2>
            <div class="card-pf-items text-center">
              <div class="card-pf-item">
                <span class="pficon pficon-screen"></span>
                <span
                  data-toggle="tooltip"
                  data-placement="top"
                  :title="$t('servers.public_ip')"
                  class="card-pf-item-text"
                  >{{ s.public_ip || "-" }}</span
                >
              </div>
              <div class="card-pf-item">
                <span
                  data-toggle="tooltip"
                  data-placement="right"
                  :title="$t('servers.active')"
                  v-if="s.status == 'active'"
                  class="pficon pficon-ok"
                ></span>
                <span
                  data-toggle="tooltip"
                  data-placement="right"
                  :title="$t('servers.no_active')"
                  v-if="s.status == 'no_active'"
                  class="pficon pficon-error-circle-o"
                ></span>
                <span
                  data-toggle="tooltip"
                  data-placement="right"
                  :title="$t('servers.no_comm')"
                  v-if="s.status == 'no_comm'"
                  class="pficon pficon-help"
                ></span>
              </div>
            </div>
            <p class="card-pf-info text-center">
              <strong>{{ $t("servers.created") }}</strong
              >{{ s.created | formatDate }}
            </p>
            <div class="divider"></div>
            <div class="card-pf-items text-center">
              <div class="card-pf-item">
                <span
                  :class="[
                    isExpired(s.subscription.valid_until)
                      ? 'fa fa-ban'
                      : 'fa fa-star',
                  ]"
                ></span>
                <span
                  data-toggle="tooltip"
                  data-placement="top"
                  :title="s.subscription.subscription_plan.description"
                  class="card-pf-item-text"
                  >{{
                    (s &&
                      s.subscription &&
                      s.subscription.subscription_plan &&
                      s.subscription.subscription_plan.name) ||
                    "-"
                  }}</span
                >
              </div>
              <div class="card-pf-item">
                <span
                  data-toggle="tooltip"
                  data-placement="right"
                  :title="$t('servers.license_valid')"
                  v-if="!isExpired(s.subscription.valid_until)"
                  class="pficon pficon-ok"
                ></span>
                <span
                  data-toggle="tooltip"
                  data-placement="right"
                  :title="$t('servers.expired')"
                  v-if="isExpired(s.subscription.valid_until)"
                  class="pficon pficon-warning-triangle-o"
                ></span>
              </div>
            </div>
            <p
              v-if="
                extractServices(s.subscription.subscription_plan.code).length >
                0
              "
              class="card-pf-info text-center"
            >
              <strong>{{ $t("payment.services") }}</strong>
              <code>{{
                extractServices(s.subscription.subscription_plan.code)
              }}</code>
            </p>
            <p class="card-pf-info text-center">
              <strong>{{ $t("servers.until") }}</strong>
              <span v-if="!isExpired(s.subscription.valid_until)">{{
                s.subscription.valid_until | formatDate(false)
              }}</span>
              <span class="gray" v-if="isExpired(s.subscription.valid_until)">
                <strong>{{ $t("servers.expired") }}</strong>
              </span>
            </p>
          </div>
          <div class="card-pf-footer">
            <div class="card-pf-time-frame-filter">
              <renew-button :obj="s" :update="listServers"></renew-button>
            </div>
            <p>
              <delete-server :obj="s" :update="listServers"></delete-server>
            </p>
          </div>
        </div>
      </div>

      <!-- <router-link to="/dashboard"> -->
    </div>

    <div
      class="modal fade"
      id="newServerModal"
      tabindex="-1"
      role="dialog"
      aria-labelledby="myModalLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button
              type="button"
              class="close"
              data-dismiss="modal"
              aria-hidden="true"
            >
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="myModalLabel">
              {{ $t("servers.new_server") }}
            </h4>
          </div>
          <div class="modal-body">
            <div class="card-pf-view">
              <div>
                <div class="card-pf-top-element">
                  <span class="pficon pficon-server card-pf-icon-circle"></span>
                </div>
                <h2 class="card-pf-title text-center">
                  {{ $t("servers.server_created") }}
                </h2>
                <div class="card-pf-items text-center">
                  {{ $t("servers.server_id") }}
                  <div class="card-pf-item">
                    <span class="card-pf-item-text"
                      ><code>{{ newServer.secret }}</code></span
                    >
                  </div>
                  <div class="servers-container">
                    <button
                      v-clipboard="newServer.secret"
                      @success="handleCopy(true)"
                      @error="handleCopy(false)"
                      type="button"
                      class="clipboard-btn"
                    >
                      <span class="fa fa-copy"></span>
                    </button>
                    <span
                      v-if="copySucceeded"
                      class="pficon pficon-ok renew-button"
                    ></span>
                  </div>
                </div>
                <div class="card-pf-items text-center">
                  <div class="alert alert-info alert-dismissable">
                    <span class="pficon pficon-info"></span>
                    {{ $t("servers.paste_it") }}
                    <a
                      href="http://docs.nethserver.org/en/v7/subscription.html"
                      class="alert-link"
                      target="_blank"
                      >{{ $t("servers.more_info") }}</a
                    >.
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-primary" data-dismiss="modal">
              {{ $t("servers.done") }}
            </button>
          </div>
        </div>
      </div>
    </div>
    <div
      class="modal fade"
      id="deleteServerModal"
      tabindex="-1"
      role="dialog"
      aria-labelledby="myModalLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button
              type="button"
              class="close"
              data-dismiss="modal"
              aria-hidden="true"
            >
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="myModalLabel">
              {{ $t("servers.delete") }}
            </h4>
          </div>
          <div class="modal-body">
            <form class="form-horizontal">
              <div class="alert alert-warning alert-dismissable">
                <span class="pficon pficon-warning-triangle-o"></span>
                <strong>{{ $t("servers.are_you_sure") }}</strong>
                {{ $t("servers.no_reversible") }}
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-default" data-dismiss="modal">
              {{ $t("servers.cancel") }}
            </button>
            <button
              @click="deleteServer()"
              type="button"
              class="btn btn-danger"
            >
              {{ $t("servers.delete") }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import LoginService from "./../services/login";
import StorageService from "./../services/storage";
import UtilService from "./../services/util";
import { setTimeout } from "timers";
import _ from "lodash";

import RenewButton from "./directives/RenewButton.vue";
import DeleteServer from "./directives/DeleteServer.vue";

export default {
  name: "servers",
  mixins: [LoginService, StorageService, UtilService],
  components: {
    renewButton: RenewButton,
    deleteServer: DeleteServer,
  },
  updated: function () {
    $('[data-toggle="tooltip"]').tooltip();
  },
  data() {
    // get plans list
    this.planList();

    //get servers list
    this.listServers();

    //get counters
    this.getCounters();

    // handle action
    var newServerAction = this.get("newServer") || false;
    if (newServerAction) {
      this.addServer("ns8");
      this.delete("newServer");
    }

    return {
      copySucceeded: false,
      isLoading: true,
      filters: {
        search: "",
        currentStatus: this.get("servers_filter_status") || "all",
        currentPlan: {
          name: "All",
          code: "all",
        },
        plans: [
          {
            name: "All",
            code: "all",
          },
        ],
        statuses: ["all", "active", "no_active", "no_comm"],
      },
      newServer: {},
      toDelete: {},
      servers: [],
      counters: {},
    };
  },
  methods: {
    extractServices(plan) {
      var context = this;
      var code = plan.split("+")[1];
      return code
        ? code
            .split(",")
            .sort()
            .map(function (el) {
              return context.$i18n.t("servers." + el);
            }).join(", ")
        : [];
    },
    getImage(s) {
      if (s.subscription.subscription_plan.base_code == "crostino") {
        return require("./../assets/crostino.svg");
      }
      if (s.subscription.subscription_plan.base_code == "lasagna") {
        return require("./../assets/lasagna.svg");
      }
      if (s.subscription.subscription_plan.base_code == "fiorentina") {
        return require("./../assets/fiorentina.svg");
      }
      if (s.subscription.subscription_plan.base_code == "pizza") {
        return require("./../assets/pizza.svg");
      }
      if (s.subscription.subscription_plan.base_code == "personal-ns8") {
        return require("./../assets/personal-ns8.svg");
      }
      if (s.subscription.subscription_plan.base_code == "personal-nsec") {
        return require("./../assets/personal-nsec.svg");
      }
      if (s.subscription.subscription_plan.base_code == "business-ns8") {
        return require("./../assets/business-ns8.svg");
      }
      if (s.subscription.subscription_plan.base_code == "business-nsec") {
        return require("./../assets/business-nsec.svg");
      }
      if (s.subscription.subscription_plan.base_code == "trial-ns8") {
        return require("./../assets/trial-ns8.svg");
      }
      if (s.subscription.subscription_plan.base_code == "trial-nsec") {
        return require("./../assets/trial-nsec.svg");
      }
      return require("./../assets/trial.svg");
    },
    handleCopy(status) {
      this.copySucceeded = status;
    },
    isExpired(date) {
      return new Date().toISOString() > date;
    },
    isUnsupported(plan) {
      return (
        plan == "crostino" ||
        plan == "pizza" ||
        plan == "fiorentina" ||
        plan == "lasagna" ||
        plan == "trial"
      );
    },
    addServer(type) {
      this.$http
        .post(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/ui/systems",
          {
            notification: {
              emails: [this.get("logged_user").email || ""],
            },
            type: type,
          },
          {
            headers: {
              Authorization: "Bearer " + this.get("access_token", false) || "",
            },
          }
        )
        .then(
          function (success) {
            this.newServer = success.body;
            setTimeout(function () {
              $("#newServerModal").modal("toggle");
            }, 0);
            this.listServers();
          },
          function (error) {
            console.error(error);
          }
        );
    },
    planList() {
      this.$http
        .get(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/ui/plans",
          {
            headers: {
              Authorization: "Bearer " + this.get("access_token", false) || "",
            },
          }
        )
        .then(
          function (success) {
            this.filters.plans = this.filters.plans.concat(success.body);
            this.filters.currentPlan =
              this.get("servers_filter_plan") || this.filters.plans[0];
          },
          function (error) {
            console.error(error);
          }
        );
    },
    listServers() {
      this.isLoading = true;
      this.$http
        .get(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/ui/systems",
          {
            headers: {
              Authorization: "Bearer " + this.get("access_token", false) || "",
            },
          }
        )
        .then(
          function (success) {
            this.servers = success.body;
            this.isLoading = false;
            this.getCounters();
          },
          function (error) {
            console.error(error);
            this.servers = [];
            this.isLoading = false;
          }
        );
    },
    getCounters() {
      this.$http
        .get(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/ui/systems/counters",
          {
            headers: {
              Authorization: "Bearer " + this.get("access_token", false) || "",
            },
          }
        )
        .then(
          function (success) {
            this.counters = success.body;
          },
          function (error) {
            console.error(error);
          }
        );
    },
    clearFilters(filter) {
      if (filter == "all") {
        this.filters.search = "";
        this.filters.currentStatus = "all";
        this.set("servers_filter_status", "all");
        this.filters.currentPlan = {
          name: "All",
          code: "all",
        };
        this.set("servers_filter_plan", {
          name: "All",
          code: "all",
        });
      }

      if (filter == "search") {
        this.filters.search = "";
      }

      if (filter == "status") {
        this.filters.currentStatus = "all";
        this.set("servers_filter_status", "all");
      }

      if (filter == "plan") {
        this.filters.currentPlan = {
          name: "All",
          code: "all",
        };
        this.set("servers_filter_plan", {
          name: "All",
          code: "all",
        });
      }
    },
    setFilterStatus(status) {
      this.filters.currentStatus = status;
      this.set("servers_filter_status", status);
      this.filteredServers();
    },
    setFilterPlan(plan) {
      this.filters.currentPlan = plan;
      this.set("servers_filter_plan", plan);
      this.filteredServers();
    },
    filteredServers() {
      var filter = {};
      if (this.filters.currentStatus != "all") {
        filter["status"] = this.filters.currentStatus;
      }
      if (this.filters.currentPlan.code != "all") {
        filter["subscription"] = {
          subscription_plan: {
            code: this.filters.currentPlan.code,
          },
        };
      }
      var filtered = _.filter(this.servers, filter);
      if (this.filters.search.length > 0) {
        var context = this;
        filtered = _.filter(filtered, function (item) {
          return (
            item.hostname
              .toLowerCase()
              .search(context.filters.search.toLowerCase()) >= 0 ||
            item.public_ip
              .toLowerCase()
              .search(context.filters.search.toLowerCase()) >= 0
          );
        });
      }
      return _.orderBy(filtered, ["hostname"], "asc");
    },
  },
};
</script>

<style scoped>
.create-server-ns8 {
  position: absolute;
  top: 75px;
  right: 200px;
  z-index: 99;
}
.create-server-nsec {
  position: absolute;
  top: 75px;
  right: 35px;
  z-index: 99;
}
.info-trial {
  margin-left: 10px;
  margin-right: 10px;
}

.adjust-height {
  min-height: 365px;
  max-height: 365px;
}
</style>
