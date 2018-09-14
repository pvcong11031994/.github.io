/**
 * Validate date
 */

// Do exec regex
function makeRegexParser(regex) {
    return regex.exec;
}

// Detect format date input
function getPatternDate(strFormat) {
    // Default is YYYYMMDD format
    var $patternRegex = /^(\d{4})(\d{2,2})(\d{2,2})$/;
    switch (strFormat.toUpperCase()) {
        case "YYYY/MM/DD":
            $patternRegex = /^(\d{4})\/(\d{2,2})\/(\d{2,2})$/;
            break;
        default :
            $patternRegex = /^(\d{4})(\d{2,2})(\d{2,2})$/;
    }
    return $patternRegex;
}

// Validate date : format
function validateDate(selectorItem) {
    // Get date format from input item : default YYYYMMDD if not exist
    var $strFormat = $(selectorItem).data("validate-format") == undefined ? "YYYYMMDD" : $(selectorItem).data("validate-format");
    // Generate regex to validate date
    var $patternDate = getPatternDate($strFormat);
    // Message for error
    var $msgError = " 期間が" + $strFormat + "形式でありません。";

    var strDate = $(selectorItem).val();
    // Step 1 : Don't validate blank data
    if (strDate == '') return "";

    // Step 2 : validate with pattern
    var dtArray = $patternDate.exec(strDate);
    if (dtArray == null) {
        return $msgError;
    }

    // Step 3 : Check valid date
    var dtYear = dtArray[1];
    var dtMonth = dtArray[2];
    var dtDay = dtArray[3];
    var date = new Date(dtYear, dtMonth - 1, dtDay);
    if (!(date.getFullYear() == dtYear && date.getMonth() + 1 == dtMonth && date.getDate() == dtDay)) {
        return $msgError;
    }
    return "";
}

/**
 * Validate date : range from-to
 *
 * @param $name     name of element need check
 * @param $title    Name item for message error
 */
function checkDateFromTo($name, $title) {
    //
    var $nameFrom = $name + "_from";
    var $nameTo = $name + "_to";

    // Validate from value
    var $msgError = "<p>" + $title + "(前)</p>";
    var $temp = validateDate($nameFrom);
    if ($temp != "") {
        $msgError += "<p>" + $temp + "</p>";
        return $msgError;
    }
    // Validate to value
    $msgError = "<p>" + $title + "(後)</p>";
    $temp = validateDate($nameTo);
    if ($temp != "") {
        $msgError += "<p>" + $temp + "</p>";
        return $msgError;
    }
    // Check from-to
    // Detect split char
    var $fromVal = $($nameFrom).val();
    var $toVal = $($nameTo).val();
    var $charSplit = $fromVal.replace(/[0-9]/gi, '').charAt(0);
    if ($fromVal != "" && $toVal != "" && $fromVal.split($charSplit).join("") > $toVal.split($charSplit).join("")) {
        $msgError = "<p>" + $title + "</p><p>開始日が終了日より後です。</p>";
        return $msgError;
    }
    return "";
}