(function ($) {
  'use strict';
  var $form = $('#form-technical-analysis');
  var $btnSave = $("#btn-save");
  var $recordId = $("#recordId");

  var pageFunction = function () {
    var initEmitenLookup = function () {
      var url = $.helper.baseApiPath("/emiten/lookup");
      $("#emiten_id").select2({
        ajax: {
          url: url,
          dataType: 'json',
          delay: 250,
          type: "GET",
          contentType: "application/json",
          data: function (params) {
            var field = JSON.stringify(["emiten_name", "emiten_code"]);
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
          var _description = data.description == undefined ? "-" : data.description;
          var html = `<div class="" style="font-size: 10pt; ">
                        <span class="fw-700">` + data.text + `</span>
                      </div>`;
          return html;
        },
        cache: true,
        placeholder: "Pilih Emiten",
        minimumInputLength: 0,
        allowClear: true,
      });
    }

    var initSignalLookup = function () {
      $("#signal").select2({
        escapeMarkup: function (markup) {
          return markup;
        },
        templateResult: function (data) {
          var html = `<div class="" style="font-size: 10pt; ">
                        <span class="fw-700">` + data.text + `</span>
                      </div>`;
          return html;
        },
        cache: false,
        placeholder: "Pilih Status",
        minimumInputLength: 0,
        allowClear: true,
      });
      $('#signal').val(null).trigger('change');
    }

    var initBandarmologyLookup = function () {
      $("#bandarmology_status").select2({
        escapeMarkup: function (markup) {
          return markup;
        },
        templateResult: function (data) {
          var html = `<div class="" style="font-size: 10pt; ">
                        <span class="fw-700">` + data.text + `</span>
                      </div>`;
          return html;
        },
        cache: false,
        placeholder: "Pilih Status",
        minimumInputLength: 0,
        allowClear: true,
      });
      $('#bandarmology_status').val(null).trigger('change');
    }

    var initTimeframeLookup = function () {
      $("#timeframe").select2({
        escapeMarkup: function (markup) {
          return markup;
        },
        templateResult: function (data) {
          var html = `<div class="" style="font-size: 10pt; ">
                        <span class="fw-700">` + data.text + `</span>
                      </div>`;
          return html;
        },
        cache: false,
        placeholder: "Pilih Status",
        minimumInputLength: 0,
        allowClear: true,
      });
      $('#timeframe').val(null).trigger('change');
    }

    var loadDetail = function () {
      if ($recordId.val() != "") {
        getById($recordId.val());
      }
    }

    $btnSave.on("click", function (event) {
      var title = "Apakah yakin ingin menambah Analisa Teknikal?";
      if ($recordId.val() != "") title = "Apakah yakin ingin mengubah Analisa Teknikal";

      
      var isvalidate = $form[0].checkValidity();
      if (isvalidate) {
        Swal.fire({
          title: title,
          text: "",
          icon: 'warning',
          showCancelButton: true,
          confirmButtonColor: '#3085d6',
          cancelButtonColor: '#d33',
          confirmButtonText: 'Ya',
          cancelButtonText: 'Tidak'
        }).then((result) => {
          if (result.value) {
            SaveOrUpdate(event);
          }
        });
      } else {
        event.preventDefault();
        event.stopPropagation();
        $form.addClass('was-validated');
      }
    });

    var SaveOrUpdate = function (e) {
      var record = $form.serializeToJSON();
      $.ajax({
        url: $.helper.baseApiPath("/technical_analysis/save"),
        type: 'POST',
        data: record,
        success: function (r) {
          console.log(r);
          if (r.status) {
            var $message = "Berhasil menambah Analisa Teknikal";
            if (record.id != "") $message = "Berhasil mengubah Analisa Teknikal";
            Swal.fire('Berhasil!', $message, 'success');
            
            $recordId.val(r.data.id)
            history.pushState('', 'ID', location.hash.split('?')[0] + '?id=' + r.data.id);
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
              toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar.", 'Peringatan!');
            } else {
              toastr.error(value, 'Error!');
            }
          });
        }
      });
    }

    
    var getById = function (id) {
      $.ajax({
        url: $.helper.baseApiPath("/technical_analysis/getById/" + id),
        type: 'GET',
        success: function (r) {
          console.log("getById", r);
          if (r.status) {
            $form.find('input').val(function () {
              return r.data[this.name];
            });
            $("textarea[name=reason_to_buy]").val(r.data.reason_to_buy);
            
            var newOption = new Option(r.data.emiten.emiten_name + " [" + r.data.emiten.emiten_code + "]", r.data.emiten_id, false, false);
            $('#emiten_id').append(newOption).trigger('change');
            $('#signal').val(r.data.signal).trigger('change');
            $('#bandarmology_status').val(r.data.bandarmology_status).trigger('change');
            $('#timeframe').val(r.data.timeframe).trigger('change');
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

    return {
      init: function () {
        initEmitenLookup();
        initSignalLookup();
        initBandarmologyLookup();
        initTimeframeLookup();
        loadDetail();
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));