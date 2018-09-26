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
import android.widget.ImageButton;
import android.widget.TextView;

import com.android.productchange.R;
import com.android.productchange.api.Config;
import com.android.productchange.common.constants.Constants;
import com.android.productchange.interfaces.OnDateSelectedListener;
import com.android.productchange.objects.CalendarDate;
import com.android.productchange.views.CustomCalendarView;

/**
 * @author tien-lv
 *         Created by tien-lv on 2017/12/19.
 *         Dialog product detail screen
 */

@SuppressWarnings("deprecation")
public class DCalendarFragment extends DialogFragment implements OnDateSelectedListener,
        View.OnClickListener {

    public interface RankDialogListener {
        void onRankSelectedDialog(int itemSelected);
    }

    private CustomCalendarView mCustomCalendar;

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
    boolean flagDate;

    ImageButton imbBack;
    TextView txvPath;

    /**
     * Init Dialog Product Detail with custom layout
     */
    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
            Bundle saveInstanceState) {

        // Init custom product detail layout
        View rootView = inflater.inflate(R.layout.layout_calendar, container, false);
        if (getDialog().getWindow() != null) {
            getDialog().getWindow().requestFeature(Window.FEATURE_NO_TITLE);
            getDialog().getWindow().setBackgroundDrawable(new ColorDrawable(Color.TRANSPARENT));
        }

        if (getArguments() != null) {
            id = getArguments().getString(Constants.COLUMN_ID);
            type = getArguments().getInt(Config.TYPE);
            rank = getArguments().getInt(Constants.RANK);
            dateFrom = getArguments().getString(Constants.COLUMN_DATE_FROM);
            dateTo = getArguments().getString(Constants.COLUMN_DATE_TO);
            flagDate = getArguments().getBoolean(Constants.FLAG_DATE);
        }

        mCustomCalendar = (CustomCalendarView) rootView.findViewById(R.id.custom_calendar);
        mCustomCalendar.setOnDateSelectedListener(this);

        imbBack = (ImageButton) rootView.findViewById(R.id.imb_back);
        txvPath = (TextView) rootView.findViewById(R.id.txv_path);

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

    @Override
    public void onDateSelected(CalendarDate date) {

        String formattedDate =
                date.yearToString() + date.monthToString() + date.dayToString();
        if (flagDate) {
            dateFrom = formattedDate;
        } else {
            dateTo = formattedDate;
        }
        DSelectDateFragment dSelectDateFragment = new DSelectDateFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();
        Bundle bundle = new Bundle();
        bundle.putInt(Config.TYPE, type);
        bundle.putString(Constants.COLUMN_ID, id);
        bundle.putInt(Constants.RANK, rank);
        bundle.putString(Constants.COLUMN_DATE_FROM, dateFrom);
        bundle.putString(Constants.COLUMN_DATE_TO, dateTo);
        dSelectDateFragment.setArguments(bundle);
        dSelectDateFragment.show(fm, null);
        dismiss();
    }

    @Override
    public void onClick(View v) {
        Bundle bundle = new Bundle();
        DSelectDateFragment dListFragment = new DSelectDateFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();
        switch (v.getId()) {
            case R.id.imb_back:
                bundle.putInt(Config.TYPE, type);
                bundle.putString(Constants.COLUMN_ID, id);
                bundle.putString(Constants.COLUMN_NAME, name);
                bundle.putInt(Constants.RANK, rank);
                bundle.putString(Constants.COLUMN_DATE_FROM, dateFrom);
                bundle.putString(Constants.COLUMN_DATE_TO, dateTo);
                dListFragment.setArguments(bundle);
                dListFragment.show(fm, null);
                dismiss();
                break;
        }
    }
}
