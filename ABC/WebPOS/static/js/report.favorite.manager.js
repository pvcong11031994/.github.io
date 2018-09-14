/**
 * Common for each report
 *      - Input order
 *      - Drill down
 *
 * @since       2017/09/29
 * @author      Thang-NQ
 */

function deleteFavoriteData() {
    var listJan = "";
    $('#tbl_report_result:first tr#favorite-old').each(function (row, tr) {
        if ($(this).find('input').is(':checked')) {
            listJan +=  ($(tr).find('td[name="jan_code"]').text()) + ";"
        }
    });
    if(listJan != "") {
        var data = {
            "csrf.Token": $("#form_common input[name='csrf.Token']").val(),
            "list_jan": listJan,
        };

        $.ajax({
            url: location.pathname + "/query_delete_ajax",
            type: "POST",
            data: data,
        }).done(function (res) {
            window.scrollTo(0, 0);
            if (res.Success) {
                $(".query-success").html(res.Msg);
                $(".query-success").css("display","block")
                $(".query-err").html();
                $(".query-err").css("display","none")
                loadFavoriteData();
            } else {
                $(".query-err").html(res.Msg);
                $(".query-err").css("display","block")
                $(".query-success").html();
                $(".query-success").css("display","none")
                loadFavoriteData();
            }
        });
    }
}
function editFavoriteData() {
    if (insertFavoriteData()){
        updateFavoriteData();
    };
}
function updateFavoriteData() {
    var data = new Array();
    $('#tbl_report_result:first tr#favorite-old').each(function (row, tr) {
        if ($(this).find('input').is(':checked')) {
            if ($(this).find('input[name="priority_number"]').val() == "") {
                return
            } else {
                data.push({
                    "JanCode": $(tr).find('td[name="jan_code"]').text(),
                    "PriorityNumber": $(tr).find('input[name="priority_number"]').val(),
                    "Memo": $(tr).find('input[name="memo"]').val()
                })
            }
        }
    });
    if (data.length > 0) {
        data = JSON.stringify(data);
        var data = {
            "csrf.Token": $("#form_common input[name='csrf.Token']").val(),
            "data": data,
        };

        $.ajax({
            url: location.pathname + "/query_update_ajax",
            type: "POST",
            data: data,
        }).done(function (res) {
            window.scrollTo(0, 0);
            if (res.Success) {
                $(".query-success").html(res.Msg);
                $(".query-success").css("display","block")
                $(".query-err").html();
                $(".query-err").css("display","none")
                loadFavoriteData();
            } else {
                $(".query-err").html(res.Msg);
                $(".query-err").css("display","block")
                $(".query-success").html();
                $(".query-success").css("display","none")
                resetHighLightErrorRecord()
                highLightErrorRecord(res.JanError)
            }
        });
    }
}

