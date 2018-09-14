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

var $sumHeaderRow = false;
var $inputOrder = false;
function initCommon() {
    putShopToArray("shop_cd");
    initAreaCommon();

    initDrillDown("#query_result table.query-result");

    if ($inputOrder) {
        initInputOrder($("#query_result table.query-result"));
    }
}

function initAreaCommon() {
    if ($("div.area_common_view").length == 0) {
        $("<div class='area_common_view'></div>")
            .insertAfter($("span.page-header"));
        $("div.area_common_view").hide();
    }
}

function showAreaCommon() {
    initAreaCommon();
    if (!$("div.area_common_view").is(':visible')) {
        $(".current_report").hide();
        $("div.area_common_view").show();
        $("div.area_common_view").html("");
    }
}
function hideAreaCommon() {
    $(".current_report").show();
    $("div.area_common_view").hide();
    $("div.area_common_view").html("");
}

// Close view
$(document).off("click", ".close_view", function () {
});
$(document).on("click", ".close_view", function () {
    $(this).parents("div.child_area_view").remove();
    if ($("div.child_area_view").length == 0) {
        hideAreaCommon();
    } else {
        $("div.child_area_view").last().show();
    }
});
$(document).off("click", ".close_all_view", function () {
});
$(document).on("click", ".close_all_view", function () {
    hideAreaCommon();
});

var vShop = [];
var $shop = null;
function putShopToArray($itemName) {
    $shop = null;
    vShop = [];
    if ($("select[name*=" + $itemName + "]").length > 0) {
        $shop = $("select[name*=" + $itemName + "]");
        vShop = $shop.val();
    } else if ($("input[name*=" + $itemName + "]").length > 0) {
        $shop = $("input[name*=" + $itemName + "]");
        vShop.push($shop.val());
    }
}

/**
 * Code common Wiki https://vjppd.backlog.jp/alias/wiki/1074337016
 * Import file in report require
 * Function initInputOrder() call each time button 集計実行 click :
 *      - Check and render html if exist field JAN/ISBN in result report
 *      - Process for 1 shop or multi shop
 *
 * @since       2016/11/20
 * @author      Huan-TN
 */

var $objInputTotal = {};
var MyStorage = {
    clear: function () {
        $objInputTotal = {};
    }
};
$(document).off("click", ".set_all_input_order", function () {
});
$(document).on("click", ".set_all_input_order", function () {
    if ($.trim($("input[name=all_input_order]").val()) != "") {
        $(this).parents("table").last().find("input[name=item_input_order]:enabled").val($("input[name=all_input_order]").val()).trigger("change");
    }
});
$(document).off("click", ".clear_all_input_order", function () {
});
$(document).on("click", ".clear_all_input_order", function () {
    $("input[name=all_input_order]").val(null).trigger("change");
    $("input[name=item_input_order]:enabled").val(null).trigger("change");
});
$(document).off("keydown", "input[name=item_input_order], input[name=all_input_order]", function () {
});
$(document).on("keydown", "input[name=item_input_order], input[name=all_input_order]", function (e) {
    var key = e.which || e.keyCode || 0;
    if (!(48 <= key && key <= 57) && !(96 <= key && key <= 105)
        && key != 0 && key != 8
        && key != 9 && key != 13
        && key != 46 && key != 37 && key != 39
    ) {
        return false;
    }
});
$(document).off("change", "input[name=item_input_order]", function () {
});
$(document).on("change", "input[name=item_input_order]", function () {
    $(this).removeClass("input_has_data");
    if ($.trim($(this).val()) != "") {
        addClass($(this), "input_has_data");
    }
    if (vShop.length == 1) {
        $objInputTotal = $objInputTotal == null ? {} : $objInputTotal;
        var keyItem = $(this).parents("td").find("input[name=jan]").val() + vShop[0];
        if ($(this).val() != "") {
            $objInputTotal[keyItem] = {
                value: $(this).val(),
                disable: false
            }
        } else if ($objInputTotal[keyItem] != null && $objInputTotal[keyItem] != undefined) {
            delete $objInputTotal[keyItem];
        }
    }
});
var $mainJanKey = ["ＪＡＮ", "JAN", "ISBN", "ＪＡＮコード", "JANコード"];

