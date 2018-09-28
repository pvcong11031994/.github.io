package com.android.returncandidate.activities;


import android.content.*;
import android.os.*;
import android.support.v7.app.*;

import com.android.returncandidate.common.helpers.*;
import com.android.returncandidate.common.utils.*;
import com.facebook.stetho.Stetho;


/**
 * Application initialization<br>
 * Flow : {@link MainActivity} ▶ {@link LoginActivity}
 *
 * @author minh-th
 * @version 2.0
 * @since 2018-05-10
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
        Stetho.initializeWithDefaults(this);
        // Redirect to login screen
        Intent intent = new Intent(this, LoginActivity.class);
        startActivity(intent);

        // Init database
        DatabaseManager.initializeInstance(new DatabaseHelper(getApplicationContext()));
        new DatabaseHelper(this);
    }
}
