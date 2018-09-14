if (!String.prototype.startsWith) {
    String.prototype.startsWith = function (searchString, position) {
        position = position || 0;
        return this.substr(position, searchString.length) === searchString;
    };
}
$(function () {
    $.fn.swapLayoutItem = function ($area) {
        const SWAP_ITEM_MAP = {
            "selectable_sum_items": "layout_area_sum",
            "layout_area_col": "selectable_row_col_items",
            "layout_area_row": "selectable_row_col_items",
            "layout_area_sum": "selectable_sum_items"
        };
        var target = ($area || "" ) == "" ? SWAP_ITEM_MAP[$(this).parent().attr("id")] : $area;
        $("#" + target).append($(this).detach());
    };

    $.each(gLayoutSelectedCols.split(","), function (i, v) {
        $("#row_col_" + v).swapLayoutItem("layout_area_col");
    });

    $.each(gLayoutSelectedRows.split(","), function (i, v) {
        $("#row_col_" + v).swapLayoutItem("layout_area_row");
    });

    $.each(gLayoutSelectedSums.split(","), function (i, v) {
        $("#sum_" + v).swapLayoutItem();
    });
    updateListItems();

    var gDragItem = null;
    $(".layout_item").prop("draggable", "true");
    $(".layout_item").on("dragstart", function (e) {
        $(this).addClass("dragging");
        gDragItem = $(this);
    });

    $(".layout_item").on("dblclick", function () {
        if ($(this).parent().attr("id") == "selectable_row_col_items") {
            return;
        }
        $(this).swapLayoutItem();
        updateListItems();
    });

    $(".layout_item").on("dragend", function (e) {
        $(this).removeClass("dragging");
        gDragItem = null;
    });

    $(".drop-area").on("dragover", function (e) {
        if (gDragItem != null) {
            if (gDragItem.parent()[0] != $(this)[0] && gDragItem.attr('id').startsWith($(this).data("drop-prefix"))) {
                e.preventDefault();
                $(this).addClass("drag-over");
            }
        }
    });

    $(".drop-area").on("dragleave", function () {
        $(this).removeClass("drag-over");
    });

    $(".drop-area").on("drop", function () {
        $(this).removeClass("drag-over");
        if (gDragItem != null) {
            $(this).append(gDragItem.detach());
            updateListItems();
        }
    });

    $(".layout_item").on("dragover", function (e) {
        if (gDragItem != null) {
            if (gDragItem[0] != $(this)[0] && (gDragItem.parent()[0] == $(this).parent()[0] ||
                gDragItem.attr('id').startsWith($(this).parent().data("drop-prefix")))) {
                e.preventDefault();
                $(this).addClass("re-order-item");
                if (e.offsetX < $(this).width() / 2) {
                    $(this).removeClass("insert-right");
                    $(this).addClass("insert-left");
                } else {
                    $(this).removeClass("insert-left");
                    $(this).addClass("insert-right");
                }
            }
        }
    });

    $(".layout_item").on("drop", function (e) {
        $(this).removeClass("re-order-item");
        $(this).removeClass("insert-left");
        $(this).removeClass("insert-right");
        if (gDragItem != null) {
            if (e.offsetX < $(this).width() / 2) {
                $(this).before(gDragItem.detach());
            } else {
                $(this).after(gDragItem.detach());
            }
            updateListItems();
            gDragItem = null;
        }
    });

    $(".layout_item").on("dragleave", function () {
        $(this).removeClass("re-order-item");
        $(this).removeClass("insert-left");
        $(this).removeClass("insert-right");
    });

    function updateListItems() {
        var listCol = [];
        $("#layout_area_col .layout_item").each(function () {
            listCol.push($(this).data("item-id"));
        });

        var listRow = [];
        $("#layout_area_row .layout_item").each(function () {
            listRow.push($(this).data("item-id"));
        });

        var listSum = [];
        $("#layout_area_sum .layout_item").each(function () {
            listSum.push($(this).data("item-id"));
        });

        $("input[name=layout_cols]").val(listCol.concat());
        $("input[name=layout_rows]").val(listRow.concat());
        $("input[name=layout_sums]").val(listSum.concat());
    }
});

$('select.range-year').each(function () {
    var $this = $(this);
    var $from = $this.data("from");
    var $to = $this.data("to");
    var $val = $this.data("val");
    var $text = $this.data("text");
    if (typeof $text == "undefined") {
        $text = "";
    } else if (typeof $text != "string") {
        $text = $text.toString();
    }

    var $thisYear = (new Date()).getFullYear();

    if ($from == "now") {
        $from = $thisYear;
    } else if ((/^[+\-]\d+$/g).test($from)) {
        var move = Number($from.toString().slice(1));
        if ($from[0] == "+") {
            $from = $thisYear + move;
        } else {
            $from = $thisYear - move;
        }
    }

    if ($to == "now") {
        $to = $thisYear;
    } else if ((/^[+\-]\d+$/g).test($to)) {
        var move = Number($to.toString().slice(1));
        if ($to[0] == "+") {
            $to = $thisYear + move;
        } else {
            $to = $thisYear - move;
        }
    }

    if ($val == "now") {
        $val = $thisYear;
    } else if ((/^[+\-]\d+$/g).test($val)) {
        var move = Number($to.toString().slice(1));
        if ($val[0] == "+") {
            $val = $thisYear + move;
        } else {
            $val = $thisYear - move;
        }
    }

    for (var $y = $to; $y >= $from; $y--) {
        var newOption = $("<option>").val($y).text($text + $y.toString());
        if ($y == $val) newOption.prop("selected", true);
        $this.append(newOption);
    }
});

$('select.range-month').each(function () {
    var $this = $(this);
    var $val = $this.data("val");
    var $text = $this.data("text");
    if (typeof $text == "undefined") {
        $text = "";
    } else if (typeof $text != "string") {
        $text = $text.toString();
    }

    var $thisMonth = (new Date()).getMonth() + 1;
    if ($val == "now") {
        $val = $thisMonth;
    }
    for (var $m = 1; $m <= 12; $m++) {
        var newOption = $("<option>").val($m).text($text + $m.toString());
        if ($m == $val) newOption.prop("selected", true);
        $this.append(newOption);
    }
});

$('select.range-day').each(function () {
    var $this = $(this);
    var $val = $this.data("val");
    var $text = $this.data("text");
    if (typeof $text == "undefined") {
        $text = "";
    } else if (typeof $text != "string") {
        $text = $text.toString();
    }

    var $thisDay = (new Date()).getDate();
    if ($val == "now") {
        $val = $thisDay;
    }
    for (var $d = 1; $d <= 31; $d++) {
        var newOption = $("<option>").val($d).text($text + $d.toString());
        if ($d == $val) newOption.prop("selected", true);
        $this.append(newOption);
    }
});

$('select.range-hour').each(function () {
    var $this = $(this);
    var $val = $this.data("val");
    var $text = $this.data("text");
    if (typeof $text == "undefined") {
        $text = "";
    } else if (typeof $text != "string") {
        $text = $text.toString();
    }

    var $thisHours = (new Date()).getHours();
    if ($val == "now") {
        $val = $thisHours;
    }
    for (var $h = 0; $h <= 23; $h++) {
        var newOption = $("<option>").val($h).text($text + $h.toString());
        if ($h == $val) newOption.prop("selected", true);
        $this.append(newOption);
    }
});
