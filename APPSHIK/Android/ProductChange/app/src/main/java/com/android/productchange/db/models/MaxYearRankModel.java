package com.android.productchange.db.models;

import android.annotation.SuppressLint;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteStatement;

import com.android.productchange.common.utils.DatabaseManager;
import com.android.productchange.db.entity.MaxYearRank;

import static com.android.productchange.common.constants.Constants.*;

/**
 *
 */
public class MaxYearRankModel {

    private MaxYearRank maxYearRank;

    /**
     * SQLite Database
     */
    private SQLiteDatabase db;

    /**
     * SQLite Statement
     */
    private SQLiteStatement stmt;

    /**
     * Constructor Model Return Book
     */
    public MaxYearRankModel() {
        maxYearRank = new MaxYearRank();
    }

    /**
     * Constructor Model Book
     *
     * @param db       SQLite Database
     * @param isInsert check is insert
     */
    public MaxYearRankModel(boolean isInsert, SQLiteDatabase db) {
        if (isInsert) {
            this.db = db;
            stmt = db.compileStatement(getSqlInsert());
        }
    }

    /**
     * Function create table max year rank
     */
    public static String createTable() {

        return String.format(
                "CREATE TABLE %s(%s INTEGER)",
                TABLE_MAX_YEAR_RANK, COLUMN_MAX_YEAR_RANK);
    }


    /**
     * String sql insert statement
     *
     * @return String
     */
    private static String getSqlInsert() {
        return String.format(
                "INSERT INTO %s (%s) VALUES (?)",
                TABLE_MAX_YEAR_RANK, COLUMN_MAX_YEAR_RANK);
    }

    public void insertBulk(MaxYearRank maxYearRank) {
        stmt.bindString(1, String.valueOf(maxYearRank.getMaxYearRank()));
        stmt.execute();
        stmt.clearBindings();

    }

    /**
     * Check Data is exist
     *
     * @return boolean
     */

    public boolean checkData() {

        boolean bool;
        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_MAX_YEAR_RANK);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        bool = cursor.moveToFirst() && cursor.getInt(cursor.getColumnIndex(ALIAS_COUNT)) > 0;

        DatabaseManager.getInstance().closeDatabase();

        return bool;
    }

    public int countMaxYearRank() {

        int count = 0;
        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_MAX_YEAR_RANK);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                count = cursor.getInt(cursor.getColumnIndex(ALIAS_COUNT));
            }
        } else {
            DatabaseManager.getInstance().closeDatabase();
            return 0;
        }
        DatabaseManager.getInstance().closeDatabase();
        return count;
    }

    /**
     * Get year rank
     */

    public MaxYearRank getMaxYearRank() {

        MaxYearRank maxYearRank = new MaxYearRank();

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT %s FROM %s", COLUMN_MAX_YEAR_RANK, TABLE_MAX_YEAR_RANK);
        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                do {
                    maxYearRank.setMaxYearRank(cursor.getInt(cursor.getColumnIndex(COLUMN_MAX_YEAR_RANK)));
                } while (cursor.moveToNext());
            }
        } else {
            DatabaseManager.getInstance().closeDatabase();
            return null;
        }
        DatabaseManager.getInstance().closeDatabase();
        return maxYearRank;
    }
}
