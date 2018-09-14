(function ($) {
    var $mainPanel = $("<div style='display: none;' class='opacity-overlay' tabindex='-1'>\
                            <div class='container'>\
                                <div class='subpage-wrapper'>\
                                    <div class='subpage-content'>\
                                    </div>\
                                </div>\
                            </div>\
                        </div>");

    var $yesNo = $("<div style='display: block;' class='confirm-delete modal fade in'>\
                        <div class='modal-dialog '>\
                            <div class='modal-content'>\
                                <div class='modal-header'>\
                                    <button type='button' class='close subpage-close'>×</button>\
                                    <h4 class='modal-title'>Modal title</h4>\
                                </div>\
                                <div class='modal-body modal-message'>\
                                    データを削除します、よろしいですか\
                                </div>\
                                <div class='modal-footer'>\
                                    <button type='button' class='btn btn-default subpage-yes'>はい</button>\
                                    <button type='button' class='btn btn-primary subpage-no'>いいえ</button>\
                                </div>\
                            </div>\
                        </div>\
                    </div>");

    var $messageShow = $("<div style='display: block;' class='show-message-box modal fade in'>\
                              <div class='modal-dialog '>\
                                  <div class='modal-content'>\
                                      <div class='modal-header'>\
                                        <button type='button' class='close subpage-close'>×</button>\
                                        <h4 class='modal-title'></h4>\
                                      </div>\
                                      <div class='modal-body modal-message overflow-multi-line'>\
                                      </div>\
                                      <div class='modal-footer'>\
                                          <button type='button' class='btn btn-primary subpage-yes'>はい</button>\
                                      </div>\
                                  </div>\
                              </div>\
                          </div>");

    $.fn.subpage = function (options) {
        if (!$("body").find(".opacity-overlay").length) {
            $("body").append($mainPanel);
        }
        options = $.extend({}, $.fn.subpage.config, options);
        return this.each(function () {

            reInitView(this, options);

            $(".subpage-content").html(options.element);
            $(".opacity-overlay").fadeIn();

            if (options.onStart != null) {
                options.onStart(options);
            }

            $(".subpage-content").find('.subpage-close').on('click', function () {
                callback(options.onClose, options);
            });
            $(".subpage-content").find('.subpage-yes').on('click', function () {
                if (options.onYes != null) {
                    callback(options.onYes);
                } else {
                    finishAll();
                }
            });
            $(".subpage-content").find('.subpage-no').on('click', function () {
                if (options.onNo != null) {
                    callback(options.onNo, options);
                } else {
                    finishAll();
                }
            });

            $.each(options.customAction, function ($key, $value) {
                $(".subpage-content").find('.' + $key).on('click', function () {
                    if ($value != null) {
                        $value(options);
                    }
                });
            });

            if (options.timer > 0) {
                setTimeout(function () {
                    finishAll();
                }, options.timer);
            }
        });
    };

    $.fn.subpageHide = function () {
        finishAll();
    }

    var callback = function (func, options) {
        func(options);
        finishAll();
    };

    function finishAll() {
        $(".opacity-overlay").fadeOut();
        $(".subpage-content").html("");
    }

    function isWindow(obj) {
        if (typeof(window.constructor) === 'undefined') {
            return obj instanceof window.constructor;
        } else {
            return obj.window === obj;
        }
    }

    $language = {
        "ja": {
            yes: "はい",
            no: "いいえ",
            close: "閉じる",
            cancel: "キャンセル",
            error: "エラー",
            warning: "警告",
            noti: "通知",
        },
        "en": {
            yes: "Yes",
            no: "No",
            close: "Close",
            cancel: "Cancel",
            error: "Error",
            warning: "Warning",
            noti: "Notification",
        }
    };

    $.fn.subpage.config = {
        type: 'message',
        reInit: false,
        title: '通知',
        message: '',
        locale: 'ja',
        timer: 0,
        element: $messageShow,
        customAction: {},
        onStart: null,
        onYes: null,
        onNo: null,
        onClose: function () {
            $(".opacity-overlay").fadeOut();
            $(".subpage-content").html("");
        },
    };

    function setInfo(options) {
        if (options.reInit) {
            options.element.find(".modal-title").html(options.title);
            options.element.find(".modal-message").html(options.message);
            options.element.find(".subpage-yes").html($language[options.locale]['yes']);
            options.element.find(".subpage-no").html($language[options.locale]['no']);
            options.element.find(".subpage-cancel").html($language[options.locale]['cancel']);
        }
    }

    function reInitView(item, options) {
        if (isWindow(item)) {
            switch (options.type) {
                case "yesNo":
                    options.element = $yesNo;
                    break;
                case "message":
                default:
                    options.element = $messageShow;
            }
        } else {
            options.element = $(item).clone();
        }
        setInfo(options);
    }
}(jQuery));

