$('input[data-form-name]').each(
    function () {
        var $THIS = $(this);
        if ($THIS.val() != '') {
            $THIS.attr('name', $THIS.data('form-name'));
        }
    }
);
$(document).on('change', 'input[data-form-name]', function () {
    var $THIS = $(this);

    if ($THIS.val() != '') {
        if ($THIS.attr('name') === undefined) {
            $THIS.attr('name', $THIS.data('form-name'));
        }
    } else {
        $THIS.removeAttr("name");
    }
});