// Single shop
var $urlSingleShopConfirm = "/report/input_order/single_shop_confirm_ajax";
function confirmSingleShop() {
    if ($("input.input_has_data").length == 0) {
        return false;
    }

    var $form = $("<form method='POST' action='" + $urlSingleShopConfirm + "'></form>");
    $form.append($("input[name='csrf.Token']:last").clone());
    $form.append($("<input name='shop_cd'/>").val(vShop[0]));
    $form.append($("input.input_has_data").parents("td").clone());

    $objInputTotal = $objInputTotal == null ? {} : $objInputTotal;
    $form.find("input.input_has_data").parents("td").each(function () {
        var jan = $(this).find("input[name=jan]").val();
        $objInputTotal[jan + vShop[0]] = {
            value: $(this).find("input[name=item_input_order]").val(),
            disable: false
        }
    });
    $form.find("input[name=all_input_order]").remove();
    $form.find("input[name=item_input_order]").remove();

    $.ajax({
        type: $form.attr('method'),
        url: $form.attr('action'),
        data: $form.serialize(),
        success: function (data) {
            if (data.success && !data.success) {
                subpage({
                    reInit: true,
                    title: "",
                    message: data.msg_err
                });
            } else {
                showAreaCommon();
                $("div.child_area_view").hide();
                $("div.area_common_view").append(data);
            }
        }
    });
}

var $urlMultiShopConfirm = "/report/input_order/multi_shop_find_ajax";
function confirmMultiShop($jan) {
    var $form = $("<form method='POST' action='" + $urlMultiShopConfirm + "'></form>");
    $form.append($("input[name='csrf.Token']:last").clone());
    $form.append($("<input name='jan'/>").val($jan));
    $.each(vShop, function (i, v) {
        $form.append($("<input name='shop_cd'/>").val(v));
    });
    addInputFromToLikeName("date", $form);
    addInputFromToLikeName("year", $form);
    addInputFromToLikeName("month", $form);
    addInputFromToLikeName("day", $form);
    addInputFromToLikeName("time", $form);
    addFromSelect("day_of_week", $form);
    addFromSelect("year", $form);
    addFromSelect("month", $form);
    $.ajax({
        type: $form.attr('method'),
        url: $form.attr('action'),
        data: $form.serialize(),
        success: function (data) {
            if (data.success && !data.success) {
                subpage({
                    reInit: true,
                    title: "",
                    message: data.msg_err
                });
            } else {
                showAreaCommon();
                $("div.child_area_view").hide();
                $("div.area_common_view").append(data);
            }
        }
    });
}

/**
 * Code common drill down : Wiki link https://vjppd.backlog.jp/alias/wiki/1074341315
 *
 * @since       2016/12/01
 * @author      Huan-TN
 */
var listTag = {
    K_LIST_PRODUCT: "list_product",
    K_SHOP_CD: "drill_shop_cd",
    K_SHOP_NAME: "drill_shop_name",
    K_KUBUN: "drill_kubun",
    K_BUMON: "drill_bumon",

    K_B_GENRE: "drill_bGenre",
    K_M_GENRE: "drill_mGenre",
    K_SM_GENRE: "drill_smGenre",
    K_S_GENRE: "drill_sGenre",

    K_YEAR: "drill_year",
    K_MONTH: "drill_month",

    K_GOODS_COUNT: "drill_goodsCount",
    K_AMOUNT: "drill_amount"
};

var listElem = {
    K_SHOP_CD: ["店舗コード", "店舗CD"],
    K_SHOP_NAME: ["店舗名"],
    K_KUBUN: ["店舗POS分類", "店舗POS分類CD", "POSジャンル"],
    K_BUMON: ["BSPOS分類", "BSPOS分類CD", "店舗POSジャンル"],

    K_B_GENRE: ["メディア大分類", "MS大分類", "メディア大分類コード"],
    K_M_GENRE: ["メディア中分類", "MS中分類", "メディア中分類コード"],
    K_SM_GENRE: ["メディア中小分類", "MS中小分類", "メディア中小分類コード"],
    K_S_GENRE: ["メディア小分類", "MS小分類", "メディア小分類コード"],

    K_YEAR: ["年"],
    K_MONTH: ["月"],

    K_GOODS_COUNT: ["数量"],
    K_AMOUNT: ["金額"]
};
var eRows = {};
var eCols = {};
var eSums = {};
var tableTHEAD = false;
var $sumHeaderText = ["総合計"];
var $drill_media = "";
var $drill_name = "";

/**
 * Init all item has event drill down
 */
