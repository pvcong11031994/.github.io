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
import com.android.returncandidate.adapters.ListViewGenre2Adapter;
import com.android.returncandidate.common.constants.Constants;
import com.android.returncandidate.common.utils.Common;
import com.android.returncandidate.common.utils.FlagSettingNew;
import com.android.returncandidate.common.utils.FlagSettingOld;
import com.android.returncandidate.db.entity.Books;
import com.android.returncandidate.db.entity.CLP;
import com.android.returncandidate.db.models.BookModel;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

/**
 * <h1>List select Dialog</h1>
 * <p>
 * Dialog list select screen
 *
 * @author cong-pv
 * @since 2018/07/10.
 */
@SuppressWarnings("deprecation")
public class DFilterListGenre2Fragment extends DialogFragment implements View.OnClickListener {

    /**
     * Interface to item selected to activity
     */
    public interface ItemSelectedListGenre2DialogListener {

        /**
         * Function send list selected data to activity
         */
        void onLitSelectedListGenre2Dialog(FlagSettingNew flagSettingNew);
    }

    /**
     * List data for Item List View
     */
    private ArrayList<HashMap<String, String>> list;

    /**
     * Book model
     */
    private BookModel bookModel;

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
    private String flagGroup1CdNew;
    private String flagGroup1NameNew;

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
        View rootView = inflater.inflate(R.layout.layout_genre2, container, false);
        if (getDialog().getWindow() != null) {
            getDialog().getWindow().requestFeature(Window.FEATURE_NO_TITLE);
            getDialog().getWindow().setBackgroundDrawable(new ColorDrawable(Color.TRANSPARENT));
        }
        flagSettingNew = new FlagSettingNew();
        flagSettingOld = new FlagSettingOld();
        common = new Common();

        if (getArguments() != null) {
            common.SetArgumentsFlagData(flagSettingNew, flagSettingOld, getArguments());
            //get group1 select new
            flagGroup1CdNew = getArguments().getString(Constants.FLAG_GROUP1_CD);
            flagGroup1NameNew = getArguments().getString(Constants.FLAG_GROUP1_NAME);
            //get flag switch OCR
            flagSwitchOCR = getArguments().getString(Constants.FLAG_SWITCH_OCR);
        }

        //Header filter set text
        txvHeaderFilter = (TextView) rootView.findViewById(R.id.txv_header_filter_2);
        txvHeaderFilter.setText(Constants.HEADER_CLASSIFICATION_2 + " (" + flagGroup1NameNew + ")");
        // list genre
        lsvList = (ListView) rootView.findViewById(R.id.lsv_list_2);
        list = new ArrayList<>();
        bookModel = new BookModel();

        //show item list genre
        loadListItemGroup2Cd(flagGroup1CdNew, flagSettingNew);

        imbBack = (ImageButton) rootView.findViewById(R.id.imb_back);

        imbBack.setOnClickListener(this);

