package com.android.productchange.adapters;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.database.DataSetObserver;
import android.graphics.Color;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.TextView;

import com.android.productchange.R;
import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.utils.FormatCommon;

import java.util.ArrayList;
import java.util.HashMap;

/**
 * <h1>List View Adapter for List Genre Screen</h1>
 * <p>
 * Adapter for list view Product Change List in list view
 *
 * @author tien-lv
 * @since 2017-12-14.
 */

public class ListViewProductAdapter extends BaseAdapter {

    /**
     * List item
     */
    public ArrayList<HashMap<String, String>> list;

    /**
     * Activity
     */
    private Activity activity;
    private FormatCommon formatCommon = new FormatCommon();
    /**
     * Rank
     */
    private int rank;

    /**
     * Constructor List View Adapter
     *
     * @param activity Main Activity call this Adapter
     * @param rank     show product has rank selected
     * @param list     list data to set adapter
     */
    public ListViewProductAdapter(Activity activity, ArrayList<HashMap<String, String>> list,
                                  int rank) {

        super();
        this.activity = activity;
        this.list = list;
        this.rank = rank;
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
        View v_jan_cd, v_publisher_name_return, v_group1_name, v_group2_name, v_classify,
                v_publisher_name;
        TextView txv_id, txv_location_id, txv_name, txv_classify, txv_publish_date, txv_rank,
                txv_inventory_number, txv_jan_cd, txv_publisher_name_return, txv_group1_name,
                txv_group2_name, txv_publisher_name;
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
            convertView = inflater.inflate(R.layout.list_product_list, null);
            holder = new ViewHolder();

            // Init column in list view
            holder.txv_id = (TextView) convertView.findViewById(R.id.txv_id);
            holder.txv_location_id = (TextView) convertView.findViewById(R.id.txv_location_id);
            holder.txv_name = (TextView) convertView.findViewById(R.id.txv_name_header);
            holder.txv_classify = (TextView) convertView.findViewById(R.id.txv_classify_header);
            holder.txv_publisher_name = (TextView) convertView.findViewById(
                    R.id.txv_publisher_name);
            holder.txv_publish_date = (TextView) convertView.findViewById(
                    R.id.txv_publish_date_header);
            holder.txv_inventory_number = (TextView) convertView.findViewById(
                    R.id.txv_inventory_number_header);
            holder.txv_rank = (TextView) convertView.findViewById(R.id.txv_rank);

            //get data view and text view
            getDataView(convertView, holder);

            //hide item return books
            hidenItemReturnBooks(holder);

            convertView.setTag(holder);
        } else {
            holder = (ViewHolder) convertView.getTag();
        }

        // Set data in list to list view
        HashMap<String, String> map = list.get(position);

        // Set data when list return book
        if (rank != Constants.RANK_RETURN) {
            holder.txv_id.setText(map.get(Constants.COLUMN_ID));
            holder.txv_name.setText(map.get(Constants.COLUMN_NAME));
            holder.txv_classify.setText(map.get(Constants.COLUMN_LARGE_CLASSIFICATION_NAME));
            holder.txv_publisher_name.setText(map.get(Constants.COLUMN_PUBLISHER_NAME));
            holder.txv_inventory_number.setText(map.get(Constants.COLUMN_INVENTORY_NUMBER));
            holder.txv_publish_date.setText(formatCommon.formatDateShowList(map.get(Constants.COLUMN_PUBLISH_DATE)));
            if (map.get(Constants.COLUMN_RANKING).equals(Constants.INT_9999999)) {
                holder.txv_rank.setText(Constants.STRING_EMPTY);
            } else {
                holder.txv_rank.setText(map.get(Constants.COLUMN_RANKING));
            }
        }
        holder.txv_location_id.setText(map.get(Constants.COLUMN_LOCATION_ID));


        // set color at each row with new category rank