function subpage(options) {
    $(this).subpage(options);
}

function subpageHide() {
    $(this).subpageHide();
}

// Waiting
var timeCheck = 500;
var startTime;
var handleShow = true;
var handleHide = true;
var $waitingObj = $("<div class='backgroundSubpage' style='display:none;z-index:999999999'>\
                      <ul id='spinners'>\
                          <li class='circle-spinner selected'>\
                              <div class='spinner-container container1'>\
                                  <div class='circle1'></div>\
                                  <div class='circle2'></div>\
                                  <div class='circle3'></div>\
                                  <div class='circle4'></div>\
                              </div>\
                              <div class='spinner-container container2'>\
                                  <div class='circle1'></div>\
                                  <div class='circle2'></div>\
                                  <div class='circle3'></div>\
                                  <div class='circle4'></div>\
                              </div>\
                              <div class='spinner-container container3'>\
                                  <div class='circle1'></div>\
                                  <div class='circle2'></div>\
                                  <div class='circle3'></div>\
                                  <div class='circle4'></div>\
                              </div>\
                          </li>\
                      </ul>\
                  </div>");

function showError($err) {
    setTimeout(function () {
        hideWaiting(true);
        subpage({
            title: "エラー",
            message: $err,
            reInit: true
        });
    }, 500);
}
function initWaitingOnAjax() {
    handleShow = true;
    $(document).ajaxSend(function (event, jqXHR, ajaxOptions) {
        if (ajaxOptions.context) {
            if (ajaxOptions.context.nosubpage) {
                return;
            }
        }
        showWaiting();
    });
    $(document).ajaxStop(function () {
        hideWaiting();
    });
    $(document).ajaxComplete(function (event, xhr, settings) {
        hideWaiting();
        var pattent = /<input type=hidden name=flag_is_login_screen_name id=flag_is_login_screen_id value=flag_is_login_screen_value>/
        var pattent2 = /<input type="hidden" name="flag_is_login_screen_name" id="flag_is_login_screen_id" value="flag_is_login_screen_value">/
        if (pattent.test(xhr.responseText) || pattent2.test(xhr.responseText)){
            location.href = "/";
        }
    });

    $(document).ajaxError(function (event, XMLHttpRequest, ajaxOptions) {
        console.log(event)
        console.log(XMLHttpRequest)
        console.log(ajaxOptions)
    });
}

function showWaiting() {
    if ( !handleShow ) {
        return;
    }
    if (!$("body").find(".backgroundSubpage").length) {
        $("body").append($waitingObj);
    }
    startTime = new Date();
    $(".backgroundSubpage").show();
}

function hideWaiting($justNow) {
    if ( !handleHide ) {
        return;
    }
    $justNow = $justNow | false;
    if ($justNow) {
        removeWaiting();
    } else {
        var endTime = new Date();
        var secondsElapsed = (endTime - startTime);
        if (secondsElapsed <= timeCheck) {
            setTimeout(function () {
                removeWaiting();
            }, timeCheck);
        } else {
            removeWaiting();
        }
    }
}

function removeWaiting() {
    $(".backgroundSubpage").fadeOut();
    $(".backgroundSubpage").remove();
}
