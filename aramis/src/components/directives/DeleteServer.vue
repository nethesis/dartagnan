<template>
  <div>
    <button @click="showDeleteModal()" type="button" class="btn btn-danger">
      {{$t('servers.delete')}}
    </button>
    <div class="modal fade" :id="'deleteServerModal-'+obj.id" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="myModalLabel">{{$t('servers.delete')}} {{obj.hostname || '-'}}</h4>
          </div>
          <div class="modal-body">
            <form class="form-horizontal">
              <div class="alert alert-warning alert-dismissable">
                <span class="pficon pficon-warning-triangle-o"></span>
                <strong>{{$t('servers.are_you_sure')}}</strong> {{$t('servers.no_reversible')}}
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-default" data-dismiss="modal">{{$t('servers.cancel')}}</button>
            <button @click="deleteServer()" type="button" class="btn btn-danger">{{$t('servers.delete')}}</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
  import StorageService from './../../services/storage';

  export default {
    name: 'DeleteServer',
    props: ['obj', 'update', 'redir'],
    mixins: [StorageService],
    data() {
      return {
      }
    },
    methods: {
      showDeleteModal() {
        $('#deleteServerModal-' + this.obj.id).modal('toggle')
      },
      deleteServer() {
        var closeId = this.obj.id
        this.$http.delete('https://' + this.$root.$options.api_host + '/api/ui/systems/' + this.obj.id, {
          headers: {
            'Authorization': 'Bearer ' + this.get('access_token', false) || ''
          }
        }).then(function (success) {
          $('#deleteServerModal-' + closeId).modal('hide')
          if (this.redir) {
            this.$router.push({
              path: '/' + this.redir
            })
          } else {
            this.update()
          }
        }, function (error) {
          console.error(error)
        });
      },
    }
  }

</script>
<style>


</style>
