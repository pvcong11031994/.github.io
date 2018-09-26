package com.android.productchange.db.models;

import static com.android.productchange.common.constants.Constants.*;

import android.annotation.SuppressLint;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteStatement;

import com.android.productchange.api.Config;
import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.utils.DatabaseManager;
import com.android.productchange.db.entity.Books;
import com.android.productchange.db.entity.CLP;
import com.android.productchange.db.entity.Regularbooks;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * <h1>Model Return Book</h1>
 *
 * @author tien-lv
 * @since 2017-12-06
 */
public class RegularbookModel {

    /**
     * Entity Return Books
     */
    private Regularbooks regularbooks;

    /**
     * SQLite Database
     */
    private SQLiteDatabase db;

    /**
     * SQLite Statement
     */
    private SQLiteStatement stmt, stmtView;

    /**
     * Constructor Model Return Book
     */
    public RegularbookModel() {
        regularbooks = new Regularbooks();
    }

    /**
     * Constructor Model Book
     *
     * @param db       SQLite Database
     * @param isInsert check is insert
     */
    public RegularbookModel(boolean isInsert, SQLiteDatabase db) {
        if (isInsert) {
            this.db = db;
            stmt = db.compileStatement(getSqlInsert());
            stmtView = db.compileStatement(getSqlViewInsert());
        }
    }

    /**
     * Function create table Regular Book
     *
     * @return String
     */
    public static String createTable() {

        return String.format(
                "CREATE TABLE %s(%s TEXT,%s TEXT,%s TEXT, %s TEXT,%s TEXT, %s TEXT,%s TEXT, " +
                        "%s TEXT, %s INTEGER, %s INTEGER)", TABLE_REGULAR_BOOKS, COLUMN_LOCATION_ID,
                COLUMN_LARGE_CLASSIFICATION_ID, COLUMN_LARGE_CLASSIFICATION_NAME, COLUMN_NAME,
                COLUMN_PUBLISHER_ID, COLUMN_PUBLISHER_NAME, COLUMN_PUBLISH_DATE, COLUMN_JAN_CODE,
                COLUMN_INVENTORY_NUMBER, COLUMN_RANKING);
    }

    /**
     * Function create table view Regular Book
     *
     * @return String
     */
    public static String createViewTable() {

        return String.format(
                "CREATE TABLE %s(%s INTEGER, %s TEXT,%s TEXT,%s TEXT, %s TEXT, %s "
                        + "TEXT, %s TEXT, %s TEXT, %s TEXT,%s INTEGER, %s INTEGER, %s INTEGER)",
                TABLE_VIEW_REGULAR_BOOKS, COLUMN_ID, COLUMN_LOCATION_ID,
                COLUMN_LARGE_CLASSIFICATION_ID, COLUMN_LARGE_CLASSIFICATION_NAME, COLUMN_JAN_CODE,
                COLUMN_NAME, COLUMN_PUBLISHER_ID, COLUMN_PUBLISHER_NAME, COLUMN_PUBLISH_DATE,
                COLUMN_INVENTORY_NUMBER, COLUMN_NEW_CATEGORY_RANK, COLUMN_RANKING);
    }

    /**
     * String sql insert statement
     *
     * @return String
     */
    private static String getSqlInsert() {
        return String.format(
                "INSERT INTO %s (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) "
                        + "VALUES (?,?,?,?,?,?,?,?,?,?)",
                TABLE_REGULAR_BOOKS, COLUMN_LOCATION_ID, COLUMN_LARGE_CLASSIFICATION_ID,
                COLUMN_LARGE_CLASSIFICATION_NAME, COLUMN_NAME, COLUMN_PUBLISHER_ID, COLUMN_PUBLISHER_NAME,
                COLUMN_PUBLISH_DATE, COLUMN_JAN_CODE, COLUMN_INVENTORY_NUMBER, COLUMN_RANKING);
    }

    /**
     * String sql view insert statement
     *
     * @return String
     */
    private static String getSqlViewInsert() {
        return String.format(
                "INSERT INTO %s (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) "
                        + "VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
                TABLE_VIEW_REGULAR_BOOKS, COLUMN_ID, COLUMN_LOCATION_ID,
                COLUMN_LARGE_CLASSIFICATION_ID, COLUMN_LARGE_CLASSIFICATION_NAME, COLUMN_JAN_CODE,
                COLUMN_NAME, COLUMN_PUBLISHER_ID, COLUMN_PUBLISHER_NAME, COLUMN_PUBLISH_DATE,
                COLUMN_INVENTORY_NUMBER, COLUMN_NEW_CATEGORY_RANK, COLUMN_RANKING);
    }

