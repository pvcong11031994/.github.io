package com.android.productchange.common.utils;

import android.app.Activity;
import android.app.ProgressDialog;
import android.os.AsyncTask;
import android.view.View;
import android.widget.ListView;

import com.android.productchange.R;
import com.android.productchange.adapters.ListViewProductReturnBooksAdapter;
import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.constants.Message;
import com.android.productchange.db.entity.Returnbooks;
import com.android.productchange.db.models.ReturnbookModel;
import com.android.productchange.interfaces.AsyncResponse;
import com.android.productchange.interfaces.AsyncResponseListReturnBooks;

import java.io.IOException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.TreeMap;

/**
 * @author cong-pv
 * @since 2018-07-09
 */

public class ProcessDialogLoadListReturnBooks extends AsyncTask<Object, String, ArrayList> {

    /**
     * Progress dialog
     */
    private ProgressDialog progressDialog;
    public FlagSettingNew flagSettingNew = new FlagSettingNew();
    private Activity contextParent;
    public AsyncResponseListReturnBooks delegate = null;

    /**
     * Offset
     */
    private int offset;
    private Map<Integer, String> treeMapOrder = new TreeMap<>();
    private FormatCommon formatCommon = new FormatCommon();
    /**
     * List data for List View
     */
    private ArrayList<HashMap<String, String>> list = new ArrayList<>();

    /**
     * List view
     */
    private ListView lstView;

    private View rootView;


    /**
     * Return book model
     */
    private ReturnbookModel returnbookModel = new ReturnbookModel();


    public ProcessDialogLoadListReturnBooks(Activity contextParent, AsyncResponseListReturnBooks delegate) {
        this.delegate = delegate;
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
    protected ArrayList doInBackground(Object... params) {

        //Get param
        offset = (int) params[0];
        treeMapOrder = (TreeMap) params[1];
        flagSettingNew = (FlagSettingNew) params[2];
        lstView = (ListView) params[3];
        //list = (ArrayList) params[3];

        List<Returnbooks> returnbooksList;

        if (flagSettingNew.getFlagClassificationGroup1Cd().size() > 0 &&
                Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup1Cd().get(0))) {
            returnbooksList = returnbookModel.getListBookInfo(offset, treeMapOrder, flagSettingNew);
        } else {
            returnbooksList = returnbookModel.getListBookInfoSelectGroup1Cd(offset, treeMapOrder);
        }
        putDataListReturn(returnbooksList);

        return list;
    }


    /**
     * Function pull data when not select all group1 cd
     */

    private void putDataListReturn(List<Returnbooks> returnbooksList) {

        for (int i = 0; i < returnbooksList.size(); i++) {
            HashMap<String, String> hashMap = new HashMap<>();
            hashMap.put(Constants.COLUMN_LOCATION_ID,
                    formatCommon.formatLocationIdNewLine(returnbooksList.get(i).getLocation_id()));
            hashMap.put(Constants.COLUMN_JAN_CD,
                    returnbooksList.get(i).getJan_cd());
            hashMap.put(Constants.COLUMN_GOODS_NAME,
                    returnbooksList.get(i).getBqgm_goods_name());
            hashMap.put(Constants.COLUMN_PUBLISHER_CD,
                    returnbooksList.get(i).getBqgm_publisher_cd());
            hashMap.put(Constants.COLUMN_PUBLISHER_NAME_RETURN,
                    returnbooksList.get(i).getBqgm_publisher_name());
            hashMap.put(Constants.COLUMN_MEDIA_GROUP1_CD,
                    returnbooksList.get(i).getBqct_media_group1_cd());
            hashMap.put(Constants.COLUMN_MEDIA_GROUP1_NAME,
                    returnbooksList.get(i).getBqct_media_group1_name());
            hashMap.put(Constants.COLUMN_MEDIA_GROUP2_CD,
                    returnbooksList.get(i).getBqct_media_group2_cd());
            hashMap.put(Constants.COLUMN_MEDIA_GROUP2_NAME,
                    returnbooksList.get(i).getBqct_media_group2_name());
            hashMap.put(Constants.COLUMN_SALES_DATE,
                    formatCommon.formatDateBlank(returnbooksList.get(i).getBqgm_sales_date()));
            hashMap.put(Constants.COLUMN_STOCK_COUNT,
                    String.valueOf(returnbooksList.get(i).getBqsc_stock_count()));
            hashMap.put(Constants.COLUMN_YEAR_RANK,
                    String.valueOf(returnbooksList.get(i).getYear_rank()));
            hashMap.put(Constants.COLUMN_WRITER_NAME,
                    returnbooksList.get(i).getBqgm_writer_name());
            hashMap.put(Constants.COLUMN_PRICE,
                    String.valueOf(returnbooksList.get(i).getBqgm_price()));
            hashMap.put(Constants.COLUMN_FIRST_SUPPLY_DATE,
                    formatCommon.formatDateBlank(String.valueOf(returnbooksList.get(i).getBqtse_first_supply_date())));
            hashMap.put(Constants.COLUMN_LAST_SUPPLY_DATE,
                    formatCommon.formatDateBlank(returnbooksList.get(i).getBqtse_last_supply_date()));
            hashMap.put(Constants.COLUMN_LAST_ORDER_DATE,
                    formatCommon.formatDateBlank(returnbooksList.get(i).getBqtse_last_order_date()));
            hashMap.put(Constants.COLUMN_LAST_SALES_DATE,
                    formatCommon.formatDateBlank(returnbooksList.get(i).getBqtse_last_sale_date()));
            hashMap.put(Constants.COLUMN_TOTAL_SALES,
                    String.valueOf(returnbooksList.get(i).getSts_total_sales()));
            hashMap.put(Constants.COLUMN_TOTAL_SUPPLY,
                    String.valueOf(returnbooksList.get(i).getSts_total_supply()));
            hashMap.put(Constants.COLUMN_TOTAL_RETURN,
                    String.valueOf(returnbooksList.get(i).getSts_total_return()));
            hashMap.put(Constants.COLUMN_JOUBI,
                    String.valueOf(returnbooksList.get(i).getJoubi()));
            list.add(hashMap);
        }
    }

    /**
     * End progress loading
     */
    @Override
    protected void onPostExecute(ArrayList result) {

        progressDialog.dismiss();

        // return null;
        //lstView = (ListView) rootView.findViewById(R.id.lsv_list);

        ListViewProductReturnBooksAdapter adapterReturnBooks = new ListViewProductReturnBooksAdapter(contextParent, list);
        lstView.setAdapter(adapterReturnBooks);

        delegate.processFinish(result, adapterReturnBooks);
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
