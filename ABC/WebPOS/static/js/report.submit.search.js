/**
 * Common for each report
 *      テキストボックスにフォーカスがあるときに、Enterキーを押したら、
 *      検索を実行して欲しい。（集計実行ボタンを押した時と同じ動きをする。）
 *
 * @since       2017/09/04
 * @author      Thang-NQ
 */

$(document).on("keydown", "input[type=text]", function (e) {
    if (13 == e.keyCode ) {
        setTimeout(function () {
            if ($("#input_max_value").is(":focus")) {
                return;
            } else {
                $("#btn_search").trigger("click");
            }
        },100);   
    }
});