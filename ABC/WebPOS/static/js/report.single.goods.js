/**
 * Common for each report
 *      - Input order
 *      - Drill down
 *
 * @since       2017/09/04
 * @author      Nhut-BM
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

function initCommonSingleGoods() {
    initRedirectJan($("#query_result table.query-result"));
    initRedirectGood($("#query_result table.query-result"));
}

var $mainJanKey = ["ＪＡＮ", "JAN", "ISBN", "ＪＡＮコード", "JANコード"];
var $urlSingleGoodsAjax = "/report/single_goods_cumulative";
function redirectSingleGood($jan) {
    var $form = $("#form_search");
    // request key for search single Item
    $("input[name=jan_code]").val($jan.trim());

    $form.attr("method","POST");
    $form.attr("action",$urlSingleGoodsAjax);
    $form.attr("target","_blank");
    $form.submit();
}

var janSelected = "";
function initRedirectJan($tableRender) {
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
            redirectSingleGood($jancd);
            janSelected = "";
        });
    });
}

var $mainGoodKey = ["品名", "品名"];
var $urlHonto = "http://honto.jp/";

function initRedirectGood($tableRender) {
    // Check require
    if ($tableRender.length == 0) {
        return;
    }

    var $thGood = null;
    $.each($mainGoodKey, function ($key, $val) {
        if ($thGood == null || $thGood.length == 0) {
            return $thGood = $tableRender.find("th:contains('" + $val + "')");
        }
    });
    if ($thGood == null || $thGood.length == 0) {
        return;
    }

    // Init field for show-hide
    appendCssFile("css/common.css");

    var $colIndexGood = $thGood.index();
    var $rowTarget = "tbody tr";
    $tableRender.find($rowTarget).each(function () {
        var $fieldGood = $(this).find("td:eq(" + $colIndexGood + ")");
        addClass($fieldGood, "link_data");
        $fieldGood.on("click", function () {
            var $jancd = $(this).prev().html().trim();
            var strClass = ""
            if ($jancd.toString().match("^9784")) {
                strClass = 'isbn/';
                window.open($urlHonto + strClass + $jancd);
            } else {
                strClass = 'jan/';
                window.open($urlHonto + strClass + $jancd);
            }
        });
    });
}
