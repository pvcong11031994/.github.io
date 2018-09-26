package com.android.productchange.fragments;

import android.app.Dialog;
import android.graphics.Color;
import android.graphics.drawable.ColorDrawable;
import android.os.Bundle;
import android.support.v4.app.DialogFragment;
import android.support.v4.app.FragmentManager;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.view.Window;
import android.widget.Button;
import android.widget.ImageButton;
import android.widget.ImageView;

import com.android.productchange.R;
import com.android.productchange.api.Config;
import com.android.productchange.common.constants.Constants;

/**
 * <h1>Select Filter</h1>
 *
 * Dialog select item filter screen
 *
 * @author tien-lv
 * @since 2017-12-19
 */
@SuppressWarnings("deprecation")
public class DSelectFragment extends DialogFragment implements View.OnClickListener {

    /**
     * ID
     */
    private String idOld, idNew;

    /**
     * Type old
     */
    private int typeOld;

    /**
     * Type new
     */
    private int typeNew;

    /**
     * Year
     */
    private int yearOld;

    /**
     * Rank
     */
    private int rank;

    /**
     * Name Classify or Publisher
     */
    private String nameClassify;

    /**
     * Back flag
     */
    private boolean flagBack;

    /**
     * Button item
     */
    Button btnCatagory, btnPublisher, btnPeriodDate;

    /**
     * Image view item
     */
    ImageView imvCatagory, imvPublisher, imvPeriodDate;

    /**
     * Button back
     */
    ImageButton imbBack;

    /**
     * Date filter
     */
    private String dateFrom, dateTo;

    /**
     * Date checked
     */
    private boolean oldDateChecked, newDateChecked;

