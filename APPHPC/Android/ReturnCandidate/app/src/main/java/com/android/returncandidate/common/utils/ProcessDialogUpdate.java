package com.android.returncandidate.common.utils;

import android.app.Activity;
import android.app.ProgressDialog;
import android.os.AsyncTask;
import android.view.View;
import android.widget.TextView;

import com.android.returncandidate.common.constants.Constants;
import com.android.returncandidate.common.constants.Message;
import com.android.returncandidate.db.entity.Books;
import com.android.returncandidate.db.models.BookModel;
import com.honeywell.barcode.HSMDecoder;
import com.honeywell.barcode.Symbology;

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
    private BookModel mBookModel = new BookModel();
    private Activity contextParent;
    private String flagSwitchOCR;

    public ProcessDialogUpdate(Activity contextParent, String flagSwitchOCR) {
        this.flagSwitchOCR = flagSwitchOCR;
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

        if (!(Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup1Cd().get(0)))) {
            mBookModel.updateTableBooks(flagSettingNew);
        } else {
            if (!Constants.ROW_ALL.equals(flagSettingNew.getFlagPublisherShowScreen().get(0))) {
                mBookModel.updateTableBooks(flagSettingNew);
            }
        }
        return null;
    }


    /**
     * End progress loading
     */
    @Override
    protected void onPostExecute(ArrayList<FlagSettingNew> result) {

        progressDialog.dismiss();
        // HSM init
        HSMDecoder hsmDecoder = HSMDecoder.getInstance(contextParent);
        if (Constants.FLAG_1.equals(flagSwitchOCR)) {
            hsmDecoder.enableSymbology(Symbology.OCR);
        } else {
            hsmDecoder.enableSymbology(Symbology.EAN13);
        }
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
