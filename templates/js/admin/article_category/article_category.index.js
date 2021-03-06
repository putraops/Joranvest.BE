(function ($) {
  'use strict';
  var $form = $('#form-basic');
  var $modalForm = $("#modal-addNew");
  var $btnAddNew = $("#btn-addNew");
  var $btnSave = $("#btn-save");
  var $btnFilter = $("#btn-filter");
  var $tree = $("#category-tree");
  var $dt;

  var pageFunction = function () {

    var loadTree = function () {
      $tree.jstree({
          contextmenu: {
              select_node: false,
              items: customMenu
          },
          core: {
              themes: {
                  "responsive": true
              },
              check_callback: true,
              data: function (obj, callback) {
                  $.ajax({
                      type: "GET",
                      dataType: 'json',
                      url: $.helper.baseApiPath("/article_category/getTree"),
                      success: function (data) {
                        if (data != null) {
                          $("#no-data").addClass("d-none");
                          $("#section-tree").removeClass("d-none");
                          callback.call(this, data);
                        }
                      }
                  });
              }
          },
          "types": {
              "default": {
                  "icon": "fa fa-folder icon-state-warning icon-lg"
              },
              "root": {
                  "class": "css-custom"
              },
              "file": {
                  "icon": "fa fa-file icon-state-warning icon-lg"
              }
          },
          state: { "key": "id" },
          plugins: ["dnd", "state", "types", "contextmenu"]
      });
  }

  
  function customMenu(node) {
    // The default set of all items
    var items = {
        Create: {
            label: "Create",
            icon: "fa fa-plus-square-o",
            action: function (n) {
              if (node.parent == "#") {
                openModal(true, node.id);
              } else {
                toastr.error("Tidak diizinkan menambahkan lebih dari 1 level kategori.", 'Error!');
              }
            }
        },
        Edit: {
            label: "Edit",
            icon: "fa fa-pencil-square-o",
            action: function () {
              openModal(false, node.id);
            }
        },
        Delete: {
            label: "Delete",
            icon: "fa fa-trash-o",
            action: function () {
                deleteById(node.id);
            }
        }
    };
    return items;
}

    $btnSave.on("click", function (event) {
      SaveOrUpdate(event);
    });

    var SaveOrUpdate = function (e) {
      var isvalidate = $form[0].checkValidity();
      if (isvalidate) {
        var record = $form.serializeToJSON();
        $.ajax({
          url: $.helper.baseApiPath("/article_category/save"),
          type: 'POST',
          data: record,
          success: function (r) {
            if (r.status) {
              $("#no-data").addClass("d-none");
              $("#section-tree").removeClass("d-none");
              $tree.jstree("refresh");
              $form.trigger("reset");
              $modalForm.modal("hide");

              var $message = "Berhasil menambah Article Category";
              if (record.id != "") $message = "Berhasil mengubah Article Category";
              Swal.fire('Berhasil!', $message, 'success');
            } else {
              $.each(r.errors, function (index, value) {
                if (value.includes(`unique constraint "uk_name"`) || value.includes("kunci ganda")) {
                  toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar Article Category.", 'Peringatan!');
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
                toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar Article Category.", 'Peringatan!');
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
    
    var getById = function (is_new, id) {
      $.ajax({
        url: $.helper.baseApiPath("/article_category/getById/" + id),
        type: 'GET',
        success: function (r) {
          if (r.status) {
            if (!is_new) {
              $form.find('input').val(function () {
                return r.data[this.name];
              });             
            } else {
              $("#parent_id").val(r.data.id);
            }

            if (r.data.parent_id == null) {
              $("#section-category").addClass("d-none");
              $("#parent_name").text("");
            } else {
              $("#section-category").removeClass("d-none");
              $("#parent_name").text(r.data.parent_name);
            }

            $modalForm.modal("show");
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

    var deleteById = function (id) {
      Swal.fire({
        title: 'Apakah yakin ingin menghapus category?',
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
            url: $.helper.baseApiPath("/article_category/deleteById/" + id),
            type: 'DELETE',
            success: function (r) {
              if (r.status) {
                $tree.jstree("refresh");
                Swal.fire('Berhasil!', 'Category berhasil dihapus', 'success');
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
      openModal();
    });

    var openModal = function (is_new, node_id) {
      if (is_new) {
        $('#recordId').val(null).trigger('change');
        $form.trigger("reset");
        $form.removeClass('was-validated');
      }

      if (node_id == undefined) {
        $("#parent_id").val("");
        $("#section-category").addClass("d-none");
        $("#parent_name").text("-");
        $modalForm.modal("show");
      } else {
        getById(is_new, node_id);
      }
    }
    

    return {
      init: function () {
        loadTree();
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));