/**
 * Created by thang-nq on 17/10/13.
 */
//Event click init_sales_compare
function redirectInitSalses() {
    var $urlSalesComparisonsAjax = "/report/init_sales_compare";
    var selectedElement = $("input[name=checkboxJan]:checked")
    if (selectedElement.length > 15){
        return false
    }else {
        var listJanSelected = "";
        $("input[name=checkboxJan]:checked").each(function () {
            var rowParent = $(this).parents().eq(2);
            var janCode = "\n" + $(rowParent).find('td[name="jan_code"]').text();
            listJanSelected += janCode;
        });
        $("textarea[name=jan_cd_array]").val(listJanSelected);
        $("input[name=key_search]").val("1");
        var $form = $("#form_search");
        $form.attr("method","POST");
        $form.attr("action",$urlSalesComparisonsAjax);
        $form.attr("target","_blank");
        $form.submit();
    }
    return true
}

//Event click sales_comparison
function redirectSalesComparison() {
    var $urlSalesComparisonsAjax = "/report/sales_comparison";
    var selectedElement = $("input[name=checkboxJan]:checked")
    if (selectedElement.length > 15){
        return false
    }else {
        var listJanSelected = "";
        $("input[name=checkboxJan]:checked").each(function () {
            var rowParent = $(this).parents().eq(2);
            var janCode = "\n" + $(rowParent).find('td[name="jan_code"]').text();
            listJanSelected += janCode;
        });
        $("textarea[name=jan_cd_array]").val(listJanSelected);
        $("input[name=key_search]").val("1");
        var $form = $("#form_search");
        $form.attr("method","POST");
        $form.attr("action",$urlSalesComparisonsAjax);
        $form.attr("target","_blank");
        $form.submit();
    }
    return true
}