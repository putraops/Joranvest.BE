(function ($) {
  'use strict';
  var $btnSave = $("#btn-save");
  var $formCustomer = $('#form-customer');
  var $formDescription = $('#form-description');
  var $formPayment = $('#form-payment');
  var $formSatuan = $('#form-satuan');
  var $formKiloan = $('#form-kiloan');
  var $modalForm = $("#modal-addNew");
  var $recordId = $("#recordId");
  var $btnAddNew = $(".btn-addNew");
  var $btnUpdateCustomer = $("#btn-UpdateCustomer");

  var $OrderNumber = null;

  var cart = new Object();
  cart.order_detail = new Object();

  var pageFunction = function () {
    toastr.options = {
      "closeButton": false,
      "debug": false,
      "newestOnTop": false,
      "progressBar": false,
      "positionClass": "toast-top-right",
      "preventDuplicates": false,
      "onclick": null,
      "showDuration": "1500",
      "hideDuration": "1500",
      "timeOut": "1500",
      "extendedTimeOut": "1500",
      "showEasing": "swing",
      "hideEasing": "linear",
      "showMethod": "fadeIn",
      "hideMethod": "fadeOut"
    }

    var initThousandSeparator = function () {
      $(".numeric-money").inputmask({
        digits: 2,
        greedy: true,
        definitions: {
          '*': {
            validator: "[0-9]"
          }
        },
        rightAlign: false
      });
      $(".numeric").inputmask({
        digits: 0,
        greedy: true,
        definitions: {
          '*': {
            validator: "[0-9]"
          }
        },
        rightAlign: false
      });
    }

    var initCustomerLookup = function () {
      var url = $.helper.baseApiPath("/customer/lookup");
      $("#customer_id").select2({
        ajax: {
          url: url,
          dataType: 'json',
          delay: 250,
          type: "GET",
          contentType: "application/json",
          data: function (params) {
            var field = JSON.stringify(["first_name", "last_name", "phone"]);
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
                        <br />
                        <span>No Hp: ` + _description + ` </span>
                      </div>`;
          return html;
        },
        cache: true,
        placeholder: "Masukkan Nama atau No Hp",
        minimumInputLength: 0,
        allowClear: true,
      });
      $("#customer_id").on('select2:select', function (e) {
        var res = e.params.data;
        $("#section-LookupPelanggan").addClass("d-none");
        $("#section-detailPelanggan").removeClass("d-none");
        $("#lbl-customerName").text(res.text);
        $("#lbl-custmerNoHp").text(res.description);
      });
    }

    var loadDetail = function () {
      if ($recordId.val() != "") {
        getById($recordId.val());
      }
    }

    var getById = function (id) {
      $.ajax({
        url: $.helper.baseApiPath("/order/getById/" + id),
        type: 'GET',
        success: function (r) {
          if (r.status) {
            console.log(r);

            $formDescription.find('input').val(function () {
              return r.data[this.name];
            });
            $formPayment.find('input').val(function () {
              return r.data[this.name];
            });
            $("textarea[name=description]").val(r.data.description);

            $("#section-LookupPelanggan").addClass("d-none");
            $("#section-detailPelanggan").removeClass("d-none");
            $("#lbl-customerName").text(r.data.customer.first_name + " " + r.data.customer.last_name);
            $("#lbl-custmerNoHp").text(r.data.customer.phone);
            var newOption = new Option(r.data.customer.first_name + " " + r.data.customer.last_name, r.data.customer_id, false, false);
            $('#customer_id').append(newOption).trigger('change');

            $OrderNumber = r.data.order_number;
            $("#order-number").html($OrderNumber);
            $("#section-order-number").removeClass("d-none");
            $("#section-detailPelanggan").css("margin-top", "-45px");

            getOrderDetailByOrderId(r.data.id);

          }
        },
        error: function (r) {
          if (r.status == 404) {
            Swal.fire({
              title: 'Data tidak ditemukan',
              text: '',
              icon: 'warning',
              showCancelButton: false,
              confirmButtonColor: '#3085d6',
              cancelButtonColor: '#d33',
              confirmButtonText: 'Kembali Ke Daftar Order',
              cancelButtonText: 'Tidak'
            }).then((result) => {
              if (result.value) {
                window.location.assign($.helper.basePath("/order"));
              }
            });
          } else {
            toastr.error(r.responseText, "Warning!");
          }
        }
      });
    }

    var getOrderDetailByOrderId = function (id) {
      $.ajax({
        url: $.helper.baseApiPath("/orderdetail/getViewAll?order_id=" + id),
        type: 'GET',
        success: function (r) {
          if (r.status) {
            console.log(r);
            if (r.data != null && r.data.length > 0) {
              $.each(r.data, function (index, value) {
                var html = renderProductToCartHtml(value, value.quantity);
                if (value.is_unit) {
                  $("#list-cartSatuan").find(".pre-content").addClass("d-none");
                  $("#list-cartSatuan").append(html);
                } else {
                  $("#list-cartKiloan").find(".pre-content").addClass("d-none");
                  $("#list-cartKiloan").append(html);
                }
              });
              calculate();
              initThousandSeparator();

              $(".input-quantity").on("keyup", function () {
                calculate();
              });
              $(".btn-deleteProduct").on("click", function () {
                $("#item-" + $(this).data("id")).remove();
                calculate();
              });
            }
          }
        },
        error: function (r) {
          if (r.status == 404) {
            Swal.fire({
              title: 'Data tidak ditemukan',
              text: '',
              icon: 'warning',
              showCancelButton: false,
              confirmButtonColor: '#3085d6',
              cancelButtonColor: '#d33',
              confirmButtonText: 'Kembali Ke Daftar Order',
              cancelButtonText: 'Tidak'
            }).then((result) => {
              if (result.value) {
                window.location.assign($.helper.basePath("/order"));
              }
            });
          } else {
            toastr.error(r.responseText, "Warning!");
          }
        }
      });
    }

    $btnAddNew.on("click", function () {
      $modalForm.modal("show");
    });

    var getProducts = function (type) {
      var url = $.helper.baseApiPath("/product/getAll");
      if (type != "" && type != null) {
        url += "?type=" + type;
      }
      $.ajax({
        url: url,
        type: 'GET',
        success: function (r) {
          if (r.status) {
            $("#list-products > #pre-content").remove();
            if (r.data.length > 0) {
              var html = loadProductHtml(r.data, type);
              if (type == "0") {
                $("#list-kiloan").html(html);
                $(".btn-addKiloan").on("click", function () {
                  addToCartById($(this).data("id"));
                });
              } else {
                $("#list-satuan").html(html);
                $(".btn-addSatuan").on("click", function () {
                  addToCartById($(this).data("id"));
                });
              }
              initThousandSeparator();
            }
          } else {
            toastr.error(r.message, "Error");
          }

          if (r.status) {} else {
            toastr.error(r.message, "Error");
          }
        },
        error: function () {
          toastr.error(r.responseText, "Error");
        }
      });
    }

    var addToCartById = function (id) {
      var qty = $("#productId-" + id).val();

      if (qty == "" || qty == null) {
        toastr.error("Jumlah tidak boleh kosong!", "Error!");
        return;
      }

      $.ajax({
        url: $.helper.baseApiPath("/product/getById/" + id),
        type: 'GET',
        success: function (r) {
          if (r.status) {
            $("#productId-" + id).val("");
            toastr.success("Berhasil menambahkan " + r.data.name);
            var html = renderProductToCartHtml(r.data, qty);
            if (r.data.is_unit) {
              $("#list-cartSatuan").find(".pre-content").addClass("d-none");
              $("#list-cartSatuan").append(html);
            } else {
              $("#list-cartKiloan").find(".pre-content").addClass("d-none");
              $("#list-cartKiloan").append(html);
            }
            calculate();
            initThousandSeparator();

            $(".input-quantity").on("keyup", function () {
              calculate();
            });
            $(".btn-deleteProduct").on("click", function () {
              $("#item-" + $(this).data("id")).remove();
              calculate();
            });
          } else {
            toastr.error("Product is not found! ", "Error!");
          }
        },
        error: function () {
          toastr.error(r.responseText, "Error!");
        }
      });
    }

    $btnSave.on("click", function (event) {
      SaveOrUpdate(event);
    });

    var SaveOrUpdate = function (e) {
      var isValid = true;
      if (cart.order_detail.length == 0) {
        toastr.error("Tambahkan Produk terlebih dahulu.", "Peringatan!");
        isValid = false;
        return;
      }

      var isValidateCustomer = $formCustomer[0].checkValidity();
      var isValidateKiloan = $formKiloan[0].checkValidity();
      var isValidateSatuan = $formSatuan[0].checkValidity();

      if (!isValidateKiloan) {
        e.preventDefault();
        e.stopPropagation();
        $formKiloan.addClass('was-validated');
      }

      if (!isValidateSatuan) {
        e.preventDefault();
        e.stopPropagation();
        $formSatuan.addClass('was-validated');
      }

      calculate();

      if (isValidateCustomer) {
        var customerRecord = $formCustomer.serializeToJSON();
        $.extend(cart, customerRecord);
      } else {
        toastr.error("Silahkan Pilih Data Pelanggan terlebih dahulu. <br/>Jika tidak tersedia daftarkan segera.", "Peringatan!");
        e.preventDefault();
        e.stopPropagation();
        $formKiloan.addClass('was-validated');
      }

      var isValidatePayment = $formPayment[0].checkValidity();
      if (isValidatePayment) {
        var paymentRecord = $formPayment.serializeToJSON();
        $.extend(cart, paymentRecord);
      } else {
        toastr.error("Silahkan Masukkan Total Bayar!", "Peringatan!");
        e.preventDefault();
        e.stopPropagation();
        $formPayment.addClass('was-validated');
      }

      var descriptionRecord = $formDescription.serializeToJSON();
      $.extend(cart, descriptionRecord);

      if (isValid && isValidateCustomer && isValidateSatuan && isValidateKiloan && isValidatePayment) {
        if ($recordId.val() != "") {
          cart.id = $recordId.val();
        }
        console.log("Request: ", cart);

        Swal.fire({
          title: 'Apakah kamu yakin ingin menyimpan?',
          text: '',
          icon: 'warning',
          showCancelButton: true,
          confirmButtonColor: '#3085d6',
          cancelButtonColor: '#d33',
          confirmButtonText: 'Ya',
          cancelButtonText: 'Tidak'
        }).then((result) => {
          if (result.value) {
            $.ajax({
              url: $.helper.baseApiPath("/order/save"),
              type: 'POST',
              data: cart,
              success: function (r) {
                console.log("Response: ", r);
                if (r.status) {
                  var title = $recordId.val() == "" ? "Berhasil membuat Pesanan baru" : "Berhasil mengubah Pesanan!";
                  var subtitle = 'Apakah kamu ingin kembali ke Daftar Pesanan?';

                  $OrderNumber = r.data.order_number;
                  if ($recordId.val() == "") {
                    $("#order-number").html($OrderNumber);
                    $("#section-order-number").removeClass("d-none");
                    $("#section-detailPelanggan").css("margin-top", "-45px");
                    history.pushState('', 'ID', location.hash.split('?')[0] + '?id=' + r.data.id);
                  }

                  $recordId.val(r.data.id)
                  redirectToListOrder(title, subtitle);

                } else {
                  console.log("result: ", r);
                  $.each(r.errors, function (index, value) {
                    if (value.indexOf("unique index 'uk_name_entity'")) {
                      toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar produk.", 'Peringatan!');
                    } else {
                      toastr.error(value, 'Error!');
                    }
                  });
                }
              },
              error: function (r) {
                toastr.error(r.responseText, 'Error!');
              }
            });
          }
        });
      }
    }

    var loadProductHtml = function (data, type) {
      var html = ``;
      $.each(data, function (key, value) {
        var initAddProductClass = value.is_unit ? "btn-addSatuan" : "btn-addKiloan";
        var badgeStatus = "";
        if ((value.category.name.toLowerCase()).indexOf('express') >= 0) {
          badgeStatus = "badge-success";
        } else {
          badgeStatus = "badge-secondary"
        }
        html += ` <div class="list-group-item flex-column align-items-start p-2" style="font-size: 11pt;">
                    <div class="d-flex w-100 justify-content-between">
                        <span class="mb-0 font-weight-bold">` + value.name + `</span>`;
        html += `       <span class="badge ` + badgeStatus + ` no-radius font-weight-normal" style="font-size: 10pt; min-width: 70px;">` + value.category.name + `</span>`;
        html += `   </div>
                    <div style="margin-top: 0px;">
                        <span class="no-radius product-price">Rp ` + thousandSeparator(value.price) + `</span>
                        <div class="row mt-1">
                            <div class="col-8 mt-0">
                                <div class="form-group mb-1">`;
        if (value.is_unit) {
          html += `   <div class="input-group input-group-sm">
                        <div class="input-group-prepend">
                          <span class="input-group-text">Jumlah</span>
                        </div>
                        <input type="text" class="form-control numeric text-right" value=""  id="productId-` + value.id + `"
                          data-inputmask="'alias': 'currency'" required="" im-insert="true" inputmode="numeric">
                      </div>`;
        } else {
          html += `   <div class="input-group input-group-sm">
                        <input type="text" class="form-control numeric text-left" data-inputmask="'alias': 'currency'"
                          required="" value="" im-insert="true" inputmode="numeric" id="productId-` + value.id + `">
                        <div class="input-group-prepend">
                            <span class="input-group-text">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Kg&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>
                        </div>
                      </div>`;
        }
        html += ` </div>
                            </div>
                            <div class="col-4 mt-0">
                                <button class="btn btn-outline-success btn-sm no-radius w-100 ` + initAddProductClass + ` " data-id="` + value.id + `">Tambah 
                                </button>
                            </div>
                        </div>
                    </div>
                  </div>`;
      });
      return html;
    }

    var renderProductToCartHtml = function (data, qty) {
      console.log("renderProduct: ", data);
      var badgeStatusTipe = "";
      var productId = "";
      var productName = "";
      var categoryName = "";
      if (data.category == undefined) {
        badgeStatusTipe = data.category_name.toLowerCase().indexOf('express') >= 0 ? "badge-success" : "badge-secondary";
        productName = data.product_name;
        categoryName = data.category_name;
        productId = data.product_id;
      } else {
        badgeStatusTipe = data.category.name.toLowerCase().indexOf('express') >= 0 ? "badge-success" : "badge-secondary";
        productName = data.name;
        categoryName = data.category.name;
        productId = data.id;
      }

      var html = ``;
      html += `<div class="list-group-item flex-column align-items-start pt-2 pl-3 pr-3 pb-2" style="font-size: 11pt;" id="item-` + data.id + `">
                <div class="d-flex w-100 justify-content-between">
                    <span class="mb-0 font-weight-bold">` + productName + `</span>
                    <div class="text-right">
                        <button type="button" class="btn btn-danger btn-circle btn-sm btn-deleteProduct" data-id="` + data.id + `" data-name="3" data-toggle="tooltip" data-placement="bottom" title="" data-original-title="Hapus">
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                </div>
                <div style="margin-top: -5px;">
                    <div class="row">
                        <div class="col-md-12 mt-0">
                            <div class="form-group row mb-1">
                                <div class="col-12">
                                    <span class="no-radius product-price">Rp ` + thousandSeparator(data.price) + `</span>
                                </div>
                                <div class="col-12">
                                  <div class="row">
                                    <div class="col-5">
                                      <span class="badge ` + badgeStatusTipe + ` no-radius mr-2 mt-1" style="font-size: 10pt; min-width: 70px; font-weight: 500;">` + categoryName + `</span>
                                    </div>
                                    <div class="col-7">
                                        <input type="text" class="text d-none" name="satuanId[]" value="` + productId + `">
                                        <input type="text" class="text d-none" name="satuanPrice[]" value="` + data.price + `">
                                        <div class="input-group input-group-sm">
                                            <div class="input-group-prepend">
                                                <span class="input-group-text">Jumlah</span>
                                            </div>
                                            <input type="text" class="form-control numeric text-right input-quantity" name="satuanQty[]" data-inputmask="'alias': 'currency'" required="" value="` + qty + `" im-insert="true" inputmode="numeric">
                                        </div>
                                    </div>
                                  </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>`
      return html;
    }

    var calculate = function () {
      var sub_total = 0;
      cart.order_detail = new Object();
      var records = [];
      var total = document.getElementsByName('satuanId[]').length;
      for (var i = 0; i < total; i++) {
        var id = document.getElementsByName('satuanId[]')[i].value;
        var qty = document.getElementsByName('satuanQty[]')[i].value == "" ? 0 : parseFloat(document.getElementsByName('satuanQty[]')[i].value.replace(/,/g, ""));
        var price = parseFloat(document.getElementsByName('satuanPrice[]')[i].value.replace(/,/g, ""));
        sub_total += (qty * price);
        records.push({
          id: "",
          product_id: id,
          quantity: qty,
          price: price,
        });
      }

      $("input[name=total_price]").val(sub_total);
      cart.detail = JSON.stringify(records);
      cart.order_detail = records;
    }

    $btnUpdateCustomer.on("click", function () {
      $('#customer_id').val(null).trigger('change');
      $("#section-LookupPelanggan").removeClass("d-none");
      $("#section-detailPelanggan").addClass("d-none");
    });

    var redirectToListOrder = function (title, subtitle) {
      Swal.fire({
        title: title,
        text: subtitle,
        icon: 'success',
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: 'Ya',
        cancelButtonText: 'Tidak'
      }).then((result) => {
        if (result.value) {
          window.location.assign($.helper.basePath("/order"));
        }
      });
    }

    return {
      init: function () {
        loadDetail();
        initCustomerLookup();
        getProducts("0");
        getProducts("1");
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));