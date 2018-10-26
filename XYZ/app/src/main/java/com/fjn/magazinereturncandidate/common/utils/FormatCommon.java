package com.fjn.magazinereturncandidate.common.utils;

import android.icu.text.DecimalFormat;
import android.icu.text.NumberFormat;
import android.icu.text.SimpleDateFormat;

import com.fjn.magazinereturncandidate.common.constants.Constants;

import java.text.ParseException;
import java.util.Locale;

/**
 * Function format common
 * Created by cong-pv on 2018/10/23.
 */

public class FormatCommon {

    /**
     * Function format money Japan
     *
     * @param money String value format
     * @return Result format success
     */
    public String formatMoney(String money) {

        //Convert String to Float
        NumberFormat numberFormat = NumberFormat.getNumberInstance(Locale.JAPANESE);
        String strFormat = money;
        String pattern = "#,###,###";
        DecimalFormat decimalFormat = (DecimalFormat) numberFormat;
        decimalFormat.applyPattern(pattern);
        try {
            strFormat = decimalFormat.format(Float.parseFloat(money));
        } catch (Exception e) {
            e.printStackTrace();
        }
        return strFormat;
    }

    /**
     * license
     * Format String date
     *
     * @param date date
     */
    public String formatDate(String date) {

        if (Constants.VALUE_DEFAULT_DATE.equals(date)) {
            return Constants.BLANK;
        }
        String result = date;
        SimpleDateFormat inputUser = new SimpleDateFormat("yyyyMMdd");
        SimpleDateFormat resultFormat = new SimpleDateFormat("yyyy/MM/dd");
        try {
            result = resultFormat.format(inputUser.parse(date));
        } catch (ParseException e) {
            e.printStackTrace();
        }
        return result;
    }

    /**
     * Function format number is 0008 => 8
     *
     * @param value string value input
     * @return Number format success
     */
    public String formatNumber(String value) {

        return String.valueOf(Integer.parseInt(value));
    }

    /**
     * Get time milliseconds current
     *
     * @return milliseconds current
     */
    public long getCurrentTimeLong() {

        return System.currentTimeMillis();
    }
}


