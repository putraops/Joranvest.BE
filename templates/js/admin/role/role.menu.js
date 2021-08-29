(function ($) {
  'use strict';
  var $recordId = $("#recordId");
  var $form = $('#form-basic');
  var $modalForm = $("#modal-addNew");
  var $btnSave = $("#btn-save");
  var $btnFilter = $("#btn-filter");
  var $tree = $("#category-tree");
  var $dt;

  var pageFunction = function () {
    var loadData = function () {
      $.handler.setLoading($('#loading-sm-template').html(), $("#role-name"));
      $.ajax({
        url: $.helper.baseApiPath("/role/getById/" + $recordId.val()),
        type: 'GET',
        success: function (r) {
          if (r.status) {
            $("#role-name").html(r.data.name);
          } else {
            toastr.error(r.status.message, "Warning!");  
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

    var loadTree = function () {
      $tree.jstree("destroy").empty();
      $tree.jstree({
        core: {
            themes: {
                "responsive": true
            },
            check_callback: true,
            data: function (obj, callback) {
                $.ajax({
                    type: "GET",
                    dataType: 'json',
                    url: $.helper.baseApiPath("/application_menu/getTreeByRoleId/" + $recordId.val()),
                    success: function (data) {
                      console.log(data);
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
        plugins: ["dnd", "types", "checkbox"], 
        checkbox : {
          "keep_selected_style" : false
        },
        checkbox: {
          three_state: true, // to avoid that fact that checking a node also check others
          whole_node: false,  // to avoid checking the box just clicking the node 
          tie_selection: false // for checking without selecting and selecting without checking
      },
      }).on("check_node.jstree uncheck_node.jstree", function (e, result) {
        console.log(result);
        var loading = `<span class="jstree-loading"><i class="jstree-icon jstree-ocl"></i></span>`;
        var data = {
            role_id: $recordId.val(),
            application_menu_id: result.node.id,
        };
        console.log(data);
        $("#" + result.node.id + ">a").append(loading);

        if (result.node.state.checked) {
          if (result.node.children.length > 0) {
            data.has_children = true;
            // data.children = result.node.children.

            console.log(JSON.stringify(result.node.children));
            data.children = JSON.stringify(result.node.children);
            //-- Do something to insert children into database
          }

          console.log(data);

          $.post($.helper.baseApiPath("/role_menu/save"), data, function (r) {
            $("#" + result.node.id + ">a>span").remove();
          });
        } else {
          if (result.node.parent == "#") {
            //-- Do something to delete children from database
            data.is_parent = true;
          } else {
            data.is_parent = false;
          }

          console.log(data);
          $.post($.helper.baseApiPath('/role_menu/deleteByRoleAndMenuId'), data, function (r) {
            $("#" + result.node.id + ">a>span").remove();
          });
        }
      });
    }

    return {
      init: function () {
        loadTree();
        if ($recordId.val() != "") {
          loadData();
        } else {
          toastr.error("404 not found!", "Warning!");  
          window.location.assign($.helper.basePath("/admin/role"));
        }
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));