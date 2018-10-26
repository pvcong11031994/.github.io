package com.fjn.magazinereturncandidate.activities;

import android.app.ProgressDialog;
import android.content.DialogInterface;
import android.content.Intent;
import android.os.Bundle;
import android.support.v7.app.AlertDialog;
import android.support.v7.app.AppCompatActivity;
import android.view.Gravity;
import android.widget.TextView;
import android.widget.Toast;

import com.fjn.magazinereturncandidate.api.Config;
import com.fjn.magazinereturncandidate.api.HttpPostDataReturn;
import com.fjn.magazinereturncandidate.api.HttpPostMaxYearRank;
import com.fjn.magazinereturncandidate.api.HttpResponse;
import com.fjn.magazinereturncandidate.common.constants.Constants;
import com.fjn.magazinereturncandidate.common.constants.Message;
import com.fjn.magazinereturncandidate.common.helpers.DatabaseHelper;
import com.fjn.magazinereturncandidate.common.utils.DatabaseManagerCommon;
import com.fjn.magazinereturncandidate.common.utils.LogManagerCommon;
import com.fjn.magazinereturncandidate.db.models.MaxYearRankModel;
import com.fjn.magazinereturncandidate.db.models.ReturnMagazineModel;

/**
 * Data retriever class load data from API server
 *
 * @author cong-pv
 * @version 1.0
 * @since 2018-10-15
 */
public class DataLoaderActivity extends AppCompatActivity implements HttpResponse {

    private String TAG = Constants.TAG_APPLICATION_NAME;

    private String userID;
    private String license;
    private String shopID;
    private String loginKey;
    private String serverName;
    private boolean flagLogin;
    /**
     * Progress dialog
     */
    ProgressDialog progressDialog;

    /**
     * Initialize screen layout
     *
     * @param state {@link Bundle }
     */
    @Override
    public void onCreate(Bundle state) {

        super.onCreate(state);

        // Database init
        DatabaseManagerCommon.initializeInstance(new DatabaseHelper(getApplicationContext()));

        Bundle bundle = getIntent().getExtras();
        if (bundle != null) {
            userID = bundle.getString(Constants.COLUMN_USER_ID);
            shopID = bundle.getString(Constants.COLUMN_SHOP_ID);
            loginKey = bundle.getString(Constants.COLUMN_LOGIN_KEY);
            serverName = bundle.getString(Constants.COLUMN_SERVER_NAME);
            flagLogin = bundle.getBoolean(Constants.FLAG_LOGIN);
            license = bundle.getString(Constants.COLUMN_LICENSE);
        }

        // init progress dialog
        progressDialog = new ProgressDialog(this);
        progressDialog.setMessage(Message.MESSAGE_IMPORT_DATA_SCREEN);
        progressDialog.setCancelable(false);
        progressDialog.setCanceledOnTouchOutside(false);

        // Logged in status confirmation
        if (!flagLogin) {
            moveNextActivity();
        } else {

            // Loading return candidate data from API server
            loadDataReturnZasshi();
        }
    }

