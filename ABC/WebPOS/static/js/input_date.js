$(function () {
    $('.input_date').bindInputDate();
    $('.input_month').bindInputMonth();
    $('.input_week').bindInputWeek();
});

$.fn.bindInputDate = function () {
    $(this).each(function () {
        var dateFormat = $(this).data("dateformat") || $(this).attr("placeholder") || "yy/mm/dd";
        dateFormat = dateFormat.toLocaleLowerCase().replace(/yyyy/g, "yy");
        $(this).datepicker({
            dateFormat: dateFormat,
            changeMonth: true,
            changeYear: true,
            showButtonPanel: true,
            currentText: "当月",
            closeText: "閉じる",
            dayNames: ["日曜日", "月曜日", "火曜日", "水曜日", "木曜日", "金曜日", "土曜日"],
            dayNamesShort: ["日", "月", "火", "水", "木", "金", "土"],
            dayNamesMin: ["日", "月", "火", "水", "木", "金", "土"],
            monthNames: ["01月", "02月", "03月", "04月", "05月", "06月", "07月", "08月", "09月", "10月", "11月", "12月"],
            monthNamesShort: ["01月", "02月", "03月", "04月", "05月", "06月", "07月", "08月", "09月", "10月", "11月", "12月"]
        });
    });
};

$.fn.bindInputMonth = function () {
    $(this).each(function () {
        var dateFormat = $(this).data("dateformat") || $(this).attr("placeholder") || "yy/mm";
        dateFormat = dateFormat.toLocaleLowerCase().replace(/yyyy/g, "yy");

        $(this).MonthPicker({
            Button: false,
            MonthFormat:dateFormat,
            i18n:{
                year: "",
                prevYear: "前年",
                nextYear: "次年",
                next12Years: "",
                prev12Years: "",
                nextLabel: "次",
                prevLabel: "前",
                buttonText: "次",
                jumpYears: "年選択",
                backTo: "戻る",
                months: ["01月", "02月", "03月", "04月", "05月", "06月", "07月", "08月", "09月", "10月", "11月", "12月"],
            },
        });

    });
};

$.fn.bindInputWeek = function () {
    $(this).each(function () {
        var dateFormat = $(this).data("dateformat") || $(this).attr("placeholder") || "yy/mm/dd";
        dateFormat = dateFormat.toLocaleLowerCase().replace(/yyyy/g, "yy");
        $(this).datepicker({
            dateFormat: dateFormat,
            changeMonth: true,
            changeYear: true,
            showButtonPanel: true,
            currentText: "当月",
            closeText: "閉じる",
            dayNames: ["日曜日", "月曜日", "火曜日", "水曜日", "木曜日", "金曜日", "土曜日"],
            dayNamesShort: ["日", "月", "火", "水", "木", "金", "土"],
            dayNamesMin: ["日", "月", "火", "水", "木", "金", "土"],
            monthNames: ["01月", "02月", "03月", "04月", "05月", "06月", "07月", "08月", "09月", "10月", "11月", "12月"],
            monthNamesShort: ["01月", "02月", "03月", "04月", "05月", "06月", "07月", "08月", "09月", "10月", "11月", "12月"],
            onSelect: function() {
                var date = $(this).datepicker('getDate');
                date.setDate(date.getDate() - 1);
                var startDate = new Date(date.getFullYear(), date.getMonth(), date.getDate() - date.getDay() +1);
                var endDate = new Date(date.getFullYear(), date.getMonth(), date.getDate() - date.getDay() + 7);
                $(this).val($.datepicker.formatDate('yy/mm/dd', startDate)+"～"+$.datepicker.formatDate('mm/dd', endDate));
            },
        });

    });
};