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

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * <h1>Model Book</h1>
 *
 * @author tien-lv
 * @since 2017/12/06.
 */
public class BookModel {

    /**
     * SQLite Database
     */
    private SQLiteDatabase db;

    /**
     * SQlite Statement
     */
    private SQLiteStatement stmt;

    /**
     * Constructor Model Book
     */
    public BookModel() {

    }

    /**
     * Constructor Model Book
     *
     * @param db       is SQLite database
     * @param isInsert check is insert
     */
    public BookModel(boolean isInsert, SQLiteDatabase db) {
        if (isInsert) {
            this.db = db;
            stmt = db.compileStatement(getSqlInsert());
        }
    }

    /**
     * String Insert for statement
     *
     * @return String
     */
    private static String getSqlInsert() {
        return String.format(
                "INSERT INTO %s (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)",
                TABLE_BOOKS, COLUMN_LOCATION_ID, COLUMN_LARGE_CLASSIFICATION_ID, COLUMN_LARGE_CLASSIFICATION_NAME,
                COLUMN_NAME, COLUMN_PUBLISHER_ID, COLUMN_PUBLISHER_NAME, COLUMN_PUBLISH_DATE, COLUMN_JAN_CODE,
                COLUMN_INVENTORY_NUMBER, COLUMN_OLD_CATEGORY_RANK, COLUMN_NEW_CATEGORY_RANK,
                COLUMN_FLAG_ORDER_NOW, COLUMN_RANKING);
    }

    /**
     * Function create table books
     *
     * @return String
     */
    public static String createTable() {

        return String.format(
                "CREATE TABLE %s(%s TEXT,%s TEXT, %s TEXT, %s TEXT, %s TEXT,%s TEXT, %s TEXT," +
                        "%s TEXT, %s INTEGER, %s INTEGER, %s INTEGER, %s INTEGER, %s INTEGER)",
                TABLE_BOOKS, COLUMN_LOCATION_ID, COLUMN_LARGE_CLASSIFICATION_ID, COLUMN_LARGE_CLASSIFICATION_NAME,
                COLUMN_NAME, COLUMN_PUBLISHER_ID, COLUMN_PUBLISHER_NAME, COLUMN_PUBLISH_DATE, COLUMN_JAN_CODE,
                COLUMN_INVENTORY_NUMBER, COLUMN_OLD_CATEGORY_RANK, COLUMN_NEW_CATEGORY_RANK,
                COLUMN_FLAG_ORDER_NOW, COLUMN_RANKING);
    }

    /**
     * Function insert into table books
     */
    public void insertData(SQLiteDatabase db, int indexListString, List<String> listValue) {

        StringBuilder valuesBuilder = new StringBuilder();
        for (int i = 0; i < indexListString; i += Constants.VALUE_COUNT_COLUMN_TABLE_BOOK_INSERT) {
            if (i != 0) {
                valuesBuilder.append(", ");
            }
            valuesBuilder.append("(?,?,?,?,?,?,?,?,?,?,?,?,?)");
        }
        stmt = db.compileStatement(String.format("INSERT INTO %s (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) VALUES ",
                TABLE_BOOKS, COLUMN_LOCATION_ID, COLUMN_LARGE_CLASSIFICATION_ID, COLUMN_LARGE_CLASSIFICATION_NAME,
                COLUMN_NAME, COLUMN_PUBLISHER_ID, COLUMN_PUBLISHER_NAME, COLUMN_PUBLISH_DATE, COLUMN_JAN_CODE,
                COLUMN_INVENTORY_NUMBER, COLUMN_OLD_CATEGORY_RANK, COLUMN_NEW_CATEGORY_RANK,
                COLUMN_FLAG_ORDER_NOW, COLUMN_RANKING) + valuesBuilder.toString());
        for (int i = 0; i < indexListString; i += Constants.VALUE_COUNT_COLUMN_TABLE_BOOK_INSERT) {
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
        }
        stmt.executeInsert();
        stmt.clearBindings();

    }

