package com.android.productchange.adapters;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.ImageView;
import android.widget.TextView;

import com.android.productchange.R;
import com.android.productchange.api.Config;
import com.android.productchange.common.constants.Constants;

import java.util.ArrayList;
import java.util.HashMap;

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
    private String id;

    /**
     * Id
     */
    private int ranking;

    /**
     * type
     */
    private int type;

    /**
     * position selected group 1 cd
     */
    private String positionSelectGroupCd1;

    /**
     * Constructor List View Adapter
     *
     * @param activity Main Activity call this Adapter
     * @param id       show icon checked if id has in list data
     * @param list     list data to set adapter
     */
    public ListViewGenreAdapter(Activity activity, ArrayList<HashMap<String, String>> list,
                                String id, int ranking, int type, String positionSelectGroupCd1) {

        super();
        this.activity = activity;
        this.list = list;
        this.id = id;
        this.ranking = ranking;
        this.type = type;
        this.positionSelectGroupCd1 = positionSelectGroupCd1;
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
    }

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
            if (ranking == Constants.RANK_RETURN && type == Config.TYPE_CLASSIFY) {
                convertView = inflater.inflate(R.layout.list_genre_return_book, null);
            } else {
                convertView = inflater.inflate(R.layout.list_genre, null);
            }
            holder = new ViewHolder();

            // Init column in list view
            holder.txv_id = (TextView) convertView.findViewById(R.id.txv_id);
            holder.txv_name = (TextView) convertView.findViewById(R.id.txv_name_header);
            holder.imv_checked = (ImageView) convertView.findViewById(R.id.imv_checked);

            convertView.setTag(holder);
        } else {
            holder = (ViewHolder) convertView.getTag();
        }

        // Set data in list to list view
        HashMap<String, String> map = list.get(position);
        if (ranking == Constants.RANK_RETURN && type == Config.TYPE_CLASSIFY) {
            holder.txv_id.setText(map.get(Constants.COLUMN_MEDIA_GROUP1_CD));
            holder.txv_name.setText(map.get(Constants.COLUMN_MEDIA_GROUP1_NAME));
            // check id at position in list
            if (map.get(Constants.COLUMN_MEDIA_GROUP1_CD).equals(positionSelectGroupCd1)) {
                holder.imv_checked.setVisibility(View.VISIBLE);
            } else {
                holder.imv_checked.setVisibility(View.INVISIBLE);
            }
        } else {
            holder.txv_id.setText(map.get(Constants.COLUMN_ID));
            holder.txv_name.setText(map.get(Constants.COLUMN_NAME));
            // check id at position in list
            if (map.get(Constants.COLUMN_ID).equals(id)) {
                holder.imv_checked.setVisibility(View.VISIBLE);
            } else {
                holder.imv_checked.setVisibility(View.INVISIBLE);
            }
        }

        return convertView;
    }

}
