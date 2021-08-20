(function ($) {
  'use strict';
  var $recordId = $("#recordId");

  var $formDescription = $('#form-description');
  var $formPayment = $('#form-payment');
  var $formPaymentDetail = $('#form-payment-detail');
  
  var $btnSave = $("#btn-save");
  var $btnShowPaymentModal = $("#btn-showPaymentModal");
  var $btnPayment = $("#btn-payment");
  var $btnFinish = $(".btn-finish");
  var $btnToList = $(".btn-toList");

  var $modalPayment = $("#modal-payment");
  var $modalChange = $("#modal-change");

  var $OrderNumber = null;

  var cart = {};
  cart.order_detail = new Object();

  var $insufficientPayment = 0;
  var $totalPayment = 0;

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
            console.log("getById", r);

            $.helper.form.fill($formDescription, r.data);
            $.helper.form.fill($formPayment, r.data);
            $.helper.form.fill($formPaymentDetail, r.data);
            
            $insufficientPayment = r.data.insufficient_payment;
            $OrderNumber = r.data.order_number;
            $("#order-number").html($OrderNumber);
            $("#section-detailPelanggan").css("margin-top", "-45px");
            $("#lbl-customerName").text((r.data.customer.first_name).toUpperCase() + " " + r.data.customer.last_name);
            $("#lbl-custmerNoHp").text(r.data.customer.phone);

            getOrderDetailByOrderId(r.data.id);

            if (r.data.payment_status == 200) {
              $btnSave.remove();
              $btnShowPaymentModal.remove();
              $btnFinish.removeClass("d-none");
              $("#lbl-total-payment").text("Total Bayar :");
              $("input[name=insufficient_payment]").attr("disabled", "disabled");
              
              if (r.data.insufficient_payment == 0) {
                $("#section-total-bayar").addClass("d-none")
              }
            }

            if (r.data.order_status == 200) {
              $btnFinish.remove();
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
                  $("#section-satuan").removeClass("d-none");
                  $("#list-cartSatuan").find(".pre-content").addClass("d-none");
                  $("#list-cartSatuan").append(html);
                } else {
                  $("#section-kiloan").removeClass("d-none");
                  $("#list-cartKiloan").find(".pre-content").addClass("d-none");
                  $("#list-cartKiloan").append(html);
                }
              });
              initThousandSeparator();
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

    $btnShowPaymentModal.on("click", function () {
      $modalPayment.modal("show");
    });

    $btnSave.on("click", function (event) {
      SaveOrUpdate(event);
    });
    $("input[name=total_payment]").keypress(function (e) {
      if (e.which == 13) {
        $btnPayment.click();
        return false;   
      }
    });

    $btnPayment.on("click", function (e) {
      doPayment($btnPayment);
    });
    
    var doPayment = function (e) {
      var isValid = $formPayment[0].checkValidity();
      if (isValid) {
        $totalPayment = $("input[name=total_payment]").val().replace(/,/g, "");
        if (parseFloat($totalPayment) >= $insufficientPayment) {
          
        } else {
          toastr.error("Nominal yang dibayarkan tidak boleh kurang.");
          return;
        }

        var data = {
          id: $recordId.val(),
          total_payment: $totalPayment
        }

        Swal.fire({
          title: 'Apakah kamu yakin ingin melakukan Pembayaran?',
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
              url: $.helper.baseApiPath("/order/payment"),
              type: 'POST',
              data: data,
              success: function (r) {
                console.log("Response: ", r);

                if (r.status) {
                  toastr.success("Pembayaran Berhasil!");
                  $btnFinish.removeClass("d-none");
                  $btnShowPaymentModal.remove();

                  $("#txt-change").text(thousandSeparatorFromValueWithComma(($insufficientPayment - $totalPayment) * -1))
                  $modalPayment.modal("hide");
                  $modalChange.modal("show");
                }
              },
              error: function (r) {
                toastr.error(r.responseText, 'Error!');
              }
            });
          }
        });
      } else {
        e.preventDefault();
        e.stopPropagation();
        $formPayment.addClass('was-validated');
      }
    }

    $btnFinish.on("click", function () {
      UpdateStatus($btnFinish, 200);
    });

    var UpdateStatus = function (e, status) {
      var data = {
        id: $recordId.val(),
        status: status
      }

      Swal.fire({
        title: 'Apakah kamu yakin ingin menyelesaikan Pesanan?',
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
            url: $.helper.baseApiPath("/order/updateStatus"),
            type: 'POST',
            data: data,
            success: function (r) {
              console.log("Response: ", r);

              if (r.status) {
                e.remove();
                redirectToListOrder("Pesanan Selesai", "Apakah kamu ingin kembali ke Daftar Pesanan?");
              }
            },
            error: function (r) {
              toastr.error(r.responseText, 'Error!');
            }
          });
        }
      });
    }

    var renderProductToCartHtml = function (data, qty) {
      // console.log("renderProduct: ", data);
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
                                    <div class="col-7 text-right">`;
      if (data.is_unit) {
        html += `                   <span>Jumlah: </span><span class="fw-500">`+ thousandSeparatorWithoutComma(data.quantity) +`</span>`;
      } else {
        html += `                   <span class="fw-500">`+ thousandSeparatorDecimal(data.quantity) +`</span><span> Kg</span>`;  
      }
      html += `                     </div>
                                  </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>`
      return html;
    }

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

    $btnToList.on("click", function () {
      window.location.assign($.helper.basePath("/order"));
    });
    
    return {
      init: function () {
        loadDetail();
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));