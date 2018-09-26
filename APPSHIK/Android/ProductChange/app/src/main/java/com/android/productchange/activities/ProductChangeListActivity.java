package com.android.productchange.activities;

import android.annotation.SuppressLint;
import android.app.ProgressDialog;
import android.content.DialogInterface;
import android.content.Intent;
import android.graphics.Color;
import android.os.Bundle;
import android.os.Parcelable;
import android.support.v4.app.FragmentManager;
import android.support.v7.app.AppCompatActivity;
import android.view.Gravity;
import android.view.View;
import android.widget.AdapterView;
import android.widget.FrameLayout;
import android.widget.ImageButton;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.ListView;
import android.widget.TextView;

import com.android.productchange.R;
import com.android.productchange.adapters.ListViewProductAdapter;
import com.android.productchange.adapters.ListViewProductReturnBooksAdapter;
import com.android.productchange.api.Config;
import com.android.productchange.common.constants.Constants;
import com.android.productchange.common.constants.Message;
import com.android.productchange.common.helpers.DatabaseHelper;
import com.android.productchange.common.utils.Common;
import com.android.productchange.common.utils.DatabaseManager;
import com.android.productchange.common.utils.EndlessScrollListener;
import com.android.productchange.common.utils.FlagSettingNew;
import com.android.productchange.common.utils.FlagSettingOld;
import com.android.productchange.common.utils.FormatCommon;
import com.android.productchange.common.utils.LogManager;
import com.android.productchange.common.utils.ProcessDialogLoadListReturnBooks;
import com.android.productchange.db.entity.Books;
import com.android.productchange.db.entity.CLP;
import com.android.productchange.db.entity.MaxYearRank;
import com.android.productchange.db.entity.Returnbooks;
import com.android.productchange.db.models.BookModel;
import com.android.productchange.db.models.MaxYearRankModel;
import com.android.productchange.db.models.PeriodbookModel;
import com.android.productchange.db.models.RegularbookModel;
import com.android.productchange.db.models.ReturnbookModel;
import com.android.productchange.fragments.DFilterSettingFragment;
import com.android.productchange.fragments.DListFragment;
import com.android.productchange.fragments.DProductDetailFragment;
import com.android.productchange.fragments.DRankFragment;
import com.android.productchange.fragments.DSelectDateFragment;
import com.android.productchange.fragments.DSelectFragment;
import com.android.productchange.fragments.DSelectYearFragment;
import com.android.productchange.interfaces.AsyncResponseListReturnBooks;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.Date;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;


/**
 * <h1>Product Change List Activity</h1>
 * <p>
 * Load data from table into list view
 * Filter data on item selected
 *
 * @author tien-lv
 * @since 2018-02-08
 */
