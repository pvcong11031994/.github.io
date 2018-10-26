package com.fjn.magazinereturncandidate.db.models;

import android.annotation.SuppressLint;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteStatement;


import com.fjn.magazinereturncandidate.common.constants.Constants;
import com.fjn.magazinereturncandidate.common.utils.DatabaseManagerCommon;
import com.fjn.magazinereturncandidate.db.entity.ReturnMagazineEntity;

import java.util.List;

import static com.fjn.magazinereturncandidate.common.constants.Constants.ALIAS_COUNT;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_FIRST_SUPPLY_DATE;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_FLAG_SALES;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_GOODS_NAME;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_JAN_CD;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_JOUBI;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_LAST_ORDER_DATE;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_LAST_SALES_DATE;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_LAST_SUPPLY_DATE;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_LOCATION_ID;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_MEDIA_GROUP1_CD;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_MEDIA_GROUP1_NAME;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_MEDIA_GROUP2_CD;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_MEDIA_GROUP2_NAME;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_PERCENT;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_PRICE;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_PUBLISHER_CD;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_PUBLISHER_NAME;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_SALES_DATE;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_STOCK_COUNT;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_TOTAL_RETURN;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_TOTAL_SALES;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_TOTAL_SUPPLY;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_TRN_DATE;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_WRITER_NAME;
import static com.fjn.magazinereturncandidate.common.constants.Constants.COLUMN_YEAR_RANK;
import static com.fjn.magazinereturncandidate.common.constants.Constants.TABLE_RETURN_MAGAZINE;


/**
 * Return Magazine model
 *
 * @author cong-pv
 * @since 2018-10-18
 */

public class ReturnMagazineModel {

    /**
     * SQLite Database
     */
    private SQLiteDatabase db;

    private SQLiteStatement stmt;

    /**
     * Constructor Model Book
     */
    public ReturnMagazineModel() {

    }

    /**
     * Constructor Model Book
     */
    public ReturnMagazineModel(boolean isInsert, SQLiteDatabase db) {
        if (isInsert) {
            this.db = db;
            stmt = db.compileStatement(getSqlInsert());
        }
    }

    private static String getSqlInsert() {
        return String.format(
                "INSERT INTO %s (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) "
                        + "VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
                TABLE_RETURN_MAGAZINE, COLUMN_JAN_CD, COLUMN_STOCK_COUNT, COLUMN_GOODS_NAME,
                COLUMN_WRITER_NAME, COLUMN_PUBLISHER_CD, COLUMN_PUBLISHER_NAME, COLUMN_PRICE,
                COLUMN_FIRST_SUPPLY_DATE, COLUMN_LAST_SUPPLY_DATE, COLUMN_LAST_SALES_DATE,
                COLUMN_LAST_ORDER_DATE, COLUMN_MEDIA_GROUP1_CD, COLUMN_MEDIA_GROUP1_NAME,
                COLUMN_MEDIA_GROUP2_CD, COLUMN_MEDIA_GROUP2_NAME, COLUMN_SALES_DATE,
                COLUMN_TRN_DATE, COLUMN_YEAR_RANK, COLUMN_TOTAL_SALES, COLUMN_TOTAL_SUPPLY,
                COLUMN_TOTAL_RETURN, COLUMN_LOCATION_ID);
    }

    /**
     * Function create table books
     */
    public static String createTable() {

        return String.format(
                "CREATE TABLE %s(%s TEXT,%s INTEGER,%s TEXT,%s TEXT, %s TEXT, %s TEXT, %s FLOAT, %s TEXT, " +
                        "%s TEXT, %s TEXT,%s TEXT, %s TEXT, %s TEXT, %s TEXT, %s TEXT, %s TEXT, %s TEXT," +
                        " %s INTEGER, %s INTEGER, %s INTEGER, %s INTEGER, %s TEXT)",
                TABLE_RETURN_MAGAZINE, COLUMN_JAN_CD, COLUMN_STOCK_COUNT, COLUMN_GOODS_NAME,
                COLUMN_WRITER_NAME, COLUMN_PUBLISHER_CD, COLUMN_PUBLISHER_NAME, COLUMN_PRICE,
                COLUMN_FIRST_SUPPLY_DATE, COLUMN_LAST_SUPPLY_DATE, COLUMN_LAST_SALES_DATE,
                COLUMN_LAST_ORDER_DATE, COLUMN_MEDIA_GROUP1_CD, COLUMN_MEDIA_GROUP1_NAME,
                COLUMN_MEDIA_GROUP2_CD, COLUMN_MEDIA_GROUP2_NAME, COLUMN_SALES_DATE,
                COLUMN_TRN_DATE, COLUMN_YEAR_RANK, COLUMN_TOTAL_SALES, COLUMN_TOTAL_SUPPLY,
                COLUMN_TOTAL_RETURN, COLUMN_LOCATION_ID);
    }

