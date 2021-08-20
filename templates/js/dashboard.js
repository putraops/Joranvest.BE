(function ($) {
  'use strict';
  var pageFunction = function () {
    var GetTotalOrder = function () {
      $.ajax({
        url: $.helper.baseApiPath("/order/getTotalOrder"),
        type: 'GET',
        success: function (r) {
          console.log(r);
          if (r.status){
            $("#total-order").html(r.data.TotalOrder);
            $("#total-order-inprocess").html(r.data.TotalOrderInProcess);
            $("#total-order-done").html(r.data.TotalOrderDone);
            $("#total-progress").html(r.data.Progress + "%");
            $("#progress-bar-order").attr("aria-valuenow", r.data.Progress);
            $("#progress-bar-order").css("width", r.data.Progress + "%");
          }
        },
        error: function (r) {
          console.log(r.responseText)
          if (r.status == 404) {
          } else {
          }
        }
      });
    }

    var getTotalIncome = function () {
      $.ajax({
        url: $.helper.baseApiPath("/payment/getTotalIncome"),
        type: 'GET',
        success: function (r) {
          console.log("getTotalIncome: ", r);
          if (r.status){
            $("#total-income-today").html("Rp " + thousandSeparatorFromValueWithComma(r.data.Today));
            $("#total-income-this-month").html("Rp " + thousandSeparatorFromValueWithComma(r.data.ThisMonth));
          }
        },
        error: function (r) {
          console.log(r.responseText)
          if (r.status == 404) {
          } else {
          }
        }
      });
    }

    return {
      init: function () {
        GetTotalOrder();
        getTotalIncome();
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));