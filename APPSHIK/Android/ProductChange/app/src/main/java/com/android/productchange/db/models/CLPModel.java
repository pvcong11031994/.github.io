package com.android.productchange.db.models;

import static com.android.productchange.common.constants.Constants.*;

import android.annotation.SuppressLint;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteStatement;

import com.android.productchange.api.Config;
import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.utils.DatabaseManager;
import com.android.productchange.db.entity.CLP;

import java.util.List;

/**
 * <h1>Model of Large_classifications, Publishers</h1>
 *
 * @author tien-lv
 * @since 2017-12-25
 */
public class CLPModel {

    /**
     * SQLite Database
     */
    private SQLiteDatabase db;

    /**
     * SQLite Statement
     */
    private SQLiteStatement stmtClassify, stmtPublisher;

    /**
     * Constructor Model CLP
     */
    public CLPModel() {

    }

    /**
     * Constructor Model CLP
     *
     * @param db       SQLite database
     * @param isInsert check is insert
     */
    public CLPModel(boolean isInsert, SQLiteDatabase db) {
        if (isInsert) {
            this.db = db;
            stmtClassify = db.compileStatement(getSqlInsertClassify());
            stmtPublisher = db.compileStatement(getSqlInsertPublisher());
        }
    }

    /**
     * Function create table locations
     *
     * @return String
     */
    public static String createLocationsTable() {

        return String.format("CREATE TABLE %s(%s TEXT PRIMARY KEY, %s TEXT)",
                TABLE_LOCATIONS, COLUMN_ID, COLUMN_NAME);
    }

    /**
     * Function create table large classification
     *
     * @return String
     */
    public static String createLagreClassificationsTable() {

        return String.format("CREATE TABLE %s(%s TEXT PRIMARY KEY, %s TEXT)",
                TABLE_LARGE_CLASSIFICATIONS, COLUMN_ID, COLUMN_NAME);
    }

    /**
     * Function create table publishers
     *
     * @return String
     */
    public static String createPublishersTable() {

        return String.format("CREATE TABLE %s(%s TEXT PRIMARY KEY, %s TEXT)",
                TABLE_PUBLISHERS, COLUMN_ID, COLUMN_NAME);
    }

    /**
     * String sql insert statement of Classify
     *
     * @return String
     */
    private static String getSqlInsertClassify() {
        return String.format(
                "INSERT INTO %s (%s,%s) VALUES (?,?)",
                TABLE_LARGE_CLASSIFICATIONS, COLUMN_ID, COLUMN_NAME);
    }

    /**
     * String sql insert statement of Classify
     *
     * @return String
     */
    private static String getSqlInsertPublisher() {
        return String.format(
                "INSERT INTO %s (%s,%s) VALUES (?,?)",
                TABLE_PUBLISHERS, COLUMN_ID, COLUMN_NAME);
    }

    /**
     * Function insert into table classify
     */
    public void insertData(SQLiteDatabase db, int indexListString, List<String> listValue) {

        StringBuilder valuesBuilder = new StringBuilder();
        for (int i = 0; i < indexListString; i += Constants.VALUE_COUNT_COLUMN_TABLE_CLASSIFY_INSERT) {
            if (i != 0) {
                valuesBuilder.append(", ");
            }
            valuesBuilder.append("(?,?)");
        }
        stmtClassify = db.compileStatement(String.format("INSERT INTO %s (%s,%s) VALUES ",
                TABLE_LARGE_CLASSIFICATIONS, COLUMN_ID, COLUMN_NAME) + valuesBuilder.toString());
        for (int i = 0; i < indexListString; i += Constants.VALUE_COUNT_COLUMN_TABLE_CLASSIFY_INSERT) {
            stmtClassify.bindString(i + 1, listValue.get(i));
            stmtClassify.bindString(i + 2, listValue.get(i + 1));
        }
        stmtClassify.executeInsert();
        stmtClassify.clearBindings();

    }
    /**
     * Function insert into table publisher
     */
    public void insertDataPublisher(SQLiteDatabase db, int indexListString, List<String> listValue) {

        StringBuilder valuesBuilder = new StringBuilder();
        for (int i = 0; i < indexListString; i += Constants.VALUE_COUNT_COLUMN_TABLE_PUBLISHER_INSERT) {
            if (i != 0) {
                valuesBuilder.append(", ");
            }
            valuesBuilder.append("(?,?)");
        }
        stmtPublisher = db.compileStatement(String.format("INSERT INTO %s (%s,%s) VALUES ",
                TABLE_PUBLISHERS, COLUMN_ID, COLUMN_NAME) + valuesBuilder.toString());
        for (int i = 0; i < indexListString; i += Constants.VALUE_COUNT_COLUMN_TABLE_PUBLISHER_INSERT) {
            stmtPublisher.bindString(i + 1, listValue.get(i));
            stmtPublisher.bindString(i + 2, listValue.get(i + 1));
        }
        stmtPublisher.executeInsert();
        stmtPublisher.clearBindings();

    }

    /**
     * Check data is exist
     *
     * @param type is Classify or Publisher
     * @return boolean
     */
    public boolean checkData(int type) {

        boolean bool;

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery;
        switch (type) {
            case Config.TYPE_CLASSIFY:
                selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                        TABLE_LARGE_CLASSIFICATIONS);
                break;
            case Config.TYPE_PUBLISHER:
                selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                        TABLE_PUBLISHERS);
                break;
            default:
                selectQuery = "";
                break;
        }
        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        bool = cursor.moveToFirst() && cursor.getInt(cursor.getColumnIndex(ALIAS_COUNT)) > 0;

        DatabaseManager.getInstance().closeDatabase();

        return bool;
    }

    /**
     * Count data in table Regular Books
     *
     * @return int
     */
    public int countBooks() {

        int count = 0;
        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_LARGE_CLASSIFICATIONS);

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
     * Count data in table Publisher Books
     *
     * @return int
     */
    public int countDataTablePublisher() {

        int count = 0;
        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_PUBLISHERS);

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
}
