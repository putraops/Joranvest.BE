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
      // var filterOrderStatus = $("#filter-order-status");
      // var filterPaymentStatus = $("#filter-payment-status");
      // var default_order = {
      //   "column": "r.updated_at",
      //   "dir": "DESC"
      // }

      var filter = [];
      // if (filterOrderStatus.val() != -1) {
      //   filter.push({
      //     "column": filterOrderStatus.data("field"),
      //     "value": filterOrderStatus.val().toString()
      //   })
      // }
      // if (filterPaymentStatus.val() != -1) {
      //   filter.push({
      //     "column": filterPaymentStatus.data("field"),
      //     "value": filterPaymentStatus.val().toString()
      //   })
      // }

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
          url: $.helper.baseApiPath("/membership/getDatatables"),
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
            data: "price",
            name: "price",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              var html = "";
              if (type === 'display') {
                return "Rp " + thousandSeparatorDecimal(data)
              }
              return html;
            }
          },
          
          {
            data: "description",
            name: "description"
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
                html += `<button type="button" class="btn btn-primary btn-xs font-weight-bold d-sm-inline-block shadow-md mr-1 editRow" data-id="` + data + `" data-name="` + row.name + `" style="min-width: 50px;">Lihat</button>`;
                html += `<button type="button" class="btn btn-danger btn-xs font-weight-bold d-sm-inline-block shadow-md deleteRow mr-1" data-id="` + data + `" data-name="` + row.name + `" style="min-width: 50px;">Hapus</button>`;
              }
              return html;
            }
          },
        ],
        initComplete: function (settings, json) {
          $(this).on('click', '.editRow', function () {
            var recordId = $(this).data('id');
            window.location.assign($.helper.basePath("/order/payment?id=") + recordId);
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
        $.ajax({
          url: $.helper.baseApiPath("/membership/save"),
          type: 'POST',
          data: record,
          success: function (r) {
            if (r.status) {
              $("#list-products > #pre-content").remove();
              $form.trigger("reset");
              $modalForm.modal("hide");

              if (record.id == "") {
                toastr.success("Berhasil menambah Membership", 'Information!');
                $dt.ajax.reload();
              } else {
                toastr.success("Berhasil mengubah Produk", 'Information!');
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
            url: $.helper.baseApiPath("/membership/deleteById/" + id),
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