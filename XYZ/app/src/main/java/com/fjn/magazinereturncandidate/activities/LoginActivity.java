package com.fjn.magazinereturncandidate.activities;

import android.app.ProgressDialog;
import android.content.DialogInterface;
import android.content.Intent;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.text.TextUtils;
import android.view.Gravity;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;

import com.fjn.magazinereturncandidate.R;
import com.fjn.magazinereturncandidate.api.Config;
import com.fjn.magazinereturncandidate.api.HttpPostUserLogin;
import com.fjn.magazinereturncandidate.api.HttpResponse;
import com.fjn.magazinereturncandidate.common.constants.Constants;
import com.fjn.magazinereturncandidate.common.constants.Message;
import com.fjn.magazinereturncandidate.common.helpers.DatabaseHelper;
import com.fjn.magazinereturncandidate.common.utils.DatabaseManagerCommon;
import com.fjn.magazinereturncandidate.common.utils.EnDecryptInfoCommon;
import com.fjn.magazinereturncandidate.common.utils.LogManagerCommon;
import com.fjn.magazinereturncandidate.db.entity.UsersEntity;
import com.fjn.magazinereturncandidate.db.models.UserModel;

import org.json.JSONException;
import org.json.JSONObject;

import static com.fjn.magazinereturncandidate.common.constants.Message.MESSAGE_CANCEL;
import static com.fjn.magazinereturncandidate.common.constants.Message.MESSAGE_RETRY;
import static com.fjn.magazinereturncandidate.common.utils.EnDecryptInfoCommon.encryptMD5;
import static com.fjn.magazinereturncandidate.common.utils.NetworkCommon.isNetworkConnection;

/**
 * Login activity<br>
 * Flow : {@link LoginActivity} ▶ {@link DataLoaderActivity}
 *
 * @author cong-pv
 * @version 1.0
 * @since 2018-10-15
 */
public class LoginActivity extends AppCompatActivity implements View.OnClickListener, HttpResponse {

    private String TAG = Constants.TAG_APPLICATION_NAME;

    private EditText edtUserID;
    private EditText edtPassword;

    private String loginKey;
    private String license;
    private String serverName;
    private String shopID;
    private String userID;
    private int userRole;

    private boolean flagLogin;

    private UserModel userModel;
    private ProgressDialog progress;

    private EnDecryptInfoCommon enDecryptInfoCommon;

    /**
     * Initialize screen layout.<br>
     * Base on login status will activity be redirected to unlock screen or login screen
     *
     * @param savedInstanceState {@link Bundle }
     */
    @Override
    protected void onCreate(Bundle savedInstanceState) {

        LogManagerCommon.i(TAG, Message.TAG_LOGIN_ACTIVITY + Message.MESSAGE_ACTIVITY_START);
        super.onCreate(savedInstanceState);

        setContentView(R.layout.activity_login);

        // Login button init
        Button btnLogin = (Button) findViewById(R.id.btn_login);
        btnLogin.setOnClickListener(this);

        // Element assignment
        edtUserID = (EditText) findViewById(R.id.edt_userid);
        edtPassword = (EditText) findViewById(R.id.edt_password);

        // Progress dialog init
        progress = new ProgressDialog(this);
        progress.setMessage(Message.MESSAGE_LOADING_SCREEN);
        progress.setCancelable(false);

        // Database init
        DatabaseManagerCommon.initializeInstance(new DatabaseHelper(getApplicationContext()));
        userModel = new UserModel();
        enDecryptInfoCommon = new EnDecryptInfoCommon();

        // Check if user is login previously
        // If user is already logged in -> transition to unlock screen
        // Else -> transition to login screen
        if (userModel.checkIsData()) {

            UsersEntity userInfo = userModel.getUserInfo();
            loginKey = userInfo.getLogin_key();
            shopID = userInfo.getShop_id();
            userID = userInfo.getUserid();
            userRole = userInfo.getRole();
            serverName = userInfo.getServer_name();
            flagLogin = false;
            license = userInfo.getLicense();

            moveNextActivity();
        }
    }

