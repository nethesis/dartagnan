<template>
  <span>
    <button
      v-if="
        isExpired(obj.subscription.valid_until) ||
        (plans.length > 0 &&
          plans[plans.length - 1].code !=
            obj.subscription.subscription_plan.code)
      "
      @click="showRenewModal(obj.id)"
      type="button"
      class="btn btn-primary"
    >
      <span class="fa fa-shopping-cart"></span>
      {{
        isExpired(obj.subscription.valid_until)
          ? $t("payment.renew_button")
          : $t("payment.upgrade_button")
      }}
    </button>
    <div
      class="modal fade"
      :id="'paymentModalRenew-' + obj.id"
      data-backdrop="static"
      tabindex="-1"
      role="dialog"
      aria-labelledby="myModalLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button
              @click="hideRenewModal()"
              type="button"
              class="close"
              data-dismiss="modal"
              aria-hidden="true"
            >
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="myModalLabel">
              {{ onUpgrade ? $t("payment.upgrade") : $t("payment.renew") }}
            </h4>
          </div>
          <div class="modal-body">
            <div class="card-pf-view">
              <div>
                <div class="card-pf-top-element">
                  <span
                    :class="[
                      'fa',
                      onUpgrade ? 'fa-arrow-up' : 'fa-refresh',
                      'card-pf-icon-circle',
                      'adjust-icon-size',
                    ]"
                  ></span>
                </div>
                <h2 class="card-pf-title text-center">
                  {{
                    onUpgrade
                      ? $t("payment.upgrade_proceed")
                      : $t("payment.renew_proceed")
                  }}
                </h2>
                <div class="card-pf-items text-center">
                  <div class="card-pf-item details-pay-item">
                    <span class="card-pf-item-text">
                      <strong>{{ $t("payment.plan") }}</strong>
                    </span>
                  </div>
                  <div class="card-pf-item details-pay-item">
                    <div class="dropdown">
                      <button
                        class="btn btn-default dropdown-toggle"
                        type="button"
                        id="dropdownMenu1"
                        data-toggle="dropdown"
                      >
                        {{ currentPlan.name }}
                        <span class="caret"></span>
                      </button>
                      <ul
                        class="dropdown-menu"
                        role="menu"
                        aria-labelledby="dropdownMenu1"
                      >
                        <li
                          v-for="p in plans"
                          v-if="
                            p.code !== 'trial' &&
                            p.code !== 'trial-ns8' &&
                            p.code !== 'trial-nsec'
                          "
                          v-bind:key="p.code"
                          role="presentation"
                          :class="[
                            p.code == currentPlan.code ? 'selected' : '',
                            p.code != obj.subscription.subscription_plan.code &&
                            (p.base_price || p.price || p.full_price) <
                              obj.subscription.subscription_plan.base_price
                              ? 'disabled'
                              : '',
                          ]"
                        >
                          <a
                            @click="
                              p.code !=
                                obj.subscription.subscription_plan.code &&
                              (p.full_price || p.price || p.base_price) <
                                obj.subscription.subscription_plan.base_price
                                ? null
                                : changePlan(p, false)
                            "
                            role="menuitem"
                            tabindex="-1"
                            href="#"
                          >
                            <img
                              :src="getImage(p.base_code)"
                              class="plan-icon-mini"
                            />
                            {{ p.name }}</a
                          >
                        </li>
                      </ul>
                    </div>
                  </div>
                </div>

                <div
                  v-show="
                    currentPlan.code != 'crostino' &&
                    currentPlan.code != 'pizza' &&
                    currentPlan.code != 'lasagna' &&
                    currentPlan.code != 'fiorentina'
                  "
                  class="card-pf-items text-center"
                >
                  <div
                    v-if="services[extractProduct(currentPlan.code)] && services[extractProduct(currentPlan.code)].length > 0"
                    class="card-pf-item details-pay-item"
                  >
                    <span class="card-pf-item-text">
                      <strong>{{ $t("payment.services") }}</strong>
                    </span>
                  </div>
                  <div class="card-pf-item details-pay-item">
                    <ul class="services-list">
                      <li
                        v-for="s in services[extractProduct(currentPlan.code)]"
                        v-bind:key="s.code"
                      >
                        <input
                          type="checkbox"
                          :name="s.code"
                          :value="s.code"
                          v-model="selectedServices"
                          @change="changePlan(currentPlan, true)"
                          :disabled="currentServices.indexOf(s.code) != -1"
                        />
                        <label :for="s.code">{{ s.name }} (+{{s.price}}€)</label
                        ><br />
                      </li>
                    </ul>
                  </div>
                </div>

                <div class="card-pf-items text-center">
                  <div class="card-pf-item details-pay-item">
                    <span class="card-pf-item-text">
                      <strong>{{ $t("payment.details") }}</strong>
                    </span>
                  </div>
                  <div class="card-pf-item details-pay-item">
                    <div
                      v-if="!onUpgradePriceCalc"
                      class="details-markdown text-left"
                      v-html="markdownDescription"
                    ></div>
                    <div
                      v-if="onUpgradePriceCalc"
                      class="spinner spinner-sm"
                    ></div>
                  </div>
                </div>
                <div class="card-pf-items text-center">
                  <div class="card-pf-item details-pay-item">
                    <span class="card-pf-item-text">
                      <strong>{{ $t("payment.price") }}</strong>
                    </span>
                  </div>
                  <div class="card-pf-item details-pay-item">
                    <span v-if="!onUpgradePriceCalc" class="card-pf-item-text">
                      <strong
                        >{{
                          currentPlan.full_price > 0
                            ? currentPlan.full_price
                            : 0
                        }}€</strong
                      >
                    </span>
                    <div
                      v-if="onUpgradePriceCalc"
                      class="spinner spinner-sm"
                    ></div>
                  </div>
                </div>

                <div
                  v-if="
                    onUpgrade &&
                    currentPlan.price != currentPlan.full_price &&
                    discounts.annualDiscount > 0
                  "
                  class="card-pf-items text-center"
                >
                  <div class="card-pf-item details-pay-item">
                    <span class="card-pf-item-text">
                      <strong>{{
                        $t("payment.previous_license_discount")
                      }}</strong>
                    </span>
                  </div>
                  <div class="card-pf-item details-pay-item">
                    <span v-if="!onUpgradePriceCalc" class="card-pf-item-text">
                      <strong>{{ discounts.annualDiscount }}€</strong>
                    </span>
                    <div
                      v-if="onUpgradePriceCalc"
                      class="spinner spinner-sm"
                    ></div>
                  </div>
                </div>

                <div
                  v-if="discounts.volumeDiscount != 0"
                  class="card-pf-items text-center"
                >
                  <div class="card-pf-item details-pay-item">
                    <span class="card-pf-item-text">
                      <strong>{{ $t("payment.volume_discount") }}</strong>
                    </span>
                  </div>
                  <div class="card-pf-item details-pay-item">
                    <span v-if="!onUpgradePriceCalc" class="card-pf-item-text">
                      <strong>{{ discounts.volumeDiscount }}%</strong>
                      <span>
                        (<strong>{{ discounts.count }}</strong>
                        {{ $t("payment.active_licenses") }})</span
                      >
                    </span>
                    <div
                      v-if="onUpgradePriceCalc"
                      class="spinner spinner-sm"
                    ></div>
                  </div>
                </div>

                <div class="card-pf-items text-center">
                  <div class="card-pf-item details-pay-item">
                    <span class="card-pf-item-text">
                      <strong>{{ $t("payment.period") }}</strong>
                    </span>
                  </div>
                  <div class="card-pf-item details-pay-item">
                    <span
                      v-if="!onUpgrade && !onUpgradePriceCalc"
                      class="card-pf-item-text"
                      >{{ obj.subscription.valid_until | formatDate(false) }} -
                      {{
                        calculateSubscription(
                          obj.subscription.valid_until,
                          obj.subscription.subscription_plan
                        ) | formatDate(false)
                      }}</span
                    >
                    <span
                      v-if="onUpgrade && !onUpgradePriceCalc"
                      class="card-pf-item-text"
                      >{{ new Date().toISOString() | formatDate(false) }} -
                      {{
                        calculateSubscription(
                          new Date().toISOString(),
                          currentPlan
                        ) | formatDate(false)
                      }}</span
                    >
                    <div
                      v-if="onUpgradePriceCalc"
                      class="spinner spinner-sm"
                    ></div>
                  </div>
                </div>

                <div class="card-pf-items text-center">
                  <div class="card-pf-item details-pay-item">
                    <span class="card-pf-item-text">
                      <strong>{{ $t("payment.final_price") }}</strong>
                    </span>
                  </div>
                  <div class="card-pf-item details-pay-item">
                    <span v-if="!onUpgradePriceCalc" class="card-pf-item-text">
                      <strong>{{ currentPlan.price }}€</strong>
                      <span>+ {{ $t("payment.taxes") }}</span>
                    </span>
                    <div
                      v-if="onUpgradePriceCalc"
                      class="spinner spinner-sm"
                    ></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="divider"></div>
          <div class="modal-footer text-center">
            <div v-if="payment.onProgress" class="spinner"></div>
            <div
              v-if="payment.done && !errors.state && !payment.onProgress"
              class="alert alert-success alert-dismissable"
            >
              <span class="pficon pficon-ok"></span>
              <strong>{{ $t("payment.confirmed") }}</strong
              >. {{ $t("payment.payment_id_ref") }}
              <pre
                class="filters-container"
              ><strong>{{payment.details.paymentID}}</strong></pre>
            </div>
            <button
              v-if="payment.done && !errors.state && !payment.onProgress"
              type="button"
              class="btn btn-primary done-payment-button"
              @click="hideRenewModal()"
            >
              {{ $t("servers.done") }}
            </button>
            <div
              v-if="
                !payment.done &&
                !errors.state &&
                !payment.onProgress &&
                currentPlan.price > 0
              "
              class="card-pf-item details-pay-item"
            >
              <span class="card-pf-item-text">
                <strong>{{ $t("payment.pay_with") }}</strong>
              </span>
            </div>
            <div
              v-show="currentPlan.price == 0"
              class="alert alert-warning alert-dismissable"
            >
              <span class="pficon pficon-warning-triangle-o"></span>
              <strong>{{ $t("payment.invalid_subscription") }}</strong>
              {{ $t("payment.invalid_subscription_text") }}
            </div>
            <div
              v-show="
                !payment.done &&
                !errors.state &&
                !payment.onProgress &&
                currentPlan.price > 0
              "
              :id="'paypal-button-container-' + obj.id"
            ></div>
            <div
              v-if="payment.done && errors.state && !payment.onProgress"
              class="alert alert-danger alert-dismissable"
            >
              <span class="pficon pficon-error-circle-o"></span>
              <strong>{{ $t("payment.payment_error") }}</strong
              >. {{ errors.message || $t("payment.payment_error_details") }}
              <p>{{ $t("payment.payment_id_ref") }}</p>
              <pre
                class="filters-container"
              ><strong>{{payment.details.paymentID}}</strong></pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  </span>