    /**
     * Function insert into table books
     */
    public void insertData(SQLiteDatabase db, int indexListString, List<String> listValue) {

        StringBuilder valuesBuilder = new StringBuilder();
        for (int i = 0; i < indexListString; i += Constants.VALUE_COUNT_COLUMN_TABLE_RETURN_MAGAZINE_INSERT) {
            if (i != 0) {
                valuesBuilder.append(", ");
            }
            valuesBuilder.append("(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)");
        }
        stmt = db.compileStatement(String.format("INSERT INTO %s (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) VALUES ",
                TABLE_RETURN_MAGAZINE, COLUMN_JAN_CD, COLUMN_STOCK_COUNT, COLUMN_GOODS_NAME,
                COLUMN_WRITER_NAME, COLUMN_PUBLISHER_CD, COLUMN_PUBLISHER_NAME, COLUMN_PRICE,
                COLUMN_FIRST_SUPPLY_DATE, COLUMN_LAST_SUPPLY_DATE, COLUMN_LAST_SALES_DATE,
                COLUMN_LAST_ORDER_DATE, COLUMN_MEDIA_GROUP1_CD, COLUMN_MEDIA_GROUP1_NAME,
                COLUMN_MEDIA_GROUP2_CD, COLUMN_MEDIA_GROUP2_NAME, COLUMN_SALES_DATE,
                COLUMN_TRN_DATE, COLUMN_YEAR_RANK, COLUMN_TOTAL_SALES, COLUMN_TOTAL_SUPPLY,
                COLUMN_TOTAL_RETURN, COLUMN_LOCATION_ID) + valuesBuilder.toString());
        for (int i = 0; i < indexListString; i += Constants.VALUE_COUNT_COLUMN_TABLE_RETURN_MAGAZINE_INSERT) {
            stmt.bindString(i + 1, listValue.get(i));
            stmt.bindString(i + 2, listValue.get(i + 1));
            stmt.bindString(i + 3, listValue.get(i + 2));
            stmt.bindString(i + 4, listValue.get(i + 3));
            stmt.bindString(i + 5, listValue.get(i + 4));
            stmt.bindString(i + 6, listValue.get(i + 5));
            stmt.bindString(i + 7, listValue.get(i + 6));
            stmt.bindString(i + 8, listValue.get(i + 7));
            stmt.bindString(i + 9, listValue.get(i + 8));
            stmt.bindString(i + 10, listValue.get(i + 9));
            stmt.bindString(i + 11, listValue.get(i + 10));
            stmt.bindString(i + 12, listValue.get(i + 11));
            stmt.bindString(i + 13, listValue.get(i + 12));
            stmt.bindString(i + 14, listValue.get(i + 13));
            stmt.bindString(i + 15, listValue.get(i + 14));
            stmt.bindString(i + 16, listValue.get(i + 15));
            stmt.bindString(i + 17, listValue.get(i + 16));
            stmt.bindString(i + 18, listValue.get(i + 17));
            stmt.bindString(i + 19, listValue.get(i + 18));
            stmt.bindString(i + 20, listValue.get(i + 19));
            stmt.bindString(i + 21, listValue.get(i + 20));
            stmt.bindString(i + 22, listValue.get(i + 21));
        }
        stmt.executeInsert();
        stmt.clearBindings();
    }


