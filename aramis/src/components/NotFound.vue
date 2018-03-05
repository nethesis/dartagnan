<template>
  <div>
    <h2>Not Found</h2>
  </div>
</template>

<script>
  import LoginService from './../services/login';
  import StorageService from './../services/storage';
  import UtilService from './../services/util';

  export default {
    name: 'notfound',
    mixins: [LoginService, StorageService, UtilService],
    data() {
      return {}
    },
    methods: {}
  }

</script>

<style scoped>

</style>
