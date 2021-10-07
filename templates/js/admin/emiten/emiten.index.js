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

    var initSectorLookup = function () {
      var url = $.helper.baseApiPath("/sector/lookup");
      $("#sector_id").select2({
        ajax: {
          url: url,
          dataType: 'json',
          delay: 250,
          type: "GET",
          contentType: "application/json",
          data: function (params) {
            var field = JSON.stringify(["name"]);
            var req = {
              q: params.term, // search term
              page: params.page,
              field: field
            };

            console.log(req);

            return req;
          },
          processResults: function (r) {
            return r.data;
          },
        },
        escapeMarkup: function (markup) {
          return markup;
        },
        templateResult: function (data) {
          var _description = data.description == undefined ? "-" : data.description;
          var html = `<div class="" style="font-size: 10pt; ">
                        <span class="fw-700">` + data.text + `</span>
                      </div>`;
          return html;
        },
        cache: true,
        placeholder: "Pilih Sector",
        minimumInputLength: 0,
        allowClear: true,
      });
    }

    var initEmitenCategoryLookup = function () {
      var url = $.helper.baseApiPath("/emiten_category/lookup");
      $("#emiten_category_id").select2({
        ajax: {
          url: url,
          dataType: 'json',
          delay: 250,
          type: "GET",
          contentType: "application/json",
          data: function (params) {
            var field = JSON.stringify(["name"]);
            var req = {
              q: params.term, // search term
              page: params.page,
              field: field
            };

            return req;
          },
          processResults: function (r) {
            return r.data;
          },
        },
        escapeMarkup: function (markup) {
          return markup;
        },
        templateResult: function (data) {
          var html = `<div class="" style="font-size: 10pt; ">
                        <span class="fw-700">` + data.text + `</span>
                      </div>`;
          return html;
        },
        cache: true,
        placeholder: "Pilih Kategori",
        minimumInputLength: 0,
        allowClear: true,
      });
    }

    // dropdownParent: $modalForm,
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
          url: $.helper.baseApiPath("/emiten/getDatatables"),
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
              var html = "";
              if (type === 'display') {
                //html = `<a class="font-weight-bold" href="/order/detail?id=`+ row.id +`" style="text-decoration: none; font-size: 10pt;">`+ data +`</a>`;
                html = `<span class="font-weight-bold" style="font-size: 10pt;">` + data + `</span>`;
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
              console.log(row);
              var html = "";
              if (type === 'display') {
                html = `<span class="font-weight-bold" style="font-size: 10pt;">` + data + `</span>`;
              }
              return html;
            }
          },
          {
            data: "emiten_category_name",
            name: "c.name",
            orderable: false,
            searchable: false,
            class: "text-left",
            render: function (data, type, row) {
              var html = "";
              if (type == 'display') {
                if ((row.emiten_category_name).toUpperCase() == "SYARIAH") {
                  html += `<span class="badge badge-info">` + row.emiten_category_name + `</span>`;
                } else {
                  html += `<span class="badge badge-warning">` + row.emiten_category_name + `</span>`;
                }
              }
              return html;
            }
          },
          {
            data: "sector_name",
            name: "s.name"
          },
          {
            data: "description",
            name: "r.description"
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
                html += `<button type="button" class="btn btn-primary btn-xs font-weight-bold d-sm-inline-block shadow-md mr-1 editRow" data-id="` + data + `" data-name="` + row.emiten_name + `" style="min-width: 50px;">Ubah</button>`;
                html += `<button type="button" class="btn btn-danger btn-xs font-weight-bold d-sm-inline-block shadow-md deleteRow mr-1" data-id="` + data + `" data-name="` + row.emiten_name + `" style="min-width: 50px;">Hapus</button>`;
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
          console.log(r);
          if (r.status) {
            $form.find('input').val(function () {
              return r.data[this.name];
            });
            
            var newOption = new Option(r.data.sector.name, r.data.sector_id, true, true);
            $('#sector_id').append(newOption).trigger('change');
            var newOption = new Option(r.data.emiten_category.name, r.data.emiten_category_id, true, true);
            $('#emiten_category_id').append(newOption).trigger('change');
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
            url: $.helper.baseApiPath("/emiten/deleteById/" + id),
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
      
      $('#emiten_category_id').val(null).trigger('change');
      $('#sector_id').val(null).trigger('change');
      $modalForm.modal("show");
    });
    

    return {
      init: function () {
        loadDatatables();
        initSectorLookup();
        initEmitenCategoryLookup();
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));