    /**
     * btnLogin onClick event handler
     *
     * @param v {@link View }
     */
    @Override
    public void onClick(View v) {

        connect();
    }

    /**
     * Connect to API server
     */
    private void connect() {

        // Validate input
        if (!checkInput()) {
            return;
        }

        if (isNetworkConnection(LoginActivity.this)) {

            LogManagerCommon.i(TAG, Message.TAG_LOGIN_ACTIVITY + Message.LOADING_DATA_FROM_SERVER);
            progress.show();

            // Confirm login user
            String[] params = new String[]{Config.CODE_LOGIN, Config.API_LOGIN,
                    edtUserID.getText().toString(), edtPassword.getText().toString()};
            new HttpPostUserLogin(this).execute(params);
        } else {

            // Network failed
            progress.dismiss();
            showDialog();
        }
    }

    /**
     * Input validation
     */
    private boolean checkInput() {

        // Validate user id and password is empty
        if (TextUtils.isEmpty(edtUserID.getText().toString())) {
            edtUserID.setError(
                    String.format(Message.MESSAGE_CHECK_INPUT_EMPTY, getString(R.string.userid)));
            edtUserID.requestFocus();
            return false;
        } else if (TextUtils.isEmpty(edtPassword.getText().toString())) {
            edtPassword.setError(
                    String.format(Message.MESSAGE_CHECK_INPUT_EMPTY, getString(R.string.password)));
            edtPassword.requestFocus();
            return false;
        }

        return true;
    }

    /**
     * Back event handler
     */
    @Override
    public void onBackPressed() {
        LogManagerCommon.i(TAG, Message.TAG_LOGIN_ACTIVITY + Message.MESSAGE_ACTIVITY_END);

        finishAffinity();
    }

    /**
     * Handling response from API server
     *
     * @param output       {@link String }
     * @param typeLocation int
     * @param fileName     {@link String}
     */
    @Override
    public void progressFinish(String output, int typeLocation, String fileName) {

        LogManagerCommon.i(TAG, Message.TAG_LOGIN_ACTIVITY + Message.LOADING_DATA_FROM_SERVER_SUCCESS);
        if (output.contains(Message.CODE_200)) {
            // Process OK
            try {
                JSONObject response = new JSONObject(output);
                if (!response.isNull(Config.RESULT_LIST)) {
                    JSONObject res = response.getJSONObject(Config.RESULT_LIST);
                    if (response.getJSONObject(Config.RESULT_LIST).has(Config.LOGIN_KEY)) {

                        UsersEntity usersEntity = new UsersEntity();
                        usersEntity.setUserID(res.getString(Constants.COLUMN_USER_ID));
                        usersEntity.setName(res.getString(Constants.COLUMN_NAME));
                        usersEntity.setUid(res.getString(Constants.COLUMN_UID));
                        usersEntity.setShop_id(enDecryptInfoCommon.encryptString(res.getString(Constants.COLUMN_SHOP_ID)));
                        usersEntity.setLogin_key(enDecryptInfoCommon.encryptString(res.getString(Constants.COLUMN_LOGIN_KEY)));
                        usersEntity.setShop_name(res.getString(Constants.COLUMN_SHOP_NAME));
                        usersEntity.setServer_name(enDecryptInfoCommon.encryptString(res.getString(Constants.COLUMN_SERVER_NAME)));
                        usersEntity.setRole(Integer.parseInt(res.getString(Constants.COLUMN_USER_ROLE)));
                        usersEntity.setPassword(encryptMD5(edtPassword.getText().toString()));
                        usersEntity.setLicense(enDecryptInfoCommon.encryptString(res.getString(Constants.COLUMN_LICENSE)));

                        userModel.insert(usersEntity);

                        loginKey = res.getString(Constants.COLUMN_LOGIN_KEY);
                        shopID = res.getString(Constants.COLUMN_SHOP_ID);
                        userID = res.getString(Constants.COLUMN_USER_ID);
                        userRole = Integer.parseInt(res.getString(Constants.COLUMN_USER_ROLE));
                        serverName = res.getString(Constants.COLUMN_SERVER_NAME);
                        flagLogin = true;
                        license = res.getString(Constants.COLUMN_LICENSE);
                    }
                }
            } catch (JSONException e) {
                e.printStackTrace();
                LogManagerCommon.e(Constants.TAG_APPLICATION_NAME, e.toString());
            }
            moveNextActivity();

        } else {
            // Process ERROR
            progress.dismiss();
            android.support.v7.app.AlertDialog.Builder alertDialogBuilder =
                    new android.support.v7.app.AlertDialog.Builder(this);

            // Error login
            if (output.contains(Message.CODE_401)) {
                alertDialogBuilder.setMessage(Message.MESSAGE_401);
                // Error API URL not found
            } else if (output.contains(Message.CODE_404)) {
                alertDialogBuilder.setMessage(Message.MESSAGE_404);
                // Error server
            } else if (output.contains(Message.CODE_500)) {
                alertDialogBuilder.setMessage(Message.MESSAGE_500);
                // No response result
            } else {
                alertDialogBuilder.setMessage(Message.MESSAGE_RESULT_EMPTY);
            }

            // Show warning dialog
            alertDialogBuilder.setNegativeButton(Message.CONFIRM_OK, new DialogInterface.OnClickListener() {
                @Override
                public void onClick(DialogInterface dialog, int which) {
                    dialog.dismiss();
                }
            });
            android.support.v7.app.AlertDialog alert = alertDialogBuilder.show();
            TextView messageText = (TextView) alert.findViewById(android.R.id.message);
            assert messageText != null;
            messageText.setGravity(Gravity.CENTER);
        }
    }

