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
import android.widget.TextView;

import com.android.productchange.R;
import com.android.productchange.api.Config;
import com.android.productchange.common.constants.Constants;

/**
 * <h1>Select Year Dialog</h1>
 *
 * Dialog select year filter
 *
 * @author tien-lv
 * @since 2017-12-19
 */
@SuppressWarnings("deprecation")
public class DSelectYearFragment extends DialogFragment implements View.OnClickListener {

    /**
     * Interface for filter item selected to activity
     */
    public interface SelectedFilterDialogListener {

        /**
         * Function send list selected data to activity
         *
         * @param typeSelected         is type selected
         * @param idSelected           is id selected
         * @param nameClassifySelected is name Classify or Publisher selected
         * @param yearSelected         is year selected
         */
        void onSelectedFilterDialog(int typeSelected, String idSelected,
                String nameClassifySelected, int yearSelected);
    }

    /**
     * ID
     */
    private String idOld, idNew, id;

    /**
     * Type
     */
    private int typeOld, typeNew, type, yearOld, year, rank;

    /**
     * name Classify or Publisher
     */
    private String nameClassify;

    /**
     * Item Button
     */
    Button btnAllYear, btn1YearAgo, btn2YearAgo, btn3YearAgo, btn4YearAgo, btn5YearAgo, btnNull;

    /**
     * Item Image view
     */
    ImageView imvAllYear, imv1YearAgo, imv2YearAgo, imv3YearAgo, imv4YearAgo, imv5YearAgo, imvNull;

    /**
     * Button back
     */
    ImageButton imbBack;

    /**
     * Title
     */
    TextView txvPath;

    /**
     * Date checked
     */
    boolean dateChecked;

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
        View rootView = inflater.inflate(R.layout.layout_select_year, container, false);
        if (getDialog().getWindow() != null) {
            getDialog().getWindow().requestFeature(Window.FEATURE_NO_TITLE);
            getDialog().getWindow().setBackgroundDrawable(new ColorDrawable(Color.TRANSPARENT));
        }

        // get data from list dialog
        if (getArguments() != null) {
            idOld = getArguments().getString(Constants.COLUMN_ID);
            idNew = getArguments().getString(Constants.COLUMN_ID_NEW);
            typeOld = getArguments().getInt(Config.TYPE);
            typeNew = getArguments().getInt(Config.TYPE_NEW);
            yearOld = getArguments().getInt(Constants.YEAR_AGO);
            rank = getArguments().getInt(Constants.RANK);
            dateChecked = getArguments().getBoolean(Constants.FLAG_DATE_CHECK);
            nameClassify = getArguments().getString(Constants.COLUMN_LARGE_CLASSIFICATION_NAME);
        }

        if (typeNew != 0) {
            if (typeOld != typeNew) {
                type = typeNew;
                id = idNew;
                year = Constants.SELECT_ALL_YEAR;
            } else {
                type = typeOld;
                if (!idOld.equals(idNew)) {
                    id = idNew;
                    year = Constants.SELECT_ALL_YEAR;
                } else {
                    id = idOld;
                    year = yearOld;
                }
            }
        } else {
            type = typeOld;
            id = idOld;
            year = yearOld;
        }

        btnAllYear = (Button) rootView.findViewById(R.id.btn_all_year);
        btn1YearAgo = (Button) rootView.findViewById(R.id.btn_1_year_ago);
        btn2YearAgo = (Button) rootView.findViewById(R.id.btn_2_year_ago);
        btn3YearAgo = (Button) rootView.findViewById(R.id.btn_3_year_ago);
        btn4YearAgo = (Button) rootView.findViewById(R.id.btn_4_year_ago);
        btn5YearAgo = (Button) rootView.findViewById(R.id.btn_5_year_ago);
        btnNull = (Button) rootView.findViewById(R.id.btn_null);
        imbBack = (ImageButton) rootView.findViewById(R.id.imb_back);
        imvAllYear = (ImageView) rootView.findViewById(R.id.imv_all_year);
        imv1YearAgo = (ImageView) rootView.findViewById(R.id.imv_1_year_ago);
        imv2YearAgo = (ImageView) rootView.findViewById(R.id.imv_2_year_ago);
        imv3YearAgo = (ImageView) rootView.findViewById(R.id.imv_3_year_ago);
        imv4YearAgo = (ImageView) rootView.findViewById(R.id.imv_4_year_ago);
        imv5YearAgo = (ImageView) rootView.findViewById(R.id.imv_5_year_ago);
        imvNull = (ImageView) rootView.findViewById(R.id.imv_null);

        txvPath = (TextView) rootView.findViewById(R.id.txv_path);

        // load item checked and title
        loadChecked();

