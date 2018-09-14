$('input.input-number').each(function () {
    var $THIS = $(this);

    var max = $THIS.data("max");
    var min = $THIS.data("min");

    var lenmax = 0;
    var lenmin = 0;

    if (typeof max == "number") {
        lenmax = max.toString().length;
    }
    if (typeof min == "number") {
        lenmin = min.toString().length;
    }

    if (lenmin > lenmax) {
        lenmax = lenmin;
    }
    if (lenmax > 0) {
        $THIS.attr("maxlength", lenmax);
    }
});

$(document).on('keydown', 'input.input-number', function (e) {
    var $THIS = $(this);

    var min = $THIS.data("min");
    var enableMinus = false;
    if (typeof min == "number" && min < 0) {
        enableMinus = true;
    }

    // Allow: backspace, delete, tab, escape, enter and .
    if ($.inArray(e.keyCode, [46, 8, 9, 27, 13]) !== -1 ||
        // Allow: Ctrl+A
        (e.keyCode == 65 && e.ctrlKey === true) ||
        // Allow: Ctrl+C
        (e.keyCode == 67 && e.ctrlKey === true) ||
        // Allow: Ctrl+X
        (e.keyCode == 88 && e.ctrlKey === true) ||
        // Allow: home, end, left, right
        (e.keyCode >= 35 && e.keyCode <= 39) ||
        // Allow: F1->F12
        (e.keyCode >= 112 && e.keyCode <= 123) ||
        // Allow: start with minus
        (enableMinus && (e.keyCode == 109 || e.keyCode == 189) && $THIS.val().length == 0)) {
        // let it happen, don't do anything
        return;
    }
    // Ensure that it is a number and stop the keypress
    if ((e.shiftKey || (e.keyCode < 48 || e.keyCode > 57)) && (e.keyCode < 96 || e.keyCode > 105)) {
        e.preventDefault();
    }
});

$(document).on('change', 'input.input-number', function (e) {
    var $THIS = $(this);

    if ($THIS.val() == "") {
        return;
    }

    var max = $THIS.data("max");
    var min = $THIS.data("min");

    if (typeof max == "number") {
        if ($THIS.val() > max) {
            $THIS.val(max);
        }
    }
    if (typeof min == "number") {
        if ($THIS.val() < min) {
            $THIS.val(min);
        }
    }
});

function dropTest() {
    var elem = $(this);
    setTimeout(function() {
        // gets the copied text after a specified time (100 milliseconds)
        var janArrayValue = elem.val();
        if(!isNaN(janArrayValue))
        {
            return
        }else{
            elem.val("")
        }
        elem.trigger("input");
    }, 100);
}
$(document).on('drop', 'input.input-number', dropTest);
