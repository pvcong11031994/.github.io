if (!String.prototype.startsWith) {
    String.prototype.startsWith = function (searchString, position) {
        position = position || 0;
        return this.substr(position, searchString.length) === searchString;
    };
}
$(function () {
    $.fn.swapLayoutItem = function () {
        const SWAP_ITEM_MAP = {
            "selectable_col_items": "layout_area_col",
            "selectable_row_items": "layout_area_row",
            "selectable_sum_items": "layout_area_sum",
            "layout_area_col": "selectable_col_items",
            "layout_area_row": "selectable_row_items",
            "layout_area_sum": "selectable_sum_items"
        };
        var target = SWAP_ITEM_MAP[$(this).parent().attr("id")];
        $("#" + target).append($(this).detach());
    };

    $.each(gLayoutSelectedCols.split(","), function (i, v) {
        $("#col_" + v).swapLayoutItem();
    });

    $.each(gLayoutSelectedRows.split(","), function (i, v) {
        $("#row_" + v).swapLayoutItem();
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

