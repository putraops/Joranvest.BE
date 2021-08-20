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
          url: $.helper.baseApiPath("/category/getAll"),
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
            data: "description"
          },
          {
            data: "is_active",
            orderable: true,
            searchable: true,
            class: "text-center",
            render: function (data, type, row) {
              var html = "";
              if (type == 'display') {
                html = ``
                if (data == true) {
                  html = `<span class="badge badge-secondary">Aktif</span>`;
                } else {
                  html = `<span class="badge badge-danger">Tidak Aktif</span>`;
                }
              }
              return html;
            }
          },
          {
            data: "id",
            orderable: true,
            searchable: true,
            class: "text-center",
            render: function (data, type, row) {
              var html = "";
              if (type == 'display') {
                html = `
                    <button type="button" class="btn btn-warning btn-xs font-weight-bold text-dark editRow" data-id="` + data + `" style="min-width: 50px;">Edit</button>
                    <button type="button" class="btn btn-danger btn-xs font-weight-bold d-sm-inline-block shadow-md deleteRow" data-id="` + data + `" data-name="` + row.name + `" style="min-width: 50px;">Delete</button>`;
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
            deleteById(recordId, recordName);
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

    $btnSaveOrUpdate.on("click", function (event) {
      SaveOrUpdate(event);
    });

    var SaveOrUpdate = function (e) {
      var isvalidate = $form[0].checkValidity();
      if (isvalidate) {
        var record = $form.serializeToJSON();
        console.log(record);
        $.ajax({
          url: $.helper.baseApiPath("/category/save"),
          type: 'POST',
          data: record,
          success: function (r) {
            if (r.status) {
              $dt.ajax.reload();
              $form.trigger("reset");
              $modalForm.modal("hide");

              if (record.id == "") {
                toastr.success("Berhasil menambah Kategori", 'Information!');
              } else {
                toastr.success("Berhasil mengubah Kategori", 'Information!');
              }
            } else {
              if (r.message.indexOf("unique index 'idx_categories_name'") >= 0) {
                toastr.error(record.name + " sudah terdaftar.", 'Peringatan!');
              } else {
                toastr.error(r.errors, 'Error!');
              }
            }
          },
          error: function (r) {
            toastr.error(r.errors, 'Error!');
          }
        });
      } else {
        e.preventDefault();
        e.stopPropagation();
        $form.addClass('was-validated');
      }
    }

    var getById = function (id) {
      var url = $.helper.baseApiPath("/category/getById/" + id);
      $.ajax({
        url: url,
        type: 'GET',
        success: function (r) {
          if (r.status) {
            $form.find('input').val(function () {
              return r.data[this.name];
            });
            $modalForm.modal("show");
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

    var deleteById = function (id, name) {
      Swal.fire({
        title: 'Apakah yakin ingin menghapus ' + name + '?',
        text: "",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: 'Ya',
        cancelButtonText: 'Tidak'
      }).then((result) => {
        if (result.value) {
          $.ajax({
            url: $.helper.baseApiPath("/category/deleteById/" + id),
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
      });
    }

    $btnNewTenant.on("click", function () {
      $form.trigger("reset");
      $form.removeClass('was-validated')
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