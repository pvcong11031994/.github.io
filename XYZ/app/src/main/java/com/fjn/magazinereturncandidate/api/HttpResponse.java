package com.fjn.magazinereturncandidate.api;

/**
 * HTTP response interface
 *
 * @author cong-pv
 * @since 2018-10-15
 */

public interface HttpResponse {
    void progressFinish(String output, int typeLocation, String fileName);
}
