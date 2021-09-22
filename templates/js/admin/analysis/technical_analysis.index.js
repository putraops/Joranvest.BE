(function ($) {
  'use strict';
  var $dtBasic = $("#dtBasic");
  var $form = $('#form-basic');
  var $btnSave = $("#btn-save");
  var $btnFilter = $("#btn-filter");
  var $dt;

  var pageFunction = function () {
    $("#filter-order-status, #filter-payment-status").select2({
      cache: true,
      placeholder: "Pilihs Status",
      minimumInputLength: 0,
      allowClear: true,
    });

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
          url: $.helper.baseApiPath("/technical_analysis/getDatatables"),
          type: "POST",
          contentType: "application/json",
          data: function (d) {
            // if (d.draw == 1) {
            //   d.default_order = default_order
            // }
            // if (filter.length > 0) {
            //   d.filter = filter;
            // }
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
            data: "emiten_name",
            name: "emiten_name",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              console.log(row);
              var html = "";
              if (type === 'display') {
                //html = `<a class="font-weight-bold" href="/order/detail?id=`+ row.id +`" style="text-decoration: none; font-size: 10pt;">`+ data +`</a>`;
                html = `<a href="/admin/technical_analysis/detail?id=` + row.id + `" class="font-weight-bold" style="font-size: 10pt;">` + data + `</a>`;
              }
              return html;
            }
          },
          {
            data: "emiten_code",
            name: "emiten_code",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              var html = "";
              if (type === 'display') {
                html = `<span class="font-weight-bold" style="font-size: 10pt;">` + data + `</span>`;
              }
              return html;
            }
          },
          {
            data: "created_by_fullname",
            name: "u1.first_name",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              var html = "";
              if (type === 'display') {
                html = `<span class="font-weight-bold" style="font-size: 10pt;">` + data + `</span>`;
              }
              return html;
            }
          },
          {
            data: "signal",
            name: "signal",
            orderable: false,
            searchable: false,
            class: "text-left",
            render: function (data, type, row) {
              var html = `<span style="min-width: 80px; font-weight: 500;" `
              if (type === 'display') {
                if (data == "Netral") html += `class="badge badge-secondary">`;
                else if (data == "Sideway") html += `class="badge badge-warning">`;
                else if (data == "Uptrend") html += `class="badge badge-success">`;
                else if (data == "Downtrend") html += `class="badge badge-danger">`;
              }
              html += data + `</span>`;
              return html;
            }
          },
          {
            data: "bandarmology_status",
            name: "bandarmology_status"
          },
          {
            data: "start_buy",
            name: "start_buy",
            orderable: false,
            searchable: false,
            class: "text-left",
            render: function (data, type, row) {
              var html = thousandSeparatorInteger(row.start_buy);
              if (type === 'display') {
                if ($.isNumeric(row.end_buy) && row.end_buy != 0 && (row.start_buy != row.end_buy)) {
                  html += " - " + thousandSeparatorInteger(row.end_buy)
                }
              }
              return html;
            }
          },
          {
            data: "start_sell",
            name: "start_sell",
            orderable: false,
            searchable: false,
            class: "text-left",
            render: function (data, type, row) {
              var html = thousandSeparatorInteger(row.start_sell);
              if (type === 'display') {
                if ($.isNumeric(row.end_sell) && row.end_sell != 0 && (row.start_sell != row.end_sell)) {
                  html += " - " + thousandSeparatorInteger(row.end_sell)
                }
              }
              return html;
            }
          },
          {
            data: "start_cut",
            name: "start_cut",
            orderable: false,
            searchable: false,
            class: "text-left",
            render: function (data, type, row) {
              var html = thousandSeparatorInteger(row.start_cut);
              if (type === 'display') {
                if ($.isNumeric(row.end_cut) && row.end_cut != 0 && (row.start_cut != row.end_cut)) {
                  html += " - " + thousandSeparatorInteger(row.end_cut)
                }
              }
              return html;
            }
          },
          {
            data: "start_ratio",
            name: "start_ratio",
            orderable: false,
            searchable: false,
            class: "text-left",
            render: function (data, type, row) {
              var html = "";
              if (type === 'display') {
                html += thousandSeparatorInteger(row.start_ratio) + " : " + thousandSeparatorInteger(row.end_ratio)
              }
              return html;
            }
          },
          {
            data: "reason_to_buy",
            name: "reason_to_buy",
            orderable: false,
            searchable: false,
            render: function (data, type, row) {
              var html = "";
              if (type === 'display') {
                html  = "<div class='text-wrap width-200'>" + data + "</div>";
              }
              return html;
            }
          },
          {
            data: "id",
            name: "id",
            orderable: false,
            searchable: false,
            class: "text-left",
            render: function (data, type, row) {
              var html = "";
              if (type == 'display') {
                // html += `<button type="button" class="btn btn-primary btn-xs font-weight-bold d-sm-inline-block shadow-md mr-1 detailRow" data-id="` + data + `" data-name="` + row.name + `" style="min-width: 50px;">Lihat</button>`;
                html += `<a href="/admin/technical_analysis/detail?id=` + row.id + `" type="button" class="btn btn-primary btn-xs font-weight-bold d-sm-inline-block shadow-md mr-1" data-id="` + data + `" data-name="` + row.name + `" style="min-width: 50px;">Ubah</a>`;
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
          url: $.helper.baseApiPath("/emiten/save"),
          type: 'POST',
          data: record,
          success: function (r) {
            console.log(r);
            if (r.status) {
              $dt.ajax.reload();
              $form.trigger("reset");
              $modalForm.modal("hide");

              if (record.id == "") {
                toastr.success("Berhasil menambah Emiten", 'Information!');
                $dt.ajax.reload();
              } else {
                toastr.success("Berhasil mengubah Emiten", 'Information!');
              }

              $('[data-toggle="tooltip"]').tooltip();
            } else {
              $.each(r.errors, function (index, value) {
                console.log(value);
                if (value.includes(`unique constraint "uk_name"`)) {
                  toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar produk.", 'Peringatan!');
                } else {
                  toastr.error(value, 'Error!');
                }
              });
            }
          },
          error: function (r) {
            var obj = JSON.parse(r.responseText);
            $.each(obj.errors, function (index, value) {
              if (value.includes("unique index 'uk_name_entity'")) {
                console.log(value);
                toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar produk 123.", 'Peringatan!');
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
        url: $.helper.baseApiPath("/emiten/getById/" + id),
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
        title: 'Apakah yakin ingin menghapus analisa teknikal ini?',
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
            url: $.helper.baseApiPath("/technical_analysis/deleteById/" + id),
            type: 'DELETE',
            success: function (r) {
              console.log(r);
              if (r.status) {
                $dt.ajax.reload();
                Swal.fire('Berhasil!', 'Berhasil menghapus Analisa Teknikal', 'success');
              }
            },
            error: function (r) {
              toastr.error(r.responseText, "Warning");
            }
          });
        }
      });
    }

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