function initDrillDown($tableRender) {
    $drill_media = $($tableRender).data("drill-item");
    $drill_name = $($tableRender).data("drill-name");
    tableTHEAD = $($tableRender).find("thead").length == 1;
    eRows = {};
    eCols = {};
    $.each(listElem, function ($key, $val) {
        var $ele = null;
        $.each($val, function ($k, $v) {
            if ($ele == null || $ele.length == 0) {
                $ele = $($tableRender + " th:contains('" + $v + "')");
                if (($key == "K_YEAR" || $key == "K_MONTH") && $ele.html() != $v) {
                    $ele = null;
                }
                if ($drill_media != "" && ($key == "K_GOODS_COUNT" || $key == "K_AMOUNT")) {
                    $ele = null;
                }
                return $ele
            }
        });
        if ($ele != null && $ele != undefined && $ele.length > 0) {
            if ($ele.hasClass("row-name")) {
                eRows[$key] = $ele;
            } else if ($ele.hasClass("col-name")) {
                eCols[$key] = $ele;
            } else if ($ele.hasClass("sum-name")) {
                eSums[$key] = $ele;
            }
        }
    });

    if ($.isEmptyObject(eRows) && $.isEmptyObject(eCols)) {
        return;
    }

    putShopToArray("shop_cd");
    appendCssFile("css/report/handle.common.css");


    $.each($sumHeaderText, function ($key, $val) {
        if (!$sumHeaderRow) {
            return $sumHeaderRow = $($tableRender + " th:contains('" + $val + "'):eq(0)").length > 0;
        }
    });

    var $rowIndex = 0;
    var $eleCheck = $($tableRender + " th.row-name:eq(0)");
    var $strItem = $tableRender + " tbody tr";
    if (!tableTHEAD) {
        var $thRowspan = parseInt($eleCheck.attr('rowspan') == undefined ? 1 : $eleCheck.attr('rowspan')) - 1;
        $rowIndex = $eleCheck.parents("tr").index() + $thRowspan;
        if ($sumHeaderRow) {
            $rowIndex += 1;
        }
        $strItem += ":gt(" + $rowIndex + ")";
    }

    if (!$.isEmptyObject(eRows)) {
        $.each(eRows, function ($key, $val) {
            $($strItem).each(function () {
                var $ele = $(this).find("td:eq(" + $val.index() + ")");
                if ($.trim($ele.html()) == "") {
                    return;
                }
                $ele.addClass("allow_drill_down");
                $ele.addClass(listTag[$key]);
                $ele.data("info", listTag[$key]);
            });
        });
    }

    if (!$.isEmptyObject(eSums)) {
        $.each(eSums, function ($key, $sum_header) {
            $.each($sum_header, function ($i, $val) {
                var $index = $($val).index();
                $($strItem).each(function () {
                    var $ele = $(this).find("td:eq(" + $index + ")");
                    if ($.trim($ele.html()) == "") {
                        return;
                    }
                    $ele.addClass("allow_drill_down");
                    $ele.addClass(listTag[$key]);
                    $ele.data("info", listTag[$key]);
                });
            });
        });
    }

    if (!$.isEmptyObject(eCols)) {
        $strItem = $tableRender + " tbody th.col-header-";
        if (tableTHEAD) {
            $strItem = $tableRender + " thead th.col-header-";
        }
        $.each(eCols, function ($key, $val) {
            $($strItem + $val.parents("tr").index()).each(function () {
                if ($.trim($(this).html()) == "") {
                    return;
                }
                $(this).addClass("allow_drill_down");
                $(this).addClass(listTag[$key]);
                $(this).data("info", listTag[$key]);
            });
        });
    }
}