function insertFavoriteData() {
    var data = new Array();
    var isUniqueError = false
    var isInvalidJanError = false

    resetHighLightErrorRecord()
    $('#tbl_report_result:first tr#favorite-new').each(function (row, tr) {
        if ($(this).find('input').is(':checked')) {
            var janCodeNew = $(this).find('input[name="jan_code_new"]').val()
            if (checkUniqueJan(janCodeNew)) {
                highLightErrorRecord(janCodeNew);
                isUniqueError = true;
            }
        }
    });
    if  (!isUniqueError){
        $('#tbl_report_result:first tr#favorite-new').each(function (row, tr) {
            if ($(this).find('input').is(':checked')) {
                var janCodeNew = $(this).find('input[name="jan_code_new"]').val()
                if (janCodeNew.length < 13) {
                    highLightErrorRecord(janCodeNew);
                    isInvalidJanError = true;
                }
            }
        });
    }
    if (!isUniqueError && !isInvalidJanError) {
        $('#tbl_report_result:first tr#favorite-new').each(function (row, tr) {
            if ($(this).find('input').is(':checked')) {
                if ($(this).find('input[name="priority_number"]').val() != "" &&
                    $(this).find('input[name="jan_code_new"]').val() != "") {
                    data.push({
                        "JanCode": $(tr).find('input[name="jan_code_new"]').val(),
                        "PriorityNumber": $(tr).find('input[name="priority_number"]').val(),
                        "ProductName": $(tr).find('input[name="product_name"]').val(),
                        "AuthorName": $(tr).find('input[name="author_name"]').val(),
                        "MakerName": $(tr).find('input[name="maker_name"]').val(),
                        "UnitPrice": $(tr).find('input[name="unit_price"]').val(),
                        "ReleaseDate": $(tr).find('input[name="release_date"]').val(),
                        "Memo": $(tr).find('input[name="memo"]').val(),
                    })
                }
            }
        });
        if (data.length > 0) {
            data = JSON.stringify(data);
            var data = {
                "csrf.Token": $("#form_common input[name='csrf.Token']").val(),
                "data": data,
            };

            $.ajax({
                url: location.pathname + "/query_insert_ajax",
                type: "POST",
                data: data,
            }).done(function (res) {
                window.scrollTo(0, 0);
                if (res.Success) {
                    $(".query-success").html(res.Msg);
                    $(".query-success").css("display", "block")
                    $(".query-err").html();
                    $(".query-err").css("display", "none")
                    loadFavoriteData()
                } else {
                    $(".query-err").html(res.Msg);
                    $(".query-err").css("display", "block")
                    $(".query-success").html();
                    $(".query-success").css("display", "none")
                    resetHighLightErrorRecord()
                    highLightErrorRecord(res.JanError)
                }
            });
        }
    }else if(isUniqueError){
        window.scrollTo(0, 0);
        $(".query-err").html("入力したJANコードは既に登録されています。");
        $(".query-err").css("display", "block")
        $(".query-success").html();
        $(".query-success").css("display", "none")
        return false
    }else if(isInvalidJanError){
        window.scrollTo(0, 0);
        $(".query-err").html("JANコードは13桁で入力してください。");
        $(".query-err").css("display", "block")
        $(".query-success").html();
        $(".query-success").css("display", "none")
        return false
    }
    return true
}

function loadFavoriteData() {
    var data = {
        "csrf.Token": $("#form_common input[name='csrf.Token']").val(),
    };
    $.ajax({
        url: location.pathname + "/query_load_ajax",
        type: "POST",
        data: data,
        success: function (responeHtml) {
            $('#favorite-query-result').html(responeHtml).show();
            setTimeout(function () {
                $("#tbl_report_result").tablesorter();
                FixedMidashi.create();
                initRedirectJan($("table.query-result"));
                initRedirectGood($("table.query-result"));
                addRequriedToCheckbox();
            }, 100);
        },
    });
}

function addNewFavoriteRow()
{
    var table = document.getElementsByClassName('query-result')[0];
    var row = table.insertRow();
    var cellCheckBox = row.insertCell(0);
    var cellPriorityNumber = row.insertCell(1);
    var cellJan = row.insertCell(2);
    var cellProduct = row.insertCell(3);
    var cellAuthor = row.insertCell(4);
    var cellMaker = row.insertCell(5);
    var cellPrice = row.insertCell(6);
    var cellReleaseDate = row.insertCell(7);
    var cellMemo = row.insertCell(8);
    var cellCreateDate = row.insertCell(9);
    var cellUpdateDate = row.insertCell(10);


    $(row).attr("id","favorite-new")
    $(cellCheckBox).css("text-align","center");
    $(cellPriorityNumber).css("text-align","center");
    $(cellJan).css("text-align","center");
    $(cellProduct).css("text-align","center");
    $(cellAuthor).css("text-align","center");
    $(cellMaker).css("text-align","center");
    $(cellPrice).css("text-align","center");
    $(cellReleaseDate).css("text-align","center");
    $(cellMemo).css("text-align","center");

    cellCheckBox.innerHTML = "<label class='check_label'><input type='checkbox' name='checkboxJanNew' onclick='setRequired(this)' checked></label>";
    cellPriorityNumber.innerHTML = "<input name='priority_number' style='width: 50px; padding-right: 2px; padding-left: 2px; text-align: right' onkeypress='validate(event)' required/>";
    cellJan.innerHTML = "<input name='jan_code_new' maxlength='13' style='width: 95px; padding-right: 2px; padding-left: 2px;' onkeypress='validate(event)' required/>";
    cellProduct.innerHTML = "<input name='product_name' style='width: 100%; padding-right: 2px; padding-left: 2px;'/>";
    cellAuthor.innerHTML = "<input name='author_name' style='width: 150px; padding-right: 2px; padding-left: 2px;'/>";
    cellMaker.innerHTML = "<input name='maker_name' style='width: 100%; padding-right: 2px; padding-left: 2px;'/>";
    cellPrice.innerHTML = "<input name='unit_price' style='width: 46px; padding-right: 2px; padding-left: 2px; text-align: right' onkeypress='validate(event)' />";
    cellReleaseDate.innerHTML = "<input name='release_date' class='input_date' placeholder='YYYY/MM/DD' style='width: 85px; padding-right: 2px; padding-left: 2px;'/>";
    cellMemo.innerHTML = "<input name='memo' style='width: 165px; padding-right: 2px; padding-left: 2px;'/>";

    $('.input_date').bindInputDate();
    window.scrollTo(0,document.body.scrollHeight);
    FixedMidashi.create();
}

