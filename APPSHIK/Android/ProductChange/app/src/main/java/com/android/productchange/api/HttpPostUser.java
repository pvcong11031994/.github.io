package com.android.productchange.api;

import android.app.ProgressDialog;
import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.os.AsyncTask;

import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.constants.Message;
import com.android.productchange.common.utils.DatabaseManager;
import com.android.productchange.common.utils.LogManager;
import com.android.productchange.db.entity.Returnbooks;
import com.android.productchange.db.models.ReturnbookModel;
import com.android.productchange.interfaces.HttpResponse;

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
import java.util.ArrayList;
import java.util.Calendar;
import java.util.Date;
import java.util.List;

import javax.net.ssl.HttpsURLConnection;

/**
 * <h1>Http Post User</h1>
 * <p>
 * Task http Post to API get list return books data
 *
 * @author tien-lv
 * @since 2017-12-18
 */

public class HttpPostUser extends AsyncTask<String, String, String> {

    /**
     * Http response
     */
    private HttpResponse response;

    /**
     * Progress dialog
     */
    ProgressDialog progressDialog;

    /**
     * Multi Thread Count
     */
    private int multiThreadCount;
    /**
     * Constructor HttpPost
     *
     * @param c context
     */
    public HttpPostUser(Context c) {

        this.response = (HttpResponse) c;
        progressDialog = new ProgressDialog(c);
        multiThreadCount = 0;
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
     * @throws IOException from insert error
     * @throws Exception   from request fail
     * @see IOException
     * @see Exception
     */
    @Override
    protected String doInBackground(String... params) {

        try {

            Date cr = Calendar.getInstance().getTime();
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

            InputStream getData = httpsURLConnection.getInputStream();

            JsonReader jsonReader = new JsonReader(new InputStreamReader(getData));

            Gson gson = new GsonBuilder().create();
            SQLiteDatabase db = DatabaseManager.getInstance().openDatabase();
            ReturnbookModel returnbookModel = new ReturnbookModel(true, db);
            Returnbooks returnbooks;
            db.beginTransaction();
            try {
                jsonReader.beginArray();
                int index = 0;
                int indexListString = 0;
                publishProgress(Message.MESSAGE_IMPORT_DATA_SCREEN);
                List<String> listValue = new ArrayList<>();
                while (jsonReader.hasNext()) {
                    returnbooks = gson.fromJson(jsonReader, Returnbooks.class);
                    listValue.add(returnbooks.getJan_cd());
                    listValue.add(String.valueOf(returnbooks.getBqsc_stock_count()));
                    listValue.add(returnbooks.getBqgm_goods_name());
                    listValue.add(returnbooks.getBqgm_writer_name());
                    listValue.add(returnbooks.getBqgm_publisher_cd());
                    listValue.add(returnbooks.getBqgm_publisher_name());
                    listValue.add(returnbooks.getBqgm_price().toString());
                    listValue.add(returnbooks.getBqtse_first_supply_date());
                    listValue.add(returnbooks.getBqtse_last_supply_date());
                    listValue.add(returnbooks.getBqtse_last_sale_date());
                    listValue.add(returnbooks.getBqtse_last_order_date());
                    listValue.add(returnbooks.getBqct_media_group1_cd());
                    listValue.add(returnbooks.getBqct_media_group1_name());
                    listValue.add(returnbooks.getBqct_media_group2_cd());
                    listValue.add(returnbooks.getBqct_media_group2_name());
                    listValue.add(returnbooks.getBqgm_sales_date());
                    listValue.add(returnbooks.getBqio_trn_date());
                    listValue.add(returnbooks.getPercent().toString());
                    listValue.add(returnbooks.getFlag_sales());
                    listValue.add(String.valueOf(returnbooks.getYear_rank()));
                    listValue.add(String.valueOf(returnbooks.getJoubi()));
                    listValue.add(String.valueOf(returnbooks.getSts_total_sales()));
                    listValue.add(String.valueOf(returnbooks.getSts_total_supply()));
                    listValue.add(String.valueOf(returnbooks.getSts_total_return()));
                    listValue.add(returnbooks.getLocation_id());
                    index++;
                    indexListString += Constants.VALUE_COUNT_COLUMN_TABLE_RETURN_BOOK_INSERT; //25: count column table return books insert
                    if (indexListString == 975) { // 975: Maximum record import multi (1000 - 1000 % 25)
                        returnbookModel.insertData(db, indexListString, listValue);
                        listValue.clear();
                        indexListString = 0;
                    }
                }
                //Insert data < 975
                if (indexListString >= Constants.VALUE_COUNT_COLUMN_TABLE_RETURN_BOOK_INSERT) {
                    returnbookModel.insertData(db, indexListString, listValue);
                }

                db.setTransactionSuccessful();
                LogManager.i(Constants.TAG_APPLICATION_NAME,
                        String.format(Message.MESSAGE_TIME_EXECUTE,
                                String.valueOf(
                                        (Calendar.getInstance().getTime().getTime() - cr.getTime())
                                                / 1000),
                                String.valueOf(index)));
                // end transaction
                //db.endTransaction();
            } catch (IOException | IllegalStateException e) {
                e.printStackTrace();
                LogManager.e(Constants.TAG_APPLICATION_NAME, e.toString());

                try {
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
                db.endTransaction();
                DatabaseManager.getInstance().closeDatabase();
                multiThreadCount = 2;
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
            switch (params[0]) {
                case Config.CODE_GET_LIST_BY_USER:
                    jsonObject.put(Constants.COLUMN_SHOP_ID, params[2]);
                    jsonObject.put(Constants.COLUMN_LOGIN_KEY, params[3]);
                    jsonObject.put(Constants.COLUMN_SERVER_NAME, params[4]);
                    break;
            }
            jsonObject.put(Config.API_KEY, Config.API_KEY_VALUE);

        } catch (JSONException e) {
            e.printStackTrace();
            LogManager.e(Constants.TAG_APPLICATION_NAME, e.toString());
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