    /**
     * Get book info with filter
     *
     * @param id       {@link String}
     * @param type     {@link Integer}
     * @param offset   {@link Integer}
     * @param rank     {@link Integer}
     * @param mapOrder {@link HashMap}
     * @return list has entity Books
     */
    public List<Books> getListBookInfo(String id, int type, int offset, int rank,
                                       Map<Integer, String> mapOrder) {

        List<Books> booksList = new ArrayList<>();

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT * FROM %s", TABLE_BOOKS);
        String selectQueryWhere;
        if (rank != RANK_ARRIVAL) {
            selectQueryWhere = String.format(" WHERE (%s = %s OR %s = %s)",
                    COLUMN_OLD_CATEGORY_RANK,
                    rank, COLUMN_NEW_CATEGORY_RANK, rank);
        } else {
            selectQueryWhere = String.format(" WHERE %s = %s", COLUMN_FLAG_ORDER_NOW,
                    FLAG_ORDER_NOW);
        }

        if (!id.equals("-1")) {
            switch (type) {
                case Config.TYPE_CLASSIFY:
                    selectQueryWhere = selectQueryWhere + String.format(" AND %s = '%s'",
                            COLUMN_LARGE_CLASSIFICATION_ID, id);
                    break;
                case Config.TYPE_PUBLISHER:
                    selectQueryWhere = selectQueryWhere + String.format(" AND %s = '%s'",
                            COLUMN_PUBLISHER_ID, id);
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

        String selectQueryLimit = String.format(" LIMIT 1000 OFFSET %s", offset);

        @SuppressLint("Recycle")
        Cursor cursor = db.rawQuery(
                selectQuery + selectQueryWhere + selectQueryOrder + selectQueryLimit, null);

        if (cursor != null)

        {
            if (cursor.moveToFirst()) {
                do {
                    Books books = new Books();
                    books.setLocation_id(
                            cursor.getString(cursor.getColumnIndex(COLUMN_LOCATION_ID)));
                    books.setLarge_classifications_id(
                            cursor.getString(
                                    cursor.getColumnIndex(COLUMN_LARGE_CLASSIFICATION_ID)));
                    books.setLarge_classifications_name(
                            cursor.getString(
                                    cursor.getColumnIndex(COLUMN_LARGE_CLASSIFICATION_NAME)));
                    books.setJan_code(cursor.getString(cursor.getColumnIndex(COLUMN_JAN_CODE)));
                    books.setName(cursor.getString(cursor.getColumnIndex(COLUMN_NAME)));
                    books.setPublisher_id(
                            cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_ID)));
                    books.setPublisher_name(
                            cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_NAME)));
                    books.setPublish_date(
                            cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISH_DATE)));
                    books.setInventory_number(
                            cursor.getInt(cursor.getColumnIndex(COLUMN_INVENTORY_NUMBER)));
                    books.setOld_catagory_rank(
                            cursor.getInt(cursor.getColumnIndex(COLUMN_OLD_CATEGORY_RANK)));
                    books.setNew_catagory_rank(
                            cursor.getInt(cursor.getColumnIndex(COLUMN_NEW_CATEGORY_RANK)));
                    books.setFlag_order_now(
                            cursor.getInt(cursor.getColumnIndex(COLUMN_FLAG_ORDER_NOW)));
                    books.setRanking(cursor.getInt(cursor.getColumnIndex(COLUMN_RANKING)));

                    booksList.add(books);
                } while (cursor.moveToNext());
            }
        } else

        {
            DatabaseManager.getInstance().closeDatabase();
            return null;
        }
        DatabaseManager.getInstance().

                closeDatabase();
        return booksList;
    }

    /**
     * Check data is exits in table Books
     *
     * @return boolean
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
     * Count data in table Books
     *
     * @return int
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
     * Get list info item is Classify or Publisher
     *
     * @param type is Publisher or Classify
     * @return list is entity CLP
     */
    public List<CLP> getInfo(int type) {

        List<CLP> clpList = new ArrayList<>();

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery;
        switch (type) {
            case Config.TYPE_CLASSIFY:
                selectQuery = String.format(
                        "SELECT %s , %s  FROM %s ", COLUMN_ID, COLUMN_NAME,
                        TABLE_LARGE_CLASSIFICATIONS);
                break;
            case Config.TYPE_PUBLISHER:
                selectQuery = String.format(
                        "SELECT %s , %s  FROM %s ", COLUMN_PUBLISHER_ID, COLUMN_PUBLISHER_NAME,
                        TABLE_BOOKS);
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
     * Get list info Publisher Or Classify
     *
     * @param type type is Publisher or Classify
     * @return list is entity CLP
     */
    public List<CLP> getInfoClassifyPublisher(int type, int ranking) {

        List<CLP> clpList = new ArrayList<>();

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery;

        switch (type) {
            case Config.TYPE_CLASSIFY:
                if (ranking == Constants.RANK_RETURN) {
                    selectQuery = String.format(
                            "SELECT %s , %s FROM %s GROUP BY %s", COLUMN_MEDIA_GROUP1_CD,
                            COLUMN_MEDIA_GROUP1_NAME, TABLE_GENRE_RETURN_BOOK, COLUMN_MEDIA_GROUP1_CD);
                } else {
                    selectQuery = String.format(
                            "SELECT %s ,%s FROM %s ", COLUMN_ID, COLUMN_NAME,
                            TABLE_LARGE_CLASSIFICATIONS);
                }
                break;
            case Config.TYPE_PUBLISHER:
                selectQuery = String.format(
                        "SELECT %s ,%s FROM %s ", COLUMN_ID, COLUMN_NAME, TABLE_PUBLISHERS);
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
                    if (type == Config.TYPE_CLASSIFY && ranking == Constants.RANK_RETURN) {
                        clp.setId(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP1_CD)));
                        clp.setName(cursor.getString(cursor.getColumnIndex(COLUMN_MEDIA_GROUP1_NAME)));
                    } else {
                        clp.setId(cursor.getString(cursor.getColumnIndex(COLUMN_ID)));
                        clp.setName(cursor.getString(cursor.getColumnIndex(COLUMN_NAME)));
                    }
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
}
