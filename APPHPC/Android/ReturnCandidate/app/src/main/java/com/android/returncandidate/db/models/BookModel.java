package com.android.returncandidate.db.models;

import static com.android.returncandidate.common.constants.Constants.*;

import android.annotation.SuppressLint;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteStatement;

import com.android.returncandidate.common.constants.Constants;
import com.android.returncandidate.common.utils.Common;
import com.android.returncandidate.common.utils.DatabaseManager;
import com.android.returncandidate.common.utils.FlagSettingNew;
import com.android.returncandidate.db.entity.Books;
import com.android.returncandidate.db.entity.CLP;

import java.util.ArrayList;
import java.util.List;

/**
 * Book model
 *
 * @author cong-pv
 * @since 2018-06-20
 */

public class BookModel {

    /**
     * SQLite Database
     */
    private SQLiteDatabase db;

    private SQLiteStatement stmt;
    private Common common = new Common();

    /**
     * Constructor Model Book
     */
    public BookModel() {

    }

    /**
     * Constructor Model Book
     */
    public BookModel(boolean isInsert, SQLiteDatabase db) {
        if (isInsert) {
            this.db = db;
            stmt = db.compileStatement(getSqlInsert());
        }
    }

    private static String getSqlInsert() {
        return String.format(
                "INSERT INTO %s (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) "
                        + "VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
                TABLE_BOOKS, COLUMN_JAN_CD, COLUMN_STOCK_COUNT, COLUMN_GOODS_NAME,
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
                TABLE_BOOKS, COLUMN_JAN_CD, COLUMN_STOCK_COUNT, COLUMN_GOODS_NAME,
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
                TABLE_BOOKS, COLUMN_JAN_CD, COLUMN_STOCK_COUNT, COLUMN_GOODS_NAME,
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
     * Get Item Book
     */
    public Books getItemBook(String jan_code) {

        Books books = new Books();

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT * FROM %s WHERE %s = '%s'", TABLE_BOOKS, COLUMN_JAN_CD, jan_code);
        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                do {
                    books.setJan_cd(cursor.getString(cursor.getColumnIndex(COLUMN_JAN_CD)));
                    books.setBqsc_stock_count(cursor.getInt(cursor.getColumnIndex(COLUMN_STOCK_COUNT)));
                    books.setBqgm_goods_name(cursor.getString(cursor.getColumnIndex(COLUMN_GOODS_NAME)));
                    books.setBqgm_writer_name(cursor.getString(cursor.getColumnIndex(COLUMN_WRITER_NAME)));
                    books.setBqgm_publisher_cd(cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_CD)));
                    books.setBqgm_publisher_name(cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_NAME)));
                    books.setBqgm_price(cursor.getFloat(cursor.getColumnIndex(COLUMN_PRICE)));
                    books.setBqtse_first_supply_date(cursor.getString(cursor.getColumnIndex(COLUMN_FIRST_SUPPLY_DATE)));
                    books.setBqtse_last_supply_date(cursor.getString(cursor.getColumnIndex(COLUMN_LAST_SUPPLY_DATE)));
                    books.setBqtse_last_sale_date(cursor.getString(cursor.getColumnIndex(COLUMN_LAST_SALES_DATE)));
                    books.setBqtse_last_order_date(cursor.getString(cursor.getColumnIndex(COLUMN_LAST_ORDER_DATE)));
                    books.setBqct_media_group1_cd(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP1_CD)));
                    books.setBqct_media_group1_name(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP1_NAME)));
                    books.setBqct_media_group2_cd(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP2_CD)));
                    books.setBqct_media_group2_name(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP2_NAME)));
                    books.setBqgm_sales_date(cursor.getString(cursor.getColumnIndex(COLUMN_SALES_DATE)));
                    books.setBqio_trn_date(cursor.getString(cursor.getColumnIndex(COLUMN_TRN_DATE)));
                    books.setPercent(cursor.getFloat(cursor.getColumnIndex(COLUMN_PERCENT)));
                    books.setFlag_sales(cursor.getString(cursor.getColumnIndex(COLUMN_FLAG_SALES)));
                    books.setYear_rank(cursor.getInt(cursor.getColumnIndex(COLUMN_YEAR_RANK)));
                    books.setJoubi(cursor.getInt(cursor.getColumnIndex(COLUMN_JOUBI)));
                    books.setSts_total_sales(cursor.getInt(cursor.getColumnIndex(COLUMN_TOTAL_SALES)));
                    books.setSts_total_supply(cursor.getInt(cursor.getColumnIndex(COLUMN_TOTAL_SUPPLY)));
                    books.setSts_total_return(cursor.getInt(cursor.getColumnIndex(COLUMN_TOTAL_RETURN)));
                    books.setLocation_id(cursor.getString(cursor.getColumnIndex(COLUMN_LOCATION_ID)));
                } while (cursor.moveToNext());
            }
        } else {
            DatabaseManager.getInstance().closeDatabase();
            return null;
        }
        DatabaseManager.getInstance().closeDatabase();
        return books;
    }

    /**
     * Get sum and stock when select not all + click filter setting
     *
     * @return
     */
    public Books getSumStockAndCountJanIsSelectGroup2() {

        Books singleBooks = new Books();

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT SUM(bqsc_stock_count) bqsc_stock_count, " +
                "COUNT(jan_cd) jan_cd FROM %s WHERE %s = '1'", TABLE_BOOKS, COLUMN_FLAG_SALES);

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
     * Get sum and stock when select all group cd
     *
     * @param flagSettingNew
     * @return
     */

    public Books getSumStockAndCountJanIsNotSelectGroup(FlagSettingNew flagSettingNew) {

        Books singleBooks = new Books();
        Common common = new Common();
        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT SUM(bqsc_stock_count) bqsc_stock_count, " +
                "COUNT(jan_cd) jan_cd FROM %s ", TABLE_BOOKS);

        //Query condition filter
        String queryCondition = conditionFilterSetting(flagSettingNew);
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
     * Get sum and stock when click filter group cd not all + click filter not click setting
     *
     * @param flagSettingNew
     * @return
     */
    public Books getDataSelectGroupCdCountSum(FlagSettingNew flagSettingNew) {

        Books singleBooks = new Books();
        int countRecordFilter = countBooksGroupCd(flagSettingNew);
        if (countRecordFilter > 0) {
            float intStockPercent = Float.parseFloat(common.FormatPercentLocal(flagSettingNew.getFlagStockPercent()));
            int recordLimit = (int) Math.ceil(countRecordFilter * intStockPercent);

            //Select database
            db = DatabaseManager.getInstance().openDatabase();

            //Query condition filter
            String queryCondition = conditionFilterSetting(flagSettingNew);
            //Query condition filter select group cd
            String queryConditionGroupCd = conditionFilterSettingGroupCd(flagSettingNew);

            String selectQuery = String.format("SELECT SUM(%s) %s, COUNT(%s) %s FROM (SELECT %s, %s, %s, %s, %s, %s, %s, %s FROM %s WHERE 1 = 1  %s ORDER BY %s DESC LIMIT %s) ",
                    COLUMN_STOCK_COUNT, COLUMN_STOCK_COUNT, COLUMN_JAN_CD, COLUMN_JAN_CD, COLUMN_JAN_CD, COLUMN_SALES_DATE,
                    COLUMN_TRN_DATE, COLUMN_LAST_SALES_DATE, COLUMN_LAST_SUPPLY_DATE, COLUMN_STOCK_COUNT, COLUMN_PUBLISHER_NAME, COLUMN_JOUBI,
                    TABLE_BOOKS, queryConditionGroupCd, COLUMN_YEAR_RANK, recordLimit);

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
        } else {
            singleBooks.setCountJan_Cd(0);
            singleBooks.setSumStocks(0);
            return singleBooks;
        }
    }

    /**
     * Update multi row table
     *
     * @param flagSettingNew
     */
    public void updateTableBooks(FlagSettingNew flagSettingNew) {

        int countRecordFilter = countBooksGroupCd(flagSettingNew);
        if (countRecordFilter > 0) {
            float intStockPercent = Float.parseFloat(common.FormatPercentLocal(flagSettingNew.getFlagStockPercent()));
            int recordLimit = (int) Math.ceil(countRecordFilter * intStockPercent);

            //Select database
            db = DatabaseManager.getInstance().openDatabase();
            db.beginTransaction();

            //Query condition filter
            String queryCondition = conditionFilterSetting(flagSettingNew);
            //Query condition filter select group cd
            String queryConditionGroupCd = conditionFilterSettingGroupCd(flagSettingNew);

            String selectQuery = String.format("SELECT %s FROM (SELECT %s, %s, %s, %s, %s, %s, %s, %s " +
                            "FROM %s WHERE 1 = 1  %s ORDER BY %s DESC LIMIT %s) ", COLUMN_JAN_CD, COLUMN_JAN_CD,
                    COLUMN_SALES_DATE, COLUMN_TRN_DATE, COLUMN_LAST_SALES_DATE, COLUMN_LAST_SUPPLY_DATE, COLUMN_STOCK_COUNT,
                    COLUMN_PUBLISHER_NAME, COLUMN_JOUBI, TABLE_BOOKS, queryConditionGroupCd, COLUMN_YEAR_RANK, recordLimit);

            String strUpdate = String.format("UPDATE %s SET %s = CASE WHEN %s IN (%s) THEN 1 ELSE 0 END", TABLE_BOOKS,
                    COLUMN_FLAG_SALES, COLUMN_JAN_CD, selectQuery + queryCondition);
            db.execSQL(strUpdate);
            db.setTransactionSuccessful();
            db.endTransaction();
            DatabaseManager.getInstance().closeDatabase();
        }
    }

    /**
     * Check Data
     */

    public boolean checkData() {

        boolean bool;
        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_BOOKS);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        bool = cursor.moveToFirst() && cursor.getInt(cursor.getColumnIndex(ALIAS_COUNT)) > 0;

        DatabaseManager.getInstance().closeDatabase();

        return bool;
    }

    /**
     * Count data select group cd
     */
    private int countBooksGroupCd(FlagSettingNew flagSettingNew) {

        int count = 0;
        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s WHERE 1=1 ", ALIAS_COUNT,
                TABLE_BOOKS);
        //Query condition filter select group cd
        String queryConditionGroupCd = conditionFilterSettingGroupCd(flagSettingNew);

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

    /**
     * Count data
     */
    public int countBooks() {

        int count = 0;
        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_BOOKS);

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
                COLUMN_MEDIA_GROUP1_NAME, TABLE_LARGE_CLASSIFICATIONS, COLUMN_MEDIA_GROUP1_CD);

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
     * Get list info filter group 2 cd
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
                    COLUMN_MEDIA_GROUP2_NAME, TABLE_LARGE_CLASSIFICATIONS, COLUMN_MEDIA_GROUP2_CD);
        } else {
            selectQuery = String.format(
                    "SELECT %s ,%s FROM %s WHERE %s = '%s' GROUP BY %s", COLUMN_MEDIA_GROUP2_CD,
                    COLUMN_MEDIA_GROUP2_NAME, TABLE_LARGE_CLASSIFICATIONS, COLUMN_MEDIA_GROUP1_CD, selectGroup1Cd, COLUMN_MEDIA_GROUP2_CD);
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
                COLUMN_MEDIA_GROUP2_NAME, TABLE_LARGE_CLASSIFICATIONS, COLUMN_MEDIA_GROUP1_CD, strCondition, COLUMN_MEDIA_GROUP2_CD);

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
     * Condition select filter other
     *
     * @param flagSettingNew
     * @return
     */
    private String conditionFilterSetting(FlagSettingNew flagSettingNew) {

        String queryCondition = "";

        //Format Date time
        String strReleaseDate = common.FormatDateTime(flagSettingNew.getFlagReleaseDate());
        String strUndisturbed = common.FormatDateTime(flagSettingNew.getFlagUndisturbed());

        queryCondition += String.format(" WHERE bqgm_sales_date < '%s' ", strReleaseDate);
        //Condition undisturbed
        queryCondition += String.format(" AND bqio_trn_date <= '%s' AND bqtse_last_sale_date <= '%s' " +
                        "AND bqtse_last_supply_date <= '%s' ", strUndisturbed,
                strUndisturbed, strUndisturbed);
        queryCondition += String.format(" AND bqsc_stock_count >= %s ", Integer.parseInt(flagSettingNew.getFlagNumberOfStocks()));
        //Condition

        if (flagSettingNew.getFlagPublisherShowScreen().size() > 0 &&
                !Constants.ROW_ALL.equals(flagSettingNew.getFlagPublisherShowScreen().get(0))) {
            String strPublisher = "";
            for (int i = 0; i < flagSettingNew.getFlagPublisherShowScreen().size(); i++) {
                if (i > 0) {
                    strPublisher += ",";
                }
                strPublisher += "'";
                strPublisher += flagSettingNew.getFlagPublisherShowScreen().get(i);
                strPublisher += "'";
            }
            queryCondition += String.format(" AND bqgm_publisher_name IN (%s) ", strPublisher);
        }
        if (Constants.VALUE_YES_STANDING.equals(flagSettingNew.getFlagJoubi())) {
            queryCondition += String.format(" AND joubi != %s", Constants.VALUE_JOUBI);
        }

        return queryCondition;
    }

    /**
     * Return condition select group1cd/group2cd
     *
     * @param flagSettingNew
     * @return
     */
    private String conditionFilterSettingGroupCd(FlagSettingNew flagSettingNew) {

        String queryConditionGroupCd = "";
        if (flagSettingNew.getFlagClassificationGroup1Cd().size() >= 1) {
            if (Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup1Cd().get(0))) {
                return queryConditionGroupCd;
            } else {
                if (flagSettingNew.getFlagClassificationGroup2Cd().size() == 1) {
                    if (Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup2Cd().get(0))) {
                        return queryConditionGroupCd;
                    } else {
                        queryConditionGroupCd += String.format(" AND bqct_media_group2_cd = '%s' ", flagSettingNew.getFlagClassificationGroup2Cd().get(0));
                    }
                } else if (flagSettingNew.getFlagClassificationGroup2Cd().size() > 1) {
                    String strCondition = "";
                    for (int i = 0; i < flagSettingNew.getFlagClassificationGroup2Cd().size(); i++) {
                        if (i > 0) {
                            strCondition += ",";
                        }
                        strCondition += "'";
                        strCondition += flagSettingNew.getFlagClassificationGroup2Cd().get(i);
                        strCondition += "'";

                    }
                    queryConditionGroupCd += String.format(" AND bqct_media_group2_cd IN (%s) ", strCondition);
                }
            }
        }
        return queryConditionGroupCd;
    }
}
