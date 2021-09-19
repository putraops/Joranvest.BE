(function ($) {
  'use strict';
  var $formInformation = $('#form-information');
  var $formSpeaker = $('#form-speaker');
  var $formTime = $('#form-time');
  var $formPriceReward = $('#form-priceReward');
  var $btnSave = $("#btn-save");
  var $btnSubmit = $("#btn-submit");
  var $recordId = $("#recordId");
  var $btnNavNext = $(".btnNav-next");
  var $btnNavPrevious = $(".btnNav-previous");

  var pageFunction = function () {
    // $('#myTab a').on('click', function (event) {
    //   event.preventDefault()
    // })
    
    $(".nav").click(function() {
      // var next = $('.nav.nav-tabs > li > .active');
      // console.log(next);
      //  return false;
    });
    $btnNavNext.on("click", function (){
      var tabPane = $(this).closest('.tab-pane');
      var tabId = ($(tabPane).attr("id"));

      var form = $(this).closest('.needs-validation');
      var isvalidate = form[0].checkValidity();
      if (isvalidate) {
        var next = $("[aria-controls="+ tabId +"]").parent().next('li');
        if(next.length > 0) {
          next.find("a").removeClass("disabled");
          next.find('a').trigger('click');
        }
      } else {
        toastr.error("Silahkan Periksa kembali Form", "Peringatan!")
        form.addClass('was-validated');
      }
    });
    $btnNavPrevious.on("click", function (){
      var tabPane = $(this).closest('.tab-pane');
      var tabId = ($(tabPane).attr("id"));
      
      var next = $("[aria-controls="+ tabId +"]").parent().prev('li');
      if(next.length > 0) {
        next.find('a').trigger('click');
      }
    });
    

    // $('a[data-toggle="tab"]').on('shown.bs.tab', function (e) {
    //   var next = $('.nav-tabs > li > .active').parent().next('li');
    //   console.log(next);
    //   if(next.length){
    //     next.find('a').trigger('click');
    //   }else{
    //     //jQuery('#myTabs a:first').tab('show');
    //   }
    //   e.target // activated tab
    //   e.relatedTarget // previous tab
    // })

    $(".input-date").datepicker({
      format: 'yyyy-mm-dd',
      autoHide: true
    });

    var initWebinarCategoryLookup = function () {
      var url = $.helper.baseApiPath("/webinar_category/lookup");
      $("#webinar_category_id").select2({
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
        $("#validation-webinar_category_id").css("display", "none");
      });
    }
    
    var initWebinarLevelLookup = function () {
      $("#webinar_level").select2({
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
        placeholder: "Pilih Level",
        minimumInputLength: 0,
        allowClear: true,
      });
      $('#webinar_level').val(null).trigger('change');
    }

    var initSpeakerTypeLookup = function () {
      $("#speaker_type").select2({
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
        placeholder: "Pilih Tipe Pembicara",
        minimumInputLength: 0,
        allowClear: true,
      }).on('select2:select', function (e) {
        var value = e.params.data;
        $("#section-speaker-organization").addClass("d-none");
        $("#section-speaker-user").addClass("d-none");
        $('#organizer_organization_id').val(null).trigger('change');
        $('#organizer_user_id').val(null).trigger('change');
        
        $("#validation-organizer_user_id").css("display", "none");
        $("#validation-organization_id").css("display", "none");
        $("#organizer_organization_id").removeAttr("required");
        $("#organizer_user_id").removeAttr("required");
        
        if (value.text.toLowerCase() == "personal") {
          $("#section-speaker-user").removeClass("d-none");
          $("#organizer_user_id").attr("required", "");
          $("#validation-organizer_user_id").css("display", "block");
        } else {
          $("#section-speaker-organization").removeClass("d-none");
          $("#organizer_organization_id").attr("required", "");
          $("#validation-organization_id").css("display", "block");
        }
      });
      $('#speaker_type').val(null).trigger('change');
    }

    var initOrganizationLookup = function () {
      var url = $.helper.baseApiPath("/organization/lookup");
      $("#organizer_organization_id").select2({
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
          var html = `<div class="" style="font-size: 10pt; ">
                        <span class="fw-700">` + data.text + `</span>
                      </div>`;
          return html;
        },
        cache: true,
        placeholder: "Pilih Organisasi",
        minimumInputLength: 0,
        allowClear: true,
      }).on('select2:select', function (e) {
        $("#validation-organization_id").css("display", "none");
      });
    }

    var initApplicationUserLookup = function () {
      var url = $.helper.baseApiPath("/application_user/lookup");
      $("#organizer_user_id").select2({
        ajax: {
          url: url,
          dataType: 'json',
          delay: 250,
          type: "GET",
          contentType: "application/json",
          data: function (params) {
            var field = JSON.stringify(["first_name", "last_name"]);
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
        placeholder: "Pilih Pembicara",
        minimumInputLength: 0,
        allowClear: true,
      }).on('select2:select', function (e) {
        $("#validation-organizer_user_id").css("display", "none");
      });
    }
    

    var loadDetail = function () {
      if ($recordId.val() != "") {
        getById($recordId.val());
        loadAttachments();
      }
    }

    $btnSave.on("click", function (event) {
      var title = "Apakah yakin ingin menambah Webinar?";
      if ($recordId.val() != "") title = "Apakah yakin ingin mengubah Webinar";
      var isvalidate = $formPriceReward[0].checkValidity();
      if (isvalidate) {

        SaveOrUpdate(event);

        // Swal.fire({
        //   title: title,
        //   text: "",
        //   icon: 'warning',
        //   showCancelButton: true,
        //   confirmButtonColor: '#3085d6',
        //   cancelButtonColor: '#d33',
        //   confirmButtonText: 'Ya',
        //   cancelButtonText: 'Tidak'
        // }).then((result) => {
        //   if (result.value) {
        //   }
        // });
      } else {
        toastr.error("Silahkan Periksa kembali Form", "Peringatan!")
        event.preventDefault();
        event.stopPropagation();
        $formPriceReward.addClass('was-validated');
      }
    });

    var SaveOrUpdate = function (e) {
      $btnSave.html(`<span class="spinner-border spinner-border-sm mr-2" role="status" aria-hidden="true"></span>Loading...`);
      $btnSave.attr("disabled", "disabled");
      var record = {};
      $.extend(record, $formInformation.serializeToJSON(), $formSpeaker.serializeToJSON(), $formTime.serializeToJSON(), $formPriceReward.serializeToJSON());

      if (record.webinar_end_date == "") {
        record.webinar_end_date = record.webinar_start_date;
      }
      record.webinar_end_date += "T"+ $("#end_time").val() + ":00Z";
      record.webinar_start_date += "T"+ $("#start_time").val() + ":00Z";
      record.webinar_speaker = JSON.stringify(record.webinar_speaker);

      console.log(record);
      $.ajax({
        url: $.helper.baseApiPath("/webinar/save"),
        type: 'POST',
        data: record,
        success: function (r) {
          console.log(r);
          if (r.status) {
            var $message = "Berhasil menambah Webinar";
            if (record.id != "") $message = "Berhasil mengubah Webinar";
            Swal.fire('Berhasil!', $message, 'success');
            
            $recordId.val(r.data.id)
            history.pushState('', 'ID', location.hash.split('?')[0] + '?id=' + r.data.id);

            {
              var tabPane = $btnSave.closest('.tab-pane');
              console.log(tabPane);
              var tabId = ($(tabPane).attr("id"));
              console.log(tabId);

              var next = $("[aria-controls="+ tabId +"]").parent().next('li');
              console.log(next);
              if(next.length > 0) {
                next.find("a").removeClass("disabled");
                next.find('a').trigger('click');
              }
            }
          } else {
            $.each(r.errors, function (index, value) {
              console.log(value);
              if (value.includes(`unique constraint "uk_name"`)) {
                toastr.error(record.name + " sudah terdaftar. Silahkan cek kembali daftar webinar.", 'Peringatan!');
              } else {
                toastr.error(value, 'Error!');
              }
            });
          }
          
          $btnSave.html(`Simpan & Lanjut`);
          $btnSave.removeAttr("disabled");
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

    $btnSubmit.on("click", function (event) {
      var title = "Apakah yakin ingin submit Webinar?";
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
            Submit(event);
          }
      });
    });

    var Submit = function (e) {
      $btnSubmit.html(`<span class="spinner-border spinner-border-sm mr-2" role="status" aria-hidden="true"></span>Loading...`);
      $btnSubmit.attr("disabled", "disabled");

      $.ajax({
        url: $.helper.baseApiPath("/webinar/submit/" + $recordId.val()),
        type: 'POST',
        success: function (r) {
          console.log(r);
          if (r.status) {
            Swal.fire('Berhasil!', "Berhasil submit Webinar", 'success');
            //SubmitControlForm(true);
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
          
          $btnSubmit.html(`Submit`);
          $btnSubmit.removeAttr("disabled");
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
          
          $btnSubmit.html(`Submit`);
          $btnSubmit.removeAttr("disabled");
        }
      });
    }

    
    var getById = function (id) {
      $.ajax({
        url: $.helper.baseApiPath("/webinar/getById/" + id),
        type: 'GET',
        success: function (r) {
          console.log("getById", r);
          if (r.status) {
            $(".nav-link").removeClass("disabled");
            $formInformation.find('input').val(function () {
              return r.data[this.name];
            });
            $formSpeaker.find('input').val(function () {
              return r.data[this.name];
            });
            $formTime.find('input').val(function () {
              return r.data[this.name];
            });
            $formPriceReward.find('input').val(function () {
              return r.data[this.name];
            });


            $("textarea[name=description]").val(r.data.description);
            
            $('#webinar_level').val(r.data.webinar_level).trigger('change');
            var newOption = new Option(r.data.webinar_category.name, r.data.webinar_category_id, false, false);
            $('#webinar_category_id').append(newOption).trigger('change');
            if (r.data.organizer_organization_id == "") {
              getWebinarSpeaker(r.data.id);
              $('#speaker_type').val("Personal").trigger('change');
              $("#section-speaker-user").removeClass("d-none");
            } else {
              $('#speaker_type').val("Organisasi").trigger('change');
              $("#section-speaker-organization").removeClass("d-none");
              getOrganization(r.data.organizer_organization_id);
            }

            //-- Price Config
            if (r.data.price == 0) {
              $("#cbx-event-charge").prop("checked", true);
              $("input[name=price]").attr("readonly", "");
              $("input[name=discount]").attr("readonly", "");
            }

            //-- Min Age Config
            if (r.data.min_age == 0) {
              $("#cbx-min-age").prop("checked", true);
            }

            //-- Certificate Config
            if (r.data.is_certificate) {
              $("#is_certificate").prop('checked', true);
            }

            //-- Date Configuration
            $("input[name=webinar_start_date]").val(moment(r.data.webinar_start_date.Time).format('YYYY-MM-DD'));
            // $("input[name=webinar_last_start_date]").val('');
            // $("#date-start").datepicker('setDate', moment(r.data.webinar_first_start_date.Time).format('YYYY-MM-DD'));    
            $("#start_time").val(moment(r.data.webinar_start_date.Time).utc().format('HH:mm'));
            $("#end_time").val(moment(r.data.webinar_end_date.Time).utc().format('HH:mm'));
            if (moment(r.data.webinar_start_date.Time).utc().format('YYYY-MM-DD') == moment(r.data.webinar_end_date.Time).utc().format('YYYY-MM-DD')) {
              $("input[name=webinar_end_date]").val('');  
            } else {
              $("input[name=webinar_end_date]").val(moment(r.data.webinar_end_date.Time).utc().format('YYYY-MM-DD'));
              $("#section-range-end").removeClass("d-none");
              $('#event_range').prop('checked', true);
            }
              

            // $("#section-speaker-organization").addClass("d-none");
            // $("#section-speaker-user").addClass("d-none");
            // $('#organizer_organization_id').val(null).trigger('change');
            // $('#organizer_user_id').val(null).trigger('change');
            
            // $("#validation-organizer_user_id").css("display", "none");
            // $("#validation-organization_id").css("display", "none");
            // if (value.text.toLowerCase() == "personal") {
            //   $("#section-speaker-user").removeClass("d-none");
            // } else {
            //   $("#section-speaker-organization").removeClass("d-none");
            // }
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

    $("#cbx-min-age").change(function () {
      if ($(this).is(":checked")) {
        $("input[name=min_age]").val("0");
      } else {
        $("input[name=min_age]").val("");
      }
    });

    $("#cbx-event-charge").change(function () {
      if ($(this).is(":checked")) {
        $("input[name=price]").val("0")
        $("input[name=discount]").val("0");
        $("input[name=price]").attr("readonly", "");
        $("input[name=discount]").attr("readonly", "");
      } else {
        $("input[name=price]").val("");
        $("input[name=discount]").val("");
        $("input[name=price]").removeAttr("readonly");
        $("input[name=discount]").removeAttr("readonly");
      }
    });


    $("#event_range").change(function() {
      triggerEventRange();
    });

    var triggerEventRange = function () {
      if ($("#event_range").is(":checked")) {
        $("#section-range-end").removeClass("d-none")
      } else {
        $("#date-end").val("");
        $("#webinar_last_start_time").val("");
        $("#webinar_last_end_time").val("");
        $("#section-range-end").addClass("d-none")
      }
    }

    var getWebinarSpeaker = function (webinar_id) {
      $.ajax({
        url: $.helper.baseApiPath("/webinar_speaker/getAll?webinar_id=" + webinar_id),
        type: 'GET',
        success: function (r) {
          console.log("webinarSpeaker", r);
          if (r.status) {
            if (r.data != null && r.data.length > 0) {
              $.each(r.data, function( index, value ) {
                var newOption = new Option(value.speaker_full_name, value.speaker_id, true, true);
                $('#organizer_user_id').append(newOption).trigger('change');
              });
            }
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

    var getOrganization = function (organization_id) {
      $.ajax({
        url: $.helper.baseApiPath("/organization/getById/" + organization_id),
        type: 'GET',
        success: function (r) {
          console.log("webinarSpeaker", r);
          if (r.status && r.data != null) {
            var newOption = new Option(r.data.name, organization_id, false, false);
            $('#organizer_organization_id').append(newOption).trigger('change');
          } else {
            toastr.error(r.status.message, "Warning!");
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

    Dropzone.autoDiscover = false;
    var myDropzone = new Dropzone("div#my-dropzone", { 
      paramName: "file", 
      maxFilesize: 2, //-- 2Mb
      headers:
      {
        "Authorization": $('meta[name=x-token]').attr("content")
      },
      accept: function(file, done) {
        console.log(file);
        if (!file.type.match('\.jpeg')  && !file.type.match('\.jpg') && !file.type.match('\.png'))
                    // && !file.type.match('\.xls')
                    // && !file.type.match('\.docx')
                    // && !file.type.match('\.xlsx')
                    // && !(file.type == 'application/vnd.ms-excel')
                    // && !(file.type == 'application/msword')
                    // && !(allowedFileFormat.includes(fileExtension))
                    {
                    Swal.fire('Hanya file .jpeg / .jpg / .png yang diizinkan.', "", 'error');
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
      init: function (file) {
        console.log(file)
        var checkFile = false;
        this.on("error", function (file, response) {
            toastr.error("error upload : ", response, "Peringatan")
            //loadAttachment();
            this.removeAllFiles();
        });
        this.on("processing", function (file) {
          this.options.url = $.helper.baseApiPath("/filemaster/uploadByType/webinar/1/") + $recordId.val();
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
            var data = r.data[0];
            if (r.data != null && r.data.length > 0) {
              var img = `<img src="/`+ data.filepath +`" title="` + data.filename + `" title="` + data.filename + `" class="mr-1" style="width: 100%;"/>`;
              $("#section-cover-image").html(img);
            }
          }
        },
        error: function (r) {
          toastr.error(r.responseText, "Warning!");
        }
      });
    }

    return {
      init: function () {
        initWebinarCategoryLookup();
        initWebinarLevelLookup();
        initSpeakerTypeLookup();
        initOrganizationLookup();
        initApplicationUserLookup();
        loadDetail();
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));