        //Event click back device
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
     * Get data filter group2_cd with table local book
     */
    private void loadListItemGroup2Cd(String flagGroup1Cd, FlagSettingNew flagSettingNew) {

        List<CLP> rlist = bookModel.getInfoGroupCd2(flagGroup1Cd);

        //Check all or not all
        int intcheck = 0;
        for (int i = 0; i < rlist.size(); i++) {
            for (int j = 0; j < flagSettingNew.getFlagClassificationGroup2Cd().size(); j++) {
                if (rlist.get(i).getId().equals(flagSettingNew.getFlagClassificationGroup2Cd().get(j))) {
                    intcheck++;
                    break;
                }
            }
        }

        HashMap<String, String> hashMapAll = new HashMap<>();
        if (rlist.size() > 1) {
            hashMapAll.put(Constants.COLUMN_MEDIA_GROUP2_CD, Constants.ID_ROW_ALL);
            hashMapAll.put(Constants.COLUMN_MEDIA_GROUP2_NAME, Constants.ROW_ALL);
            for (String valueGroup2Cd : flagSettingNew.getFlagClassificationGroup2Cd()) {
                if (Constants.ID_ROW_ALL.equals(valueGroup2Cd) || intcheck == rlist.size()) {
                    hashMapAll.put(Constants.FLAG_SELECT, Constants.VALUE_STR_CHECK);
                } else {
                    hashMapAll.put(Constants.FLAG_SELECT, Constants.VALUE_STR_NO_CHECK);
                }
            }
            list.add(hashMapAll);
        }
        // Set data into list adapter -- Check select or not // TODO
        if (rlist.size() == 1) {
            HashMap<String, String> hashMap = new HashMap<>();
            hashMap.put(Constants.COLUMN_MEDIA_GROUP2_CD, String.valueOf(rlist.get(0).getId()));
            hashMap.put(Constants.COLUMN_MEDIA_GROUP2_NAME, rlist.get(0).getName());
            hashMap.put(Constants.FLAG_SELECT, Constants.VALUE_STR_CHECK);
            list.add(hashMap);
        } else {
            if (flagSettingNew.getFlagClassificationGroup2Cd().size() == 1 &&
                    Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup2Cd().get(0))) {
                for (int i = 0; i < rlist.size(); i++) {
                    HashMap<String, String> hashMap = new HashMap<>();
                    hashMap.put(Constants.COLUMN_MEDIA_GROUP2_CD, String.valueOf(rlist.get(i).getId()));
                    hashMap.put(Constants.COLUMN_MEDIA_GROUP2_NAME, rlist.get(i).getName());
                    hashMap.put(Constants.FLAG_SELECT, Constants.VALUE_STR_CHECK);
                    list.add(hashMap);
                }
            } else {
                for (int i = 0; i < rlist.size(); i++) {
                    boolean check = false;
                    HashMap<String, String> hashMap = new HashMap<>();
                    hashMap.put(Constants.COLUMN_MEDIA_GROUP2_CD, String.valueOf(rlist.get(i).getId()));
                    hashMap.put(Constants.COLUMN_MEDIA_GROUP2_NAME, rlist.get(i).getName());
                    for (int j = 0; j < flagSettingNew.getFlagClassificationGroup2Cd().size(); j++) {
                        if (rlist.get(i).getId().equals(flagSettingNew.getFlagClassificationGroup2Cd().get(j))) {
                            check = true;
                            break;
                        } else {
                            check = false;
                        }
                    }
                    if (check) {
                        hashMap.put(Constants.FLAG_SELECT, Constants.VALUE_STR_CHECK);
                    } else {
                        hashMap.put(Constants.FLAG_SELECT, Constants.VALUE_STR_NO_CHECK);
                    }
                    list.add(hashMap);
                }
            }
        }

        //If no check -> select all
        Boolean checks = true;
        for (int i = 0; i < list.size(); i++) {
            if (Constants.VALUE_STR_CHECK.equals(list.get(i).get(Constants.FLAG_SELECT))) {
                checks = false;
            }
        }
        if (checks) {
            for (int i = 0; i < list.size(); i++) {
                list.get(i).put(Constants.FLAG_SELECT, Constants.VALUE_STR_CHECK);
            }
        }
        // Adapter init
        // Set data adapter to list view
        ListViewGenre2Adapter adapter = new ListViewGenre2Adapter(getActivity(), list);
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
     * Back dialog group 1 cd
     */
    private void backDialogReturn() {

        flagSettingNew = putToListItemCheckToFlagSettingNew();
        flagSettingNew = checkListGroup2Cd();
        DFilterListGenreFragment dFilterListGenreFragment = new DFilterListGenreFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();
        //set old group1
        //flagSettingNew.setFlagClassificationGroup1Cd(flagSettingNew.getFlagClassificationGroup1CdOld());
        //flagSettingNew.setFlagClassificationGroup1Name(flagSettingNew.getFlagClassificationGroup1NameOld());

        Bundle bundle = common.DataPutActivity(flagSettingNew, flagSettingOld);
        //put flag switch OCR
        bundle.putString(Constants.FLAG_SWITCH_OCR, flagSwitchOCR);

        dFilterListGenreFragment.setArguments(bundle);
        dFilterListGenreFragment.show(fm, null);
        dismiss();
    }

