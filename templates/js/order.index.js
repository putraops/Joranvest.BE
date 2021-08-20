(function ($) {
  'use strict';
  var $dtBasic = $("#dtBasic");
  var $form = $('#form-basic');
  var $modalForm = $("#modal-addNew");
  var $btnNewOrder = $("#btn-newOrder");
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
      var filterOrderStatus = $("#filter-order-status");
      var filterPaymentStatus = $("#filter-payment-status");
      var default_order = {
        "column": "r.updated_at",
        "dir": "DESC"
      }

      var filter = [];
      if (filterOrderStatus.val() != -1) {
        filter.push({
          "column": filterOrderStatus.data("field"),
          "value": filterOrderStatus.val().toString()
        })
      }
      if (filterPaymentStatus.val() != -1) {
        filter.push({
          "column": filterPaymentStatus.data("field"),
          "value": filterPaymentStatus.val().toString()
        })
      }

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
          url: $.helper.baseApiPath("/order/getDatatables"),
          type: "POST",
          contentType: "application/json",
          data: function (d) {
            if (d.draw == 1) {
              d.default_order = default_order
            }
            if (filter.length > 0) {
              d.filter = filter;
            }
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
            data: "order_number",
            name: "order_number",
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
            data: "order_status",
            name: "order_status",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              var html = "";
              if (type === 'display') {
                if (data == 1) {
                  html = `<span class="badge badge-secondary" style="min-width: 100px;">Belum dikerjakan</span>`;
                } else if (data == 2) {
                  html = `<span class="badge badge-warning"  style="min-width: 100px;">Sedang diproses</span>`;
                } else if (data == 200) {
                  html = `<span class="badge badge-success" style="min-width: 100px;">Selesai</span>`;
                } else if (data == 3) {
                  html = `<span class="badge badge-danger" style="min-width: 100px;">Batal</span>`;
                }
              }
              return html;
            }
          },
          {
            data: "payment_status",
            name: "payment_status",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              var html = "";
              if (type === 'display') {
                if (data == 200) {
                  html = `<span class="badge badge-success" style="min-width: 100px;">Lunas</span>`;
                } else {
                  html = `<span class="badge badge-secondary" style="min-width: 100px;">Belum Lunas</span>`;
                }
              }
              return html;
            }
          },
          {
            data: "customer_first_name",
            name: "first_name",
            orderable: true,
            searchable: true,
            class: "text-left ",
            render: function (data, type, row) {
              var html = "";
              if (type == 'display') {
                html = row.customer_first_name.toUpperCase() + " " + row.customer_last_name
              }
              return html;
            }
          },
          {
            data: "customer_phone",
            name: "phone",
            orderable: true,
            searchable: true,
            class: "text-left ",
            render: function (data, type, row) {
              var html = "";
              if (type == 'display') {
                html = data
              }
              return html;
            }
          },
          {
            data: "total_price",
            name: "total_price",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              var html = "";
              if (type === 'display') {
                return "Rp " + thousandSeparatorWithoutComma(data)
              }
              return html;
            }
          },
          {
            data: "total_payment",
            name: "total_payment",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              var html = "";
              if (type === 'display') {
                return "Rp " + thousandSeparatorWithoutComma(data)
              }
              return html;
            }
          },
          {
            data: "insufficient_payment",
            name: "insufficient_payment",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              var html = "";
              if (type === 'display') {
                return "Rp " + thousandSeparatorWithoutComma(row.insufficient_payment)
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
                html += `<button type="button" class="btn btn-primary btn-sm font-weight-bold d-sm-inline-block shadow-md mr-1 editRow" data-id="` + data + `" data-name="` + row.name + `" style="min-width: 50px;">Lihat</button>`;
                html += `<button type="button" class="btn btn-danger btn-sm font-weight-bold d-sm-inline-block shadow-md deleteRow mr-1" data-id="` + data + `" data-name="` + row.name + `" style="min-width: 50px;">Hapus</button>`;
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
        var $table = e; // table selector 
      });

      $dt.on('processing.dt', function (e, settings, processing) {
        if (processing) {} else {}
      })
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

    $btnNewOrder.on("click", function () {
      window.location.assign($.helper.basePath("/order/detail"));
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