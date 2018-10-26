package com.fjn.magazinereturncandidate.db.models;

import android.annotation.SuppressLint;
import android.content.ContentValues;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;

import com.fjn.magazinereturncandidate.common.utils.DatabaseManagerCommon;
import com.fjn.magazinereturncandidate.db.entity.UsersEntity;

import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_LICENSE;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_LOGIN_KEY;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_NAME;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_PASSWORD;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_SERVER_NAME;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_SHOP_ID;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_SHOP_NAME;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_UID;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_USER_ID;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_USER_ROLE;
import static com.fjn.magazinereturncandidate.common.constants.Constants.TABLE_USER;
import static com.fjn.magazinereturncandidate.common.utils.EnDecryptInfoCommon.encryptMD5;
import static com.fjn.magazinereturncandidate.common.utils.EnDecryptInfoCommon.decryptString;


/**
 * User model
 *
 * @author cong-pv
 * @since 2018-10-15
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
     */
    public static String createTable() {

        return String.format(
                "CREATE TABLE %s(%s TEXT PRIMARY KEY, %s TEXT, %s TEXT,%s TEXT, %s "
                        + "TEXT, %s TEXT, %s TEXT, %s TEXT, %s INTEGER, %s TEXT)",
                TABLE_USER, COLUMN_USER_ID, COLUMN_PASSWORD, COLUMN_NAME, COLUMN_UID,
                COLUMN_SHOP_ID, COLUMN_SHOP_NAME, COLUMN_LOGIN_KEY, COLUMN_SERVER_NAME,
                COLUMN_USER_ROLE, COLUMN_LICENSE);
    }

    /**
     * Function insert into table Users
     */
    public int insert(UsersEntity usersEntity) {

        int courseId;
        db = DatabaseManagerCommon.getInstance().openDatabase();
        ContentValues values = new ContentValues();
        values.put(COLUMN_USER_ID, usersEntity.getUserid());
        values.put(COLUMN_PASSWORD, usersEntity.getPassword());
        values.put(COLUMN_NAME, usersEntity.getName());
        values.put(COLUMN_UID, usersEntity.getUid());
        values.put(COLUMN_SHOP_ID, usersEntity.getShop_id());
        values.put(COLUMN_SHOP_NAME, usersEntity.getShop_name());
        values.put(COLUMN_LOGIN_KEY, usersEntity.getLogin_key());
        values.put(COLUMN_SERVER_NAME, usersEntity.getServer_name());
        values.put(COLUMN_USER_ROLE, usersEntity.getRole());
        values.put(COLUMN_LICENSE, usersEntity.getLicense());

        // Inserting Row
        courseId = (int) db.insert(TABLE_USER, null, values);
        DatabaseManagerCommon.getInstance().closeDatabase();

        return courseId;
    }

    /**
     * Get user info
     */
    public UsersEntity getUserInfo() {

        db = DatabaseManagerCommon.getInstance().openDatabase();
        String selectQuery = String.format("SELECT * FROM %s LIMIT 1", TABLE_USER);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        UsersEntity users = null;
        if (cursor != null) {
            if (cursor.moveToFirst()) {
                users = new UsersEntity();
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
                users.setLicense(decryptString(cursor.getString(cursor.getColumnIndex(COLUMN_LICENSE))));
            }
        } else {
            return null;
        }

        DatabaseManagerCommon.getInstance().closeDatabase();
        return users;
    }

    /**
     * Check data is exist in table
     */
    public boolean checkIsData() {

        db = DatabaseManagerCommon.getInstance().openDatabase();
        String sqlQuery = String.format("SELECT * FROM %s", TABLE_USER);

        Cursor cursor = db.rawQuery(sqlQuery, null);
        if (cursor.getCount() <= 0) {
            cursor.close();
            return false;
        }
        cursor.close();
        DatabaseManagerCommon.getInstance().closeDatabase();
        return true;
    }

    /**
     * Check data is exist in table
     */
    public boolean checkDataIsExist(String password) {

        db = DatabaseManagerCommon.getInstance().openDatabase();
        String sqlQuery = String.format("SELECT * FROM %s WHERE %s = '%s'", TABLE_USER,
                COLUMN_PASSWORD, encryptMD5(password));

        Cursor cursor = db.rawQuery(sqlQuery, null);
        if (cursor.getCount() <= 0) {
            cursor.close();
            return false;
        }
        cursor.close();
        DatabaseManagerCommon.getInstance().closeDatabase();
        return true;
    }

}