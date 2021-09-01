(function ($) {
  'use strict';
  var $dtBasic = $("#dtBasic");
  var $form = $('#form-basic');
  var $btnSave = $("#btn-save");
  var $btnFilter = $("#btn-filter");
  var $dt;

  var pageFunction = function () {
    $("#filter-order-status, #filter-payment-status").select2({
      cache: true,
      placeholder: "Pilihs Status",
      minimumInputLength: 0,
      allowClear: true,
    });

    $btnFilter.on("click", function () {
      loadDatatables();
    });

    var loadDatatables = function () {
      var filter = [];

      $dt = $dtBasic.DataTable({
        destroy: true,
        serverSide: true,
        pageLength: 10,
        pagingType: "full_numbers",
        responsive: true,
        processing: true,
        language: {
          processing: `<div class="spinner-border text-primary" role="status">
                        <span class="sr-only">Loading...</span>
                       </div>`
        },
        ajax: {
          url: $.helper.baseApiPath("/webinar/getDatatables"),
          type: "POST",
          contentType: "application/json",
          data: function (d) {
            // if (d.draw == 1) {
            //   d.default_order = default_order
            // }
            // if (filter.length > 0) {
            //   d.filter = filter;
            // }
            console.log(d);
            return JSON.stringify(d);
          }
        },
        columns: [
          // {
          //   data: "id",
          //   orderable: false,
          //   searchable: false,
          //   class: "text-center",
          //   render: function (data, type, row, meta) {
          //     return meta.row + meta.settings._iDisplayStart + 1;
          //   }
          // },


          {
            data: "title",
            name: "title",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              var html = "";
              if (type === 'display') {
                //html = `<a class="font-weight-bold" href="/order/detail?id=`+ row.id +`" style="text-decoration: none; font-size: 10pt;">`+ data +`</a>`;
                html = `<a href="/admin/webinar/detail?id=` + row.id + `" class="font-weight-bold" style="font-size: 10pt;">` + data + `</a>`;
              }
              return html;
            }
          },
          
          {
            data: "description",
            name: "r.description"
          },
          {
            data: "webinar_level",
            name: "webinar_level"
          },
          {
            data: "webinar_first_start_date",
            name: "webinar_first_start_date",
            orderable: true,
            searchable: true,
            class: "text-left min-date-width",
            render: function (data, type, row) {
              var html = `<div class='text-wrap'>`;
              if (type === 'display') {
                html  += moment(row.webinar_first_start_date.Time).format('MMM DD, YYYY');
                html  += "<br/><span>Jam: " + moment(row.webinar_first_start_date.Time).format('HH:mm') + "</span>";
                if (row.webinar_first_start_date.Time != row.webinar_first_end_date.Time) {
                  html  += "<span> - " + moment(row.webinar_first_end_date.Time).format('HH:mm') + "</span>";
                }
                if (row.webinar_last_start_date.Time != "0001-01-01T00:00:00Z") {
                  html  += "<br /><span>" + moment(row.webinar_last_start_date.Time).format('MMM DD, YYYY') + "</span>";
                  html  += "<br/><span>Jam: " + moment(row.webinar_last_start_date.Time).format('HH:mm') + "</span>";
                  if (row.webinar_last_start_date.Time != row.webinar_last_end_date.Time) {
                    html  += "<span> - " + moment(row.webinar_first_end_date.Time).format('HH:mm') + "</span>";
                  }
                }

                html += "</div>"
              }
              return html;
            }
          },
          {
            data: "organizer_organization_name",
            name: "organizer_organization_name",
            orderable: false,
            searchable: false,
            render: function (data, type, row) {
              console.log(row);
              var html = "";
              var speaker = "";
              if (type === 'display') {
                if(row.organizer_organization_id != "") {
                  speaker = row.organizer_organization_name;
                } else {
                  speaker = row.speaker_name;
                }
                html =  `<span class="font-weight-bold" style="font-size: 10pt;">` + speaker + `</span>`;
              }
              return html;
            }
          },
          {
            data: "webinar_level",
            name: "webinar_level"
          },
          {
            data: "min_age",
            name: "min_age",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              var html = "";
              if (type === 'display') {
                if (data == 0){
                  html = `<span class="font-weight-bold" style="font-size: 10pt;">Semua Umur</span>`;
                } else {
                  html = `<span class="font-weight-bold" style="font-size: 10pt;">` + data + `</span>`;
                }
              }
              return html;
            }
          },          
          {
            data: "price",
            name: "price",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              var html = thousandSeparatorInteger(data);
              return html;
            }
          },
          {
            data: "discount",
            name: "discount",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              var html = thousandSeparatorInteger(data);
              return html;
            }
          },
          {
            data: "is_certificate",
            name: "is_certificate",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              var html = thousandSeparatorInteger(data);
              if (data) {
                html = `<span class="badge badge-success font-weight-bold no-radius" style="min-width: 80px; font-weight: 500; font-size: 12px;" >Ya</span>`;
              } else {
                html = `<span class="badge badge-warning font-weight-bold no-radius" style="min-width: 80px; font-weight: 500; font-size: 12px;" >Tidak</span>`;
              }
              return html;
            }
          },
          {
            data: "reward",
            name: "reward",
            orderable: true,
            searchable: true,
            class: "text-left",
            render: function (data, type, row) {
              var html = thousandSeparatorInteger(data);
              return html;
            }
          },          
          {
            data: "id",
            name: "id",
            orderable: false,
            searchable: false,
            class: "text-left",
            render: function (data, type, row) {
              var html = "";
              if (type == 'display') {
                // html += `<button type="button" class="btn btn-primary btn-xs font-weight-bold d-sm-inline-block shadow-md mr-1 detailRow" data-id="` + data + `" data-name="` + row.name + `" style="min-width: 50px;">Lihat</button>`;
                html += `<a href="/admin/webinar/detail?id=` + row.id + `" type="button" class="btn btn-primary btn-xs font-weight-bold d-sm-inline-block shadow-md mr-1" data-id="` + data + `" data-name="` + row.name + `" style="min-width: 50px;">Ubah</a>`;
                html += `<button type="button" class="btn btn-danger btn-xs font-weight-bold d-sm-inline-block shadow-md deleteRow mr-1" data-id="` + data + `" data-name="` + row.name + `" style="min-width: 50px;">Hapus</button>`;
              }
              return html;
            }
          },
        ],
        initComplete: function (settings, json) {
          $(this).on('click', '.editRow', function () {
            var recordId = $(this).data('id');
            getById(recordId);
          });

          $(this).on('click', '.deleteRow', function () {
            var recordId = $(this).data('id');
            var recordName = $(this).data('name');
            
            deleteById(recordId, recordName);
          });
        }
      }, function (e, settings, json) {
        var $table = e; // table selector 
      });

      $dt.on('processing.dt', function (e, settings, processing) {
        if (processing) {} else {}
      })
    }

    var deleteById = function (id, name) {
      Swal.fire({
        title: 'Apakah yakin ingin menghapus webinar ini?',
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
            url: $.helper.baseApiPath("/webinar/deleteById/" + id),
            type: 'DELETE',
            success: function (r) {
              console.log(r);
              if (r.status) {
                $dt.ajax.reload();
                Swal.fire('Berhasil!', 'Berhasil menghapus Webinar', 'success');
              }
            },
            error: function (r) {
              toastr.error(r.responseText, "Warning");
            }
          });
        }
      });
    }

    return {
      init: function () {
        loadDatatables();
      }
    }
  }();

  $(document).ready(function () {
    pageFunction.init();
  });
}(jQuery));