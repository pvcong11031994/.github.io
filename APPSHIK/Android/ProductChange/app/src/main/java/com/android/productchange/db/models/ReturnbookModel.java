package com.android.productchange.db.models;

import static com.android.productchange.common.constants.Constants.*;

import android.annotation.SuppressLint;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteStatement;

import com.android.productchange.api.Config;
import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.utils.Common;
import com.android.productchange.common.utils.ConditionQueryCommon;
import com.android.productchange.common.utils.DatabaseManager;
import com.android.productchange.common.utils.FlagSettingNew;
import com.android.productchange.db.entity.Books;
import com.android.productchange.db.entity.CLP;
import com.android.productchange.db.entity.Returnbooks;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * <h1>Model Return Book</h1>
 *
 * @author cong-pv
 * @since 2018-08-31
 */
public class ReturnbookModel {

    /**
     * Entity Return Books
     */
    private Returnbooks returnbooks;

    /**
     * SQLite Database
     */
    private SQLiteDatabase db;

    /**
     * SQLite Statement
     */
    private SQLiteStatement stmt;
    private Common common = new Common();
    private ConditionQueryCommon conditionQueryCommon = new ConditionQueryCommon();

    /**
     * Constructor Model Return Book
     */
    public ReturnbookModel() {
        returnbooks = new Returnbooks();
    }

    /**
     * Constructor Model Book
     *
     * @param db       SQLite Database
     * @param isInsert check is insert
     */
    public ReturnbookModel(boolean isInsert, SQLiteDatabase db) {
        if (isInsert) {
            this.db = db;
            stmt = db.compileStatement(getSqlInsert());
        }
    }

    private static String getSqlInsert() {
        return String.format(
                "INSERT INTO %s (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) "
                        + "VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
                TABLE_RETURN_BOOKS, COLUMN_JAN_CD, COLUMN_STOCK_COUNT, COLUMN_GOODS_NAME,
                COLUMN_WRITER_NAME, COLUMN_PUBLISHER_CD, COLUMN_PUBLISHER_NAME, COLUMN_PRICE,
                COLUMN_FIRST_SUPPLY_DATE, COLUMN_LAST_SUPPLY_DATE, COLUMN_LAST_SALES_DATE,
                COLUMN_LAST_ORDER_DATE, COLUMN_MEDIA_GROUP1_CD, COLUMN_MEDIA_GROUP1_NAME,
                COLUMN_MEDIA_GROUP2_CD, COLUMN_MEDIA_GROUP2_NAME, COLUMN_SALES_DATE,
                COLUMN_TRN_DATE, COLUMN_PERCENT, COLUMN_FLAG_SALES, COLUMN_YEAR_RANK, COLUMN_JOUBI,
                COLUMN_TOTAL_SALES, COLUMN_TOTAL_SUPPLY, COLUMN_TOTAL_RETURN, COLUMN_LOCATION_ID);
    }

    /**
     * Function create table books
     */
    public static String createTable() {

        return String.format(
                "CREATE TABLE %s(%s TEXT,%s INTEGER,%s TEXT,%s TEXT, %s TEXT, %s TEXT, %s FLOAT, %s TEXT, " +
                        "%s TEXT, %s TEXT,%s TEXT, %s TEXT, %s TEXT, %s TEXT, %s TEXT, %s TEXT, %s TEXT," +
                        " %s FLOAT, %s TEXT, %s INTEGER, %s INTEGER, %s INTEGER, %s INTEGER, %s INTEGER, %s TEXT)",
                TABLE_RETURN_BOOKS, COLUMN_JAN_CD, COLUMN_STOCK_COUNT, COLUMN_GOODS_NAME,
                COLUMN_WRITER_NAME, COLUMN_PUBLISHER_CD, COLUMN_PUBLISHER_NAME, COLUMN_PRICE,
                COLUMN_FIRST_SUPPLY_DATE, COLUMN_LAST_SUPPLY_DATE, COLUMN_LAST_SALES_DATE,
                COLUMN_LAST_ORDER_DATE, COLUMN_MEDIA_GROUP1_CD, COLUMN_MEDIA_GROUP1_NAME,
                COLUMN_MEDIA_GROUP2_CD, COLUMN_MEDIA_GROUP2_NAME, COLUMN_SALES_DATE,
                COLUMN_TRN_DATE, COLUMN_PERCENT, COLUMN_FLAG_SALES, COLUMN_YEAR_RANK, COLUMN_JOUBI,
                COLUMN_TOTAL_SALES, COLUMN_TOTAL_SUPPLY, COLUMN_TOTAL_RETURN, COLUMN_LOCATION_ID);
    }

