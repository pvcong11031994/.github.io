package com.android.productchange.fragments;

import android.app.Dialog;
import android.content.DialogInterface;
import android.graphics.Color;
import android.graphics.drawable.ColorDrawable;
import android.os.Bundle;
import android.support.v4.app.DialogFragment;
import android.support.v4.app.FragmentManager;
import android.view.KeyEvent;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.view.Window;
import android.widget.AdapterView;
import android.widget.ImageButton;
import android.widget.ListView;
import android.widget.TextView;

import com.android.productchange.R;
import com.android.productchange.adapters.ListViewMonthReturnBookAdapter;
import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.utils.Common;
import com.android.productchange.common.utils.FlagSettingNew;
import com.android.productchange.common.utils.FlagSettingOld;

import java.util.ArrayList;
import java.util.HashMap;

/**
 * <h1>List select Dialog</h1>
 * <p>
 * Dialog list select screen
 *
 * @author cong-pv
 * @since 2018-07-09.
 */
@SuppressWarnings("deprecation")
public class DFilterListStocksPercentFragment extends DialogFragment implements View.OnClickListener {

    /**
     * Interface to item selected to activity
     */
    public interface ItemSelectedStocksPercentDialogListener {

        /**
         * Function send list selected data to activity
         */
        void onLitSelectedStocksPercentDialog(FlagSettingNew flagSettingNew);
    }

    /**
     * List data for Item List View
     */
    private ArrayList<HashMap<String, String>> list;

    /**
     * List view
     */
    private ListView lsvList;

    /**
     * Button back
     */
    private ImageButton imbBack;

    private FlagSettingNew flagSettingNew;
    private FlagSettingOld flagSettingOld;
    private Common common;

    private TextView txvHeaderFilter;
    private String flagSwitchOCR;

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
        View rootView = inflater.inflate(R.layout.layout_genre_return_book, container, false);
        if (getDialog().getWindow() != null) {
            getDialog().getWindow().requestFeature(Window.FEATURE_NO_TITLE);
            getDialog().getWindow().setBackgroundDrawable(new ColorDrawable(Color.TRANSPARENT));
        }
        flagSettingNew = new FlagSettingNew();
        flagSettingOld = new FlagSettingOld();
        common = new Common();

        //Header filter set text
        txvHeaderFilter = (TextView) rootView.findViewById(R.id.txv_header_filter);
        txvHeaderFilter.setText(Constants.HEADER_STOCKS_PERCENT);

        if (getArguments() != null) {
            common.SetArgumentsFlagData(flagSettingNew, flagSettingOld, getArguments());
            //get flag switch OCR
           // flagSwitchOCR = getArguments().getString(Constants.FLAG_SWITCH_OCR);
        }

        // list group1 cd
        lsvList = (ListView) rootView.findViewById(R.id.lsv_list);
        list = new ArrayList<>();

        //show item list group 1 cd
        loadListItemStocksPercent(flagSettingNew.getFlagStockPercent());

        imbBack = (ImageButton) rootView.findViewById(R.id.imb_back);

        lsvList.setOnItemClickListener(new AdapterView.OnItemClickListener() {
            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {

                flagSettingNew.setFlagStockPercent(list.get(position).get(Constants.SELECT_POISITION));
                flagSettingNew.setFlagStockPercentShowScreen(list.get(position).get(Constants.SELECT_VALUE));
                // move to select year dialog
                backDialogReturn();
            }
        });

        imbBack.setOnClickListener(this);
        //Event lick back device
        getDialog().setOnKeyListener(new DialogInterface.OnKeyListener() {
            @Override
            public boolean onKey(DialogInterface dialog, int keyCode, KeyEvent event) {
                if (keyCode == KeyEvent.KEYCODE_BACK) {
                    backDialogReturn();
                    return true;
                }
                return false;
            }
        });
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
     * Get data filter stocks percent
     */
    private void loadListItemStocksPercent(String positionSelectStocksPercent) {

        // Load data from .xml
        final String arrReleaseDate[] = getResources().getStringArray(R.array.stocks_percent);

        // Set data into list adapter
        for (int i = 0; i < arrReleaseDate.length; i++) {
            HashMap<String, String> hashMap = new HashMap<>();
            hashMap.put(Constants.SELECT_POISITION, String.valueOf(i + 1));
            hashMap.put(Constants.SELECT_VALUE, String.valueOf(arrReleaseDate[i]));
            list.add(hashMap);
        }
        // Adapter init
        // Set data adapter to list view
        ListViewMonthReturnBookAdapter adapter = new ListViewMonthReturnBookAdapter(getActivity(), list, positionSelectStocksPercent);
        lsvList.setAdapter(adapter);
    }


    /**
     * On click
     *
     * @param v is View on click listener
     */
    @Override
    public void onClick(View v) {
        if (v.getId() == R.id.imb_back) {
            //TODO
            backDialogReturn();
        }
    }

    /**
     * Back dialog filter one
     */
    private void backDialogReturn() {

        DFilterSettingFragment dFilterSettingFragment = new DFilterSettingFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();
        Bundle bundle = common.DataPutActivity(flagSettingNew, flagSettingOld);

        dFilterSettingFragment.setArguments(bundle);
        dFilterSettingFragment.show(fm, null);
        dismiss();
    }
}