    //Check list group2 cd when back
    private FlagSettingNew checkListGroup2Cd() {

        ArrayList<String> arrayList = new ArrayList<>();
        ArrayList<String> arrayListName = new ArrayList<>();
        Boolean check = false;
        for (int j = 0; j < flagSettingNew.getFlagClassificationGroup2Cd().size(); j++) {
            for (int i = 0; i < list.size(); i++) {
                if ((list.get(i).get(Constants.COLUMN_MEDIA_GROUP2_CD).equals(flagSettingNew.getFlagClassificationGroup2Cd().get(j)))
                        && (Constants.VALUE_STR_NO_CHECK.equals(list.get(i).get(Constants.FLAG_SELECT)))) {
                    check = true;
                    break;
                } else {
                    check = false;
                }
            }
            if (!check) {
                arrayList.add(flagSettingNew.getFlagClassificationGroup2Cd().get(j));
                arrayListName.add(flagSettingNew.getFlagClassificationGroup2Name().get(j));
            }
        }

        flagSettingNew.setFlagClassificationGroup2Cd(arrayList);
        flagSettingNew.setFlagClassificationGroup2Name(arrayListName);
        return flagSettingNew;
    }

    private FlagSettingNew putToListItemCheckToFlagSettingNew() {

        //Boolean check select all
        int flagCheckSelectAll = -1;

        ArrayList<String> arrayList = new ArrayList<>();
        ArrayList<String> arrListEnd = new ArrayList<>();
        ArrayList<String> arrListNameEnd = new ArrayList<>();
        ArrayList<String> arrayListGroup1Cd = new ArrayList<>();
        ArrayList<String> arrayListGroup1Name = new ArrayList<>();

        //Get list new check group2 cd
        for (int i = 0; i < list.size(); i++) {
            if (Constants.VALUE_STR_CHECK.equals(list.get(i).get(Constants.FLAG_SELECT))) {
                //flagSettingNew.setFlagPublisher();
                arrayList.add(list.get(i).get(Constants.COLUMN_MEDIA_GROUP2_CD));
                arrListEnd.add(list.get(i).get(Constants.COLUMN_MEDIA_GROUP2_CD));
                arrListNameEnd.add(list.get(i).get(Constants.COLUMN_MEDIA_GROUP2_NAME));
            }
        }

        //Check list new into list old
        for (int i = 0; i < flagSettingNew.getFlagClassificationGroup2Cd().size(); i++) {
            boolean check = false;
            for (int j = 0; j < arrayList.size(); j++) {
                if (!flagSettingNew.getFlagClassificationGroup2Cd().get(i).equals(arrayList.get(j))) {
                    check = true;
                } else {
                    check = false;
                    break;
                }
            }
            if (check) {
                arrListEnd.add(flagSettingNew.getFlagClassificationGroup2Cd().get(i));
                arrListNameEnd.add(flagSettingNew.getFlagClassificationGroup2Name().get(i));
            }
        }

        //Check flag group1 cd
        //If group2 not select
        if (arrListEnd.size() == 0) {
            for (int i = 0; i < flagSettingNew.getFlagClassificationGroup1Cd().size(); i++) {
                //Check select in Array<String> ?
                if (Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup1Cd().get(i))) {
                    flagCheckSelectAll = i;
                }
                //put array list group1 new (When group2cd not select)
                if (!flagGroup1CdNew.equals(flagSettingNew.getFlagClassificationGroup1Cd().get(i))) {
                    arrayListGroup1Cd.add(flagSettingNew.getFlagClassificationGroup1Cd().get(i));
                    arrayListGroup1Name.add(flagSettingNew.getFlagClassificationGroup1Name().get(i));

                }
            }
            flagSettingNew.setFlagClassificationGroup1Cd(arrayListGroup1Cd);
            flagSettingNew.setFlagClassificationGroup1Name(arrayListGroup1Name);
        } else {
            flagSettingNew.setFlagClassificationGroup2Cd(arrListEnd);
            flagSettingNew.setFlagClassificationGroup2Name(arrListNameEnd);
        }

        //Check if select all -> uncheck
        if (flagCheckSelectAll != -1) {
            flagSettingNew.getFlagClassificationGroup1Cd().remove(flagCheckSelectAll);
        }

        return flagSettingNew;
    }
}
