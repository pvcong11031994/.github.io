package com.fjn.magazinereturncandidate.activities;

import android.content.Intent;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;

import com.fjn.magazinereturncandidate.common.helpers.DatabaseHelper;
import com.fjn.magazinereturncandidate.common.utils.DatabaseManagerCommon;


/**
 * Application initialization<br>
 * Flow : {@link MainActivity} ▶ {@link LoginActivity}
 *
 * @author cong-pv
 * @version 2.0
 * @since 2018-10-15
 */

public class MainActivity extends AppCompatActivity {
    /**
     * Initiate database and redirects to {@link LoginActivity}
     *
     * @param savedInstanceState {@link Bundle }
     */
    @Override
    protected void onCreate(Bundle savedInstanceState) {

        super.onCreate(savedInstanceState);
        //Stetho.initializeWithDefaults(this);
        // Redirect to login screen
        Intent intent = new Intent(this, LoginActivity.class);
        startActivity(intent);

        // Init database
        DatabaseManagerCommon.initializeInstance(new DatabaseHelper(getApplicationContext()));
        new DatabaseHelper(this);

    }
}
