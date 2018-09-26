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
import android.widget.AdapterView;
import android.widget.ImageButton;
import android.widget.ListView;
import android.widget.TextView;

import com.android.productchange.R;
import com.android.productchange.adapters.ListViewGenreAdapter;
import com.android.productchange.api.Config;
import com.android.productchange.common.constants.Constants;
import com.android.productchange.db.entity.CLP;
import com.android.productchange.db.models.BookModel;
import com.android.productchange.db.models.ReturnbookModel;

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
public class DListFragment extends DialogFragment implements View.OnClickListener {

    /**
     * Interface to item selected to activity
     */
    public interface ItemSelectedDialogListener {

        /**
         * Function send list selected data to activity
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
     * Return book model
     */
    private ReturnbookModel returnbookModel;

    /**
     * List view
     */
    private ListView lsvList;

    /**
     * ID
     */
    private String idOld, idNew, id;

    /**
     * Type
     */
    private int typeOld, typeNew, type;

    /**
     * year
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
     * Button back
     */
    private ImageButton imbBack;

    /**
     * Title
     */
    private TextView txvPath;

    /**
     * Date filter
     */
    private String dateFrom, dateTo;

    /**
     * Date checked
     */
    private boolean oldDateChecked, newDateChecked, dateChecked;

    /**
     * Save flag filter genre selected
     */
    private String flagGroup1Cd;
    private String flagGroup2Cd;
    private String flagGroup2Name;

    /**
     * Save flag filter percent selected
     */
    private int flagPercent;

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

        if (getArguments() != null) {
            idOld = getArguments().getString(Constants.COLUMN_ID);
            idNew = getArguments().getString(Constants.COLUMN_ID_NEW);
            typeOld = getArguments().getInt(Config.TYPE);
            typeNew = getArguments().getInt(Config.TYPE_NEW);
            yearOld = getArguments().getInt(Constants.YEAR_AGO);
            rank = getArguments().getInt(Constants.RANK);
            oldDateChecked = getArguments().getBoolean(Constants.FLAG_DATE_CHECK);
            newDateChecked = getArguments().getBoolean(Constants.FLAG_DATE_CHECK_NEW);
            nameClassify = getArguments().getString(Constants.COLUMN_LARGE_CLASSIFICATION_NAME);
            dateFrom = getArguments().getString(Constants.COLUMN_DATE_FROM);
            dateTo = getArguments().getString(Constants.COLUMN_DATE_TO);

            //Return book
            flagGroup1Cd = getArguments().getString(Constants.FLAG_SELECT_GROUP1_CD);
            flagGroup2Cd = getArguments().getString(Constants.FLAG_SELECT_GROUP2_CD);
            flagGroup2Name = getArguments().getString(Constants.FLAG_SELECT_GROUP2_NAME);
            flagPercent = getArguments().getInt(Constants.FLAG_PERCENT_SELECTED);
        }

        txvPath = (TextView) rootView.findViewById(R.id.txv_path);

        lsvList = (ListView) rootView.findViewById(R.id.lsv_list);
        list = new ArrayList<>();
        bookModel = new BookModel();


        if (typeNew != 0) {
            if (typeOld != typeNew) {
                type = typeNew;
                id = idNew;
                dateChecked = newDateChecked;
            } else {
                type = typeOld;
                id = idOld;
                dateChecked = oldDateChecked;
            }
        } else {
            type = typeOld;
            id = idOld;
            dateChecked = oldDateChecked;
        }

        if (dateChecked) {
            id = idNew;
        }

        loadNameSelected(type);
        //loadItems(id);
        loadItems(id, flagGroup1Cd);

        imbBack = (ImageButton) rootView.findViewById(R.id.imb_back);

