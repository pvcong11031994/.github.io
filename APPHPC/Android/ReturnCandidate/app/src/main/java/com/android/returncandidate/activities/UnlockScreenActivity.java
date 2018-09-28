package com.android.returncandidate.activities;

import android.annotation.*;
import android.app.*;
import android.content.*;
import android.net.*;
import android.os.*;
import android.support.v7.app.*;
import android.text.*;
import android.util.Log;
import android.view.*;
import android.widget.*;

import com.android.returncandidate.*;
import com.android.returncandidate.api.*;
import com.android.returncandidate.common.constants.*;
import com.android.returncandidate.common.constants.Message;
import com.android.returncandidate.common.helpers.*;
import com.android.returncandidate.common.utils.*;
import com.android.returncandidate.db.models.*;

import java.io.*;
import java.text.*;
import java.util.*;

/**
 * Unlock Screen
 */

public class UnlockScreenActivity extends AppCompatActivity implements View.OnClickListener,
        HttpResponse {

    /**
     * TAG
     */
    private String TAG = Constants.TAG_APPLICATION_NAME;

    /**
     * User id
     */
    private String userID;

    /**
     * User license
     */
    private String license;

    /**
     * List scan
     */
    private HashMap<String, LinkedList<String[]>> hashMapArrBooks;

    /**
     * Shop id
     */
    private String shopID;

    /**
     * Server name
     */
    private String serverName;

    /**
     * Edit text password.
     */
    private EditText edtPassword;

    /**
     * Count wrong password input
     */
    private int count = 0;

    /**
     * Count file send response
     */
    private int countFile = 0;

    /**
     * List file
     */
    private File[] files;

    /**
     * Progress dialog.
     */
    private ProgressDialog progress;

    /**
     * Init
     */
    @Override
    protected void onCreate(Bundle savedInstanceState) {

        // print log begin process
        LogManager.i(TAG, Message.TAG_UNLOCK_ACTIVITY + Message.MESSAGE_ACTIVITY_START);
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_unlock_screen);

        Bundle bundle = getIntent().getExtras();
        if (bundle != null) {
            userID = bundle.getString(Constants.COLUMN_USER_ID);
            shopID = bundle.getString(Constants.COLUMN_SHOP_ID);
            serverName = bundle.getString(Constants.COLUMN_SERVER_NAME);
            license = bundle.getString(Constants.COLUMN_LICENSE);
            hashMapArrBooks = (HashMap<String, LinkedList<String[]>>) bundle.getSerializable(Constants.COLUMN_INFOR_LIST_SCAN);
        }

        // init process loading screen
        progress = new ProgressDialog(this);
        progress.setMessage(Message.MESSAGE_UPLPOAD_LOG_SCREEN);
        progress.setCancelable(false);

        TextView txvUserID = (TextView) findViewById(R.id.txv_userid);
        edtPassword = (EditText) findViewById(R.id.edt_password);
        Button btnUnlock = (Button) findViewById(R.id.btn_unlock);
        Button btnLogout = (Button) findViewById(R.id.btn_logout);

        txvUserID.setText(userID);

        btnUnlock.setOnClickListener(this);
        btnLogout.setOnClickListener(this);
    }

    /**
     * Onclick event.
     */
    @Override
    public void onClick(View v) {

        switch (v.getId()) {
            case R.id.btn_unlock:
                unLock();
                break;
            case R.id.btn_logout:
                showAlertDialog(true);
                break;
        }
    }

    /**
     * Unlock function
     */
    public void unLock() {

        UserModel userModel = new UserModel();

        // Validate input
        if (TextUtils.isEmpty(edtPassword.getText().toString())) {
            edtPassword.setError(
                    String.format(Message.MESSAGE_CHECK_INPUT_EMPTY, getString(R.string.password)));
            edtPassword.requestFocus();
        } else {
            if (!userModel.checkDataIsExist(edtPassword.getText().toString())) {
                android.support.v7.app.AlertDialog.Builder alertDialogBuilder =
                        new android.support.v7.app.AlertDialog.Builder(this);
                alertDialogBuilder.setMessage(Message.MESSAGE_PASSWORD_ERR);
                alertDialogBuilder.setCancelable(false);
                LogManager.e(TAG, Message.TAG_UNLOCK_ACTIVITY + Message.MESSAGE_PASSWORD_ERR);
                alertDialogBuilder.setNegativeButton("OK", new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        dialog.dismiss();
                        count++;
                        if (count >= 3) {
                            progress.show();
                            compressFile();
                            sendFileLog();
                        }
                    }
                });
                android.support.v7.app.AlertDialog alert = alertDialogBuilder.show();
                TextView messageText = (TextView) alert.findViewById(android.R.id.message);
                assert messageText != null;
                messageText.setGravity(Gravity.CENTER);
            } else {

                LogManager.i(TAG, Message.TAG_UNLOCK_ACTIVITY + Message.MESSAGE_ACTIVITY_END);

                LogManager.i(TAG, String.format(Message.MESSAGE_ACTIVITY_MOVE,
                        Message.UNLOCK_ACTIVITY_NAME, Message.SCANNER_ACTIVITY_NAME));

                // move to process scanner barcode
                Intent intent = new Intent(this, SdmScannerActivity.class);
                Bundle bundle = new Bundle();
                Log.d("LINCONGPV", license);
                bundle.putString(Constants.COLUMN_USER_ID, userID);
                bundle.putString(Constants.COLUMN_SHOP_ID, shopID);
                bundle.putString(Constants.COLUMN_SERVER_NAME, serverName);
                bundle.putString(Constants.COLUMN_LICENSE, license);
                bundle.putSerializable(Constants.COLUMN_INFOR_LIST_SCAN, hashMapArrBooks);
                intent.putExtras(bundle);
                startActivity(intent);
                finish();
            }
        }
    }

    /**
     * Show dialog warning logout
     */
    public void showAlertDialog(boolean isLogout) {

        progress.show();
        android.support.v7.app.AlertDialog.Builder dialog =
                new android.support.v7.app.AlertDialog.Builder(this);
        dialog.setCancelable(false);
        if (isLogout) {
            dialog
                    .setMessage(getString(R.string.logout_msg))
                    .setPositiveButton(getString(R.string.logout_yes),
                            new DialogInterface.OnClickListener() {
                                @Override
                                public void onClick(DialogInterface dialog, int which) {

                                    LogManager.i(TAG,
                                            Message.TAG_UNLOCK_ACTIVITY + Message.MESSAGE_LOGOUT);
                                    // print log end process
                                    LogManager.i(TAG, Message.TAG_UNLOCK_ACTIVITY
                                            + Message.MESSAGE_ACTIVITY_END);
                                    // print log move screen
                                    LogManager.i(TAG,
                                            String.format(Message.MESSAGE_ACTIVITY_MOVE,
                                                    Message.UNLOCK_ACTIVITY_NAME,
                                                    Message.LOGIN_ACTIVITY_NAME));
                                    compressFile();
                                    sendFileLog();
                                }
                            })
                    .setNegativeButton(getString(R.string.logout_no),
                            new DialogInterface.OnClickListener() {
                                @Override
                                public void onClick(DialogInterface dialog, int which) {
                                    dialog.dismiss();
                                    progress.dismiss();
                                }
                            });
        } else {

            dialog
                    .setMessage(Message.MESSAGE_NETWORK_ERR)
                    .setPositiveButton(getString(R.string.retry),
                            new DialogInterface.OnClickListener() {
                                @Override
                                public void onClick(DialogInterface dialog, int which) {
                                    sendFileLog();
                                }
                            })
                    .setNegativeButton(getString(R.string.cancel),
                            new DialogInterface.OnClickListener() {
                                @Override
                                public void onClick(DialogInterface dialog, int which) {
                                    dialog.dismiss();
                                    delete(Constants.STRING_EMPTY);
                                    clearAndLogout();
                                }
                            });
        }

        android.support.v7.app.AlertDialog alert = dialog.show();
        TextView messageText = (TextView) alert.findViewById(android.R.id.message);
        assert messageText != null;
        messageText.setGravity(Gravity.CENTER);
    }

    /**
     * Clear data and log out
     */
    private void clearAndLogout() {

        // Process logout
        // clear table
        DatabaseManager.initializeInstance(
                new DatabaseHelper(getApplicationContext()));
        DatabaseHelper ds = new DatabaseHelper(this);
        ds.clearTables();

        // stop process loading screen
        progress.dismiss();

        finishAffinity();
        // move to login screen
        Intent intent = new Intent(this, LoginActivity.class);

        startActivity(intent);
    }

    /**
     * Delete file
     */
    public void delete(String fileName) {

        // clear log file for new session user
        File file = new File(Log4JHelper.fileName);
        if (file.exists()) {
            file.delete();
        }
        if (!fileName.isEmpty()) {
            File fileGz = new File(fileName);
            if (fileGz.exists()) {
                fileGz.delete();
            }
        }
    }

    @Override
    public void progressFinish(String output, int typeLocation, String fileName) {

        if (output.contains(Message.CODE_200)) {
            countFile++;
            delete(fileName);
        } else {
            countFile++;
        }
        if (countFile == files.length) {

            clearAndLogout();
        }
    }

    /**
     * Send log file
     */
    public void sendFileLog() {

        // check network status
        if (checkNetwork()) {

            File root = new File(Environment.getExternalStorageDirectory(), "returncandidate_log");
            files = root.listFiles();

            for (File file : files) {
                String[] params =
                        new String[]{Config.API_POST_FILE, file.toString()};
                new HttpPostFile(this).execute(params);
            }
        } else {
            // stop process loading screen
            progress.dismiss();
            showAlertDialog(false);
        }
    }

    /**
     * Check network
     *
     * @return network already : true
     * network not ready : false
     */
    public boolean checkNetwork() {
        ConnectivityManager connectivityManager = (ConnectivityManager) this.getSystemService(
                Context.CONNECTIVITY_SERVICE);
        NetworkInfo networkInfo = connectivityManager.getActiveNetworkInfo();
        return !(networkInfo == null || !networkInfo.isConnected() || !networkInfo.isAvailable());
    }

    /**
     * even when click back
     */
    @Override
    public void onBackPressed() {

        super.onBackPressed();
        // print log end process
        LogManager.i(TAG, Message.TAG_UNLOCK_ACTIVITY + Message.MESSAGE_ACTIVITY_END);

        finishAffinity();
    }

    /**
     * compress file
     */
    private void compressFile() {

        @SuppressLint("SimpleDateFormat") SimpleDateFormat dateFormat = new SimpleDateFormat(
                "yyyyMMddHHmmss");
        Calendar cal = Calendar.getInstance();
        String strDate = dateFormat.format(cal.getTime());
        // compress file log
        String filepath = Log4JHelper.fileName;

        File root = new File(Environment.getExternalStorageDirectory(), "returncandidate_log");
        if (!root.exists()) {
            root.mkdirs();
        }

        userID = userID.replaceAll("__", "_");
        serverName = serverName.replaceAll("__", "_");
        shopID = shopID.replaceAll("__", "_");
        //Gzip file name
        String gzipFileName = root + "/" + serverName + "__" + shopID + "__" + userID + "__"
                + strDate + ".log.gz";
        GzipFile gzipFile = new GzipFile();
        gzipFile.compressGzipFile(filepath, gzipFileName);
    }
}
