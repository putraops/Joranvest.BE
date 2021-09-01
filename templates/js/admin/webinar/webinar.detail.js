(function ($) {
  'use strict';
  var $form = $('#form-main');
  var $btnSave = $("#btn-save");
  var $recordId = $("#recordId");

  var pageFunction = function () {
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
            console.log(r);
            return r.data;
          },
        },
        escapeMarkup: function (markup) {
          return markup;
        },
        templateResult: function (data) {
          console.log("templateResult", data);
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
        if (value.text.toLowerCase() == "personal") {
          $("#section-speaker-user").removeClass("d-none");
        } else {
          $("#section-speaker-organization").removeClass("d-none");
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
      }
    }

    $btnSave.on("click", function (event) {
      var title = "Apakah yakin ingin menambah Webinar?";
      if ($recordId.val() != "") title = "Apakah yakin ingin mengubah Webinar";

      var isvalidate = $form[0].checkValidity();

      if ($('#speaker_type').val() != null) {
        if (($('#speaker_type').val()).toLowerCase() == "personal") {
          if ($('#organizer_user_id').val() == "") {
            isvalidate = false;
            $("#validation-organizer_user_id").css("display", "block");
          }
        } else {
          if ($('#organizer_organization_id').val() == null || $('#organizer_organization_id').val() == "") {
            isvalidate = false;
            $("#validation-organization_id").css("display", "block");
          }
        }
      }

      if (isvalidate) {
        //SaveOrUpdate(event);
        //return;

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
        toastr.error("Silahkan Periksa kembali Form", "Peringatan!")
        event.preventDefault();
        event.stopPropagation();
        $form.addClass('was-validated');
      }
    });

    var SaveOrUpdate = function (e) {
      var record = $form.serializeToJSON();
      console.log("record.webinar_first_start_date: ", record.webinar_first_start_date);
      console.log("record.webinar_first_start_time: ", $("#webinar_first_start_time").val());
      console.log("record.webinar_first_end_time: ", $("#webinar_first_end_time").val());

      if (record.webinar_first_start_date){
        if ($("#webinar_first_end_time").val() != "") {
          record.webinar_first_end_date = record.webinar_first_start_date + "T"+ $("#webinar_first_end_time").val() + ":00Z";
        }
        record.webinar_first_start_date += "T"+ $("#webinar_first_start_time").val() + ":00Z";
        if ($("#webinar_first_end_time").val() == "") {
          record.webinar_first_end_date = record.webinar_first_start_date;
        }
      }
      if (record.webinar_last_start_date){
        if ($("#webinar_last_end_time").val() != "") {
          record.webinar_last_end_date = record.webinar_last_start_date + "T"+ $("#webinar_last_end_time").val() + ":00Z";
        }
        record.webinar_last_start_date += "T"+ $("#webinar_last_start_time").val() + ":00Z";
        if ($("#webinar_last_end_time").val() == "") {
          record.webinar_last_end_date = record.webinar_last_start_date;
        }
      }
      console.log(record);

      record.webinar_speaker = JSON.stringify(record.webinar_speaker);

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
        url: $.helper.baseApiPath("/webinar/getById/" + id),
        type: 'GET',
        success: function (r) {
          console.log("getById", r);
          if (r.status) {
            $form.find('input').val(function () {
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
            $("input[name=webinar_first_start_date]").val(moment(r.data.webinar_first_start_date.Time).format('YYYY-MM-DD'));
            $("#date-start").datepicker('setDate', moment(r.data.webinar_first_start_date.Time).format('YYYY-MM-DD'));    
            $("#webinar_first_start_time").val(moment(r.data.webinar_first_start_date.Time).format('HH:mm'));
            $("#webinar_first_end_time").val(moment(r.data.webinar_first_end_date.Time).format('HH:mm'));
            if (r.data.webinar_last_start_date.Time != "0001-01-01T00:00:00Z") {
              $("#section-range-end").removeClass("d-none");
              $('#event_range').prop('checked', true);
              
              $("input[name=webinar_last_start_date]").val(moment(r.data.webinar_last_start_date.Time).format('YYYY-MM-DD'));
              $("#date-end").datepicker('setDate', moment(r.data.webinar_last_start_date.Time).format('YYYY-MM-DD'));    
              $("#webinar_last_start_time").val(moment(r.data.webinar_last_start_date.Time).format('HH:mm'));
              $("#webinar_last_end_time").val(moment(r.data.webinar_last_end_date.Time).format('HH:mm'));
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