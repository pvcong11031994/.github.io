package com.android.productchange.api;

import android.app.ProgressDialog;
import android.content.Context;
import android.os.AsyncTask;

import com.android.productchange.db.models.GenreReturnBooksModel;
import com.android.productchange.db.models.PublisherReturnBooksModel;
import com.android.productchange.interfaces.HttpResponse;
import com.android.productchange.common.constants.Message;

import java.io.IOException;

/**
 * <h1>Http Post Publisher</h1>
 * <p>
 * Task http Post to API get list Publisher data
 *
 * @author cong-pv
 * @since 2018-08-31
 */

public class HttpPostPublisherReturnBooks extends AsyncTask<String, String, String> {

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
    public HttpPostPublisherReturnBooks(Context c) {

        this.response = (HttpResponse) c;
        progressDialog = new ProgressDialog(c);
        multiThreadCount = 6;
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


        PublisherReturnBooksModel publisherReturnBooksModel = new PublisherReturnBooksModel();
        GenreReturnBooksModel genreReturnBooksModel = new GenreReturnBooksModel();

        //Insert data to table publisher
        publisherReturnBooksModel.getSqlInsertPublisherReturnBooks();

        //Insert data to table classify
        genreReturnBooksModel.getSqlInsertClassify();


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
