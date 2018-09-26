package com.android.productchange.fragments;

import android.app.Dialog;
import android.graphics.Color;
import android.graphics.drawable.ColorDrawable;
import android.icu.text.DecimalFormat;
import android.icu.text.SimpleDateFormat;
import android.os.Bundle;
import android.support.v4.app.DialogFragment;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.view.Window;
import android.widget.ArrayAdapter;
import android.widget.TextView;

import com.android.productchange.R;
import com.android.productchange.common.constants.Constants;

import java.text.ParseException;

/**
 * Created by cong-pv on 2018/06/23.
 * Dialog product detail screen
 */

public class DProductDetailFragment extends DialogFragment {

    /**
     * Init array item in comboBox
     */
    String arr[] = new String[11];

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState) {

        //Init custom product detail layout
        View rootView = inflater.inflate(R.layout.product_detail, container, false);
        if (getDialog().getWindow() != null) {
            getDialog().getWindow().requestFeature(Window.FEATURE_NO_TITLE);
            getDialog().getWindow().setBackgroundDrawable(new ColorDrawable(Color.TRANSPARENT));
        }

        //Get Id textview Product detailk
        TextView txv_jan_cd = (TextView) rootView.findViewById(R.id.txv_jan_cd);
        TextView txv_group1_name = (TextView) rootView.findViewById(R.id.txv_group1_name);
        TextView txv_group2_name = (TextView) rootView.findViewById(R.id.txv_group2_name);
        TextView txv_location_id = (TextView) rootView.findViewById(R.id.txv_location_id);
        TextView txv_product_name = (TextView) rootView.findViewById(R.id.txv_product_name);
        TextView txv_writer_name = (TextView) rootView.findViewById(R.id.txv_writer_name);
        TextView txv_publisher_name = (TextView) rootView.findViewById(R.id.txv_publisher_name);
        TextView txv_publish_date = (TextView) rootView.findViewById(R.id.txv_publish_date);
        TextView txv_price = (TextView) rootView.findViewById(R.id.txv_price);
        TextView txv_inventory_number = (TextView) rootView.findViewById(R.id.txv_inventory_number);
        TextView txv_first_supply_date = (TextView) rootView.findViewById(R.id.txv_first_supply_date);
        TextView txv_last_supply_date = (TextView) rootView.findViewById(R.id.txv_last_supply_date);
        TextView txv_last_sales_date = (TextView) rootView.findViewById(R.id.txv_last_sales_date);
        TextView txv_last_order_date = (TextView) rootView.findViewById(R.id.txv_last_order_date);
        TextView txv_total_sales = (TextView) rootView.findViewById(R.id.txv_total_sales);
        TextView txv_total_supply = (TextView) rootView.findViewById(R.id.txv_total_supply);
        TextView txv_total_return = (TextView) rootView.findViewById(R.id.txv_total_return);
        TextView txv_year_rank = (TextView) rootView.findViewById(R.id.txv_year_rank);
        TextView txv_joubi = (TextView) rootView.findViewById(R.id.txv_joubi);
        Bundle bundle = getArguments();
        if (bundle != null) {
            //Show year_rank
            String valueYearRank = bundle.getString(Constants.COLUMN_YEAR_RANK);
            String showYearRank;
            if (Constants.VALUE_MAX_YEAR_RANK.equals(valueYearRank)) {
                showYearRank = Constants.SHOW_MAX_YEAR_RANK;
            } else {
                showYearRank = String.format("%s位/%s中", formatMoney(valueYearRank),
                        formatMoney(String.valueOf(bundle.getInt(Constants.COLUMN_MAX_YEAR_RANK))));
            }
            txv_group1_name.setText(bundle.getString(Constants.COLUMN_MEDIA_GROUP1_NAME));
            txv_group2_name.setText(bundle.getString(Constants.COLUMN_MEDIA_GROUP2_NAME));
            txv_jan_cd.setText(bundle.getString(Constants.COLUMN_JAN_CD));
            txv_product_name.setText(bundle.getString(Constants.COLUMN_GOODS_NAME));
            txv_writer_name.setText(bundle.getString(Constants.COLUMN_WRITER_NAME));
            txv_publisher_name.setText(bundle.getString(Constants.COLUMN_PUBLISHER_NAME_RETURN));
            txv_publish_date.setText(
                    formatDate(bundle.getString(Constants.COLUMN_SALES_DATE, Constants.BLANK)));
            if (Constants.BLANK.equals(bundle.getString(Constants.COLUMN_PRICE))) {
                txv_price.setText(Constants.BLANK);

            } else {
                txv_price.setText(String.format("%s%s", Constants.SYMBOL,
                        formatMoney((bundle.getString(Constants.COLUMN_PRICE)))));
            }
            if (Constants.NULL.equals(bundle.getString(Constants.COLUMN_STOCK_COUNT))) {
                txv_inventory_number.setText(Constants.BLANK);
            } else {
                txv_inventory_number.setText(
                        String.valueOf(bundle.getString(Constants.COLUMN_STOCK_COUNT)));
            }
            txv_first_supply_date.setText(
                    formatDate(bundle.getString(Constants.COLUMN_FIRST_SUPPLY_DATE,
                            Constants.BLANK)));
            txv_last_supply_date.setText(
                    formatDate(bundle.getString(Constants.COLUMN_LAST_SUPPLY_DATE,
                            Constants.BLANK)));
            txv_last_sales_date.setText(
                    formatDate(bundle.getString(Constants.COLUMN_LAST_SALES_DATE,
                            Constants.BLANK)));
            txv_last_order_date.setText(
                    formatDate(bundle.getString(Constants.COLUMN_LAST_ORDER_DATE,
                            Constants.BLANK)));
            txv_year_rank.setText(showYearRank);
            txv_joubi.setText(formatJoubi(bundle.getString(Constants.COLUMN_JOUBI)));
            txv_total_sales.setText(formatMoney(bundle.getString(Constants.COLUMN_TOTAL_SALES)));
            txv_total_supply.setText(formatMoney(bundle.getString(Constants.COLUMN_TOTAL_SUPPLY)));
            txv_total_return.setText(formatMoney(bundle.getString(Constants.COLUMN_TOTAL_RETURN)));
            txv_location_id.setText(bundle.getString(Constants.COLUMN_LOCATION_ID));
        }
        loadItems();

        // Set adapter for combo box with array
        ArrayAdapter<String> adapter1 = new ArrayAdapter<>(
                getContext(), android.R.layout.simple_list_item_1, arr
        );

        // Set item selected
        adapter1.setDropDownViewResource(android.R.layout.simple_list_item_single_choice);

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

    /**
     * Set item for combo box
     */
    public void loadItems() {
        for (int i = 0; i < 11; i++) {
            arr[i] = String.valueOf(i);
        }
    }

    /**
     * Format String date
     *
     * @param date date
     */
    public String formatDate(String date) {

        if (Constants.VALUE_DEFAULT_DATE.equals(date)) {
            return Constants.BLANK;
        } else if (Constants.VALUE_ERROR_DATE.equals(date)) {
            return Constants.VALUE_ERROR_DATE_FORMAT;
        }
        String result = date;
        SimpleDateFormat inputUser = new SimpleDateFormat("yyyyMMdd");
        SimpleDateFormat resultFormat = new SimpleDateFormat("yyyy/MM/dd");
        try {
            result = resultFormat.format(inputUser.parse(date));
        } catch (ParseException e) {
            e.printStackTrace();
        }
        return result;
    }

    public String formatMoney(String money) {

        //Convert String to Float
        String strFormat = money;
        String pattern = "#,###.#";
        DecimalFormat decimalFormat = new DecimalFormat(pattern);
        try {
            strFormat = decimalFormat.format(Float.parseFloat(money));
        } catch (Exception e) {
            e.printStackTrace();
        }
        return strFormat;
    }

    private String formatJoubi(String joubi) {

        if (Constants.VALUE_JOUBI.equals(joubi)) {
            return Constants.VALUE_JOUBI_SHOW;
        }
        return Constants.BLANK;
    }
}