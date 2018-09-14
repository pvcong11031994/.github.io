/**
 * Created 2017/06/13.
 */

$(function () {
    // event chang group_type
    $("input[name=group_type]").on("change", function () {
        if($("input[name=group_type]:checked").val() == '4') {
            $('#date_from').prop('disabled', true);
            $('#date_to').prop('disabled', true);
        } else {
            $('#date_from').prop('disabled', false);
            $('#date_to').prop('disabled', false);
        }

        // change show date seach
        switch ($("input[name=group_type]:checked").val()){
            case "1":
                $("tr.select_date").css("display","none");
                $("tr.select_week").css("display","table-row");
                $("tr.select_month").css("display","none");
                // ASO-5288 Fix bug No.01 - Edit - START
                // if ($("input[name=date_from]").val() != "") {
                //     $("input[name=week_from]").val(addDaysForWeek(0,0,0,$("input[name=date_from]").val()));
                //     $("input[name=date_from]").val("");
                // }
                // if ($("input[name=date_to]").val() != "") {
                //     $("input[name=week_to]").val(addDaysForWeek(0,0,-1,$("input[name=date_to]").val()));
                //     $("input[name=date_to]").val("");
                // }
                //
                // if ($("input[name=month_from]").val() != "") {
                //     $("input[name=week_from]").val(addDaysForWeek(0,0,0,$("input[name=month_from]").val()+"/01"));
                //     $("input[name=month_from]").val("");
                // }
                // if ($("input[name=month_to]").val() != "") {
                //     $("input[name=week_to]").val(addDaysForWeek(0,1,-2,$("input[name=month_to]").val()+"/01"));
                //     $("input[name=month_to]").val("");
                // }
                $("input[name=week_from]").val(addDaysForWeek(0,0,0,addDays(0,-1,0,"")));
                $("input[name=week_to]").val(addDaysForWeek(0,0,-1,addDays(0,0,0,"")));
                // ASO-5288 Fix bug No.01 - Edit - END
                break;
            case "2":
                var dateFrom;
                var dateTo;
                $("tr.select_date").css("display","none");
                $("tr.select_week").css("display","none");
                $("tr.select_month").css("display","table-row");

                // ASO-5288 Fix bug No.01 - Edit - START
                // if ($("input[name=date_from]").val() != "") {
                //     dateFrom = $("input[name=date_from]").val().split("/");
                //     $("input[name=month_from]").val(dateFrom[0]+"/"+dateFrom[1]);
                //     $("input[name=date_from]").val("");
                // }
                // if ($("input[name=date_to]").val() != "") {
                //     dateTo = $("input[name=date_to]").val().split("/");
                //     $("input[name=month_to]").val(dateTo[0]+"/"+dateTo[1]);
                //     $("input[name=date_to]").val("");
                // }
                //
                // if ($("input[name=week_from]").val() != "") {
                //     dateFrom = $("input[name=week_from]").val().split("/");
                //     $("input[name=month_from]").val(dateFrom[0]+"/"+dateFrom[1]);
                //     $("input[name=week_from]").val("");
                // }
                // if ($("input[name=week_to]").val() != "") {
                //     dateTo = $("input[name=week_to]").val().split("～");
                //     var dateFirst = dateTo[0].split("/")
                //     $("input[name=month_to]").val(dateFirst[0]+"/"+dateTo[1].split("/")[0]);
                //     $("input[name=week_to]").val("");
                // }
                dateFrom = addDays(0,-1,0,"").split("/");
                dateTo = addDays(0,0,0,"").split("/");
                $("input[name=month_from]").val(dateFrom[0]+"/"+dateFrom[1]);
                $("input[name=month_to]").val(dateTo[0]+"/"+dateTo[1]);
                // ASO-5288 Fix bug No.01 - Edit - END
                break;
            default:
                $("tr.select_date").css("display","table-row");
                $("tr.select_week").css("display","none");
                $("tr.select_month").css("display","none");
                // ASO-5288 Fix bug No.01 - Edit - START
                // if ($("input[name=month_from]").val() != "") {
                //     $("input[name=date_from]").val($("input[name=month_from]").val()+"/01");
                //     $("input[name=month_from]").val("");
                // }
                // if ($("input[name=month_to]").val() != "") {
                //     $("input[name=date_to]").val(addDays(0,1,-1,$("input[name=month_to]").val()+"/01"));
                //     $("input[name=month_to]").val("");
                // }
                //
                // if ($("input[name=week_from]").val() != "") {
                //     $("input[name=date_from]").val($("input[name=week_from]").val().split("～")[0]);
                //     $("input[name=week_from]").val("");
                // }
                // if ($("input[name=week_to]").val() != "") {
                //     var dateTo = $("input[name=week_to]").val().split("～");
                //     var dateFirst = dateTo[0].split("/")
                //     $("input[name=date_to]").val(dateFirst[0]+"/"+dateTo[1]);
                //     $("input[name=week_to]").val("");
                // }
                $("input[name=date_from]").val(addDays(0,-1,0,""));
                $("input[name=date_to]").val(addDays(0,0,0,""));
            // ASO-5288 Fix bug No.01 - Edit - END
        }
    });
});

