package com.android.productchange.activities;

import android.app.ProgressDialog;
import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;
import android.net.ConnectivityManager;
import android.net.NetworkInfo;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.text.TextUtils;
import android.util.Base64;
import android.view.Gravity;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;

import com.android.productchange.R;
import com.android.productchange.api.Config;
import com.android.productchange.api.HttpPost;
import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.constants.Message;
import com.android.productchange.common.helpers.DatabaseHelper;
import com.android.productchange.common.utils.DatabaseManager;
import com.android.productchange.common.utils.LogManager;
import com.android.productchange.db.entity.Users;
import com.android.productchange.db.models.UserModel;
import com.android.productchange.interfaces.HttpResponse;

import org.json.JSONException;
import org.json.JSONObject;

import java.math.BigInteger;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.text.SimpleDateFormat;
import java.util.Calendar;

import javax.crypto.Cipher;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;


/**
 * <h1>Login Screen Activity</h1>
 * <p>
 * Validate user and password input
 * Send data login to server
 * Get result from server
 *
 * @author tien-lv
 * @since 2018-02-08
 */
public class LoginActivity extends AppCompatActivity implements View.OnClickListener, HttpResponse {

    /**
     * TAG
     */
    private String TAG = Constants.TAG_APPLICATION_NAME;
    /**
     * Edit text user id.
     */
    private EditText edtUserID;

    /**
     * Edit text password.
     */
    private EditText edtPassword;

    /**
     * Response login key.
     */
    private String loginKey;

    /**
     * Response server name.
     */
    private String serverName;

    /**
     * Response shop id.
     */
    private String shopID;

    /**
     * Response user id.
     */
    private String userID;

    /**
     * Response user role.
     */
    private int userRole;

    /**
     * Response create date.
     */
    private String createDate;

    /**
     * Check login first
     */
    private boolean flagLogin;

    /**
     * Model user.
     */
    UserModel userModel;

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

        LogManager.i(TAG, Message.TAG_LOGIN_ACTIVITY + Message.MESSAGE_ACTIVITY_START);
        super.onCreate(savedInstanceState);

        //Call layout login screen.
        setContentView(R.layout.activity_login);

        Button btnLogin = (Button) findViewById(R.id.btn_login);
        //Button Login on click event.
        btnLogin.setOnClickListener(this);

        edtUserID = (EditText) findViewById(R.id.edt_userid);
        edtPassword = (EditText) findViewById(R.id.edt_password);

        // init process loading screen
        progress = new ProgressDialog(this);
        progress.setMessage(Message.MESSAGE_LOADING_SCREEN);
        progress.setCancelable(false);

        // init database
        DatabaseManager.initializeInstance(new DatabaseHelper(getApplicationContext()));
        userModel = new UserModel();

