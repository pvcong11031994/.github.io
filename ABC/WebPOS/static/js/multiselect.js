$(".multiselect").multiselect({
    selectedText: function (numChecked, numTotal, checkedItems) {
        if (numChecked == 1) {
            for (var i = 0; i < checkedItems.length; i++) {
                return $(checkedItems[i]).siblings().text();
            }
        } else if (numChecked == numTotal) {
            return "すべて";
        } else {
            return "複数選択済";
        }

    },
    noneSelectedText: "選択してください",
    checkAllText: 'すべて',
    uncheckAllText: 'なし',
});
$("select").parent("td").find(".ui-multiselect").click(function () {
    editViewSelect();
});
$("select").parent("td").find(".ui-multiselect").find("span").click(function () {
    editViewSelect();
});
function editViewSelect() {
    var $element = $(".ui-multiselect-menu:visible");
    if ($element.is(':visible') && ( $(window).height() - ($element.offset().top - $(document).scrollTop())) < $element.height() + 30) {
        var $height = $element.height() + 30;
        var top = $element.position().top - $height;
        $element.css({"top": top});
    }
}
