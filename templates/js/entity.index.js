(function ($) {
  'use strict';
  var $dtBasic = $("#dtBasic");
  var $btnNewTenant = $("#btn-newTenant");
  var $btnSaveOrUpdate = $("#btn-save");
  var $form = $('#form-basic');
  var $modalForm = $("#modal-addNew");
  var $recordId = $("#recordId");

  var pageFunction = function () {
    var $dt;
    var loadDatatables = function () {
      $dt = $dtBasic.DataTable({
        responsive: true,
        pageLength: 10,
        ajax: {
          url: $.helper.baseApiPath("/entity/getall"),
          type: "GET",
          contentType: "application/json",
          data: function (d) {
            return JSON.stringify(d);
          }
        },
        columns: [{
            data: "id",
            orderable: false,
            searchable: false,
            class: "text-center",
            render: function (data, type, row, meta) {
              return meta.row + meta.settings._iDisplayStart + 1;
            }
          },
          {
            data: "name"
          },
          {
            data: "owner_name"
          },
          {
            data: "address"
          },
          {
            data: "email"
          },
          {
            data: "phone"
          },
          {
            data: "description"
          },
          {
            data: "id",
            orderable: true,
            searchable: true,
            class: "text-center top-0 start-0",
            render: function (data, type, row) {
              var html = "";
              if (type == 'display') {
                html = `<button type="button" class="btn btn-danger btn-xs font-weight-bold d-sm-inline-block shadow-md deleteRow" data-id="` + data + `" data-name="` + row.name + `" style="min-width: 50px;">Delete</button>`;
              }
              return html;
            }
          },
        ],
        initComplete: function (settings, json) {
          $(this).on('click', '.editRow', function () {
            var recordId = $(this).data('id');
            getById(recordId);
          });

          $(this).on('click', '.deleteRow', function () {
            var recordId = $(this).data('id');
            var recordName = $(this).data('name');
            Swal.fire({
              title: 'Apakah yakin ingin menghapus ' + recordName + '?',
              text: "",
              icon: 'warning',
              showCancelButton: true,
              confirmButtonColor: '#3085d6',
              cancelButtonColor: '#d33',
              confirmButtonText: 'Ya',
              cancelButtonText: 'Tidak'
            }).then((result) => {
              if (result.value) {
                deleteById(recordId, recordName);
              }
            });
          });
        }
      }, function (e, settings, json) {
        var $table = e; 
      });

      $dt.on('processing.dt', function (e, settings, processing) {
        if (processing) {
        } else {
        }
      })
    }
    
    $btnSaveOrUpdate.on("click", function(event) {
      SaveOrUpdateTenant(event);
    });

    var SaveOrUpdateTenant = function (e) {
      var isvalidate = $form[0].checkValidity();
      if (isvalidate) {
        var record = $form.serializeToJSON();
        $.ajax({
          url: $.helper.baseApiPath("/entity/save"),
          type: 'POST',
          data: record,
          success: function (r) {
            if (r.status) {
              $dt.ajax.reload();
              $form.trigger("reset");
              toastr.success("Berhasil menambah Tenant",'Information!');
              $modalForm.modal("hide");
            } else {
              toastr.error(r.message, 'Error!');
            }
          },
          error: function (r) {
            if (r.status == 400 || r.status == 409) {
              toastr.error(r.responseText,'Bad Request!');
            } else {
              toastr.error(r.responseText,'Information!');
            }
          }
        });
      } else {
        e.preventDefault();
        e.stopPropagation();
        $form.addClass('was-validated');
      }
    }

    var getById = function (id) {
      $.ajax({
        url: $.helper.baseApiPath("/entity/getById/" + id),
        type: 'GET',
        success: function (r) {
          if (r.status) {
            $form.find('input').val(function () {
                return r.data[this.name];
            });
            $("textarea[name=description]").val(r.data.description);
            $modalForm.modal("show");
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

    var deleteById = function (id, name) {
      $.ajax({
        url:  $.helper.baseApiPath("/entity/deleteById/" + id),
        type: 'DELETE',
        success: function (r) {
          if (r.status) {
            $dt.ajax.reload();
            Swal.fire('Berhasil!', name + ' berhasil dihapus', 'success');
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

    $btnNewTenant.on("click", function () {
      $form.trigger("reset");
      $modalForm.modal("show");
    });

    return {
      init: function () {
        loadDatatables();
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));