        if (!map.isEmpty()) {
            if (map.get(Constants.COLUMN_NEW_CATEGORY_RANK).equals(
                    String.valueOf(Constants.RANK_PLATFORM1))) {
                holder.txv_location_id.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_PLATFORM1));
                holder.txv_name.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_PLATFORM1));
                holder.txv_classify.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_PLATFORM1));
                holder.txv_publisher_name.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_PLATFORM1));
                holder.txv_publish_date.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_PLATFORM1));
                holder.txv_inventory_number.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_PLATFORM1));
                holder.txv_rank.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_PLATFORM1));
            } else if (map.get(Constants.COLUMN_NEW_CATEGORY_RANK).equals(
                    String.valueOf(Constants.RANK_PLATFORM2))) {
                holder.txv_location_id.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_PLATFORM2));
                holder.txv_name.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_PLATFORM2));
                holder.txv_classify.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_PLATFORM2));
                holder.txv_publisher_name.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_PLATFORM2));
                holder.txv_publish_date.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_PLATFORM2));
                holder.txv_inventory_number.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_PLATFORM2));
                holder.txv_rank.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_PLATFORM2));
            } else if (map.get(Constants.COLUMN_NEW_CATEGORY_RANK).equals(
                    String.valueOf(Constants.RANK_SURFACE))) {
                holder.txv_location_id.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_SURFACE));
                holder.txv_name.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_SURFACE));
                holder.txv_classify.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_SURFACE));
                holder.txv_publisher_name.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_SURFACE));
                holder.txv_publish_date.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_SURFACE));
                holder.txv_inventory_number.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_SURFACE));
                holder.txv_rank.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_SURFACE));
            } else if (map.get(Constants.COLUMN_NEW_CATEGORY_RANK).equals(
                    String.valueOf(Constants.RANK_SHELDER))) {
                holder.txv_location_id.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_SHELDER));
                holder.txv_name.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_SHELDER));
                holder.txv_classify.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_SHELDER));
                holder.txv_publisher_name.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_SHELDER));
                holder.txv_publish_date.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_SHELDER));
                holder.txv_inventory_number.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_SHELDER));
                holder.txv_rank.setBackgroundColor(
                        Color.parseColor(Constants.COLOR_RANK_SHELDER));
            } else {
                holder.txv_location_id.setBackgroundColor(Color.WHITE);
                holder.txv_name.setBackgroundColor(Color.WHITE);
                holder.txv_classify.setBackgroundColor(Color.WHITE);
                holder.txv_publisher_name.setBackgroundColor(Color.WHITE);
                holder.txv_publish_date.setBackgroundColor(Color.WHITE);
                holder.txv_inventory_number.setBackgroundColor(Color.WHITE);
                holder.txv_rank.setBackgroundColor(Color.WHITE);
            }
        }

        return convertView;
    }

    private void getDataView(View rootView, ViewHolder holder) {

        holder.v_jan_cd = rootView.findViewById(R.id.v_jan_cd);
        holder.v_publisher_name_return = rootView.findViewById(R.id.v_publisher_name_return);
        holder.v_group1_name = rootView.findViewById(R.id.v_group1_name);
        holder.v_group2_name = rootView.findViewById(R.id.v_group2_name);
        holder.v_classify = rootView.findViewById(R.id.v_classify_header);
        holder.v_publisher_name = rootView.findViewById(R.id.v_publisher_name);
        holder.txv_jan_cd = (TextView) rootView.findViewById(R.id.txv_jan_cd);
        holder.txv_publisher_name_return = (TextView) rootView.findViewById(R.id.txv_publisher_name_return);
        holder.txv_group1_name = (TextView) rootView.findViewById(R.id.txv_group1_name);
        holder.txv_group2_name = (TextView) rootView.findViewById(R.id.txv_group2_name);

    }

    private void hidenItemReturnBooks(ViewHolder holder) {

        //hide
        holder.v_jan_cd.setVisibility(View.GONE);
        holder.v_publisher_name_return.setVisibility(View.GONE);
        holder.v_group1_name.setVisibility(View.GONE);
        holder.v_group2_name.setVisibility(View.GONE);
        holder.txv_jan_cd.setVisibility(View.GONE);
        holder.txv_publisher_name_return.setVisibility(View.GONE);
        holder.txv_group1_name.setVisibility(View.GONE);
        holder.txv_group2_name.setVisibility(View.GONE);

        //show
        holder.v_classify.setVisibility(View.VISIBLE);
        holder.v_publisher_name.setVisibility(View.VISIBLE);
        holder.txv_classify.setVisibility(View.VISIBLE);
        holder.txv_publisher_name.setVisibility(View.VISIBLE);
    }

    /**
     * Register data set
     *
     * @param observer {@link DataSetObserver}
     */
    @Override
    public void registerDataSetObserver(DataSetObserver observer) {
        super.registerDataSetObserver(observer);
    }

    /**
     * Unregister data set
     *
     * @param observer {@link DataSetObserver}
     */
    @Override
    public void unregisterDataSetObserver(DataSetObserver observer) {
        super.unregisterDataSetObserver(observer);
    }

    /**
     * Notify data set changed
     */
    @Override
    public void notifyDataSetChanged() {
        super.notifyDataSetChanged();
    }
}
