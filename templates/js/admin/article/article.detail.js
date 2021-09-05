(function ($) {
  'use strict';
  var $formInformation = $('#form-information');
  var $formSpeaker = $('#form-speaker');
  var $formTime = $('#form-time');
  var $formPriceReward = $('#form-priceReward');
  var $btnSave = $("#btn-save");
  var $recordId = $("#recordId");
  var $btnNavNext = $(".btnNav-next");
  var $btnNavPrevious = $(".btnNav-previous");

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

    var loadDetail = function () {
      if ($recordId.val() != "") {
        getById($recordId.val());
        loadAttachments();
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