package com.android.returncandidate.db.models;

import android.database.sqlite.*;
import com.android.returncandidate.common.utils.*;
import static com.android.returncandidate.common.constants.Constants.*;

/**
 * Large_classifications, Locations, Publishers model
 *
 * @author cong-pv
 * @since 2018-06-18
 */

public class CLPModel {

    /**
     * Constructor Model CLP
     */
    public CLPModel() {

    }

    /**
     * Function create table large classification
     */
    public static String createLagerClassificationsTable() {

        return String.format("CREATE TABLE %s(%s TEXT, %s TEXT, %s TEXT, %s TEXT)",
                TABLE_LARGE_CLASSIFICATIONS, COLUMN_MEDIA_GROUP1_CD, COLUMN_MEDIA_GROUP1_NAME,
                COLUMN_MEDIA_GROUP2_CD, COLUMN_MEDIA_GROUP2_NAME);
    }

    /**
     * Function create table large classification
     */
    public static String createPublisherTable() {

        return String.format("CREATE TABLE %s(%s TEXT, %s TEXT, %s TEXT)",
                TABLE_PUBLISHERS, COLUMN_PUBLISHER_NAME, COLUMN_MEDIA_GROUP2_CD,
                COLUMN_COUNT_PUBLISHER_NAME);
    }
    /**
     * String sql insert statement of Classify
     */
    public void getSqlInsertClassify() {
        SQLiteDatabase db = DatabaseManager.getInstance().openDatabase();
        String query = String.format(
                "INSERT INTO %s (%s, %s, %s, %s) SELECT %s, %s, %s, %s FROM %s WHERE %s != '' GROUP BY %s, %s",
                TABLE_LARGE_CLASSIFICATIONS, COLUMN_MEDIA_GROUP1_CD, COLUMN_MEDIA_GROUP1_NAME,
                COLUMN_MEDIA_GROUP2_CD, COLUMN_MEDIA_GROUP2_NAME, COLUMN_MEDIA_GROUP1_CD,
                COLUMN_MEDIA_GROUP1_NAME, COLUMN_MEDIA_GROUP2_CD, COLUMN_MEDIA_GROUP2_NAME,
                TABLE_BOOKS, COLUMN_MEDIA_GROUP1_CD, COLUMN_MEDIA_GROUP1_CD, COLUMN_MEDIA_GROUP2_CD);
        db.execSQL(query);
        DatabaseManager.getInstance().closeDatabase();
    }

    /**
     * String sql insert statement of publisher
     */
    public void getSqlInsertPublisher() {
        SQLiteDatabase db = DatabaseManager.getInstance().openDatabase();
        String query = String.format(
                "INSERT INTO %s (%s, %s, %s) SELECT %s, %s, count(*) AS %s FROM %s " +
                        "WHERE TRIM(TRIM(%s), '　') != '' AND TRIM(TRIM(%s), '　') != '' GROUP BY %s, %s",
                TABLE_PUBLISHERS, COLUMN_PUBLISHER_NAME, COLUMN_MEDIA_GROUP2_CD, COLUMN_COUNT_PUBLISHER_NAME,
                COLUMN_PUBLISHER_NAME, COLUMN_MEDIA_GROUP2_CD, COLUMN_COUNT_PUBLISHER_NAME, TABLE_BOOKS, COLUMN_MEDIA_GROUP2_CD,
                COLUMN_PUBLISHER_NAME, COLUMN_PUBLISHER_NAME, COLUMN_MEDIA_GROUP2_CD);
        db.execSQL(query);
        DatabaseManager.getInstance().closeDatabase();
    }
}
