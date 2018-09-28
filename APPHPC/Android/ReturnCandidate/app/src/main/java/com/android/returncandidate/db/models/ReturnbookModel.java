package com.android.returncandidate.db.models;

import android.content.*;
import android.database.sqlite.*;
import android.os.*;
import android.support.annotation.*;

import com.android.returncandidate.common.utils.*;
import com.android.returncandidate.db.entity.*;

import static com.android.returncandidate.common.constants.Constants.*;

/**
 * Returnbook model
 *
 * @author tien-lv
 * @since 2017-12-06
 */

public class ReturnbookModel {

    /**
     * Entity Return Books
     */
    private Returnbooks returnbooks;

    /**
     * Constructor Model Return Book
     */
    public ReturnbookModel() {
        returnbooks = new Returnbooks();
    }

    /**
     * Function create table Return Book
     */
    public static String createTable() {

        return String.format(
                "CREATE TABLE %s(%s INTEGER PRIMARY KEY, %s TEXT, %s TEXT, %s TEXT, %s "
                        + "TEXT, %s INTEGER, %s INTEGER)",
                TABLE_RETURN_BOOKS, COLUMN_ID, COLUMN_SHOP_ID, COLUMN_USER_ID, COLUMN_RETURN_DATE,
                COLUMN_JAN_CODE, COLUMN_NUMBER, COLUMN_LIST_STATUS);
    }

    /**
     * Function insert into table Return Books
     */
    @RequiresApi(api = Build.VERSION_CODES.KITKAT)
    public int insert(Returnbooks returnbooks) {

        int courseId;
        // SQLite Database
        try (SQLiteDatabase db = DatabaseManager.getInstance().openDatabase()) {
            ContentValues values = new ContentValues();
            values.put(COLUMN_ID, returnbooks.getId());
            values.put(COLUMN_SHOP_ID, returnbooks.getShop_id());
            values.put(COLUMN_USER_ID, returnbooks.getUserid());
            values.put(COLUMN_RETURN_DATE, returnbooks.getReturn_date());
            values.put(COLUMN_JAN_CODE, returnbooks.getJan_code());
            values.put(COLUMN_NUMBER, returnbooks.getNumber());
            values.put(COLUMN_LIST_STATUS, returnbooks.getList_status());

            // Inserting Row
            courseId = (int) db.insert(TABLE_RETURN_BOOKS, null, values);
        }
        DatabaseManager.getInstance().closeDatabase();

        return courseId;
    }
}
