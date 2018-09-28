package com.android.returncandidate.common.utils;

import com.android.returncandidate.common.helpers.*;

import org.apache.log4j.*;

/**
 * Init Log for activity
 *
 * @author nhut-bm
 * @since 2018-01-15
 */

public class LogManager {
    public static Logger logger;

    public static void i(String tag, String message) {
        logger = Log4JHelper.getLogger(tag, false);
        logger.info(message);
    }

    public static void e(String tag, String message) {
        logger = Log4JHelper.getLogger(tag, true);
        logger.error(message);
    }
}
