(function ($) {
  'use strict';
  var $form = $('#form-analysis');
  var $btnSave = $("#btn-save");
  var $recordId = $("#recordId");

  var pageFunction = function () {
    $(".input-date").datepicker({
      format: 'yyyy-mm-dd',
      autoHide: true
    });

    Dropzone.autoDiscover = false;
    var myDropzone = new Dropzone("div#my-dropzone", { 
      paramName: "file", 
      maxFilesize: 2, 
      headers:
      {
        "Authorization": $('meta[name=x-token]').attr("content")
      },
      accept: function(file, done) {
        if (!file.type.match('\.pdf')
                    // && !file.type.match('\.doc')
                    // && !file.type.match('\.xls')
                    // && !file.type.match('\.docx')
                    // && !file.type.match('\.xlsx')
                    // && !(file.type == 'application/vnd.ms-excel')
                    // && !(file.type == 'application/msword')
                    // && !(allowedFileFormat.includes(fileExtension))
                    ) {
                    // swal("Only allow pdf, doc, xls, msg, jpeg, png and gif file", {
                    //     icon: "warning",
                    // });
    
                    
                    Swal.fire('Hanya file .pdf yang diizinkan.', "", 'error');
                    myDropzone.removeAllFiles();
                    return;
                }
                if ($recordId.val() == '' || $recordId.val() == undefined) {
                    myDropzone.removeAllFiles();
                    alert('Please save activity, and try again');
                    done('Please save activity, and try again');
                }
        else { done(); }
      },
      init: function () {
        var checkFile = false;
        this.on("error", function (file, response) {
            console.log("error upload : ", response)
            //loadAttachment();
            this.removeAllFiles();
        });
        this.on("processing", function (file) {
            this.options.url = $.helper.baseApiPath("/filemaster/upload/") + $recordId.val();
        });
      }
    });
    myDropzone.on("complete", function(file) {
      loadAttachments();
      myDropzone.removeFile(file);
    });

    var loadAttachments = function () {
      $.ajax({
        url: $.helper.baseApiPath("/filemaster/getAll?record_id=" + $recordId.val()),
        type: 'GET',
        success: function (r) {
          console.log("loadAttachments", r);
          if (r.status) {
            if (r.data != null && r.data.length > 0) {
              $("#total_attachment  ").text(r.data.length);
              $.each(r.data, function( index, value ) {
                var html = "";
                if ((value.extension).toUpperCase() == ".PDF") {
                  html += `<img src="/assets/gallery/icon/iconPDF.png" title="` + value.filename + `" class="mr-1" style="width: 30px;"/>`;
                }
                $("#section-files").append(html);
              });
            }
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

   
    // Dropzone.autoDiscover = false;
    // var myDropzone = new Dropzone("div#my-dropzone", { // Make the whole body a dropzone
    //     thumbnailWidth: 50,
    //     thumbnailHeight: 50,
    //     parallelUploads: 20,
    //     // autoQueue: true, // Make sure the files aren't queued until manually added,
    //     autoProcessQueue: false,
    //     // headers:
    //     // {
    //     //     "Authorization": 'Bearer ' + $('meta[name=x-token]').attr('content'),
    //     //     "Server-Type": 'Axa.ASP.NET Web API'
    //     // },
    //     accept: function (file, done) {
    //       alert();
    //       if (file.name == "justinbieber.jpg") {
    //         done("Naha, you don't.");
    //       }
    //       else { done(); }
    //       // console.log(file);
    //       // console.log(done);
    //       //   // let allowedFileFormat = ["msg", "jpeg", "jpg", "png", "gif", "xls", "xlsx", "doc", "docx", "pdf"];
    //       //   let allowedFileFormat = ["pdf"];
    //       //   let fileExtension = file.name.split('.').pop();
    //       //   //console.log("fileExtension ", fileExtension);
    //       //   if (!file.type.match('\.pdf')
    //       //       && !file.type.match('\.doc')
    //       //       && !file.type.match('\.xls')
    //       //       && !file.type.match('\.docx')
    //       //       && !file.type.match('\.xlsx')
    //       //       && !(file.type == 'application/vnd.ms-excel')
    //       //       && !(file.type == 'application/msword')
    //       //       && !(allowedFileFormat.includes(fileExtension))) {
    //       //       // swal("Only allow pdf, doc, xls, msg, jpeg, png and gif file", {
    //       //       //     icon: "warning",
    //       //       // });

                
    //       //       Swal.fire('Hanya file .pdf yang diizinkan.', "", 'error');
    //       //       myDropzone.removeAllFiles();
    //       //       return;
    //       //   }
    //       //   if ($recordId.val() == '' || $recordId.val() == undefined) {
    //       //       myDropzone.removeAllFiles();
    //       //       alert('Please save activity, and try again');
    //       //       done('Please save activity, and try again');
    //       //   }
    //       //   done();
    //     },
    // });

   

    var initEmitenLookup = function () {
      var url = $.helper.baseApiPath("/emiten/lookup");
      $("#emiten_id").select2({
        ajax: {
          url: url,
          dataType: 'json',
          delay: 250,
          type: "GET",
          contentType: "application/json",
          data: function (params) {
            var field = JSON.stringify(["emiten_name", "emiten_code"]);
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
          var html = `<div class="" style="font-size: 10pt; ">
                        <span class="fw-700">` + data.text + `</span>
                      </div>`;
          return html;
        },
        cache: true,
        placeholder: "Pilih Emiten",
        minimumInputLength: 0,
        allowClear: true,
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

    var loadDetail = function () {
      if ($recordId.val() != "") {
        $("#section-attachments").removeClass("d-none");
        getById($recordId.val());
        loadAttachments();
      }
    }

    $btnSave.on("click", function (event) {
      var title = "Apakah yakin ingin menambah Analisa Teknikal?";
      if ($recordId.val() != "") title = "Apakah yakin ingin mengubah Analisa Teknikal";

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
      record.research_date += "T00:00:00Z";

      console.log(record);

      $.ajax({
        url: $.helper.baseApiPath("/fundamental_analysis/save"),
        type: 'POST',
        data: record,
        success: function (r) {
          console.log(r);
          if (r.status) {
            var $message = "Berhasil menambah Analisa Teknikal. Silahkan tambah attachment jika diperlukan.";
            if (record.id != "") $message = "Berhasil mengubah Analisa Teknikal";
            Swal.fire('Berhasil!', $message, 'success');
            
            $("#section-attachments").removeClass("d-none");
            $recordId.val(r.data.id)
            history.pushState('', 'ID', location.hash.split('?')[0] + '?id=' + r.data.id);
          } else {
            $.each(r.errors, function (index, value) {
              console.log(value);
              if (value.includes(`unique constraint "uk_name"`)) {
                toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar produk.", 'Peringatan!');
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
        url: $.helper.baseApiPath("/fundamental_analysis/getById/" + id),
        type: 'GET',
        success: function (r) {
          console.log("getById", r);
          if (r.status) {
            $form.find('input').val(function () {
              return r.data[this.name];
            });
            $("textarea[name=research_data]").val(r.data.research_data);
            $("input[name=research_date]").val(moment(r.data.research_date.Time).format('YYYY-MM-DD'));
            $(".input-date").datepicker('setDate', moment(r.data.research_date.Time).format('YYYY-MM-DD'));    
            
            var newOption = new Option(r.data.emiten.emiten_name + " [" + r.data.emiten.emiten_code + "]", r.data.emiten_id, true, true);
            $('#emiten_id').append(newOption).trigger('change');

            calculateMarginOfSafety();
            getFundamentalTagByFundamentalAnalysisId(r.data.id);
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

    var getFundamentalTagByFundamentalAnalysisId = function (fundamental_analysis_id) {
      $.ajax({
        url: $.helper.baseApiPath("/fundamental_analysis_tag/getAll?fundamental_analysis_id=" + fundamental_analysis_id),
        type: 'GET',
        success: function (r) {
          console.log("getFundamentalTagByFundamentalAnalysisId", r);
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

    $(".calculation").on("change", function () {
      calculateMarginOfSafety();
    });

    var calculateMarginOfSafety = function () {
      var current_price = $("input[name=current_price]").val().replace(/,/g, "");
      var normal_price = $("input[name=normal_price]").val().replace(/,/g, "");

      if (current_price == "") current_price = 0;
      if (normal_price == "") normal_price = 0;

      var margin_of_safety = ((parseInt(normal_price) - parseInt(current_price)) / parseInt(current_price)) * 100;

      if (parseInt(current_price) > 0 && parseInt(normal_price) > 0) {
        $("#lbl-marginOfSafety").text(margin_of_safety.toFixed(0) + "%");
        $("input[name=margin_of_safety]").val(margin_of_safety);

        if (parseInt(margin_of_safety) > 0) {
          $("#lbl-marginOfSafety").removeClass("normal-price");
          $("#lbl-marginOfSafety").removeClass("expensive");
          $("#lbl-marginOfSafety").addClass("cheap");
        } else if (parseInt(margin_of_safety) < 0) {
          $("#lbl-marginOfSafety").removeClass("normal-price");
          $("#lbl-marginOfSafety").removeClass("cheap");
          $("#lbl-marginOfSafety").addClass("expensive");
        }
      } else {
        $("#lbl-marginOfSafety").text("-");
        $("input[name=margin_of_safety]").val(0);
        $("#lbl-marginOfSafety").removeClass("cheap");
        $("#lbl-marginOfSafety").removeClass("expensive");
        $("#lbl-marginOfSafety").addClass("normal-price");
      }
    }

    return {
      init: function () {
        initEmitenLookup();
        initTagLookup();
        loadDetail();
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));