package com.android.returncandidate.common.constants;

/**
 * Constant
 *
 * @author tien-lv
 * @since 2017-12-01
 */

public class Constants {

    /**
     * Table name book
     */
    public static final String TABLE_BOOKS = "books";
    /**
     * Table max year rank
     */
    public static final String TABLE_MAX_YEAR_RANK = "maxyearrank";

    /**
     * Database version 1
     */
    public static final int DATABASE_VERSION = 1;

    /**
     * Database name
     */
    public static final String DATABASE_NAME = "Return_Candidate_Manager";

    /**
     * Table name return books
     */
    public static final String TABLE_RETURN_BOOKS = "returnbooks";

    /**
     * Table name users
     */
    public static final String TABLE_USER = "users";


    /**
     * Table name large classifications
     */
    public static final String TABLE_LARGE_CLASSIFICATIONS = "large_classifications";

    /**
     * Table name publishers
     */
    public static final String TABLE_PUBLISHERS = "publishers";

    /**
     * Column id
     */
    public static final String COLUMN_ID = "id";

    /**
     * Column name
     */
    public static final String COLUMN_NAME = "name";

    /**
     * Column user id
     */
    public static final String COLUMN_USER_ID = "userid";

    /**
     * Column password
     */
    public static final String COLUMN_PASSWORD = "password";

    /**
     * Column user role
     */
    public static final String COLUMN_USER_ROLE = "role";

    /**
     * Column Uid
     */
    public static final String COLUMN_UID = "uid";

    /**
     * Column shop id
     */
    public static final String COLUMN_SHOP_ID = "shop_id";

    /**
     * Column shop name
     */
    public static final String COLUMN_SHOP_NAME = "shop_name";

    /**
     * Column login key
     */
    public static final String COLUMN_LOGIN_KEY = "login_key";

    /**
     * Column license
     */
    public static final String COLUMN_LICENSE = "license";

    /**
     * Column infor list scan
     */
    public static final String COLUMN_INFOR_LIST_SCAN = "infor_scan";

    /**
     * Column server name
     */
    public static final String COLUMN_SERVER_NAME = "server_name";

    /**
     * Column location id
     */
    public static final String COLUMN_LOCATION_ID = "location_id";

    /**
     * Column jan code
     */
    public static final String COLUMN_JAN_CODE = "jan_code";

    /**
     * Column return date
     */
    public static final String COLUMN_RETURN_DATE = "return_date";

    /**
     * Column number
     */
    public static final String COLUMN_NUMBER = "number";

    /**
     * Column list status
     */
    public static final String COLUMN_LIST_STATUS = "list_status";

    /**
     * Alias count
     */
    public static final String ALIAS_COUNT = "count";

    /**
     * Query drop table is exist
     */
    public static final String QUERY_DROP_TABLE_EXIST = "DROP TABLE IF EXISTS %s";

    /**
     * TAG Application
     */
    public static final String TAG_APPLICATION_NAME = "返品候補";

    /**
     * Check Login
     */
    public static final String FLAG_LOGIN = "flag_login";

    /**
     * String empty
     */
    public static final String STRING_EMPTY = "";

    /**
     * Code barcode first 3 character
     */
    public static final String PREFIX_JAN_CODE_IGNORE_1 = "192";
    public static final String PREFIX_JAN_CODE_IGNORE_2 = "99";
    public static final String PREFIX_JAN_CODE_IGNORE_3 = "191";
    public static final String PREFIX_JAN_CODE_978 = "978";

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

    /**
     * Row Data All
     */
    public static final String ROW_ALL = "すべて";

    /**
     * Id Row Data All
     */
    public static final String ID_ROW_ALL = "-1";

    /**
     * List Column table return book
     */
    public static final String COLUMN_JAN_CD = "jan_cd";
    public static final String COLUMN_STOCK_COUNT = "bqsc_stock_count";
    public static final String COLUMN_GOODS_NAME = "bqgm_goods_name";
    public static final String COLUMN_WRITER_NAME = "bqgm_writer_name";
    public static final String COLUMN_PUBLISHER_CD = "bqgm_publisher_cd";
    public static final String COLUMN_PUBLISHER_NAME = "bqgm_publisher_name";
    public static final String COLUMN_COUNT_PUBLISHER_NAME = "count_bqgm_publisher_name";
    public static final String FLAG_SELECT = "flag_select";
    public static final String COLUMN_PRICE = "bqgm_price";
    /**
     * Column table max year rank
     */
    public static final String COLUMN_MAX_YEAR_RANK = "max_year_rank";
    /**
     * Column first supply date
     */
    public static final String COLUMN_FIRST_SUPPLY_DATE = "bqtse_first_supply_date";

