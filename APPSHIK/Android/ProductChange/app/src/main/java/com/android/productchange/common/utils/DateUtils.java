package com.android.productchange.common.utils;

import com.android.productchange.common.constants.Constants;

import java.util.Calendar;

/**
 * Created by ShirlyKadosh on 4/20/17.
 */

public class DateUtils {

    /**
     * Return a string with the month name like it appears in the Julian and Gregorian calendars.
     *
     * @param month value is 0-based: 0 for January, 11 for December.
     * @return the month name like it appears in the Julian and Gregorian calendars as a string.
     */
    public static String monthToString(int month) {
        switch (month) {
            case Calendar.JANUARY:
                return Constants.JANUARY;

            case Calendar.FEBRUARY:
                return Constants.FEBRUARY;

            case Calendar.MARCH:
                return Constants.MARCH;

            case Calendar.APRIL:
                return Constants.APRIL;

            case Calendar.MAY:
                return Constants.MAY;

            case Calendar.JUNE:
                return Constants.JUNE;

            case Calendar.JULY:
                return Constants.JULY;

            case Calendar.AUGUST:
                return Constants.AUGUST;

            case Calendar.SEPTEMBER:
                return Constants.SEPTEMBER;

            case Calendar.OCTOBER:
                return Constants.OCTOBER;

            case Calendar.NOVEMBER:
                return Constants.NOVEMBER;

            case Calendar.DECEMBER:
                return Constants.DECEMBER;
        }
        return "";
    }

    /**
     * Return a string with the day name like it appears in the Julian and Gregorian calendars.
     *
     * @param day is the day of the week value.
     * @return the day of week name like it appears in the Julian and Gregorian calendars as a
     * string.
     */
    public static String dayOfWeekToString(int day) {
        switch (day) {
            case Calendar.SUNDAY:
                return Constants.SUNDAY;

            case Calendar.MONDAY:
                return Constants.MONDAY;

            case Calendar.TUESDAY:
                return Constants.TUESDAY;

            case Calendar.WEDNESDAY:
                return Constants.WEDNESDAY;

            case Calendar.THURSDAY:
                return Constants.THURSDAY;

            case Calendar.FRIDAY:
                return Constants.FRIDAY;

            case Calendar.SATURDAY:
                return Constants.SATURDAY;
        }
        return "";
    }
}
