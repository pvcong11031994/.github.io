package com.android.productchange.interfaces;

/**
 * <h1>Http Response</h1>
 *
 * Interface http response result
 *
 * @author tien-lv
 * @since 2017-12-22
 */
public interface HttpResponse {

    /**
     * Function get result from Asyntask and send to activity
     *
     * @param output           result from API
     * @param multiThreadCount thread count
     * @param fileName         file name
     */
    void progressFinish(String output, int multiThreadCount, String fileName);

}
