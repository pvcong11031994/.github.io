$('select.select-year').each(function () {
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

    for (var $y = $from; $y <= $to; $y++) {
        var newOption = $("<option>").val($y).text($y.toString() + $text);
        if ($y == $val) newOption.prop("selected", true);
        $this.append(newOption);
    }
});

$('select.select-month').each(function () {
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
        var newOption = $("<option>").val($m).text($m.toString() + $text);
        if ($m == $val) newOption.prop("selected", true);
        $this.append(newOption);
    }
});


$('select.select-day').each(function () {
    var $this = $(this);
    var $val = $this.data("val");
    var $text = $this.data("text");
    if (typeof $text == "undefined") {
        $text = "";
    } else if (typeof $text != "string") {
        $text = $text.toString();
    }

    var $thisDate = (new Date()).getDate();
    if ($val == "now") {
        $val = $thisDate;
    }
    for (var $d = 1; $d <= 31; $d++) {
        var newOption = $("<option>").val($d).text($d.toString() + $text);
        if ($d == $val) newOption.prop("selected", true);
        $this.append(newOption);
    }
});
