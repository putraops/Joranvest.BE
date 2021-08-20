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
          url: $.helper.baseApiPath("/customer/getAll"),
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
            data: "first_name",
            orderable: false,
            searchable: false,
            class: "text-left",
            render: function (data, type, row, meta) {
              return row.first_name.toUpperCase() + " " + row.last_name
            }
          },
          {
            data: "username"
          },
          {
            data: "address"
          },
          {
            data: "phone"
          },
          {
            data: "email"
          },
          {
            data: "id",
            orderable: true,
            searchable: true,
            class: "text-center",
            render: function (data, type, row) {
              var html = "";
              if (type == 'display') {
                html = `<button type="button" class="btn btn-warning btn-xs text-dark font-weight-bold editRow" data-id="` + data + `" style="min-width: 50px;">Edit</button>`;
                html += `<button type="button" class="btn btn-danger btn-xs font-weight-bold d-sm-inline-block shadow-md ml-1 deleteRow" data-id="` + data + `" data-name="` + row.first_name.toUpperCase() + " " + row.last_name + `" style="min-width: 50px;">Delete</button>`;
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

    $btnSaveOrUpdate.on("click", function (event) {
      SaveOrUpdate(event);
    });

    var SaveOrUpdate = function (e) {
      var isvalidate = $form[0].checkValidity();
      if (isvalidate) {
        var record = $form.serializeToJSON();
        $.ajax({
          url: $.helper.baseApiPath("/customer/save"),
          type: 'POST',
          data: record,
          success: function (r) {
            if (r.status) {
              $dt.ajax.reload();
              $form.trigger("reset");
              toastr.success(record.id == "" ? "Berhasil Menambah Data Pelanggan." : "Berhasil Update Data Pelanggan.", 'Information!');
              $modalForm.modal("hide");
            } else {
              var obj = JSON.parse(r.message);
              if (obj.errors.length > 0) {
                validationResponse(obj.errors);
              } else {
                toastr.error(r.responseText, 'Error!');
              }             
            }
          },
          error: function (r) {
            var obj = JSON.parse(r.responseText);
            if (obj.errors.length > 0) {
              validationResponse(obj.errors);
            } else {
              toastr.error(r.responseText, 'Error!');
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
        url: $.helper.baseApiPath("/customer/getById/" + id),
        type: 'GET',
        success: function (r) {
          if (r.status) {
            $form.find('input').val(function () {
              return r.data[this.name];
            });
            $modalForm.modal("show");
            $form.removeClass('was-validated');
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

    var deleteById = function (id, name) {
      $.ajax({
        url: $.helper.baseApiPath("/customer/deleteById/" + id),
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

    var validationResponse = function (data) {
      $.each(data, function (index, value) {
        if (value.indexOf("validation for 'Email' failed") >= 0) {
          toastr.error("Masukkan Format Email yang benar.", 'Peringatan!');
        } else if (value.indexOf("unique index 'idx_customers_phone'") >= 0) {
            toastr.error("Nomor Hp sudah terdaftar.", 'Peringatan!');
        } else if (value.indexOf("Duplicate Email") >= 0) {
          toastr.error("Email sudah terdaftar", 'Peringatan!');
        } else {
          toastr.error(value, 'Error!');
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