        // check use is login previous
        // case login previous -> by pass login
        // case not login previous -> show login screen
        if (userModel.checkIsData()) {
            // get user info
            Users userInfo = userModel.getUserInfo();
            loginKey = userInfo.getLogin_key();
            shopID = userInfo.getShop_id();
            userID = userInfo.getUserid();
            userRole = userInfo.getRole();
            serverName = userInfo.getServer_name();
            createDate = userInfo.getCreate_date();
            flagLogin = false;

            moveNextActivity();
        }
    }

    /**
     * Override On click event
     *
     * @param v is View on click listener
     */
    @Override
    public void onClick(View v) {

        connect();
    }

    /**
     * This is function connect to server
     */
    private void connect() {

        // Validate input
        if (!checkInput()) {
            return;
        }

        // Check network is connected
        if (checkNetwork()) {
            // case network ok: check user with API server
            LogManager.i(TAG, Message.TAG_LOGIN_ACTIVITY + Message.LOADING_DATA_FROM_SERVER);
            progress.show();
            // Request params
            String[] params = new String[]{Config.CODE_LOGIN, Config.API_LOGIN,
                    edtUserID.getText().toString(), edtPassword.getText().toString()};
            new HttpPost(this).execute(params);
        } else {
            // case network not ok: show dialog warning
            progress.dismiss();
            showDialog();
        }
    }

    /**
     * Function validate
     * user id input
     * password input
     *
     * @return data input is correct: true
     * data input not correct: false
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
     * Override press back for exit application.
     */
    @Override
    public void onBackPressed() {
        LogManager.i(TAG, Message.TAG_LOGIN_ACTIVITY + Message.MESSAGE_ACTIVITY_END);
        finishAffinity();
    }

    /**
     * Override process when get json result from Server
     *
     * @param fileName         is file name
     * @param output           is result code form server
     * @param multiThreadCount is count multi thread
     * @throws JSONException on json result error
     * @see JSONException
     */
    @Override
    public void progressFinish(String output, int multiThreadCount, String fileName) {

        LogManager.i(TAG, Message.TAG_LOGIN_ACTIVITY + Message.LOADING_DATA_FROM_SERVER_SUCCESS);
        if (output.contains(Message.CODE_200)) {
            // process OK
            try {
                JSONObject response = new JSONObject(output);
                if (!response.isNull(Config.RESULT_LIST)) {
                    JSONObject res = response.getJSONObject(Config.RESULT_LIST);
                    // check and insert data form result to table Users
                    if (response.getJSONObject(Config.RESULT_LIST).has(Config.LOGIN_KEY)) {

                        Users users = new Users();

                        //Get info shop_cd, server_name and login_key
                        users.setUserID(res.getString(Constants.COLUMN_USER_ID));
                        users.setName(res.getString(Constants.COLUMN_NAME));
                        users.setUid(res.getString(Constants.COLUMN_UID));
                        users.setShop_id(encryptString(res.getString(Constants.COLUMN_SHOP_ID)));
                        users.setLogin_key(encryptString(res.getString(Constants.COLUMN_LOGIN_KEY)));
                        users.setShop_name(res.getString(Constants.COLUMN_SHOP_NAME));
                        users.setServer_name(encryptString(res.getString(Constants.COLUMN_SERVER_NAME)));
                        users.setRole(Integer.parseInt(res.getString(Constants.COLUMN_USER_ROLE)));
                        users.setPassword(encryptMD5(edtPassword.getText().toString()));
                        Calendar dateDefault = Calendar.getInstance();
                        SimpleDateFormat df = new SimpleDateFormat(Constants.DATE_FORMAT_STRING);
                        String formattedDate = df.format(dateDefault.getTime());

                        users.setCreate_date(formattedDate);

                        userModel.insert(users);

                        loginKey = res.getString(Constants.COLUMN_LOGIN_KEY);
                        shopID = res.getString(Constants.COLUMN_SHOP_ID);
                        userID = res.getString(Constants.COLUMN_USER_ID);
                        userRole = Integer.parseInt(res.getString(Constants.COLUMN_USER_ROLE));
                        serverName = res.getString(Constants.COLUMN_SERVER_NAME);
                        createDate = formattedDate;
                        flagLogin = true;
                    }
                }
            } catch (JSONException e) {
                e.printStackTrace();
                LogManager.e(Constants.TAG_APPLICATION_NAME, e.toString());
            }
            moveNextActivity();

        } else {
            // process ERROR
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
                // no response result
            } else {
                alertDialogBuilder.setMessage(Message.MESSAGE_RESULT_EMPTY);
            }

            // show dialog warning
            alertDialogBuilder.setNegativeButton("OK", new DialogInterface.OnClickListener() {
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

    //encrypt and decrypt DES
    private String encryptString(String strEnDeCrypt) {

        String TOKEN_KEY = "fqJfdzGDvfwbedsKSUGty3VZ9taXxMVw";
        try {
            byte[] iv = new byte[16];
            new SecureRandom().nextBytes(iv);
            Cipher cipher = Cipher.getInstance("AES/CBC/PKCS5Padding");
            cipher.init(Cipher.ENCRYPT_MODE, new SecretKeySpec(TOKEN_KEY.getBytes("utf-8"), "AES"), new IvParameterSpec(iv));
            byte[] cipherText = cipher.doFinal(strEnDeCrypt.getBytes("utf-8"));
            byte[] ivAndCipherText = getCombinedArray(iv, cipherText);
            return Base64.encodeToString(ivAndCipherText, Base64.NO_WRAP);
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }

    //Encrypt MD5 Password
    private static String encryptMD5(String strEncrypt) {

        try {
            MessageDigest md = MessageDigest.getInstance("MD5");
            byte[] messageDigest = md.digest(strEncrypt.getBytes());
            BigInteger number = new BigInteger(1, messageDigest);
            String hashtext = number.toString(16);
            while (hashtext.length() < 32) {
                hashtext = "0" + hashtext;
            }
            return hashtext;
        } catch (NoSuchAlgorithmException e) {
            e.printStackTrace();
            return null;
        }
    }


    private static byte[] getCombinedArray(byte[] one, byte[] two) {

        byte[] combined = new byte[one.length + two.length];
        for (int i = 0; i < combined.length; ++i) {
            combined[i] = i < one.length ? one[i] : two[i - one.length];
        }
        return combined;
    }

    /**
     * Move to Classify Screen with user is admin
     * Move to Barcode Screen with user is employees
     */
    public void moveNextActivity() {
        // stop process loading screen
        progress.dismiss();
        LogManager.i(TAG, Message.TAG_LOGIN_ACTIVITY
                + String.format(Message.MESSAGE_LOGIN_SUCCESS, userID, shopID, serverName));
        LogManager.i(TAG, Message.TAG_LOGIN_ACTIVITY + Message.MESSAGE_ACTIVITY_END);
        Intent intent;

        // role User
        LogManager.i(TAG, String.format(Message.MESSAGE_ACTIVITY_MOVE,
                Message.LOGIN_ACTIVITY_NAME, Message.SCANNER_ACTIVITY_NAME));
        intent = new Intent(this, ProductChangeMainActivity.class);

        // setting info for next screen
        Bundle bundle = new Bundle();
        bundle.putString(Config.LOGIN_KEY, loginKey);
        bundle.putString(Constants.COLUMN_USER_ID, userID);
        bundle.putString(Constants.COLUMN_SHOP_ID, shopID);
        bundle.putString(Constants.COLUMN_SERVER_NAME, serverName);
        bundle.putInt(Constants.COLUMN_USER_ROLE, userRole);
        bundle.putBoolean(Constants.FLAG_LOGIN, flagLogin);
        bundle.putString(Constants.COLUMN_CREATE_DATE, createDate);
        intent.putExtras(bundle);

        startActivity(intent);
        finish();
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
     * Show dialog warning logout
     */
    private void showDialog() {

        progress.show();
        android.support.v7.app.AlertDialog.Builder dialog =
                new android.support.v7.app.AlertDialog.Builder(this);
        dialog
                .setMessage(Message.MESSAGE_NETWORK_ERR)
                .setCancelable(false)
                .setPositiveButton(getString(R.string.retry),
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                connect();
                            }
                        })
                .setNegativeButton(getString(R.string.cancel),
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
