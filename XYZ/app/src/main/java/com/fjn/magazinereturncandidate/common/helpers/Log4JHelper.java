package com.fjn.magazinereturncandidate.common.helpers;

import android.os.Environment;

import org.apache.log4j.Logger;

import de.mindpipe.android.logging.log4j.LogConfigurator;

/**
 * Init and config Log4J
 *
 * @author cong-pv
 * @since 2018-10-15
 */
public class Log4JHelper {

    private final static LogConfigurator mLogConfigrator = new LogConfigurator();

    public static String fileName = "";
    private static String filePatternInfo = "%d I/%c： %m%n";

    // Init for static request
    static {
        configureLog4j();
    }

    private static void configureLog4j() {

        fileName = Environment.getExternalStorageDirectory()
                + "/MagazineReturnCandidate/"
                + "log.log";
        String filePattern = filePatternInfo;
        int maxBackupSize = 10;
        long maxFileSize = 1024 * 1024;

        configure(fileName, filePattern, maxBackupSize, maxFileSize);
    }


    /**
     * Set config for log
     *
     * @param fileName      file name
     * @param filePattern   pattern log by line
     * @param maxBackupSize size file log backup
     * @param maxFileSize   size file log
     */
    private static void configure(String fileName, String filePattern, int maxBackupSize,
                                  long maxFileSize) {

        mLogConfigrator.setFileName(fileName);
        mLogConfigrator.setMaxFileSize(maxFileSize);
        mLogConfigrator.setFilePattern(filePattern);
        mLogConfigrator.setMaxBackupSize(maxBackupSize);
        mLogConfigrator.setUseLogCatAppender(true);
        mLogConfigrator.setUseFileAppender(true);
        mLogConfigrator.configure();
    }

    /**
     * Call logger
     *
     * @param name  name activity
     * @param isErr bool var when err
     * @return Logger
     */
    public static Logger getLogger(String name, boolean isErr) {

        // init pattern error
        String filePatternErr = "%d E/%c： %m%n";
        if (isErr) {
            mLogConfigrator.setFilePattern(filePatternErr);
        } else {
            mLogConfigrator.setFilePattern(filePatternInfo);
        }
        mLogConfigrator.configure();

        return org.apache.log4j.Logger.getLogger(name);
    }
}