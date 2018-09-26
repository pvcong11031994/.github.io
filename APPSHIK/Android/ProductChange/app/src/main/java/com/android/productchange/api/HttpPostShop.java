package com.android.productchange.api;

import android.app.ProgressDialog;
import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.os.AsyncTask;

import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.constants.Message;
import com.android.productchange.common.utils.DatabaseManager;
import com.android.productchange.common.utils.LogManager;
import com.android.productchange.db.entity.Books;
import com.android.productchange.db.entity.Returnbooks;
import com.android.productchange.db.models.BookModel;
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
 * <h1>Http Post Shop</h1>
 * <p>
 * Task http Post to API get list books data by user
 *
 * @author tien-lv
 * @since 2017-12-18
 */

public class HttpPostShop extends AsyncTask<String, String, String> {

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
    public HttpPostShop(Context c) {

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
            SQLiteDatabase db = DatabaseManager.getInstance().openDatabase();
            BookModel bookModel = new BookModel(true, db);
            Books books;

            db.beginTransaction();

            try {
                jsonReader.beginArray();
                int index = 0;
                int indexListString = 0;

                // update progress bar
                publishProgress(Message.MESSAGE_IMPORT_DATA_SCREEN);
                List<String> listValue = new ArrayList<>();

                // insert data from json response
                while (jsonReader.hasNext()) {
                    books = gson.fromJson(jsonReader, Books.class);
                    listValue.add(books.getLocation_id());
                    listValue.add(books.getLarge_classifications_id());
                    listValue.add(books.getLarge_classifications_name());
                    listValue.add(books.getName());
                    listValue.add(books.getPublisher_id());
                    listValue.add(books.getPublisher_name());
                    listValue.add(books.getPublish_date());
                    listValue.add(books.getJan_code());
                    listValue.add(String.valueOf(books.getInventory_number()));
                    listValue.add(String.valueOf(books.getOld_catagory_rank()));
                    listValue.add(String.valueOf(books.getNew_catagory_rank()));
                    listValue.add(String.valueOf(books.getFlag_order_now()));
                    listValue.add(String.valueOf(books.getRanking()));
                    index++;
                    indexListString += Constants.VALUE_COUNT_COLUMN_TABLE_BOOK_INSERT; //13: count column table books insert
                    if (indexListString == 988) { // 988: Maximum record import multi (1000 - 1000 % 25)
                        bookModel.insertData(db, indexListString, listValue);
                        listValue.clear();
                        indexListString = 0;
                    }
                }
                //Insert data < 988
                if (indexListString >= Constants.VALUE_COUNT_COLUMN_TABLE_BOOK_INSERT) {
                    bookModel.insertData(db, indexListString, listValue);
                }

                // submit transaction
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
                // end transaction
                db.endTransaction();
                // close database
                DatabaseManager.getInstance().closeDatabase();
                multiThreadCount = 1;
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
                case Config.CODE_GET_LIST_SHOP_BY_USER:
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
