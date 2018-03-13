<template>
  <div>
    <h2>{{$t('profile.profile')}}</h2>
    <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
      <div class="card-pf card-pf-view">
        <div class="card-pf-body">
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
        <div class="card-pf-view-checkbox">
          <input type="checkbox">
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
      return {
        profile: this.get('logged_user')
      }
    },
    methods: {}
  }

</script>

<style scoped>


</style>
