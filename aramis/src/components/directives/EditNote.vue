<template>
  <span class="inline-block">
    <button class="btn btn-default right" type="button" @click="openAlertNoteModal()">
      <span class="pficon pficon-edit"></span>
    </button>
    <div class="modal fade" :id="'noteAlert-'+obj.id" tabindex="-1" role="dialog" aria-labelledby="noteAlert" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="myModalLabel">{{$t('servers.note_alert_for')}} {{obj.namei18n}}</h4>
          </div>
          <div class="modal-body">
            <form class="form-horizontal">
              <div class="form-group">
                <label class="col-sm-3 control-label" for="textInput-modal-markup">{{$t('servers.note')}}</label>
                <div class="col-sm-9">
                  <textarea v-model="obj.note" type="text" id="textInput-modal-markup" class="form-control"></textarea>
                </div>
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-default" data-dismiss="modal">{{$t('servers.cancel')}}</button>
            <button @click="saveAlertNote()" type="button" class="btn btn-primary">{{$t('servers.save')}}</button>
          </div>
        </div>
      </div>
    </div>
  </span>
</template>

<script>
  import StorageService from './../../services/storage';
  import _ from 'lodash'

  export default {
    name: 'editNote',
    mixins: [StorageService],
    props: ['obj'],
    data() {
      return {
      }
    },
    methods: {
      openAlertNoteModal() {
        $('#noteAlert-'+this.obj.id).modal('toggle')
      },
      saveAlertNote() {
        var closeId = this.obj.id
        this.$http.put(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/ui/alerts/' + this.obj.id, {
          system_id: this.obj.system.id.toString(),
          note: this.obj.note
        }, {
          headers: {
            'Authorization': 'Bearer ' + this.get('access_token', false) || ''
          }
        }).then(function (success) {
          $('#noteAlert-'+closeId).modal('hide')
        }, function (error) {
          console.error(error)
        });
      }
    }
  }

</script>
