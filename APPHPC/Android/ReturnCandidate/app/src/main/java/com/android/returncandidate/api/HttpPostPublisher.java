package com.android.returncandidate.api;

import android.app.ProgressDialog;
import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.os.AsyncTask;

import com.android.returncandidate.common.constants.Constants;
import com.android.returncandidate.common.constants.Message;
import com.android.returncandidate.common.utils.DatabaseManager;
import com.android.returncandidate.common.utils.LogManager;
import com.android.returncandidate.db.entity.CLP;
import com.android.returncandidate.db.entity.Publisher;
import com.android.returncandidate.db.models.CLPModel;
import com.android.returncandidate.db.models.PublisherModel;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.stream.JsonReader;

import org.json.JSONException;
import org.json.JSONObject;

import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.io.OutputStreamWriter;
import java.net.URL;
import java.util.Calendar;
import java.util.Date;

import javax.net.ssl.HttpsURLConnection;

/**
 * <h1>Http Post Publisher</h1>
 * <p>
 * Task http Post to API get list Publisher data
 *
 * @author cong-pv
 * @since 2018-07-09
 */

public class HttpPostPublisher extends AsyncTask<String, String, String> {

    /**
     * Http response
     */
    private HttpResponse response;

    /**
     * Progress dialog
     */
    ProgressDialog progressDialog;

    /**
     * Multi thread count
     */
    private int multiThreadCount;

    /**
     * Constructor HttpPost
     *
     * @param c context
     */
    public HttpPostPublisher(Context c) {

        this.response = (HttpResponse) c;
        progressDialog = new ProgressDialog(c);
        multiThreadCount = 4;
    }

    /**
     * Set progress dialog loading
     */
    protected void onPreExecute() {

        progressDialog.setMessage(Message.MESSAGE_DOWNLOAD_DATA_SCREEN);
        progressDialog.setIndeterminate(false);
        progressDialog.setCancelable(false);
        progressDialog.setCanceledOnTouchOutside(false);
        progressDialog.show();
    }

    /**
     * Send request and get response to services API
     *
     * @param params String params for activity
     * @return result from API
     * @throws IOException from insert error
     * @throws Exception   from request fail
     * @see IOException
     * @see Exception
     */
    @Override
    protected String doInBackground(String... params) {

        CLPModel clpModel = new CLPModel();

        //Insert data to table publisher
        clpModel.getSqlInsertPublisher();

        clpModel.getSqlInsertClassify();


        return null;
    }

    /**
     * End progress loading
     *
     * @param result from API
     */
    @Override
    protected void onPostExecute(String result) {

        progressDialog.dismiss();
        response.progressFinish(result, multiThreadCount, null);
    }

    /**
     * Update progress bar
     *
     * @param values updated
     */
    @Override
    protected void onProgressUpdate(String... values) {
        progressDialog.setMessage(values[0]);
    }
}
