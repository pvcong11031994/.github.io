package com.android.productchange.fragments;

import android.annotation.SuppressLint;
import android.app.Dialog;
import android.content.DialogInterface;
import android.graphics.Color;
import android.graphics.drawable.ColorDrawable;
import android.os.Bundle;
import android.support.v4.app.DialogFragment;
import android.support.v4.app.FragmentManager;
import android.view.Gravity;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.view.Window;
import android.widget.Button;
import android.widget.ImageButton;
import android.widget.LinearLayout;
import android.widget.TextView;

import com.android.productchange.R;
import com.android.productchange.api.Config;
import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.constants.Message;
import com.android.productchange.db.entity.Users;
import com.android.productchange.db.models.UserModel;
import com.android.productchange.views.CustomCalendarView;

import java.io.File;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.Date;

/**
 * @author tien-lv
 *         Created by tien-lv on 2017/12/19.
 *         Dialog product detail screen
 */

@SuppressWarnings("deprecation")
public class DSelectDateFragment extends DialogFragment implements View.OnClickListener {

    public interface SelectDateDialogListener {
        void onSelectedDateDialog(String idSelected, int typeSelected, String dateFromSelected,
                String dateToSelected, boolean dateChecked);
    }

    LinearLayout llDateFrom, llDateTo;
    TextView txvDateFrom, txvDateTo, txvPath;

    /**
     * ID
     */
    private String id, name;

    /**
     * Rank
     */
    private int rank;

    /**
     * Type
     */
    private int type;

    private String dateFrom, dateTo;

    ImageButton imbBack;
    Button btnDecision;

    boolean dateChecked;
    Date calDateFrom = null, calDateTo = null;

    /**
     * Model user.
     */
    UserModel userModel;

    /**
     * Response create date.
     */
    private String createDate, dateFromFormat, dateToFormat;

    Calendar dateFromDefault = Calendar.getInstance(), dateToDefault = Calendar.getInstance();

