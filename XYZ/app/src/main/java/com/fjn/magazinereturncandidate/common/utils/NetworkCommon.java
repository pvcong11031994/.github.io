package com.fjn.magazinereturncandidate.common.utils;

import android.app.Activity;
import android.content.Context;
import android.net.ConnectivityManager;
import android.net.NetworkInfo;

/**
 * Function common network
 * Created by cong-pv on 2018/10/17.
 */

public class NetworkCommon {

    /**
     * Check network
     *
     * @return This returns true if network connectivity is okay
     */
    public static boolean isNetworkConnection(Activity activity) {
        ConnectivityManager connectivityManager = (ConnectivityManager) activity.getSystemService(
                Context.CONNECTIVITY_SERVICE);
        NetworkInfo networkInfo = connectivityManager.getActiveNetworkInfo();
        return !(networkInfo == null || !networkInfo.isConnected() || !networkInfo.isAvailable());
    }
}
