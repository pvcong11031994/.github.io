package com.fjn.magazinereturncandidate.common.utils;

import com.honeywell.barcode.HSMDecodeResult;
import com.honeywell.plugins.PluginResultListener;

/**
 * Created by cong-pv on 2018/10/26.
 */

public interface MyCustomPluginResultListener extends PluginResultListener {
    public void onMyCustomPluginResult(HSMDecodeResult[] barcodeData);
}
