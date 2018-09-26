package com.android.productchange.common.helpers;

import static com.android.productchange.common.constants.Constants.*;

import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteOpenHelper;

import com.android.productchange.db.models.BookModel;
import com.android.productchange.db.models.CLPModel;
import com.android.productchange.db.models.GenreReturnBooksModel;
import com.android.productchange.db.models.MaxYearRankModel;
import com.android.productchange.db.models.PeriodbookModel;
import com.android.productchange.db.models.PublisherReturnBooksModel;
import com.android.productchange.db.models.RegularbookModel;
import com.android.productchange.db.models.ReturnbookModel;
import com.android.productchange.db.models.UserModel;

/**
 * <h1>Database Helper</h1>
 * <p>
 * Connect to Database SQLite
 *
 * @author tien-lv
 * @since 2017-11-30
 */

public class DatabaseHelper extends SQLiteOpenHelper {

    /**
     * Constructor DatabaseHelper
     *
     * @param context Context
     */
    public DatabaseHelper(Context context) {
        super(context, DATABASE_NAME, null, DATABASE_VERSION);
    }

    /**
     * Init to Create table
     *
     * @param db SQLite database
     */
    @Override
    public void onCreate(SQLiteDatabase db) {

        db.execSQL(UserModel.createTable());
        db.execSQL(CLPModel.createLocationsTable());
        db.execSQL(CLPModel.createLagreClassificationsTable());
        db.execSQL(CLPModel.createPublishersTable());
        db.execSQL(PeriodbookModel.createTable());
        db.execSQL(RegularbookModel.createTable());
        db.execSQL(PeriodbookModel.createViewTable());
        db.execSQL(RegularbookModel.createViewTable());
        db.execSQL(BookModel.createTable());
        //Create table return books
        db.execSQL(ReturnbookModel.createTable());
        db.execSQL(ReturnbookModel.createTableTemp());
        db.execSQL(GenreReturnBooksModel.createClassifyReturnBooksTable());
        db.execSQL(PublisherReturnBooksModel.createPublisherReturnBooksTable());
        db.execSQL(MaxYearRankModel.createTable());
    }

    /**
     * Upgrade if table is exist
     *
     * @param db         SQLite database
     * @param newVersion new version
     * @param oldVersion old version
     */
    @Override
    public void onUpgrade(SQLiteDatabase db, int oldVersion, int newVersion) {

        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_USER));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_LOCATIONS));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_LARGE_CLASSIFICATIONS));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_PUBLISHERS));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_BOOKS));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_PERIOD_BOOKS));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_VIEW_PERIOD_BOOKS));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_REGULAR_BOOKS));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_VIEW_REGULAR_BOOKS));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_RETURN_BOOKS));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_RETURN_BOOKS_TEMP));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_GENRE_RETURN_BOOK));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_PUBLISHERS_RETURN_BOOKS));
        db.execSQL(String.format(QUERY_DROP_TABLE_EXIST, TABLE_MAX_YEAR_RANK));
        onCreate(db);
    }

    /**
     * Delete data in all tables
     */
    public void clearTables() {

        //SQLite database
        SQLiteDatabase db = this.getWritableDatabase();
        db.delete(TABLE_USER, null, null);
        db.delete(TABLE_LOCATIONS, null, null);
        db.delete(TABLE_LARGE_CLASSIFICATIONS, null, null);
        db.delete(TABLE_PUBLISHERS, null, null);
        db.delete(TABLE_BOOKS, null, null);
        db.delete(TABLE_PERIOD_BOOKS, null, null);
        db.delete(TABLE_VIEW_PERIOD_BOOKS, null, null);
        db.delete(TABLE_REGULAR_BOOKS, null, null);
        db.delete(TABLE_VIEW_REGULAR_BOOKS, null, null);
        db.delete(TABLE_RETURN_BOOKS, null, null);
        db.delete(TABLE_RETURN_BOOKS_TEMP, null, null);
        db.delete(TABLE_GENRE_RETURN_BOOK, null, null);
        db.delete(TABLE_PUBLISHERS_RETURN_BOOKS, null, null);
        db.delete(TABLE_MAX_YEAR_RANK, null, null);
    }
}