// Set handle event
$(document).off("click", ".allow_drill_down", function () {
});
$(document).on("click", ".allow_drill_down", function (event) {
    var $requestData = {};
    var $name = "";
    var $ele = "";
    var $valueItem = $.trim($(this).html());
    switch ($(this).data("info")) {
        case listTag["K_SHOP_NAME"]:
        case listTag["K_SHOP_CD"]:
            $name = "shop_cd";
            $ele = "select[name=shop_cd] option:contains('" + $valueItem + "')";
            $requestData["drill_name"] = listTag["K_LIST_PRODUCT"];
            break;
        case listTag["K_KUBUN"]:
            $name = "kubun_cd";
            $ele = "select[name=kubun_cd] option:contains('" + $valueItem + "')";
            $requestData["drill_name"] = listTag["K_LIST_PRODUCT"];
            break;
        case listTag["K_BUMON"]:
            $name = "bumon_cd";
            $ele = "select[name=bumon_cd] option:contains('" + $valueItem + "')";
            $requestData["drill_name"] = listTag["K_LIST_PRODUCT"];
            break;
        case listTag["K_B_GENRE"]:
            $name = "media_group1_cd";
            $requestData["drill_name"] = $name;
            $ele = "select[name=media_group1_cd] option:contains('" + $valueItem + "')";
            break;
        case listTag["K_M_GENRE"]:
            $name = "media_group2_cd";
            $requestData["drill_name"] = $name;
            $ele = "select[name=media_group2_cd] option:contains('" + $valueItem + "')";
            break;
        case listTag["K_SM_GENRE"]:
            $name = "media_group3_cd";
            $requestData["drill_name"] = $name;
            $ele = "select[name=media_group3_cd] option:contains('" + $valueItem + "')";
            break;
        case listTag["K_S_GENRE"]:
            $name = "media_group4_cd";
            $requestData["drill_name"] = $name;
            $ele = "select[name=media_group4_cd] option:contains('" + $valueItem + "')";
            $requestData["drill_name"] = listTag["K_LIST_PRODUCT"];
            break;
        case listTag["K_YEAR"]:
            $name = "year";
            $requestData["drill_name"] = $name;
            $requestData["drill_item"] = $name;
            $requestData[$name] = $valueItem;
            break;
        case listTag["K_MONTH"]:
            $name = "month";
            $requestData["drill_name"] = $name;
            $requestData["drill_item"] = $name;
            $requestData[$name] = $valueItem;
            break;
        case listTag["K_GOODS_COUNT"]:
        case listTag["K_AMOUNT"]:
            $ele = "";
            $name = $drill_media;
            $valueItem = $(this).parent("tr").find("td:eq(0)").html();
            $requestData["drill_name"] = listTag["K_LIST_PRODUCT"];
            break;
        default:
            $name = "";
            $ele = "";
    }

    if ($($ele).val() != undefined && $($ele).val() != null) {
        $requestData["drill_item"] = $name;
        $requestData[$name] = $($ele).val();
    } else if ($valueItem != "") {
        $requestData["drill_item"] = $name;
        $requestData[$name] = $valueItem.split(" ")[0];
    } else {
        $requestData = {};
    }

    if (!$.isEmptyObject($requestData)) {
        handleDrillDown($requestData);
    }
});

var $urlDrillDownAction = "/report/drill_down/drill_ajax";
function handleDrillDown($data) {
    var $form = $("<form method='POST' action='" + $urlDrillDownAction + "'></form>");
    $form.append($("input[name='csrf.Token']:last").clone());
    var $drillItem = $data["drill_item"];
    if ($drillItem != "shop_cd") {
        $.each(vShop, function (i, v) {
            $form.append($("<input name='shop_cd'/>").val(v));
        });
        $.each($data, function (key, value) {
            $form.append($("<input name='" + key + "'/>").val(value));
        });
    } else {
        $form.append($("<input name='drill_name'/>").val(listTag["K_LIST_PRODUCT"]));
        $form.append($("<input name='drill_item'/>").val($drillItem));
        $form.append($("<input name='shop_cd'/>").val($data["shop_cd"]));
    }

    addInputFromToLikeName("date", $form);
    addInputFromToLikeName("year", $form);
    addInputFromToLikeName("month", $form);
    addInputFromToLikeName("day", $form);
    addInputFromToLikeName("time", $form);
    addFromSelect("day_of_week", $form);
    addFromSelect("year", $form);
    addFromSelect("month", $form);

    $.ajax({
        type: $form.attr('method'),
        url: $form.attr('action'),
        data: $form.serialize(),
        success: function (data) {
            showAreaCommon();
            $("div.child_area_view").hide();
            $("div.area_common_view").append(data);
            var $groupBtn = $("div.child_area_view").last().find(".foot_btn_group");
            $groupBtn.attr("colspan", $("div.child_area_view").last().find("table thead").find(".row-name").parent("tr").find("th").length);

            initDrillDown("div.child_area_view:last table.query-result");
            initInputOrder($("div.child_area_view").last().find("table.query-result"));
        }
    });
}

