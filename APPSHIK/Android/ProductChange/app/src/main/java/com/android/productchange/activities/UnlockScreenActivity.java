package com.android.productchange.activities;

import android.app.ProgressDialog;
import android.content.DialogInterface;
import android.content.Intent;
import android.content.pm.ActivityInfo;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.text.TextUtils;
import android.view.Gravity;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;

import com.android.productchange.R;
import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.constants.Message;
import com.android.productchange.common.helpers.DatabaseHelper;
import com.android.productchange.common.utils.DatabaseManager;
import com.android.productchange.common.utils.LogManager;
import com.android.productchange.db.models.UserModel;

/**
 * <h1>Unlock Screen Activity</h1>
 *
 * Activity show after table Users has data
 * Check Users login again into app when app finish
 *
 * @author tien-lv
 * @since 2018-02-08
 */
public class UnlockScreenActivity extends AppCompatActivity implements View.OnClickListener {

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
     * Server name
     */
    private String serverName;

    /**
     * Create date
     */
    private String createDate;

    /**
     * Edit text password.
     */
    private EditText edtPassword;

    /**
     * Count wrong password input
     */
    private int count = 0;

    /**
     * Progress dialog.
     */
    private ProgressDialog progress;

    /**
     * Init on Create Activity
     *
     * @param savedInstanceState is Bundle of activity
     */
    @Override
    protected void onCreate(Bundle savedInstanceState) {

        // print log begin process
        LogManager.i(TAG, Message.TAG_UNLOCK_ACTIVITY + Message.MESSAGE_ACTIVITY_START);
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_unlock_screen);
        setRequestedOrientation(ActivityInfo.SCREEN_ORIENTATION_PORTRAIT);

        Bundle bundle = getIntent().getExtras();
        if (bundle != null) {
            userID = bundle.getString(Constants.COLUMN_USER_ID);
            shopID = bundle.getString(Constants.COLUMN_SHOP_ID);
            serverName = bundle.getString(Constants.COLUMN_SERVER_NAME);
            createDate=bundle.getString(Constants.COLUMN_CREATE_DATE);
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
     * Override on click event.
     *
     * @param v is View on click listener
     */
    @Override
    public void onClick(View v) {

        switch (v.getId()) {
            case R.id.btn_unlock:
                unLock();
                break;
            case R.id.btn_logout:
                showAlertDialog();
                break;
        }
    }

    /**
     * Unlock user function
     */
    public void unLock() {

        UserModel userModel = new UserModel();

        // Validate input
        if (TextUtils.isEmpty(edtPassword.getText().toString())) {
            edtPassword.setError(
                    String.format(Message.MESSAGE_CHECK_INPUT_EMPTY, getString(R.string.password)));
            edtPassword.requestFocus();
        } else {
            // Validate password is not correct for 3 times
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
                        // 3 times password input not correct clear data and logout
                        if (count >= 3) {
                            progress.show();
                            clearAndLogout();
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

                // move to process product list
                Intent intent = new Intent(this, ProductChangeListActivity.class);
                Bundle bundle = new Bundle();
                bundle.putString(Constants.COLUMN_USER_ID, userID);
                bundle.putString(Constants.COLUMN_SHOP_ID, shopID);
                bundle.putString(Constants.COLUMN_SERVER_NAME, serverName);
                bundle.putString(Constants.COLUMN_CREATE_DATE,createDate);
                intent.putExtras(bundle);
                startActivity(intent);
                finish();
            }
        }
    }

    /**
     * Show dialog warning logout
     */
    public void showAlertDialog() {

        progress.show();
        android.support.v7.app.AlertDialog.Builder dialog =
                new android.support.v7.app.AlertDialog.Builder(this);
        dialog.setCancelable(false);

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
                                clearAndLogout();
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
        DatabaseManager.initializeInstance(new DatabaseHelper(getApplicationContext()));
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
     * Close app when click back press
     */
    @Override
    public void onBackPressed() {

        super.onBackPressed();
        // print log end process
        LogManager.i(TAG, Message.TAG_UNLOCK_ACTIVITY + Message.MESSAGE_ACTIVITY_END);
        finishAffinity();
    }
}
