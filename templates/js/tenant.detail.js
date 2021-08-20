(function ($) {
  'use strict';
  var $recorId = $("#recordId");
  var $dtBasic = $("#dtBasic");
  var $btnNewTenant = $("#btn-newTenant");
  var $btnSaveOrUpdate = $("#btn-save");
  var $formDetails = $('#form-details');
  var $modalForm = $("#modal-addNew");

  var pageFunction = function () {
    var loadDetails = function (id) {
      $.ajax({
        url: $.helper.baseApiPath("/tenants/getById/" + $recorId.val()),
        type: 'GET',
        success: function (r) {
          if (r.status) {
            $formDetails.find('input').val(function () {
              return r.data[this.name];
            });
            $("textarea[name=description]").val(r.data.description);
            $(".tenant-name").text(r.data.name);
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning");
        }
      });
    }


    return {
      init: function () {
        loadDetails()
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));