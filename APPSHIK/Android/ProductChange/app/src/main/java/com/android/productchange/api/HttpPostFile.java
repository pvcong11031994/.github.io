package com.android.productchange.api;

import android.content.Context;
import android.os.AsyncTask;

import com.android.productchange.interfaces.HttpResponse;

import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.InputStreamReader;
import java.net.URL;

import javax.net.ssl.HttpsURLConnection;

/**
 * <h1>Http Post File</h1>
 *
 * Task Post File log to API Log
 *
 * @author tien-lv
 * @since 2018-01-11.
 */

public class HttpPostFile extends AsyncTask<String, String, String> {

    /**
     * Http response
     */
    private HttpResponse response;

    /**
     * File name
     */
    private String fileName = "";

    /**
     * Constructor HttpPost
     *
     * @param c context
     */
    public HttpPostFile(Context c) {

        this.response = (HttpResponse) c;
    }

    /**
     * Set progress dialog loading.
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

        String lineEnd = "\r\n";
        String twoHyphens = "--";
        String boundary = "*****";
        int maxBufferSize = 1024 * 1024;
        int bytesRead, bytesAvailable, bufferSize;
        byte[] buffer;

        try {

            // set url from params
            URL url = new URL(params[0]);
            FileInputStream fileInputStream = new FileInputStream(
                    new File(params[1]));

            // init connection to server with https
            HttpsURLConnection httpsURLConnection = (HttpsURLConnection) url.openConnection();
            httpsURLConnection.setRequestMethod(Config.METHOD_POST);
            httpsURLConnection.setDoInput(true);
            httpsURLConnection.setDoOutput(true);
            httpsURLConnection.setUseCaches(false);
            httpsURLConnection.setRequestProperty(Config.PROPERTY_KEY,
                    Config.PROPERTY_VALUE_POST_FILE);

            httpsURLConnection.setRequestProperty(Config.CONNECTION_KEY, Config.CONNECTION_VALUE);
            httpsURLConnection.setRequestProperty(Config.ENCTYPE_KEY,
                    Config.PROPERTY_VALUE_POST_FILE);
            httpsURLConnection.setRequestProperty(Config.PROPERTY_KEY,
                    Config.PROPERTY_VALUE_POST_FILE + ";" + Config.BOUNDARY + "=" + boundary);
            httpsURLConnection.setRequestProperty(Config.UPLOADFILE, params[1]);
            httpsURLConnection.setRequestProperty(Config.API_KEY, Config.API_KEY_VALUE);


            DataOutputStream dataOutputStream = new DataOutputStream(
                    httpsURLConnection.getOutputStream());

            // writing bytes to data output stream
            dataOutputStream.writeBytes(twoHyphens + boundary + lineEnd);
            dataOutputStream.writeBytes(Config.CONTENT_DISPOSITION + params[1] + "\"" + lineEnd);
            dataOutputStream.writeBytes(lineEnd);

            // create a buffer of  maximum size
            bytesAvailable = fileInputStream.available();
            bufferSize = Math.min(bytesAvailable, maxBufferSize);
            buffer = new byte[bufferSize];
            bytesRead = fileInputStream.read(buffer, 0, bufferSize);

            while (bytesRead > 0) {
                dataOutputStream.write(buffer, 0, bufferSize);
                bytesAvailable = fileInputStream.available();
                bufferSize = Math.min(bytesAvailable, maxBufferSize);
                bytesRead = fileInputStream.read(buffer, 0, bufferSize);
            }

            // send multipart form data necessary after file data...
            dataOutputStream.writeBytes(lineEnd);
            dataOutputStream.writeBytes(twoHyphens + boundary + twoHyphens + lineEnd);

            fileInputStream.close();
            dataOutputStream.flush();
            dataOutputStream.close();

            int responseCode = httpsURLConnection.getResponseCode();

            fileName = params[1];
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

        response.progressFinish(result, 0, fileName);
    }

}