    /**
     * Init Dialog Product Detail with custom layout
     */
    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
            Bundle saveInstanceState) {

        // Init custom product detail layout
        View rootView = inflater.inflate(R.layout.layout_select_date, container, false);
        if (getDialog().getWindow() != null) {
            getDialog().getWindow().requestFeature(Window.FEATURE_NO_TITLE);
            getDialog().getWindow().setBackgroundDrawable(new ColorDrawable(Color.TRANSPARENT));
        }

        if (getArguments() != null) {
            id = getArguments().getString(Constants.COLUMN_ID);
            type = getArguments().getInt(Config.TYPE);
            rank = getArguments().getInt(Constants.RANK);
            dateChecked = getArguments().getBoolean(Constants.FLAG_DATE_CHECK);
            dateFrom = getArguments().getString(Constants.COLUMN_DATE_FROM);
            dateTo = getArguments().getString(Constants.COLUMN_DATE_TO);
            createDate = getArguments().getString(Constants.COLUMN_CREATE_DATE);
        }

        llDateFrom = (LinearLayout) rootView.findViewById(R.id.ll_dateFrom);
        llDateTo = (LinearLayout) rootView.findViewById(R.id.ll_dateTo);

        txvDateFrom = (TextView) rootView.findViewById(R.id.txv_dateFrom);
        txvDateTo = (TextView) rootView.findViewById(R.id.txv_dateTo);
        txvPath = (TextView) rootView.findViewById(R.id.txv_path);

        imbBack = (ImageButton) rootView.findViewById(R.id.imb_back);
        btnDecision = (Button) rootView.findViewById(R.id.btn_decision);

        if (createDate == null) {
            userModel = new UserModel();
            Users userInfo = userModel.getUserInfo();
            createDate = userInfo.getCreate_date();
        }

        @SuppressLint("SimpleDateFormat") SimpleDateFormat sdf = new SimpleDateFormat(
                Constants.DATE_FORMAT_STRING);
        Date calDateFrom = null, calDateTo = null;
        try {
            calDateFrom = sdf.parse(createDate);
            calDateTo = sdf.parse(createDate);
        } catch (ParseException e) {
            e.printStackTrace();
        }

        dateFromDefault.setTime(calDateFrom);
        dateToDefault.setTime(calDateTo);
        dateFromDefault.add(Calendar.DATE, Constants.DATE_FROM);
        dateToDefault.add(Calendar.DATE, Constants.DATE_TO);

        SimpleDateFormat df = new SimpleDateFormat(Constants.DATE_FORMAT_STRING);
        dateFromFormat = df.format(dateFromDefault.getTime());
        dateToFormat = df.format(dateToDefault.getTime());

        dateFromFormat = dateFromFormat.substring(0, 4) + File.separator + dateFromFormat.substring(
                4, 6) + File.separator + dateFromFormat.substring(6, 8);
        dateToFormat = dateToFormat.substring(0, 4) + File.separator + dateToFormat.substring(4, 6)
                + File.separator + dateToFormat.substring(6, 8);

        loadData();

        imbBack.setOnClickListener(this);
        btnDecision.setOnClickListener(this);
        txvDateFrom.setOnClickListener(this);
        txvDateTo.setOnClickListener(this);

        return rootView;
    }

    private void loadData() {

        Calendar dateFromDefault = Calendar.getInstance(), dateToDefault = Calendar.getInstance();
        dateFromDefault.add(Calendar.DATE, Constants.DATE_FROM);
        dateToDefault.add(Calendar.DATE, Constants.DATE_TO);

        CustomCalendarView.calendarDateFrom = Calendar.getInstance();
        CustomCalendarView.calendarDateTo = Calendar.getInstance();

        @SuppressLint("SimpleDateFormat") SimpleDateFormat sdf = new SimpleDateFormat(
                Constants.DATE_FORMAT_STRING);
        try {
            calDateFrom = sdf.parse(dateFrom);
            calDateTo = sdf.parse(dateTo);
        } catch (ParseException e) {
            e.printStackTrace();
        }

        if (calDateFrom == null || calDateTo == null) {
            CustomCalendarView.calendarDateFrom = dateFromDefault;
            CustomCalendarView.calendarDateTo = dateToDefault;
        } else {
            CustomCalendarView.calendarDateFrom.setTime(calDateFrom);
            CustomCalendarView.calendarDateTo.setTime(calDateTo);
        }

        txvDateFrom.setText(
                dateFrom.substring(0, 4) + File.separator + dateFrom.substring(4, 6)
                        + File.separator + dateFrom.substring(6, 8));


        txvDateTo.setText(dateTo.substring(0, 4) + File.separator + dateTo.substring(4, 6)
                + File.separator + dateTo.substring(6, 8));

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

    @Override
    public void onClick(View v) {
        Bundle bundle = new Bundle();
        DSelectFragment dSelectFragment = new DSelectFragment();
        DCalendarFragment dCalendarFragment = new DCalendarFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();

        Date dateFromSelected = CustomCalendarView.calendarDateFrom.getTime();
        Date dateToSelected = CustomCalendarView.calendarDateTo.getTime();

        switch (v.getId()) {
            case R.id.imb_back:
                bundle.putInt(Config.TYPE, type);
                bundle.putString(Constants.COLUMN_ID, id);
                bundle.putInt(Constants.RANK, rank);
                bundle.putBoolean(Constants.FLAG_DATE_CHECK, dateChecked);
                bundle.putString(Constants.COLUMN_DATE_FROM, dateFrom);
                bundle.putString(Constants.COLUMN_DATE_TO, dateTo);
                dSelectFragment.setArguments(bundle);
                dSelectFragment.show(fm, null);
                dismiss();
                break;
            case R.id.txv_dateFrom:
                CustomCalendarView.flag = true;
                bundle.putInt(Config.TYPE, type);
                bundle.putString(Constants.COLUMN_ID, id);
                bundle.putInt(Constants.RANK, rank);
                bundle.putString(Constants.COLUMN_DATE_FROM, dateFrom);
                bundle.putString(Constants.COLUMN_DATE_TO, dateTo);
                bundle.putString(Constants.COLUMN_CREATE_DATE, createDate);
                bundle.putBoolean(Constants.FLAG_DATE, true);
                dCalendarFragment.setArguments(bundle);
                dCalendarFragment.show(fm, null);
                dismiss();
                break;
            case R.id.txv_dateTo:
                CustomCalendarView.flag = false;
                bundle.putInt(Config.TYPE, type);
                bundle.putString(Constants.COLUMN_ID, id);
                bundle.putInt(Constants.RANK, rank);
                bundle.putString(Constants.COLUMN_DATE_FROM, dateFrom);
                bundle.putString(Constants.COLUMN_DATE_TO, dateTo);
                bundle.putString(Constants.COLUMN_CREATE_DATE, createDate);
                bundle.putBoolean(Constants.FLAG_DATE, false);
                dCalendarFragment.setArguments(bundle);
                dCalendarFragment.show(fm, null);
                dismiss();
                break;
            case R.id.btn_decision:

                if (!dateChecked) {
                    dateChecked = true;
                }

                if (compareTo(dateFromSelected, dateFromDefault.getTime()) >= 0 && compareTo(
                        dateToSelected, dateToDefault.getTime()) <= 0) {
                    SelectDateDialogListener activity =
                            (SelectDateDialogListener) getActivity();
                    activity.onSelectedDateDialog(id, type, dateFrom, dateTo, dateChecked);
                    dismiss();
                } else {
                    showDialog();
                }

                break;
        }

    }

    /**
     * Dialog show
     */

    private void showDialog() {
        android.support.v7.app.AlertDialog.Builder dialog =
                new android.support.v7.app.AlertDialog.Builder(getActivity());
        dialog.setCancelable(false);

        dialog
                .setMessage(String.format(Message.MESSAGE_CALENDAR, dateFromFormat, dateToFormat))
                .setNegativeButton(getString(R.string.ok_button),
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                dialog.dismiss();
                            }
                        });

        android.support.v7.app.AlertDialog alert = dialog.show();
        TextView messageText = (TextView) alert.findViewById(android.R.id.message);
        assert messageText != null;
        messageText.setGravity(Gravity.CENTER);
    }

    /**
     * Compare 2 date
     *
     * @param date1 {@link Date}
     * @param date2 {@link Date}
     * @return long
     */
    private long compareTo(Date date1, Date date2) {
        //returns negative value if date1 is before date2
        //returns 0 if dates are even
        //returns positive value if date1 is after date2
        return date1.getTime() - date2.getTime();
    }
}