        lsvList.setOnItemClickListener(new AdapterView.OnItemClickListener() {
            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {

                dateChecked = false;
                // check rank is return
                // move to select year dialog
                // send data to activity
                ItemSelectedDialogListener activity =
                        (ItemSelectedDialogListener) getActivity();
                activity.onLitSelectedDialog(type, list.get(position).get(Constants.COLUMN_ID),
                        list.get(position).get(Constants.COLUMN_NAME), dateChecked);
                dismiss();

            }
        });

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
     * Set data for list item of Classify or Publisher
     *
     * @param id is id selected
     */
    public void loadItems(String id, String positionSelectGroup1Cd) {

        List<CLP> rlist = bookModel.getInfoClassifyPublisher(type, rank);

        HashMap<String, String> hashMapAll = new HashMap<>();
        if (rank == Constants.RANK_RETURN && type == Config.TYPE_CLASSIFY) {
            hashMapAll.put(Constants.COLUMN_MEDIA_GROUP1_CD, Constants.ID_ROW_ALL);
            hashMapAll.put(Constants.COLUMN_MEDIA_GROUP1_NAME, Constants.ROW_ALL);
        } else {
            hashMapAll.put(Constants.COLUMN_NAME, Constants.ROW_ALL);
            hashMapAll.put(Constants.COLUMN_ID, Constants.ID_ROW_ALL);
        }
        list.add(hashMapAll);

        // Set data into list adapter
        for (int i = 0; i < rlist.size(); i++) {
            HashMap<String, String> hashMap = new HashMap<>();
            if (rank == Constants.RANK_RETURN && type == Config.TYPE_CLASSIFY) {
                hashMap.put(Constants.COLUMN_MEDIA_GROUP1_CD, String.valueOf(rlist.get(i).getId()));
                hashMap.put(Constants.COLUMN_MEDIA_GROUP1_NAME, rlist.get(i).getName());
            } else {
                hashMap.put(Constants.COLUMN_NAME, rlist.get(i).getName());
                hashMap.put(Constants.COLUMN_ID, String.valueOf(rlist.get(i).getId()));
            }
            list.add(hashMap);
        }

        // Adapter init
        // Set data adapter to list view
        ListViewGenreAdapter adapter = new ListViewGenreAdapter(getActivity(), list, id, rank, type, positionSelectGroup1Cd);
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
            DSelectFragment dSelectFragment = new DSelectFragment();
            FragmentManager fm = getActivity().getSupportFragmentManager();
            Bundle bundle = new Bundle();
            bundle.putInt(Config.TYPE, typeOld);
            bundle.putInt(Config.TYPE_NEW, typeNew);
            bundle.putInt(Constants.RANK, rank);
            bundle.putBoolean(Constants.FLAG_DATE_CHECK, oldDateChecked);
            bundle.putBoolean(Constants.FLAG_DATE_CHECK_NEW, newDateChecked);
            bundle.putInt(Constants.YEAR_AGO, yearOld);
            bundle.putString(Constants.COLUMN_ID, idOld);
            bundle.putString(Constants.COLUMN_ID_NEW, idNew);
            bundle.putString(Constants.COLUMN_LARGE_CLASSIFICATION_NAME, nameClassify);
            bundle.putBoolean(Constants.FLAG_BACK, true);
            bundle.putString(Constants.COLUMN_DATE_FROM, dateFrom);
            bundle.putString(Constants.COLUMN_DATE_TO, dateTo);

            //Flag return book
            bundle.putInt(Constants.FLAG_PERCENT_SELECTED, flagPercent);
            bundle.putString(Constants.FLAG_SELECT_GROUP1_CD, flagGroup1Cd);
            bundle.putString(Constants.FLAG_SELECT_GROUP2_CD, flagGroup2Cd);
            bundle.putString(Constants.FLAG_SELECT_GROUP2_NAME, flagGroup2Name);
            dSelectFragment.setArguments(bundle);
            dSelectFragment.show(fm, null);
            dismiss();
        }
    }

    /**
     * Load Classify or Publisher name selected
     *
     * @param type is Classify or Publisher selected
     */
    private void loadNameSelected(int type) {
        if (type == Config.TYPE_CLASSIFY) {
            txvPath.setText(getResources().getString(R.string.select_cagtagory));
        } else {
            txvPath.setText(getResources().getString(R.string.select_publisher));
        }
    }
}
