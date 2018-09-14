$(function () {
    $("#btn_save_report_layout").click(function () {
        $.ajax({
            url: location.pathname + "/save_layout_ajax",
            type: "POST",
            data: $("#form_save_report_layout").serialize()
        }).done(function (res) {
            if (res.Success) {
                showSuccess(res.Msg);
            } else {
                showError(res.Msg);
            }
        });
    });
});

$(function () {
    $("#btn_save_menu_report").click(function () {
        $(".modal_menu_report").subpage({
            onStart: function () {
                $(".subpage-content .menu_report_name").focus();
            },
            onYes: function () {
                var reportName = $.trim($(".subpage-content .menu_report_name").val());
                $("#form_save_report_layout input[name=report_name]").val(reportName);

                if (reportName != "") {
                    $.ajax({
                        url: location.pathname + "/save_layout_ajax",
                        type: "POST",
                        data: $("#form_save_report_layout").serialize()
                    }).done(function (res) {
                        if (res.Success) {
                            showSuccess(res.Msg);
                            $("#sub_report_name").text(" - " + reportName);
                            var newUrl = window.location.href.replace(/\?menu=\d+/, "") + "?menu=" + res.MenuId;
                            history.replaceState(null, null, newUrl);
                            insertReportMenu([res]);
                        } else {
                            showError(res.Msg);
                        }
                    });
                }
            }
        });
    });
});
