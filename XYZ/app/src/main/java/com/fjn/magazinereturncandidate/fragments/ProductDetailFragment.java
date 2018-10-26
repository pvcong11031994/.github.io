package com.fjn.magazinereturncandidate.fragments;

import android.app.Dialog;
import android.graphics.Color;
import android.graphics.drawable.ColorDrawable;
import android.os.Bundle;
import android.support.v4.app.DialogFragment;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.view.Window;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.TextView;

import com.fjn.magazinereturncandidate.R;
import com.fjn.magazinereturncandidate.common.constants.Constants;
import com.fjn.magazinereturncandidate.common.utils.CheckDataCommon;
import com.fjn.magazinereturncandidate.common.utils.FormatCommon;
import com.fjn.magazinereturncandidate.common.utils.RegisterLicenseCommon;
import com.honeywell.barcode.HSMDecoder;


/**
 * Created by cong-pv on 2018/10/17.
 * Dialog product detail screen
 */

public class ProductDetailFragment extends DialogFragment implements View.OnClickListener {

    /**
     * Create interface item submit edit
     */
    public interface SubmitEditDataNumberReturnMagazine {
        /**
         * Function send data
         */
        void onDataEdit(int positionEdit, String valueEdit);
    }


    /**
     * Init array item in comboBox
     */
    String arr[] = new String[15];

    /**
     * Save position edit
     */
    int positionEdit;

    /**
     * Save value textview
     */
    TextView txv_inventory_number;

