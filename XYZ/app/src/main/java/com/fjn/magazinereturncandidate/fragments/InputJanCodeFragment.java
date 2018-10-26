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
import android.view.WindowManager;
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
 * Fragment input jan_code
 * Created by cong-pv on 2018/10/18.
 */

public class InputJanCodeFragment extends DialogFragment implements View.OnClickListener {

    /**
     * Create interface item submit edit
     */
    public interface SubmitInputJanDataReturnMagazine {
        /**
         * Function send data
         */
        void onDataInput(String valueJanCode, String valueNumberReturn);
    }

    /**
     * Init array item in comboBox
     */
    String arr[] = new String[15];

    /**
     * Save value textview
     */
    TextView txv_inventory_number;

    TextView txv_jan_cd;

    private String flagSwitchOCR;

    private CheckDataCommon checkDataCommon;
    private RegisterLicenseCommon registerLicenseCommon;
    private HSMDecoder hsmDecoder;
    private FormatCommon formatCommon;

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState) {

        //Init custom product detail layout
        View rootView = inflater.inflate(R.layout.fragment_input_jan, container, false);
        if (getDialog().getWindow() != null) {
            getDialog().getWindow().requestFeature(Window.FEATURE_NO_TITLE);
            getDialog().getWindow().setBackgroundDrawable(new ColorDrawable(Color.TRANSPARENT));
        }

        //Get flag setting OCR
        if (getArguments() != null) {
            flagSwitchOCR = getArguments().getString(Constants.FLAG_SWITCH_OCR);
        }

        //Get Id textView Product detail
        txv_jan_cd = (TextView) rootView.findViewById(R.id.txv_jan_cd_input);
        txv_inventory_number = (TextView) rootView.findViewById(R.id.txv_inventory_number_input);
        Button btn_submit_edit = (Button) rootView.findViewById(R.id.btn_submit_edit_input);
        // Set adapter for combo box with array
        ArrayAdapter<String> adapter1 = new ArrayAdapter<>(getContext(), android.R.layout.simple_list_item_1, arr);
        // Set item selected
        adapter1.setDropDownViewResource(android.R.layout.simple_list_item_single_choice);
        // Set button click
        btn_submit_edit.setOnClickListener(this);

        //Set default end text
        txv_jan_cd.setText("");
        txv_jan_cd.append(Constants.PREFIX_JAN_CODE_MAGAZINE); //focus end text
        //Show keyboard
        txv_jan_cd.requestFocus();
        getDialog().getWindow().setSoftInputMode(WindowManager.LayoutParams.SOFT_INPUT_STATE_ALWAYS_VISIBLE);

        checkDataCommon = new CheckDataCommon();
        registerLicenseCommon = new RegisterLicenseCommon();
        formatCommon = new FormatCommon();
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
            case R.id.btn_submit_edit_input:
                saveDataInputJanReturnMagazine();
                break;
        }

    }

    //save data when click submit
    private void saveDataInputJanReturnMagazine() {

        //Check jan input
        if (checkDataCommon.validateFields(txv_jan_cd)) {
            //Check digit jan_code
            String strJanInput = txv_jan_cd.getText().toString();
            String strResult = checkDataCommon.validateCheckDigit(strJanInput, txv_jan_cd);

            if (strResult != null && checkDataCommon.validateFieldsNotNull(txv_inventory_number)) {
                InputJanCodeFragment.SubmitInputJanDataReturnMagazine activity = (InputJanCodeFragment.SubmitInputJanDataReturnMagazine) getActivity();
                activity.onDataInput(strResult, formatCommon.formatNumber(txv_inventory_number.getText().toString()));
                dismiss();
            }
        }

    }
}
