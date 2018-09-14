var $validateBefore = null;

$(function () {
    if ($("#dialog_loading").length == 0) {
        $("<div>").attr("id", "dialog_loading")
            .append($("<span>").addClass("processing-text").text("集計実行中 ..."))
            .append($("<p>").attr("id", "processing_time"))
            .hide().appendTo($("body"));
    }

    var timeCounter = function () {
        var obj = new Object();

        obj.start = function () {
            obj.time_begin = (new Date()).getTime();

            obj.time_interval = setInterval(function () {
                obj.time_now = (new Date()).getTime();
                obj.time_elapse = obj.time_now - obj.time_begin;

                $('#processing_time').text((obj.time_elapse / 1000).toFixed(1) + ' 秒');

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

    var loadingDialog = function () {
        var obj = new Object();

        obj.dialog = $('#dialog_loading').dialog({
            width: 280,
            height: 180,
            modal: true,
            resizable: false,
            closeOnEscape: false,
            autoOpen: false,
            buttons: {
                'キャンセル': function () {
                    queryHandler.cancel();
                    loadingDialog.close();
                    timeCounter.stop();
                }
            },
            open: function () {
                $("#dialog_loading").parent().find(".ui-dialog-titlebar-close").remove();
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
                    var csvDownload = $("input[name=search_handle_type]").val();
                    $('#processing_time').html('データを受信しました。<br/>ブラウザでレンダリング...');
                    $("#form_search").find("input[name=search_handle_type]").val("");
                    setTimeout(function () {
                        // check session timeout
                        var pattent = /<input type=hidden name=flag_is_login_screen_name id=flag_is_login_screen_id value=flag_is_login_screen_value>/
                        var pattent2 = /<input type="hidden" name="flag_is_login_screen_name" id="flag_is_login_screen_id" value="flag_is_login_screen_value">/

                        if (!pattent.test(responeHtml) && !pattent2.test(responeHtml)) {
                            $('#query_result').html(responeHtml).show();
                        }
                        loadingDialog.close();
                        timeCounter.setNow();
                        console.log('BROWSER: Render html time: ' + timeCounter.htmlTime() + ' seconds');
                        // if (typeof initCommon === "function") {
                        //     //What common?
                        //     initCommon();
                        // }
                    }, 100);
                    setTimeout(function () {
                        if ($('#tbl_report_result').length == 0 && $("table.feeze-header").length == 0){
                            return;
                        }
                        if (location.pathname == "/report/best_sales_goods_transition") {
                            $("#tbl_report_result").tablesorter();
                            FixedMidashi.create();
                            initCommonSingleItemTransition();
                        } else if (location.pathname == "/report/best_sales") {
                            $("#tbl_report_result").tablesorter();
                            $("#tbl_report_result").floatThead({top : 101});
                            initCommonSingleGoods();
                        }  else if (location.pathname == "/report/best_sales_maria") {
                            $("#tbl_report_result").tablesorter();
                            $("#tbl_report_result").floatThead({top : 101});
                            initCommonSingleGoods();
                        } else if (location.pathname == "/report/best_sales_cloud") {
                            $("#tbl_report_result").tablesorter();
                            $("#tbl_report_result").floatThead({top : 101});
                            initCommonSingleGoods();
                        } else if(location.pathname == "/report/sales_comparison") {
                            $("#tbl_report_result").tablesorter();
                            FixedMidashi.create();
                            initCommonSingleGoods();
                        } else if(location.pathname == "/report/init_sales_compare"){
                            $("#tbl_report_result").tablesorter();
                            FixedMidashi.create();
                            initCommonSingleGoods();
                        } else if(location.pathname == "/report/single_goods_cumulative") {
                            FixedMidashi.create();
                        } else if(location.pathname == "/report/single_goods_stock_x") {
                            FixedMidashi.create();
                        } else if (location.pathname == "/report/search_goods") {
                            $("#tbl_report_result").tablesorter();
                            FixedMidashi.create();
                            initCommonSingleGoods();
                        } else if ($("table.feeze-header")) {
                            $("#tbl_report_result").tablesorter();
                            FixedMidashi.create();
                        }
                    },100);
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
                    loadingDialog.close();
                }
            });
        };

        obj.cancel = function () {
            if (obj.queryAjax) {
                obj.queryAjax.abort();
                $("#form_search").find("input[name=search_handle_type]").val("");
            }
        };

        return obj;
    }();


    $("#btn_search").click(function (event) {
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

        if ($("input[name='flag']").val() == ""){
            $("#search_condition_area input[name=page]").val(1);
        } else {
            $("input[name='flag']").val("");
        }
        $('#query_result').html("").hide();
        loadingDialog.show();
    });

    $("body").on("click", "#page_prev", function () {
        var nextVal = Number($("#page_select").val()) - 1;
        if ($("#page_select option[value=" + nextVal + "]").length > 0) {
            $("#page_select").val(nextVal);
            changePage();
        }
    });

    $("body").on("click", "#page_next", function () {
        var nextVal = Number($("#page_select").val()) + 1;
        if ($("#page_select option[value=" + nextVal + "]").length > 0) {
            $("#page_select").val(nextVal);
            changePage();
        }
    });

    $("body").on("change", "#page_select", changePage);

    function changePage() {
        var val = Number($("#page_select").val());
        $("#search_condition_area input[name=page]").val(val);

        $('#query_result').html("").hide();
        loadingDialog.show();
    }
});
