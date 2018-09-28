package com.android.returncandidate.api;

import android.app.*;
import android.content.*;
import android.database.sqlite.*;
import android.os.*;

import com.android.returncandidate.common.constants.*;
import com.android.returncandidate.common.constants.Message;
import com.android.returncandidate.common.utils.*;
import com.android.returncandidate.db.entity.*;
import com.android.returncandidate.db.models.*;
import com.google.gson.*;
import com.google.gson.stream.*;

import org.json.*;

import java.io.*;
import java.net.*;
import java.util.*;

import javax.net.ssl.*;

/**
 * HTTP POST request
 *
 * @author tien-lv
 * @since 2017-12-18
 * Created by tien-lv on 2017/12/18.
 */

public class HttpPostUser extends AsyncTask<String, String, String> {

    /**
     * Http response.
     */
    private HttpResponse response;

    /**
     * Progress dialog.
     */
    ProgressDialog progressDialog;

    /**
     * Constructor HttpPost.
     */
    public HttpPostUser(Context c) {

        this.response = (HttpResponse) c;
        progressDialog = new ProgressDialog(c);
    }

    /**
     * Set progress dialog loading.
     */
    protected void onPreExecute() {

        progressDialog.setMessage(Message.MESSAGE_DOWNLOAD_DATA_SCREEN);
        progressDialog.setIndeterminate(false);
        progressDialog.setCancelable(false);
        progressDialog.setCanceledOnTouchOutside(false);
        progressDialog.show();
    }

    /**
     * Send request and get response to services API.
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
            BookModel bookModel = new BookModel(true, db);
            Books books;
            db.beginTransaction();
            try {
                jsonReader.beginArray();
                int index = 0;
                int indexListString = 0;
                publishProgress(Message.MESSAGE_IMPORT_DATA_SCREEN);
                List<String> listValue = new ArrayList<>();
                while (jsonReader.hasNext()) {
                    books = gson.fromJson(jsonReader, Books.class);
                    listValue.add(books.getJan_cd());
                    listValue.add(String.valueOf(books.getBqsc_stock_count()));
                    listValue.add(books.getBqgm_goods_name());
                    listValue.add(books.getBqgm_writer_name());
                    listValue.add(books.getBqgm_publisher_cd());
                    listValue.add(books.getBqgm_publisher_name());
                    listValue.add(books.getBqgm_price().toString());
                    listValue.add(books.getBqtse_first_supply_date());
                    listValue.add(books.getBqtse_last_supply_date());
                    listValue.add(books.getBqtse_last_sale_date());
                    listValue.add(books.getBqtse_last_order_date());
                    listValue.add(books.getBqct_media_group1_cd());
                    listValue.add(books.getBqct_media_group1_name());
                    listValue.add(books.getBqct_media_group2_cd());
                    listValue.add(books.getBqct_media_group2_name());
                    listValue.add(books.getBqgm_sales_date());
                    listValue.add(books.getBqio_trn_date());
                    listValue.add(books.getPercent().toString());
                    listValue.add(books.getFlag_sales());
                    listValue.add(String.valueOf(books.getYear_rank()));
                    listValue.add(String.valueOf(books.getJoubi()));
                    listValue.add(String.valueOf(books.getSts_total_sales()));
                    listValue.add(String.valueOf(books.getSts_total_supply()));
                    listValue.add(String.valueOf(books.getSts_total_return()));
                    listValue.add(books.getLocation_id());
                    index++;
                    indexListString += Constants.VALUE_COUNT_COLUMN_TABLE_RETURN_BOOK_INSERT; //25: count column table return books insert
                    if (indexListString == 975) { // 975: Maximum record import multi
                        bookModel.insertData(db, indexListString, listValue);
                        listValue.clear();
                        indexListString = 0;
                    }
                }
                //Insert data < 975
                if (indexListString >= Constants.VALUE_COUNT_COLUMN_TABLE_RETURN_BOOK_INSERT) {
                    bookModel.insertData(db, indexListString, listValue);
                }

                db.setTransactionSuccessful();
                LogManager.i(Constants.TAG_APPLICATION_NAME,
                        String.format(Message.MESSAGE_TIME_EXECUTE,
                                String.valueOf(
                                        (Calendar.getInstance().getTime().getTime() - cr.getTime())
                                                / 1000),
                                String.valueOf(index)));
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
            }

        } catch (Exception e) {
            return "Exception: " + e.getMessage();
        }
        return null;
    }

    /**
     * End progress loading.
     */
    @Override
    protected void onPostExecute(String result) {

        progressDialog.dismiss();
        //move get data classify when complete get data return books
        response.progressFinish(result, 2, null);
    }

    /**
     * Set request params.
     *
     * @return JSONObject json form response
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

    @Override
    protected void onProgressUpdate(String... values) {
        progressDialog.setMessage(values[0]);
    }
}
