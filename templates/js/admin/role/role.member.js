(function ($) {
  'use strict';
  var $recordId = $("#recordId");
  var $btnSearch = $("#btn-search");
  var $txtSearch = $("#txt-search");

  var pageFunction = function () {
    $btnSearch.on("click", function () { 
      $.handler.setLoading($('#loading-template').html(), $("#section-users"));
      loadAvailableUsers();
      // if($txtSearch.val().length >= 3) {
      // } else {
      //   toastr.error("Minimal 3 characters to search..", "Warning!");
      // }
    });

    var loadAvailableUsers = function () {
      $.ajax({
        url: $.helper.baseApiPath("/role_member/getUsersNotInRole/" + $recordId.val()),
        type: 'GET',
        success: function (r) {
          console.log(r);
          
          var temp = $.handler.setTemplate($('#user-template').html());
          if (r.status) {
            if (r.data.length > 0) {
              $("#section-users").html("");
              $.each(r.data, function( index, value ) {
                var html = $.handler.template.fill(temp, value);
                $("#section-users").append(html);
              });
              $(".addToRole").on("click", function () {
                SaveOrUpdate($(this).data("user_id"));
              });
            } else {
              $.handler.setLoading($('#userNotAvailable-template').html(), $("#section-users"));
            }
          } else {
            toastr.error(r.status.message, "Warning!");  
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

    var loadRoleMembers = function () {
      $.handler.setLoading($('#loading-template').html(), $("#section-roleMembers"));
      $.ajax({
        url: $.helper.baseApiPath("/role_member/getUsersInRole/" + $recordId.val()),
        type: 'GET',
        success: function (r) {
          console.log(r);
          
          var temp = $.handler.setTemplate($('#roleMember-template').html());
          if (r.status) {
            if (r.data.length > 0) {
              $("#section-roleMembers").html("");
              $.each(r.data, function( index, value ) {
                var html = $.handler.template.fill(temp, value);
                $("#section-roleMembers").append(html);
              });
              $(".removeFromRole").on("click", function () {
                deleteById($(this).data("user_id"));
              });
            } else {
              $.handler.setLoading($('#userNotAvailable-template').html(), $("#section-roleMembers"));
            }
          } else {
            toastr.error(r.status.message, "Warning!");  
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

    var SaveOrUpdate = function (user_id) {
      var data = {};
      data.id = "";
      data.application_user_id = user_id;
      data.role_id = $recordId.val();
      console.log(data);

      $.ajax({
        url: $.helper.baseApiPath("/role_member/save"),
        type: 'POST',
        data: data,
        success: function (r) {
          console.log(r);
          if (r.status) {
            $("#user-" + user_id).remove();
            toastr.success("Berhasil menambahkan ke dalam Role", "Success!");
            loadRoleMembers();
          } else {
            $.each(r.errors, function (index, value) {
              console.log("else: ", value);
              if (value.includes(`unique constraint "uk_name"`) || value.includes("kunci ganda")) {
                toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar role.", 'Peringatan!');
              } else {
                toastr.error(value, 'Error!');
              }
            });
          }
        },
        error: function (r) {
          var obj = JSON.parse(r.responseText);
          $.each(obj.errors, function (index, value) {
            if (value.includes("unique index 'uk_name_entity'") || value.includes("kunci ganda")) {
              console.log("error", value);
              toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar role.", 'Peringatan!');
            } else {
              toastr.error(value, 'Error!');
            }
          });
        }
      });
    }

    var deleteById = function (id) {
      $.ajax({
        url: $.helper.baseApiPath("/role_member/deleteById/" + id),
        type: 'DELETE',
        success: function (r) {
          console.log(r);
          if (r.status) {
            loadRoleMembers();
            toastr.success("Success", "Success!")
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning");
        }
      });
    }

    return {
      init: function () {
        loadRoleMembers();
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));