package com.android.returncandidate.api;

/**
 * API server configuration
 *
 * @author tien-lv
 * @since 2017-12-22
 */

public class Config {

    /**
     * Method POST
     */
    static final String METHOD_POST = "POST";

    /**
     * Property key
     */
    static final String PROPERTY_KEY = "Content-Type";

    /**
     * Property value
     */
    //static final String PROPERTY_VALUE = "application/x-www-form-urlencoded";
    static final String PROPERTY_VALUE = "application/json";

    /**
     * Property value post file
     */
    static final String PROPERTY_VALUE_POST_FILE = "multipart/form-data";

    /**
     * Api key
     */
    static final String API_KEY = "api_key";

    /**
     * Api key value
     */
    static final String API_KEY_VALUE = "fjn_wsa_20180723_14";

    /**
     * Charset UTF-8
     */
    static final String CHARSET_UTF_8 = "UTF-8";

    /**
     * Http server
     */
    private static final String HTTP_SERVER = "https://bigdata.webpos-cloud.com:44312";

    /**
     * API Login
     */
    public static final String API_LOGIN = HTTP_SERVER + "/api/v1/login";

    /**
     * API Get list data by user
     */
    public static final String API_GET_LIST_BY_USER = HTTP_SERVER + "/api/v1/get_list_bq_by_user_ochanomizu_store_02";
    public static final String API_GET_MAX_YEAR_RANK = HTTP_SERVER + "/api/v1/get_max_year_rank";
    public static final String API_GET_LIST_PUBLISHER = HTTP_SERVER + "/api/v1/get_list_publisher";

    /**
     * API Post file by user
     */
    public static final String API_POST_FILE = HTTP_SERVER + "/api/v1/post_log";

    /**
     * API Get list data classify
     */
    public static final String API_GET_LIST_CLASSIFY =
            HTTP_SERVER + "/api/v1/get_list_classify";

    /**
     * API code login
     */
    public static final String CODE_LOGIN = "1";

    /**
     * API code get list data
     */
    public static final String CODE_GET_LIST_DATA = "2";

    /**
     * API code get list data by user
     */
    public static final String CODE_GET_LIST_BY_USER = "3";


    /**
     * API code get list data shop by user
     */
    public static final String CODE_GET_LIST_CLASSIFY = "5";

    /**
     * API code get list data publisher
     */
    public static final String CODE_GET_LIST_PUBLISHER = "6";


    /**
     * API code get max year rank
     */
    public static final String CODE_GET_MAX_YEAR_RANK = "7";

    /**
     * Login key
     */
    public static final String LOGIN_KEY = "login_key";

    /**
     * maker
     * Type
     */
    public static final String TYPE = "type";

    /**
     * Type location
     */
    public static final int TYPE_LOCATION = 1;

    /**
     * Type classify
     */
    public static final int TYPE_CLASSIFY = 2;

    /**
     * Type publisher
     */
    public static final int TYPE_PUBLISHER = 3;

    /**
     * Result list
     */
    public static final String RESULT_LIST = "resultList";

    /**
     * Boundary
     */
    public static final String BOUNDARY = "boundary";

    /**
     * Connection key
     */
    public static final String CONNECTION_KEY = "Connection";

    /**
     * Connection value
     */
    public static final String CONNECTION_VALUE = "Keep-Alive";

    /**
     * ENCTYPE key
     */
    public static final String ENCTYPE_KEY = "ENCTYPE";

    /**
     * Upload file
     */
    public static final String UPLOADFILE = "upload_file";

    /**
     * Content Disposition
     */
    public static final String CONTENT_DISPOSITION = "Content-Disposition: form-data; "
            + "name=\"" + UPLOADFILE + "\";filename=\"";
}