    /**
     * Response output.
     *
     * @param output           {@link String }
     * @param multiThreadCount int
     * @param fileName         {@link String }
     */
    @Override
    public void progressFinish(String output, int multiThreadCount, String fileName) {

        LogManagerCommon.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.LOADING_DATA_FROM_SERVER_SUCCESS);
        if (output == null) {
            if (multiThreadCount == 2) {
                // Process success
                ReturnMagazineModel returnMagazineModel = new ReturnMagazineModel();

                //Check count record and show dialog
                showDialogRecordDate(returnMagazineModel.countReturnMagazine());

                loadMaxYearRank();
            } else {
                moveNextActivity();
            }
        } else {

            AlertDialog.Builder alertDialogBuilder =
                    new AlertDialog.Builder(this);

            // Error login require
            if (output.contains(Message.CODE_401)) {
                alertDialogBuilder.setMessage(Message.MESSAGE_401 + Message.MESSAGE_RELOAD);
                // Error API URL not found
            } else if (output.contains(Message.CODE_404)) {
                alertDialogBuilder.setMessage(Message.MESSAGE_404 + Message.MESSAGE_RELOAD);
                // Error server
            } else if (output.contains(Message.CODE_500)) {
                alertDialogBuilder.setMessage(Message.MESSAGE_500 + Message.MESSAGE_RELOAD);
                // No response result
            } else {
                alertDialogBuilder.setMessage(
                        Message.MESSAGE_RESULT_EMPTY + Message.MESSAGE_RELOAD);
            }
            alertDialogBuilder.setCancelable(false);

            // Configure alert dialog button
            alertDialogBuilder.setPositiveButton(Message.MESSAGE_SELECT_YES,
                    new DialogInterface.OnClickListener() {
                        @Override
                        public void onClick(DialogInterface dialog, int which) {
                            dialog.dismiss();

                            // reload data again form server
                            loadDataReturnZasshi();
                        }
                    });
            alertDialogBuilder.setNegativeButton(Message.MESSAGE_SELECT_NO,
                    new DialogInterface.OnClickListener() {
                        @Override
                        public void onClick(DialogInterface dialog, int which) {
                            dialog.dismiss();

                            // Process logout
                            // clear table
                            DatabaseManagerCommon.initializeInstance(
                                    new DatabaseHelper(getApplicationContext()));
                            DatabaseHelper ds = new DatabaseHelper(DataLoaderActivity.this);
                            ds.clearTables();

                            // move to login screen
                            Intent intent = new Intent(DataLoaderActivity.this,
                                    LoginActivity.class);
                            startActivity(intent);
                            finish();
                        }
                    });
            AlertDialog alert = alertDialogBuilder.show();
            TextView messageText = (TextView) alert.findViewById(android.R.id.message);
            assert messageText != null;
            messageText.setGravity(Gravity.CENTER);
        }

    }

    /**
     * Loading data from API server
     */

    public void loadDataReturnZasshi() {

        //Import data JSON to database
        ReturnMagazineModel returnMagazineModel = new ReturnMagazineModel();

        // Get data from API if table is empty
        if (!returnMagazineModel.checkData()) {
            LogManagerCommon.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.LOADING_DATA_FROM_SERVER);
            new HttpPostDataReturn(this).execute();
        } else {
            loadMaxYearRank();
        }

    }

    /**
     * Load max year rank
     */
    public void loadMaxYearRank() {

        MaxYearRankModel maxYearRankModel = new MaxYearRankModel();

        // Get data from API if table is empty
        if (!maxYearRankModel.checkData()) {
            LogManagerCommon.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.LOADING_DATA_FROM_SERVER);
            String[] params =
                    new String[]{Config.CODE_GET_MAX_YEAR_RANK, Config.API_GET_MAX_YEAR_RANK, shopID, loginKey, serverName};
            new HttpPostMaxYearRank(this).execute(params);
        } else {
            moveNextActivity();
        }
    }


    /**
     * Transition to unlock screen if not logged in.<br>
     * Transition to scanner when data is loaded
     */
    public void moveNextActivity() {

        Intent intent;
        if (!flagLogin) {
            intent = new Intent(this, UnlockScreenActivity.class);
        } else {
            intent = new Intent(this, SdmScannerActivity.class);
        }
        Bundle bundle = new Bundle();
        bundle.putString(Constants.COLUMN_USER_ID, userID);
        bundle.putString(Constants.COLUMN_SHOP_ID, shopID);
        bundle.putString(Constants.COLUMN_SERVER_NAME, serverName);
        bundle.putString(Constants.COLUMN_LICENSE, license);
        intent.putExtras(bundle);
        startActivity(intent);
        finish();
    }

    /**
     * Back event handler
     */
    @Override
    public void onBackPressed() {
        super.onBackPressed();

        finishAffinity();
    }

    private void showDialogRecordDate(int countRecordData) {
        if (countRecordData > 0) {
            Toast.makeText(this, String.format(Message.MESSAGE_LOADING_DATA_NUMBER,
                    String.valueOf(countRecordData)), Toast.LENGTH_LONG).show();
            LogManagerCommon.i(TAG, String.format(Message.MESSAGE_LOADING_DATA_NUMBER,
                    String.valueOf(countRecordData)));
        }
    }

}