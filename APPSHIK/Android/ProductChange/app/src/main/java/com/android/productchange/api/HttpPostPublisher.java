package com.android.productchange.api;

import android.annotation.SuppressLint;
import android.app.ProgressDialog;
import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.os.AsyncTask;

import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.constants.Message;
import com.android.productchange.common.utils.DatabaseManager;
import com.android.productchange.db.entity.Books;
import com.android.productchange.db.entity.CLP;
import com.android.productchange.db.models.BookModel;
import com.android.productchange.db.models.CLPModel;
import com.android.productchange.db.models.PeriodbookModel;
import com.android.productchange.db.models.RegularbookModel;
import com.android.productchange.interfaces.HttpResponse;

import java.io.IOException;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.Date;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;

/**
 * <h1>Http Post Publisher</h1>
 * <p>
 * Task http Post to API get list Publisher data
 *
 * @author cong-pv
 * @since 2018-09-06
 */

public class HttpPostPublisher extends AsyncTask<String, String, String> {

    /**
     * Http response
     */
    private HttpResponse response;

    /**
     * Progress dialog
     */
    ProgressDialog progressDialog;

    /**
     * Multi thread count
     */
    private int multiThreadCount;
    private CLPModel clpModel = new CLPModel();
    private BookModel bookModel = new BookModel();
    private PeriodbookModel periodbookModel = new PeriodbookModel();
    private RegularbookModel regularbookModel = new RegularbookModel();
    /**
     * Create date
     */
    private String createDate;

    /**
     * Constructor HttpPost
     *
     * @param c context
     */
    public HttpPostPublisher(Context c) {

        this.response = (HttpResponse) c;
        progressDialog = new ProgressDialog(c);
        multiThreadCount = 8;
    }

    /**
     * Set progress dialog loading
     */
    protected void onPreExecute() {

        progressDialog.setMessage(Message.MESSAGE_IMPORT_DATA_SCREEN);
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

        createDate = params[0];
        List<CLP> rlistPublisher;

        // check table Publisher is data exist
        if (!clpModel.checkData(Config.TYPE_PUBLISHER)) {
            rlistPublisher = bookModel.getInfo(Config.TYPE_PUBLISHER);
            HashMap<String, String> map = new HashMap<>();
            for (int i = 0; i < rlistPublisher.size(); i++) {
                if (!map.containsKey(rlistPublisher.get(i).getId())) {
                    map.put(rlistPublisher.get(i).getId(), rlistPublisher.get(i).getName());
                }
            }

            rlistPublisher = periodbookModel.getInfo(Config.TYPE_PUBLISHER);
            for (int i = 0; i < rlistPublisher.size(); i++) {
                if (!map.containsKey(rlistPublisher.get(i).getId())) {
                    map.put(rlistPublisher.get(i).getId(), rlistPublisher.get(i).getName());
                }

            }

            rlistPublisher = regularbookModel.getInfo(Config.TYPE_PUBLISHER);
            for (int i = 0; i < rlistPublisher.size(); i++) {
                if (!map.containsKey(rlistPublisher.get(i).getId())) {
                    map.put(rlistPublisher.get(i).getId(), rlistPublisher.get(i).getName());
                }

            }
            // open database
            SQLiteDatabase db = DatabaseManager.getInstance().openDatabase();
            CLPModel clpPublisherModel = new CLPModel(true, db);

            // begin transaction
            db.beginTransaction();
            int indexListString = 0;
            List<String> listValue = new ArrayList<>();
            // insert data into table Publisher
            for (String key : map.keySet()) {
                if (!((map.get(key)).trim()).isEmpty()) {
                    listValue.add(key);
                    listValue.add(map.get(key));
                    indexListString += Constants.VALUE_COUNT_COLUMN_TABLE_PUBLISHER_INSERT; //2: count column table period books insert
                    if (indexListString == 998) { // 998: Maximum record import multi (1000 - 1000 % 25)
                        clpPublisherModel.insertDataPublisher(db, indexListString, listValue);
                        listValue.clear();
                        indexListString = 0;
                    }
                    //clpPublisherModel.insertBulkPublisher(key, map.get(key));
                }
            }
            //Insert data < 998
            if (indexListString >= Constants.VALUE_COUNT_COLUMN_TABLE_PUBLISHER_INSERT) {
                clpPublisherModel.insertDataPublisher(db, indexListString, listValue);
            }
            // submit transaction
            db.setTransactionSuccessful();
            // end transaction
            db.endTransaction();
            // close database
            DatabaseManager.getInstance().closeDatabase();

        }

        loadTableView();


        return null;
    }

    /**
     * Load data publisher from table Books and ReturnBooks
     */
    private void loadTableView() {

        List<Books> listTableView;

        @SuppressLint("SimpleDateFormat") SimpleDateFormat sdf = new SimpleDateFormat(
                Constants.DATE_FORMAT_STRING);
        Date calDateFrom = null, calDateTo = null;
        try {
            calDateFrom = sdf.parse(createDate);
            calDateTo = sdf.parse(createDate);
        } catch (ParseException e) {
            e.printStackTrace();
        }

        Calendar dateFromDefault = Calendar.getInstance(), dateToDefault = Calendar.getInstance();
        dateFromDefault.setTime(calDateFrom);
        dateToDefault.setTime(calDateTo);
        dateFromDefault.add(Calendar.DATE, Constants.DATE_FROM);
        dateToDefault.add(Calendar.DATE, Constants.DATE_TO);

        @SuppressLint("SimpleDateFormat") SimpleDateFormat df = new SimpleDateFormat(
                Constants.DATE_FORMAT_STRING);
        String formattedDateFrom = df.format(dateFromDefault.getTime());
        String formattedDateTo = df.format(dateToDefault.getTime());

        LinkedHashMap<String, String> mapOrder = new LinkedHashMap<>();

        listTableView = periodbookModel.getListBookInfo(Constants.ID_ROW_ALL,
                Config.TYPE_CLASSIFY, 0, formattedDateFrom, formattedDateTo, mapOrder);

        // open database
        SQLiteDatabase db = DatabaseManager.getInstance().openDatabase();
        PeriodbookModel mPeriodbookModel = new PeriodbookModel(true, db);
        // begin transaction
        db.beginTransaction();
        // insert data into table view Period Books
        for (int i = 0; i < listTableView.size(); i++) {
            mPeriodbookModel.insertViewBulk(listTableView.get(i));
        }

        // submit transaction
        db.setTransactionSuccessful();
        // end transaction
        db.endTransaction();
        // close database
        DatabaseManager.getInstance().closeDatabase();

        listTableView.clear();


        listTableView = regularbookModel.getListBookInfo(Constants.ID_ROW_ALL,
                Config.TYPE_CLASSIFY, 0, mapOrder);
        // open database
        db = DatabaseManager.getInstance().openDatabase();

        RegularbookModel mRegularbookModel = new RegularbookModel(true, db);
        // begin transaction
        db.beginTransaction();
        // insert data into table view Period Books
        for (int i = 0; i < listTableView.size(); i++) {
            mRegularbookModel.insertViewBulk(listTableView.get(i));
        }

        // submit transaction
        db.setTransactionSuccessful();
        // end transaction
        db.endTransaction();
        // close database
        DatabaseManager.getInstance().closeDatabase();

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
     * Update progress bar
     *
     * @param values updated
     */
    @Override
    protected void onProgressUpdate(String... values) {
        progressDialog.setMessage(values[0]);
    }
}
