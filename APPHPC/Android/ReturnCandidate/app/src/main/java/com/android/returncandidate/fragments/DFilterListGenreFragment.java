package com.android.returncandidate.fragments;

import android.app.AlertDialog;
import android.app.Dialog;
import android.content.DialogInterface;
import android.graphics.Color;
import android.graphics.drawable.ColorDrawable;
import android.os.Bundle;
import android.support.v4.app.DialogFragment;
import android.support.v4.app.FragmentManager;
import android.view.Gravity;
import android.view.KeyEvent;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.view.Window;
import android.widget.AdapterView;
import android.widget.ImageButton;
import android.widget.ListView;
import android.widget.TextView;

import com.android.returncandidate.R;
import com.android.returncandidate.adapters.ListViewGenreAdapter;
import com.android.returncandidate.common.constants.Constants;
import com.android.returncandidate.common.constants.Message;
import com.android.returncandidate.common.utils.Common;
import com.android.returncandidate.common.utils.FlagSettingNew;
import com.android.returncandidate.common.utils.FlagSettingOld;
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
 * @author tien-lv
 * @since 2017/12/19.
 */
@SuppressWarnings("deprecation")
public class DFilterListGenreFragment extends DialogFragment implements View.OnClickListener {

    /**
     * Interface to item selected to activity
     */
    public interface ItemSelectedDialogListener {

        /**
         * Function send list selected data to activity
         * MOG
         *
         * @param typeSelected         is type selected
         * @param idSelected           is id selected
         * @param nameClassifySelected is name Classify or Publisher selected
         * @param dateChecked          {@link boolean}
         */
        void onLitSelectedDialog(int typeSelected, String idSelected, String nameClassifySelected,
                                 boolean dateChecked);
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
        txvHeaderFilter.setText(Constants.HEADER_CLASSIFICATION_1);

        if (getArguments() != null) {
            common.SetArgumentsFlagData(flagSettingNew, flagSettingOld, getArguments());
            //get flag switch OCR
            flagSwitchOCR = getArguments().getString(Constants.FLAG_SWITCH_OCR);
        }
        // list group1 cd
        lsvList = (ListView) rootView.findViewById(R.id.lsv_list);
        list = new ArrayList<>();
        bookModel = new BookModel();

        //show item list group 1 cd
        loadListItemGroup1Cd();

        imbBack = (ImageButton) rootView.findViewById(R.id.imb_back);

        lsvList.setOnItemClickListener(new AdapterView.OnItemClickListener() {
            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {

                if (position != 0) {
                    putToListItemCheckToFlagSettingNew();
                    moveGroup2Cd(position);
                }
            }
        });

