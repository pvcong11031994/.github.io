package com.android.productchange.db.models;

import static com.android.productchange.common.constants.Constants.COLUMN_CREATE_DATE;
import static com.android.productchange.common.constants.Constants.COLUMN_LOGIN_KEY;
import static com.android.productchange.common.constants.Constants.COLUMN_NAME;
import static com.android.productchange.common.constants.Constants.COLUMN_PASSWORD;
import static com.android.productchange.common.constants.Constants.COLUMN_SERVER_NAME;
import static com.android.productchange.common.constants.Constants.COLUMN_SHOP_ID;
import static com.android.productchange.common.constants.Constants.COLUMN_SHOP_NAME;
import static com.android.productchange.common.constants.Constants.COLUMN_UID;
import static com.android.productchange.common.constants.Constants.COLUMN_USER_ID;
import static com.android.productchange.common.constants.Constants.COLUMN_USER_ROLE;
import static com.android.productchange.common.constants.Constants.TABLE_USER;

import android.annotation.SuppressLint;
import android.content.ContentValues;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.util.Base64;

import com.android.productchange.common.utils.DatabaseManager;
import com.android.productchange.db.entity.Users;

import java.math.BigInteger;
import java.security.InvalidKeyException;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Arrays;

import javax.crypto.BadPaddingException;
import javax.crypto.Cipher;
import javax.crypto.IllegalBlockSizeException;
import javax.crypto.NoSuchPaddingException;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;

/**
 * <h1>Model User</h1>
 *
 * @author tien-lv
 * @since 2017-12-06
 */
public class UserModel {

    /**
     * SQLite Database
     */
    private SQLiteDatabase db;

    /**
     * Constructor Model User
     */
    public UserModel() {

    }

    /**
     * Function create table User
     *
     * @return String
     */
    public static String createTable() {

        return String.format(
                "CREATE TABLE %s(%s TEXT PRIMARY KEY, %s TEXT, %s TEXT,%s TEXT, %s "
                        + "TEXT, %s TEXT, %s TEXT, %s TEXT, %s INTEGER,%s TEXT)",
                TABLE_USER, COLUMN_USER_ID, COLUMN_PASSWORD, COLUMN_NAME, COLUMN_UID,
                COLUMN_SHOP_ID, COLUMN_SHOP_NAME, COLUMN_LOGIN_KEY, COLUMN_SERVER_NAME,
                COLUMN_USER_ROLE, COLUMN_CREATE_DATE);
    }

    /**
     * Function insert into table Users
     *
     * @param users is entity Users
     * @return int
     */
    public int insert(Users users) {

        int courseId;
        db = DatabaseManager.getInstance().openDatabase();
        ContentValues values = new ContentValues();
        values.put(COLUMN_USER_ID, users.getUserid());
        values.put(COLUMN_PASSWORD, users.getPassword());
        values.put(COLUMN_NAME, users.getName());
        values.put(COLUMN_UID, users.getUid());
        values.put(COLUMN_SHOP_ID, users.getShop_id());
        values.put(COLUMN_SHOP_NAME, users.getShop_name());
        values.put(COLUMN_LOGIN_KEY, users.getLogin_key());
        values.put(COLUMN_SERVER_NAME, users.getServer_name());
        values.put(COLUMN_USER_ROLE, users.getRole());
        values.put(COLUMN_CREATE_DATE, users.getCreate_date());

        // Inserting Row
        courseId = (int) db.insert(TABLE_USER, null, values);
        DatabaseManager.getInstance().closeDatabase();

        return courseId;
    }

    /**
     * Get user info in table users
     *
     * @return Users
     */
    public Users getUserInfo() {

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT * FROM %s LIMIT 1", TABLE_USER);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        Users users = null;
        if (cursor != null) {
            if (cursor.moveToFirst()) {
                users = new Users();
                users.setUserID(cursor.getString(cursor.getColumnIndex(COLUMN_USER_ID)));
                users.setName(cursor.getString(cursor.getColumnIndex(COLUMN_NAME)));
                users.setUid(cursor.getString(cursor.getColumnIndex(COLUMN_UID)));
                users.setShop_id(
                        decryptString(cursor.getString(cursor.getColumnIndex(COLUMN_SHOP_ID))));
                users.setShop_name(cursor.getString(cursor.getColumnIndex(COLUMN_SHOP_NAME)));
                users.setLogin_key(decryptString(cursor.getString(cursor.getColumnIndex(COLUMN_LOGIN_KEY))));
                users.setServer_name(
                        decryptString(cursor.getString(cursor.getColumnIndex(COLUMN_SERVER_NAME))));
                users.setRole(Integer.parseInt(
                        cursor.getString(cursor.getColumnIndex(COLUMN_USER_ROLE))));
                users.setCreate_date(
                        cursor.getString(cursor.getColumnIndex(COLUMN_CREATE_DATE)));

            }
        } else {
            return null;
        }

        DatabaseManager.getInstance().closeDatabase();
        return users;
    }

    private String decryptString(String strEnDeCrypt) {
        String TOKEN_KEY = "fqJfdzGDvfwbedsKSUGty3VZ9taXxMVw";
        try {
            byte[] ivAndCipherText = Base64.decode(strEnDeCrypt, Base64.NO_WRAP);
            byte[] iv = Arrays.copyOfRange(ivAndCipherText, 0, 16);
            byte[] cipherText = Arrays.copyOfRange(ivAndCipherText, 16, ivAndCipherText.length);

            Cipher cipher = Cipher.getInstance("AES/CBC/PKCS5Padding");
            cipher.init(Cipher.DECRYPT_MODE, new SecretKeySpec(TOKEN_KEY.getBytes("utf-8"), "AES"), new IvParameterSpec(iv));
            return new String(cipher.doFinal(cipherText), "utf-8");
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }

    /**
     * Check data is exist in table
     *
     * @return boolean
     */
    public boolean checkIsData() {

        db = DatabaseManager.getInstance().openDatabase();
        String sqlQuery = String.format("SELECT * FROM %s", TABLE_USER);

        Cursor cursor = db.rawQuery(sqlQuery, null);
        if (cursor.getCount() <= 0) {
            cursor.close();
            return false;
        }
        cursor.close();
        DatabaseManager.getInstance().closeDatabase();
        return true;
    }

    /**
     * Check data is exist in table
     *
     * @return boolean
     */
    public boolean checkDataIsExist(String password) {


        db = DatabaseManager.getInstance().openDatabase();
        String sqlQuery = String.format("SELECT * FROM %s WHERE %s = '%s'", TABLE_USER,
                COLUMN_PASSWORD, encryptMD5(password));

        Cursor cursor = db.rawQuery(sqlQuery, null);
        if (cursor.getCount() <= 0) {
            cursor.close();
            return false;
        }
        cursor.close();
        DatabaseManager.getInstance().closeDatabase();
        return true;
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
}