package com.android.returncandidate.adapters;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.ImageView;
import android.widget.TextView;

import com.android.returncandidate.R;
import com.android.returncandidate.common.constants.Constants;

import java.util.ArrayList;
import java.util.HashMap;

/**
 * <h1>List View Adapter for List Genre Screen</h1>
 * <p>
 *
 * @author cong-pv
 * @since 2018-07-09.
 */
public class ListViewMonthAdapter extends BaseAdapter {

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
    private String positionSelect;

    /**
     * Constructor List View Adapter
     *
     * @param activity Main Activity call this Adapter
     * @param list     list data to set adapter
     */
    public ListViewMonthAdapter(Activity activity, ArrayList<HashMap<String, String>> list, String positionSelect) {

        super();
        this.positionSelect = positionSelect;
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
            convertView = inflater.inflate(R.layout.list_genre, null);
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
        holder.txv_id.setText(map.get(Constants.SELECT_POISITION));
        holder.txv_name.setText(map.get(Constants.SELECT_VALUE));
        // check id at position in list
        if (map.get(Constants.SELECT_POISITION).equals(positionSelect)) {
            holder.imv_checked.setVisibility(View.VISIBLE);
        } else {
            holder.imv_checked.setVisibility(View.INVISIBLE);
        }
        return convertView;
    }

}
