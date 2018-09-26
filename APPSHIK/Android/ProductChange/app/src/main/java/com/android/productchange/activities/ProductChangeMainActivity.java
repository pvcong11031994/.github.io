package com.android.productchange.activities;

import android.app.ProgressDialog;
import android.content.DialogInterface;
import android.content.Intent;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.view.Gravity;
import android.widget.TextView;
import android.widget.Toast;

import com.android.productchange.api.Config;
import com.android.productchange.api.HttpPostClassify;
import com.android.productchange.api.HttpPostMaxYearRank;
import com.android.productchange.api.HttpPostPeriod;
import com.android.productchange.api.HttpPostPublisher;
import com.android.productchange.api.HttpPostPublisherReturnBooks;
import com.android.productchange.api.HttpPostRegular;
import com.android.productchange.api.HttpPostShop;
import com.android.productchange.api.HttpPostUser;
import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.constants.Message;
import com.android.productchange.common.helpers.DatabaseHelper;
import com.android.productchange.common.utils.DatabaseManager;
import com.android.productchange.common.utils.LogManager;
import com.android.productchange.db.models.BookModel;
import com.android.productchange.db.models.CLPModel;
import com.android.productchange.db.models.MaxYearRankModel;
import com.android.productchange.db.models.PeriodbookModel;
import com.android.productchange.db.models.PublisherModel;
import com.android.productchange.db.models.PublisherReturnBooksModel;
import com.android.productchange.db.models.RegularbookModel;
import com.android.productchange.db.models.ReturnbookModel;
import com.android.productchange.interfaces.HttpResponse;


/**
 * <h1>Product Change Main Activity</h1>
 * <p>
 * This activity check data has in tables Books, ReturnBooks, Classify, Publisher
 * if tables are null, activity call to server for insert data and move to next activity
 *
 * @author tien-lv
 * @since 2018-02-08
 */
public class ProductChangeMainActivity extends AppCompatActivity implements HttpResponse {

    /**
     * TAG
     */
    private String TAG = Constants.TAG_APPLICATION_NAME;

    /**
     * User id
     */
    private String userID;

    /**
     * Shop id
     */
    private String shopID;

    /**
     * Login key
     */
    private String loginKey;

    /**
     * Server name
     */
    private String serverName;

    /**
     * Create date
     */
    private String createDate;

    /**
     * Check login first
     */
    private boolean flagLogin;

    /**
     * Progress dialog.
     */
    ProgressDialog progressDialog;

    /**
     * Init on Create Activity
     *
     * @param state bundle data to request API server
     */
    @Override
    public void onCreate(Bundle state) {

        super.onCreate(state);

        // init database
        DatabaseManager.initializeInstance(new DatabaseHelper(getApplicationContext()));

        // get data from bundle
        Bundle bundle = getIntent().getExtras();
        if (bundle != null) {
            userID = bundle.getString(Constants.COLUMN_USER_ID);
            shopID = bundle.getString(Constants.COLUMN_SHOP_ID);
            loginKey = bundle.getString(Constants.COLUMN_LOGIN_KEY);
            serverName = bundle.getString(Constants.COLUMN_SERVER_NAME);
            flagLogin = bundle.getBoolean(Constants.FLAG_LOGIN);
            createDate = bundle.getString(Constants.COLUMN_CREATE_DATE);
        }

        // init progress dialog
        progressDialog = new ProgressDialog(this);
        progressDialog.setMessage(Message.MESSAGE_IMPORT_DATA_SCREEN);
        progressDialog.setCancelable(false);
        progressDialog.setCanceledOnTouchOutside(false);

        // check is login
        if (!flagLogin) {

            moveNextActivity();
        } else {
            // Loading data books from API server
            loadBooksData();

        }
    }

