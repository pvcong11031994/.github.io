package com.fjn.magazinereturncandidate.api;

/**
 * API server configuration
 *
 * @author cong-pv
 * @since 2018-10-15
 */

public class Config {

    /**
     * API code login
     */
    public static final String CODE_LOGIN = "1";

    /**
     * API code get list data
     */
    public static final String CODE_GET_LIST_DATA = "2";

    /**
     * API code get max year rank
     */
    public static final String CODE_GET_MAX_YEAR_RANK = "7";


    /**
     * Http server
     */
    private static final String HTTP_SERVER = "https://bigdata.webpos-cloud.com:44312";

    /**
     * API Login
     */
    public static final String API_LOGIN = HTTP_SERVER + "/api/v1/login";
    public static final String API_GET_MAX_YEAR_RANK = HTTP_SERVER + "/api/v1/get_max_year_rank";

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
    static final String PROPERTY_VALUE = "application/json";

    /**
     * Charset UTF-8
     */
    static final String CHARSET_UTF_8 = "UTF-8";

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
     * Api key
     */
    static final String API_KEY = "api_key";

    /**
     * Api key value
     */
    static final String API_KEY_VALUE = "fjn_wsa_20180723_14";

    /**
     * Result list
     */
    public static final String RESULT_LIST = "resultList";

    /**
     * API Post file by user
     */
    public static final String API_POST_FILE = HTTP_SERVER + "/api/v1/post_log";

    /**
     * API Post file data return magazine
     */
    public static final String API_POST_FILE_DATA_RETURN_MAGAZINE = HTTP_SERVER + "/api/v1/post_files_data_return_magazine";

    /**
     * Property value post file
     */
    static final String PROPERTY_VALUE_POST_FILE = "multipart/form-data";

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
     * Boundary
     */
    public static final String BOUNDARY = "boundary";


    /**
     * Upload file
     */
    public static final String UPLOADFILE = "upload_file";

    public static final String UPLOADFILEDATA = "upload_file_data";

    /**
     * Content Disposition
     */
    public static final String CONTENT_DISPOSITION = "Content-Disposition: form-data; "
            + "name=\"" + UPLOADFILE + "\";filename=\"";

    public static final String FOLDER_ANDROID_SAVE_DATA = "/MagazineReturnCandidate/datasend/";
    /**
     * Content Disposition Send Data
     */
    public static final String CONTENT_DISPOSITION_SEND_DATA = "Content-Disposition: form-data; "
            + "name=\"" + UPLOADFILEDATA + "\";filename=\"";
}
