package com.fjn.magazinereturncandidate.common.constants;

/**
 * Constant
 *
 * @author cong-pv
 * @since 2018-10-15
 */

public class Constants {

    /**
     * GROUP CONSTANTS NAME
     */
    public static final String TAG_APPLICATION_NAME = "雑誌返品チェック";
    //Database name
    public static final String DATABASE_NAME = "MAGAZINE_RETURN_CANDIDATE";
    //Database version
    public static final int DATABASE_VERSION = 1;


    /**
     * GROUP TABLE USER
     */
    public static final String TABLE_USER = "users";
    public static final String COLUMN_USER_ID = "userid";
    public static final String COLUMN_PASSWORD = "password";
    public static final String COLUMN_NAME = "name";
    public static final String COLUMN_UID = "uid";
    public static final String COLUMN_SHOP_ID = "shop_id";
    public static final String COLUMN_SHOP_NAME = "shop_name";
    public static final String COLUMN_LOGIN_KEY = "login_key";
    public static final String COLUMN_SERVER_NAME = "server_name";
    public static final String COLUMN_USER_ROLE = "role";
    public static final String COLUMN_LICENSE = "license";
    public static final String FLAG_LOGIN = "flag_login";


    /**
     * GROUP TABLE RETURN MAGAZINE
     */
    public static final String TABLE_RETURN_MAGAZINE = "return_magazine";
    public static final String COLUMN_JAN_CD = "jan_cd";
    public static final String COLUMN_STOCK_COUNT = "bqsc_stock_count";
    public static final String COLUMN_GOODS_NAME = "bqgm_goods_name";
    public static final String COLUMN_WRITER_NAME = "bqgm_writer_name";
    public static final String COLUMN_PUBLISHER_CD = "bqgm_publisher_cd";
    public static final String COLUMN_PUBLISHER_NAME = "bqgm_publisher_name";
    public static final String COLUMN_PRICE = "bqgm_price";
    public static final String COLUMN_FIRST_SUPPLY_DATE = "bqtse_first_supply_date";
    public static final String COLUMN_LAST_SUPPLY_DATE = "bqtse_last_supply_date";
    public static final String COLUMN_LAST_SALES_DATE = "bqtse_last_sale_date";
    public static final String COLUMN_LAST_ORDER_DATE = "bqtse_last_order_date";
    public static final String COLUMN_MEDIA_GROUP1_CD = "bqct_media_group1_cd";
    public static final String COLUMN_MEDIA_GROUP1_NAME = "bqct_media_group1_name";
    public static final String COLUMN_MEDIA_GROUP2_CD = "bqct_media_group2_cd";
    public static final String COLUMN_MEDIA_GROUP2_NAME = "bqct_media_group2_name";
    public static final String COLUMN_SALES_DATE = "bqgm_sales_date";
    public static final String COLUMN_TRN_DATE = "bqio_trn_date";
    public static final String COLUMN_PERCENT = "percent";
    public static final String COLUMN_FLAG_SALES = "flag_sales";
    public static final String COLUMN_YEAR_RANK = "year_rank";
    public static final String COLUMN_JOUBI = "joubi";
    public static final String COLUMN_TOTAL_SALES = "sts_total_sales";
    public static final String COLUMN_TOTAL_SUPPLY = "sts_total_supply";
    public static final String COLUMN_TOTAL_RETURN = "sts_total_return";
    public static final String COLUMN_LOCATION_ID = "location_id";
    public static final int VALUE_COUNT_COLUMN_TABLE_RETURN_MAGAZINE_INSERT = 22;

    /**
     * GROUP TABLE MAX YEAR RANK
     */
    public static final String TABLE_MAX_YEAR_RANK = "maxyearrank";
    public static final String COLUMN_MAX_YEAR_RANK = "max_year_rank";


    /**
     * GROUP JANCODE IGNORE
     */
    public static final String PREFIX_JAN_CODE_IGNORE_1 = "192";
    public static final String PREFIX_JAN_CODE_IGNORE_2 = "99";
    public static final String PREFIX_JAN_CODE_IGNORE_3 = "191";
    public static final String PREFIX_JAN_CODE_978 = "978";
    public static final String PREFIX_JAN_CODE_MAGAZINE = "4910";
    public static final int JAN_12_CHAR = 12;
    public static final int JAN_13_CHAR = 13;
    public static final int JAN_18_CHAR = 18;


    /**
     * QUERY
     */
    public static final String QUERY_DROP_TABLE_EXIST = "DROP TABLE IF EXISTS %s";


    /**
     * String empty
     */
    public static final String STRING_EMPTY = "";
    /**
     * Symbol ￥
     */
    public static final String SYMBOL = "￥";

    /**
     * String Null for empty data
     */
    public static final String BLANK = "";
    /**
     * Alias count
     */
    public static final String ALIAS_COUNT = "count";

    public static final String VALUE_DEFAULT_DATE = "10000101";

    /**
     * Flag 0
     */
    public static final String FLAG_0 = "0";

    /**
     * Flag 1
     */
    public static final String FLAG_1 = "1";

    /**
     * Time out
     */
    public static final long TIME_OUT = 5 * 60 * 1000;

    public static final String VALUE_MAX_YEAR_RANK = "99999999"; // if year_rank = null
    public static final String SHOW_MAX_YEAR_RANK = "売上なし"; // if year_rank = null

    public static final String VALUE_JOUBI = "5";
    public static final String VALUE_JOUBI_SHOW = "常備";

    public static final String POSITION_EDIT_PRODUCT = "position_edit";

    /**
     * Column infor list scan
     */
    public static final String COLUMN_INFOR_LIST_SCAN = "infor_scan";

    public static final String FLAG_SWITCH_OCR = "flag_switch_ocr";

    public static final int MAXIMUM_INPUT_NUMBER_JAN_CODE = 999;


}