    /**
     * Response result from server
     *
     * @param output           string response
     * @param multiThreadCount multithread count
     * @param fileName         file name
     */
    @Override
    public void progressFinish(String output, int multiThreadCount, String fileName) {

        LogManager.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.LOADING_DATA_FROM_SERVER_SUCCESS);
        if (output == null) {
            if (multiThreadCount == 1) {
                // process success
                BookModel bookModel = new BookModel();

                //Check count record and show dialog
                showDialogRecordDate(bookModel.countBooks());

                loadReturnbookData();

            } else if (multiThreadCount == 2) {
                // process success
                ReturnbookModel returnbookModel = new ReturnbookModel();

                //Check count record and show dialog
                showDialogRecordDate(returnbookModel.countBooks());

                loadPeriodbookData();
            } else if (multiThreadCount == 3) {
                // process success
                PeriodbookModel periodbookModel = new PeriodbookModel();

                //Check count record and show dialog
                showDialogRecordDate(periodbookModel.countBooks());

                loadRegularbookData();
            } else if (multiThreadCount == 4) {

                // process success
                RegularbookModel regularbookModel = new RegularbookModel();

                //Check count record and show dialog
                showDialogRecordDate(regularbookModel.countBooks());

                loadClassifyData();
            } else if (multiThreadCount == 5) {

                // process success
                CLPModel clpModel = new CLPModel();

                //Check count record and show dialog
                showDialogRecordDate(clpModel.countBooks());

                loadPublisherReturnBooksData();
            } else if (multiThreadCount == 6) {

                // process success
                PublisherReturnBooksModel publisherReturnBooksModel = new PublisherReturnBooksModel();

                //Check count record and show dialog
                showDialogRecordDate(publisherReturnBooksModel.countDataTablePublisherReturn());

                loadMaxYearRank();
            } else if (multiThreadCount == 7) {
                // process success
                MaxYearRankModel maxYearRankModel = new MaxYearRankModel();

                //Check count record and show dialog
                showDialogRecordDate(maxYearRankModel.countMaxYearRank());

                loadPublisher();
            } else if (multiThreadCount == 8) {
                // Process success
                CLPModel publisherModel = new CLPModel();

                //Check count record and show dialog
                showDialogRecordDate(publisherModel.countDataTablePublisher());

                moveNextActivity();
            } else {
                moveNextActivity();
            }
        } else {

            android.support.v7.app.AlertDialog.Builder alertDialogBuilder =
                    new android.support.v7.app.AlertDialog.Builder(this);

            // Error login require
            if (output.contains(Message.CODE_401)) {
                alertDialogBuilder.setMessage(Message.MESSAGE_401 + Message.MESSAGE_RELOAD);
                // Error API URL not found
            } else if (output.contains(Message.CODE_404)) {
                alertDialogBuilder.setMessage(Message.MESSAGE_404 + Message.MESSAGE_RELOAD);
                // Error server
            } else if (output.contains(Message.CODE_500)) {
                alertDialogBuilder.setMessage(Message.MESSAGE_500 + Message.MESSAGE_RELOAD);
                // no response result
            } else {
//                alertDialogBuilder.setMessage(
//                       Message.MESSAGE_RESULT_EMPTY + Message.MESSAGE_RELOAD);
                return;
            }
            alertDialogBuilder.setCancelable(false);
            // show dialog which message error form server
            alertDialogBuilder.setPositiveButton(Message.MESSAGE_SELECT_YES,
                    new DialogInterface.OnClickListener() {
                        @Override
                        public void onClick(DialogInterface dialog, int which) {
                            dialog.dismiss();

                            // reload data again form server
                            loadBooksData();
                        }
                    });
            alertDialogBuilder.setNegativeButton(Message.MESSAGE_SELECT_NO,
                    new DialogInterface.OnClickListener() {
                        @Override
                        public void onClick(DialogInterface dialog, int which) {
                            dialog.dismiss();

                            // Process logout
                            // clear table
                            DatabaseManager.initializeInstance(
                                    new DatabaseHelper(getApplicationContext()));
                            DatabaseHelper ds = new DatabaseHelper(ProductChangeMainActivity.this);
                            ds.clearTables();

                            // move to login screen
                            Intent intent = new Intent(ProductChangeMainActivity.this,
                                    LoginActivity.class);
                            startActivity(intent);
                            finish();
                        }
                    });
            android.support.v7.app.AlertDialog alert = alertDialogBuilder.show();
            TextView messageText = (TextView) alert.findViewById(android.R.id.message);
            messageText.setGravity(Gravity.CENTER);
        }
    }

    private void showDialogRecordDate(int countRecordData) {
        if (countRecordData > 0) {
            Toast.makeText(this, String.format(Message.MESSAGE_LOADING_DATA_NUMBER,
                    String.valueOf(countRecordData)), Toast.LENGTH_LONG).show();
            LogManager.i(TAG, String.format(Message.MESSAGE_LOADING_DATA_NUMBER,
                    String.valueOf(countRecordData)));
        }
    }

    /**
     * Load max year rank
     */
    public void loadMaxYearRank() {

        MaxYearRankModel maxYearRankModel = new MaxYearRankModel();

        // Get data from API if table is empty
        if (!maxYearRankModel.checkData()) {
            LogManager.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.LOADING_DATA_FROM_SERVER);
            String[] params =
                    new String[]{Config.CODE_GET_MAX_YEAR_RANK, Config.API_GET_MAX_YEAR_RANK, shopID, loginKey, serverName};
            new HttpPostMaxYearRank(this).execute(params);
        } else {
            loadPublisher();
        }
    }

    /**
     * Loading data from API server
     */
    public void loadBooksData() {

        BookModel bookModel = new BookModel();

        // Get data from API if table is empty
        if (!bookModel.checkData()) {
            LogManager.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.LOADING_DATA_FROM_SERVER);
            String[] params =
                    new String[]{Config.CODE_GET_LIST_SHOP_BY_USER,
                            Config.API_GET_LIST_SHOP_BY_USER,
                            shopID, loginKey, serverName};
            new HttpPostShop(this).execute(params);
        } else {
            loadReturnbookData();
        }
    }

    /**
     * Loading data Returnbook from API server
     */
    public void loadReturnbookData() {

        ReturnbookModel returnbookModel = new ReturnbookModel();

        // Get data from API if table is empty
        if (!returnbookModel.checkData()) {
            LogManager.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.LOADING_DATA_FROM_SERVER);
            String[] params =
                    new String[]{Config.CODE_GET_LIST_BY_USER, Config.API_GET_LIST_BY_USER,
                            shopID, loginKey, serverName};
            new HttpPostUser(this).execute(params);
        } else {
            loadPeriodbookData();
        }
    }

    /**
     * Loading data Returnbook from API server
     */
    public void loadPeriodbookData() {

        PeriodbookModel periodbookModel = new PeriodbookModel();

        // Get data from API if table is empty
        if (!periodbookModel.checkData()) {
            LogManager.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.LOADING_DATA_FROM_SERVER);
            String[] params =
                    new String[]{Config.CODE_GET_LIST_PERIOD, Config.API_GET_LIST_PERIOD,
                            shopID, loginKey, serverName};
            new HttpPostPeriod(this).execute(params);
        } else {
            loadRegularbookData();
        }
    }

    /**
     * Loading data Returnbook from API server
     */
    public void loadRegularbookData() {

        RegularbookModel regularbookModel = new RegularbookModel();

        // Get data from API if table is empty
        if (!regularbookModel.checkData()) {
            LogManager.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.LOADING_DATA_FROM_SERVER);
            String[] params =
                    new String[]{Config.CODE_GET_LIST_REGULAR, Config.API_GET_LIST_REGULAR,
                            shopID, loginKey, serverName};
            new HttpPostRegular(this).execute(params);
        } else {
            loadClassifyData();
        }
    }

    /**
     * Loading data Classify from API server
     */
    public void loadClassifyData() {

        CLPModel clpModel = new CLPModel();

        // Get data from API if table is empty
        if (!clpModel.checkData(Config.TYPE_CLASSIFY)) {
            LogManager.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.LOADING_DATA_FROM_SERVER);
            String[] params =
                    new String[]{Config.CODE_GET_LIST_CLASSIFY, Config.API_GET_LIST_CLASSIFY,
                            shopID, loginKey, serverName};
            new HttpPostClassify(this).execute(params);
        } else {
            loadPublisherReturnBooksData();
        }
    }

    /**
     * Loading data publisher return books
     */
    public void loadPublisherReturnBooksData() {

        PublisherReturnBooksModel publisherReturnBooksModel = new PublisherReturnBooksModel();

        // Get data from API if table is empty
        if (!publisherReturnBooksModel.checkData()) {
            LogManager.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.MESSAGE_IMPORT_DATA_SCREEN);
            new HttpPostPublisherReturnBooks(this).execute();
        } else {
            loadMaxYearRank();
        }
    }

    /**
     * Load data publisher from table Books and ReturnBooks
     */
    private void loadPublisher() {

        PublisherModel publisherModel = new PublisherModel();

        // Get data from API if table is empty
        if (!publisherModel.checkData()) {
            LogManager.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.MESSAGE_IMPORT_DATA_SCREEN);
            String[] params =
                    new String[]{createDate};
            new HttpPostPublisher(this).execute(params);
        } else {
            moveNextActivity();
        }
    }

    /**
     * Move to Classify Screen with user is admin
     * Move to Barcode Screen with user is employees
     */

    public void moveNextActivity() {

        Intent intent;
        if (!flagLogin) {
            // move to process scanner barcode
            intent = new Intent(this, UnlockScreenActivity.class);
        } else {
            // move to process scanner barcode
            intent = new Intent(this, ProductChangeListActivity.class);
        }
        Bundle bundle = new Bundle();
        bundle.putString(Constants.COLUMN_USER_ID, userID);
        bundle.putString(Constants.COLUMN_SHOP_ID, shopID);
        bundle.putString(Constants.COLUMN_SERVER_NAME, serverName);
        bundle.putString(Constants.COLUMN_CREATE_DATE, createDate);
        intent.putExtras(bundle);
        startActivity(intent);
        finish();
    }

    /**
     * even when click back
     */
    @Override
    public void onBackPressed() {
        super.onBackPressed();
        finishAffinity();
    }

}
