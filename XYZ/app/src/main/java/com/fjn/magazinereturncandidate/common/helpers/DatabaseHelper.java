package com.fjn.magazinereturncandidate.common.helpers;

import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteOpenHelper;


import com.fjn.magazinereturncandidate.db.models.MaxYearRankModel;
import com.fjn.magazinereturncandidate.db.models.ReturnMagazineModel;
import com.fjn.magazinereturncandidate.db.models.UserModel;

import static com.fjn.magazinereturncandidate.common.constants.Constants.DATABASE_NAME;
import static com.fjn.magazinereturncandidate.common.constants.Constants.DATABASE_VERSION;
import static com.fjn.magazinereturncandidate.common.constants.Constants.QUERY_DROP_TABLE_EXIST;
import static com.fjn.magazinereturncandidate.common.constants.Constants.TABLE_MAX_YEAR_RANK;
import static com.fjn.magazinereturncandidate.common.constants.Constants.TABLE_RETURN_MAGAZINE;
import static com.fjn.magazinereturncandidate.common.constants.Constants.TABLE_USER;

/**
 * Connect to Database SQLite
 *
 * @author cong-pv
 * @since 2018-10-15
 */

public class DatabaseHelper extends SQLiteOpenHelper {

    /**
     * Constructor DatabaseHelper
     */
    public DatabaseHelper(Context context) {

        super(context, DATABASE_NAME, null, DATABASE_VERSION);
    }

    /**
     * Init to Create table
     */
    @Override
    public void onCreate(SQLiteDatabase db) {

        db.execSQL(UserModel.createTable());
        db.execSQL(ReturnMagazineModel.createTable());
        db.execSQL(MaxYearRankModel.createTable());
    }

    /**
     * Upgrade if table is exist
     */
    @Override
    public void onUpgrade(SQLiteDatabase db, int oldVersion, int newVersion) {

        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_USER));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_RETURN_MAGAZINE));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_MAX_YEAR_RANK));

        onCreate(db);
    }

    /**
     * Delete data in all tables
     */
    public void clearTables() {

        //SQLite database
        SQLiteDatabase db = this.getWritableDatabase();
        db.delete(TABLE_USER, null, null);
        db.delete(TABLE_RETURN_MAGAZINE, null, null);
        db.delete(TABLE_MAX_YEAR_RANK, null, null);
    }

}
