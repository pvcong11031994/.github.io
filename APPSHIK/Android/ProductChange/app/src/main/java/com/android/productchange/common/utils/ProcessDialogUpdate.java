package com.android.productchange.common.utils;

import android.app.Activity;
import android.app.ProgressDialog;
import android.os.AsyncTask;

import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.constants.Message;
import com.android.productchange.db.models.ReturnbookModel;

import java.io.IOException;
import java.util.ArrayList;

/**
 * @author cong-pv
 * @since 2018-07-09
 */

public class ProcessDialogUpdate extends AsyncTask<ArrayList<FlagSettingNew>, String, ArrayList<FlagSettingNew>> {

    /**
     * Progress dialog
     */
    private ProgressDialog progressDialog;
    public FlagSettingNew flagSettingNew = new FlagSettingNew();
    private ReturnbookModel mBookModel = new ReturnbookModel();
    private Activity contextParent;

    public ProcessDialogUpdate(Activity contextParent) {
        this.contextParent = contextParent;
        progressDialog = new ProgressDialog(contextParent);
    }


    /**
     * Set progress dialog loading
     */
    protected void onPreExecute() {

        progressDialog.setMessage(Message.MESSAGE_WAITING_UPDATE);
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



        /*if (!(Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup1Cd().get(0)))) {
            mBookModel.updateTableBooks(flagSettingNew);
        } else {
            if (!Constants.ROW_ALL.equals(flagSettingNew.getFlagPublisherShowScreen().get(0))) {
                mBookModel.updateTableBooks(flagSettingNew);
            }
        }*/
        return null;
    }


    /**
     * End progress loading
     */
    @Override
    protected void onPostExecute(ArrayList<FlagSettingNew> result) {
        progressDialog.dismiss();
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
