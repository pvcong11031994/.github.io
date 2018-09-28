package com.android.returncandidate.db.models;

import android.annotation.SuppressLint;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;

import com.android.returncandidate.common.constants.Constants;
import com.android.returncandidate.common.utils.DatabaseManager;
import com.android.returncandidate.common.utils.FlagSettingNew;
import com.android.returncandidate.db.entity.Publisher;

import java.util.ArrayList;
import java.util.List;

import static com.android.returncandidate.common.constants.Constants.*;

/**
 * Publishers model
 *
 * @author cong-pv
 * @since 2018-07-09
 */

public class PublisherModel {

    /**
     * SQLite Database
     */
    private SQLiteDatabase db;

    /**
     * Constructor Model CLP
     */
    public PublisherModel() {

    }

    /**
     * Get list info Publisher
     */
    public List<Publisher> getInfoPublisher(FlagSettingNew flagSettingNew) {

        List<Publisher> publisherList = new ArrayList<>();

        db = DatabaseManager.getInstance().openDatabase();
        //Query condition filter select group cd
        String queryConditionGroupCd = conditionFilterSettingGroupCd(flagSettingNew);

        String selectQuery = String.format(
                "SELECT %s, SUM(%s) cnt FROM %s WHERE 1=1 %s GROUP BY %s ORDER BY cnt DESC, %s",
                COLUMN_PUBLISHER_NAME, COLUMN_COUNT_PUBLISHER_NAME, TABLE_PUBLISHERS, queryConditionGroupCd, COLUMN_PUBLISHER_NAME, COLUMN_PUBLISHER_NAME);

        @SuppressLint("Recycle")
        Cursor cursor = db.rawQuery(selectQuery, null);

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

    /**
     * Check data
     */
    public boolean checkData() {

        boolean bool;

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery;
        selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_LARGE_CLASSIFICATIONS);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        bool = cursor.moveToFirst() && cursor.getInt(cursor.getColumnIndex(ALIAS_COUNT)) > 0;

        DatabaseManager.getInstance().closeDatabase();

        return bool;
    }

    /**
     * Count data in table Publisher
     *
     * @return int
     */
    public int countPublisher() {

        int count = 0;
        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_LARGE_CLASSIFICATIONS);

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
