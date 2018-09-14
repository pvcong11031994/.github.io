function insertReportMenu(userMenu) {
    var ul = null;
    for (var i = 0; i < userMenu.length; i++) {
        var menu = userMenu[i];
        var href = "/report/" + menu.ReportId;
        if (ul == null || ul.data("reportId") != menu.ReportId) {
            var a = $("a[href='" + href + "']");

            if (a.length > 0) {
                var li = a.parent();
                if (!li.hasClass("expand")) {
                    li.addClass("expand");
                }

                ul = li.children("ul");
                if (ul.length == 0) {
                    li.append($("<ul>"));
                    ul = li.children("ul");
                }

                ul.data("reportId", menu.ReportId);
            } else {
                ul = null;
            }
        }
        if (ul != null && ul.data("reportId") == menu.ReportId) {
            ul.append(
                $("<li>").addClass("custom-report-menu").append($("<a>").attr("href", href + "?menu=" + menu.MenuId).text(menu.ReportName))
                    .append($("<span>").addClass("btn-delete-menu"))
            );
        }
    }
}

$(function () {
    $("body").on("click", ".btn-delete-menu", function () {

        var li = $(this).parent();
        var a = li.children("a");
        var href = a.attr("href");

        var menuName = a.text();
        var menuId = href.match(/[\?&]menu=(\d+)/)[1];
        var reportId = href.match(/^\/report\/([^\/?&]+)/)[1];

        $(window).subpage({
            reInit: true,
            type: "yesNo",
            title: "メニュー削除確認",
            message: "「" + menuName + "」メニューを削除します、よろしいですか。",
            onYes: function () {

                var menuData = {
                    "csrf.Token": $("#form_common input[name='csrf.Token']").val(),
                    "report_id": reportId,
                    "menu_id": menuId
                };

                $.ajax({
                    url: "/maintenance/reportmenu/delete_menu_ajax",
                    type: "POST",
                    data: menuData,
                }).done(function (res) {
                    if (res.Success) {
                        showSuccess(res.Msg);

                        var ul = li.parent();
                        li.remove();
                        if (ul.children("li").length == 0) {
                            var expand = ul.parent();
                            ul.remove();
                            expand.removeClass("expand");
                        }

                        if (window.location.href.endsWith(href)) {
                            var noMenuIdHref = href.replace(/\?menu=\d+/, "");
                            window.location.replace(noMenuIdHref);
                        }
                    } else {
                        showError(res.Msg);
                    }
                });
            }
        });
    });
});