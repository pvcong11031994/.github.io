package com.android.returncandidate.adapters;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.os.Bundle;
import android.support.v4.app.FragmentManager;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.CheckBox;
import android.widget.ImageView;
import android.widget.TextView;

import com.android.returncandidate.R;
import com.android.returncandidate.common.constants.Constants;
import com.android.returncandidate.common.utils.Common;
import com.android.returncandidate.common.utils.FlagSettingNew;
import com.android.returncandidate.db.entity.CLP;
import com.android.returncandidate.db.models.BookModel;
import com.android.returncandidate.fragments.DFilterListGenre2Fragment;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

/**
 * <h1>List View Adapter for List Genre Screen</h1>
 * <p>
 * Adapter for list view Classify or Publisher in filter dialog
 *
 * @author tien-lv
 * @since 2017-12-14.
 */
public class ListViewGenreAdapter extends BaseAdapter {

    /**
     * List item
     */
    public ArrayList<HashMap<String, String>> list;

    /**
     * Activity
     */
    private Activity activity;

    /**
     * Id
     */
    private FlagSettingNew flagSettingNew;

    /**
     * Constructor List View Adapter
     *
     * @param activity Main Activity call this Adapter
     * @param list     list data to set adapter
     */
    public ListViewGenreAdapter(Activity activity, ArrayList<HashMap<String, String>> list, FlagSettingNew flagSettingNew) {

        super();
        this.flagSettingNew = flagSettingNew;
        this.activity = activity;
        this.list = list;
    }

    /**
     * Get count item
     *
     * @return list size is int
     */
    @Override
    public int getCount() {
        return list.size();
    }

    /**
     * Get item with position
     *
     * @param position position selected
     * @return object in list by position selected
     */
    @Override
    public Object getItem(int position) {
        return list.get(position);
    }

    /**
     * Get Item Id with position
     *
     * @param position position selected
     * @return id for default
     */
    @Override
    public long getItemId(int position) {
        return 0;
    }

    /**
     * Init View Holder
     */
    private class ViewHolder {
        TextView txv_id;
        TextView txv_name;
        ImageView imv_checked;
        CheckBox checkBox;
    }

    private Common common = new Common();

    /**
     * Set custom layout for list view
     *
     * @param position    in list data
     * @param convertView set layout from custom layout
     * @param parent      view group
     * @return view from layout custom
     */
    @SuppressLint("InflateParams")
    @Override
    public View getView(int position, View convertView, ViewGroup parent) {

        ViewHolder holder;
        LayoutInflater inflater = activity.getLayoutInflater();

        if (convertView == null) {

            // Init custom layout list genre
            convertView = inflater.inflate(R.layout.list_genre_multi, null);
            holder = new ViewHolder();


            // Init column in list view
            holder.txv_id = (TextView) convertView.findViewById(R.id.txv_id);
            holder.txv_name = (TextView) convertView.findViewById(R.id.txv_name_header);
            //holder.imv_checked = (ImageView) convertView.findViewById(R.id.imv_checked);
            holder.checkBox = (CheckBox) convertView.findViewById(R.id.rowCheckBox);
            convertView.setTag(holder);
        } else {
            holder = (ViewHolder) convertView.getTag();
        }

        // Set data in list to list view
        HashMap<String, String> map = list.get(position);
        holder.txv_id.setText(map.get(Constants.COLUMN_MEDIA_GROUP1_CD));
        holder.txv_name.setText(map.get(Constants.COLUMN_MEDIA_GROUP1_NAME));
        if (Constants.VALUE_STR_CHECK.equals(list.get(position).get(Constants.FLAG_SELECT))) {
            holder.checkBox.setChecked(true);
        } else {
            holder.checkBox.setChecked(false);
        }
        holder.checkBox.setOnClickListener(onStateChangedListener(holder.checkBox, position));
        return convertView;
    }

    private View.OnClickListener onStateChangedListener(final CheckBox checkBoxs, final int positions) {
        return new View.OnClickListener() {
            @Override
            public void onClick(View v) {

                //If check all
                if (positions == 0) {
                    //when select all is check
                    if (!checkBoxs.isChecked()) {
                        for (int i = 0; i < list.size(); i++) {
                            list.get(i).put(Constants.FLAG_SELECT, Constants.VALUE_STR_NO_CHECK);
                            list.set(i, list.get(i));
                        }
                    } else { // when select all is not check
                        for (int i = 0; i < list.size(); i++) {
                            list.get(i).put(Constants.FLAG_SELECT, Constants.VALUE_STR_CHECK);
                            list.set(i, list.get(i));
                        }
                    }
                    //Get list default group1 cd
                    BookModel mBookModel = new BookModel();
                    ArrayList<String> arrGroup1Cd = new ArrayList<>();
                    ArrayList<String> arrGroup1Name = new ArrayList<>();
                    ArrayList<String> arrGroup2Cd = new ArrayList<>();
                    ArrayList<String> arrGroup2Name = new ArrayList<>();
                    arrGroup1Cd.add(Constants.ID_ROW_ALL);
                    arrGroup1Name.add(Constants.ROW_ALL);
                    arrGroup2Cd.add(Constants.ID_ROW_ALL);
                    arrGroup2Name.add(Constants.ROW_ALL);
                    List<CLP> listDefaultGroup1 = mBookModel.getInfoGroupCd1();
                    List<CLP> listDefaultGroup2 = mBookModel.getInfoGroupCd2(Constants.ID_ROW_ALL);
                    for (int i = 0; i < listDefaultGroup1.size(); i++) {
                        arrGroup1Cd.add(listDefaultGroup1.get(i).getId());
                        arrGroup1Name.add(listDefaultGroup1.get(i).getName());
                    }
                    for (int i = 0; i < listDefaultGroup2.size(); i++) {
                        arrGroup2Cd.add(listDefaultGroup2.get(i).getId());
                        arrGroup2Name.add(listDefaultGroup2.get(i).getName());
                    }
                    flagSettingNew.setFlagClassificationGroup1Cd(arrGroup1Cd);
                    flagSettingNew.setFlagClassificationGroup1Name(arrGroup1Name);
                    flagSettingNew.setFlagClassificationGroup2Cd(arrGroup2Cd);
                    flagSettingNew.setFlagClassificationGroup2Name(arrGroup2Name);
                }
                notifyDataSetChanged();

            }
        };
    }

}
