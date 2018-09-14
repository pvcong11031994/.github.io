// $(function () {

    var MG1 = {};
    var MG12 = {};
    var MG123 = {};
    //var MG1234 = {};

    for (var i= 0; i < arrGenre.length; i++) {
        MG1[arrGenre[i].MediaGroup1Cd] = {Code: arrGenre[i].MediaGroup1Cd ,Name: arrGenre[i].MediaGroup1Name, MediaGroup2: arrGenre[i].MediaGroup2,JanGroup:arrGenre[i].JanGroup};
        for (var j = 0; j < arrGenre[i].MediaGroup2.length ; j++ ){
            MG12[arrGenre[i].MediaGroup1Cd + arrGenre[i].MediaGroup2[j].MediaGroup2Cd] = {
                Code: arrGenre[i].MediaGroup2[j].MediaGroup2Cd ,
                Name: arrGenre[i].MediaGroup2[j].MediaGroup2Name,
                MediaGroup3: arrGenre[i].MediaGroup2[j].MediaGroup3 };
            for (var k = 0; k < arrGenre[i].MediaGroup2[j].MediaGroup3.length ; k++ ){
                MG123[arrGenre[i].MediaGroup1Cd + arrGenre[i].MediaGroup2[j].MediaGroup2Cd + arrGenre[i].MediaGroup2[j].MediaGroup3[k].MediaGroup3Cd] = {
                    Code:arrGenre[i].MediaGroup2[j].MediaGroup3[k].MediaGroup3Cd,
                    Name:arrGenre[i].MediaGroup2[j].MediaGroup3[k].MediaGroup3Name,
                    MediaGroup4:arrGenre[i].MediaGroup2[j].MediaGroup3[k].MediaGroup4
                }
            }
        }
    }


    $('#media_group1_cd').on("change", function(e) {
        if($(this).val() != null){
            Get_MG2();
        }
        if (!$("#media_group1_cd option:selected").length) {
            $('#media_group2_cd').find('option').remove();
            $('#media_group3_cd').find('option').remove();
            $('#media_group4_cd').find('option').remove();
            $('#media_group2_cd').multiselect("refresh");
            $('#media_group3_cd').multiselect("refresh");
            $('#media_group4_cd').multiselect("refresh");
        }
    });

    $('#media_group2_cd').on("change", function(e) {
        if($(this).val() != null){
            Get_MG3();
        }
        if (!$("#media_group2_cd option:selected").length) {
            $('#media_group3_cd').find('option').remove();
            $('#media_group4_cd').find('option').remove();
            $('#media_group3_cd').multiselect("refresh");
            $('#media_group4_cd').multiselect("refresh");
        }
    });

    $('#media_group3_cd').on("change", function(e) {
        if($(this).val() != null){
            Get_MG4();
        }
        if (!$("#media_group3_cd option:selected").length) {
            $('#media_group4_cd').find('option').remove();
            $('#media_group4_cd').multiselect("refresh");
        }
    });

    function Get_MG2() {
        $('#media_group2_cd').find('option').remove();
        $( "#media_group1_cd option:selected" ).each(function() {
            var MG1Genre = $(this).val();
            for (var i = 0; i < MG1[MG1Genre].MediaGroup2.length; i++) {
                $('#media_group2_cd').append($('<option>', {
                    value: MG1[MG1Genre].MediaGroup2[i].MediaGroup2Cd,
                    text : MG1[MG1Genre].MediaGroup2[i].MediaGroup2Cd + " " + MG1[MG1Genre].MediaGroup2[i].MediaGroup2Name
                }));
            }
        });
        $('#media_group2_cd').multiselect("refresh");
    }
    function Get_MG3() {
        $('#media_group3_cd').find('option').remove();
        $( "#media_group1_cd option:selected" ).each(function() {
            var MG1Genre = $(this).val();
            $( "#media_group2_cd option:selected" ).each(function() {
                var MG12Genre = $(this).val();
                if ( MG12[MG1Genre + MG12Genre] != null) {
                    for (var i = 0; i < MG12[MG1Genre + MG12Genre].MediaGroup3.length; i++) {
                        $('#media_group3_cd').append($('<option>', {
                            value: MG12[MG1Genre + MG12Genre].MediaGroup3[i].MediaGroup3Cd,
                            text: MG12[MG1Genre + MG12Genre].MediaGroup3[i].MediaGroup3Cd + " " + MG12[MG1Genre + MG12Genre].MediaGroup3[i].MediaGroup3Name
                        }));
                    }
                }

            });
        });
        $('#media_group3_cd').multiselect("refresh");
    }
    function Get_MG4() {
        $('#media_group4_cd').find('option').remove();
        $( "#media_group1_cd option:selected" ).each(function() {
            var MG1Genre = $(this).val();
            $( "#media_group2_cd option:selected" ).each(function() {
                var MG12Genre = $(this).val();
                $( "#media_group3_cd option:selected" ).each(function() {
                    var MG123Genre = $(this).val();
                    if ( MG123[MG1Genre + MG12Genre + MG123Genre] != null) {
                        for (var i = 0; i < MG123[MG1Genre + MG12Genre + MG123Genre].MediaGroup4.length; i++) {
                            $('#media_group4_cd').append($('<option>', {
                                value: MG123[MG1Genre + MG12Genre + MG123Genre].MediaGroup4[i].MediaGroup4Cd,
                                text: MG123[MG1Genre + MG12Genre + MG123Genre].MediaGroup4[i].MediaGroup4Cd + " " + MG123[MG1Genre + MG12Genre + MG123Genre].MediaGroup4[i].MediaGroup4Name
                            }));
                        }
                    }
                });
            });
        });
        $('#media_group4_cd').multiselect("refresh");
    }

    function set_MG1_tab(control_type) {
        $('#media_group1_cd').find('option').remove();
        $('#media_group2_cd').find('option').remove();
        $('#media_group3_cd').find('option').remove();
        $('#media_group4_cd').find('option').remove();

        for (var i= 0; i < arrGenre.length; i++) {
            if(arrGenre[i].JanGroup == control_type) {
                var MG1Genre = arrGenre[i].MediaGroup1Cd
                $('#media_group1_cd').append($('<option>', {
                    value: MG1[MG1Genre].Code,
                    text: MG1[MG1Genre].Code + MG1[MG1Genre].Name
                }));
            }
        }

        $('#media_group1_cd').multiselect("refresh");
        $('#media_group2_cd').multiselect("refresh");
        $('#media_group3_cd').multiselect("refresh");
        $('#media_group4_cd').multiselect("refresh");
    }

    function setBackMG(group , value) {
        if (value != null) {
            if(value.length > 0){
                $( "#media_group"+group+"_cd option" ).each(function() {
                    var v = $(this).val();
                    $.each(value, function( index, val ) {
                        if(v == val){
                            $("#media_group"+group+"_cd option[value="+ val +"]").attr("selected","selected");
                        }
                    });
                });
                $("#media_group"+group+"_cd").multiselect("refresh");
            }
            if(group == "1") {
                Get_MG2();
            } else if(group == "2"){
                Get_MG3();
            } else if(group == "3"){
                Get_MG4();
            }
        }
    }
// });