@SuppressWarnings("deprecation")
public class ProductChangeListActivity extends AppCompatActivity implements View.OnClickListener,
        DRankFragment.RankDialogListener, DListFragment.ItemSelectedDialogListener,
        DSelectYearFragment.SelectedFilterDialogListener,
        DSelectDateFragment.SelectDateDialogListener,
        DFilterSettingFragment.ItemSelectedFilterSettingDialogListener,
        AsyncResponseListReturnBooks {

    /**
     * TAG
     */
    private String TAG = Constants.TAG_APPLICATION_NAME;

    /**
     * List data for List View
     */
    private ArrayList<HashMap<String, String>> list, datalist;

    /**
     * User id
     */
    String userID;

    /**
     * Shop id
     */
    private String shopID;

    /**
     * Server name
     */
    String serverName;

    /**
     * Type
     */
    private int type;

    /**
     * Offset
     */
    private int offset;

    /**
     * Rank
     */
    private int rank;

    /**
     * year ago
     */
    private int yearAgo;

    /**
     * Linear layout title and selection
     */
    private LinearLayout llTitle, llSelection;

    /**
     * Frame layout close, filter, selection
     */
    private FrameLayout flClose, flFilter, flSelection;

    /**
     * Text view header of list
     */
    private TextView txvLocarionHeader, txvNameHeader, txvClassifyHeader, txvPublisherHeader,
            txvPublishDateHeader, txvInventoryNumberHeader, txvRankHeader, txvJanCdHeader,
            txvPublisherReturnHeader, txvGroup1NameHeader, txvGroup2NameHeader;
    /**
     * Image view header of list
     */
    private ImageView imvLocationHeader, imvNameHeader, imvClassifyHeader, imvPublisherHeader,
            imvPublishDateHeader, imvInventoryNumberHeader, imvRankHeader, imvJanCdHeader,
            imvPublisherReturnHeader, imvGroup1NameHeader, imvGroup2NameHeader;

    private FrameLayout fl_jan_cd_header, fl_publisher_name_return_header, fl_group1_name_header,
            fl_group2_name_header, fl_classify_header, fl_publisher_header;

    private View v_jan_cd_header, v_publisher_name_return_header, v_group1_name_header,
            v_group2_name_header, v_classify_header, v_publisher_header;
    /**
     * Image view Icon
     */
    private ImageView imvIcon;

    /**
     * Text view Rank name and classify name
     */
    private TextView txvRankName, txvClassifyName;

    /**
     * Name Classify
     */
    private String nameClassify;

    /**
     * ID
     */
    private String id;

    /**
     * Book model
     */
    private BookModel bookModel;

    /**
     * Return book model
     */
    private ReturnbookModel returnbookModel;

    /**
     * Period book model
     */
    private PeriodbookModel periodbookModel;

    /**
     * Regular book model
     */
    private RegularbookModel regularbookModel;

    /**
     * List view
     */
    private ListView lstView;

    /**
     * Adapter
     */
    private ListViewProductAdapter adapter;
    private ListViewProductReturnBooksAdapter adapterReturnBooks;

    /**
     * Progress dialog.
     */
    private ProgressDialog progress;

    /**
     * Date filter
     */
    private String dateFrom, dateTo, createDate;

    Parcelable stateListArrival = null, stateListPlatform1 = null, stateListPlatform2 = null,
            stateListSurface = null, stateListShelder = null, stateListReturn = null,
            stateListPeriod = null, stateListRegular = null;

    /**
     * Date checked
     */
    private boolean checkedDate;

    private int countLocationClick, countNameClick, countClassifyClick, countPublisherClick,
            countPublishDateClick, countInventoryClick, countRankClick, countJanCdClick, countGroup1NameClick,
            countGroup2NameClick;
    private LinkedHashMap<String, String> mapOrder;
    private Map<Integer, String> treeMapOrder;

    long timeout;

    private int percent;
    private String flagGroup1Cd;
    private String flagGroup2Cd;
    private String flagGroup2Name;
    private Boolean flagFilterSubmit;
    /**
     * Save flag filter
     */
    private FlagSettingNew flagSettingNew = new FlagSettingNew();
    private FlagSettingOld flagSettingOld = new FlagSettingOld();
    private Common common = new Common();

    private MaxYearRank maxYearRank = new MaxYearRank();
    private MaxYearRankModel maxYearRankModel = new MaxYearRankModel();

    private FormatCommon formatCommon;

    private ImageButton imbRank, imbSelect, imbClose;

    /**
     * Init on Create Activity
     *
     * @param savedInstanceState is Bundle of activity
     */
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_product_change_list);

        // init process loading screen
        progress = new ProgressDialog(this);
        progress.setMessage(Message.MESSAGE_LOADING_SCREEN);
        progress.setCancelable(false);


        // get data form bundle
        Bundle bundle = getIntent().getExtras();
        if (bundle != null) {
            userID = bundle.getString(Constants.COLUMN_USER_ID);
            shopID = bundle.getString(Constants.COLUMN_SHOP_ID);
            serverName = bundle.getString(Constants.COLUMN_SERVER_NAME);
            rank = bundle.getInt(Constants.RANK);
            createDate = bundle.getString(Constants.COLUMN_CREATE_DATE);
        }
        // set flag return book default
        //-----------------START-----------------
        percent = Constants.SELECT_PERCENT_30;
        flagGroup1Cd = Constants.ID_ROW_ALL;
        flagGroup2Cd = Constants.ID_ROW_ALL;
        flagGroup2Name = Constants.ROW_ALL;
        //------------------END------------------
        // set default rank
        if (rank == Constants.SELECT_NULL) {
            rank = Constants.RANK_ARRIVAL;
        }
        maxYearRank = maxYearRankModel.getMaxYearRank();

        dateFrom = Constants.STRING_EMPTY;
        dateTo = Constants.STRING_EMPTY;
        checkedDate = false;
        flagFilterSubmit = true;

        @SuppressLint("SimpleDateFormat") SimpleDateFormat sdf = new SimpleDateFormat(
                Constants.DATE_FORMAT_STRING);
        Date calDateFrom = null, calDateTo = null;
        try {
            calDateFrom = sdf.parse(createDate);
            calDateTo = sdf.parse(createDate);
        } catch (ParseException e) {
            e.printStackTrace();
        }

        Calendar dateFromDefault = Calendar.getInstance(), dateToDefault = Calendar.getInstance();
        dateFromDefault.setTime(calDateFrom);
        dateToDefault.setTime(calDateTo);
        dateFromDefault.add(Calendar.DATE, Constants.DATE_FROM);
        dateToDefault.add(Calendar.DATE, Constants.DATE_TO);

        @SuppressLint("SimpleDateFormat") SimpleDateFormat df = new SimpleDateFormat(
                Constants.DATE_FORMAT_STRING);
        String formattedDateFrom = df.format(dateFromDefault.getTime());
        String formattedDateTo = df.format(dateToDefault.getTime());

        if (dateFrom.isEmpty()) {
            dateFrom = formattedDateFrom;
        }

        if (dateTo.isEmpty()) {
            dateTo = formattedDateTo;
        }

        txvRankName = (TextView) findViewById(R.id.txv_selectionName);
        txvClassifyName = (TextView) findViewById(R.id.txv_classifyName);
        TextView txvID = (TextView) findViewById(R.id.txv_shopId);
        imbRank = (ImageButton) findViewById(R.id.imb_rank);
        imbSelect = (ImageButton) findViewById(R.id.imb_select);
        imbClose = (ImageButton) findViewById(R.id.imb_close);
        imvIcon = (ImageView) findViewById(R.id.imv_icon);

        llTitle = (LinearLayout) findViewById(R.id.ll_title);
        llSelection = (LinearLayout) findViewById(R.id.ll_selection);
        flClose = (FrameLayout) findViewById(R.id.fl_close);
        flFilter = (FrameLayout) findViewById(R.id.fl_filter);
        flSelection = (FrameLayout) findViewById(R.id.fl_selection);
        lstView = (ListView) findViewById(R.id.lsv_list);

        txvLocarionHeader = (TextView) findViewById(R.id.txv_location_header);
        txvNameHeader = (TextView) findViewById(R.id.txv_name_header);
        txvPublisherHeader = (TextView) findViewById(R.id.txv_publisher_header);
        txvPublishDateHeader = (TextView) findViewById(R.id.txv_publish_date_header);
        txvInventoryNumberHeader = (TextView) findViewById(R.id.txv_inventory_number_header);
        txvRankHeader = (TextView) findViewById(R.id.txv_rank_header);

        //Item text view return book
        txvJanCdHeader = (TextView) findViewById(R.id.txv_jan_cd_header);
        txvPublisherReturnHeader = (TextView) findViewById(R.id.txv_publisher_name_return_header);
        txvGroup1NameHeader = (TextView) findViewById(R.id.txv_group1_name_header);
        txvGroup2NameHeader = (TextView) findViewById(R.id.txv_group2_name_header);
        txvClassifyHeader = (TextView) findViewById(R.id.txv_classify_header);

        //Item frame layout return book:
        fl_jan_cd_header = (FrameLayout) findViewById(R.id.fl_jan_cd_header);
        fl_publisher_name_return_header = (FrameLayout) findViewById(R.id.fl_publisher_name_return_header);
        fl_group1_name_header = (FrameLayout) findViewById(R.id.fl_group1_name_header);
        fl_group2_name_header = (FrameLayout) findViewById(R.id.fl_group2_name_header);
        fl_classify_header = (FrameLayout) findViewById(R.id.fl_classify_header);
        fl_publisher_header = (FrameLayout) findViewById(R.id.fl_publisher_header);

        //Item Image view return book
        imvJanCdHeader = (ImageView) findViewById(R.id.imv_jan_cd_header);
        imvPublisherReturnHeader = (ImageView) findViewById(R.id.imv_publisher_name_return_header);
        imvGroup1NameHeader = (ImageView) findViewById(R.id.imv_group1_name_header);
        imvGroup2NameHeader = (ImageView) findViewById(R.id.imv_group2_name_header);

        //Item init
        imvLocationHeader = (ImageView) findViewById(R.id.imv_location_header);
        imvClassifyHeader = (ImageView) findViewById(R.id.imv_classify_header);
        imvPublisherHeader = (ImageView) findViewById(R.id.imv_publisher_header);
        imvPublishDateHeader = (ImageView) findViewById(R.id.imv_publish_date_header);
        imvInventoryNumberHeader = (ImageView) findViewById(R.id.imv_inventory_number_header);
        imvRankHeader = (ImageView) findViewById(R.id.imv_rank_header);
        imvNameHeader = (ImageView) findViewById(R.id.imv_name_header);

        //Item view return book
        v_jan_cd_header = findViewById(R.id.v_jan_cd_header);
        v_publisher_name_return_header = findViewById(R.id.v_publisher_name_return_header);
        v_group1_name_header = findViewById(R.id.v_group1_name_header);
        v_group2_name_header = findViewById(R.id.v_group2_name_header);
        v_classify_header = findViewById(R.id.v_classify_header);
        v_publisher_header = findViewById(R.id.v_publisher_header);

        //Show/hide header
        if (rank == Constants.RANK_RETURN) {
            //show item
            showItemHeaderReturnBook();
        } else {
            hideItemHeaderReturnBook();
        }

        // init data for all show in list view
        id = Constants.ID_ROW_ALL;
        type = Config.TYPE_CLASSIFY;
        offset = 0;
        nameClassify = Constants.STRING_EMPTY;
        yearAgo = Constants.SELECT_ALL_YEAR;
        txvClassifyName.setText(Constants.ROW_ALL);

        txvID.setText(shopID);
        bookModel = new BookModel();
        returnbookModel = new ReturnbookModel();
        periodbookModel = new PeriodbookModel();
        regularbookModel = new RegularbookModel();

        list = new ArrayList<>();
        datalist = new ArrayList<>();
        mapOrder = new LinkedHashMap<>();
        formatCommon = new FormatCommon();

        countLocationClick = Constants.SORT_NULL;
        countNameClick = Constants.SORT_NULL;
        countClassifyClick = Constants.SORT_NULL;
        countPublisherClick = Constants.SORT_NULL;
        countPublishDateClick = Constants.SORT_NULL;
        countInventoryClick = Constants.SORT_NULL;
        countRankClick = Constants.SORT_NULL;

        //return book
        countJanCdClick = Constants.SORT_NULL;
        countGroup1NameClick = Constants.SORT_NULL;
        countGroup2NameClick = Constants.SORT_NULL;

        //Init filter setting return books
        setDefaultFilterSettingReturnBooks();

        // Load data for header title
        loadHeader();

        // Load data for list view
        loadList();

        datalist.addAll(list);

        // Adapter init
        // Set data adapter to list view
        adapter = new ListViewProductAdapter(this, datalist, rank);
        lstView.setAdapter(adapter);


        // set Scroll list view for load more data
        lstView.setOnScrollListener(new EndlessScrollListener() {
            @Override
            public boolean onLoadMore(int page, int totalItemsCount) {
                offset += 1000;
                datalist.clear();
                datalist.addAll(list);
                if (rank == Constants.RANK_RETURN) {
                    loadListReturnWhenSrcoll();
                    adapterReturnBooks.notifyDataSetChanged();
                } else {
                    loadList();
                    adapter.notifyDataSetChanged();
                }
                lstView.requestLayout();
                return true;
            }

        });

        imbSelect.setOnClickListener(this);
        imbClose.setOnClickListener(this);
        imbRank.setOnClickListener(this);
        flClose.setOnClickListener(this);
        txvLocarionHeader.setOnClickListener(this);
        txvPublisherHeader.setOnClickListener(this);
        txvPublishDateHeader.setOnClickListener(this);
        txvClassifyHeader.setOnClickListener(this);
        txvInventoryNumberHeader.setOnClickListener(this);
        txvRankHeader.setOnClickListener(this);
        txvJanCdHeader.setOnClickListener(this);
        txvGroup1NameHeader.setOnClickListener(this);
        txvGroup2NameHeader.setOnClickListener(this);
        txvNameHeader.setOnClickListener(this);
        txvPublisherReturnHeader.setOnClickListener(this);
    }

    /**
     * On Pause app when go use another app
     */
    @Override
    protected void onPause() {
        super.onPause();
        timeout = System.currentTimeMillis();
    }

    /**
     * On Restart when call back this app from another screen move to back Unlock Screen
     */
    @Override
    protected void onRestart() {

        super.onRestart();
        if (timeout < (System.currentTimeMillis() - Constants.TIME_OUT)) {
            Intent intent = new Intent(this, UnlockScreenActivity.class);
            Bundle bundle = new Bundle();
            bundle.putString(Constants.COLUMN_USER_ID, userID);
            bundle.putString(Constants.COLUMN_SHOP_ID, shopID);
            bundle.putString(Constants.COLUMN_SERVER_NAME, serverName);
            bundle.putString(Constants.COLUMN_CREATE_DATE, createDate);
            intent.putExtras(bundle);
            startActivity(intent);
            finish();
        }
    }

    /**
     * even when click back
     */
    @Override
    public void onBackPressed() {
        super.onBackPressed();
        finishAffinity();
    }

    /**
     * Event on click for item logout, rank selected, filter selected
     * or sort by column header
     */
    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.imb_select:
                //Disable item click
                disableItemsSetting();
                //Load itme select
                loadSelect();
                break;
            case R.id.imb_rank:
