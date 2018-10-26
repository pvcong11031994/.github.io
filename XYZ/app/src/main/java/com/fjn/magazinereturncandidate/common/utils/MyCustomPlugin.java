package com.fjn.magazinereturncandidate.common.utils;

import android.content.Context;
import android.graphics.Canvas;
import android.view.MotionEvent;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;

import com.fjn.magazinereturncandidate.R;
import com.honeywell.barcode.HSMDecodeResult;
import com.honeywell.plugins.DecodeBasePlugin;
import com.honeywell.plugins.PluginResultListener;
import com.honeywell.plugins.decode.PanoramicDecodePlugin;

import java.util.List;
import java.util.Map;

/**
 * Created by cong-pv on 2018/10/26.
 */

public class MyCustomPlugin extends PanoramicDecodePlugin {
    private TextView tvMessage;
    private int clickCount = 0;
    public MyCustomPlugin(Context context, Integer countJan) {
        super(context, countJan);

        //inflate the base UI layer
        View.inflate(context, R.layout.my_custom_plugin, this);

        tvMessage = (TextView) findViewById(R.id.textViewMsg);

        Button buttonHello = (Button) findViewById(R.id.buttonHello);
        buttonHello.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View arg0) {
                tvMessage.setText("Custom Plugin: UI button clicked " + ++clickCount + " times");
            }
        });
    }

    @Override
    protected void onStart() {
        super.onStart();
//do something
    }

    @Override
    protected void onStop() {
        super.onStop();
//do something
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
//do something
    }

    @Override
    protected void onDecode(HSMDecodeResult[] results) {
        super.onDecode(results);
//tells all plug-in monitor listeners we have a result
//this is used to signal HSMCameraPreview (if you are using it) that a result has been found
//and control should be returned to the caller. This call is not necessary if you are using an HSMDecodeComponent.
        this.signalResultFound(1);
//notifies all plug-in listeners we have a result
        List<PluginResultListener> listeners = this.getResultListeners();
        for (PluginResultListener listener : listeners)
            ((MyCustomPluginResultListener) listener).onMyCustomPluginResult(results);
    }

    @Override
    protected void onDecodeFailed() {
        super.onDecodeFailed();
//do something
    }

    @Override
    protected void onImage(byte[] image, int width, int height) {
        super.onImage(image, width, height);
        //do something
    }

    @Override
    protected void onSizeChanged(int width, int height, int oldWidth, int oldHeight) {
        super.onSizeChanged(width, height, oldWidth, oldHeight);
//do something
    }

    @Override
    protected void onDraw(Canvas canvas) {
        super.onDraw(canvas);
    }

    @Override
    public boolean onTouchEvent(MotionEvent event) {
        return super.onTouchEvent(event);
    }
}