    /**
     * Function insert into table regular books
     */
    public void insertData(SQLiteDatabase db, int indexListString, List<String> listValue) {

        StringBuilder valuesBuilder = new StringBuilder();
        for (int i = 0; i < indexListString; i += Constants.VALUE_COUNT_COLUMN_TABLE_REGULAR_BOOK_INSERT) {
            if (i != 0) {
                valuesBuilder.append(", ");
            }
            valuesBuilder.append("(?,?,?,?,?,?,?,?,?,?)");
        }
        stmt = db.compileStatement(String.format("INSERT INTO %s (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) VALUES ",
                TABLE_REGULAR_BOOKS, COLUMN_LOCATION_ID, COLUMN_LARGE_CLASSIFICATION_ID,
                COLUMN_LARGE_CLASSIFICATION_NAME, COLUMN_NAME, COLUMN_PUBLISHER_ID, COLUMN_PUBLISHER_NAME,
                COLUMN_PUBLISH_DATE, COLUMN_JAN_CODE, COLUMN_INVENTORY_NUMBER, COLUMN_RANKING) + valuesBuilder.toString());
        for (int i = 0; i < indexListString; i += Constants.VALUE_COUNT_COLUMN_TABLE_REGULAR_BOOK_INSERT) {
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
        }
        stmt.executeInsert();
        stmt.clearBindings();

    }

    /**
     * Function insert into table regularbooks
     *
     * @param books is entity return books
     */
    public void insertViewBulk(Books books) {

        stmtView.bindString(1, String.valueOf(books.getId()));
        stmtView.bindString(2, books.getLocation_id());
        stmtView.bindString(3, books.getLarge_classifications_id());
        stmtView.bindString(4, books.getLarge_classifications_name());
        stmtView.bindString(5, books.getJan_code());
        stmtView.bindString(6, books.getName());
        stmtView.bindString(7, books.getPublisher_id());
        stmtView.bindString(8, books.getPublisher_name());
        stmtView.bindString(9, books.getPublish_date());
        stmtView.bindString(10, String.valueOf(books.getInventory_number()));
        stmtView.bindString(11, String.valueOf(books.getNew_catagory_rank()));
        stmtView.bindString(12, String.valueOf(books.getRanking()));
        stmtView.execute();
        stmtView.clearBindings();
    }

