package com.fjn.magazinereturncandidate.common.utils;

import com.fjn.magazinereturncandidate.common.helpers.Log4JHelper;

import org.apache.log4j.Logger;

/**
 * Init Log for activity
 *
 * @author cong-pv
 * @since 2018-10-15
 */

public class LogManagerCommon {

    private static Logger logger;

    public static void i(String tag, String message) {
        logger = Log4JHelper.getLogger(tag, false);
        logger.info(message);
    }

    public static void e(String tag, String message) {
        logger = Log4JHelper.getLogger(tag, true);
        logger.error(message);
    }
}
