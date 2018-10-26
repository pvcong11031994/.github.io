package com.fjn.magazinereturncandidate.db.models;

import android.annotation.SuppressLint;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteStatement;

import com.fjn.magazinereturncandidate.common.utils.DatabaseManagerCommon;
import com.fjn.magazinereturncandidate.db.entity.MaxYearRankEntity;

import static com.fjn.magazinereturncandidate.common.constants.Constants.ALIAS_COUNT;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_MAX_YEAR_RANK;
import static com.fjn.magazinereturncandidate.common.constants.Constants.TABLE_MAX_YEAR_RANK;


/**
 *
 */
public class MaxYearRankModel {

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

    public void insertBulk(MaxYearRankEntity maxYearRankEntity) {
        stmt.bindString(1, String.valueOf(maxYearRankEntity.getMaxYearRank()));
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
        db = DatabaseManagerCommon.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_MAX_YEAR_RANK);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        bool = cursor.moveToFirst() && cursor.getInt(cursor.getColumnIndex(ALIAS_COUNT)) > 0;

        DatabaseManagerCommon.getInstance().closeDatabase();

        return bool;
    }

    /**
     * Get year rank
     */

    public MaxYearRankEntity getMaxYearRank() {

        MaxYearRankEntity maxYearRank = new MaxYearRankEntity();

        db = DatabaseManagerCommon.getInstance().openDatabase();
        String selectQuery = String.format("SELECT %s FROM %s", COLUMN_MAX_YEAR_RANK, TABLE_MAX_YEAR_RANK);
        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                do {
                    maxYearRank.setMaxYearRank(cursor.getInt(cursor.getColumnIndex(COLUMN_MAX_YEAR_RANK)));
                } while (cursor.moveToNext());
            }
        } else {
            DatabaseManagerCommon.getInstance().closeDatabase();
            return null;
        }
        DatabaseManagerCommon.getInstance().closeDatabase();
        return maxYearRank;
    }
}
