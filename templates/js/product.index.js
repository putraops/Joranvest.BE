(function ($) {
  'use strict';
  var $btnAddNew = $("#btn-addNew");
  var $btnSave = $("#btn-save");
  var $form = $('#form-basic');
  var $formFilter = $('#form-filter');
  var $modalForm = $("#modal-addNew");
  var $recordId = $("#recordId");
  var $btnEditProduct = $(".btn-editProduct");

  var pageFunction = function () {
    var initCategoryLookup = function () {
      var url = $.helper.baseApiPath("/category/lookup");
      $("#categoryId").select2({
        ajax: {
          url: url,
          dataType: 'json',
          delay: 250,
          type: "GET",
          contentType: "application/json",
          data: function (params) {
            return {
              q: params.term, // search term
              page: params.page
            };
          },
          escapeMarkup: function (markup) {
            return markup;
          },
          templateResult: function (data) {
            return data.html;
          },
          processResults: function (r) {
            return r.data
          },
          cache: true,
        },
        placeholder: "Pilih Kategori",
        minimumInputLength: 0,
        dropdownParent: $modalForm,
        allowClear: true,
      });
      $("#categoryIdFilter").select2({
        ajax: {
          url: url,
          dataType: 'json',
          delay: 250,
          type: "GET",
          contentType: "application/json",
          data: function (params) {
            return {
              q: params.term, // search term
              page: params.page
            };
          },
          escapeMarkup: function (markup) {
            return markup;
          },
          templateResult: function (data) {
            return data.html;
          },
          processResults: function (r) {
            return r.data
          },
          cache: true,
        },
        placeholder: "Pilih Kategori",
        minimumInputLength: 0,
        allowClear: true,
      });
    }

    $("#btn-filter").on("click", function () {
      var record = $formFilter.serializeToJSON();
      if (!jQuery.isEmptyObject(record)) {
        getAll(record);
      }
    })

    var getAll = function (param) {
      var url = $.helper.baseApiPath("/product/getAll");
      if (param != null) {
        url += "?";
        if (param.category_id != undefined) {
          url += "category_id=" + param.category_id + "&"
        }
        if (param.product_type != undefined) {
          url += "product_type=" + param.product_type + "&"
        }
      }

      $.ajax({
        url: url,
        type: 'GET',
        success: function (r) {
          if (r.status) {
            $("#list-products > #pre-content").remove();
            if (r.data.length == 0) {
              loadProductByData(null);
            } else if (r.data.length > 0) {
              var html = loadProductByData(r.data);
              $("#list-products").html(html);
              $('[data-toggle="tooltip"]').tooltip();
              $(".btn-editProduct").on("click", function () {
                getById($(this).data("id"));
              });
              $(".btn-deleteProduct").on("click", function () {
                deleteById($(this).data("id"), $(this).data("name"));
              });
            }
          } else {
            toastr.error(r.message, "Error");
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }

      });
    }

    var getById = function (id) {
      $.ajax({
        url: $.helper.baseApiPath("/product/getById/" + id),
        type: 'GET',
        success: function (r) {
          if (r.status) {
            console.log(r);
            $form.find('input').val(function () {
              return r.data[this.name];
            });

            if (r.data.is_unit == false) {
              $('#satuanRadio').prop("checked", false);
              $('#kiloanRadio').prop("checked", true);
            } else {
              $('#satuanRadio').prop("checked", true);
              $('#kiloanRadio').prop("checked", false);
            }
            var newOption = new Option(r.data.category.name, r.data.category.id, false, false);
            $('#categoryId').append(newOption).trigger('change');
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
            url: $.helper.baseApiPath("/product/deleteById/" + id),
            type: 'DELETE',
            success: function (r) {
              if (r.status) {
                getAll(null);
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
      $('#categoryId').val(null).trigger('change');
      $form.trigger("reset");
      $form.removeClass('was-validated');
      $modalForm.modal("show");
    });

    $btnSave.on("click", function (event) {
      SaveOrUpdate(event);
    });

    var SaveOrUpdate = function (e) {
      var isvalidate = $form[0].checkValidity();
      if (isvalidate) {
        var record = $form.serializeToJSON();
        $.ajax({
          url: $.helper.baseApiPath("/product/save"),
          type: 'POST',
          data: record,
          success: function (r) {
            if (r.status) {
              $("#list-products > #pre-content").remove();
              $form.trigger("reset");
              $modalForm.modal("hide");

              if (record.id == "") {
                toastr.success("Berhasil menambah Produk", 'Information!');
                var datas = new Array();
                datas.push(r.data);
                var html = loadProductByData(datas);
                $("#list-products").append(html);
                $(".btn-editProduct").on("click", function () {
                  getById($(this).data("id"));
                });
                $(".btn-deleteProduct").on("click", function () {
                  deleteById($(this).data("id"), $(this).data("name"));
                });
              } else {
                toastr.success("Berhasil mengubah Produk", 'Information!');
                $("#product-" + r.data.id).find(".product-name").text(r.data.name);
                $("#product-" + r.data.id).find(".product-description").text(r.data.description == "" ? "-" : r.data.description);
                $("#product-" + r.data.id).find(".product-price").text(thousandSeparator(r.data.price));
                $("#product-" + r.data.id).find(".product-is_unit").text(r.data.is_unit == "true" ? "Satuan" : "Kiloan");
                $("#product-" + r.data.id).find(".product-category_name").text(r.data.category.name);
              }

              $('[data-toggle="tooltip"]').tooltip();
            } else {
              $.each(r.errors, function (index, value) {
                if (value.indexOf(`unique constraint "uk_product_name"`)) {
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
              if (value.indexOf("unique index 'uk_name_entity'")) {
                toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar produk.", 'Peringatan!');
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

    $('#filter-type').select2({
      placeholder: "Pilih Type",
      minimumInputLength: 0,
      allowClear: true,
    }).select2('val', '3');

    $('input[type=radio][name=is_unit]').change(function () {
      if (this.value == "1") {
        $("#label-kategori-price").text("satuan");
      } else {
        $("#label-kategori-price").text("kilo");
      }
    });

    var loadProductByData = function (data) {
      var html = ``;
      if (data == null) {
        var html = `
        <li class="list-group-item d-flex p-0 justify-content-between lh-condensed pre-content" style="height: 80px;" id="pre-content">
            <div class="vertical-center w-100 text-center">
                <span class="font-weight-bold text-dark">Tidak ada produk tersedia.</span>
            </div>
        </li>`;
        $("#list-products").html(html);
      } else {
        $.each(data, function (key, value) {
          var badgeStatusTipe = value.category.name.toLowerCase().indexOf('express') >= 0 ? "badge-success" : "badge-secondary";
          var badgeCategory = value.is_unit ? "badge-danger" : "badge-warning";

          html += `<li class="list-group-item d-flex justify-content-between lh-condensed pt-2 pr-2 pb-2 pl-3" id="product-` + value.id + `">
                        <div class="w-80">
                            <h6 class="my-0 font-weight-bold text-dark product-name">` + value.name + `</h6>
                            <div>
                              <span class="badge pl-0 pr-0">Rp <span class="product-price">` + thousandSeparator(value.price) + `</span></span> <br />
                              <span class="badge ` + badgeCategory + ` no-radius" style="min-width: 70px;"><span class="product-is_unit">` + (value.is_unit ? "Satuan" : "Kiloan") + `</span></span>
                              <span class="badge ` + badgeStatusTipe + ` no-radius" style="min-width: 70px;"><span class="product-category_name">` + value.category.name + `</span></span>
                            </div>
                        </div>
                        <div>
                            <div class="text-right mb-1">
                                <button type="button" class="btn btn-warning btn-circle btn-sm mb-1 btn-editProduct" data-id="` + value.id + `"  data-name="` + value.name + `" data-toggle="tooltip" data-placement="bottom" title="Ubah">
                                    <i class="fas fa-pencil-alt"></i>
                                </button>
                                <button type="button" class="btn btn-danger btn-circle btn-sm mb-1 btn-deleteProduct" data-id="` + value.id + `" data-name="` + value.name + `" data-toggle="tooltip" data-placement="bottom" title="Hapus">
                                    <i class="fas fa-trash"></i>
                                </button>
                            </div>
                        </div>
                    </li>`;
        });
        return html;
      }
      return false;
    }

    return {
      init: function () {
        initCategoryLookup();
        getAll(null);
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));