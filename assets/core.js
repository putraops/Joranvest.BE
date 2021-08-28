// function dateFormat(date) {
//     return moment(date).format('MMM DD, YYYY');;
// }
// function timeFormat(date) {
//     return moment(date).format('HH:mm');;
// }

function thousandSeparator(nStr) {
    nStr += '';
    var x = nStr.split('.');
    var x1 = x[0];
    var x2 = x.length > 1 ? '.' + x[1] : '';
    var rgx = /(\d+)(\d{3})/;
    while (rgx.test(x1)) {
        x1 = x1.replace(rgx, '$1' + ',' + '$2');
    }
    return (x1 + x2) + ".00";
}

function thousandSeparatorDecimal(nStr) {
    nStr = (nStr) + ''; //Math.ceil(nStr);;
    var x = nStr.split('.');
    var x1 = x[0];
    var x2 = x.length > 1 ? '.' + x[1] : '.00';
    var rgx = /(\d+)(\d{3})/;
    while (rgx.test(x1)) {
        x1 = x1.replace(rgx, '$1' + ',' + '$2');
    }
    if (x2.length == 2) {
        return x1 + x2 + "0" 
    }
    return x1 + x2;
}

function thousandSeparatorFromValueWithComma(nStr) {
    nStr = Math.ceil(nStr) + ''; //Math.ceil(nStr);;
    var x = nStr.split('.');
    var x1 = x[0];
    var x2 = x.length > 1 ? '.' + x[1] : '';
    var rgx = /(\d+)(\d{3})/;
    while (rgx.test(x1)) {
        x1 = x1.replace(rgx, '$1' + ',' + '$2');
    }
    return x1 + x2;
}

function thousandSeparatorInteger(nStr) {
    nStr += '';
    var x = nStr.split('.');
    var x1 = x[0];
    var x2 = x.length > 1 ? '.' + x[1] : '';
    var rgx = /(\d+)(\d{3})/;
    while (rgx.test(x1)) {
        x1 = x1.replace(rgx, '$1' + ',' + '$2');
    }
    return (x1 + x2);
}

(function ($) {
    var hostName = $('meta[name=hostName]').attr("content");
    var routes = {
        basePath: function (param) {
            return hostName + param;
        },
        baseApiPath: function (param) {
            return hostName + "/api" + param;
        }
    }

    var coreHelper = {
        getJqueryObject: function (obj) {
            return (typeof obj === 'string') ? $('#' + obj) : (obj.jquery ? obj : $(obj));
        },
        toDefault: function (data, defaultValue) {
            return typeof data === 'undefined' || data === null || data === '' ? defaultValue : data;
        },
    }

    var controlHelper = {
        form : {
            fill: function (form, viewData, excludes) {
                if (typeof viewData === 'undefined' || viewData == null) return;

                excludes = typeof excludes === 'undefined' ? [] : excludes;
                var $el = coreHelper.getJqueryObject(form);
                var formEls = $el.find('input, select, textarea, label, span, h1, h2, h3, h4, h5, h6');
                var self = this,
                    tempObj = {},
                    type = '',
                    field = '';

                    
                $.each(formEls, function (key, obj) {
                    try {
                        if ($.inArray(obj.id, excludes) == -1 || $.inArray(obj.name, excludes) == -1) {
                            tempObj = $(obj);
                            field = typeof tempObj.attr('data-field') !== 'undefined' ? tempObj.data('field') : null;
                            //type = obj.type || typeof obj;
                            type = obj.type !== undefined ? obj.type : tempObj.attr('data-objecttype') !== 'undefined' ? tempObj.data('objecttype') : null;
                            //only attribute data-field broo
                            if (field == null) return;

                            //console.log('type : ' + type);
                            //console.log(tempObj.attr('data-objecttype') !== 'undefined' ? tempObj.data('objecttype') : null);

                            switch (type) {
                                case 'radio':
                                    console.log(obj.value);
                                    if (String(obj.value) == String(viewData[field])) {
                                        obj.checked = true;
                                    }

                                    break;
                                case 'select-one':

                                    /*Single Selection*/
                                    console.log(tempObj);
                                    var formattedValue = coreHelper.toDefault(viewData[field], '');
                                    var fieldText = typeof tempObj.attr('data-text') !== 'undefined' ? tempObj.data('text') : null;
                                    var formattedText = coreHelper.toDefault(viewData[fieldText], '');
                                    //tempObj.val({
                                    //    id: formattedValue,
                                    //    text: formattedText
                                    //}).trigger('change.select2');
                                    tempObj.select2("trigger", "select", {
                                        data: { id: formattedValue, text: formattedText }
                                    });
                                    break;
                                case 'hidden':
                                case 'text':
                                case 'number':
                                case 'email':
                                case 'textarea':
                                    var formattedValue = coreHelper.toDefault(viewData[field], '');
                                    if (tempObj.hasClass("date-picker")) {
                                        var format = tempObj.data("format") || "mm/dd/yyyy";
                                        //formattedValue = moment(formattedValue).format("yyyy/MM/dd");
                                        var resultFormat = moment(viewData[field]).format("MM/DD/YYYY");
                                        tempObj.datepicker('setDate', resultFormat);
                                        break;
                                    }
                                    else {
                                        switch (tempObj.data("inputmask")) {
                                            case "'alias': 'currency'":
                                                formattedValue = thousandSeparatorWithoutComma(formattedValue);
                                                break;
                                        }
                                    }
                                    tempObj.val(formattedValue);
                                    break;
                                case 'label':
                                case 'span':
                                    var formattedValue = coreHelper.toDefault(viewData[field], '');
                                    if (tempObj.hasClass("date-picker")) {
                                        var format = tempObj.data("format") || "MM/DD/YYYY";
                                        formattedValue = moment(formattedValue).format(format)
                                    }
                                    if (tempObj.hasClass("money")) {
                                        formattedValue = thousandSeparatorDecimal(formattedValue);
                                    } else if (tempObj.hasClass("money-decimal")) {
                                        formattedValue = thousandSeparatorFromValueWithComma(formattedValue);
                                    } else if (tempObj.hasClass("money-nondecimal")) {
                                        formattedValue = thousandSeparatorWithoutComma(formattedValue);
                                    }
                                    tempObj.text(formattedValue);
                                    break;
                                case 'select2':
                                    alert();
                                    break;
                                //default:
                                //    tempObj.html(formattedValue);
                                //    break;
                            }
                        }
                    } catch (e) {
                        
                    }
                });
            }
        }
    }

    $.helper = $.extend({}, routes, controlHelper);


    var template = {
        setLoading: function ($element, $to) {
            return $to.html($element);
        },
        setTemplate: function ($el) {
            return $el;
        }
    }
    var handler = {
        template : {
            fill: function (form, data, excludes) {
                var $html = form;
                var isFound = false;
                var start = 0;
                var end = 0;
                for (let i = 0; i < form.length; i++) {
                    if (form[i] == "$" && form[i + 1] == "." && form[i + 2] == "{") {
                        isFound = true;
                        start = i;
                    } 

                    if (isFound) {
                        if (form[i] == "}") {
                            end = i;
                            var variable = form.substring((start + 3), end);
                            var value = data[variable];
                            variable = "$.{" + variable + "}";
                            $html = $html.replace(variable, value);
                            isFound = false;
                            start = 0;
                            end = 0;
                        }  
                    }
                }
                return $html;
            }
        }
    }
    $.handler = $.extend({}, template, handler);
}(jQuery));