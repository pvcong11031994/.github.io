package com.fjn.magazinereturncandidate.common.utils;

import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteOpenHelper;

/**
 * Database manager
 *
 * @author cong-pv
 * @since 2018-10-15
 */

public class DatabaseManagerCommon {

    /**
     * Open counter
     */
    private boolean mOpenCounter = false;

    /**
     * Instance of Database manager
     */
    private static DatabaseManagerCommon instance;
    /**
     * SQLite helper
     */
    private static SQLiteOpenHelper mDatabaseHelper;
    /**
     * SQLite database
     */
    private SQLiteDatabase mDatabase;

    /**
     * Init
     */
    public static synchronized void initializeInstance(SQLiteOpenHelper helper) {

        if (instance == null) {
            instance = new DatabaseManagerCommon();
            mDatabaseHelper = helper;
        }
    }

    /**
     * Get instance
     */
    public static synchronized DatabaseManagerCommon getInstance() {

        if (instance == null) {
            throw new IllegalStateException(DatabaseManagerCommon.class.getSimpleName() +
                    " is not initialized, call initializeInstance(..) method first.");
        }
        return instance;
    }

    /**
     * Open database
     */
    public synchronized SQLiteDatabase openDatabase() {

        if (!mOpenCounter) {
            // Opening new database
            mDatabase = mDatabaseHelper.getWritableDatabase();
            mOpenCounter = true;
        }
        return mDatabase;
    }

    /**
     * Close database
     */
    public synchronized void closeDatabase() {

        if (mOpenCounter) {
            // Closing database
            mDatabase.close();
            mOpenCounter = false;
        }
    }
}
