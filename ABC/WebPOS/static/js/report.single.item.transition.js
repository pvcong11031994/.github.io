/**
 * Common for each report
 *      - Input order
 *      - Drill down
 *
 * @since       2016/12/13
 * @author      Huan-TN
 */

/**
 * Check and load file if not exist
 * @param $filePath location of file
 */
function appendCssFile($filePath) {
    if (!$("link[href='/static/" + $filePath + "']").length) {
        $('head').append('<link rel="stylesheet" type="text/css" href="/static/' + $filePath + '"/>');
    }
}

/**
 * Like toggleClass but just addClass
 * @param $elm      Element need add class
 * @param strClass  class name
 */
function addClass($elm, strClass) {
    if (!$elm.hasClass("current_report")) {
        $elm.addClass(strClass);
    }
}

/**
 * Append data selected/input to form submit or object data
 * @param $name     name of element get data
 * @param $form     form submit
 * @param $data     data submit
 */
function addFromSelect($name, $form, $data) {
    if ($("select[name*=" + $name + "]").length > 0) {
        var valueSelect = $("select[name*=" + $name + "]").val();
        if (valueSelect != null && Array.isArray(valueSelect)) {
            $.each(valueSelect, function (i, v) {
                addDataToFormOrDataObj($name, $form, $data, v);
            });
        } else if (valueSelect != null) {
            addDataToFormOrDataObj($name, $form, $data, valueSelect);
        }
    } else if ($("input[name*=" + $name + "]").length > 0) {
        addDataToFormOrDataObj($name, $form, $data, $("input[name*=" + $name + "]").val());
    }
}
function addInputFromToLikeName($name, $form, $data) {
    if ($("input[name*=" + $name + "_from]").length > 0) {
        addDataToFormOrDataObj($name + "_from", $form, $data, $("input[name*=" + $name + "_from]").val());
        addDataToFormOrDataObj($name + "_to", $form, $data, $("input[name*=" + $name + "_to]").val());
    }
}
function addDataToFormOrDataObj($name, $form, $data, $value) {
    if ($form != undefined && $form != null) {
        $form.append($("<input name='" + $name + "'/>").val($value));
    } else if ($data != undefined && $data != null) {
        $data[$name] = $value;
    }
}

function initCommonSingleItemTransition() {
    initInputOrder($("#query_result table.query-result"));
}

function initAreaCommon() {
    if ($("div.area_single_item_view").length == 0) {
        $("<div class='area_single_item_view'></div>").insertAfter($("table.table_search"));
        $("div.area_single_item_view").hide();
    }
}

function showAreaCommon($jan) {

    // display and init single item screen
    $("span.main-report").css("display","none");
    $("span.single-item").css("display","inherit");
    $("table tr.single-item").css("display","table-row");
    $("table td.single-item").css("display","table-row");
    $(".jan-item").prop("disabled",false);
    $("div.single-item").css("display","inherit");
    if ($jan != ""){
        $("input[data-form-name=flag_single_item]").val($jan);
        $("input[data-form-name=jan_cd]").val($jan);
    }
    initAreaCommon();
    if (!$("div.area_single_item_view").is(':visible')) {
        $(".current_report").hide();
        $("div.area_single_item_view").css("display","grid");
        $("div.area_single_item_view").html("");
    }
}
function hideAreaCommon() {

    // hide single item screen
    $("span.main-report").css("display","inherit");
    $("span.single-item").css("display","none");
    $("table tr.single-item").css("display","none");
    $("table td.single-item").css("display","none");
    $("div.single-item").css("display","none");

    // reset value
    $("input[name=jan_cd]").val("");
    $(".jan-item").prop("disabled",true);
    $("input[data-form-name=flag_single_item]").val("");

    // show main search screen
    $(".current_report").show();
    $("div.area_single_item_view").hide();
    $("div.area_single_item_view").html("");

    // fix column result table
    FixedMidashi.create();


    initCommonSingleItemTransition();
}

// Close view
$(document).off("click", ".close_view", function () {
});
$(document).on("click", ".close_view", function () {
    hideAreaCommon();
});
$(document).off("click", ".close_all_view", function () {
});
$(document).on("click", ".close_all_view", function () {
    hideAreaCommon();
});

var $mainJanKey = ["ＪＡＮ", "JAN", "ISBN", "ＪＡＮコード", "JANコード"];
var $urlMultiShopConfirm = "/report/single_item_transition/single_item_transition_find_ajax";
function confirmSingleItemTransition($jan) {
    var $form = $("<form method='POST' action='" + $urlMultiShopConfirm + "'></form>");
    $form.append($("input[name='csrf.Token']:last").clone());
    // request key for search single Item
    $form.append($("<input name='jan'/>").val($jan));
    var sqlData = $("input[name='rand_string_select']").val();
    var flagSuper = $("input[name='super_flag_select']").val();

    $form.append($("<input name='rand_string'/>").val(sqlData));
    $form.append($("<input name='super_flag'/>").val(flagSuper));

    $.ajax({
        type: $form.attr('method'),
        url: $form.attr('action'),
        cache: false,
        data: $form.serialize(),
        success: function (data) {
            if (data.success == false ||
                data.success == "false" ) {
                showError(data.msg_err);
            } else {
                showAreaCommon($jan);
                $("div.area_single_item_view").html(data);
                FixedMidashi.create();
            }
        }
    });
}

