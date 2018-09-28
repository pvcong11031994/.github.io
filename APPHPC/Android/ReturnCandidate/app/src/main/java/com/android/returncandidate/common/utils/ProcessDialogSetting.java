package com.android.returncandidate.common.utils;

import android.app.Activity;
import android.app.ProgressDialog;
import android.os.AsyncTask;
import android.view.View;
import android.widget.TextView;

import com.android.returncandidate.R;
import com.android.returncandidate.common.constants.Constants;
import com.android.returncandidate.common.constants.Message;
import com.android.returncandidate.db.entity.Books;
import com.android.returncandidate.db.models.BookModel;

import java.io.IOException;
import java.util.ArrayList;

/**
 * @author cong-pv
 * @since 2018-07-09
 */

public class ProcessDialogSetting extends AsyncTask<ArrayList<FlagSettingNew>, String, ArrayList<FlagSettingNew>> {

    /**
     * Progress dialog
     */
    ProgressDialog progressDialog;

    public FlagSettingNew flagSettingNew = new FlagSettingNew();
    public FlagSettingOld flagSettingOld;

    private View rootView;
    private Activity contextParent;
    TextView txv_number_of_candidates, txv_number_of_books;
    Books book;
    BookModel mBookModel = new BookModel();

    public ProcessDialogSetting(Activity contextParent, View rootView) {
        this.rootView = rootView;
        this.contextParent = contextParent;
        progressDialog = new ProgressDialog(contextParent);
    }


    /**
     * Set progress dialog loading
     */
    protected void onPreExecute() {

        progressDialog.setMessage(Message.MESSAGE_WAITING_FILTER);
        progressDialog.setIndeterminate(false);
        progressDialog.setCancelable(false);
        progressDialog.setCanceledOnTouchOutside(false);
        progressDialog.show();
    }

    /**
     * @param params String params for activity
     * @return result from API
     * @throws IOException from insert error
     * @throws Exception   from request fail
     * @see IOException
     * @see Exception
     */
    @Override
    protected ArrayList<FlagSettingNew> doInBackground(ArrayList<FlagSettingNew>... params) {

        //Get param
        //ArrayList<FlagSettingNew> result = new ArrayList<FlagSettingNew>();
        ArrayList<FlagSettingNew> param = params[0];

        flagSettingNew.setFlagClassificationGroup1Cd(param.get(0).getFlagClassificationGroup1Cd());
        flagSettingNew.setFlagClassificationGroup2Cd(param.get(0).getFlagClassificationGroup2Cd());
        flagSettingNew.setFlagPublisherShowScreen(param.get(0).getFlagPublisherShowScreen());
        flagSettingNew.setFlagReleaseDate(param.get(0).getFlagReleaseDate());
        flagSettingNew.setFlagUndisturbed(param.get(0).getFlagUndisturbed());
        flagSettingNew.setFlagNumberOfStocks(param.get(0).getFlagNumberOfStocks());
        flagSettingNew.setFlagStockPercent(param.get(0).getFlagStockPercent());
        flagSettingNew.setFlagClickSetting(param.get(0).getFlagClickSetting());
        flagSettingNew.setFlagJoubi(param.get(0).getFlagJoubi());

        //Check get sum/count condition

        if (flagSettingNew.getFlagClassificationGroup1Cd().size() >= 1) {
            if (Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup1Cd().get(0))) {
                book = mBookModel.getSumStockAndCountJanIsNotSelectGroup(flagSettingNew);
            } else {
                if (Constants.VALUE_CHECK_ONCLICK_SETTING.equals(flagSettingNew.getFlagClickSetting())) {
                    book = mBookModel.getSumStockAndCountJanIsSelectGroup2();
                } else {
                    //When click filter other
                    book = mBookModel.getDataSelectGroupCdCountSum(flagSettingNew);
                }
            }
        } else {
            book = mBookModel.getSumStockAndCountJanIsSelectGroup2();
        }
        return null;
    }


    /**
     * End progress loading
     */
    @Override
    protected void onPostExecute(ArrayList<FlagSettingNew> result) {
        progressDialog.dismiss();

        //Update stocks and jan_cd
        txv_number_of_candidates = (TextView) rootView.findViewById(R.id.txv_number_of_candidates);
        txv_number_of_books = (TextView) rootView.findViewById(R.id.txv_number_of_books);
        txv_number_of_candidates.setText(book.getCountJan_Cd() + " " + Constants.VALUE_COUNT_JAN_CD);
        txv_number_of_books.setText(book.getSumStocks() + " " + Constants.VALUE_SUM_STOCKS);
    }

    /**
     * Update progress bar
     *
     * @param values updated
     */
    @Override
    protected void onProgressUpdate(String... values) {
    }
}
