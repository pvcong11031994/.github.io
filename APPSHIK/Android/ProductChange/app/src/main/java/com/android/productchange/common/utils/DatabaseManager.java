package com.android.productchange.common.utils;

import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteOpenHelper;

/**
 * <h1>Database Manager</h1>
 *
 * @author tien-lv
 * @since 2017-12-06
 */

public class DatabaseManager {

    /**
     * Open counter
     */
    private boolean mOpenCounter = false;

    /**
     * Instance of Database manager
     */
    private static DatabaseManager instance;
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
     *
     * @param helper SQLite open helper
     */
    public static synchronized void initializeInstance(SQLiteOpenHelper helper) {

        if (instance == null) {
            instance = new DatabaseManager();
            mDatabaseHelper = helper;
        }
    }

    /**
     * Get instance
     */
    public static synchronized DatabaseManager getInstance() {

        if (instance == null) {
            throw new IllegalStateException(DatabaseManager.class.getSimpleName() +
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