    private String flagSwitchOCR;
    private RegisterLicenseCommon registerLicenseCommon;
    private HSMDecoder hsmDecoder;
    private FormatCommon formatCommon;
    private CheckDataCommon checkDataCommon;

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState) {

        formatCommon = new FormatCommon();
        registerLicenseCommon = new RegisterLicenseCommon();
        checkDataCommon = new CheckDataCommon();

        //Init custom product detail layout
        View rootView = inflater.inflate(R.layout.fragment_product_detail, container, false);
        if (getDialog().getWindow() != null) {
            getDialog().getWindow().requestFeature(Window.FEATURE_NO_TITLE);
            getDialog().getWindow().setBackgroundDrawable(new ColorDrawable(Color.TRANSPARENT));
        }

        //Get Id textview Product detail
        TextView txv_jan_cd = (TextView) rootView.findViewById(R.id.txv_jan_cd);
        TextView txv_group1_name = (TextView) rootView.findViewById(R.id.txv_group1_name);
        TextView txv_group2_name = (TextView) rootView.findViewById(R.id.txv_group2_name);
        TextView txv_product_name = (TextView) rootView.findViewById(R.id.txv_product_name);
        TextView txv_writer_name = (TextView) rootView.findViewById(R.id.txv_writer_name);
        TextView txv_publisher_name = (TextView) rootView.findViewById(R.id.txv_publisher_name);
        TextView txv_publish_date = (TextView) rootView.findViewById(R.id.txv_publish_date);
        TextView txv_price = (TextView) rootView.findViewById(R.id.txv_price);
        txv_inventory_number = (TextView) rootView.findViewById(R.id.txv_inventory_number);
        TextView txv_first_supply_date = (TextView) rootView.findViewById(R.id.txv_first_supply_date);
        TextView txv_last_supply_date = (TextView) rootView.findViewById(R.id.txv_last_supply_date);
        TextView txv_last_sales_date = (TextView) rootView.findViewById(R.id.txv_last_sales_date);
        TextView txv_last_order_date = (TextView) rootView.findViewById(R.id.txv_last_order_date);
        TextView txv_year_rank = (TextView) rootView.findViewById(R.id.txv_year_rank);
        //TextView txv_joubi = (TextView) rootView.findViewById(R.id.txv_joubi);
        TextView txv_total_sales = (TextView) rootView.findViewById(R.id.txv_total_sales);
        TextView txv_total_supply = (TextView) rootView.findViewById(R.id.txv_total_supply);
        TextView txv_total_return = (TextView) rootView.findViewById(R.id.txv_total_return);
        TextView txv_location_id = (TextView) rootView.findViewById(R.id.txv_location_id);
        Button btn_submit_edit = (Button) rootView.findViewById(R.id.btn_submit_edit);

        Bundle bundle = getArguments();
        if (bundle != null) {

            //Show year_rank
            String valueYearRank = bundle.getString(Constants.COLUMN_YEAR_RANK);
            String showYearRank;
            if (Constants.VALUE_MAX_YEAR_RANK.equals(valueYearRank)) {
                showYearRank = Constants.SHOW_MAX_YEAR_RANK;
            } else {
                showYearRank = String.format("%s位/%s中", formatCommon.formatMoney(valueYearRank),
                        formatCommon.formatMoney(String.valueOf(bundle.getInt(Constants.COLUMN_MAX_YEAR_RANK))));
            }

            //Get bundle janCode
            String janCode = bundle.getString(Constants.COLUMN_JAN_CD);
            int lenJanCode = 0;
            String janCodeResult = "";
            if (janCode != null) {
                lenJanCode = janCode.length();
                janCodeResult = janCode.substring(0, lenJanCode - 5);
            }
            if (lenJanCode != Constants.JAN_18_CHAR) {
                janCodeResult = janCode;
            }
            txv_jan_cd.setText(janCodeResult);
            txv_group1_name.setText(bundle.getString(Constants.COLUMN_MEDIA_GROUP1_NAME));
            txv_group2_name.setText(bundle.getString(Constants.COLUMN_MEDIA_GROUP2_NAME));
            txv_product_name.setText(bundle.getString(Constants.COLUMN_GOODS_NAME));
            txv_writer_name.setText(bundle.getString(Constants.COLUMN_WRITER_NAME));
            txv_publisher_name.setText(bundle.getString(Constants.COLUMN_PUBLISHER_NAME));
            txv_publish_date.setText(
                    formatCommon.formatDate(bundle.getString(Constants.COLUMN_SALES_DATE, Constants.BLANK)));
            if (Constants.BLANK.equals(bundle.getString(Constants.COLUMN_PRICE))) {
                txv_price.setText(Constants.BLANK);
            } else {
                txv_price.setText(String.format("%s%s", Constants.SYMBOL,
                        formatCommon.formatMoney((bundle.getString(Constants.COLUMN_PRICE)))));
            }
            if (bundle.getString(Constants.COLUMN_STOCK_COUNT) == null) {
                txv_inventory_number.setText("");
                txv_inventory_number.append(Constants.BLANK); //focus end text
            } else {
                txv_inventory_number.setText("");
                txv_inventory_number.append(
                        formatCommon.formatMoney(bundle.getString(Constants.COLUMN_STOCK_COUNT)));//focus end text
            }
            txv_first_supply_date.setText(
                    formatCommon.formatDate(bundle.getString(Constants.COLUMN_FIRST_SUPPLY_DATE,
                            Constants.BLANK)));
            txv_last_supply_date.setText(
                    formatCommon.formatDate(bundle.getString(Constants.COLUMN_LAST_SUPPLY_DATE,
                            Constants.BLANK)));
            txv_last_sales_date.setText(
                    formatCommon.formatDate(bundle.getString(Constants.COLUMN_LAST_SALES_DATE,
                            Constants.BLANK)));
            txv_last_order_date.setText(
                    formatCommon.formatDate(bundle.getString(Constants.COLUMN_LAST_ORDER_DATE,
                            Constants.BLANK)));
            txv_year_rank.setText(showYearRank);
            //txv_joubi.setText(formatJoubi(bundle.getString(Constants.COLUMN_JOUBI)));
            txv_total_sales.setText(formatCommon.formatMoney(bundle.getString(Constants.COLUMN_TOTAL_SALES)));
            txv_total_supply.setText(formatCommon.formatMoney(bundle.getString(Constants.COLUMN_TOTAL_SUPPLY)));
            txv_total_return.setText(formatCommon.formatMoney(bundle.getString(Constants.COLUMN_TOTAL_RETURN)));
            txv_location_id.setText(bundle.getString(Constants.COLUMN_LOCATION_ID));

            //get position
            positionEdit = bundle.getInt(Constants.POSITION_EDIT_PRODUCT);

            //get flag OCR
            flagSwitchOCR = getArguments().getString(Constants.FLAG_SWITCH_OCR);
        }

        loadItems();

        // Set adapter for combo box with array
        ArrayAdapter<String> adapter1 = new ArrayAdapter<>(
                getContext(), android.R.layout.simple_list_item_1, arr
        );

        // Set item selected
        adapter1.setDropDownViewResource(android.R.layout.simple_list_item_single_choice);

        // Set button click
        btn_submit_edit.setOnClickListener(this);

        hsmDecoder = HSMDecoder.getInstance(getActivity());
        return rootView;
    }

    @Override
    public void onStart() {
        super.onStart();
        Dialog dialog = getDialog();
        if (dialog != null) {
            if (dialog.getWindow() != null) {
                dialog.getWindow().setLayout(ViewGroup.LayoutParams.MATCH_PARENT,
                        ViewGroup.LayoutParams.WRAP_CONTENT);
            }
        }
    }


    @Override
    public void onStop() {
        registerLicenseCommon.EnableOCROrJanCode(flagSwitchOCR, hsmDecoder);
        super.onStop();
    }

    /**
     * Event click button submit edit number return magazine
     */
    @Override
    public void onClick(View v) {

        switch (v.getId()) {
            //Event click button 保存
            case R.id.btn_submit_edit:
                saveDataEditNumberReturnMagazine();
                break;
        }

    }

    //save data when click submit
    private void saveDataEditNumberReturnMagazine() {

        if (checkDataCommon.validateFieldsNotNull(txv_inventory_number)) {
            SubmitEditDataNumberReturnMagazine activity = (SubmitEditDataNumberReturnMagazine) getActivity();
            activity.onDataEdit(positionEdit, formatCommon.formatNumber(txv_inventory_number.getText().toString()));
            dismiss();
        }
    }

    /**
     * Set item for combo box
     */
    public void loadItems() {
        for (int i = 0; i < 15; i++) {
            arr[i] = String.valueOf(i);
        }
    }
}