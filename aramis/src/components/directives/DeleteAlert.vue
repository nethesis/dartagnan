<template>
  <div>
    <button @click="showDeleteModal()" type="button" class="btn btn-danger">
      {{$t('servers.delete')}}
    </button>
    <div class="modal fade" :id="'deleteAlertModal-'+obj.id" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="myModalLabel">{{$t('servers.delete')}} {{obj.namei18n || '-'}}</h4>
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
            <button @click="deleteAlert()" type="button" class="btn btn-danger">{{$t('servers.delete')}}</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import StorageService from "./../../services/storage";

export default {
  name: "DeleteAlert",
  props: ["obj", "update"],
  mixins: [StorageService],
  data() {
    return {};
  },
  methods: {
    showDeleteModal() {
      $("#deleteAlertModal-" + this.obj.id).modal("toggle");
    },
    deleteAlert() {
      var closeId = this.obj.id;
      this.$http
        .delete(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/ui/alerts/" +
            this.obj.id,
          {
            headers: {
              Authorization: "Bearer " + this.get("access_token", false) || ""
            }
          }
        )
        .then(
          function(success) {
            $("#deleteAlertModal-" + closeId).modal("hide");
            this.update();
          },
          function(error) {
            console.error(error);
          }
        );
    }
  }
};
</script>
<style>

</style>
