package com.fjn.magazinereturncandidate.api;

import android.app.ProgressDialog;
import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.os.AsyncTask;

import com.fjn.magazinereturncandidate.common.constants.Constants;
import com.fjn.magazinereturncandidate.common.constants.Message;
import com.fjn.magazinereturncandidate.common.utils.DatabaseManagerCommon;
import com.fjn.magazinereturncandidate.common.utils.LogManagerCommon;
import com.fjn.magazinereturncandidate.db.entity.MaxYearRankEntity;
import com.fjn.magazinereturncandidate.db.models.MaxYearRankModel;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.stream.JsonReader;

import org.json.JSONException;
import org.json.JSONObject;

import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.io.OutputStreamWriter;
import java.net.URL;
import java.util.Calendar;
import java.util.Date;

import javax.net.ssl.HttpsURLConnection;

/**
 * <h1>Http Post User</h1>
 *
 * @author cong-pv
 * @since 2018-10-17
 * <p>
 */

public class HttpPostMaxYearRank extends AsyncTask<String, String, String> {

    /**
     * Http response
     */
    private HttpResponse response;

    /**
     * Progress dialog
     */
    private ProgressDialog progressDialog;

    /**
     * Constructor HttpPost
     *
     * @param c context
     */
    public HttpPostMaxYearRank(Context c) {

        this.response = (HttpResponse) c;
        progressDialog = new ProgressDialog(c);
    }

    /**
     * Set progress dialog loading
     */
    protected void onPreExecute() {

        progressDialog.setMessage(Message.MESSAGE_DOWNLOAD_DATA_SCREEN);
        progressDialog.setIndeterminate(false);
        progressDialog.setCancelable(false);
        progressDialog.setCanceledOnTouchOutside(false);
        progressDialog.show();
    }

    /**
     * Send request and get response to services API
     *
     * @param params String params for activity
     * @return result from API
     * @see IOException
     * @see Exception
     */
    @Override
    protected String doInBackground(String... params) {

        try {

            Date cr = Calendar.getInstance().getTime();

            // set url from params
            URL url = new URL(params[1]);

            // init connection to server with https
            HttpsURLConnection httpsURLConnection = (HttpsURLConnection) url.openConnection();

            httpsURLConnection.setRequestMethod(Config.METHOD_POST);
            httpsURLConnection.setDoInput(true);
            httpsURLConnection.setDoOutput(true);
            httpsURLConnection.setUseCaches(false);
            httpsURLConnection.setRequestProperty(Config.PROPERTY_KEY, Config.PROPERTY_VALUE);

            JSONObject ob = setParams(params);

            // request param to API
            OutputStream os = httpsURLConnection.getOutputStream();
            BufferedWriter writer = new BufferedWriter(
                    new OutputStreamWriter(os, Config.CHARSET_UTF_8));
            writer.write(ob.toString());
            writer.flush();
            writer.close();
            os.close();

            InputStream getData = httpsURLConnection.getInputStream();

            JsonReader jsonReader = new JsonReader(new InputStreamReader(getData));

            Gson gson = new GsonBuilder().create();
            // open database
            SQLiteDatabase db = DatabaseManagerCommon.getInstance().openDatabase();
            MaxYearRankModel maxYearRankModel = new MaxYearRankModel(true, db);

            // update progress bar
            publishProgress(Message.MESSAGE_IMPORT_DATA_SCREEN);

            // begin transaction
            db.beginTransaction();
            try {
                jsonReader.beginArray();
                int index = 0;
                while (jsonReader.hasNext()) {
                    index++;
                    maxYearRankModel.insertBulk(
                            (MaxYearRankEntity) gson.fromJson(jsonReader, MaxYearRankEntity.class));

                }
                // submit transaction
                db.setTransactionSuccessful();
                LogManagerCommon.i(Constants.TAG_APPLICATION_NAME,
                        String.format(Message.MESSAGE_TIME_EXECUTE,
                                String.valueOf(
                                        (Calendar.getInstance().getTime().getTime() - cr.getTime())
                                                / 1000),
                                String.valueOf(index)));
                // end transaction
                db.endTransaction();
            } catch (IOException | IllegalStateException e) {
                e.printStackTrace();
                LogManagerCommon.e(Constants.TAG_APPLICATION_NAME, e.toString());

                try {
                    // get result code from API
                    int responseCode = httpsURLConnection.getResponseCode();

                    if (responseCode == HttpsURLConnection.HTTP_OK) {

                        BufferedReader in = new BufferedReader(new InputStreamReader(getData));

                        String line = in.readLine();

                        return (line != null ? line : "");
                    }
                } catch (IOException e1) {
                    e1.printStackTrace();
                }
            } finally {
                // close database
                DatabaseManagerCommon.getInstance().closeDatabase();

            }

        } catch (Exception e) {
            return "Exception: " + e.getMessage();
        }
        return null;
    }

    /**
     * End progress loading
     *
     * @param result from API
     */
    @Override
    protected void onPostExecute(String result) {

        progressDialog.dismiss();
        response.progressFinish(result, 3, null);
    }

    /**
     * Set request params
     *
     * @param params from activity
     * @return Json object from params
     * @see JSONObject
     */
    private JSONObject setParams(String... params) {

        JSONObject jsonObject = new JSONObject();

        try {
            switch (params[0]) {
                case Config.CODE_GET_MAX_YEAR_RANK:
                    jsonObject.put(Constants.COLUMN_SHOP_ID, params[2]);
                    jsonObject.put(Constants.COLUMN_LOGIN_KEY, params[3]);
                    jsonObject.put(Constants.COLUMN_SERVER_NAME, params[4]);
                    break;
            }
            jsonObject.put(Config.API_KEY, Config.API_KEY_VALUE);

        } catch (JSONException e) {
            e.printStackTrace();
            LogManagerCommon.e(Constants.TAG_APPLICATION_NAME, e.toString());
        }
        return jsonObject;
    }

    /**
     * Update progress bar
     *
     * @param values updated
     */
    @Override
    protected void onProgressUpdate(String... values) {
        progressDialog.setMessage(values[0]);
    }
}
