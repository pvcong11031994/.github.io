package com.android.productchange.db.models;

import android.annotation.SuppressLint;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;

import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.utils.ConditionQueryCommon;
import com.android.productchange.common.utils.DatabaseManager;
import com.android.productchange.common.utils.FlagSettingNew;
import com.android.productchange.db.entity.Publisher;

import java.util.ArrayList;
import java.util.List;

import static com.android.productchange.common.constants.Constants.*;

/**
 * Publishers model
 *
 * @author cong-pv
 * @since 2018-08-30
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
     * Check data
     */
    public boolean checkData() {

        boolean bool;

        db = DatabaseManager.getInstance().openDatabase();
        String selectQuery;
        selectQuery = String.format("SELECT COUNT(*) AS %s FROM %s", ALIAS_COUNT,
                TABLE_PUBLISHERS);

        @SuppressLint("Recycle") Cursor cursor = db.rawQuery(selectQuery, null);

        bool = cursor.moveToFirst() && cursor.getInt(cursor.getColumnIndex(ALIAS_COUNT)) > 0;

        DatabaseManager.getInstance().closeDatabase();

        return bool;
    }

}