/**
 * 過去7日
 * */
$(".select-past-seven-date").click(function (event) {
    $("input[name=date_from]").val(addDays(0,0,-7,""));
    $("input[name=date_to]").val(addDays(0,0,0,""));

    $("input[name=month_from]").val(addDaysForMonth(0,0,0));
    $("input[name=month_to]").val(addDaysForMonth(0,0,0));

    $("input[name=week_from]").val(addDaysForWeek(0,0,-7,""));
    $("input[name=week_to]").val(addDaysForWeek(0,0,0,""));
});
/**
 * 過去1ヶ月
 * */
$(".select-past-month").click(function (event) {
    $("input[name=date_from]").val(addDays(0,-1,0,""));
    $("input[name=date_to]").val(addDays(0,0,0,""));

    $("input[name=month_from]").val(addDaysForMonth(0,-1,0));
    $("input[name=month_to]").val(addDaysForMonth(0,0,0));

    $("input[name=week_from]").val(addDaysForWeek(0,-1,-1,""));
    $("input[name=week_to]").val(addDaysForWeek(0,0,0,""));
});
/**
 * 過去1年
 * */
$(".select-past-year").click(function (event) {
    $("input[name=date_from]").val(addDays(-1,0,0,""));
    $("input[name=date_to]").val(addDays(0,0,0,""));

    $("input[name=month_from]").val(addDaysForMonth(-1,0,0));
    $("input[name=month_to]").val(addDaysForMonth(0,0,0));

    $("input[name=week_from]").val(addDaysForWeek(-1,0,0,""));
    $("input[name=week_to]").val(addDaysForWeek(0,0,0,""));
});
/**
 * ○月
 * */
$(".select-current-month").click(function (event) {
    $currentDate = addDays(0,0,0,"").split("/");
    $("input[name=date_from]").val($currentDate[0] + "/" + $currentDate[1] + "/01" );
    $("input[name=date_to]").val(addDays(0,0,0,""));

    $("input[name=month_from]").val($currentDate[0] + "/" + $currentDate[1]);
    $("input[name=month_to]").val(addDaysForMonth(0,0,0,""));

    $("input[name=week_from]").val(addDaysForWeek(0,0,0,$currentDate[0] + "/" + $currentDate[1] + "/01"));
    $("input[name=week_to]").val(addDaysForWeek(0,0,0,""));
});
/**
 * ○-1月
 * */
$(".select-one-past-month").click(function (event) {
    $currentDate = addDays(0,-1,0,"").split("/");
    $("input[name=date_from]").val($currentDate[0] + "/" + $currentDate[1] + "/01" );
    $pastDate = new Date().getDate();
    $("input[name=date_to]").val(addDays(0,0,-$pastDate,""));


    $("input[name=month_from]").val($currentDate[0] + "/" + $currentDate[1]);
    $("input[name=month_to]").val(addDaysForMonth(0,0,-$pastDate));

    $("input[name=week_from]").val(addDaysForWeek(0,0,0,$currentDate[0] + "/" + $currentDate[1] + "/01"));
    $("input[name=week_to]").val(addDaysForWeek(0,0,-$pastDate,""));
});
/**
 * ○年
 * */
$(".select-current-year").click(function (event) {
    $currentDate = addDays(0,0,0,"").split("/");
    $("input[name=date_from]").val($currentDate[0] + "/01/01" );
    $("input[name=date_to]").val(addDays(0,0,0,""));

    $("input[name=month_from]").val($currentDate[0] + "/01");
    $("input[name=month_to]").val(addDaysForMonth(0,0,0));

    $("input[name=week_from]").val(addDaysForWeek(0,0,0,$currentDate[0] + "/01/01"));
    $("input[name=week_to]").val(addDaysForWeek(0,0,0,""));
});
/**
 * ○-1年
 * */
