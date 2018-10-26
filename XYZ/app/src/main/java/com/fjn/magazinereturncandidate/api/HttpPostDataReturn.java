package com.fjn.magazinereturncandidate.api;

import android.app.ProgressDialog;
import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.os.AsyncTask;

import com.fjn.magazinereturncandidate.common.constants.Constants;
import com.fjn.magazinereturncandidate.common.constants.Message;
import com.fjn.magazinereturncandidate.common.utils.DatabaseManagerCommon;
import com.fjn.magazinereturncandidate.common.utils.LogManagerCommon;
import com.fjn.magazinereturncandidate.db.models.ReturnMagazineModel;
import com.google.gson.stream.JsonReader;

import java.io.InputStream;
import java.io.InputStreamReader;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.Date;
import java.util.List;

/**
 * HTTP POST request
 *
 * @author cong-pv
 * @since 2018-10-16
 */

public class HttpPostDataReturn extends AsyncTask<String, String, String> {

    /**
     * Http response.
     */
    private HttpResponse response;

    /**
     * Progress dialog.
     */
    private ProgressDialog progressDialog;

    /**
     * Constructor HttpPost.
     */
    public HttpPostDataReturn(Context c) {

        this.response = (HttpResponse) c;
        progressDialog = new ProgressDialog(c);
    }

    /**
     * Set progress dialog loading.
     */
    protected void onPreExecute() {

        progressDialog.setMessage(Message.MESSAGE_IMPORT_DATA_SCREEN);
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

            InputStream getData = this.getClass().getClassLoader().getResourceAsStream("zasshi491.json");

            JsonReader jsonReader = new JsonReader(new InputStreamReader(getData));

            SQLiteDatabase db = DatabaseManagerCommon.getInstance().openDatabase();
            ReturnMagazineModel bookModel = new ReturnMagazineModel(true, db);
            db.beginTransaction();
            try {
                jsonReader.beginArray();
                int index = 0;
                int indexListString = 0;
                publishProgress(Message.MESSAGE_IMPORT_DATA_SCREEN);
                List<String> listValue = new ArrayList<>();
                while (jsonReader.hasNext()) {
                    jsonReader.beginObject();
                    while (jsonReader.hasNext()) {
                        String value = jsonReader.nextName();
                        switch (value) {
                            case Constants.COLUMN_JAN_CD:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_STOCK_COUNT:
                                listValue.add(String.valueOf(jsonReader.nextInt()));
                                break;
                            case Constants.COLUMN_GOODS_NAME:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_WRITER_NAME:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_PUBLISHER_CD:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_PUBLISHER_NAME:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_PRICE:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_FIRST_SUPPLY_DATE:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_LAST_SUPPLY_DATE:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_LAST_SALES_DATE:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_LAST_ORDER_DATE:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_MEDIA_GROUP1_CD:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_MEDIA_GROUP1_NAME:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_MEDIA_GROUP2_CD:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_MEDIA_GROUP2_NAME:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_SALES_DATE:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_TRN_DATE:
                                listValue.add(jsonReader.nextString());
                                break;
                            case Constants.COLUMN_YEAR_RANK:
                                listValue.add(String.valueOf(jsonReader.nextInt()));
                                break;
                            case Constants.COLUMN_TOTAL_SALES:
                                listValue.add(String.valueOf(jsonReader.nextInt()));
                                break;
                            case Constants.COLUMN_TOTAL_SUPPLY:
                                listValue.add(String.valueOf(jsonReader.nextInt()));
                                break;
                            case Constants.COLUMN_TOTAL_RETURN:
                                listValue.add(String.valueOf(jsonReader.nextInt()));
                                break;
                            case Constants.COLUMN_LOCATION_ID:
                                listValue.add(jsonReader.nextString());
                                break;
                            default:
                                jsonReader.skipValue();
                                break;
                        }
                    }
                    jsonReader.endObject();
                    index++;
                    indexListString += Constants.VALUE_COUNT_COLUMN_TABLE_RETURN_MAGAZINE_INSERT; //22: count column table return books insert
                    if (indexListString == 990) { // 990: Maximum record import multi
                        bookModel.insertData(db, indexListString, listValue);
                        listValue.clear();
                        indexListString = 0;
                    }
                }
                //Insert data < 990
                if (indexListString >= Constants.VALUE_COUNT_COLUMN_TABLE_RETURN_MAGAZINE_INSERT) {
                    bookModel.insertData(db, indexListString, listValue);
                }

                db.setTransactionSuccessful();
                LogManagerCommon.i(Constants.TAG_APPLICATION_NAME,
                        String.format(Message.MESSAGE_TIME_EXECUTE,
                                String.valueOf(
                                        (Calendar.getInstance().getTime().getTime() - cr.getTime())
                                                / 1000),
                                String.valueOf(index)));
                jsonReader.endArray();
            } catch (IllegalStateException e) {
                e.printStackTrace();
                LogManagerCommon.e(Constants.TAG_APPLICATION_NAME, e.toString());
            } finally {
                db.endTransaction();
                DatabaseManagerCommon.getInstance().closeDatabase();
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

    @Override
    protected void onProgressUpdate(String... values) {
        progressDialog.setMessage(values[0]);
    }
}
