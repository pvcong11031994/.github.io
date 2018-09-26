package com.android.productchange.common.utils;

import com.android.productchange.common.helpers.Log4JHelper;

import org.apache.log4j.Logger;

/**
 * <h1>Log Manager</h1>
 *
 * Init Log for activity
 *
 * @author nhut-bm
 * @since 2018-01-15
 */
public class LogManager {
    /**
     * Logger
     */
    public static Logger logger;

    /**
     * Log info
     *
     * @param tag     is tag name
     * @param message is message info
     */
    public static void i(String tag, String message) {
        logger = Log4JHelper.getLogger(tag, false);
        logger.info(message);
    }

    /**
     * Log error
     *
     * @param tag     is tag name
     * @param message is message error
     */
    public static void e(String tag, String message) {
        logger = Log4JHelper.getLogger(tag, true);
        logger.error(message);
    }
}