$(".select-one-past-year").click(function (event) {
    $currentDate = addDays(-1,0,0,"").split("/");
    $("input[name=date_from]").val($currentDate[0] + "/01/01" );
    $("input[name=date_to]").val($currentDate[0] + "/12/31");

    $("input[name=month_from]").val($currentDate[0] + "/01");
    $("input[name=month_to]").val($currentDate[0] + "/12");

    $("input[name=week_from]").val(addDaysForWeek(0,0,0,$currentDate[0] + "/01/01"));
    $("input[name=week_to]").val(addDaysForWeek(0,0,0,$currentDate[0] + "/12/31"));
});
/**
 * ○-2年
 * */
$(".select-two-past-year").click(function (event) {
    $currentDate = addDays(-2,0,0,"").split("/");
    $("input[name=date_from]").val($currentDate[0] + "/01/01" );
    $("input[name=date_to]").val($currentDate[0] + "/12/31");

    $("input[name=month_from]").val($currentDate[0] + "/01");
    $("input[name=month_to]").val($currentDate[0] + "/12");

    $("input[name=week_from]").val(addDaysForWeek(0,0,0,$currentDate[0] + "/01/01"));
    $("input[name=week_to]").val(addDaysForWeek(0,0,0,$currentDate[0] + "/12/31"));
});

function addDays(years, months , days, dateFormat) {
    $date = new Date();
    if (dateFormat != ""){
        $date = new Date(dateFormat);
    }
    if (isNaN($date.getDate())){
        return "";
    }
    $current_date = $date.getDate();
    $current_month = $date.getMonth() + 1;
    $current_year = $date.getFullYear();
    var result = new Date();
    result.setYear($date.getFullYear() + years);
    result.setMonth($date.getMonth() +  months);
    result.setDate($date.getDate() + days);

    $resultYear = result.getFullYear();
    $resultMonth = 0;
    $resultDay = 0;
    if( result.getMonth() < 9) {
        $resultMonth = "0" + (result.getMonth() + 1);
    } else {
        $resultMonth = result.getMonth() + 1;
    }
    if( result.getDate() < 10) {
        $resultDay = "0" + result.getDate();
    } else {
        $resultDay = result.getDate();
    }

    return $resultYear + "/" + $resultMonth + "/" + $resultDay;
}
function addDaysForMonth(years, months , days) {
    $date = new Date();
    if (isNaN($date.getDate())){
        return "";
    }
    $current_date = $date.getDate();
    $current_month = $date.getMonth() + 1;
    $current_year = $date.getFullYear();
    var result = new Date();
    result.setYear($date.getFullYear() + years);
    result.setMonth($date.getMonth() +  months);
    result.setDate($date.getDate() + days);

    $resultYear = result.getFullYear();
    $resultMonth = 0;
    $resultDay = 0;
    if( result.getMonth() < 9) {
        $resultMonth = "0" + (result.getMonth() + 1);
    } else {
        $resultMonth = result.getMonth() + 1;
    }
    if( result.getDate() < 10) {
        $resultDay = "0" + result.getDate();
    } else {
        $resultDay = result.getDate();
    }

    return $resultYear + "/" + $resultMonth;
}

function addDaysForWeek(years, months , days, dateFormat) {
    $date = new Date();
    if (dateFormat != ""){
        $date = new Date(dateFormat);
    }
    if (isNaN($date.getDate())){
        return "";
    }

    $current_date = $date.getDate();
    $current_month = $date.getMonth() + 1;
    $current_year = $date.getFullYear();
    var result = new Date();
    result.setYear($date.getFullYear() + years);
    result.setMonth($date.getMonth() +  months);
    result.setDate($date.getDate() + days);

    $resultYear = result.getFullYear();
    $resultMonth = 0;
    $resultDay = 0;
    if( result.getMonth() < 9) {
        $resultMonth = "0" + (result.getMonth() + 1);
    } else {
        $resultMonth = result.getMonth() + 1;
    }
    if( result.getDate() < 10) {
        $resultDay = "0" + result.getDate();
    } else {
        $resultDay = result.getDate();
    }

    var startDate = new Date(result.getFullYear(), result.getMonth(), result.getDate() - result.getDay() +1);
    var endDate = new Date(result.getFullYear(), result.getMonth(), result.getDate() - result.getDay() + 7);
    return $.datepicker.formatDate('yy/mm/dd', startDate)+"～"+$.datepicker.formatDate('mm/dd', endDate);
}