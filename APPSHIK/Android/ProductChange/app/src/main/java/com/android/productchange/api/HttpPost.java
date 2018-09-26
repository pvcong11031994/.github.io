package com.android.productchange.api;

import android.content.Context;
import android.os.AsyncTask;

import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.utils.LogManager;
import com.android.productchange.interfaces.HttpResponse;

import org.json.JSONException;
import org.json.JSONObject;

import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.io.OutputStreamWriter;
import java.net.URL;

import javax.net.ssl.HttpsURLConnection;

/**
 * <h1>Http Post</h1>
 *
 * Task http Post to API Login
 *
 * @author tien-lv
 * @since 2017-12-18
 */

public class HttpPost extends AsyncTask<String, String, String> {

    /**
     * Http response
     */
    private HttpResponse response;

    /**
     * Multi thread count
     */
    private int multiThreadCount;

    /**
     * Constructor HttpPost
     *
     * @param c context
     */
    public HttpPost(Context c) {
        this.response = (HttpResponse) c;
        multiThreadCount = 0;
    }

    /**
     * Set progress dialog loading
     */
    protected void onPreExecute() {
    }

    /**
     * Send request and get response to services API
     *
     * @param params String params for activity
     * @return result from API
     * @throws Exception from request fail
     * @see Exception
     */
    @Override
    protected String doInBackground(String... params) {

        try {

            // set url from params
            URL url = new URL(params[1]);

            // init connection to server with https
            HttpsURLConnection httpsURLConnection = (HttpsURLConnection) url.openConnection();

            httpsURLConnection.setRequestMethod(Config.METHOD_POST);
            httpsURLConnection.setDoInput(true);
            httpsURLConnection.setDoOutput(true);
            httpsURLConnection.setUseCaches(false);
            httpsURLConnection.setRequestProperty(Config.PROPERTY_KEY,
                    Config.PROPERTY_VALUE);

            JSONObject ob = setParams(params);

            // request param to API
            OutputStream os = httpsURLConnection.getOutputStream();
            BufferedWriter writer = new BufferedWriter(
                    new OutputStreamWriter(os, Config.CHARSET_UTF_8));
            writer.write(ob.toString());
            writer.flush();
            writer.close();
            os.close();

            // get response code from API
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
     * End progress loading
     *
     * @param result from API
     */
    @Override
    protected void onPostExecute(String result) {
        response.progressFinish(result, multiThreadCount, null);
    }

    /**
     * Set request params
     *
     * @param params from activity
     * @return Json object from params
     * @throws JSONObject format json error
     * @see JSONObject
     */
    private JSONObject setParams(String... params) {

        JSONObject jsonObject = new JSONObject();
        try {
            // set param request
            switch (params[0]) {
                case Config.CODE_LOGIN:
                    jsonObject.put(Constants.COLUMN_USER_ID, params[2]);
                    jsonObject.put(Constants.COLUMN_PASSWORD, params[3]);
                    break;
                default:
                    break;
            }
            jsonObject.put(Config.API_KEY, Config.API_KEY_VALUE);
        } catch (JSONException e) {
            e.printStackTrace();
            LogManager.e(Constants.TAG_APPLICATION_NAME, e.toString());
        }

        return jsonObject;
    }
}
