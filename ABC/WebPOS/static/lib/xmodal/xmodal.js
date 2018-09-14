(function () {
    const backgroundID = 'xmodal_background';
    const backgroundClass = 'xmodal-background';
    const xModalClass = 'xmodal';
    const closeButtonClass = 'xmodal-btn-close';

    //Create xmodal background
    var background = $('#' + backgroundID);
    if (background) {
        background.remove();
    }
    background = $('<div>').attr('id', backgroundID).addClass(backgroundClass).appendTo('body');

    var observeCenterModal = null;

    var resizeTimer;
    $(window).on('resize', function () {
        if (observeCenterModal != null) {
            clearTimeout(resizeTimer);
            resizeTimer = setTimeout(function () {
                centering(observeCenterModal);
            }, 250);
        }
    });

    const centering = function (jQueryNode) {
        if (jQueryNode) {
            var top = (window.innerHeight - jQueryNode.outerHeight()) / 2;
            var left = (window.innerWidth - jQueryNode.outerWidth()) / 2;
            var newLocation = {top: top + 'px', left: left + 'px'};
            jQueryNode.css(newLocation);
        }
    };

    /*
     * Modal position is center default
     * click on child element has xmodal-btn-close class to close modal
     * to determine which button has been clicked, on the beforeClose and afterClose get e.closeButton
     * ops {
     *   outerClickClose: true/false (default false)
     *   beforeShow: event function(e)
     *   afterShown: event function(e)
     *   beforeClose: event function(e)
     *   afterClose: event function(e)
     * }
     * */
    $.fn.xmodal = function (ops) {
        var t = $(this);
        t.addClass(xModalClass);

        if (ops && ops.beforeShow) {
            ops.beforeShow(t);
        }
        centering(t);
        background.show();
        t.show();

        if (ops && ops.afterShown) {
            ops.afterShown(t);
        }
        centering(t);

        observeCenterModal = t;

        var tHide = t.find('.' + closeButtonClass);

        const closeProcess = function (e) {
            t.closeButton = e.target;

            if (ops && ops.beforeClose) {
                ops.beforeClose(t);
            }

            t.hide();

            if (ops && ops.afterClose) {
                ops.afterClose(t);
            }

            observeCenterModal = null;
            background.hide();

            tHide.unbind('click', closeProcess);
            if (ops && ops.outerClickClose) {
                background.unbind('click', closeProcess);
            }
        };

        tHide.click(closeProcess);

        if (ops && ops.outerClickClose) {
            background.click(closeProcess);
        }
    };
})();