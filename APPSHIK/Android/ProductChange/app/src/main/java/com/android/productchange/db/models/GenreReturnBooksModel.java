package com.android.productchange.db.models;

import android.database.sqlite.SQLiteDatabase;

import com.android.productchange.common.utils.DatabaseManager;

import static com.android.productchange.common.constants.Constants.*;

/**
 * Model classify return books
 * Created by cong-pv on 2018/08/31.
 */

public class GenreReturnBooksModel {

    public GenreReturnBooksModel() {
    }

    /**
     * Create table genre of return books
     *
     * @return
     */
    public static String createClassifyReturnBooksTable() {

        return String.format("CREATE TABLE %s(%s TEXT, %s TEXT, %s TEXT, %s TEXT)",
                TABLE_GENRE_RETURN_BOOK, COLUMN_MEDIA_GROUP1_CD, COLUMN_MEDIA_GROUP1_NAME,
                COLUMN_MEDIA_GROUP2_CD, COLUMN_MEDIA_GROUP2_NAME);
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
}

