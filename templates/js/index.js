(function ($) {
  'use strict';
  var pageFunction = function () {
    var login = function () {
      $.ajax({
        url: $.helper.baseApiPath("/getUsers"),
        type: "GET",
        dataType: 'json',
        contentType: 'application/json',
        success: function (r) {
          console.log(r);
        },
        error: function (r) {
          alert(r.statusText)
        }
      });
    }
    return {
      init: function () {
        login();
      }
    }
  }();

  $(document).ready(function () {
      pageFunction.init();
  });
}(jQuery)); 