</template>
<script>
import StorageService from "./../../services/storage";

import paypal from "paypal-checkout";
import { setTimeout } from "timers";
import marked from "marked";

export default {
  name: "RenewButton",
  props: ["obj", "update"],
  mixins: [StorageService],
  mounted() {
    var context = this;
    paypal.Button.render(
      {
        env: CONFIG.PAYPAL_PRODUCTION ? "production" : "sandbox",
        style: {
          layout: "vertical", // horizontal | vertical
          size: "medium", // medium | large | responsive
          shape: "rect", // pill | rect
          color: "gold", // gold | blue | silver | black
        },
        funding: {
          allowed: [paypal.FUNDING.CARD],
          disallowed: [paypal.FUNDING.CREDIT],
        },
        client: {
          sandbox: CONFIG.PAYPAL_SANDBOX,
          production: CONFIG.PAYPAL_PRODUCTION,
        },
        payment: function (data, actions) {
          return actions.payment.create({
            payment: {
              transactions: [
                {
                  amount: {
                    total:
                      Math.round(
                        (context.currentPlan.price +
                          (context.currentPlan.price *
                            context.billingInfo.tax) /
                            100) *
                          100
                      ) / 100,
                    currency: "EUR",
                    details: {
                      subtotal:
                        Math.round(context.currentPlan.price * 100) / 100,
                      tax:
                        Math.round(
                          ((context.currentPlan.price *
                            context.billingInfo.tax) /
                            100) *
                            100
                        ) / 100,
                    },
                  },
                  item_list: {
                    items: [
                      {
                        name: context.currentPlan.code,
                        description: context.currentPlan.name,
                        sku: context.obj.uuid,
                        price: context.currentPlan.price,
                        currency: "EUR",
                        quantity: "1",
                      },
                    ],
                  },
                },
              ],
            },
            experience: {
              input_fields: {
                no_shipping: 1,
              },
            },
          });
        },

        onAuthorize: function (data, actions) {
          return actions.payment.execute().then(function () {
            context.payment.onProgress = true;
            if (context.onUpgrade) {
              context.upgradeCheck(data);
            } else {
              context.renewCheck(data);
            }
          });
        },
      },
      "#paypal-button-container-" + this.obj.id
    );
  },
  data() {
    // get plans
    this.plansList();

    // read upgrade ref and show modal
    if (this.get("upgrade_ref", false)) {
      var context = this;
      setTimeout(function () {
        context.showRenewModal(context.get("upgrade_ref", false));
        context.delete("upgrade_ref");
      }, 0);
    }

    return {
      payment: {
        done: false,
        details: {},
        onProgress: false,
      },
      errors: {
        message: "",
        state: false,
      },
      plans: [],
      services: {
        ns8: _.orderBy([], "name", "asc"),
        nsec: _.orderBy(
          [
            { name: "Hotspot", code: "hotspot", price: 60.0 },
            { name: "Threat shield", code: "threat_shield", price: 96.0 },
            { name: "Flashstart Lite", code: "flashstart_lite", price: 18.0 },
          ],
          "name",
          "asc"
        ),
      },
      selectedServices: this.extractServices(
        this.obj.subscription.subscription_plan.code
      ),
      selectedPlan: {},
      markdownDescription: "",
      currentPlan: this.obj.subscription.subscription_plan,
      currentServices: this.extractServices(
        this.obj.subscription.subscription_plan.code
      ),
      billingInfo: {},
      discounts: {
        volumeDiscount: 0,
        count: 0,
        annualDiscount: 0,
      },
      onUpgradePriceCalc: false,
      onUpgrade: false,
    };
  },
  methods: {
    getImage(s) {
      if (s == "personal-ns8") {
        return require("./../../assets/fiorentina.svg");
      }
      if (s == "personal-nsec") {
        return require("./../../assets/pizza.svg");
      }
      if (s == "business-ns8") {
        return require("./../../assets/crostino.svg");
      }
      if (s == "business-nsec") {
        return require("./../../assets/lasagna.svg");
      }
    },
    extractProduct(license) {
      var product = license.split("-")[1];
      return product ? product.split("+")[0] : "";
    },
    extractServices(plan) {
      var code = plan.split("+")[1];
      return code ? code.split(",").sort() : [];
    },
    isExpired(date) {
      return new Date().toISOString() > date;
    },
    showRenewModal(id) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/ui/billings",
          {
            headers: {
              Authorization: "Bearer " + this.get("access_token", false) || "",
            },
          }
        )
        .then(
          function (success) {
            this.billingInfo = success.body;
            this.payment.done = false;
            this.payment.onProgress = false;
            this.payment.details = {};
            this.errors.message = "";
            this.errors.state = false;
            this.currentPlan =
              this.obj.subscription.subscription_plan.code != "trial" &&
              this.obj.subscription.subscription_plan.code != "trial-ns8" &&
              this.obj.subscription.subscription_plan.code != "trial-nsec"
                ? this.obj.subscription.subscription_plan
                : this.plans[1];
            if (
              this.obj.subscription.subscription_plan.code == "trial" ||
              this.obj.subscription.subscription_plan.code == "trial-ns8" ||
              this.obj.subscription.subscription_plan.code == "trial-nsec"
            ) {
              this.currentPlan.full_price = this.plans[1].price;
            } else {
              this.currentPlan.full_price = this.currentPlan.price;
            }
            this.markdownDescription = marked(this.currentPlan.description, {
              sanitize: true,
            });
            this.onUpgradePriceCalc = false;
            this.onUpgrade = !this.isExpired(this.obj.subscription.valid_until);

            // check volume discounts
            this.$http
              .get(
                this.$root.$options.api_scheme +
                  this.$root.$options.api_host +
                  "/api/ui/plans/volume_discount",
                {
                  headers: {
                    Authorization:
                      "Bearer " + this.get("access_token", false) || "",
                  },
                }
              )
              .then(
                function (success) {
                  this.currentPlan.price =
                    this.currentPlan.price -
                    (this.currentPlan.full_price * success.body.discount) / 100;
                  this.discounts.volumeDiscount = success.body.discount;
                  this.discounts.count = success.body.count;
                  $("#paymentModalRenew-" + id).modal("toggle");
                },
                function (error) {
                  console.error(error);
                }
              );
          },
          function (error) {
            this.$parent.$parent.action = "updateBilling";
            this.set("upgrade_ref", id);
            this.$router.push({
              path: "/profile",
            });
            console.error(error);
          }
        );
    },
    hideRenewModal() {
      this.update();
      $("#paymentModalRenew-" + this.obj.id).modal("hide");
    },
    plansList() {
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
            this.plans = _.orderBy(success.body, "price", "asc");
            this.plans = this.plans.filter(
              (plan) =>
                plan.code != "trial" &&
                plan.code != "trial-ns8" &&
                plan.code != "trial-nsec" &&
                plan.code != "crostino" &&
                plan.code != "lasagna" &&
                plan.code != "pizza" &&
                plan.code != "fiorentina"
            );

            // check current subscription
            if (this.obj.subscription.subscription_plan.code.includes("ns8")) {
              this.plans = this.plans.filter(
                (plan) =>
                  plan.code != "trial-nsec" &&
                  plan.code != "personal-nsec" &&
                  plan.code != "business-nsec"
              );
            }
            if (this.obj.subscription.subscription_plan.code.includes("nsec")) {
              this.plans = this.plans.filter(
                (plan) =>
                  plan.code != "trial-ns8" &&
                  plan.code != "personal-ns8" &&
                  plan.code != "business-ns8"
              );
            }
          },
          function (error) {
            console.error(error);
          }
        );
    },
    changePlan(plan, withServices) {
      if (this.currentServices.length > 0 && !withServices) {
        this.selectedServices = this.currentServices;
        withServices = true;
      }

      var newCode = plan.base_code;
      if (withServices && this.selectedServices.length > 0) {
        newCode = "!" + newCode + "+" + this.selectedServices.sort().join(",");
      } else {
        this.selectedServices = [];
      }

      if (
        newCode !== this.obj.subscription.subscription_plan.code ||
        withServices
      ) {
        // handle upgrade
        var context = this;
        this.calculateUpgradePrice(newCode, function (data) {
          context.onUpgrade = true;
          context.currentPlan = plan;
          context.currentPlan.price = data.price;
          context.currentPlan.full_price = data.full_price;
          context.discounts.annualDiscount =
            Math.round(data.discount * 100) / 100;
          context.markdownDescription = marked(plan.description, {
            sanitize: true,
          });

          // check volume discounts
          context.$http
            .get(
              context.$root.$options.api_scheme +
                context.$root.$options.api_host +
                "/api/ui/plans/volume_discount",
              {
                headers: {
                  Authorization:
                    "Bearer " + context.get("access_token", false) || "",
                },
              }
            )
            .then(
              function (success) {
                context.currentPlan.price =
                  context.currentPlan.price -
                  (context.currentPlan.full_price * success.body.discount) /
                    100;
                context.discounts.volumeDiscount = success.body.discount;
                context.discounts.count = success.body.count;
              },
              function (error) {
                console.error(error);
              }
            );
        });
      } else {
        this.onUpgrade = false;
        this.currentPlan = plan;
        this.currentPlan.full_price = plan.price;
        this.markdownDescription = marked(plan.description, { sanitize: true });

        // check volume discounts
        this.$http
          .get(
            this.$root.$options.api_scheme +
              this.$root.$options.api_host +
              "/api/ui/plans/volume_discount",
            {
              headers: {
                Authorization:
                  "Bearer " + this.get("access_token", false) || "",
              },
            }
          )
          .then(
            function (success) {
              this.currentPlan.price =
                this.currentPlan.price -
                (this.currentPlan.full_price * success.body.discount) / 100;
              this.discounts.volumeDiscount = success.body.discount;
              this.discounts.count = success.body.count;
            },
            function (error) {
              console.error(error);
            }
          );
      }
    },
    calculateUpgradePrice(plan, callback) {
      this.onUpgradePriceCalc = true;
      this.$http
        .get(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/ui/systems/" +
            this.obj.id +
            "/upgrade_price?plan=" +
            encodeURIComponent(plan),
          {
            headers: {
              Authorization: "Bearer " + this.get("access_token", false) || "",
            },
          }
        )
        .then(
          function (success) {
            this.onUpgradePriceCalc = false;
            callback(success.body);
          },
          function (error) {
            console.error(error);
            this.onUpgradePriceCalc = false;
          }
        );
    },
    calculateSubscription(date, subscription) {
      var moment = require("moment/moment.js");
      return moment(date, "YYYY-MM-DDTHH:mm:ss").add(
        subscription.period,
        "days"
      );
    },
    renewCheck(payment) {
      this.$http
        .post(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/ui/systems/" +
            this.obj.id +
            "/renewal",
          {
            payment_id: payment.paymentID,
          },
          {
            headers: {
              Authorization: "Bearer " + this.get("access_token", false) || "",
            },
          }
        )
        .then(
          function (success) {
            this.payment.onProgress = false;
            this.payment.details = payment;
            this.payment.done = true;
          },
          function (error) {
            console.error(error);
            this.payment.onProgress = false;
            this.payment.details = payment;
            this.payment.done = true;
            this.errors.state = true;
            this.errors.message = error.body.message;
          }
        );
    },
    upgradeCheck(payment) {
      this.$http
        .post(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/ui/systems/" +
            this.obj.id +
            "/upgrade",
          {
            payment_id: payment.paymentID,
            subscription_plan_id: this.currentPlan.id,
          },
          {
            headers: {
              Authorization: "Bearer " + this.get("access_token", false) || "",
            },
          }
        )
        .then(
          function (success) {
            this.payment.onProgress = false;
            this.payment.details = payment;
            this.payment.done = true;
          },
          function (error) {
            console.error(error);
            this.payment.onProgress = false;
            this.payment.details = payment;
            this.payment.done = true;
            this.errors.state = true;
            this.errors.message = error.body.message;
          }
        );
    },
  },
};
</script>
<style>
.services-list {
  list-style: none;
  padding-left: 5px;
  text-align: left;
  margin-bottom: 0px;
}

.plan-icon-mini {
  width: 30px;
  margin-right: 0px;
}
</style>