var janSelected = "";
function initInputOrder($tableRender) {
    // Check require
    if ($tableRender.length == 0) {
        return;
    }

    var $thJAN = null;
    $.each($mainJanKey, function ($key, $val) {
        if ($thJAN == null || $thJAN.length == 0) {
            return $thJAN = $tableRender.find("th:contains('" + $val + "')");
        }
    });
    if ($thJAN == null || $thJAN.length == 0) {
        return;
    }

    // Init field for show-hide
    appendCssFile("css/common.css");

    var $colIndexJAN = $thJAN.index();
    var $rowTarget = "tbody tr";
    janSelected = "";
    $tableRender.find($rowTarget).each(function () {
        var $fieldJAN = $(this).find("td:eq(" + $colIndexJAN + ")");
        addClass($fieldJAN, "click_for_single_item_transition");
        $fieldJAN.on("click", function () {
            var $jancd = $(this).html();
            if (janSelected == "") {
                janSelected = $jancd;
            } else if (janSelected == $jancd){
                return;
            }
            $(this).addClass("jan_has_click");
            confirmSingleItemTransition($(this).html());
            janSelected = "";
        });
    });
}


// search screen
var $validateBefore = null;

$(function () {
    if ($("#dialog_loading_single_item").length == 0) {
        $("<div>").attr("id", "dialog_loading_single_item")
            .append($("<span>").addClass("processing-text").text("集計実行中 ..."))
            .append($("<p>").attr("id", "processing_time_single_item"))
            .hide().appendTo($("body"));
    }

    var timeCounter = function () {
        var obj = new Object();

        obj.start = function () {
            obj.time_begin = (new Date()).getTime();

            obj.time_interval = setInterval(function () {
                obj.time_now = (new Date()).getTime();
                obj.time_elapse = obj.time_now - obj.time_begin;

                $('#processing_time_single_item').text((obj.time_elapse / 1000).toFixed(1) + ' 秒');

            }, 50);
        };

        obj.stop = function () {
            clearInterval(obj.time_interval);
        };

        obj.setBeginHtml = function () {
            obj.time_begin_html = (new Date()).getTime();
        };

        obj.setNow = function () {
            obj.time_now = (new Date()).getTime();
        };

        obj.htmlTime = function () {
            return ((obj.time_now - obj.time_begin_html) / 1000).toFixed(0);
        };

        return obj;
    }();

    var loadingDialogSingleItem = function () {
        var obj = new Object();

        obj.dialog = $('#dialog_loading_single_item').dialog({
            width: 280,
            height: 180,
            modal: true,
            resizable: false,
            closeOnEscape: false,
            autoOpen: false,
            buttons: {
                'キャンセル': function () {
                    queryHandler.cancel();
                    loadingDialogSingleItem.close();
                    timeCounter.stop();
                }
            },
            open: function () {
                $("#dialog_loading_single_item").parent().find(".ui-dialog-titlebar-close").remove();
                timeCounter.start();
                queryHandler.submit();
            }
        });

        obj.show = function () {
            obj.dialog.dialog('open');
        };

        obj.close = function () {
            obj.dialog.dialog('close');
        };

        return obj;
    }();

    var queryHandler = function () {
        var obj = new Object();
        obj.submit = function () {
            obj.queryAjax = $.ajax({
                url: location.pathname + "/query_ajax",
                type: "POST",
                data: $("#form_search").serialize(),
                context: {nosubpage: true},
                success: function (responeHtml) {
                    timeCounter.stop();
                    timeCounter.setBeginHtml();
                    $('#processing_time_single_item').html('データを受信しました。<br/>ブラウザでレンダリング...');
                    $("#form_search").find("input[name=search_handle_type]").val("");
                    setTimeout(function () {
                        showAreaCommon("");
                        $("div.area_single_item_view").html(responeHtml);
                        $("table td.single-item").css("display","inherit");
                        FixedMidashi.create();
                        loadingDialogSingleItem.close();
                        timeCounter.setNow();
                        console.log('BROWSER: Render html time: ' + timeCounter.htmlTime() + ' seconds');
                        // if (typeof initCommon === "function") {
                        //     //What common?
                        //     initCommon();
                        // }
                    }, 100);
                },
                error: function (event) {
                    if (event.statusText == "abort") return;
                    var message;
                    try {
                        var response = JSON.parse(event.responseText);
                        message = "エラー " + response.code + " (" + response.status + ")<br/>" + response.text;
                    } catch (e) {
                        message = event.responseText;
                    }
                    $('#query_result').html("<p class='query-error'>" + message + "</p>").show();
                    timeCounter.stop();
                    loadingDialogSingleItem.close();
                }
            });
        };

        obj.cancel = function () {
            if (obj.queryAjax) {
                obj.queryAjax.abort();
            }
        };

        return obj;
    }();


    $("#btn_search_single_item").click(function (event) {
        if ($validateBefore != null && $validateBefore != undefined && $.isFunction($validateBefore) && !$validateBefore()) {
            event.preventDefault();
            return;
        }

        /* MyStorage is used for "Input order"? */
        if (typeof MyStorage !== 'undefined') {
            if (typeof MyStorage.clear === "function") {
                MyStorage.clear();
            }
        }

        $('div.area_single_item_view').html("").hide();
        $('div.graph_area').html("");
        $("table td.single-item").css("display","none");
        loadingDialogSingleItem.show();
    });
});
