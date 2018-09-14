$(function () {

    $("#media_group1_cd").on("change", function () {

        var selectedG1 = $(this).val();

        if (selectedG1 == null || selectedG1.length == 0) {
            selectedG1 = $.map($("#media_group1_cd option"), function (option) {
                return option.value;
            });
            $("#media_group2_cd option,#media_group3_cd option,#media_group4_cd option").each(function () {
                $(this).prop("selected", false);
            });
        }

        $("#media_group2_cd option,#media_group3_cd option,#media_group4_cd option").each(function () {
            var g1 = $(this).val().substr(0, 2);
            if ($.inArray(g1, selectedG1) !== -1) {
                $(this).prop("disabled", false);
            } else {
                $(this).prop("selected", false);
                $(this).prop("disabled", true);
            }
        });
        $('#media_group2_cd').multiselect("refresh");
        $('#media_group3_cd').multiselect("refresh");
        $('#media_group4_cd').multiselect("refresh");
    });

    $("#media_group2_cd").on("change", function () {
        var selectedG2 = $(this).val();

        if (selectedG2 == null || selectedG2.length == 0) {
            selectedG2 = $.map($("#media_group2_cd option"), function (option) {
                return option.value;
            });
            $("#media_group3_cd option,#media_group4_cd option").each(function () {
                $(this).prop("selected", false);
            });
        }

        $("#media_group3_cd option,#media_group4_cd option").each(function () {
            var g2 = $(this).val().substr(0, 2);
            if ($.inArray(g2, selectedG2) !== -1) {
                $(this).prop("disabled", false);
            } else {
                $(this).prop("selected", false);
                $(this).prop("disabled", true);
            }
        });

        $('#media_group3_cd').multiselect("refresh");
        $('#media_group4_cd').multiselect("refresh");
    });

    $("#media_group3_cd").on("change", function () {
        var selectedG3 = $(this).val();

        if (selectedG3 == null || selectedG3.length == 0) {
            selectedG3 = $.map($("#media_group3_cd option"), function (option) {
                return option.value;
            });
            $("#media_group4_cd option").each(function () {
                $(this).prop("selected", false);
            });
        }

        $("#media_group4_cd option").each(function () {
            var g3 = $(this).val().substr(0, 4);
            if ($.inArray(g3, selectedG3) !== -1) {
                $(this).prop("disabled", false);
            } else {
                $(this).prop("selected", false);
                $(this).prop("disabled", true);
            }
        });

        $('#media_group4_cd').multiselect("refresh");
    });
});