    private int flagPercent;
    private String flagGroup1Cd;
    private String flagGroup2Cd;
    private String flagGroup2Name;
    /**
     * Init Dialog Product Detail with custom layout
     *
     * @param container         is view group
     * @param inflater          is layout inflater
     * @param saveInstanceState bundle of dialog
     * @return show dialog
     */
    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
            Bundle saveInstanceState) {

        // Init custom product detail layout
        View rootView = inflater.inflate(R.layout.layout_select, container, false);
        if (getDialog().getWindow() != null) {
            getDialog().getWindow().requestFeature(Window.FEATURE_NO_TITLE);
            getDialog().getWindow().setBackgroundDrawable(new ColorDrawable(Color.TRANSPARENT));
        }

        flagBack = false;
        if (getArguments() != null) {
            idOld = getArguments().getString(Constants.COLUMN_ID);
            typeOld = getArguments().getInt(Config.TYPE);
            typeNew = getArguments().getInt(Config.TYPE_NEW);
            yearOld = getArguments().getInt(Constants.YEAR_AGO);
            rank = getArguments().getInt(Constants.RANK);
            oldDateChecked = getArguments().getBoolean(Constants.FLAG_DATE_CHECK);
            nameClassify = getArguments().getString(Constants.COLUMN_LARGE_CLASSIFICATION_NAME);
            flagBack = getArguments().getBoolean(Constants.FLAG_BACK);
            dateFrom = getArguments().getString(Constants.COLUMN_DATE_FROM);
            dateTo = getArguments().getString(Constants.COLUMN_DATE_TO);

            //Return book
            flagGroup1Cd = getArguments().getString(Constants.FLAG_SELECT_GROUP1_CD);
            flagGroup2Cd = getArguments().getString(Constants.FLAG_SELECT_GROUP2_CD);
            flagGroup2Name = getArguments().getString(Constants.FLAG_SELECT_GROUP2_NAME);
            flagPercent = getArguments().getInt(Constants.FLAG_PERCENT_SELECTED);
        }

        if (!flagBack) {
            typeNew = typeOld;
        } else {
            typeNew = Constants.SELECT_NULL;
        }

        btnCatagory = (Button) rootView.findViewById(R.id.btn_catagory);
        btnPublisher = (Button) rootView.findViewById(R.id.btn_publisher);
        btnPeriodDate = (Button) rootView.findViewById(R.id.btn_period_date);
        imbBack = (ImageButton) rootView.findViewById(R.id.imb_back);
        imvPublisher = (ImageView) rootView.findViewById(R.id.imv_publisher_header);
        imvCatagory = (ImageView) rootView.findViewById(R.id.imv_catagory);
        imvPeriodDate = (ImageView) rootView.findViewById(R.id.imv_period_date);

        // load item checked and title
        loadChecked();

        if (rank != Constants.RANK_PERIOD) {
            btnPeriodDate.setVisibility(View.GONE);
        }

        btnCatagory.setOnClickListener(this);
        btnPublisher.setOnClickListener(this);
        btnPeriodDate.setOnClickListener(this);
        imbBack.setOnClickListener(this);

        return rootView;
    }

    /**
     * Set size layout on start
     */
    @Override
    public void onStart() {
        super.onStart();
        Dialog dialog = getDialog();
        if (dialog != null) {
            if (dialog.getWindow() != null) {
                dialog.getWindow().setLayout(ViewGroup.LayoutParams.FILL_PARENT,
                        ViewGroup.LayoutParams.FILL_PARENT);
            }
        }
    }

    /**
     * On click event
     *
     * @param v is View on click listener
     */
    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.btn_publisher:
                loadList(Config.TYPE_PUBLISHER);
                dismiss();
                break;
            case R.id.btn_catagory:
                loadList(Config.TYPE_CLASSIFY);
                dismiss();
                break;
            case R.id.btn_period_date:
                loadPeriodDate();
                dismiss();
                break;
            case R.id.imb_back:
                dismiss();
                break;
        }
    }

    /**
     * load item checked
     */
    private void loadChecked() {

        if (rank != 7) {
            switch (typeOld) {
                case Config.TYPE_PUBLISHER:
                    imvCatagory.setVisibility(View.GONE);
                    imvPublisher.setVisibility(View.VISIBLE);
                    break;
                default:
                    break;
            }
        } else {
            switch (typeOld) {
                case Config.TYPE_PUBLISHER:
                    imvCatagory.setVisibility(View.GONE);
                    imvPublisher.setVisibility(View.VISIBLE);
                    break;
                default:
                    break;
            }

            if (oldDateChecked) {
                imvPeriodDate.setVisibility(View.VISIBLE);
                imvPublisher.setVisibility(View.GONE);
                imvCatagory.setVisibility(View.GONE);
            }
        }


    }

    /**
     * Move list select dialog
     *
     * @param typeSelected is Classify or Publisher
     */
    private void loadList(int typeSelected) {

        idNew = Constants.ID_ROW_ALL;

        if (typeOld != typeSelected) {
            idNew = Constants.ID_ROW_ALL;
            typeNew = typeSelected;
            newDateChecked = false;
        }
        DListFragment dListFragment = new DListFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();
        Bundle bundle = new Bundle();
        bundle.putInt(Config.TYPE, typeOld);
        bundle.putInt(Config.TYPE_NEW, typeNew);
        bundle.putInt(Constants.YEAR_AGO, yearOld);
        bundle.putInt(Constants.RANK, rank);
        bundle.putBoolean(Constants.FLAG_DATE_CHECK, oldDateChecked);
        bundle.putBoolean(Constants.FLAG_DATE_CHECK_NEW, newDateChecked);
        bundle.putString(Constants.COLUMN_ID, idOld);
        bundle.putString(Constants.COLUMN_ID_NEW, idNew);
        bundle.putString(Constants.COLUMN_LARGE_CLASSIFICATION_NAME, nameClassify);
        bundle.putString(Constants.COLUMN_DATE_FROM, dateFrom);
        bundle.putString(Constants.COLUMN_DATE_TO, dateTo);

        //Flag return book
        bundle.putInt(Constants.FLAG_PERCENT_SELECTED, flagPercent);
        bundle.putString(Constants.FLAG_SELECT_GROUP1_CD, flagGroup1Cd);
        bundle.putString(Constants.FLAG_SELECT_GROUP2_CD, flagGroup2Cd);
        bundle.putString(Constants.FLAG_SELECT_GROUP2_NAME, flagGroup2Name);
        dListFragment.setArguments(bundle);
        dListFragment.show(fm, null);
    }

    /**
     * Move list period date dialog
     */
    private void loadPeriodDate() {

        DSelectDateFragment dSelectDateFragment = new DSelectDateFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();
        Bundle bundle = new Bundle();
        bundle.putInt(Config.TYPE, typeOld);
        bundle.putInt(Config.TYPE_NEW, typeNew);
        bundle.putInt(Constants.YEAR_AGO, yearOld);
        bundle.putInt(Constants.RANK, rank);
        bundle.putBoolean(Constants.FLAG_DATE_CHECK, oldDateChecked);
        bundle.putBoolean(Constants.FLAG_DATE_CHECK_NEW, newDateChecked);
        bundle.putString(Constants.COLUMN_ID, idOld);
        bundle.putString(Constants.COLUMN_ID_NEW, idNew);
        bundle.putString(Constants.COLUMN_LARGE_CLASSIFICATION_NAME, nameClassify);
        bundle.putString(Constants.COLUMN_DATE_FROM, dateFrom);
        bundle.putString(Constants.COLUMN_DATE_TO, dateTo);
        dSelectDateFragment.setArguments(bundle);
        dSelectDateFragment.show(fm, null);
    }
}