    /**
     * Transition to {@link DataLoaderActivity}<br>
     * Transfer user information to next screen
     */
    public void moveNextActivity() {

        // Stop progress dialog
        progress.dismiss();
        LogManagerCommon.i(TAG, Message.TAG_LOGIN_ACTIVITY
                + String.format(Message.MESSAGE_LOGIN_SUCCESS, userID, shopID, serverName));
        LogManagerCommon.i(TAG, Message.TAG_LOGIN_ACTIVITY + Message.MESSAGE_ACTIVITY_END);

        LogManagerCommon.i(TAG, String.format(Message.MESSAGE_ACTIVITY_MOVE,
                Message.LOGIN_ACTIVITY_NAME, Message.SCANNER_ACTIVITY_NAME));
        Intent intent = new Intent(this, DataLoaderActivity.class);

        // Setting info for next screen
        Bundle bundle = new Bundle();
        bundle.putString(Config.LOGIN_KEY, loginKey);
        bundle.putString(Constants.COLUMN_USER_ID, userID);
        bundle.putString(Constants.COLUMN_SHOP_ID, shopID);
        bundle.putString(Constants.COLUMN_SERVER_NAME, serverName);
        bundle.putInt(Constants.COLUMN_USER_ROLE, userRole);
        bundle.putBoolean(Constants.FLAG_LOGIN, flagLogin);
        bundle.putString(Constants.COLUMN_LICENSE, license);
        intent.putExtras(bundle);

        startActivity(intent);
        finish();
    }

    /**
     * Warning dialog when failed to connect to API server
     */
    private void showDialog() {

        progress.show();
        android.support.v7.app.AlertDialog.Builder dialog =
                new android.support.v7.app.AlertDialog.Builder(this);
        dialog
                .setMessage(Message.MESSAGE_NETWORK_ERR)
                .setCancelable(false)
                .setPositiveButton(MESSAGE_RETRY,
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                connect();
                            }
                        })
                .setNegativeButton(MESSAGE_CANCEL,
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
}