    /**
     * Get book info filter in table regular books
     *
     * @param id       {@link String}
     * @param type     {@link Integer}
     * @param offset   {@link Integer}
     * @param mapOrder {@link HashMap}
     * @return list has entity is return books
     */
    public List<Books> getListBookInfo(String id, int type, int offset,
                                       HashMap<String, String> mapOrder) {

        List<Books> regularbooksList = new ArrayList<>();

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format(
                "SELECT CASE WHEN B.%s IS NULL THEN RB.%s ELSE B.%s END %s,"
                        + "CASE WHEN B.%s IS NULL THEN RB.%s ELSE B.%s END %s,"
                        + "CASE WHEN B.%s IS NULL THEN RB.%s ELSE B.%s END %s,"
                        + "CASE WHEN B.%s IS NULL THEN RB.%s ELSE B.%s END %s,"
                        + "CASE WHEN B.%s IS NULL THEN RB.%s ELSE B.%s END %s,"
                        + "CASE WHEN B.%s IS NULL THEN RB.%s ELSE B.%s END %s,"
                        + "CASE WHEN B.%s IS NULL THEN RB.%s ELSE B.%s END %s,"
                        + "CASE WHEN B.%s IS NULL THEN RB.%s ELSE B.%s END %s,"
                        + " RB.%s,"
                        + "B.%s AS %s,CASE WHEN B.%s IS NULL "
                        + "THEN %s ELSE B.%s END %s FROM %s AS RB LEFT JOIN %s AS B ON RB.%s = B"
                        + ".%s",
                COLUMN_NAME, COLUMN_NAME, COLUMN_NAME, COLUMN_NAME,
                COLUMN_LOCATION_ID, COLUMN_LOCATION_ID, COLUMN_LOCATION_ID, COLUMN_LOCATION_ID,

                COLUMN_LARGE_CLASSIFICATION_NAME, COLUMN_LARGE_CLASSIFICATION_NAME,
                COLUMN_LARGE_CLASSIFICATION_NAME, COLUMN_LARGE_CLASSIFICATION_NAME,

                COLUMN_PUBLISHER_NAME, COLUMN_PUBLISHER_NAME, COLUMN_PUBLISHER_NAME,
                COLUMN_PUBLISHER_NAME,

                COLUMN_PUBLISH_DATE, COLUMN_PUBLISH_DATE,
                COLUMN_PUBLISH_DATE, COLUMN_PUBLISH_DATE,

                COLUMN_INVENTORY_NUMBER, COLUMN_INVENTORY_NUMBER, COLUMN_INVENTORY_NUMBER,
                COLUMN_INVENTORY_NUMBER,

                COLUMN_LARGE_CLASSIFICATION_ID, COLUMN_LARGE_CLASSIFICATION_ID,
                COLUMN_LARGE_CLASSIFICATION_ID, COLUMN_LARGE_CLASSIFICATION_ID,

                COLUMN_PUBLISHER_ID, COLUMN_PUBLISHER_ID, COLUMN_PUBLISHER_ID, COLUMN_PUBLISHER_ID,

                COLUMN_JAN_CODE,

                COLUMN_NEW_CATEGORY_RANK, COLUMN_NEW_CATEGORY_RANK,
                COLUMN_RANKING, INT_9999999, COLUMN_RANKING, COLUMN_RANKING,
                TABLE_REGULAR_BOOKS, TABLE_BOOKS,
                COLUMN_JAN_CODE, COLUMN_JAN_CODE);
        String selectQueryWhere = " WHERE 1=1 ";

        if (!id.equals("-1")) {
            switch (type) {
                case Config.TYPE_CLASSIFY:
                    selectQueryWhere += String.format(" AND RB.%s = '%s'",
                            COLUMN_LARGE_CLASSIFICATION_ID, id);
                    break;
                case Config.TYPE_PUBLISHER:
                    selectQueryWhere += String.format(" AND RB.%s = '%s'", COLUMN_PUBLISHER_ID,
                            id);
                    break;
                default:
                    break;
            }
        }
        String selectQueryGroupBy = String.format(" GROUP BY RB.%s ", COLUMN_JAN_CODE);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(
                selectQuery + selectQueryWhere + selectQueryGroupBy, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                do {
                    Books regularbooks = new Books();
                    regularbooks.setLocation_id(
                            cursor.getString(cursor.getColumnIndex(COLUMN_LOCATION_ID)));
                    regularbooks.setLarge_classifications_id(cursor.getString(
                            cursor.getColumnIndex(COLUMN_LARGE_CLASSIFICATION_ID)));
                    regularbooks.setLarge_classifications_name(cursor.getString(
                            cursor.getColumnIndex(COLUMN_LARGE_CLASSIFICATION_NAME)));
                    regularbooks.setName(cursor.getString(cursor.getColumnIndex(COLUMN_NAME)));
                    regularbooks.setPublisher_id(
                            cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_ID)));
                    regularbooks.setPublisher_name(
                            cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_NAME)));
                    regularbooks.setPublish_date(
                            cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISH_DATE)));
                    regularbooks.setJan_code(
                            cursor.getString(cursor.getColumnIndex(COLUMN_JAN_CODE)));
                    regularbooks.setInventory_number(Integer.parseInt(
                            cursor.getString(cursor.getColumnIndex(COLUMN_INVENTORY_NUMBER))));
                    regularbooks.setNew_catagory_rank(
                            cursor.getInt(cursor.getColumnIndex(COLUMN_NEW_CATEGORY_RANK)));
                    regularbooks.setRanking(Integer.parseInt(
                            cursor.getString(cursor.getColumnIndex(COLUMN_RANKING))));

                    regularbooksList.add(regularbooks);
                } while (cursor.moveToNext());
            }
        } else {
            DatabaseManager.getInstance().closeDatabase();
            return null;
        }
        DatabaseManager.getInstance().closeDatabase();
        return regularbooksList;
    }

    /**
     * Get book info filter in table view regular books
     *
     * @param id       {@link String}
     * @param type     {@link Integer}
     * @param offset   {@link Integer}
     * @param mapOrder {@link HashMap}
     * @return list has entity is return books
     */
    public List<Books> getListViewBookInfo(String id, int type, int offset,
                                           Map<Integer, String> mapOrder) {

        List<Books> regularbooksList = new ArrayList<>();

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format(
                "SELECT %s,"
                        + " %s,"
                        + " %s,"
                        + " %s,"
                        + " %s,"
                        + " %s,"
                        + " %s,"
                        + " %s,"
                        + " %s,"
                        + " %s FROM %s ",
                COLUMN_NAME,
                COLUMN_LOCATION_ID,
                COLUMN_LARGE_CLASSIFICATION_NAME,
                COLUMN_PUBLISHER_NAME,
                COLUMN_PUBLISH_DATE,
                COLUMN_INVENTORY_NUMBER,
                COLUMN_LARGE_CLASSIFICATION_ID,
                COLUMN_PUBLISHER_ID,
                COLUMN_NEW_CATEGORY_RANK,
                COLUMN_RANKING,
                TABLE_VIEW_REGULAR_BOOKS);
        String selectQueryWhere = " WHERE 1=1 ";

        if (!id.equals("-1")) {
            switch (type) {
                case Config.TYPE_CLASSIFY:
                    selectQueryWhere += String.format(" AND %s = '%s'",
                            COLUMN_LARGE_CLASSIFICATION_ID, id);
                    break;
                case Config.TYPE_PUBLISHER:
                    selectQueryWhere += String.format(" AND %s = '%s'", COLUMN_PUBLISHER_ID,
                            id);
                    break;
                default:
                    break;
            }
        }

        String selectQueryOrder = ORDER_BY;
        int count = 0;
        if (mapOrder != null) {
            for (Integer key : mapOrder.keySet()) {
                if (!mapOrder.get(key).isEmpty()) {
                    switch (key) {
                        case Constants.NUMBER_1:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_LOCATION_ID, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_3:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_NAME, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_4:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_LARGE_CLASSIFICATION_ID, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_5:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_PUBLISHER_ID, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_8:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_PUBLISH_DATE, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_9:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_INVENTORY_NUMBER, mapOrder.get(key));
                            break;
                        case Constants.NUMBER_10:
                            selectQueryOrder += String.format(" %s %s,", Constants.COLUMN_RANKING, mapOrder.get(key));
                            break;
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
                selectQueryOrder += String.format(" %s ", Constants.COLUMN_RANKING);
            }
        }
        String selectQueryLimit = "";
//        String.format(" LIMIT 1000 OFFSET %s", offset);

        @SuppressLint("Recycle")
        Cursor cursor = db.rawQuery(
                selectQuery + selectQueryWhere + selectQueryOrder + selectQueryLimit, null);

        if (cursor != null)

        {
            if (cursor.moveToFirst()) {
                do {
                    Books regularbooks = new Books();
                    regularbooks.setLocation_id(
                            cursor.getString(cursor.getColumnIndex(COLUMN_LOCATION_ID)));
                    regularbooks.setLarge_classifications_name(
                            cursor.getString(
                                    cursor.getColumnIndex(COLUMN_LARGE_CLASSIFICATION_NAME)));
                    regularbooks.setName(cursor.getString(cursor.getColumnIndex(COLUMN_NAME)));
                    regularbooks.setPublisher_name(
                            cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_NAME)));
                    regularbooks.setPublish_date(
                            cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISH_DATE)));
                    regularbooks.setInventory_number(Integer.parseInt(
                            cursor.getString(cursor.getColumnIndex(COLUMN_INVENTORY_NUMBER))));
                    regularbooks.setNew_catagory_rank(
                            cursor.getInt(cursor.getColumnIndex(COLUMN_NEW_CATEGORY_RANK)));
                    regularbooks.setRanking(Integer.parseInt(
                            cursor.getString(cursor.getColumnIndex(COLUMN_RANKING))));

                    regularbooksList.add(regularbooks);
                } while (cursor.moveToNext());
            }
        } else

        {
            DatabaseManager.getInstance().closeDatabase();
            return null;
        }
        DatabaseManager.getInstance().

                closeDatabase();
        return regularbooksList;
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
                TABLE_REGULAR_BOOKS);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        bool = cursor.moveToFirst() && cursor.getInt(cursor.getColumnIndex(ALIAS_COUNT)) > 0;

        DatabaseManager.getInstance().closeDatabase();

        return bool;
    }

    /**
     * Check View Data is exist
     *
     * @return boolean
     */
    public boolean checkViewData() {

        boolean bool;
        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_VIEW_REGULAR_BOOKS);

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
                        "SELECT %s , %s FROM %s ", COLUMN_ID, COLUMN_NAME, TABLE_REGULAR_BOOKS);
                break;
            case Config.TYPE_PUBLISHER:
                selectQuery = String.format(
                        "SELECT %s ,%s FROM %s ", COLUMN_PUBLISHER_ID, COLUMN_PUBLISHER_NAME,
                        TABLE_REGULAR_BOOKS);
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
                    clp.setId(cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_ID)));
                    clp.setName(cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_NAME)));

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
     * Count data in table Regular Books
     *
     * @return int
     */
    public int countBooks() {

        int count = 0;
        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_REGULAR_BOOKS);

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