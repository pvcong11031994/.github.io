package com.android.returncandidate.activities;

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
import com.android.returncandidate.db.entity.*;
import com.android.returncandidate.db.models.*;

import javax.crypto.Cipher;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;

import java.math.BigInteger;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;

import android.util.Base64;

import org.json.*;

/**
 * Login activity<br>
 * Flow : {@link LoginActivity} ▶ {@link DataLoaderActivity}
 *
 * @author minh-th
 * @version 2.0
 * @since 2018-05-10
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

    /**
     * Initialize screen layout.<br>
     * Base on login status will activity be redirected to unlock screen or login screen
     *
     * @param savedInstanceState {@link Bundle }
     */
    @Override
    protected void onCreate(Bundle savedInstanceState) {

        LogManager.i(TAG, Message.TAG_LOGIN_ACTIVITY + Message.MESSAGE_ACTIVITY_START);
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
        DatabaseManager.initializeInstance(new DatabaseHelper(getApplicationContext()));
        userModel = new UserModel();

        // Check if user is login previously
        // If user is already logged in -> transition to unlock screen
        // Else -> transition to login screen
        if (userModel.checkIsData()) {

            Users userInfo = userModel.getUserInfo();
            loginKey = userInfo.getLogin_key();
            shopID = userInfo.getShop_id();
            userID = userInfo.getUserid();
            userRole = userInfo.getRole();
            serverName = userInfo.getServer_name();
            flagLogin = false;
            license = userInfo.getLicense();
            Log.d("LINCONGPLOGIN", license);
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

        if (checkNetwork()) {

            LogManager.i(TAG, Message.TAG_LOGIN_ACTIVITY + Message.LOADING_DATA_FROM_SERVER);
            progress.show();

            // Confirm login user
            String[] params = new String[]{Config.CODE_LOGIN, Config.API_LOGIN,
                    edtUserID.getText().toString(), edtPassword.getText().toString()};
            new HttpPost(this).execute(params);
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
        LogManager.i(TAG, Message.TAG_LOGIN_ACTIVITY + Message.MESSAGE_ACTIVITY_END);

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

        LogManager.i(TAG, Message.TAG_LOGIN_ACTIVITY + Message.LOADING_DATA_FROM_SERVER_SUCCESS);
        if (output.contains(Message.CODE_200)) {
            // Process OK
            try {
                JSONObject response = new JSONObject(output);
                if (!response.isNull(Config.RESULT_LIST)) {
                    JSONObject res = response.getJSONObject(Config.RESULT_LIST);
                    if (response.getJSONObject(Config.RESULT_LIST).has(Config.LOGIN_KEY)) {

                        Users users = new Users();
                        users.setUserID(res.getString(Constants.COLUMN_USER_ID));
                        users.setName(res.getString(Constants.COLUMN_NAME));
                        users.setUid(res.getString(Constants.COLUMN_UID));
                        users.setShop_id(encryptString(res.getString(Constants.COLUMN_SHOP_ID)));
                        users.setLogin_key(encryptString(res.getString(Constants.COLUMN_LOGIN_KEY)));
                        users.setShop_name(res.getString(Constants.COLUMN_SHOP_NAME));
                        users.setServer_name(encryptString(res.getString(Constants.COLUMN_SERVER_NAME)));
                        users.setRole(Integer.parseInt(res.getString(Constants.COLUMN_USER_ROLE)));
                        users.setPassword(encryptMD5(edtPassword.getText().toString()));
                        users.setLicense(encryptString(res.getString(Constants.COLUMN_LICENSE)));

                        userModel.insert(users);

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
                LogManager.e(Constants.TAG_APPLICATION_NAME, e.toString());
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

    /**
     * Transition to {@link DataLoaderActivity}<br>
     * Transfer user information to next screen
     */
    public void moveNextActivity() {

        // Stop progress dialog
        progress.dismiss();
        LogManager.i(TAG, Message.TAG_LOGIN_ACTIVITY
                + String.format(Message.MESSAGE_LOGIN_SUCCESS, userID, shopID, serverName));
        LogManager.i(TAG, Message.TAG_LOGIN_ACTIVITY + Message.MESSAGE_ACTIVITY_END);

        LogManager.i(TAG, String.format(Message.MESSAGE_ACTIVITY_MOVE,
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
     * Check network
     *
     * @return This returns true if network connectivity is okay
     */
    public boolean checkNetwork() {
        ConnectivityManager connectivityManager = (ConnectivityManager) this.getSystemService(
                Context.CONNECTIVITY_SERVICE);
        NetworkInfo networkInfo = connectivityManager.getActiveNetworkInfo();
        return !(networkInfo == null || !networkInfo.isConnected() || !networkInfo.isAvailable());
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

}
