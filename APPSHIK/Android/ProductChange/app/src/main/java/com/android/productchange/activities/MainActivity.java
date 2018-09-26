package com.android.productchange.activities;

import android.content.Intent;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;

import com.android.productchange.common.helpers.DatabaseHelper;
import com.android.productchange.common.utils.DatabaseManager;
import com.facebook.stetho.Stetho;

/**
 * <h1>Main Screen Activity</h1>
 *
 * Create all table
 *
 * @author tien-lv
 * @since 2018-02-08
 */
public class MainActivity extends AppCompatActivity {


    /**
     * Init on Create Activity
     *
     * @param savedInstanceState bundle of activity
     */
    @Override
    protected void onCreate(Bundle savedInstanceState) {

        super.onCreate(savedInstanceState);
        //Stetho.initializeWithDefaults(this);
        // Call Login Screen
        Intent intent = new Intent(this, LoginActivity.class);
        startActivity(intent);

        // init database
        DatabaseManager.initializeInstance(new DatabaseHelper(getApplicationContext()));
        new DatabaseHelper(this);
    }
}
