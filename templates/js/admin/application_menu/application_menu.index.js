(function ($) {
  'use strict';
  var $dtBasic = $("#dtBasic");
  var $form = $('#form-basic');
  var $modalForm = $("#modal-addNew");
  var $btnAddNew = $("#btn-addNew");
  var $btnSave = $("#btn-save");
  var $btnFilter = $("#btn-filter");
  var $dt;

  var pageFunction = function () {
    $btnFilter.on("click", function () {
      loadDatatables();
    });

    var loadDatatables = function () {
      var filter = [];

      $dt = $dtBasic.DataTable({
        destroy: true,
        serverSide: true,
        pageLength: 10,
        pagingType: "full_numbers",
        responsive: true,
        processing: true,
        language: {
          processing: `<div class="spinner-border text-primary" role="status">
                        <span class="sr-only">Loading...</span>
                       </div>`
        },
        ajax: {
          url: $.helper.baseApiPath("/application_menu_category/getDatatables"),
          type: "POST",
          contentType: "application/json",
          data: function (d) {
            console.log(d);
            return JSON.stringify(d);
          }
        },
        columns: [
          // {
          //   data: "id",
          //   orderable: false,
          //   searchable: false,
          //   class: "text-center",
          //   render: function (data, type, row, meta) {
          //     return meta.row + meta.settings._iDisplayStart + 1;
          //   }
          // },
          {
            data: "name",
            name: "name",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              var html = "";
              if (type === 'display') {
                //html = `<a class="font-weight-bold" href="/order/detail?id=`+ row.id +`" style="text-decoration: none; font-size: 10pt;">`+ data +`</a>`;
                html = `<span class="font-weight-bold" style="font-size: 10pt;">` + data + `</span>`;
              }
              return html;
            }
          },
          {
            data: "id",
            name: "id",
            orderable: false,
            searchable: false,
            class: "text-center",
            render: function (data, type, row) {
              var html = "";
              if (type == 'display') {
                html += `<a href="/admin/application_menu/detail?id=` + row.id + `" type="button" class="btn btn-primary btn-xs font-weight-bold d-sm-inline-block shadow-md mr-1" data-id="` + data + `" data-name="` + row.name + `" style="min-width: 50px;">Detail</a>`;
                html += `<button type="button" class="btn btn-warning btn-xs font-weight-bold d-sm-inline-block shadow-md mr-1 editRow" data-id="` + data + `" data-name="` + row.name + `" style="min-width: 50px;">Ubah</button>`;
                html += `<button type="button" class="btn btn-danger btn-xs font-weight-bold d-sm-inline-block shadow-md deleteRow mr-1" data-id="` + data + `" data-name="` + row.name + `" style="min-width: 50px;">Hapus</button>`;
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
        var $table = e; // table selector 
      });

      $dt.on('processing.dt', function (e, settings, processing) {
        if (processing) {} else {}
      })
    }

    $btnSave.on("click", function (event) {
      SaveOrUpdate(event);
    });

    var SaveOrUpdate = function (e) {
      var isvalidate = $form[0].checkValidity();
      if (isvalidate) {
        var record = $form.serializeToJSON();
        console.log(record);
        $.ajax({
          url: $.helper.baseApiPath("/application_menu_category/save"),
          type: 'POST',
          data: record,
          success: function (r) {
            console.log(r);
            if (r.status) {
              $dt.ajax.reload();
              $form.trigger("reset");
              $modalForm.modal("hide");

              var $message = "Berhasil menambah Emiten Category";
              if (record.id != "") $message = "Berhasil mengubah Emiten Category";
              Swal.fire('Berhasil!', $message, 'success');

              $('[data-toggle="tooltip"]').tooltip();
            } else {
              $.each(r.errors, function (index, value) {
                console.log("else: ", value);
                if (value.includes(`unique constraint "uk_name"`) || value.includes("kunci ganda")) {
                  toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar Emiten Category.", 'Peringatan!');
                } else {
                  toastr.error(value, 'Error!');
                }
              });
            }
          },
          error: function (r) {
            var obj = JSON.parse(r.responseText);
            $.each(obj.errors, function (index, value) {
              if (value.includes("unique index 'uk_name_entity'") || value.includes("kunci ganda")) {
                console.log("error", value);
                toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar Emiten Category.", 'Peringatan!');
              } else {
                toastr.error(value, 'Error!');
              }
            });
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
        url: $.helper.baseApiPath("/application_menu_category/getById/" + id),
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
            url: $.helper.baseApiPath("/application_menu_category/deleteById/" + id),
            type: 'DELETE',
            success: function (r) {
              console.log(r);
              if (r.status) {
                $dt.ajax.reload();
                Swal.fire('Berhasil!', name + ' berhasil dihapus', 'success');
              }
            },
            error: function (r) {
              toastr.error(r.responseText, "Warning");
            }
          });
        }
      });
    }

    $btnAddNew.on("click", function () {
      $('#recordId').val(null).trigger('change');
      $form.trigger("reset");
      $form.removeClass('was-validated');
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