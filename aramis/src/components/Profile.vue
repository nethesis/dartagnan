<template>
  <div>
    <h2>{{$t('profile.profile')}}</h2>
    <div v-if="updateBilling" class="row row-cards-pf no-padding-top row-divider blank-slate-pf">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
        <div class="alert alert-warning alert-dismissable">
          <span class="pficon pficon-warning-triangle-o"></span>
          <strong>{{$t('profile.billing_not_configured')}}</strong>. {{$t('profile.billing_not_configured_desc')}}.
        </div>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
      <div class="card-pf card-pf-view">
        <div class="card-pf-body">
          <h2 class="card-pf-title">
            {{$t('profile.info')}}
          </h2>
          <div class="card-pf-top-element">
            <img :src="profile.picture" class="fa fa-birthday-cake card-pf-icon-circle profile-picture-main">
          </div>
          <h2 class="card-pf-title text-center">
            {{profile.name}}
          </h2>
          <div class="card-pf-items text-center">
            <div class="card-pf-item">
              <span class="fa fa-envelope"></span>
              <span class="card-pf-item-text">{{profile.email}}</span>
            </div>
            <div class="card-pf-item">
              <span data-toggle="tooltip" data-placement="top" :title="$t('profile.email_is') + ' '+ (profile.email_verified ? $t('profile.verified') : $t('profile.not_verified'))"
                :class="['fa', profile.email_verified ? 'fa-check' : 'fa-remove red']"></span>
            </div>
          </div>
          <p class="card-pf-info text-center">
            <strong>{{$t('profile.last_login')}}</strong> {{profile.updated_at | formatDate}}
            <div class="text-center">ID:
              <strong class="soft">{{profile.sub}}</strong>
            </div>
          </p>
        </div>
      </div>
    </div>

    <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
      <div class="card-pf card-pf-view">
        <h2 class="card-pf-title">
          {{$t('profile.billing')}}
        </h2>
        <div class="card-pf-body">
          <form class="form-horizontal" v-on:submit.prevent="updateBillingInfo()">
            <div class="form-group">
              <label class="col-sm-3 control-label" for="textInput-markup">{{$t('profile.name')}}</label>
              <div class="col-sm-9">
                <input v-model="billingInfo.name" required type="text" id="textInput-markup" class="form-control">
              </div>
            </div>
            <div class="form-group">
              <label class="col-sm-3 control-label" for="textInput-markup">{{$t('profile.address')}}</label>
              <div class="col-sm-9">
                <input v-model="billingInfo.address" required type="text" id="textInput-markup" class="form-control">
              </div>
            </div>
            <div class="form-group">
              <label class="col-sm-3 control-label" for="textInput-markup">{{$t('profile.city')}}</label>
              <div class="col-sm-9">
                <input v-model="billingInfo.city" required type="text" id="textInput-markup" class="form-control">
              </div>
            </div>
            <div class="form-group">
              <label class="col-sm-3 control-label" for="textInput-markup">{{$t('profile.postal_code')}}</label>
              <div class="col-sm-9">
                <input v-model="billingInfo.postal_code" required type="text" id="textInput-markup" class="form-control">
              </div>
            </div>
            <div class="form-group">
              <label class="col-sm-3 control-label" for="textInput-markup">{{$t('profile.country')}}</label>
              <div class="col-sm-9">
                <select required v-model="billingInfo.country" class="form-control">
                  <option v-for="c in countries" v-bind:key="c.id" v-bind:value="c.country">
                    {{ c.country }}
                  </option>
                </select>
              </div>
            </div>
            <div class="form-group">
              <label class="col-sm-3 control-label" for="textInput-markup">{{$t('profile.type')}}</label>
              <div class="col-sm-9">
                <span class="span-radio">
                  <input required v-model="billingInfo.type" class="form-check-input" type="radio" name="exampleRadios" id="exampleRadios1"
                    value="business" checked>
                  <label class="form-check-label" for="exampleRadios1">
                    {{$t('profile.type_business')}}
                  </label>
                </span>
                <span class="span-radio">
                  <input required v-model="billingInfo.type" class="form-check-input" type="radio" name="exampleRadios" id="exampleRadios2"
                    value="person">
                  <label class="form-check-label" for="exampleRadios2">
                    {{$t('profile.type_person')}}
                  </label>
                </span>
              </div>
            </div>
            <div v-if="billingInfo.type == 'business'" class="form-group">
              <label class="col-sm-3 control-label" for="textInput-markup">{{$t('profile.vat')}}</label>
              <div class="col-sm-9">
                <input v-model="billingInfo.vat" required type="text" id="textInput-markup" class="form-control">
              </div>
            </div>
            <div class="modal-footer">
              <div v-if="isSaving" class="spinner spinner-sm left-spinner"></div>
              <button :disabled="isSaving" type="submit" class="btn btn-primary">Save</button>
            </div>
          </form>
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
    name: 'alerts',
    mixins: [LoginService, StorageService, UtilService],
    data() {
      setTimeout(function () {
        $('[data-toggle="tooltip"]').tooltip()
      }, 500)

      var updateBilling = false
      if (this.$parent.action == 'updateBilling') {
        updateBilling = true
        this.delete('query_params')
      }

      // get country list
      this.getCountryList()

      // get billing info
      this.getBillingInfo()

      return {
        profile: this.get('logged_user'),
        updateBilling: updateBilling,
        billingInfo: {
          type: 'person'
        },
        billingEmpty: true,
        countries: [],
        isSaving: false
      }
    },
    methods: {
      getBillingInfo() {
        this.$http.get('https://' + this.$root.$options.api_host + '/api/ui/billings', {
          headers: {
            'Authorization': 'Bearer ' + this.get('access_token', false) || ''
          }
        }).then(function (success) {
          this.billingInfo = success.body
          if (this.billingInfo.vat && this.billingInfo.vat.length > 0) {
            this.billingInfo.type = 'business'
          } else {
            this.billingInfo.type = 'person'
          }
          this.billingEmpty = false
        }, function (error) {
          console.error(error)
        });
      },
      getCountryList() {
        this.$http.get('https://' + this.$root.$options.api_host + '/api/ui/taxes', {
          headers: {
            'Authorization': 'Bearer ' + this.get('access_token', false) || ''
          }
        }).then(function (success) {
          this.countries = success.body
        }, function (error) {
          console.error(error)
        });
      },
      updateBillingInfo() {
        this.isSaving = true
        this.$http[this.billingEmpty ? 'post' : 'put']('https://' + this.$root.$options.api_host + '/api/ui/billings',
          this.billingInfo, {
            headers: {
              'Authorization': 'Bearer ' + this.get('access_token', false) || ''
            }
          }).then(function (success) {
          this.isSaving = false
          this.updateBilling = false
        }, function (error) {
          this.isSaving = false
          console.error(error)
        });
      }
    }
  }

</script>

<style scoped>
  .left-spinner {
    right: 110px;
    position: absolute;
    bottom: 62px;
  }

</style>