    /**
     * Column last supply date
     */
    public static final String COLUMN_LAST_SUPPLY_DATE = "bqtse_last_supply_date";

    /**
     * Column last sales date
     */
    public static final String COLUMN_LAST_SALES_DATE = "bqtse_last_sale_date";

    /**
     * Column last order date
     */
    public static final String COLUMN_LAST_ORDER_DATE = "bqtse_last_order_date";
    public static final String COLUMN_MEDIA_GROUP1_CD = "bqct_media_group1_cd";
    public static final String COLUMN_MEDIA_GROUP1_NAME = "bqct_media_group1_name";
    public static final String COLUMN_MEDIA_GROUP2_CD = "bqct_media_group2_cd";
    public static final String COLUMN_MEDIA_GROUP2_NAME = "bqct_media_group2_name";
    public static final String COLUMN_SALES_DATE = "bqgm_sales_date";
    public static final String COLUMN_PERCENT = "percent";
    public static final String COLUMN_TRN_DATE = "bqio_trn_date";
    public static final String COLUMN_FLAG_SALES = "flag_sales";
    public static final String COLUMN_YEAR_RANK = "year_rank";
    public static final String COLUMN_JOUBI = "joubi";
    public static final String COLUMN_TOTAL_SALES = "sts_total_sales";
    public static final String COLUMN_TOTAL_SUPPLY = "sts_total_supply";
    public static final String COLUMN_TOTAL_RETURN = "sts_total_return";
    public static final String COLUMN_ALIAS_COUNT = "cnt";

    public static final int VALUE_COUNT_COLUMN_TABLE_RETURN_BOOK_INSERT = 25;

    /**
     * Symbol ￥
     */
    public static final String SYMBOL = "￥";

    /**
     * String Null for empty data
     */
    public static final String NULL = "null";
    /**
     * String Null for empty data
     */
    public static final String BLANK = "";


    public static final String FLAG_CLASSIFICATION_GROUP1_CD = "flag_classification_group1_cd";
    public static final String FLAG_CLASSIFICATION_GROUP1_NAME = "flag_classification_group1_name";
    public static final String FLAG_CLASSIFICATION_GROUP2_CD = "flag_classification_group2_cd";
    public static final String FLAG_CLASSIFICATION_GROUP2_NAME = "flag_classification_group2_name";
    public static final String FLAG_PUBLISHER = "flag_publisher";
    public static final String FLAG_PUBLISHER_SHOW_SCREEN = "flag_publisher_show_screen";
    public static final String FLAG_RELEASE_DATE = "flag_release_date";
    public static final String FLAG_RELEASE_DATE_SHOW_SCREEN = "flag_release_date_show_screen";
    public static final String FLAG_UNDISTURBED = "flag_undisturbed";
    public static final String FLAG_UNDISTURBED_SHOW_SCREEN = "flag_undisturbed_show_screen";
    public static final String FLAG_NUMBER_OF_STOCKS = "flag_number_of_stocks";
    public static final String FLAG_NUMBER_OF_STOCKS_SHOW_SCREEN = "flag_number_of_stocks_show_screen";
    public static final String FLAG_STOCKS_PERCENT = "flag_stocks_percent";
    public static final String FLAG_STOCKS_PERCENT_SHOW_SCREEN = "flag_stocks_percent_show_screen";
    public static final String FLAG_JOUBI = "flag_joubi";
    public static final String FLAG_CLICK_SETTING = "flag_click_setting";
    public static final String FLAG_GROUP1_CD = "flag_group1_cd";
    public static final String FLAG_GROUP1_NAME = "flag_group1_name";


