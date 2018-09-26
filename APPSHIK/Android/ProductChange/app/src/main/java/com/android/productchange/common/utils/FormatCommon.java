package com.android.productchange.common.utils;

import android.icu.text.SimpleDateFormat;

import com.android.productchange.common.constants.Constants;

import java.text.ParseException;

/**
 * Created by cong-pv
 * on 2018/09/06.
 */

public class FormatCommon {

    /**
     * Function replace "," in "\n"
     */

    public String formatLocationIdNewLine(String str) {

        return str.replace(",", "\n");
    }


    /**
     * Function replace "\n" in ","
     */

    public String formatLocationIdComma(String str) {

        return str.replace("\n", ",");
    }

    /**
     * Function format date 1000/01/01 to BLANK
     */

    public String formatDateBlank(String str) {

        if (Constants.VALUE_DEFAULT_DATE.equals(str)) {
            return Constants.BLANK;
        }
        return str;
    }

    /**
     * Format date
     *
     * @param date date from list
     * @return date has formatted
     */
    public String formatDateShowList(String date) {

        if (date.length() < 8) return date;
        String dateFormatted = date.substring(0, 4) + "\n" + date.substring(4, 6) + "/"
                + date.substring(6, 8);
        return dateFormatted;
    }
}
