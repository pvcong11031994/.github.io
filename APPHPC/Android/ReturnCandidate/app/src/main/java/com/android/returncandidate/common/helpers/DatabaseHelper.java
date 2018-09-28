package com.android.returncandidate.common.helpers;

import android.content.*;
import android.database.sqlite.*;

import com.android.returncandidate.db.models.*;

import static com.android.returncandidate.common.constants.Constants.*;

/**
 * Connect to Database SQLite
 *
 * @author tien-lv
 * @since 2017-11-30
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
//        db.execSQL(CLPModel.createLocationsTable());
        db.execSQL(CLPModel.createLagerClassificationsTable());
        db.execSQL(CLPModel.createPublisherTable());
//        db.execSQL(ReturnbookModel.createTable());
        db.execSQL(BookModel.createTable());
        db.execSQL(MaxYearRankModel.createTable());
    }

    /**
     * Upgrade if table is exist
     */
    @Override
    public void onUpgrade(SQLiteDatabase db, int oldVersion, int newVersion) {

        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_USER));
//        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_LOCATIONS));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_LARGE_CLASSIFICATIONS));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_PUBLISHERS));
//        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_RETURN_BOOKS));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_BOOKS));
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
        // db.delete(TABLE_LOCATIONS, null, null);
        db.delete(TABLE_LARGE_CLASSIFICATIONS, null, null);
        db.delete(TABLE_PUBLISHERS, null, null);
        db.delete(TABLE_BOOKS, null, null);
        db.delete(TABLE_MAX_YEAR_RANK, null, null);
    }

}
