<template>
  <div>
    <h2>Profile</h2>
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
      return {}
    },
    methods: {
    }
  }

</script>

<style scoped>

</style>