    /**
     * Function create table return books temp
     */
    public static String createTableTemp() {

        return String.format(
                "CREATE TABLE %s(%s TEXT,%s INTEGER,%s TEXT,%s TEXT, %s TEXT, %s TEXT, %s FLOAT, %s TEXT, " +
                        "%s TEXT, %s TEXT,%s TEXT, %s TEXT, %s TEXT, %s TEXT, %s TEXT, %s TEXT, %s TEXT," +
                        " %s FLOAT, %s TEXT, %s INTEGER, %s INTEGER, %s INTEGER, %s INTEGER, %s INTEGER, %s TEXT)",
                TABLE_RETURN_BOOKS_TEMP, COLUMN_JAN_CD, COLUMN_STOCK_COUNT, COLUMN_GOODS_NAME,
                COLUMN_WRITER_NAME, COLUMN_PUBLISHER_CD, COLUMN_PUBLISHER_NAME, COLUMN_PRICE,
                COLUMN_FIRST_SUPPLY_DATE, COLUMN_LAST_SUPPLY_DATE, COLUMN_LAST_SALES_DATE,
                COLUMN_LAST_ORDER_DATE, COLUMN_MEDIA_GROUP1_CD, COLUMN_MEDIA_GROUP1_NAME,
                COLUMN_MEDIA_GROUP2_CD, COLUMN_MEDIA_GROUP2_NAME, COLUMN_SALES_DATE,
                COLUMN_TRN_DATE, COLUMN_PERCENT, COLUMN_FLAG_SALES, COLUMN_YEAR_RANK, COLUMN_JOUBI,
                COLUMN_TOTAL_SALES, COLUMN_TOTAL_SUPPLY, COLUMN_TOTAL_RETURN, COLUMN_LOCATION_ID);
    }

    /**
     * Function insert into table books
     */
    public void insertData(SQLiteDatabase db, int indexListString, List<String> listValue) {

        StringBuilder valuesBuilder = new StringBuilder();
        for (int i = 0; i < indexListString; i += Constants.VALUE_COUNT_COLUMN_TABLE_RETURN_BOOK_INSERT) {
            if (i != 0) {
                valuesBuilder.append(", ");
            }
            valuesBuilder.append("(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)");
        }
        stmt = db.compileStatement(String.format("INSERT INTO %s (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) VALUES ",
                TABLE_RETURN_BOOKS, COLUMN_JAN_CD, COLUMN_STOCK_COUNT, COLUMN_GOODS_NAME,
                COLUMN_WRITER_NAME, COLUMN_PUBLISHER_CD, COLUMN_PUBLISHER_NAME, COLUMN_PRICE,
                COLUMN_FIRST_SUPPLY_DATE, COLUMN_LAST_SUPPLY_DATE, COLUMN_LAST_SALES_DATE,
                COLUMN_LAST_ORDER_DATE, COLUMN_MEDIA_GROUP1_CD, COLUMN_MEDIA_GROUP1_NAME,
                COLUMN_MEDIA_GROUP2_CD, COLUMN_MEDIA_GROUP2_NAME, COLUMN_SALES_DATE,
                COLUMN_TRN_DATE, COLUMN_PERCENT, COLUMN_FLAG_SALES, COLUMN_YEAR_RANK, COLUMN_JOUBI,
                COLUMN_TOTAL_SALES, COLUMN_TOTAL_SUPPLY, COLUMN_TOTAL_RETURN, COLUMN_LOCATION_ID) + valuesBuilder.toString());
        for (int i = 0; i < indexListString; i += Constants.VALUE_COUNT_COLUMN_TABLE_RETURN_BOOK_INSERT) {
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
            stmt.bindString(i + 23, listValue.get(i + 22));
            stmt.bindString(i + 24, listValue.get(i + 23));
            stmt.bindString(i + 25, listValue.get(i + 24));
        }
        stmt.executeInsert();
        stmt.clearBindings();

    }

