package com.fjn.magazinereturncandidate.common.constants;

/**
 * Message and error code
 *
 * @author cong-pv
 * @since 2018-10-15
 */

public class Message {

    /**
     * Code 200
     */
    public static final String CODE_200 = "200";

    /**
     * Code 401
     */
    public static final String CODE_401 = "401";

    /**
     * Message 401
     */
    public static final String MESSAGE_401 = "ユーザとパスワードが不正です。";

    /**
     * Code 404
     */
    public static final String CODE_404 = "404";

    /**
     * Message 404
     */
    public static final String MESSAGE_404 = "サービスが見つかりません。";

    /**
     * Code 500
     */
    public static final String CODE_500 = "500";

    /**
     * Message 500
     */
    public static final String MESSAGE_500 = "インターナルサーバエラー。";

    /**
     * Message result empty
     */
    public static final String MESSAGE_RESULT_EMPTY = "戻りデータがありません。";

    /**
     * Login Activity tag
     */
    public static final String TAG_LOGIN_ACTIVITY = "ログイン画面：";

    /**
     * Message TAG Activity start
     */
    public static final String MESSAGE_ACTIVITY_START = "処理開始。";


    /**
     * Message when loading data
     */
    public static final String MESSAGE_LOADING_SCREEN = "ロード中。。。";

    /**
     * Message loading data from server
     */
    public static final String LOADING_DATA_FROM_SERVER = "サーバからのロード中。";


    /**
     * Check input empty
     */
    public static final String MESSAGE_CHECK_INPUT_EMPTY = "%sを入力してください。";

    /**
     * Message Activity end
     */
    public static final String MESSAGE_ACTIVITY_END = "処理完了。";

    /**
     * Message complete loading data from server
     */
    public static final String LOADING_DATA_FROM_SERVER_SUCCESS = "サーバからのロード完了。";

    /**
     * Message use when login success
     */
    public static final String MESSAGE_LOGIN_SUCCESS = "ロードに成功しました。（ユーザ：%s、店舗コード：%s、サーバー名：%s）";

    /**
     * Message TAG Activity move
     */
    public static final String MESSAGE_ACTIVITY_MOVE = "%sから%sに遷移する。";

    /**
     * Login Activity Name
     */
    public static final String LOGIN_ACTIVITY_NAME = "ログイン画面";

    /**
     * SCANNER ACTIVITY Name
     */
    public static final String SCANNER_ACTIVITY_NAME = "バーコードスキャン用カメラ画面";


    /**
     * Message when import data
     */
    public static final String MESSAGE_IMPORT_DATA_SCREEN = "データ取り込み中ー。。。";

    /**
     * Message connect to network
     */
    public static final String MESSAGE_NETWORK_ERR = "ネットワークに\n接続できませんでした。";

    /**
     * SCANNER ACTIVITY tag
     */
    public static final String TAG_SCANNER_ACTIVITY = "バーコードスキャン用カメラ画面：";

    /**
     * Message reload
     */
    public static final String MESSAGE_RELOAD = "\nリロードを行います。よろしいですか？";

    /**
     * Message select  Yes
     */
    public static final String MESSAGE_SELECT_YES = "はい";

    /**
     * Message select No
     */
    public static final String MESSAGE_SELECT_NO = "いいえ";

    /**
     * UNLOCK ACTIVITY tag
     */
    public static final String TAG_UNLOCK_ACTIVITY = "ロック解除画面：";

    /**
     * Message when upload log to GCS
     */
    public static final String MESSAGE_UPLOAD_LOG_SCREEN = "ログファイルをアップロードしています。。。";

    /**
     * Message wrong password
     */
    public static final String MESSAGE_PASSWORD_ERR = "パスワードが不正の場合。";

    /**
     * UNLOCK ACTIVITY Name
     */
    public static final String UNLOCK_ACTIVITY_NAME = "ロック解除画面";

    /**
     * Message data count summary
     */
    public static final String MESSAGE_LOADING_DATA_NUMBER = "データが%s件ロードできました。";

    /**
     * Show message when click button logout
     */
    public static final String MESSAGE_LOGOUT = "ログアウト。";

    /**
     * TAG message for remove latest scanned object
     */
    public static final String TAG_SCANNER_ACTIVITY_CANCEL = "バーコードスキャン用カメラ画面＿取消：%s";

    /**
     * TAG message for edit data scan
     */
    public static final String TAG_SCANNER_ACTIVITY_EDIT = "バーコードスキャン用カメラ画面＿編集：%s";

    /**
     * Message when download data
     */
    public static final String MESSAGE_DOWNLOAD_DATA_SCREEN = "データダウンロード中ー。。。";

    /**
     * Message load time
     */
    public static final String MESSAGE_TIME_EXECUTE = "実行時間：%s秒（件数：%s）。";

    /**
     * TAG message for product out list return candidate
     */
    public static final String TAG_SCANNER_ACTIVITY_OUT_LIST = "バーコードスキャン用カメラ画面＿リスト外：%s";

    /**
     * TAG message for product in list return candidate
     */
    public static final String TAG_SCANNER_ACTIVITY_INLIST = "バーコードスキャン用カメラ画面＿リスト内：%s　%s";

    public static final String MESSAGE_CONFIRM_DELETE_JAN_CD = "商品JANコード%sを削除しますか。";

    public static final String MESSAGE_SEND_DATA = "データを送信します。";

    public static final String CONFIRM_OK = "ＯＫ";

    public static final String MESSAGE_ERROR_INPUT_JANCODE = "JANコードは12桁、13桁、18桁のいずれかで入力してください。";

    public static final String MESSAGE_ERROR_CHECK_DIGIT_JANCODE = "JANコードが不正です。再度入力してください。";

    public static final String MESSAGE_INFO_PATH_CSV = "ネットワーク接続できません。\nCSVファイルが /MagazineReturnCandidate/datasend フォルダに保存されています。";

    public static final String MESSAGE_ERROR_BLANK_INPUT = "数量を入力してください。";

    public static final String MESSAGE_SEND_LIST_FILE_CSV = "/MagazineReturnCandidate/datasendフォルダには未送信のCSVファイルが存在しています。\\n送信を行います。よろしいですか？";

    public static final String MESSAGE_YES = "ＯＫ";

    public static final String MESSAGE_NO = "キャンセル";

    public static final String MESSAGE_CONFIRM_SEND_DATA = "データを送信します。\n送信完了時にデータはクリアされます。";

    public static final String MESSAGE_CONFIRM_LOGOUT = "ログアウトを行います。\nよろしいですか？";

    public static final String MESSAGE_YES_JP = "はい";

    public static final String MESSAGE_NO_JP = "いいえ";

    public static final String MESSAGE_RETRY = "リトライ";

    public static final String MESSAGE_CANCEL = "キャンセル";




}