    /**
     * Get Item Return Magazine
     */
    public ReturnMagazineEntity getItemReturnMagazine(String jan_code) {

        ReturnMagazineEntity returnMagazineEntity = new ReturnMagazineEntity();

        db = DatabaseManagerCommon.getInstance().openDatabase();
        String selectQuery = String.format("SELECT * FROM %s WHERE %s = '%s'", TABLE_RETURN_MAGAZINE, COLUMN_JAN_CD, jan_code);
        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                do {
                    returnMagazineEntity.setJan_cd(cursor.getString(cursor.getColumnIndex(COLUMN_JAN_CD)));
                    returnMagazineEntity.setBqsc_stock_count(cursor.getInt(cursor.getColumnIndex(COLUMN_STOCK_COUNT)));
                    returnMagazineEntity.setBqgm_goods_name(cursor.getString(cursor.getColumnIndex(COLUMN_GOODS_NAME)));
                    returnMagazineEntity.setBqgm_writer_name(cursor.getString(cursor.getColumnIndex(COLUMN_WRITER_NAME)));
                    returnMagazineEntity.setBqgm_publisher_cd(cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_CD)));
                    returnMagazineEntity.setBqgm_publisher_name(cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_NAME)));
                    returnMagazineEntity.setBqgm_price(cursor.getFloat(cursor.getColumnIndex(COLUMN_PRICE)));
                    returnMagazineEntity.setBqtse_first_supply_date(cursor.getString(cursor.getColumnIndex(COLUMN_FIRST_SUPPLY_DATE)));
                    returnMagazineEntity.setBqtse_last_supply_date(cursor.getString(cursor.getColumnIndex(COLUMN_LAST_SUPPLY_DATE)));
                    returnMagazineEntity.setBqtse_last_sale_date(cursor.getString(cursor.getColumnIndex(COLUMN_LAST_SALES_DATE)));
                    returnMagazineEntity.setBqtse_last_order_date(cursor.getString(cursor.getColumnIndex(COLUMN_LAST_ORDER_DATE)));
                    returnMagazineEntity.setBqct_media_group1_cd(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP1_CD)));
                    returnMagazineEntity.setBqct_media_group1_name(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP1_NAME)));
                    returnMagazineEntity.setBqct_media_group2_cd(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP2_CD)));
                    returnMagazineEntity.setBqct_media_group2_name(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP2_NAME)));
                    returnMagazineEntity.setBqgm_sales_date(cursor.getString(cursor.getColumnIndex(COLUMN_SALES_DATE)));
                    returnMagazineEntity.setBqio_trn_date(cursor.getString(cursor.getColumnIndex(COLUMN_TRN_DATE)));
                    returnMagazineEntity.setYear_rank(cursor.getInt(cursor.getColumnIndex(COLUMN_YEAR_RANK)));
                    returnMagazineEntity.setSts_total_sales(cursor.getInt(cursor.getColumnIndex(COLUMN_TOTAL_SALES)));
                    returnMagazineEntity.setSts_total_supply(cursor.getInt(cursor.getColumnIndex(COLUMN_TOTAL_SUPPLY)));
                    returnMagazineEntity.setSts_total_return(cursor.getInt(cursor.getColumnIndex(COLUMN_TOTAL_RETURN)));
                    returnMagazineEntity.setLocation_id(cursor.getString(cursor.getColumnIndex(COLUMN_LOCATION_ID)));
                } while (cursor.moveToNext());
            }
        } else {
            DatabaseManagerCommon.getInstance().closeDatabase();
            return null;
        }
        DatabaseManagerCommon.getInstance().closeDatabase();
        return returnMagazineEntity;
    }


    /**
     * Check Data
     */

    public boolean checkData() {

        boolean bool;
        db = DatabaseManagerCommon.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_RETURN_MAGAZINE);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        bool = cursor.moveToFirst() && cursor.getInt(cursor.getColumnIndex(ALIAS_COUNT)) > 0;

        DatabaseManagerCommon.getInstance().closeDatabase();

        return bool;
    }


    /**
     * Count data
     */
    public int countReturnMagazine() {

        int count = 0;
        db = DatabaseManagerCommon.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_RETURN_MAGAZINE);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                count = cursor.getInt(cursor.getColumnIndex(ALIAS_COUNT));
            }
        } else {
            DatabaseManagerCommon.getInstance().closeDatabase();
            return 0;
        }
        DatabaseManagerCommon.getInstance().closeDatabase();
        return count;
    }
}
