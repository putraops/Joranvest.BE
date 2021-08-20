(function ($) {
  'use strict';
  var $btnRegister = $("#btn-register");
  var $form = $('#form-basic');

  var pageFunction = function () {

    $btnRegister.on("click", function (event) {
      register(event);
    });

    var register = function (e) {
      var isvalidate = $form[0].checkValidity();
      if (isvalidate) {
        var record = $form.serializeToJSON();
        console.log(record.password);
        console.log($("#retype-password").val());
        if (record.password != $("#retype-password").val()) {
          toastr.error("Password dan Ulangi Password harus sama", "Peringatan!")
          return;
        }

        $.ajax({
          url: $.helper.basePath("/user/register"),
          type: 'POST',
          data: record,
          success: function (r) {
            if (r.status) {
              window.location.assign($.helper.basePath("/user/register/success"));
            } else {
              toastr.error(r.message, 'Error!');
            }
          },
          error: function (r) {
            if (r.status == 400 || r.status == 409) {
              var obj = JSON.parse(r.responseText);
              if (obj.errors.length == 1 && obj.errors[0].indexOf("failed on the 'email' tag")) {
                toastr.error("Masukkan Format Email yang benar.", 'Information!');
              } else {
                toastr.error(r.responseText, 'Error!');
              }
            } else {
              toastr.error(r.responseText, 'Error!');
            }
          }
        });
      } else {
        e.preventDefault();
        e.stopPropagation();
        $form.addClass('was-validated');
      }

    return;
     



      if (record.id == "") {
       

        
      }
    }

    return {
      init: function () {
        console.log($.helper.baseApiPath("/employee/getById/"));
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));