//                setDefaultColumn();
                loadRank();
                break;
            case R.id.fl_close:
                showAlertDialog();
                break;
            case R.id.imb_close:
                showAlertDialog();
                break;
            case R.id.txv_location_header:
                sortList(Constants.COLUMN_LOCATION_ID);
                break;
            case R.id.txv_name_header:
                sortList(Constants.COLUMN_NAME);
                break;
            case R.id.txv_classify_header:
                sortList(Constants.COLUMN_LARGE_CLASSIFICATION_ID);
                break;
            case R.id.txv_publisher_header:
                sortList(Constants.COLUMN_PUBLISHER_ID);
                break;
            case R.id.txv_publish_date_header:
                sortList(Constants.COLUMN_PUBLISH_DATE);
                break;
            case R.id.txv_inventory_number_header:
                sortList(Constants.COLUMN_INVENTORY_NUMBER);
                break;
            case R.id.txv_rank_header:
                sortList(Constants.COLUMN_RANKING);
                break;
            case R.id.txv_jan_cd_header:
                sortList(Constants.COLUMN_JAN_CD);
                break;
            case R.id.txv_group1_name_header:
                sortList(Constants.COLUMN_MEDIA_GROUP1_CD);
                break;
            case R.id.txv_group2_name_header:
                sortList(Constants.COLUMN_MEDIA_GROUP2_CD);
                break;
            case R.id.txv_publisher_name_return_header:
                sortList(Constants.COLUMN_PUBLISHER_CD);
                break;
        }
    }

    /**
     * Function default filter
     */

    private void setDefaultFilterSettingReturnBooks() {

        //Default flag setting new
        ArrayList<String> arrPublisherCd = new ArrayList<>();
        ArrayList<String> arrPublisherName = new ArrayList<>();
        ArrayList<String> arrGroup1Cd = new ArrayList<>();
        ArrayList<String> arrGroup1Name = new ArrayList<>();
        ArrayList<String> arrGroup2Cd = new ArrayList<>();
        ArrayList<String> arrGroup2Name = new ArrayList<>();
        arrPublisherCd.add(Constants.ID_ROW_ALL);
        arrPublisherName.add(Constants.ROW_ALL);
        arrGroup1Cd.add(Constants.ID_ROW_ALL);
        arrGroup1Name.add(Constants.ROW_ALL);
        arrGroup2Cd.add(Constants.ID_ROW_ALL);
        arrGroup2Name.add(Constants.ROW_ALL);

        //Get list default group1 cd
        List<CLP> listDefaultGroup1 = returnbookModel.getInfoGroupCd1();
        List<CLP> listDefaultGroup2 = returnbookModel.getInfoGroupCd2(Constants.ID_ROW_ALL);
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


        flagSettingNew.setFlagPublisher(arrPublisherCd);
        flagSettingNew.setFlagPublisherShowScreen(arrPublisherName);

        flagSettingNew.setFlagReleaseDate(Constants.FLAG_DEFAULT);
        flagSettingNew.setFlagReleaseDateShowScreen(Constants.FLAG_DEFAULT_RELEASE_DATE_SHOW);

        flagSettingNew.setFlagUndisturbed(Constants.FLAG_DEFAULT);
        flagSettingNew.setFlagUndisturbedShowScreen(Constants.FLAG_DEFAULT_UNDISTURBED_SHOW);

        flagSettingNew.setFlagNumberOfStocks(Constants.FLAG_DEFAULT);
        flagSettingNew.setFlagNumberOfStocksShowScreen(Constants.FLAG_DEFAULT_NUMBER_OF_STOCKS_SHOW);

        flagSettingNew.setFlagStockPercent(Constants.FLAG_DEFAULT);
        flagSettingNew.setFlagStockPercentShowScreen(Constants.FLAG_DEFAULT_STOCKS_PERCENT_SHOW);

        flagSettingNew.setFlagJoubi(Constants.VALUE_YES_STANDING);

        // Default flag setting old
        flagSettingOld.setFlagClassificationGroup1Cd(arrGroup1Cd);
        flagSettingOld.setFlagClassificationGroup1Name(arrGroup1Name);
        flagSettingOld.setFlagClassificationGroup2Cd(arrGroup2Cd);
        flagSettingOld.setFlagClassificationGroup2Name(arrGroup2Name);

        flagSettingOld.setFlagPublisher(arrPublisherCd);
        flagSettingOld.setFlagPublisherShowScreen(arrPublisherName);

        flagSettingOld.setFlagReleaseDate(Constants.FLAG_DEFAULT);
        flagSettingOld.setFlagReleaseDateShowScreen(Constants.FLAG_DEFAULT_RELEASE_DATE_SHOW);

        flagSettingOld.setFlagUndisturbed(Constants.FLAG_DEFAULT);
        flagSettingOld.setFlagUndisturbedShowScreen(Constants.FLAG_DEFAULT_UNDISTURBED_SHOW);

        flagSettingOld.setFlagNumberOfStocks(Constants.FLAG_DEFAULT);
        flagSettingOld.setFlagNumberOfStocksShowScreen(Constants.FLAG_DEFAULT_NUMBER_OF_STOCKS_SHOW);

        flagSettingOld.setFlagStockPercent(Constants.FLAG_DEFAULT);
        flagSettingOld.setFlagStockPercentShowScreen(Constants.FLAG_DEFAULT_STOCKS_PERCENT_SHOW);

        flagSettingOld.setFlagJoubi(Constants.VALUE_YES_STANDING);
    }


    /**
     * Show dialog filter
     */
    private void loadSelect() {

        //Check load rank is return books
        if (rank == Constants.RANK_RETURN) {
            filterSettingListReturnBooks();
        } else {
            filterSettingListOther();
        }

    }

    /**
     * Disable item click when click load select
     */
    //TODO
    private void disableItemsSetting() {

        imbRank.setClickable(false);
        imbSelect.setClickable(false);
        imbClose.setClickable(false);
        flClose.setClickable(false);
    }

    /**
     * Enable item click when click load select
     */
    //TODO
    private void enableItemsSetting() {

        imbRank.setClickable(true);
        imbSelect.setClickable(true);
        imbClose.setClickable(true);
        flClose.setClickable(true);
    }

    private void filterSettingListReturnBooks() {

        if (!flagFilterSubmit) {
            putFlagOldToFlagNew();
        }
        //Call filter setting when click button setting
        DFilterSettingFragment dSettingFragment = new DFilterSettingFragment();
        FragmentManager fm = getSupportFragmentManager();
        Bundle bundle = common.DataPutActivity(flagSettingNew, flagSettingOld);
        //put flag click setting
        bundle.putString(Constants.FLAG_CLICK_SETTING, Constants.VALUE_CHECK_ONCLICK_SETTING);
        dSettingFragment.setArguments(bundle);
        dSettingFragment.show(fm, null);
    }

    //Save flag new into flag old
    private void putFlagOldToFlagNew() {

        flagSettingNew.setFlagClassificationGroup1Cd(flagSettingOld.getFlagClassificationGroup1Cd());
        flagSettingNew.setFlagClassificationGroup1Name(flagSettingOld.getFlagClassificationGroup1Name());
        flagSettingNew.setFlagClassificationGroup2Cd(flagSettingOld.getFlagClassificationGroup2Cd());
        flagSettingNew.setFlagClassificationGroup2Name(flagSettingOld.getFlagClassificationGroup2Name());
        //save flag publisher
        flagSettingNew.setFlagPublisher(flagSettingOld.getFlagPublisher());
        flagSettingNew.setFlagPublisherShowScreen(flagSettingOld.getFlagPublisherShowScreen());
        //save flag release date
        flagSettingNew.setFlagReleaseDate(flagSettingOld.getFlagReleaseDate());
        flagSettingNew.setFlagReleaseDateShowScreen(flagSettingOld.getFlagReleaseDateShowScreen());
        //save flag undisturbed
        flagSettingNew.setFlagUndisturbed(flagSettingOld.getFlagUndisturbed());
        flagSettingNew.setFlagUndisturbedShowScreen(flagSettingOld.getFlagUndisturbedShowScreen());
        //save flag number of stocks
        flagSettingNew.setFlagNumberOfStocks(flagSettingOld.getFlagNumberOfStocks());
        flagSettingNew.setFlagNumberOfStocksShowScreen(flagSettingOld.getFlagNumberOfStocksShowScreen());
        //save flag stocks percent
        flagSettingNew.setFlagStockPercent(flagSettingOld.getFlagStockPercent());
        flagSettingNew.setFlagStockPercentShowScreen(flagSettingOld.getFlagStockPercentShowScreen());
        //save flag joubi
        flagSettingNew.setFlagJoubi(flagSettingOld.getFlagJoubi());
    }

    private void filterSettingListOther() {

        DSelectFragment dSelectFragment = new DSelectFragment();
        FragmentManager fm = getSupportFragmentManager();
        Bundle bundle = new Bundle();
        bundle.putInt(Config.TYPE, type);
        bundle.putString(Constants.COLUMN_ID, id);
        bundle.putString(Constants.COLUMN_LARGE_CLASSIFICATION_NAME, nameClassify);
        bundle.putInt(Constants.YEAR_AGO, yearAgo);
        bundle.putInt(Constants.RANK, rank);
        bundle.putBoolean(Constants.FLAG_DATE_CHECK, checkedDate);
        bundle.putString(Constants.COLUMN_DATE_FROM, dateFrom);
        bundle.putString(Constants.COLUMN_DATE_TO, dateTo);

        //Flag return book
        bundle.putInt(Constants.FLAG_PERCENT_SELECTED, percent);
        bundle.putString(Constants.FLAG_SELECT_GROUP1_CD, flagGroup1Cd);
        bundle.putString(Constants.FLAG_SELECT_GROUP2_CD, flagGroup2Cd);
        bundle.putString(Constants.FLAG_SELECT_GROUP2_NAME, flagGroup2Name);
        dSelectFragment.setArguments(bundle);
        dSelectFragment.show(fm, null);
    }


    private void checkRankReturnBook() {
        if (rank == Constants.RANK_RETURN) {
            showItemHeaderReturnBook();
        } else {
            hideItemHeaderReturnBook();
        }
    }

    /**
     * Show dialog rank select
     */
    private void loadRank() {

        DRankFragment dRankFragment = new DRankFragment();
        FragmentManager fm = getSupportFragmentManager();
        Bundle bundle = new Bundle();
        bundle.putInt(Constants.RANK, rank);
        dRankFragment.setArguments(bundle);
        dRankFragment.show(fm, null);
    }

    private void showItemHeaderReturnBook() {
        //show item
        fl_jan_cd_header.setVisibility(FrameLayout.VISIBLE);
        fl_publisher_name_return_header.setVisibility(FrameLayout.VISIBLE);
        fl_group1_name_header.setVisibility(FrameLayout.VISIBLE);
        fl_group2_name_header.setVisibility(FrameLayout.VISIBLE);

        v_jan_cd_header.setVisibility(FrameLayout.VISIBLE);
        v_publisher_name_return_header.setVisibility(FrameLayout.VISIBLE);
        v_group1_name_header.setVisibility(FrameLayout.VISIBLE);
        v_group2_name_header.setVisibility(FrameLayout.VISIBLE);

        //hide item
        fl_classify_header.setVisibility(FrameLayout.GONE);
        fl_publisher_header.setVisibility(FrameLayout.GONE);

        v_classify_header.setVisibility(FrameLayout.GONE);
        v_publisher_header.setVisibility(FrameLayout.GONE);
    }

    private void hideItemHeaderReturnBook() {

        //hide item
        fl_jan_cd_header.setVisibility(FrameLayout.GONE);
        fl_publisher_name_return_header.setVisibility(FrameLayout.GONE);
        fl_group1_name_header.setVisibility(FrameLayout.GONE);
        fl_group2_name_header.setVisibility(FrameLayout.GONE);

        v_jan_cd_header.setVisibility(FrameLayout.GONE);
        v_publisher_name_return_header.setVisibility(FrameLayout.GONE);
        v_group1_name_header.setVisibility(FrameLayout.GONE);
        v_group2_name_header.setVisibility(FrameLayout.GONE);

        //show item
        fl_classify_header.setVisibility(FrameLayout.VISIBLE);
        fl_publisher_header.setVisibility(FrameLayout.VISIBLE);

        v_classify_header.setVisibility(FrameLayout.VISIBLE);
        v_publisher_header.setVisibility(FrameLayout.VISIBLE);
    }

    private void loadListReturnWhenSrcoll() {

        List<Returnbooks> returnbooksList;
        if (flagSettingNew.getFlagClassificationGroup1Cd().size() > 0 &&
                Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup1Cd().get(0))) {
            returnbooksList = returnbookModel.getListBookInfo(offset, treeMapOrder, flagSettingNew);
        } else {
            returnbooksList = returnbookModel.getListBookInfoSelectGroup1Cd(offset, treeMapOrder);
        }
        putDataListReturn(returnbooksList);

        //adapterReturnBooks = new ListViewProductReturnBooksAdapter(this, list);
        //lstView.setAdapter(adapterReturnBooks);
    }

    /**
     * Set item into list view
     */
    private void loadList() {

        List<Books> booksList;
        List<Books> periodbooksList;
        List<Books> regularbooksList;


        // check is rank return
        if (rank != Constants.RANK_RETURN && rank != Constants.RANK_PERIOD
                && rank != Constants.RANK_REGULAR) {
            if (!bookModel.checkData()) {
                return;
            }
            //Book list
            booksList = bookModel.getListBookInfo(id, type, offset, rank, treeMapOrder);

            for (int i = 0; i < booksList.size(); i++) {
                HashMap<String, String> hashMap = new HashMap<>();
                hashMap.put(Constants.COLUMN_LOCATION_ID, formatCommon.formatLocationIdNewLine(booksList.get(i).getLocation_id()));
                hashMap.put(Constants.COLUMN_NAME, booksList.get(i).getName());
                hashMap.put(Constants.COLUMN_LARGE_CLASSIFICATION_NAME,
                        booksList.get(i).getLarge_classifications_name());
                hashMap.put(Constants.COLUMN_PUBLISHER_NAME,
                        booksList.get(i).getPublisher_name());
                hashMap.put(Constants.COLUMN_PUBLISH_DATE, booksList.get(i).getPublish_date());
                hashMap.put(Constants.COLUMN_INVENTORY_NUMBER,
                        String.valueOf(booksList.get(i).getInventory_number()));
                hashMap.put(Constants.COLUMN_NEW_CATEGORY_RANK,
                        String.valueOf(booksList.get(i).getNew_catagory_rank()));
                hashMap.put(Constants.COLUMN_RANKING,
                        String.valueOf(booksList.get(i).getRanking()));
                list.add(hashMap);
            }
        } else if (rank == Constants.RANK_PERIOD) {
            // Period book list
            periodbooksList = periodbookModel.getListViewBookInfo(id, type, offset, dateFrom,
                    dateTo, treeMapOrder);

            for (int i = 0; i < periodbooksList.size(); i++) {
                HashMap<String, String> hashMap = new HashMap<>();
                hashMap.put(Constants.COLUMN_LOCATION_ID, formatCommon.formatLocationIdNewLine(periodbooksList.get(i).getLocation_id()));
                hashMap.put(Constants.COLUMN_NAME, periodbooksList.get(i).getName());
                hashMap.put(Constants.COLUMN_LARGE_CLASSIFICATION_NAME,
                        periodbooksList.get(i).getLarge_classifications_name());
                hashMap.put(Constants.COLUMN_PUBLISHER_NAME,
                        periodbooksList.get(i).getPublisher_name());
                hashMap.put(Constants.COLUMN_PUBLISH_DATE,
                        periodbooksList.get(i).getPublish_date());
                hashMap.put(Constants.COLUMN_INVENTORY_NUMBER,
                        String.valueOf(periodbooksList.get(i).getInventory_number()));
                hashMap.put(Constants.COLUMN_NEW_CATEGORY_RANK,
                        String.valueOf(periodbooksList.get(i).getNew_catagory_rank()));
                hashMap.put(Constants.COLUMN_RANKING,
                        String.valueOf(periodbooksList.get(i).getRanking()));
                list.add(hashMap);
            }
        } else if (rank == Constants.RANK_REGULAR) {
            // Period book list
            regularbooksList = regularbookModel.getListViewBookInfo(id, type, offset, treeMapOrder);

            for (int i = 0; i < regularbooksList.size(); i++) {
                HashMap<String, String> hashMap = new HashMap<>();
                hashMap.put(Constants.COLUMN_LOCATION_ID, formatCommon.formatLocationIdNewLine(regularbooksList.get(i).getLocation_id()));
                hashMap.put(Constants.COLUMN_NAME, regularbooksList.get(i).getName());
                hashMap.put(Constants.COLUMN_LARGE_CLASSIFICATION_NAME,
                        regularbooksList.get(i).getLarge_classifications_name());
                hashMap.put(Constants.COLUMN_PUBLISHER_NAME,
                        regularbooksList.get(i).getPublisher_name());
                hashMap.put(Constants.COLUMN_PUBLISH_DATE,
                        regularbooksList.get(i).getPublish_date());
                hashMap.put(Constants.COLUMN_INVENTORY_NUMBER,
                        String.valueOf(regularbooksList.get(i).getInventory_number()));
                hashMap.put(Constants.COLUMN_NEW_CATEGORY_RANK,
                        String.valueOf(regularbooksList.get(i).getNew_catagory_rank()));
                hashMap.put(Constants.COLUMN_RANKING,
                        String.valueOf(regularbooksList.get(i).getRanking()));
                list.add(hashMap);
            }
        } else {
            loadListItemsReturnBooks();
        }
    }

    /**
     * Function pull data when not select all group1 cd
     */

    private void putDataListReturn(List<Returnbooks> returnbooksList) {

        for (int i = 0; i < returnbooksList.size(); i++) {
            HashMap<String, String> hashMap = new HashMap<>();
            hashMap.put(Constants.COLUMN_LOCATION_ID,
                    formatCommon.formatLocationIdNewLine(returnbooksList.get(i).getLocation_id()));
            hashMap.put(Constants.COLUMN_JAN_CD,
                    returnbooksList.get(i).getJan_cd());
            hashMap.put(Constants.COLUMN_GOODS_NAME,
                    returnbooksList.get(i).getBqgm_goods_name());
            hashMap.put(Constants.COLUMN_PUBLISHER_CD,
                    returnbooksList.get(i).getBqgm_publisher_cd());
            hashMap.put(Constants.COLUMN_PUBLISHER_NAME_RETURN,
                    returnbooksList.get(i).getBqgm_publisher_name());
            hashMap.put(Constants.COLUMN_MEDIA_GROUP1_CD,
                    returnbooksList.get(i).getBqct_media_group1_cd());
            hashMap.put(Constants.COLUMN_MEDIA_GROUP1_NAME,
                    returnbooksList.get(i).getBqct_media_group1_name());
            hashMap.put(Constants.COLUMN_MEDIA_GROUP2_CD,
                    returnbooksList.get(i).getBqct_media_group2_cd());
            hashMap.put(Constants.COLUMN_MEDIA_GROUP2_NAME,
                    returnbooksList.get(i).getBqct_media_group2_name());
            hashMap.put(Constants.COLUMN_SALES_DATE,
                    formatCommon.formatDateBlank(returnbooksList.get(i).getBqgm_sales_date()));
            hashMap.put(Constants.COLUMN_STOCK_COUNT,
                    String.valueOf(returnbooksList.get(i).getBqsc_stock_count()));
            hashMap.put(Constants.COLUMN_YEAR_RANK,
                    String.valueOf(returnbooksList.get(i).getYear_rank()));
            hashMap.put(Constants.COLUMN_WRITER_NAME,
                    returnbooksList.get(i).getBqgm_writer_name());
            hashMap.put(Constants.COLUMN_PRICE,
                    String.valueOf(returnbooksList.get(i).getBqgm_price()));
            hashMap.put(Constants.COLUMN_FIRST_SUPPLY_DATE,
                    formatCommon.formatDateBlank(String.valueOf(returnbooksList.get(i).getBqtse_first_supply_date())));
            hashMap.put(Constants.COLUMN_LAST_SUPPLY_DATE,
                    formatCommon.formatDateBlank(returnbooksList.get(i).getBqtse_last_supply_date()));
            hashMap.put(Constants.COLUMN_LAST_ORDER_DATE,
                    formatCommon.formatDateBlank(returnbooksList.get(i).getBqtse_last_order_date()));
            hashMap.put(Constants.COLUMN_LAST_SALES_DATE,
                    formatCommon.formatDateBlank(returnbooksList.get(i).getBqtse_last_sale_date()));
            hashMap.put(Constants.COLUMN_TOTAL_SALES,
                    String.valueOf(returnbooksList.get(i).getSts_total_sales()));
            hashMap.put(Constants.COLUMN_TOTAL_SUPPLY,
                    String.valueOf(returnbooksList.get(i).getSts_total_supply()));
            hashMap.put(Constants.COLUMN_TOTAL_RETURN,
                    String.valueOf(returnbooksList.get(i).getSts_total_return()));
            hashMap.put(Constants.COLUMN_JOUBI,
                    String.valueOf(returnbooksList.get(i).getJoubi()));
            list.add(hashMap);
        }
    }


    /**
     * Clear data and log out
     */

    private void clearAndLogout() {

        // Process logout
        // clear table
        DatabaseManager.initializeInstance(new DatabaseHelper(getApplicationContext()));
        DatabaseHelper ds = new DatabaseHelper(this);
        ds.clearTables();

        // stop process loading screen
        progress.dismiss();

        finishAffinity();
        // move to login screen
        Intent intent = new Intent(this, LoginActivity.class);

        startActivity(intent);
    }

    /**
     * Show dialog warning logout
     */
    public void showAlertDialog() {

        progress.show();
        android.support.v7.app.AlertDialog.Builder dialog =
                new android.support.v7.app.AlertDialog.Builder(this);
        dialog.setCancelable(false);

        dialog
                .setMessage(getString(R.string.logout_msg))
                .setPositiveButton(getString(R.string.logout_yes),
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {

                                LogManager.i(TAG,
                                        Message.TAG_UNLOCK_ACTIVITY + Message.MESSAGE_LOGOUT);
                                // print log end process
                                LogManager.i(TAG, Message.TAG_UNLOCK_ACTIVITY
                                        + Message.MESSAGE_ACTIVITY_END);
                                // print log move screen
                                LogManager.i(TAG,
                                        String.format(Message.MESSAGE_ACTIVITY_MOVE,
                                                Message.UNLOCK_ACTIVITY_NAME,
                                                Message.LOGIN_ACTIVITY_NAME));
                                clearAndLogout();
                            }
                        })
                .setNegativeButton(getString(R.string.logout_no),
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                dialog.dismiss();
                                progress.dismiss();
                            }
                        });


        android.support.v7.app.AlertDialog alert = dialog.show();
        TextView messageText = (TextView) alert.findViewById(android.R.id.message);
        assert messageText != null;
        messageText.setGravity(Gravity.CENTER);
    }

    /**
     * Load icon filter
     */
    private void loadIcon() {
        switch (type) {
            case Config.TYPE_CLASSIFY:
                imvIcon.setImageResource(R.drawable.book01);
                break;
            case Config.TYPE_PUBLISHER:
                imvIcon.setImageResource(R.drawable.publish01);
        }
    }

    /**
     * Load header text and color by rank selected
     */
    private void loadHeader() {
        switch (rank) {
            case Constants.RANK_PLATFORM1:
                txvRankName.setText(getResources().getString(R.string.rank_platform1));
                flClose.setBackgroundColor(getResources().getColor(R.color.colorPink70));
                llTitle.setBackgroundColor(getResources().getColor(R.color.colorPink63));
                flFilter.setBackgroundColor(getResources().getColor(R.color.colorPink46));
                flSelection.setBackgroundColor(getResources().getColor(R.color.colorPink33));
                llSelection.setVisibility(View.VISIBLE);
                break;
            case Constants.RANK_PLATFORM2:
                txvRankName.setText(getResources().getString(R.string.rank_platform2));
                flClose.setBackgroundColor(getResources().getColor(R.color.colorBlue70));
                llTitle.setBackgroundColor(getResources().getColor(R.color.colorBlue63));
                flFilter.setBackgroundColor(getResources().getColor(R.color.colorBlue46));
                flSelection.setBackgroundColor(getResources().getColor(R.color.colorBlue33));
                llSelection.setVisibility(View.VISIBLE);
                break;
            case Constants.RANK_SURFACE:
                txvRankName.setText(getResources().getString(R.string.rank_surface));
                flClose.setBackgroundColor(getResources().getColor(R.color.colorGreen70));
                llTitle.setBackgroundColor(getResources().getColor(R.color.colorGreen63));
                flFilter.setBackgroundColor(getResources().getColor(R.color.colorGreen46));
                flSelection.setBackgroundColor(getResources().getColor(R.color.colorGreen33));
                llSelection.setVisibility(View.VISIBLE);
                break;
            case Constants.RANK_SHELDER:
                txvRankName.setText(getResources().getString(R.string.rank_shelter));
                flClose.setBackgroundColor(getResources().getColor(R.color.colorYellow70));
                llTitle.setBackgroundColor(getResources().getColor(R.color.colorYellow63));
                flFilter.setBackgroundColor(getResources().getColor(R.color.colorYellow46));
                flSelection.setBackgroundColor(getResources().getColor(R.color.colorYellow33));
                llSelection.setVisibility(View.VISIBLE);
                break;
            case Constants.RANK_RETURN:
                txvRankName.setText(getResources().getString(R.string.rank_return));
                flClose.setBackgroundColor(getResources().getColor(R.color.colorOrange70));
                llTitle.setBackgroundColor(getResources().getColor(R.color.colorOrange63));
                flFilter.setBackgroundColor(getResources().getColor(R.color.colorOrange46));
                flSelection.setBackgroundColor(getResources().getColor(R.color.colorOrange33));
                llSelection.setVisibility(View.VISIBLE);
                break;
            case Constants.RANK_ARRIVAL:
                txvRankName.setText(getResources().getString(R.string.rank_arrival));
                flClose.setBackgroundColor(getResources().getColor(R.color.colorCyan70));
                llTitle.setBackgroundColor(getResources().getColor(R.color.colorCyan63));
                flFilter.setBackgroundColor(getResources().getColor(R.color.colorCyan46));
                flSelection.setBackgroundColor(getResources().getColor(R.color.colorCyan33));
                llSelection.setVisibility(View.VISIBLE);
                break;
            case Constants.RANK_PERIOD:
                txvRankName.setText(getResources().getString(R.string.rank_period));
                flClose.setBackgroundColor(getResources().getColor(R.color.colorBrown65));
                llTitle.setBackgroundColor(getResources().getColor(R.color.colorBrown32));
                flFilter.setBackgroundColor(getResources().getColor(R.color.colorBrown19));
                flSelection.setBackgroundColor(getResources().getColor(R.color.colorBrown16));
                llSelection.setVisibility(View.VISIBLE);
                break;
            case Constants.RANK_REGULAR:
                txvRankName.setText(getResources().getString(R.string.rank_regular));
                flClose.setBackgroundColor(getResources().getColor(R.color.colorGreenGrey46));
                llTitle.setBackgroundColor(getResources().getColor(R.color.colorGreenGrey34));
                flFilter.setBackgroundColor(getResources().getColor(R.color.colorGreenGrey18));
                flSelection.setBackgroundColor(getResources().getColor(R.color.colorGreenGrey10));
                llSelection.setVisibility(View.VISIBLE);
                break;
            default:
                txvRankName.setText(getResources().getString(R.string.rank_arrival));
                flClose.setBackgroundColor(getResources().getColor(R.color.colorCyan70));
                llTitle.setBackgroundColor(getResources().getColor(R.color.colorCyan63));
                flFilter.setBackgroundColor(getResources().getColor(R.color.colorCyan46));
                flSelection.setBackgroundColor(getResources().getColor(R.color.colorCyan33));
                llSelection.setVisibility(View.VISIBLE);
                break;
        }
    }

    /**
     * Rank selected result
     *
     * @param itemSelected is rank selected
     */
    @Override
    public void onRankSelectedDialog(int itemSelected) {

        if (itemSelected != Constants.RANK_PERIOD) {
            checkedDate = false;
        }

        saveStateList(rank);
        if (itemSelected == Constants.RANK_RETURN) {
            if (itemSelected != rank) {
                //TODO
                //Reset rank return book
                percent = Constants.SELECT_PERCENT_30;
                flagGroup1Cd = Constants.ID_ROW_ALL;
                flagGroup2Cd = Constants.ID_ROW_ALL;
                flagGroup2Name = Constants.ROW_ALL;
                txvClassifyName.setText(flagSettingNew.getFlagStockPercentShowScreen());
            }
        } else {
            if (rank == Constants.RANK_RETURN) {
                // init data for all show in list view
                id = Constants.ID_ROW_ALL;
                type = Config.TYPE_CLASSIFY;
                offset = 0;
                nameClassify = Constants.STRING_EMPTY;
                yearAgo = Constants.SELECT_ALL_YEAR;
                txvClassifyName.setText(Constants.ROW_ALL);
            }
        }
        rank = itemSelected;
        offset = 0;
        yearAgo = Constants.SELECT_ALL_YEAR;

        loadHeader();

        list.clear();

        loadList();
        if (rank == Constants.RANK_RETURN) {
            imvIcon.setImageResource(R.drawable.book01);
        } else {
            loadIcon();
        }

        //Check rank return book loading
        checkRankReturnBook();

        //Check rank return
        if (rank != Constants.RANK_RETURN) {
            adapter = new ListViewProductAdapter(this, list, rank);
            lstView.setAdapter(adapter);
        }
        restoreStateList(rank);

    }

    /**
     * Sort list with column on click
     *
     * @param nameSort column name for sort
     */
    private void sortList(String nameSort) {

        //Change name when rank is return book
        if (rank == Constants.RANK_RETURN) {
            if (Constants.COLUMN_NAME.equals(nameSort)) {
                nameSort = Constants.COLUMN_GOODS_NAME;
            } else if (Constants.COLUMN_PUBLISHER_ID.equals(nameSort)) {
                nameSort = Constants.COLUMN_PUBLISHER_CD;
            } else if (Constants.COLUMN_PUBLISH_DATE.equals(nameSort)) {
                nameSort = Constants.COLUMN_SALES_DATE;
            } else if (Constants.COLUMN_INVENTORY_NUMBER.equals(nameSort)) {
                nameSort = Constants.COLUMN_STOCK_COUNT;
            } else if (Constants.COLUMN_RANKING.equals(nameSort)) {
                nameSort = Constants.COLUMN_YEAR_RANK;
            }
        }
        switch (nameSort) {
            // Column Location Product
            case Constants.COLUMN_LOCATION_ID:
                txvLocarionHeader.setBackgroundColor(getResources().getColor(R.color.colorGrey52));
                txvLocarionHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                if (countLocationClick > Constants.SORT_ASC) {
                    countLocationClick = Constants.SORT_NULL;
                }
                showImageSort(countLocationClick, imvLocationHeader, txvLocarionHeader);
                countLocationClick += 1;
                setSortString(nameSort, countLocationClick);
                break;
            // Column Name Product
            case Constants.COLUMN_NAME:
                txvNameHeader.setBackgroundColor(getResources().getColor(R.color.colorGrey52));
                txvNameHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                if (countNameClick > Constants.SORT_ASC) {
                    countNameClick = Constants.SORT_NULL;
                }
                showImageSort(countNameClick, imvNameHeader, txvNameHeader);
                countNameClick += 1;
                setSortString(nameSort, countNameClick);
                break;
            // Column Classify id
            case Constants.COLUMN_LARGE_CLASSIFICATION_ID:
                txvClassifyHeader.setBackgroundColor(getResources().getColor(R.color.colorGrey52));
                txvClassifyHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                if (countClassifyClick > Constants.SORT_ASC) {
                    countClassifyClick = Constants.SORT_NULL;
                }
                showImageSort(countClassifyClick, imvClassifyHeader, txvClassifyHeader);
                countClassifyClick += 1;
                setSortString(nameSort, countClassifyClick);
                break;
            // Column Publish Id
            case Constants.COLUMN_PUBLISHER_ID:
                txvPublisherHeader.setBackgroundColor(getResources().getColor(R.color.colorGrey52));
                txvPublisherHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                txvPublisherReturnHeader.setBackgroundColor(getResources().getColor(R.color.colorGrey52));
                txvPublisherReturnHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                if (countPublisherClick > Constants.SORT_ASC) {
                    countPublisherClick = Constants.SORT_NULL;
                }
                showImageSort(countPublisherClick, imvPublisherHeader, txvPublisherHeader);
                showImageSort(countPublisherClick, imvPublisherReturnHeader, txvPublisherReturnHeader);
                countPublisherClick += 1;
                setSortString(nameSort, countPublisherClick);
                break;
            // Column Publish Date
            case Constants.COLUMN_PUBLISH_DATE:
                txvPublishDateHeader.setBackgroundColor(
                        getResources().getColor(R.color.colorGrey52));
                txvPublishDateHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                if (countPublishDateClick > Constants.SORT_ASC) {
                    countPublishDateClick = Constants.SORT_NULL;
                }
                showImageSort(countPublishDateClick, imvPublishDateHeader, txvPublishDateHeader);
                countPublishDateClick += 1;
                setSortString(nameSort, countPublishDateClick);
                break;
            // Column Inventory Number
            case Constants.COLUMN_INVENTORY_NUMBER:
                txvInventoryNumberHeader.setBackgroundColor(
                        getResources().getColor(R.color.colorGrey52));
                txvInventoryNumberHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                if (countInventoryClick > Constants.SORT_ASC) {
                    countInventoryClick = Constants.SORT_NULL;
                }
                showImageSort(countInventoryClick, imvInventoryNumberHeader,
                        txvInventoryNumberHeader);
                countInventoryClick += 1;
                setSortString(nameSort, countInventoryClick);
                break;
            // Column Ranking
            case Constants.COLUMN_RANKING:
                txvRankHeader.setBackgroundColor(
                        getResources().getColor(R.color.colorGrey52));
                txvRankHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                if (countRankClick > Constants.SORT_ASC) {
                    countRankClick = Constants.SORT_NULL;
                }
                showImageSort(countRankClick, imvRankHeader, txvRankHeader);
                countRankClick += 1;
                setSortString(nameSort, countRankClick);
                break;

            // Column Jan Cd
            case Constants.COLUMN_JAN_CD:
                txvJanCdHeader.setBackgroundColor(getResources().getColor(R.color.colorGrey52));
                txvJanCdHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                if (countJanCdClick > Constants.SORT_ASC) {
                    countJanCdClick = Constants.SORT_NULL;
                }
                showImageSort(countJanCdClick, imvJanCdHeader, txvJanCdHeader);
                countJanCdClick += 1;
                setSortString(nameSort, countJanCdClick);
                break;
            // Column stock count
            case Constants.COLUMN_STOCK_COUNT:
                txvInventoryNumberHeader.setBackgroundColor(getResources().getColor(R.color.colorGrey52));
                txvInventoryNumberHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                if (countInventoryClick > Constants.SORT_ASC) {
                    countInventoryClick = Constants.SORT_NULL;
                }
                showImageSort(countInventoryClick, imvInventoryNumberHeader, txvInventoryNumberHeader);
                countInventoryClick += 1;
                setSortString(nameSort, countInventoryClick);
                break;
            // Column publisher return
            case Constants.COLUMN_PUBLISHER_CD:
                txvPublisherHeader.setBackgroundColor(getResources().getColor(R.color.colorGrey52));
                txvPublisherHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                txvPublisherReturnHeader.setBackgroundColor(getResources().getColor(R.color.colorGrey52));
                txvPublisherReturnHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                if (countPublisherClick > Constants.SORT_ASC) {
                    countPublisherClick = Constants.SORT_NULL;
                }
                showImageSort(countPublisherClick, imvPublisherHeader, txvPublisherHeader);
                showImageSort(countPublisherClick, imvPublisherReturnHeader, txvPublisherReturnHeader);
                countPublisherClick += 1;
                setSortString(nameSort, countPublisherClick);
                break;
            // Column publisher return
            case Constants.COLUMN_GOODS_NAME:
                txvNameHeader.setBackgroundColor(getResources().getColor(R.color.colorGrey52));
                txvNameHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                if (countNameClick > Constants.SORT_ASC) {
                    countNameClick = Constants.SORT_NULL;
                }
                showImageSort(countNameClick, imvNameHeader, txvNameHeader);
                countNameClick += 1;
                setSortString(nameSort, countNameClick);
                break;
            // Column sales date
            case Constants.COLUMN_SALES_DATE:
                txvPublishDateHeader.setBackgroundColor(getResources().getColor(R.color.colorGrey52));
                txvPublishDateHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                if (countPublishDateClick > Constants.SORT_ASC) {
                    countPublishDateClick = Constants.SORT_NULL;
                }
                showImageSort(countPublishDateClick, imvPublishDateHeader, txvPublishDateHeader);
                countPublishDateClick += 1;
                setSortString(nameSort, countPublishDateClick);
                break;
            // Column group 1 cd click
            case Constants.COLUMN_MEDIA_GROUP1_CD:
                txvGroup1NameHeader.setBackgroundColor(getResources().getColor(R.color.colorGrey52));
                txvGroup1NameHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                if (countGroup1NameClick > Constants.SORT_ASC) {
                    countGroup1NameClick = Constants.SORT_NULL;
                }
                showImageSort(countGroup1NameClick, imvGroup1NameHeader, txvGroup1NameHeader);
                countGroup1NameClick += 1;
                setSortString(nameSort, countGroup1NameClick);
                break;
            // Column group 2 cd click
            case Constants.COLUMN_MEDIA_GROUP2_CD:
                txvGroup2NameHeader.setBackgroundColor(getResources().getColor(R.color.colorGrey52));
                txvGroup2NameHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                if (countGroup2NameClick > Constants.SORT_ASC) {
                    countGroup2NameClick = Constants.SORT_NULL;
                }
                showImageSort(countGroup2NameClick, imvGroup2NameHeader, txvGroup2NameHeader);
                countGroup2NameClick += 1;
                setSortString(nameSort, countGroup2NameClick);
                break;
            // Colum year rank return
            case Constants.COLUMN_YEAR_RANK:
                txvRankHeader.setBackgroundColor(
                        getResources().getColor(R.color.colorGrey52));
                txvRankHeader.setTextColor(getResources().getColor(R.color.colorWhite));
                if (countRankClick > Constants.SORT_ASC) {
                    countRankClick = Constants.SORT_NULL;
                }
                showImageSort(countRankClick, imvRankHeader, txvRankHeader);
                countRankClick += 1;
                setSortString(nameSort, countRankClick);
                break;
        }

        offset = 0;
        lstView.invalidate();
        if (rank != Constants.RANK_RETURN) {
            list.clear();
        }

        // refresh list
        loadList();
        // Adapter init
        // Set data adapter to list view
        if (rank != Constants.RANK_RETURN) {
            ListViewProductAdapter adapter = new ListViewProductAdapter(this, list, rank);
            lstView.setAdapter(adapter);
        }
    }

    /**
     * Set icon sort
     *
     * @param sort {@link int}
     * @param v    {@link ImageView}
     * @param tv   {@link TextView}
     */

    private void showImageSort(int sort, ImageView v, TextView tv) {
        switch (sort) {
            case Constants.SORT_NULL:
                v.setVisibility(View.VISIBLE);
                break;
            case Constants.SORT_DESC:
                v.setImageResource(R.drawable.ic_arrow_drop_up);
                break;
            case Constants.SORT_ASC:
                v.setImageResource(R.drawable.ic_arrow_drop_down);
                v.setVisibility(View.GONE);
                tv.setBackgroundColor(Color.TRANSPARENT);
                tv.setTextColor(getResources().getColor(R.color.colorBlack));
                break;
        }
    }

    /**
     * Set name and order to hash map
     *
     * @param name      {@link String}
     * @param sortClick {@link int}
     */
    private void setSortString(String name, int sortClick) {

        //Change name return book and other book
        if (Constants.COLUMN_GOODS_NAME.equals(name)) {
            name = Constants.COLUMN_NAME;
        } else if (Constants.COLUMN_PUBLISHER_CD.equals(name)) {
            name = Constants.COLUMN_PUBLISHER_ID;
        } else if (Constants.COLUMN_SALES_DATE.equals(name)) {
            name = Constants.COLUMN_PUBLISH_DATE;
        } else if (Constants.COLUMN_STOCK_COUNT.equals(name)) {
            name = Constants.COLUMN_INVENTORY_NUMBER;
        }
        Map<Integer, String> defaultMapSort = new LinkedHashMap<>();
        if (sortClick == Constants.SORT_DESC) {
            mapOrder.put(name, Constants.ORDER_DESC);
        } else if (sortClick == Constants.SORT_ASC) {
            mapOrder.put(name, Constants.ORDER_ASC);
        } else {
            mapOrder.remove(name);
        }
        //Sort map
        for (String key : mapOrder.keySet()) {
            if (rank != Constants.RANK_RETURN) {
                switch (key) {
                    case Constants.COLUMN_LOCATION_ID:
                        defaultMapSort.put(Constants.NUMBER_1, mapOrder.get(key));
                        break;
                    case Constants.COLUMN_NAME:
                        defaultMapSort.put(Constants.NUMBER_3, mapOrder.get(key));
                        break;
                    case Constants.COLUMN_LARGE_CLASSIFICATION_ID:
                        defaultMapSort.put(Constants.NUMBER_4, mapOrder.get(key));
                        break;
                    case Constants.COLUMN_PUBLISHER_ID:
                        defaultMapSort.put(Constants.NUMBER_5, mapOrder.get(key));
                        break;
                    case Constants.COLUMN_PUBLISH_DATE:
                        defaultMapSort.put(Constants.NUMBER_8, mapOrder.get(key));
                        break;
                    case Constants.COLUMN_INVENTORY_NUMBER:
                        defaultMapSort.put(Constants.NUMBER_9, mapOrder.get(key));
                        break;
                    case Constants.COLUMN_RANKING:
                        defaultMapSort.put(Constants.NUMBER_10, mapOrder.get(key));
                        break;
                }
            } else {
                switch (key) {
                    case Constants.COLUMN_LOCATION_ID:
                        defaultMapSort.put(Constants.NUMBER_1, mapOrder.get(key));
                        break;
                    case Constants.COLUMN_JAN_CD:
                        defaultMapSort.put(Constants.NUMBER_2, mapOrder.get(key));
                        break;
                    case Constants.COLUMN_NAME:
                        defaultMapSort.put(Constants.NUMBER_3, mapOrder.get(key));
                        break;
                    case Constants.COLUMN_PUBLISHER_ID:
                        defaultMapSort.put(Constants.NUMBER_5, mapOrder.get(key));
                        break;
                    case Constants.COLUMN_MEDIA_GROUP1_CD:
                        defaultMapSort.put(Constants.NUMBER_6, mapOrder.get(key));
                        break;
                    case Constants.COLUMN_MEDIA_GROUP2_CD:
                        defaultMapSort.put(Constants.NUMBER_7, mapOrder.get(key));
                        break;
                    case Constants.COLUMN_PUBLISH_DATE:
                        defaultMapSort.put(Constants.NUMBER_8, mapOrder.get(key));
                        break;
                    case Constants.COLUMN_INVENTORY_NUMBER:
                        defaultMapSort.put(Constants.NUMBER_9, mapOrder.get(key));
                        break;
                    case Constants.COLUMN_YEAR_RANK:
                        defaultMapSort.put(Constants.NUMBER_10, mapOrder.get(key));
                        break;
                }
            }
        }
        treeMapOrder = new TreeMap<>(defaultMapSort);
    }

    /**
     * Select filter dialog call back
     *
     * @param idSelected           is id call back
     * @param nameClassifySelected is name Classify or Publisher call back
     * @param typeSelected         is type Classify or Publisher call back
     * @param dateChecked          is check date
     */
    @Override
    public void onLitSelectedDialog(int typeSelected, String idSelected,
                                    String nameClassifySelected, boolean dateChecked) {

        //Enable click items setting
        enableItemsSetting();

        offset = 0;
        type = typeSelected;
        nameClassify = nameClassifySelected;
        checkedDate = dateChecked;
        if (!idSelected.contains(Constants.ID_ROW_ALL)) {
            id = idSelected;
            txvClassifyName.setText(nameClassify);
        } else {
            id = Constants.ID_ROW_ALL;
            txvClassifyName.setText(Constants.ROW_ALL);
        }

        // refresh list view
        lstView.invalidate();
        list.clear();
        loadList();
        loadIcon();

        // Adapter init
        // Set data adapter to list view
        ListViewProductAdapter adapter = new ListViewProductAdapter(this, list, rank);
        lstView.setAdapter(adapter);
    }

    /**
     * Result filter selected in select dialog
     *
     * @param typeSelected         is type selected for filter
     * @param idSelected           is id selected of type Classify or Publisher for filter
     * @param nameClassifySelected is name selected of type Classify or Publisher for filter
     * @param yearSelected         if rank is return will filter with year selected
     */
    @Override
    public void onSelectedFilterDialog(int typeSelected, String idSelected,
                                       String nameClassifySelected, int yearSelected) {

        //Enable click items setting
        enableItemsSetting();

        offset = 0;
        type = typeSelected;
        yearAgo = yearSelected;
        nameClassify = nameClassifySelected;
        if (!idSelected.contains(Constants.ID_ROW_ALL)) {
            id = idSelected;
            txvClassifyName.setText(nameClassify);
        } else {
            id = Constants.ID_ROW_ALL;
            txvClassifyName.setText(Constants.ROW_ALL);
        }

        // refresh list view
        lstView.invalidate();
        list.clear();
        loadList();
        loadIcon();

        // Adapter init
        // Set data adapter to list view
        ListViewProductAdapter adapter = new ListViewProductAdapter(this, list, rank);
        lstView.setAdapter(adapter);
    }

    /**
     * Function compute date after year selected in filter dialog
     *
     * @return String
     */
    public String dateFilter() {
        Calendar calendar = Calendar.getInstance();
        String day, month, year;
        if ((calendar.get(Calendar.DATE) - 1) < 10) {
            day = "0" + String.valueOf(calendar.get(Calendar.DATE) - 1);
        } else {
            day = String.valueOf(calendar.get(Calendar.DATE) - 1);
        }

        if (calendar.get(Calendar.MONTH) < 10) {
            month = "0" + String.valueOf(calendar.get(Calendar.MONTH));
        } else {
            month = String.valueOf(calendar.get(Calendar.MONTH));
        }
        year = String.valueOf(calendar.get(Calendar.YEAR) - yearAgo);
        if (yearAgo == Constants.SELECT_NULL) {
            return Constants.STRING_EMPTY;
        } else if (yearAgo == Constants.SELECT_ALL_YEAR) {
            return String.valueOf(Constants.SELECT_ALL_YEAR);
        } else {
            return year + month + day;
        }
    }

    /**
     * Result filter date selected in select dialog
     *
     * @param idSelected       {@link String}
     * @param typeSelected     {@link int}
     * @param dateFromSelected {@link String}
     * @param dateToSelected   {@link String}
     * @param dateChecked      {@link boolean}
     */
    @Override
    public void onSelectedDateDialog(String idSelected, int typeSelected, String dateFromSelected, String dateToSelected, boolean dateChecked) {

        //Enable click items setting
        enableItemsSetting();

        offset = 0;
        type = typeSelected;
        dateFrom = dateFromSelected;
        dateTo = dateToSelected;
        checkedDate = dateChecked;
        if (!idSelected.contains(Constants.ID_ROW_ALL)) {
            id = idSelected;
            txvClassifyName.setText(nameClassify);
        } else {
            id = Constants.ID_ROW_ALL;
            txvClassifyName.setText(Constants.ROW_ALL);
        }

        // refresh list view
        lstView.invalidate();
        list.clear();
        loadList();
        loadIcon();

        // Adapter init
        // Set data adapter to list view
        ListViewProductAdapter adapter = new ListViewProductAdapter(this, list, rank);
        lstView.setAdapter(adapter);
    }

    private void loadListItemsReturnBooks() {

        //Async task load list
        ProcessDialogLoadListReturnBooks processDialogLoadListReturnBooks = new ProcessDialogLoadListReturnBooks(this, this);
        processDialogLoadListReturnBooks.delegate = this;
        new ProcessDialogLoadListReturnBooks(this, processDialogLoadListReturnBooks.delegate).execute(offset, treeMapOrder, flagSettingNew, lstView);
    }

    @Override
    public void onLitSelectedFilterSettingDialog(FlagSettingNew _flagSettingNew,
                                                 FlagSettingOld _flagSettingOld,
                                                 Boolean _flagFilterSubmit) {
        //Enable click items setting
        enableItemsSetting();

        //If click submit
        if (_flagFilterSubmit) {
            flagSettingNew = _flagSettingNew;
            // refresh list view
            txvClassifyName.setText(_flagSettingNew.getFlagStockPercentShowScreen());
        }
        flagSettingOld = _flagSettingOld;
        flagFilterSubmit = _flagFilterSubmit;

        offset = 0;
        lstView.invalidate();
        list.clear();

        //TODO
        loadListItemsReturnBooks();

        imvIcon.setImageResource(R.drawable.book01);
    }

    /**
     * Save state list view
     *
     * @param rank {@link Integer}
     */
    private void saveStateList(int rank) {
        switch (rank) {
            case Constants.RANK_ARRIVAL:
                stateListArrival = lstView.onSaveInstanceState();
                break;
            case Constants.RANK_PLATFORM1:
                stateListPlatform1 = lstView.onSaveInstanceState();
                break;
            case Constants.RANK_PLATFORM2:
                stateListPlatform2 = lstView.onSaveInstanceState();
                break;
            case Constants.RANK_SURFACE:
                stateListSurface = lstView.onSaveInstanceState();
                break;
            case Constants.RANK_SHELDER:
                stateListShelder = lstView.onSaveInstanceState();
                break;
            case Constants.RANK_RETURN:
                stateListReturn = lstView.onSaveInstanceState();
                break;
            case Constants.RANK_PERIOD:
                stateListPeriod = lstView.onSaveInstanceState();
                break;
            case Constants.RANK_REGULAR:
                stateListRegular = lstView.onSaveInstanceState();
                break;
            default:
                break;
        }
    }

    /**
     * Restore state list
     *
     * @param rank {@link Integer}
     */
    private void restoreStateList(int rank) {

        switch (rank) {
            case Constants.RANK_ARRIVAL:
                if (stateListArrival != null) {
                    lstView.onRestoreInstanceState(stateListArrival);
                    stateListArrival = null;
                    break;
                }
                break;
            case Constants.RANK_PLATFORM1:
                if (stateListPlatform1 != null) {
                    lstView.onRestoreInstanceState(stateListPlatform1);
                    stateListPlatform1 = null;
                    break;
                }
                break;
            case Constants.RANK_PLATFORM2:
                if (stateListPlatform2 != null) {
                    lstView.onRestoreInstanceState(stateListPlatform2);
                    stateListPlatform2 = null;
                    break;
                }
                break;
            case Constants.RANK_SURFACE:
                if (stateListSurface != null) {
                    lstView.onRestoreInstanceState(stateListSurface);
                    stateListSurface = null;
                    break;
                }
                break;
            case Constants.RANK_SHELDER:
                if (stateListShelder != null) {
                    lstView.onRestoreInstanceState(stateListShelder);
                    stateListShelder = null;
                    break;
                }
                break;
            case Constants.RANK_RETURN:
                if (stateListReturn != null) {
                    lstView.onRestoreInstanceState(stateListReturn);
                    stateListReturn = null;
                    break;
                }
                break;
            case Constants.RANK_PERIOD:
                if (stateListPeriod != null) {
                    lstView.onRestoreInstanceState(stateListPeriod);
                    stateListPeriod = null;
                    break;
                }
                break;
            case Constants.RANK_REGULAR:
                if (stateListRegular != null) {
                    lstView.onRestoreInstanceState(stateListRegular);
                    stateListRegular = null;
                    break;
                }
                break;
            default:
                break;
        }
    }

    @Override
    protected void onResume() {
        super.onResume();

        //loadList();
        lstView.setOnItemClickListener(new AdapterView.OnItemClickListener() {
            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
                if (rank == Constants.RANK_RETURN) {
                    DProductDetailFragment dProductDetailFragment = new DProductDetailFragment();
                    FragmentManager fm = getSupportFragmentManager();
                    Bundle bundleItems = new Bundle();
                    bundleItems.putString(Constants.COLUMN_JAN_CD, list.get(position).get(Constants.COLUMN_JAN_CD));
                    bundleItems.putString(Constants.COLUMN_GOODS_NAME, list.get(position).get(Constants.COLUMN_GOODS_NAME));
                    bundleItems.putString(Constants.COLUMN_WRITER_NAME, list.get(position).get(Constants.COLUMN_WRITER_NAME));
                    bundleItems.putString(Constants.COLUMN_PUBLISHER_NAME_RETURN, list.get(position).get(Constants.COLUMN_PUBLISHER_NAME_RETURN));
                    bundleItems.putString(Constants.COLUMN_SALES_DATE, list.get(position).get(Constants.COLUMN_SALES_DATE));
                    bundleItems.putString(Constants.COLUMN_PRICE, list.get(position).get(Constants.COLUMN_PRICE));
                    bundleItems.putString(Constants.COLUMN_STOCK_COUNT, list.get(position).get(Constants.COLUMN_STOCK_COUNT));
                    bundleItems.putString(Constants.COLUMN_FIRST_SUPPLY_DATE, list.get(position).get(Constants.COLUMN_FIRST_SUPPLY_DATE));
                    bundleItems.putString(Constants.COLUMN_LAST_SUPPLY_DATE, list.get(position).get(Constants.COLUMN_LAST_SUPPLY_DATE));
                    bundleItems.putString(Constants.COLUMN_LAST_ORDER_DATE, list.get(position).get(Constants.COLUMN_LAST_ORDER_DATE));
                    bundleItems.putString(Constants.COLUMN_LAST_SALES_DATE, list.get(position).get(Constants.COLUMN_LAST_SALES_DATE));
                    bundleItems.putString(Constants.COLUMN_YEAR_RANK, list.get(position).get(Constants.COLUMN_YEAR_RANK));
                    bundleItems.putString(Constants.COLUMN_JOUBI, list.get(position).get(Constants.COLUMN_JOUBI));
                    bundleItems.putString(Constants.COLUMN_TOTAL_SALES, list.get(position).get(Constants.COLUMN_TOTAL_SALES));
                    bundleItems.putString(Constants.COLUMN_TOTAL_SUPPLY, list.get(position).get(Constants.COLUMN_TOTAL_SUPPLY));
                    bundleItems.putString(Constants.COLUMN_TOTAL_RETURN, list.get(position).get(Constants.COLUMN_TOTAL_RETURN));
                    bundleItems.putString(Constants.COLUMN_LOCATION_ID, formatCommon.formatLocationIdComma(list.get(position).get(Constants.COLUMN_LOCATION_ID)));
                    bundleItems.putString(Constants.COLUMN_MEDIA_GROUP1_NAME, list.get(position).get(Constants.COLUMN_MEDIA_GROUP1_NAME));
                    bundleItems.putString(Constants.COLUMN_MEDIA_GROUP2_NAME, list.get(position).get(Constants.COLUMN_MEDIA_GROUP2_NAME));
                    bundleItems.putInt(Constants.COLUMN_MAX_YEAR_RANK, maxYearRank.getMaxYearRank());
                    dProductDetailFragment.setArguments(bundleItems);
                    dProductDetailFragment.show(fm, null);
                }
            }
        });
    }

    //Get result process show list items
    @Override
    public void processFinish(ArrayList output, ListViewProductReturnBooksAdapter listViewProductReturnBooksAdapter) {
        list = output;
        adapterReturnBooks = listViewProductReturnBooksAdapter;
    }
}
