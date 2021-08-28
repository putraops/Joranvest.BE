(function ($) {
  'use strict';
  var $recordId = $("#recordId");
  var $form = $('#form-basic');
  var $modalForm = $("#modal-addNew");
  var $btnAddNew = $("#btn-addNew");
  var $btnSave = $("#btn-save");
  var $btnSearch = $("#btn-search");

  var $txtSearch = $("#txt-search");

  var pageFunction = function () {
    $btnSearch.on("click", function () { 
      $.handler.setLoading($('#loading-template').html(), $("#section-users"));
      loadUsers();
      // if($txtSearch.val().length >= 3) {
      // } else {
      //   toastr.error("Minimal 3 characters to search..", "Warning!");
      // }
    });

    var loadUsers = function () {
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
                toastr.error("Error", "Error");
                $("#user-" + $(this).data("user_id")).remove();

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
            toastr.success("Berhasil menambahkan ke dalam Role", "Success!");
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

      return;
      var isvalidate = $form[0].checkValidity();
      if (isvalidate) {
        var record = $form.serializeToJSON();
        console.log(record);
        
      } else {
        e.preventDefault();
        e.stopPropagation();
        $form.addClass('was-validated');
      }
    }

    
    var getById = function (id) {
      $.ajax({
        url: $.helper.baseApiPath("/role/getById/" + id),
        type: 'GET',
        success: function (r) {
          if (r.status) {
            $form.find('input').val(function () {
              return r.data[this.name];
            });
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
            url: $.helper.baseApiPath("/role/deleteById/" + id),
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
      $modalForm.modal("show");
    });
    

    return {
      init: function () {
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));