(function () {
    var handleInputArray = function () {
        var $THIS = $(this);
        var $parent = $THIS.parent();
        var $thisVal = $THIS.val();
        var attrName = $THIS.data('form-name');

        if ($.trim($thisVal) != '') {
            var $existEmptyInput = false;
            $parent.children('input[data-form-name="' + attrName + '"]').each(function () {
                if ($existEmptyInput)
                    return;
                if ($.trim($(this).val()) == '') {
                    $existEmptyInput = true;
                }
            });

            if (!$existEmptyInput) {
                $parent.append($THIS.clone().val("").removeAttr("name"));
            }
        } else {
            var $emptyNode = $THIS.clone().val("").removeAttr("name");
            var $existEmptyInput = false;
            $parent.children('input[data-form-name="' + attrName + '"]').each(function () {
                if ($.trim($(this).val()) == '') {
                    if (!$(this).is(":last-child")) {
                        $(this).remove();
                    } else {
                        $emptyNode = $(this);
                        $existEmptyInput = true;
                    }
                }
            });
            if (!$existEmptyInput) {
                $parent.append($emptyNode);
            }

            $emptyNode.focus();
        }
    };

    var handleInputArrayLimit = function () {
        var $THIS = $(this);
        var $parent = $THIS.parent();
        var $thisVal = $THIS.val();
        var attrName = $THIS.data('form-name');
        const LIMIT_JAN_INPUT = 10;

        if ($.trim($thisVal) != '') {
            var $existEmptyInput = false;
            var numberOfJan = $parent.children('input[data-form-name="' + attrName + '"]').length;
            $parent.children('input[data-form-name="' + attrName + '"]').each(function () {
                if ($existEmptyInput)
                    return;
                if ($.trim($(this).val()) == '') {
                    $existEmptyInput = true;
                }
            });

            if (!$existEmptyInput && numberOfJan < LIMIT_JAN_INPUT) {
                $parent.append($THIS.clone().val("").removeAttr("name"));
            }
        } else {
            var $emptyNode = $THIS.clone().val("").removeAttr("name");
            var $existEmptyInput = false;
            $parent.children('input[data-form-name="' + attrName + '"]').each(function () {
                if ($.trim($(this).val()) == '') {
                    if (!$(this).is(":last-child")) {
                        $(this).remove();
                    } else {
                        $emptyNode = $(this);
                        $existEmptyInput = true;
                    }
                }
            });
            if (!$existEmptyInput) {
                $parent.append($emptyNode);
            }

            $emptyNode.focus();
        }
    };

    $(document).on('input', 'input.input-array', handleInputArray);
    $(document).on('input', 'input.input-array-limit', handleInputArrayLimit);
})();