function initInputOrder($tableRender) {
    // Check require
    if ($tableRender.length == 0 || vShop.length == 0) {
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

    initAreaCommon();

    var $thRowspan = parseInt($thJAN.attr('rowspan') == undefined ? 1 : $thJAN.attr('rowspan')) - 1;
    var $rowIndexJAN = $thJAN.parents("tr").index() + $thRowspan;
    var $colIndexJAN = $thJAN.index();
    var $rowSpanHeader = $rowIndexJAN + 1;
    var $startIndexRow = $rowIndexJAN;
    if ($sumHeaderRow > 0) {
        $rowSpanHeader += 1;
        $startIndexRow += 1;
    }

    var $tbHasTHead = $tableRender.find("thead").length == 1;
    // Check length shop search
    if (vShop.length == 1) {
        // Append field input order
        // Header input
        var $thHeaderInputOrder = $("<th style='padding: 3px 4px;' class='header-input-order' rowspan='" + $rowSpanHeader + "'></th>");
        $thHeaderInputOrder.append($("<label>発注数</label>"));
        $thHeaderInputOrder.append($("<label class='clear_all_input_order' style='cursor: pointer;position: absolute;left: 60px;bottom: 26px;background-color: #C0C0C0;padding: 2px 5px;' >クリア</label>"));
        $thHeaderInputOrder.append($("<br/>"));
        $thHeaderInputOrder.append($("<input class='all_input_order' maxlength='4' type='text' name='all_input_order'/>"));
        $thHeaderInputOrder.append($("<label class='set_all_input_order' style='cursor: pointer;margin-left: 5px;background-color: #C0C0C0;padding: 2px 5px;'>全て左の数にする</label>"));
        if ($tbHasTHead) {
            $tableRender.find("thead tr:eq(0)").append($thHeaderInputOrder);
        } else {
            $tableRender.find("tbody tr:eq(0)").append($thHeaderInputOrder);
        }

        // Child row input
        var $itemInput = $("<input class='item_input_order' maxlength='4' type='text' name='item_input_order'/>");
        //
        var $zoneAppend = $tbHasTHead ? "tbody tr" : "tbody tr:gt(" + $startIndexRow + ")";
        $tableRender.find($zoneAppend).each(function () {
            var jan = $.trim($(this).find("td:eq(" + $colIndexJAN + ")").html());
            var $curItem = $itemInput.clone(true, true);
            if ($objInputTotal[jan + vShop[0]] != undefined && $objInputTotal[jan + vShop[0]] != null) {
                var $objSet = $objInputTotal[jan + vShop[0]];
                $curItem.val($objSet.value);
                $curItem.prop("disabled", $objSet.disable);
                if (!$objSet.disable) {
                    $curItem.addClass('input_has_data');
                }
            }
            var $tdItem = $("<td class='jan_" + jan + "'></td>");
            $tdItem.append($curItem);
            $tdItem.append($("<input type='hidden' name='jan'/>").val(jan));
            $(this).append($tdItem);
        });

        var $button = $("<button type='submit' class='btn btn-success btn-me right-20' '>登録</button>");
        $button.on("click", function () {
            confirmSingleShop();
        });
        var colCount = 1;
        if ($tableRender.find("tfoot").length == 0) {
            //
            var $zoneCount = $tbHasTHead ? "thead tr:eq(0)" : "tbody tr:eq(0)";
            $tableRender.find($zoneCount + " th").each(function () {
                if ($(this).attr('colspan')) {
                    colCount += parseInt($(this).attr('colspan'));
                } else {
                    colCount++;
                }
            });
            $tableRender.find($zoneCount + " td").each(function () {
                if ($(this).attr('colspan')) {
                    colCount += parseInt($(this).attr('colspan'));
                } else {
                    colCount++;
                }
            });
            var $td = $("<td colspan='" + colCount + "' style='text-align: center;'>");
            $td.append($button);
            $tableRender.append($("<tfoot></tfoot>").append($("<tr></tr>").append($td)));

        } else {
            colCount += parseInt($tableRender.find("tfoot tr td").attr("colspan"));
            $tableRender.find("tfoot tr td").attr("colspan", colCount)
                .append($button);
        }

    } else {
        var $rowTarget = "tbody tr:gt(" + $startIndexRow + ")";
        if ($tbHasTHead) {
            $rowTarget = "tbody tr";
        }
        $tableRender.find($rowTarget).each(function () {
            var $fieldJAN = $(this).find("td:eq(" + $colIndexJAN + ")");
            addClass($fieldJAN, "click_for_order");
            $fieldJAN.on("click", function () {
                $(this).addClass("jan_has_click");
                confirmMultiShop($(this).html());
            });
        });
    }
}