    public static final String FLAG_CLASSIFICATION_GROUP1_CD_OLD = "flag_classification_group1_cd_old";
    public static final String FLAG_CLASSIFICATION_GROUP1_NAME_OLD = "flag_classification_group1_name_old";
    public static final String FLAG_CLASSIFICATION_GROUP2_CD_OLD = "flag_classification_group2_cd_old";
    public static final String FLAG_CLASSIFICATION_GROUP2_NAME_OLD = "flag_classification_group2_name_old";
    public static final String FLAG_PUBLISHER_OLD = "flag_publisher_old";
    public static final String FLAG_PUBLISHER_SHOW_SCREEN_OLD = "flag_publisher_show_screen_old";
    public static final String FLAG_RELEASE_DATE_OLD = "flag_release_date_old";
    public static final String FLAG_RELEASE_DATE_SHOW_SCREEN_OLD = "flag_release_date_show_screen_old";
    public static final String FLAG_UNDISTURBED_OLD = "flag_undisturbed_old";
    public static final String FLAG_UNDISTURBED_SHOW_SCREEN_OLD = "flag_undisturbed_show_screen_old";
    public static final String FLAG_NUMBER_OF_STOCKS_OLD = "flag_number_of_stocks_old";
    public static final String FLAG_NUMBER_OF_STOCKS_SHOW_SCREEN_OLD = "flag_number_of_stocks_show_screen_old";
    public static final String FLAG_STOCKS_PERCENT_OLD = "flag_stocks_percent_old";
    public static final String FLAG_STOCKS_PERCENT_SHOW_SCREEN_OLD = "flag_stocks_percent_show_screen_old";
    public static final String FLAG_JOUBI_OLD = "flag_joubi_old";


    public static final String SELECT_POISITION = "position";
    public static final String SELECT_VALUE = "value";

    public static final String HEADER_CLASSIFICATION_1 = "大分類選択";
    public static final String HEADER_CLASSIFICATION_2 = "中分類選択";
    public static final String HEADER_PUBLISHER = "出版社選択";
    public static final String HEADER_RELEASE_DATE = "発売日選択";
    public static final String HEADER_UNDISTURBED = "未勤期間選択";
    public static final String HEADER_NUMBER_OF_STOCKS = "在庫数選択";
    public static final String HEADER_STOCKS_PERCENT = "在庫％選択";

    public static final String FLAG_DEFAULT = "1";
    public static final String FLAG_DEFAULT_RELEASE_DATE_SHOW = "１ヶ月以上前";
    public static final String FLAG_DEFAULT_UNDISTURBED_SHOW = "１ヶ月以上";
    public static final String FLAG_DEFAULT_NUMBER_OF_STOCKS_SHOW = "１以上";
    public static final String FLAG_DEFAULT_STOCKS_PERCENT_SHOW = "５％";

    public static final String VALUE_SUM_STOCKS = "冊";
    public static final String VALUE_COUNT_JAN_CD = "商品";

    public static final String VALUE_DEFAULT_DATE = "10000101";

    public static final int VALUE_DEFAULT_DATE_INT = 10000101;

    public static final String VALUE_CHECK_ONCLICK_SETTING = "0";
    public static final String FLAG_SWITCH_OCR = "flag_switch_ocr";

    public static final String FLAG_CLASSIFICATION_GROUP1_CD_BACK = "flag_classification_group1_cd_back";
    public static final String FLAG_CLASSIFICATION_GROUP1_NAME_BACK = "flag_classification_group1_name_back";

    public static final String YES_STANDING = "常備";
    public static final String NO_STANDING = "なし";
    public static final String VALUE_YES_STANDING = "1";
    public static final String VALUE_NO_STANDING = "0";
    public static final String VALUE_JOUBI = "5";
    public static final String VALUE_JOUBI_SHOW = "常備";

    public static int VALUE_WHEN_SELECT_ONE = 1;

    public static final String VALUE_STR_WHEN_SELECT_MULTI = "複数選択";
    public static final String VALUE_STR_WHEN_SELECT_ALL = "すべて";

    public static final String VALUE_STR_CHECK = "1";
    public static final String VALUE_STR_NO_CHECK = "0";

    public static final int VALUE_INT_DEFAULT_SELECT_ALL = 0;

    public static final String VALUE_MAX_YEAR_RANK = "99999999"; // if year_rank = null
    public static final String SHOW_MAX_YEAR_RANK = "売上なし"; // if year_rank = null

}
