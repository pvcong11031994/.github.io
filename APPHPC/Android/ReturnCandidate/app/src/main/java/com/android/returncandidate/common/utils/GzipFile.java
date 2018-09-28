package com.android.returncandidate.common.utils;

import java.io.*;
import java.util.zip.*;

/**
 * GZip utilities
 *
 * @author tien-lv
 * @since 2018-01-25
 */

public class GzipFile {

    /**
     * Compress file log to gzip
     *
     * @param file     path file src
     * @param gzipFile path file des
     */
    public void compressGzipFile(String file, String gzipFile) {
        try {
            FileInputStream fis = new FileInputStream(file);
            FileOutputStream fos = new FileOutputStream(gzipFile);
            GZIPOutputStream gzipOS = new GZIPOutputStream(fos);
            byte[] buffer = new byte[1024];
            int len;
            while ((len = fis.read(buffer)) != -1) {
                gzipOS.write(buffer, 0, len);
            }
            // Close resources
            gzipOS.close();
            fos.close();
            fis.close();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
