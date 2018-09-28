package com.android.returncandidate.api;

import android.content.*;
import android.os.*;

import com.android.returncandidate.common.constants.*;
import com.android.returncandidate.common.utils.*;

import org.json.*;

import java.io.*;
import java.net.*;
import java.security.cert.*;

import javax.net.ssl.*;

/**
 * HTTP POST request
 *
 * @author tien-lv
 * @since 2017-12-18
 */

public class HttpPost extends AsyncTask<String, String, String> {

    /**
     * Http response.
     */
    private HttpResponse response;

    /**
     * Constructor HttpPost.
     */
    public HttpPost(Context c) {
        this.response = (HttpResponse) c;
    }

    /**
     * Set progress dialog loading.
     */
    protected void onPreExecute() {
    }

    private int typeLocation = 0;

    /**
     * Send request and get response to services API.
     */
    @Override
    protected String doInBackground(String... params) {

        try {

            URL url = new URL(params[1]);

            HttpsURLConnection httpsURLConnection = (HttpsURLConnection) url.openConnection();

            httpsURLConnection.setRequestMethod(Config.METHOD_POST);
            httpsURLConnection.setDoInput(true);
            httpsURLConnection.setDoOutput(true);
            httpsURLConnection.setUseCaches(false);
            httpsURLConnection.setRequestProperty(Config.PROPERTY_KEY, Config.PROPERTY_VALUE);

            JSONObject ob = setParams(params);

            OutputStream os = httpsURLConnection.getOutputStream();
            BufferedWriter writer = new BufferedWriter(
                    new OutputStreamWriter(os, Config.CHARSET_UTF_8));
            writer.write(ob.toString());
            writer.flush();
            writer.close();
            os.close();

            int responseCode = httpsURLConnection.getResponseCode();

            if (responseCode == HttpsURLConnection.HTTP_OK) {

                BufferedReader in = new BufferedReader(
                        new InputStreamReader(httpsURLConnection.getInputStream()));

                String line = in.readLine();

                return (line != null ? line : "");

            } else {
                return String.valueOf(responseCode);
            }
        } catch (Exception e) {
            return "Exception: " + e.getMessage();
        }
    }

    /**
     * End progress loading.
     */
    @Override
    protected void onPostExecute(String result) {
        response.progressFinish(result, typeLocation, null);
    }

    /**
     * Set request params.
     */
    private JSONObject setParams(String... params) {

        JSONObject jsonObject = new JSONObject();

        try {
            switch (params[0]) {
                case Config.CODE_LOGIN:
                    jsonObject.put(Constants.COLUMN_USER_ID, params[2]);
                    jsonObject.put(Constants.COLUMN_PASSWORD, params[3]);
                    break;
                case Config.CODE_GET_LIST_DATA:
                    jsonObject.put(Constants.COLUMN_SHOP_ID, params[2]);
                    jsonObject.put(Config.LOGIN_KEY, params[3]);
                    jsonObject.put(Constants.COLUMN_SERVER_NAME, params[4]);
                    jsonObject.put(Config.TYPE, params[5]);
                    typeLocation = Integer.parseInt(params[5]);
                    break;
                default:
                    break;
            }
            jsonObject.put(Config.API_KEY, Config.API_KEY_VALUE);

        } catch (JSONException e) {
            e.printStackTrace();
            LogManager.e(Constants.TAG_APPLICATION_NAME,e.toString());
        }
        return jsonObject;
    }
}
