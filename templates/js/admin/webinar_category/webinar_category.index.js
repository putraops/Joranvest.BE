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
              check_callback: false,
              data: function (obj, callback) {
                  $.ajax({
                      type: "GET",
                      dataType: 'json',
                      url: $.helper.baseApiPath("/webinar_category/getTree"),
                      success: function (data) {
                        console.log(data);
                        callback.call(this, data);
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

      $tree.bind("move_node.jstree rename_node.jstree", function (e, data) {
          if (e.type == "move_node") {
              // $.get($.helper.resolveApi("~/core/BusinessUnit/order"),
              //     {
              //         recordId: data.node.id,
              //         parentTargetId: data.parent,
              //         newOrder: data.old_position < data.position ? data.position + 1 : data.position
              //     }, function (data) {
              //         if (data.status.success) {
              //             //$tree.jstree(true).refresh();
              //         }
              //     });
          }
      });
  }

  
  function customMenu(node) {
    // The default set of all items
    var items = {
        Create: {
            label: "Create",
            icon: "fa fa-plus-square-o",
            action: function (n) {
              alert();
                // $.helper.form.clear($form),
                //     $.helper.form.fill($form, {
                //         parent_unit: node.id,
                //         parent_unit_name: node.text
                //     });
                // $businessUnitFormModal.modal('show');
            }
        },
        Edit: {
            label: "Edit",
            icon: "fa fa-pencil-square-o",
            action: function () {

                // $businessUnitFormModal.niftyOverlay('show'), $businessUnitFormModal.modal('show'),
                //     $.helper.form.clear($form),
                //     $.get($.helper.resolveApi('~core/BusinessUnit/' + node.id + '/detail'), function (r) {
                //         if (r.status.success) {
                //             $.helper.form.fill($form, r.data);
                //         }
                //         $businessUnitFormModal.niftyOverlay('hide');
                //     }).fail(function (r) {
                //         $.helper.noty.error(r.status, r.statusText);
                //         $businessUnitFormModal.niftyOverlay('hide');
                //     });
            }
        },
        Delete: {
            label: "Delete",
            icon: "fa fa-trash-o",
            action: function () {

                var b = bootbox.confirm({
                    message: "<p class='text-semibold text-main'>Are you sure ?</p><p>You won't be able to revert this!</p>",
                    buttons: {
                        confirm: {
                            className: "btn-danger",
                            label: "Confirm"
                        }
                    },
                    callback: function (result) {
                        if (result) {

                            $.ajax({
                                type: "POST",
                                dataType: 'json',
                                contentType: 'application/json',
                                url: $.helper.resolveApi("~/core/BusinessUnit/delete"),
                                data: JSON.stringify([node.id]),
                                success: function (r) {
                                    if (r.status.success) {
                                        $.helper.noty.success("Successfully", "Data has been deleted");
                                        $tree.jstree("refresh");
                                    } else {
                                        $.helper.noty.error("Information", r.status.message);
                                    }
                                    b.modal('hide');
                                },
                                error: function (r) {
                                    $.helper.noty.error(r.status, r.statusText);
                                    b.modal('hide');
                                }
                            });

                            return false;
                        }
                    },
                    animateIn: 'bounceIn',
                    animateOut: 'bounceOut'
                });

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
        console.log(record);
        $.ajax({
          url: $.helper.baseApiPath("/webinar_category/save"),
          type: 'POST',
          data: record,
          success: function (r) {
            console.log(r);
            if (r.status) {
              loadTree();
              $form.trigger("reset");
              $modalForm.modal("hide");

              var $message = "Berhasil menambah Webinar Category";
              if (record.id != "") $message = "Berhasil mengubah Webinar Category";
              Swal.fire('Berhasil!', $message, 'success');
            } else {
              $.each(r.errors, function (index, value) {
                console.log("else: ", value);
                if (value.includes(`unique constraint "uk_name"`) || value.includes("kunci ganda")) {
                  toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar Webinar Category.", 'Peringatan!');
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
                toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar Webinar Category.", 'Peringatan!');
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

    
    var getById = function (id) {
      $.ajax({
        url: $.helper.baseApiPath("/webinar_category/getById/" + id),
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
            url: $.helper.baseApiPath("/webinar_category/deleteById/" + id),
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
        loadTree();
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));