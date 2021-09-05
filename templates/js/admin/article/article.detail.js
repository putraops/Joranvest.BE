(function ($) {
  'use strict';
  var $form = $('#form-main');
  var $btnSave = $("#btn-save");
  var $recordId = $("#recordId");

  var pageFunction = function () {
    $('.summernote').summernote({
      height: 400,   //set editable area's height
      placeholder: '',
      codemirror: { // codemirror options
        theme: 'paper'
      },
      styleTags: [
        'p',
            { title: 'Blockquote', tag: 'blockquote', className: 'blockquote', value: 'blockquote' },
            'pre', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6'
      ],
      toolbar: [
        // [groupName, [list of button]]
        ['style', ['bold', 'italic', 'underline', 'clear']],
        ['font', ['strikethrough', 'superscript', 'subscript']],
        ['fontsize', ['fontsize']],
        ['color', ['color']],
        ['para', ['ul', 'ol', 'paragraph']],
        // ['insert', ['link', 'picture']],
        ['view', ['fullscreen',// 'codeview'
        , 'help']],
        // ['height', ['height']]
      ]
    });

    var initArticleCategoryLookup = function () {
      var url = $.helper.baseApiPath("/article_category/lookup");
      $("#article_category_id").select2({
        ajax: {
          url: url,
          dataType: 'json',
          delay: 250,
          type: "GET",
          contentType: "application/json",
          data: function (params) {
            var field = JSON.stringify(["name"]);
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
          var html = `<div class="" style="font-size: 10pt;">
                        <span class="fw-700">` + data.text + `</span>
                      </div>`
          return html;
        },
        cache: true,
        placeholder: "Pilih Kategori",
        minimumInputLength: 0,
        allowClear: true,
      }).on('select2:select', function (e) {
        $("#validation-article_category_id").css("display", "none");
      });
    }

    var initTagLookup = function () {
      var url = $.helper.baseApiPath("/tag/lookup");
      $("#tagLookup").select2({
        ajax: {
          url: url,
          dataType: 'json',
          delay: 250,
          type: "GET",
          contentType: "application/json",
          data: function (params) {
            var field = JSON.stringify(["name", "description"]);
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
          var html = `<div class="" style="font-size: 10pt; ">
                        <span class="fw-700">` + data.text + `</span>
                      </div>`;
          return html;
        },
        cache: true,
        placeholder: "Pilih Tag",
        minimumInputLength: 0,
        allowClear: true,
      });
    }

    var initArticleTypeLookup = function () {
      $("#article_type").select2({
        escapeMarkup: function (markup) {
          return markup;
        },
        templateResult: function (data) {
          var html = `<div class="" style="font-size: 10pt; ">
                        <span class="fw-700">` + data.text + `</span>
                      </div>`;
          return html;
        },
        cache: false,
        placeholder: "Pilih Tipe",
        minimumInputLength: 0,
        allowClear: true,
      });
      $('#article_type').val(null).trigger('change');
    }

    $btnSave.on("click", function (event) {
      var title = "Apakah yakin ingin menambah Artikel?";
      if ($recordId.val() != "") title = "Apakah yakin ingin mengubah Artikel";

      var isvalidate = $form[0].checkValidity();
      if (isvalidate) {
        Swal.fire({
          title: title,
          text: "",
          icon: 'warning',
          showCancelButton: true,
          confirmButtonColor: '#3085d6',
          cancelButtonColor: '#d33',
          confirmButtonText: 'Ya',
          cancelButtonText: 'Tidak'
          }).then((result) => {
            if (result.value) {
              SaveOrUpdate(event);
            }
        });
      } else {
        event.preventDefault();
        event.stopPropagation();
        $form.addClass('was-validated');
      }
    });

    var SaveOrUpdate = function (e) {
      var record = $form.serializeToJSON();
      var tags = $("#tagLookup").val();
      record.tag = JSON.stringify(tags);
      console.log(record);

      $.ajax({
        url: $.helper.baseApiPath("/article/save"),
        type: 'POST',
        data: record,
        success: function (r) {
          console.log(r);
          if (r.status) {
            var $message = "Berhasil menambah Artikel. Silahkan tambah attachment jika diperlukan.";
            if (record.id != "") $message = "Berhasil mengubah Artikel";
            Swal.fire('Berhasil!', $message, 'success');
            
            $recordId.val(r.data.id)
            history.pushState('', 'ID', location.hash.split('?')[0] + '?id=' + r.data.id);
          } else {
            $.each(r.errors, function (index, value) {
              console.log(value);
              if (value.includes(`unique constraint "uk_name"`)) {
                toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar artikel.", 'Peringatan!');
              } else {
                toastr.error(value, 'Error!');
              }
            });
          }
        },
        error: function (r) {
          var obj = JSON.parse(r.responseText);
          $.each(obj.errors, function (index, value) {
            if (value.includes("unique index 'uk_name_entity'")) {
              console.log(value);
              toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar.", 'Peringatan!');
            } else {
              toastr.error(value, 'Error!');
            }
          });
        }
      });
    }

    var getById = function (id) {
      $.ajax({
        url: $.helper.baseApiPath("/article/getById/" + id),
        type: 'GET',
        success: function (r) {
          console.log("getById", r);
          if (r.status) {
            $form.find('input').val(function () {
              return r.data[this.name];
            });
            // $("textarea[name=research_data]").val(r.data.research_data);
            $(".summernote").summernote("code", r.data.body);
            var newOption = new Option(r.data.article_category.name, r.data.article_category_id, true, true);
            $('#article_category_id').append(newOption).trigger('change');
            $('#article_type').val(r.data.article_type).trigger('change');

            getTagByArticleId(r.data.id);
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

    var getTagByArticleId = function (article_id) {
      $.ajax({
        url: $.helper.baseApiPath("/article_tag/getAll?article_id=" + article_id),
        type: 'GET',
        success: function (r) {
          console.log("getTagByArticleId", r);
          if (r.status) {
            if (r.data != null && r.data.length > 0) {
              $.each(r.data, function( index, value ) {
                var newOption = new Option(value.tag_name, value.tag_id, true, true);
                $('#tagLookup').append(newOption).trigger('change');
              });
            }
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

    var loadDetail = function () {
      if ($recordId.val() != "") {
        getById($recordId.val());
      }
    }
    return {
      init: function () {
        initArticleCategoryLookup();
        initTagLookup();
        initArticleTypeLookup();
        loadDetail();
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));