        //Event long press item list view
        lsvList.setOnItemLongClickListener(new AdapterView.OnItemLongClickListener() {
            @Override
            public boolean onItemLongClick(AdapterView<?> parent, View view, final int position, long id) {

                //If long press item no check
                if ((Constants.VALUE_STR_NO_CHECK).equals(list.get(position).get(Constants.FLAG_SELECT)) ||
                        position == Constants.VALUE_INT_DEFAULT_SELECT_ALL) {
                    return false;
                } else {
                    //Update list flag setting new items
                    putToListItemCheckToFlagSettingNew();
                    //Check flag setting new
                    //If flag setting new = -1 => select all
                    convertFlagSettingNewSelectAll();

                    AlertDialog.Builder alertDialogBuilder = new AlertDialog.Builder(getActivity());
                    alertDialogBuilder.setMessage(String.format(Message.MESSAGE_CONFIRM_UNCHECK_GROUP_CD, list.get(position).get(Constants.COLUMN_MEDIA_GROUP1_NAME)));

                    alertDialogBuilder.setCancelable(false);
                    // Configure alert dialog button
                    alertDialogBuilder.setPositiveButton(Message.MESSAGE_SELECT_YES,
                            new DialogInterface.OnClickListener() {
                                @Override
                                public void onClick(DialogInterface dialog, int which) {

                                    ArrayList<String> arrayListGroup1Cd = new ArrayList<>();
                                    ArrayList<String> arrayListGroup1Name = new ArrayList<>();
                                    for (int i = 0; i < flagSettingNew.getFlagClassificationGroup1Cd().size(); i++) {
                                        if (!list.get(position).get(Constants.COLUMN_MEDIA_GROUP1_CD).equals(flagSettingNew.getFlagClassificationGroup1Cd().get(i))) {
                                            arrayListGroup1Cd.add(flagSettingNew.getFlagClassificationGroup1Cd().get(i));
                                            arrayListGroup1Name.add(flagSettingNew.getFlagClassificationGroup1Name().get(i));

                                        }
                                    }

                                    flagSettingNew.setFlagClassificationGroup1Cd(arrayListGroup1Cd);
                                    flagSettingNew.setFlagClassificationGroup1Name(arrayListGroup1Name);

                                    //Remove item group2 to list
                                    removeGroup2WhenLongPress(list.get(position).get(Constants.COLUMN_MEDIA_GROUP1_CD));

                                    //If flag setting new = -1 && len (flag setting new = (list.size() -1)
                                    deleteFlagSettingNewSelectAll();
                                    list = loadListItemGroup1CdReload();
                                    dialog.dismiss();
                                    // Hidden View holder
                                }
                            });
                    alertDialogBuilder.setNegativeButton(Message.MESSAGE_SELECT_NO,
                            new DialogInterface.OnClickListener() {
                                @Override
                                public void onClick(DialogInterface dialog, int which) {
                                    dialog.dismiss();
                                }
                            });
                    AlertDialog alert = alertDialogBuilder.show();
                    TextView messageText = (TextView) alert.findViewById(android.R.id.message);
                    messageText.setGravity(Gravity.CENTER);
                    return true;
                }
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

    private void removeGroup2WhenLongPress(String valueGroup1Cd) {

        List<CLP> rlist = bookModel.getInfoGroupCd2(valueGroup1Cd);
        HashMap<String, ArrayList<String>> hashMap = common.ListGroup2CdNew(flagSettingNew, rlist);
        flagSettingNew.setFlagClassificationGroup2Cd(hashMap.get(Constants.COLUMN_MEDIA_GROUP2_CD));
        flagSettingNew.setFlagClassificationGroup2Name(hashMap.get(Constants.COLUMN_MEDIA_GROUP2_NAME));
    }

    private void moveGroup2Cd(int position) {

        // move to select year dialog
        DFilterListGenre2Fragment dFilterListGenre2Fragment = new DFilterListGenre2Fragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();

        //save flag old if back
        flagSettingNew.setFlagClassificationGroup1CdOld(flagSettingNew.getFlagClassificationGroup1Cd());
        flagSettingNew.setFlagClassificationGroup1NameOld(flagSettingNew.getFlagClassificationGroup1Name());

        //Set flag selected group 1 cd
        boolean checkValueSelect = true;
        String valueSelected = list.get(position).get(Constants.COLUMN_MEDIA_GROUP1_CD);
        String valueSelectedName = list.get(position).get(Constants.COLUMN_MEDIA_GROUP1_NAME);
        //Loop check value
        for (String valueGroup1Cd : flagSettingNew.getFlagClassificationGroup1Cd()) {
            if (valueSelected.equals(valueGroup1Cd)) {
                checkValueSelect = false;
                break;
            }
        }
        if (checkValueSelect) {
            flagSettingNew.getFlagClassificationGroup1Cd().add(valueSelected);
            flagSettingNew.setFlagClassificationGroup1Cd(flagSettingNew.getFlagClassificationGroup1Cd());
            flagSettingNew.getFlagClassificationGroup1Name().add(valueSelectedName);
            flagSettingNew.setFlagClassificationGroup1Name(flagSettingNew.getFlagClassificationGroup1Name());
        }

        //Check if len flagSettingNew = 2 -> check select all
        if (flagSettingNew.getFlagClassificationGroup2Cd().size() == 2) {
            int positionDelete = -1;
            for (int i = 0; i < flagSettingNew.getFlagClassificationGroup2Cd().size(); i++) {
                if (Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup2Cd().get(i))) {
                    positionDelete = i;
                }
            }
            if (positionDelete != -1) {
                flagSettingNew.getFlagClassificationGroup2Cd().remove(positionDelete);
            }
        }
        //Put bundle
        Bundle bundle = common.DataPutActivity(flagSettingNew, flagSettingOld);
        //Put flag group1 --> select group2 cd
        bundle.putString(Constants.FLAG_GROUP1_CD, valueSelected);
        bundle.putString(Constants.FLAG_GROUP1_NAME, valueSelectedName);
        //put flag switch OCR
        bundle.putString(Constants.FLAG_SWITCH_OCR, flagSwitchOCR);

        dFilterListGenre2Fragment.setArguments(bundle);
        dFilterListGenre2Fragment.show(fm, null);
        dismiss();
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
     * Get data filter group1_cd with table local book
     */
    private void loadListItemGroup1Cd() {

        List<CLP> rlist = bookModel.getInfoGroupCd1();

        //Check all or not all
        int intcheck = 0;
        for (int i = 0; i < rlist.size(); i++) {
            for (int j = 0; j < flagSettingNew.getFlagClassificationGroup1Cd().size(); j++) {
                if (rlist.get(i).getId().equals(flagSettingNew.getFlagClassificationGroup1Cd().get(j))) {
                    intcheck++;
                    break;
                }
            }
        }
        HashMap<String, String> hashMapAll = new HashMap<>();
        hashMapAll.put(Constants.COLUMN_MEDIA_GROUP1_CD, Constants.ID_ROW_ALL);
        hashMapAll.put(Constants.COLUMN_MEDIA_GROUP1_NAME, Constants.ROW_ALL);

        //Uncheck select all
        for (String valueGroup1Cd : flagSettingNew.getFlagClassificationGroup1Cd()) {
            if (Constants.ID_ROW_ALL.equals(valueGroup1Cd) || intcheck == rlist.size()) {
                hashMapAll.put(Constants.FLAG_SELECT, Constants.VALUE_STR_CHECK);
            } else {
                hashMapAll.put(Constants.FLAG_SELECT, Constants.VALUE_STR_NO_CHECK);
            }
        }
        list.add(hashMapAll);

        // Set data into list adapter
        for (int i = 0; i < rlist.size(); i++) {
            HashMap<String, String> hashMap = new HashMap<>();
            hashMap.put(Constants.COLUMN_MEDIA_GROUP1_CD, String.valueOf(rlist.get(i).getId()));
            hashMap.put(Constants.COLUMN_MEDIA_GROUP1_NAME, rlist.get(i).getName());
            for (String valueGroup1Cd : flagSettingNew.getFlagClassificationGroup1Cd()) {
                if (Constants.ID_ROW_ALL.equals(valueGroup1Cd)) {
                    hashMap.put(Constants.FLAG_SELECT, Constants.VALUE_STR_CHECK);
                    break;
                } else {
                    if (rlist.get(i).getId().equals(valueGroup1Cd)) {
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
        ListViewGenreAdapter adapter = new ListViewGenreAdapter(getActivity(), list, flagSettingNew);
        lsvList.setAdapter(adapter);
    }

    /**
     * Get data filter group1_cd with table local book
     */
    private ArrayList<HashMap<String, String>> loadListItemGroup1CdReload() {

        ArrayList<HashMap<String, String>> listNew = new ArrayList<>();
        List<CLP> rlist = bookModel.getInfoGroupCd1();

        HashMap<String, String> hashMapAll = new HashMap<>();
        hashMapAll.put(Constants.COLUMN_MEDIA_GROUP1_CD, Constants.ID_ROW_ALL);
        hashMapAll.put(Constants.COLUMN_MEDIA_GROUP1_NAME, Constants.ROW_ALL);

        for (String valueGroup1Cd : flagSettingNew.getFlagClassificationGroup1Cd()) {
            if (Constants.ID_ROW_ALL.equals(valueGroup1Cd)) {
                hashMapAll.put(Constants.FLAG_SELECT, Constants.VALUE_STR_CHECK);
                //If check select all -> break foreach
                // break;
            } else {
                hashMapAll.put(Constants.FLAG_SELECT, Constants.VALUE_STR_NO_CHECK);
            }
        }
        listNew.add(hashMapAll);

        // Set data into list adapter
        for (int i = 0; i < rlist.size(); i++) {
            HashMap<String, String> hashMap = new HashMap<>();
            hashMap.put(Constants.COLUMN_MEDIA_GROUP1_CD, String.valueOf(rlist.get(i).getId()));
            hashMap.put(Constants.COLUMN_MEDIA_GROUP1_NAME, rlist.get(i).getName());
            for (String valueGroup1Cd : flagSettingNew.getFlagClassificationGroup1Cd()) {
                if (Constants.ID_ROW_ALL.equals(valueGroup1Cd)) {
                    hashMap.put(Constants.FLAG_SELECT, Constants.VALUE_STR_CHECK);
                    break;
                } else {
                    if (rlist.get(i).getId().equals(valueGroup1Cd)) {
                        hashMap.put(Constants.FLAG_SELECT, Constants.VALUE_STR_CHECK);
                        break;
                    } else {
                        hashMap.put(Constants.FLAG_SELECT, Constants.VALUE_STR_NO_CHECK);
                    }
                }
            }
            listNew.add(hashMap);
        }
        // Adapter init
        // Set data adapter to list view
        ListViewGenreAdapter adapter = new ListViewGenreAdapter(getActivity(), listNew, flagSettingNew);
        lsvList.setAdapter(adapter);
        return listNew;
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

        putToListItemCheckToFlagSettingNewBack();

        //Calculator group2 cd
        ArrayList<String> arrGroup1Cd = new ArrayList<>();
        for (String value : flagSettingNew.getFlagClassificationGroup1Cd()) {
            arrGroup1Cd.add(value);
        }
        //Check list all

        List<CLP> clpList = bookModel.getInfoGroupCd2WhenGroup1CdMulti(arrGroup1Cd);
        ArrayList<String> arrGroup2CdNew = new ArrayList<>();
        ArrayList<String> arrGroup2CdOld = new ArrayList<>();
        ArrayList<String> arrGroup2NameOld = new ArrayList<>();

        for (int i = 0; i < clpList.size(); i++) {
            arrGroup2CdNew.add(clpList.get(i).getId());
        }

        //Get list group2 cd old
        for (int i = 0; i < flagSettingNew.getFlagClassificationGroup2Cd().size(); i++) {
            arrGroup2CdOld.add(flagSettingNew.getFlagClassificationGroup2Cd().get(i));
            arrGroup2NameOld.add(flagSettingNew.getFlagClassificationGroup2Name().get(i));
        }

        //Check group2 old into group2 new
        HashMap<String, ArrayList<String>> hashMap = common.ListGroup2Cd(arrGroup2CdNew, arrGroup2CdOld, arrGroup2NameOld);
        flagSettingNew.setFlagClassificationGroup2Cd(hashMap.get(Constants.COLUMN_MEDIA_GROUP2_CD));
        flagSettingNew.setFlagClassificationGroup2Name(hashMap.get(Constants.COLUMN_MEDIA_GROUP2_NAME));


        DFilterSettingFragment dFilterSettingFragment = new DFilterSettingFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();
        Bundle bundle = common.DataPutActivity(flagSettingNew, flagSettingOld);
        //put flag switch OCR
        bundle.putString(Constants.FLAG_SWITCH_OCR, flagSwitchOCR);

        dFilterSettingFragment.setArguments(bundle);
        dFilterSettingFragment.show(fm, null);
        dismiss();
    }

    //When click all -> put to list into flag setting new
    private void putToListItemCheckToFlagSettingNew() {

        ArrayList<String> arrayList = new ArrayList<>();
        ArrayList<String> arrayListName = new ArrayList<>();
        for (int i = 0; i < list.size(); i++) {
            if (Constants.VALUE_STR_CHECK.equals(list.get(i).get(Constants.FLAG_SELECT))) {
                //flagSettingNew.setFlagPublisher();
                arrayList.add(list.get(i).get(Constants.COLUMN_MEDIA_GROUP1_CD));
                arrayListName.add(list.get(i).get(Constants.COLUMN_MEDIA_GROUP1_NAME));
            }
        }
        flagSettingNew.setFlagClassificationGroup1Cd(arrayList);
        flagSettingNew.setFlagClassificationGroup1Name(arrayListName);
    }

    //Check group1 cd when back
    private void putToListItemCheckToFlagSettingNewBack() {

        ArrayList<String> arrayList = new ArrayList<>();
        ArrayList<String> arrayListName = new ArrayList<>();
        for (int i = 0; i < list.size(); i++) {
            if (Constants.VALUE_STR_CHECK.equals(list.get(i).get(Constants.FLAG_SELECT))) {
                arrayList.add(list.get(i).get(Constants.COLUMN_MEDIA_GROUP1_CD));
                arrayListName.add(list.get(i).get(Constants.COLUMN_MEDIA_GROUP1_NAME));
            }
        }
        if (arrayList.size() == 0) {
            for (int i = 0; i < list.size(); i++) {
                arrayList.add(list.get(i).get(Constants.COLUMN_MEDIA_GROUP1_CD));
                arrayListName.add(list.get(i).get(Constants.COLUMN_MEDIA_GROUP1_NAME));
            }
        }
        flagSettingNew.setFlagClassificationGroup1Cd(arrayList);
        flagSettingNew.setFlagClassificationGroup1Name(arrayListName);
    }

    //Convert flag setting new -1= > select all
    private void convertFlagSettingNewSelectAll() {

        if (flagSettingNew.getFlagClassificationGroup1Cd().size() == 1 && Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup1Cd().get(0))) {

            ArrayList<String> arrayListGroup1Cd = new ArrayList<>();
            ArrayList<String> arrayListGroup1Name = new ArrayList<>();
            for (int i = 0; i < list.size(); i++) {
                //If position select -> uncheck group1 cd
                arrayListGroup1Cd.add(list.get(i).get(Constants.COLUMN_MEDIA_GROUP1_CD));
                arrayListGroup1Name.add(list.get(i).get(Constants.COLUMN_MEDIA_GROUP1_NAME));
            }
            flagSettingNew.setFlagClassificationGroup1Cd(arrayListGroup1Cd);
            flagSettingNew.setFlagClassificationGroup1Name(arrayListGroup1Name);
        }
    }

    //Delete select all (-1) with flag setting new
    private void deleteFlagSettingNewSelectAll() {

        if (flagSettingNew.getFlagClassificationGroup1Cd().size() == list.size() - 1) {
            flagSettingNew.getFlagClassificationGroup1Cd().remove(0);
            flagSettingNew.getFlagClassificationGroup1Name().remove(0);
        }
    }
}

