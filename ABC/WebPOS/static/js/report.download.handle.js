/**
 * Created  2017/02/08.
 */
/* 20170208 Common Download File */
$("#btn_download_csv").click(function (event) {
    // Download file CSV with fortmat 集計結果＋店舗 in best_sales when selected shop_code over 100 shops
    if ((location.pathname === "/report/best_sales" || location.pathname === "/report/best_sales_maria" || location.pathname === "/report/best_sales_cloud")
            && $(".download_type").val() == 2 && $("#shop_cd").val() != null && $("#shop_cd").val().length > 100) {
            // show dialog
            if ($("#dialog_download_overload_notify").length == 0) {
                $("<div>").attr("id", "dialog_download_overload_notify")
                    .append($("<span>").text("ダウンロード対象の店舗数が100を超えているため"))
                    .append($("<br>"))
                    .append($("<span>").text("店舗プルダウンでチェックされている、上から100店舗分をダウンロードします。"))
                    .hide().appendTo($("body"));
            }
            var loadingDialogDownloadNotify = function () {
                var obj = new Object();
                obj.dialog = $('#dialog_download_overload_notify').dialog({
                    width: 500,
                    height: 150,
                    modal: true,
                    resizable: true,
                    closeOnEscape: false,
                    autoOpen: false,
                    data: $("#form_search").serialize(),
                    buttons: [
                        {
                            text: "OK",
                            "style":'margin-right:100px; width: 80px; height: 28px',
                            click: function() {
                                loadingDialogDownloadNotify.close();
                                $("#btn_search").trigger("click");
                            }
                        },
                        {
                            text: "キャンセル",
                            click: function() {
                                loadingDialogDownloadNotify.close();
                            }
                        }
                    ],
                    open: function () {
                        $("#dialog_download_overload_notify").parent().find(".ui-dialog-titlebar-close").remove();
                    }
                })

                obj.show = function () {
                    obj.dialog.dialog('open');
                };

                obj.close = function () {
                    obj.dialog.dialog('close');
                };

                return obj;
            }();
            $("input[name=search_handle_type]").val("1");
            loadingDialogDownloadNotify.show();
        } else {
        // download CSV
        $("input[name=search_handle_type]").val("1");
        $("#btn_search").trigger("click");
    }

});
$(document).on("click", ".btn-download-file", function (event) {
    $(".report-form-download-file").submit();
});

/**
 * Add new for single_item_transition screen 2017/06/27.
 */
/* 20170208 Common Download File */
$("#btn_download_csv_single_item").click(function (event) {
    $("input[name=search_handle_type]").val("1");
    $("#btn_search_single_item").trigger("click");
});