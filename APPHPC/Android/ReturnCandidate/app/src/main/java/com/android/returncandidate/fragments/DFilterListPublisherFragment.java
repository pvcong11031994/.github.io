package com.android.returncandidate.fragments;

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
import android.widget.ImageButton;
import android.widget.ListView;
import android.widget.TextView;

import com.android.returncandidate.R;
import com.android.returncandidate.adapters.ListViewPublisherAdapter;
import com.android.returncandidate.common.constants.Constants;
import com.android.returncandidate.common.utils.Common;
import com.android.returncandidate.common.utils.FlagSettingNew;
import com.android.returncandidate.common.utils.FlagSettingOld;
import com.android.returncandidate.db.entity.Publisher;
import com.android.returncandidate.db.models.PublisherModel;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

/**
 * <h1>List select Dialog</h1>
 * <p>
 * Dialog list select screen
 *
 * @author cong-pv
 * @since 2018/07/09.
 */
@SuppressWarnings("deprecation")
public class DFilterListPublisherFragment extends DialogFragment implements View.OnClickListener {

    /**
     * Interface to item selected to activity
     */
    public interface ItemSelectedPublisherDialogListener {

        /**
         * Function send list selected data to activity
         */
        void onLitSelectedPublisherDialog(FlagSettingNew flagSettingNew);
    }

    /**
     * List data for Item List View
     */
    private ArrayList<HashMap<String, String>> list;

    private PublisherModel publisherModel;

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
    private String flagSwitchOCR;
    private TextView txvHeaderFilter;

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
        View rootView = inflater.inflate(R.layout.layout_genre, container, false);
        if (getDialog().getWindow() != null) {
            getDialog().getWindow().requestFeature(Window.FEATURE_NO_TITLE);
            getDialog().getWindow().setBackgroundDrawable(new ColorDrawable(Color.TRANSPARENT));
        }
        flagSettingNew = new FlagSettingNew();
        flagSettingOld = new FlagSettingOld();
        common = new Common();

        //Header filter set text
        txvHeaderFilter = (TextView) rootView.findViewById(R.id.txv_header_filter);
        txvHeaderFilter.setText(Constants.HEADER_PUBLISHER);
        if (getArguments() != null) {
            common.SetArgumentsFlagData(flagSettingNew, flagSettingOld, getArguments());
            //get flag switch OCR
            flagSwitchOCR = getArguments().getString(Constants.FLAG_SWITCH_OCR);
        }

        lsvList = (ListView) rootView.findViewById(R.id.lsv_list);
        list = new ArrayList<>();
        publisherModel = new PublisherModel();

        //loadItems(id);
        //Check select group cd
        loadItemsPublisher(flagSettingNew);

        imbBack = (ImageButton) rootView.findViewById(R.id.imb_back);

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
     * Set data for list item of Publisher
     */
    public void loadItemsPublisher(FlagSettingNew flagSettingNew) {

        //Get info publisher from classification
        List<Publisher> rlist = publisherModel.getInfoPublisher(flagSettingNew);

        //Check all or not all
        int intcheck = 0;
        for (int i = 0; i < rlist.size(); i++) {
            for (int j = 0; j < flagSettingNew.getFlagPublisher().size(); j++) {
                if (rlist.get(i).getId().equals(flagSettingNew.getFlagPublisher().get(j))) {
                    intcheck++;
                    break;
                }
            }
        }
        HashMap<String, String> hashMapAll = new HashMap<>();
        hashMapAll.put(Constants.COLUMN_PUBLISHER_CD, Constants.ID_ROW_ALL);
        hashMapAll.put(Constants.COLUMN_PUBLISHER_NAME, Constants.ROW_ALL);
        for (String value : flagSettingNew.getFlagPublisher()) {
            if (Constants.ID_ROW_ALL.equals(value) || intcheck == rlist.size()) {
                hashMapAll.put(Constants.FLAG_SELECT, Constants.VALUE_STR_CHECK);
                break;
            } else {
                hashMapAll.put(Constants.FLAG_SELECT, Constants.VALUE_STR_NO_CHECK);
            }
        }
        list.add(hashMapAll);

        // Set data into list adapter
        for (int i = 0; i < rlist.size(); i++) {
            HashMap<String, String> hashMap = new HashMap<>();
            hashMap.put(Constants.COLUMN_PUBLISHER_CD, String.valueOf(rlist.get(i).getId()));
            hashMap.put(Constants.COLUMN_PUBLISHER_NAME, rlist.get(i).getName());
            for (String value : flagSettingNew.getFlagPublisher()) {
                if (Constants.ID_ROW_ALL.equals(value)) {
                    hashMap.put(Constants.FLAG_SELECT, Constants.VALUE_STR_CHECK);
                    break;
                } else {
                    if (rlist.get(i).getId().equals(value)) {
                        hashMap.put(Constants.FLAG_SELECT, Constants.VALUE_STR_CHECK);
                        break;
                    } else {
                        hashMap.put(Constants.FLAG_SELECT, Constants.VALUE_STR_NO_CHECK);
                    }
                }
            }
            list.add(hashMap);
        }

        // Adapter init
        // Set data adapter to list view
        ListViewPublisherAdapter adapter = new ListViewPublisherAdapter(getActivity(), list);
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
            backDialogReturn();
        }
    }

    /**
     * Back dialog setting
     */
    private void backDialogReturn() {

        flagSettingNew = putToListItemCheckToFlagSettingNew();
        DFilterSettingFragment dFilterSettingFragment = new DFilterSettingFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();
        Bundle bundle = common.DataPutActivity(flagSettingNew, flagSettingOld);
        //put flag switch OCR
        bundle.putString(Constants.FLAG_SWITCH_OCR, flagSwitchOCR);

        dFilterSettingFragment.setArguments(bundle);
        dFilterSettingFragment.show(fm, null);
        dismiss();
    }

    private FlagSettingNew putToListItemCheckToFlagSettingNew() {

        ArrayList<String> arrayList = new ArrayList<>();
        ArrayList<String> arrayListName = new ArrayList<>();
        for (int i = 0; i < list.size(); i++) {
            if (Constants.VALUE_STR_CHECK.equals(list.get(i).get(Constants.FLAG_SELECT))) {
                //flagSettingNew.setFlagPublisher();
                arrayList.add(list.get(i).get(Constants.COLUMN_PUBLISHER_CD));
                arrayListName.add(list.get(i).get(Constants.COLUMN_PUBLISHER_NAME));
            }
        }
        if (arrayList.size() == 0) {
            for (int i = 0; i < list.size(); i++) {
                //flagSettingNew.setFlagPublisher();
                arrayList.add(list.get(i).get(Constants.COLUMN_PUBLISHER_CD));
                arrayListName.add(list.get(i).get(Constants.COLUMN_PUBLISHER_NAME));
            }
        }
        flagSettingNew.setFlagPublisher(arrayList);
        flagSettingNew.setFlagPublisherShowScreen(arrayListName);
        return flagSettingNew;
    }
}