        btnAllYear.setOnClickListener(this);
        btn1YearAgo.setOnClickListener(this);
        btn2YearAgo.setOnClickListener(this);
        btn3YearAgo.setOnClickListener(this);
        btn4YearAgo.setOnClickListener(this);
        btn5YearAgo.setOnClickListener(this);
        btnNull.setOnClickListener(this);
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
            case R.id.btn_all_year:
                selectedFilter(Constants.SELECT_ALL_YEAR);
                break;
            case R.id.btn_1_year_ago:
                selectedFilter(Constants.SELECT_1_YEAR_AGO);
                break;
            case R.id.btn_2_year_ago:
                selectedFilter(Constants.SELECT_2_YEAR_AGO);
                break;
            case R.id.btn_3_year_ago:
                selectedFilter(Constants.SELECT_3_YEAR_AGO);
                break;
            case R.id.btn_4_year_ago:
                selectedFilter(Constants.SELECT_4_YEAR_AGO);
                break;
            case R.id.btn_5_year_ago:
                selectedFilter(Constants.SELECT_5_YEAR_AGO);
                break;
            case R.id.btn_null:
                selectedFilter(Constants.SELECT_NULL);
                break;
            case R.id.imb_back:
                loadList();
                break;
        }
    }

    /**
     * load item checked and title
     */
    private void loadChecked() {

        if (typeNew == Config.TYPE_PUBLISHER) {
            txvPath.setText(getResources().getString(R.string.select_publisher) + " "
                    + Constants.ARROW + " " + nameClassify);
        } else {
            txvPath.setText(getResources().getString(R.string.select_cagtagory) + " "
                    + Constants.ARROW + " " + nameClassify);
        }

        switch (year) {
            case Constants.SELECT_ALL_YEAR:
                imvAllYear.setVisibility(View.VISIBLE);
                imv1YearAgo.setVisibility(View.GONE);
                imv2YearAgo.setVisibility(View.GONE);
                imv3YearAgo.setVisibility(View.GONE);
                imv4YearAgo.setVisibility(View.GONE);
                imv5YearAgo.setVisibility(View.GONE);
                imvNull.setVisibility(View.GONE);
                break;
            case Constants.SELECT_1_YEAR_AGO:
                imvAllYear.setVisibility(View.GONE);
                imv1YearAgo.setVisibility(View.VISIBLE);
                imv2YearAgo.setVisibility(View.GONE);
                imv3YearAgo.setVisibility(View.GONE);
                imv4YearAgo.setVisibility(View.GONE);
                imv5YearAgo.setVisibility(View.GONE);
                imvNull.setVisibility(View.GONE);
                break;
            case Constants.SELECT_2_YEAR_AGO:
                imvAllYear.setVisibility(View.GONE);
                imv1YearAgo.setVisibility(View.GONE);
                imv2YearAgo.setVisibility(View.VISIBLE);
                imv3YearAgo.setVisibility(View.GONE);
                imv4YearAgo.setVisibility(View.GONE);
                imv5YearAgo.setVisibility(View.GONE);
                imvNull.setVisibility(View.GONE);
                break;
            case Constants.SELECT_3_YEAR_AGO:
                imvAllYear.setVisibility(View.GONE);
                imv1YearAgo.setVisibility(View.GONE);
                imv2YearAgo.setVisibility(View.GONE);
                imv3YearAgo.setVisibility(View.VISIBLE);
                imv4YearAgo.setVisibility(View.GONE);
                imv5YearAgo.setVisibility(View.GONE);
                imvNull.setVisibility(View.GONE);
                break;
            case Constants.SELECT_4_YEAR_AGO:
                imvAllYear.setVisibility(View.GONE);
                imv1YearAgo.setVisibility(View.GONE);
                imv2YearAgo.setVisibility(View.GONE);
                imv3YearAgo.setVisibility(View.GONE);
                imv4YearAgo.setVisibility(View.VISIBLE);
                imv5YearAgo.setVisibility(View.GONE);
                imvNull.setVisibility(View.GONE);
                break;
            case Constants.SELECT_5_YEAR_AGO:
                imvAllYear.setVisibility(View.GONE);
                imv1YearAgo.setVisibility(View.GONE);
                imv2YearAgo.setVisibility(View.GONE);
                imv3YearAgo.setVisibility(View.GONE);
                imv4YearAgo.setVisibility(View.GONE);
                imv5YearAgo.setVisibility(View.VISIBLE);
                imvNull.setVisibility(View.GONE);
                break;
            case Constants.SELECT_NULL:
                imvAllYear.setVisibility(View.GONE);
                imv1YearAgo.setVisibility(View.GONE);
                imv2YearAgo.setVisibility(View.GONE);
                imv3YearAgo.setVisibility(View.GONE);
                imv4YearAgo.setVisibility(View.GONE);
                imv5YearAgo.setVisibility(View.GONE);
                imvNull.setVisibility(View.VISIBLE);
                break;
            default:
                break;
        }
    }

    /**
     * Selected filter data send to activity
     *
     * @param year is number selected
     */
    private void selectedFilter(int year) {

        SelectedFilterDialogListener activity = (SelectedFilterDialogListener) getActivity();
        activity.onSelectedFilterDialog(type, id, nameClassify, year);
        dismiss();
    }

    /**
     * Move back list classify or publisher dialog
     */
    private void loadList() {

        if (typeOld != typeNew) {
            idNew = Constants.ID_ROW_ALL;
        }
        // init new dialog
        DListFragment dListFragment = new DListFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();
        Bundle bundle = new Bundle();
        bundle.putInt(Config.TYPE, typeOld);
        bundle.putInt(Config.TYPE_NEW, typeNew);
        bundle.putInt(Constants.YEAR_AGO, yearOld);
        bundle.putInt(Constants.RANK, rank);
        bundle.putString(Constants.COLUMN_ID, idOld);
        bundle.putString(Constants.COLUMN_ID_NEW, idNew);
        bundle.putString(Constants.COLUMN_LARGE_CLASSIFICATION_NAME, nameClassify);
        dListFragment.setArguments(bundle);
        dListFragment.show(fm, null);
        dismiss();
    }
}
