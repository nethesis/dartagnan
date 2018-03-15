<template>
  <div>
    <button @click="deleteAlert()" type="button" class="btn btn-danger">
      {{$t('servers.delete')}}
    </button>
  </div>
</template>
<script>
  import StorageService from './../../services/storage';

  export default {
    name: 'DeleteAlert',
    props: ['obj', 'update'],
    mixins: [StorageService],
    data() {
      return {}
    },
    methods: {
      deleteAlert() {
        this.$http.delete('https://' + this.$root.$options.api_host + '/api/ui/alerts/' + this.obj.id, {
          headers: {
            'Authorization': 'Bearer ' + this.get('access_token', false) || ''
          }
        }).then(function (success) {
          this.update()
        }, function (error) {
          console.error(error)
        });
      },
    }
  }

</script>
<style>


</style>
