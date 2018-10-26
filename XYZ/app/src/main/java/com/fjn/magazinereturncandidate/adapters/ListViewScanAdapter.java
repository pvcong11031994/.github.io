package com.fjn.magazinereturncandidate.adapters;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.graphics.Color;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.TextView;

import com.fjn.magazinereturncandidate.R;
import com.fjn.magazinereturncandidate.common.constants.Constants;

import java.util.LinkedList;

/**
 * List View Adapter for List Genre Screen
 *
 * @author cong-pv
 * @version 1.0
 * @since 2018-10-17
 */

public class ListViewScanAdapter extends BaseAdapter {

    private LinkedList<String[]> list;

    private Activity activity;

    /**
     * Hook {@link ListViewScanAdapter} to designated {@link Activity}
     *
     * @param activity {@link Activity}
     * @param list     {@link LinkedList}
     */
    public ListViewScanAdapter(Activity activity, LinkedList<String[]> list) {

        super();
        this.activity = activity;
        this.list = list;
    }

    /**
     * Get count item
     */
    @Override
    public int getCount() {
        return list.size();
    }

    /**
     * Get item at index
     */
    @Override
    public Object getItem(int position) {
        return list.get(position);
    }

    /**
     * Get Item Id with position
     */
    @Override
    public long getItemId(int position) {
        return 0;
    }

    /**
     * Init View Holder
     */
    private class ViewHolder {
        TextView txv_column1;
        TextView txv_column2;
        TextView txv_column3;
        TextView txv_color;
    }

    /**
     * Set custom layout for list view
     *
     * @param position    int
     * @param convertView {@link View}
     * @param parent      {@link ViewGroup}
     */
    @SuppressLint("InflateParams")
    @Override
    public View getView(int position, View convertView, ViewGroup parent) {

        ViewHolder holder;
        LayoutInflater inflater = activity.getLayoutInflater();

        if (convertView == null) {

            // Init custom layout list genre
            convertView = inflater.inflate(R.layout.adapter_list_scan, null);
            holder = new ViewHolder();


            // Init column in list view
            holder.txv_column1 = (TextView) convertView.findViewById(R.id.txv_column1);
            holder.txv_column2 = (TextView) convertView.findViewById(R.id.txv_column2);
            holder.txv_column3 = (TextView) convertView.findViewById(R.id.txv_column3);
            holder.txv_color = (TextView) convertView.findViewById(R.id.txv_color);

            convertView.setTag(holder);
        } else {
            holder = (ViewHolder) convertView.getTag();
        }

        // Set data in list to list view
        String[] item = list.get(position);
        if (item == null || item.length == 0) return convertView;
        holder.txv_column1.setText(item[0]);
        if (item.length > 1) {
            holder.txv_column2.setText(item[2]);
            holder.txv_column3.setText(item[1]);
            holder.txv_color.setText(item[17]);
        }
        if (item[17] == null) {
            convertView.setBackgroundColor(Color.RED);
            return convertView;
        }
        if (Constants.FLAG_0.contains(item[17])) {
            convertView.setBackgroundColor(Color.RED);
        } else {
            convertView.setBackgroundColor(Color.GRAY);
        }

        return convertView;
    }
}