function copyJanToClipboard() {
    var janCodeList = ""
    $('#tbl_report_result:first tr#favorite-old input[type=checkbox]:checked').each(function () {
        var rowParent = $(this).parents().eq(2);
        janCodeList += $(rowParent).find('td[name="jan_code"]').text() + "\n";
    });
    $("#jan_cd_array_clipboard_temp").val('')
    $("#jan_cd_array_clipboard_temp").val(janCodeList)
    var element =  document.getElementById("jan_cd_array_clipboard_temp");
    element.select();
    document.execCommand('copy');
}

function addRequriedToCheckbox() {
    $('input[type=checkbox]').each(function() {
        if ( this.id.match("update-checkbox") ) {
            $(this).click(function() {
                setRequired(this)
            });
        }
    });
}

function setRequired(ele) {
    if(ele.checked) {
        var rowParent = $(ele).parents().eq(2);
        $(rowParent).find('input[name="priority_number"]').attr('required',true);
        $(rowParent).find('input[name="jan_code_new"]').attr('required',true);
    }else{
        var rowParent = $(ele).parents().eq(2);
        $(rowParent).find('input[name="priority_number"]').attr('required',false);
        $(rowParent).find('input[name="jan_code_new"]').attr('required',false);
    }
    FixedMidashi.create();
}

function highLightErrorRecord(jan) {
    $('#tbl_report_result tr').each(function (row, tr) {
        if ($(this).find('input').is(':checked')) {
            if ($(this).find('input[name="jan_code_new"]').val() == jan ||
                $(tr).find('td[name="jan_code"]').text() == jan) {
                $(tr).css("background-color", "rgb(242,83,113);")
            }
        }
    });
}

function resetHighLightErrorRecord() {
    $('#tbl_report_result:first tr').each(function (row, tr) {
        $(tr).css("background-color", "")
    });
}

function checkUniqueJan(jan) {
    var result = false
    $('#tbl_report_result:first tr').each(function (row, tr) {
        if ($(this).find('td[name="jan_code"]').text() == jan) {
            result = true
        }
    });
    return result;
}

function salesCompareClick() {
    //Event click sales_comparison
    if (!redirectSalesComparison()) {
        $(".query-err").html("商品は15件以内で選択してください。");
        $(".query-err").css("display", "block")
        $(".query-success").html();
        $(".query-success").css("display", "none")
        window.scrollTo(0, 0);
    } else {
        $(".query-err").html("");
        $(".query-err").css("display", "none")
    }
}

function initSalesCompareClick() {
    //Event click init_sales_compare
    if(!redirectInitSalses()){
        $(".query-err").html("商品は15件以内で選択してください。");
        $(".query-err").css("display", "block")
        $(".query-success").html();
        $(".query-success").css("display", "none")
        window.scrollTo(0,0);
    }else{
        $(".query-err").html("");
        $(".query-err").css("display", "none")
    }
}