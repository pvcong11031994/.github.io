package com.android.productchange.db.models;

import android.annotation.SuppressLint;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;

import com.android.productchange.common.utils.ConditionQueryCommon;
import com.android.productchange.common.utils.DatabaseManager;
import com.android.productchange.common.utils.FlagSettingNew;
import com.android.productchange.db.entity.Publisher;

import java.util.ArrayList;
import java.util.List;

import static com.android.productchange.common.constants.Constants.*;

/**
 * Model classify return books
 * Created by cong-pv on 2018/08/31.
 */

public class PublisherReturnBooksModel {

    public PublisherReturnBooksModel() {
    }

    /**
     * SQLite Database
     */
    private SQLiteDatabase db;
    private ConditionQueryCommon conditionQueryCommon = new ConditionQueryCommon();


    /**
     * Function create table the publisher
     */
    public static String createPublisherReturnBooksTable() {

        return String.format("CREATE TABLE %s(%s TEXT, %s TEXT, %s TEXT)",
                TABLE_PUBLISHERS_RETURN_BOOKS, COLUMN_PUBLISHER_NAME, COLUMN_MEDIA_GROUP2_CD,
                COLUMN_COUNT_PUBLISHER_NAME);
    }

    /**
     * String sql insert statement of publisher
     */
    public void getSqlInsertPublisherReturnBooks() {
        SQLiteDatabase db = DatabaseManager.getInstance().openDatabase();
        String query = String.format(
                "INSERT INTO %s (%s, %s, %s) SELECT %s, %s, count(*) AS %s FROM %s " +
                        "WHERE TRIM(TRIM(%s), '　') != '' AND TRIM(TRIM(%s), '　') != '' GROUP BY %s, %s",
                TABLE_PUBLISHERS_RETURN_BOOKS, COLUMN_PUBLISHER_NAME, COLUMN_MEDIA_GROUP2_CD, COLUMN_COUNT_PUBLISHER_NAME,
                COLUMN_PUBLISHER_NAME, COLUMN_MEDIA_GROUP2_CD, COLUMN_COUNT_PUBLISHER_NAME, TABLE_RETURN_BOOKS, COLUMN_MEDIA_GROUP2_CD,
                COLUMN_PUBLISHER_NAME, COLUMN_PUBLISHER_NAME, COLUMN_MEDIA_GROUP2_CD);
        db.execSQL(query);
        DatabaseManager.getInstance().closeDatabase();
    }

    /**
     * Check data
     */
    public boolean checkData() {

        boolean bool;

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery;
        selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_PUBLISHERS_RETURN_BOOKS);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        bool = cursor.moveToFirst() && cursor.getInt(cursor.getColumnIndex(ALIAS_COUNT)) > 0;

        DatabaseManager.getInstance().closeDatabase();

        return bool;
    }


    /**
     * Count data in table Publisher Return Books
     *
     * @return int
     */
    public int countDataTablePublisherReturn() {

        int count = 0;
        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_PUBLISHERS_RETURN_BOOKS);

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
     * Get list info Publisher
     */
    public List<Publisher> getInfoPublisherReturnBooks(FlagSettingNew flagSettingNew) {

        List<Publisher> publisherList = new ArrayList<>();

        db = DatabaseManager.getInstance().openDatabase();


        //Query condition filter select group cd
        String queryConditionGroupCd = conditionQueryCommon.conditionFilterSettingGroupCd(flagSettingNew);

        String selectQuery = String.format(
                "SELECT %s, SUM(%s) cnt FROM %s WHERE 1=1 %s GROUP BY %s ORDER BY cnt DESC, %s",
                COLUMN_PUBLISHER_NAME, COLUMN_COUNT_PUBLISHER_NAME, TABLE_PUBLISHERS_RETURN_BOOKS,
                queryConditionGroupCd, COLUMN_PUBLISHER_NAME, COLUMN_PUBLISHER_NAME);

        @SuppressLint("Recycle")
        Cursor cursor = db.rawQuery(selectQuery + queryConditionGroupCd, null);

        if (cursor != null) {
            if (cursor.moveToFirst()) {
                do {
                    Publisher publisher = new Publisher();
                    publisher.setId(cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_NAME)));
                    publisher.setName(cursor.getString(cursor.getColumnIndex(COLUMN_PUBLISHER_NAME)));
                    publisherList.add(publisher);
                } while (cursor.moveToNext());
            }
        } else

        {
            DatabaseManager.getInstance().closeDatabase();
            return null;
        }
        DatabaseManager.getInstance().

                closeDatabase();
        return publisherList;
    }

}