    /**
     * Get list info return books when select group1 cd
     */
    public List<Returnbooks> getListBookInfoSelectGroup1Cd(int offset, Map<Integer, String> mapOrder) {

        List<Returnbooks> returnbooksList = new ArrayList<>();

        db = DatabaseManager.getInstance().openDatabase();

        String selectQuery = String.format("SELECT * FROM %s ", TABLE_RETURN_BOOKS_TEMP);

        String selectQueryOrder = ORDER_BY;
        int count = 0;

        if (mapOrder != null) {
            for (Integer key : mapOrder.keySet()) {
                if (!mapOrder.get(key).isEmpty()) {
                    switch (key) {
                        case Constants.NUMBER_1:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_LOCATION_ID, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_2:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_JAN_CD, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_3:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_GOODS_NAME, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_5:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_PUBLISHER_CD, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_6:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_MEDIA_GROUP1_CD, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_7:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_MEDIA_GROUP2_CD, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_8:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_SALES_DATE, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_9:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_STOCK_COUNT, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_10:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_YEAR_RANK, mapOrder.get(key));
                    }
                    count++;
                }
            }
        }
        if (ORDER_BY.equals(selectQueryOrder)) {
            selectQueryOrder = BLANK;
        } else {
            if (count > 0) {
                selectQueryOrder = selectQueryOrder.substring(0, selectQueryOrder.length() - 1);
            } else {
                selectQueryOrder += String.format(" (CASE" +
                        " WHEN %s = 99999999 THEN 0" +
                        " ELSE %s END) DESC ", Constants.COLUMN_YEAR_RANK, Constants.COLUMN_YEAR_RANK);
            }
        }
        String selectQueryLimit = String.format(" LIMIT 1000 OFFSET %s", offset);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(
                selectQuery + selectQueryOrder + selectQueryLimit, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                do {
                    Returnbooks returnbooks = new Returnbooks();
                    returnbooks.setLocation_id(
                            cursor.getString(cursor.getColumnIndex(COLUMN_LOCATION_ID)));
                    returnbooks.setJan_cd(
                            cursor.getString(cursor.getColumnIndex(COLUMN_JAN_CD)));
                    returnbooks.setBqsc_stock_count(
                            cursor.getInt(cursor.getColumnIndex(COLUMN_STOCK_COUNT)));
                    returnbooks.setBqgm_goods_name(
                            cursor.getString(cursor.getColumnIndex(COLUMN_GOODS_NAME)));
                    returnbooks.setBqgm_publisher_cd(
                            cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_CD)));
                    returnbooks.setBqgm_publisher_name(
                            cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_NAME)));
                    returnbooks.setBqct_media_group1_cd(
                            cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP1_CD)));
                    returnbooks.setBqct_media_group1_name(
                            cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP1_NAME)));
                    returnbooks.setBqct_media_group2_cd(
                            cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP2_CD)));
                    returnbooks.setBqct_media_group2_name(
                            cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP2_NAME)));
                    returnbooks.setBqgm_sales_date(
                            cursor.getString(cursor.getColumnIndex(COLUMN_SALES_DATE)));
                    returnbooks.setYear_rank(Integer.parseInt(
                            cursor.getString(cursor.getColumnIndex(COLUMN_YEAR_RANK))));
                    returnbooks.setBqgm_writer_name(
                            cursor.getString(cursor.getColumnIndex(COLUMN_WRITER_NAME)));
                    returnbooks.setBqgm_price(
                            cursor.getFloat(cursor.getColumnIndex(COLUMN_PRICE)));
                    returnbooks.setBqtse_first_supply_date(
                            cursor.getString(cursor.getColumnIndex(COLUMN_FIRST_SUPPLY_DATE)));
                    returnbooks.setBqtse_last_supply_date(
                            cursor.getString(cursor.getColumnIndex(COLUMN_LAST_SUPPLY_DATE)));
                    returnbooks.setBqtse_last_order_date(
                            cursor.getString(cursor.getColumnIndex(COLUMN_LAST_ORDER_DATE)));
                    returnbooks.setBqtse_last_sale_date(
                            cursor.getString(cursor.getColumnIndex(COLUMN_LAST_SALES_DATE)));
                    returnbooks.setSts_total_sales(
                            cursor.getInt(cursor.getColumnIndex(COLUMN_TOTAL_SALES)));
                    returnbooks.setSts_total_supply(
                            cursor.getInt(cursor.getColumnIndex(COLUMN_TOTAL_SUPPLY)));
                    returnbooks.setSts_total_return(
                            cursor.getInt(cursor.getColumnIndex(COLUMN_TOTAL_RETURN)));
                    returnbooksList.add(returnbooks);
                } while (cursor.moveToNext());
            }
        } else {
            DatabaseManager.getInstance().closeDatabase();
            return null;
        }
        DatabaseManager.getInstance().closeDatabase();
        return returnbooksList;
    }

    /**
     * Get book info filter in table return books
     *
     * @param flagSettingNew {@link FlagSettingNew}
     * @return list has entity is return books
     */
    public List<Returnbooks> getListBookInfo(int offset, Map<Integer, String> mapOrder, FlagSettingNew flagSettingNew) {

        List<Returnbooks> returnbooksList = new ArrayList<>();

        db = DatabaseManager.getInstance().openDatabase();

        String selectQuery = String.format("SELECT * FROM %s ", TABLE_RETURN_BOOKS);

        //Query condition filter
        String queryCondition = conditionQueryCommon.conditionFilterSetting(flagSettingNew);

        queryCondition += String.format(" AND percent >= %s ", Float.parseFloat(common.FormatPercent(flagSettingNew.getFlagStockPercent())));

        //Query condition filter select group cd
        String queryConditionGroupCd = conditionQueryCommon.conditionFilterSettingGroupCd(flagSettingNew);

        String selectQueryOrder = ORDER_BY;
        int count = 0;

        if (mapOrder != null) {
            for (Integer key : mapOrder.keySet()) {
                if (!mapOrder.get(key).isEmpty()) {
                    switch (key) {
                        case Constants.NUMBER_1:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_LOCATION_ID, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_2:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_JAN_CD, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_3:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_GOODS_NAME, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_5:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_PUBLISHER_CD, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_6:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_MEDIA_GROUP1_CD, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_7:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_MEDIA_GROUP2_CD, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_8:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_SALES_DATE, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_9:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_STOCK_COUNT, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_10:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_YEAR_RANK, mapOrder.get(key));
                    }
                    count++;
                }
            }
        }
        if (ORDER_BY.equals(selectQueryOrder)) {
            selectQueryOrder = BLANK;
        } else {
            if (count > 0) {
                selectQueryOrder = selectQueryOrder.substring(0, selectQueryOrder.length() - 1);
            } else {
                selectQueryOrder += String.format(" (CASE" +
                        " WHEN %s = 99999999 THEN 0" +
                        " ELSE %s END) DESC ", Constants.COLUMN_YEAR_RANK, Constants.COLUMN_YEAR_RANK);
            }
        }
        String selectQueryLimit = String.format(" LIMIT 1000 OFFSET %s", offset);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(
                selectQuery + queryCondition + queryConditionGroupCd + selectQueryOrder + selectQueryLimit, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                do {
                    Returnbooks returnbooks = new Returnbooks();
                    returnbooks.setLocation_id(
                            cursor.getString(cursor.getColumnIndex(COLUMN_LOCATION_ID)));
                    returnbooks.setJan_cd(
                            cursor.getString(cursor.getColumnIndex(COLUMN_JAN_CD)));
                    returnbooks.setBqsc_stock_count(
                            cursor.getInt(cursor.getColumnIndex(COLUMN_STOCK_COUNT)));
                    returnbooks.setBqgm_goods_name(
                            cursor.getString(cursor.getColumnIndex(COLUMN_GOODS_NAME)));
                    returnbooks.setBqgm_publisher_cd(
                            cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_CD)));
                    returnbooks.setBqgm_publisher_name(
                            cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_NAME)));
                    returnbooks.setBqct_media_group1_cd(
                            cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP1_CD)));
                    returnbooks.setBqct_media_group1_name(
                            cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP1_NAME)));
                    returnbooks.setBqct_media_group2_cd(
                            cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP2_CD)));
                    returnbooks.setBqct_media_group2_name(
                            cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP2_NAME)));
                    returnbooks.setBqgm_sales_date(
                            cursor.getString(cursor.getColumnIndex(COLUMN_SALES_DATE)));
                    returnbooks.setYear_rank(Integer.parseInt(
                            cursor.getString(cursor.getColumnIndex(COLUMN_YEAR_RANK))));
                    returnbooks.setBqgm_writer_name(
                            cursor.getString(cursor.getColumnIndex(COLUMN_WRITER_NAME)));
                    returnbooks.setBqgm_price(
                            cursor.getFloat(cursor.getColumnIndex(COLUMN_PRICE)));
                    returnbooks.setBqtse_first_supply_date(
                            cursor.getString(cursor.getColumnIndex(COLUMN_FIRST_SUPPLY_DATE)));
                    returnbooks.setBqtse_last_supply_date(
                            cursor.getString(cursor.getColumnIndex(COLUMN_LAST_SUPPLY_DATE)));
                    returnbooks.setBqtse_last_order_date(
                            cursor.getString(cursor.getColumnIndex(COLUMN_LAST_ORDER_DATE)));
                    returnbooks.setBqtse_last_sale_date(
                            cursor.getString(cursor.getColumnIndex(COLUMN_LAST_SALES_DATE)));
                    returnbooks.setSts_total_sales(
                            cursor.getInt(cursor.getColumnIndex(COLUMN_TOTAL_SALES)));
                    returnbooks.setSts_total_supply(
                            cursor.getInt(cursor.getColumnIndex(COLUMN_TOTAL_SUPPLY)));
                    returnbooks.setSts_total_return(
                            cursor.getInt(cursor.getColumnIndex(COLUMN_TOTAL_RETURN)));
                    returnbooksList.add(returnbooks);
                } while (cursor.moveToNext());
            }
        } else {
            DatabaseManager.getInstance().closeDatabase();
            return null;
        }
        DatabaseManager.getInstance().closeDatabase();
        return returnbooksList;
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
                TABLE_RETURN_BOOKS);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        bool = cursor.moveToFirst() && cursor.getInt(cursor.getColumnIndex(ALIAS_COUNT)) > 0;

        DatabaseManager.getInstance().closeDatabase();

        return bool;
    }

    /**
     * Get list info Classify or Publisher in table return books
     *
     * @param type is Classify or Publisher
     * @return list has entity is CLP
     */
    public List<CLP> getInfo(int type) {

        List<CLP> clpList = new ArrayList<>();

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery;
        switch (type) {
            case Config.TYPE_CLASSIFY:
                selectQuery = String.format(
                        "SELECT %s , %s FROM %s ", COLUMN_ID, COLUMN_NAME, TABLE_RETURN_BOOKS);
                break;
            case Config.TYPE_PUBLISHER:
                selectQuery = String.format(
                        "SELECT %s ,%s FROM %s ", COLUMN_PUBLISHER_CD, COLUMN_PUBLISHER_NAME_RETURN,
                        TABLE_RETURN_BOOKS);
                break;
            default:
                selectQuery = "";
                break;
        }

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                do {
                    CLP clp = new CLP();
                    clp.setId(cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_CD)));
                    clp.setName(cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_NAME_RETURN)));

                    clpList.add(clp);
                } while (cursor.moveToNext());
            }
        } else {
            DatabaseManager.getInstance().closeDatabase();
            return null;
        }
        DatabaseManager.getInstance().closeDatabase();
        return clpList;
    }

    /**
     * Count data in table Return Books
     *
     * @return int
     */
    public int countBooks() {

        int count = 0;
        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_RETURN_BOOKS);

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
     * String sql insert statement of Classify
     */
    public void getSqlInsertClassify() {
        SQLiteDatabase db = DatabaseManager.getInstance().openDatabase();
        String query = String.format(
                "INSERT INTO %s (%s, %s, %s, %s) SELECT %s, %s, %s, %s FROM %s GROUP BY %s, %s",
                TABLE_GENRE_RETURN_BOOK, COLUMN_MEDIA_GROUP1_CD, COLUMN_MEDIA_GROUP1_NAME,
                COLUMN_MEDIA_GROUP2_CD, COLUMN_MEDIA_GROUP2_NAME, COLUMN_MEDIA_GROUP1_CD,
                COLUMN_MEDIA_GROUP1_NAME, COLUMN_MEDIA_GROUP2_CD, COLUMN_MEDIA_GROUP2_NAME,
                TABLE_RETURN_BOOKS, COLUMN_MEDIA_GROUP1_CD, COLUMN_MEDIA_GROUP2_CD);
        db.execSQL(query);
        DatabaseManager.getInstance().closeDatabase();
    }

    /**
     * Get list info filter group 1 cd
     *
     * @return list is entity book
     */
    public List<CLP> getInfoGroupCd1() {

        List<CLP> clpList = new ArrayList<>();

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery;
        selectQuery = String.format(
                "SELECT %s , %s FROM %s GROUP BY %s", COLUMN_MEDIA_GROUP1_CD,
                COLUMN_MEDIA_GROUP1_NAME, TABLE_GENRE_RETURN_BOOK, COLUMN_MEDIA_GROUP1_CD);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                do {
                    CLP clp = new CLP();
                    clp.setId(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP1_CD)));
                    clp.setName(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP1_NAME)));

                    clpList.add(clp);
                } while (cursor.moveToNext());
            }
        } else {
            DatabaseManager.getInstance().closeDatabase();
            return null;
        }
        DatabaseManager.getInstance().closeDatabase();
        return clpList;
    }


    /**
     * Get list info filter group 1 cd
     *
     * @return list is entity book
     */
    public List<CLP> getInfoGroupCd2(String selectGroup1Cd) {

        List<CLP> clpList = new ArrayList<>();

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery;
        if (Constants.ID_ROW_ALL.equals(selectGroup1Cd)) {
            selectQuery = String.format(
                    "SELECT %s ,%s FROM %s GROUP BY %s", COLUMN_MEDIA_GROUP2_CD,
                    COLUMN_MEDIA_GROUP2_NAME, TABLE_GENRE_RETURN_BOOK, COLUMN_MEDIA_GROUP2_CD);
        } else {
            selectQuery = String.format(
                    "SELECT %s ,%s FROM %s WHERE %s = '%s' GROUP BY %s", COLUMN_MEDIA_GROUP2_CD,
                    COLUMN_MEDIA_GROUP2_NAME, TABLE_GENRE_RETURN_BOOK, COLUMN_MEDIA_GROUP1_CD, selectGroup1Cd, COLUMN_MEDIA_GROUP2_CD);
        }
        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                do {
                    CLP clp = new CLP();
                    clp.setId(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP2_CD)));
                    clp.setName(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP2_NAME)));

                    clpList.add(clp);
                } while (cursor.moveToNext());
            }
        } else {
            DatabaseManager.getInstance().closeDatabase();
            return null;
        }
        DatabaseManager.getInstance().closeDatabase();
        return clpList;
    }

    /**
     * Get sum and stock when select all group cd
     *
     * @param flagSettingNew
     * @return
     */

    public Returnbooks getSumStockAndCountJanIsNotSelectGroup(FlagSettingNew flagSettingNew) {

        Returnbooks singleBooks = new Returnbooks();
        Common common = new Common();
        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT SUM(bqsc_stock_count) bqsc_stock_count, " +
                "COUNT(jan_cd) jan_cd FROM %s ", TABLE_RETURN_BOOKS);

        //Query condition filter
        String queryCondition = conditionQueryCommon.conditionFilterSetting(flagSettingNew);
        queryCondition += String.format(" AND percent >= %s ", Float.parseFloat(common.FormatPercent(flagSettingNew.getFlagStockPercent())));

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery + queryCondition, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                do {
                    singleBooks.setCountJan_Cd(cursor.getInt(cursor.getColumnIndex(COLUMN_JAN_CD)));
                    singleBooks.setSumStocks(cursor.getInt(cursor.getColumnIndex(COLUMN_STOCK_COUNT)));
                } while (cursor.moveToNext());
            }
        } else {
            DatabaseManager.getInstance().closeDatabase();
            return null;
        }
        DatabaseManager.getInstance().closeDatabase();
        return singleBooks;

    }

    /**
     * Insert data filter table temp
     */

    public void insertDataTableFilter(FlagSettingNew flagSettingNew) {

        int countRecordFilter = countBooksGroupCd(flagSettingNew);

        //Select database
        db = DatabaseManager.getInstance().openDatabase();

        //Delete data table temp return books
        String sqlDeleteTable = String.format("DELETE FROM %s", TABLE_RETURN_BOOKS_TEMP);
        db.execSQL(sqlDeleteTable);

        if (countRecordFilter > 0) {
            float intStockPercent = Float.parseFloat(common.FormatPercentLocal(flagSettingNew.getFlagStockPercent()));
            int recordLimit = (int) Math.ceil(countRecordFilter * intStockPercent);

            //Query condition filter
            String queryCondition = conditionQueryCommon.conditionFilterSetting(flagSettingNew);
            //Query condition filter select group cd
            String queryConditionGroupCd = conditionQueryCommon.conditionFilterSettingGroupCd(flagSettingNew);

            String selectQuery = String.format("INSERT INTO %s SELECT * FROM (SELECT * FROM %s WHERE 1 = 1  %s ORDER BY %s DESC LIMIT %s) ",
                    TABLE_RETURN_BOOKS_TEMP, TABLE_RETURN_BOOKS, queryConditionGroupCd, COLUMN_YEAR_RANK, recordLimit);
            db.execSQL(selectQuery + queryCondition);
        }
        DatabaseManager.getInstance().closeDatabase();
    }

    /**
     * Get sum and stock when click filter group cd not all + click filter not click setting
     *
     * @return
     */
    public Returnbooks getDataSelectGroupCdCountSum() {

        Returnbooks singleBooks = new Returnbooks();

        singleBooks.setCountJan_Cd(0);
        singleBooks.setSumStocks(0);

        //Select database
        db = DatabaseManager.getInstance().openDatabase();

        String selectQuery = String.format("SELECT SUM(%s) %s, COUNT(%s) %s FROM %s",
                COLUMN_STOCK_COUNT, COLUMN_STOCK_COUNT, COLUMN_JAN_CD, COLUMN_JAN_CD, TABLE_RETURN_BOOKS_TEMP);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                do {
                    singleBooks.setCountJan_Cd(cursor.getInt(cursor.getColumnIndex(COLUMN_JAN_CD)));
                    singleBooks.setSumStocks(cursor.getInt(cursor.getColumnIndex(COLUMN_STOCK_COUNT)));
                } while (cursor.moveToNext());
            }
        } else {
            DatabaseManager.getInstance().closeDatabase();
            return null;
        }
        DatabaseManager.getInstance().closeDatabase();
        return singleBooks;
    }

    /**
     * Get list info filter group 2 cd when group1 multi
     *
     * @return list is entity book
     */
    public List<CLP> getInfoGroupCd2WhenGroup1CdMulti(ArrayList<String> selectGroup1Cd) {

        List<CLP> clpList = new ArrayList<>();

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery;

        //Add list condition group1 cd
        String strCondition = "";
        for (int i = 0; i < selectGroup1Cd.size(); i++) {
            if (i > 0) {
                strCondition += ",";
            }
            strCondition += "'";
            strCondition += selectGroup1Cd.get(i);
            strCondition += "'";
        }

        selectQuery = String.format(
                "SELECT %s ,%s FROM %s WHERE %s IN (%s) GROUP BY %s", COLUMN_MEDIA_GROUP2_CD,
                COLUMN_MEDIA_GROUP2_NAME, TABLE_GENRE_RETURN_BOOK, COLUMN_MEDIA_GROUP1_CD, strCondition, COLUMN_MEDIA_GROUP2_CD);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                do {
                    CLP clp = new CLP();
                    clp.setId(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP2_CD)));
                    clp.setName(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP2_NAME)));

                    clpList.add(clp);
                } while (cursor.moveToNext());
            }
        } else {
            DatabaseManager.getInstance().closeDatabase();
            return null;
        }
        DatabaseManager.getInstance().closeDatabase();
        return clpList;
    }

    /**
     * Count data select group cd
     */
    private int countBooksGroupCd(FlagSettingNew flagSettingNew) {

        int count = 0;
        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s WHERE 1=1 ", ALIAS_COUNT,
                TABLE_RETURN_BOOKS);
        //Query condition filter select group cd
        String queryConditionGroupCd = conditionQueryCommon.conditionFilterSettingGroupCd(flagSettingNew);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery + queryConditionGroupCd, null);

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
