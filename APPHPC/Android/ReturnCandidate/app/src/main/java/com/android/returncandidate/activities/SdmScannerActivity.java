package com.android.returncandidate.activities;

import android.annotation.SuppressLint;
import android.app.ProgressDialog;
import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;
import android.media.MediaPlayer;
import android.net.ConnectivityManager;
import android.net.NetworkInfo;
import android.os.Bundle;
import android.os.Environment;
import android.support.v4.app.FragmentManager;
import android.support.v7.app.AlertDialog;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;
import android.view.Gravity;
import android.view.KeyEvent;
import android.view.View;
import android.widget.Button;
import android.widget.CompoundButton;
import android.widget.ImageButton;
import android.widget.ListView;
import android.widget.Switch;
import android.widget.TextView;
import android.widget.Toast;
import android.widget.AdapterView;

import com.android.returncandidate.R;
import com.android.returncandidate.adapters.ListViewScanAdapter;
import com.android.returncandidate.api.Config;
import com.android.returncandidate.api.HttpPostFile;
import com.android.returncandidate.api.HttpResponse;
import com.android.returncandidate.common.constants.Constants;
import com.android.returncandidate.common.constants.Message;
import com.android.returncandidate.common.helpers.DatabaseHelper;
import com.android.returncandidate.common.helpers.Log4JHelper;
import com.android.returncandidate.common.utils.Common;
import com.android.returncandidate.common.utils.DatabaseManager;
import com.android.returncandidate.common.utils.FlagSettingNew;
import com.android.returncandidate.common.utils.FlagSettingOld;
import com.android.returncandidate.common.utils.GzipFile;
import com.android.returncandidate.common.utils.LogManager;
import com.android.returncandidate.db.entity.Books;
import com.android.returncandidate.db.entity.CLP;
import com.android.returncandidate.db.entity.MaxYearRank;
import com.android.returncandidate.db.entity.Users;
import com.android.returncandidate.db.models.BookModel;
import com.android.returncandidate.db.models.MaxYearRankModel;
import com.android.returncandidate.db.models.UserModel;

import com.android.returncandidate.fragments.DProductDetailFragment;
import com.android.returncandidate.fragments.DFilterSettingFragment;
import com.honeywell.barcode.HSMDecodeComponent;
import com.honeywell.barcode.HSMDecodeResult;
import com.honeywell.barcode.HSMDecoder;
import com.honeywell.barcode.OCRActiveTemplate;
import com.honeywell.barcode.Symbology;
import com.honeywell.barcode.WindowMode;
import com.honeywell.license.ActivationManager;
import com.honeywell.license.ActivationResult;
import com.honeywell.plugins.decode.DecodeResultListener;

import java.io.File;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.HashMap;
import java.util.LinkedList;
import java.util.List;

/**
 * <h1>Main barcode scanner activity.</h1>
 * Begin scanning barcode, if scanned barcode is matched with database will result with a beep
 * sound, otherwise will result with a buzz sound. User can remove recently scanned object via 取消
 * button
 * and logout via ログアウト button
 *
 * @author minh-th
 * @version 2.0
 * @since 2018-05-10
 */
public class SdmScannerActivity extends AppCompatActivity implements View.OnClickListener,
        HttpResponse, DecodeResultListener,
        DFilterSettingFragment.ItemSelectedFilterSettingDialogListener {

    private String TAG = Constants.TAG_APPLICATION_NAME;

    private String userID;
    private String shopID;
    private String serverName;
    private String license;
    private LinkedList<String[]> arrBookInlist;
    private HashMap<String, LinkedList<String[]>> hashMapArrBook;

    private int countFile = 0;
    private File[] files;
    private ProgressDialog progress;

    private HSMDecoder hsmDecoder;
    private HSMDecodeComponent decCom;
    private Button btnLogout;
    private Button btnCancel;
    private ImageButton imbFilter;
    private ListView lvBook;
    private MediaPlayer normalSound;
    private MediaPlayer noReturnSound;

    private BookModel mBookModel;
    private MaxYearRankModel maxYearRankModel;
    private FlagSettingNew flagSettingNew;
    private FlagSettingOld flagSettingOld;
    private Boolean flagFilterSubmit;
    private Common common;
    long timeout;

    private Switch aSwitchOCR;
    private String flagSwitchOCR;
    private Boolean isEnableScan = true;
    Boolean joubi;
    MaxYearRank maxYearRank = new MaxYearRank();

    /**
     * Initialize screen layout
     *
     * @param state {@link Bundle}
     */
    @Override
    public void onCreate(Bundle state) {

        LogManager.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.MESSAGE_ACTIVITY_START);
        super.onCreate(state);

        setContentView(R.layout.activity_sdm_scanner);

        // Load user info

        Bundle bundle = getIntent().getExtras();
        if (bundle != null) {
            userID = bundle.getString(Constants.COLUMN_USER_ID);
            shopID = bundle.getString(Constants.COLUMN_SHOP_ID);
            serverName = bundle.getString(Constants.COLUMN_SERVER_NAME);
            license = bundle.getString(Constants.COLUMN_LICENSE);
            hashMapArrBook = (HashMap<String, LinkedList<String[]>>) bundle.getSerializable(Constants.COLUMN_INFOR_LIST_SCAN);
        }
        Log.d("LINCONGPVScan", license);
        //Set default start app
        flagSettingNew = new FlagSettingNew();
        flagSettingOld = new FlagSettingOld();
        common = new Common();

        // UI init
        Users user = new UserModel().getUserInfo();
        TextView tvUserName = (TextView) findViewById(R.id.txv_user_name);
        tvUserName.setText(user.getName());

        lvBook = (ListView) findViewById(R.id.list_book);
        arrBookInlist = new LinkedList<>();
        if (hashMapArrBook != null) {
            arrBookInlist = new LinkedList<>(hashMapArrBook.get(Constants.COLUMN_INFOR_LIST_SCAN));

        }

        mBookModel = new BookModel();
        maxYearRankModel = new MaxYearRankModel();
        DatabaseManager.initializeInstance(
                new DatabaseHelper(getApplicationContext()));

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

        maxYearRank = maxYearRankModel.getMaxYearRank();

        //Get list default group1 cd
        List<CLP> listDefaultGroup1 = mBookModel.getInfoGroupCd1();
        List<CLP> listDefaultGroup2 = mBookModel.getInfoGroupCd2(Constants.ID_ROW_ALL);
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

        //Flag disable OCR
        flagSwitchOCR = Constants.FLAG_0;

        flagFilterSubmit = true;
        // Init process loading screen
        progress = new ProgressDialog(this);
        progress.setMessage(Message.MESSAGE_UPLPOAD_LOG_SCREEN);
        progress.setCancelable(false);

        // Activate HSM license
        Log.d("LINCONGPVScanACTIVE1", license);
        ActivationResult activationResult = ActivationManager.activate(this,
                license);

        Toast.makeText(this, "Activation result: " + activationResult, Toast.LENGTH_LONG).show();

        decCom = (HSMDecodeComponent) findViewById(R.id.hsm_decodeComponent);

        // HSM init
        hsmDecoder = HSMDecoder.getInstance(this);

        // Declare symbology
        hsmDecoder.enableSymbology(Symbology.EAN13);

        // Declare HSM component UI
        hsmDecoder.enableFlashOnDecode(false);
        hsmDecoder.enableSound(false);
        hsmDecoder.enableAimer(false);
        hsmDecoder.setWindowMode(WindowMode.CENTERING);
        hsmDecoder.setWindow(18, 42, 0, 100);

        // Assign listener
        hsmDecoder.addResultListener(this);

        // Sound init
        normalSound = MediaPlayer.create(this,
                R.raw.pingpong_main); // sound is inside res/raw/mysound
        noReturnSound = MediaPlayer.create(this,
                R.raw.wrong_main); // sound is inside res/raw/mysound

        // Button init
        btnLogout = (Button) findViewById(R.id.btn_logout);
        btnCancel = (Button) findViewById(R.id.btn_cancel);
        imbFilter = (ImageButton) findViewById(R.id.imb_filter);
        aSwitchOCR = (Switch) findViewById(R.id.switch_OCR);
        btnLogout.setOnClickListener(this);
        btnCancel.setOnClickListener(this);
        imbFilter.setOnClickListener(this);

        aSwitchOCR.setOnCheckedChangeListener(new CompoundButton.OnCheckedChangeListener() {
            @Override
            public void onCheckedChanged(CompoundButton buttonView, boolean isChecked) {
                if (aSwitchOCR.isChecked()) {
                    flagSwitchOCR = Constants.FLAG_1;
                    hsmDecoder.enableSymbology(Symbology.OCR);
                    //hsmDecoder.enableSymbology(Symbology.EAN13_ISBN);
                    hsmDecoder.setOCRActiveTemplate(OCRActiveTemplate.ISBN);
                    hsmDecoder.disableSymbology(Symbology.EAN13);
                } else {
                    flagSwitchOCR = Constants.FLAG_0;
                    hsmDecoder.enableSymbology(Symbology.EAN13);
                    //hsmDecoder.setOCRActiveTemplate(OCRActiveTemplate.ISBN);
                    hsmDecoder.disableSymbology(Symbology.OCR);

                }
            }
        });

        //Check if arr list books not null
        if (arrBookInlist != null) {
            // Set data adapter to list view
            ListViewScanAdapter adapterBook = new ListViewScanAdapter(this, arrBookInlist);
            lvBook.setAdapter(adapterBook);
        }
    }

    /**
     * Decoded results handler
     *
     * @param hsmDecodeResults {@link HSMDecodeResult}
     */
    @Override
    public void onHSMDecodeResult(HSMDecodeResult[] hsmDecodeResults) {

        if (hsmDecodeResults.length > 0) {
            for (int i = 0; i < hsmDecodeResults.length; i++) {
                barcodeValidation(hsmDecodeResults[i].getBarcodeData());
            }
        }
    }

    /**
     * Barcode validation
     *
     * @param barcode {@link String}
     */
    private void barcodeValidation(String barcode) {

        // Product code validation
        //Check prefix jan code ignore
        if (barcode.startsWith(Constants.PREFIX_JAN_CODE_IGNORE_1) ||
                barcode.startsWith(Constants.PREFIX_JAN_CODE_IGNORE_2) ||
                barcode.startsWith(Constants.PREFIX_JAN_CODE_IGNORE_3)) {
            return;
        }
        if (aSwitchOCR.isChecked()) {
            barcode = formatOCRJAN(barcode);
        }

        Books book;
        book = mBookModel.getItemBook(barcode);

        String strPrice = Constants.BLANK;
        if (book.getBqgm_price() != null) {
            strPrice = book.getBqgm_price().toString();
        }


        String[] item = new String[]{barcode, String.valueOf(book.getBqsc_stock_count()), book.getBqgm_goods_name(),
                book.getBqgm_writer_name(), book.getBqgm_publisher_name(), strPrice,
                book.getBqtse_first_supply_date(), book.getBqtse_last_supply_date(),
                book.getBqtse_last_sale_date(), book.getBqtse_last_order_date(),
                book.getBqct_media_group1_cd(), book.getBqct_media_group1_name(),
                book.getBqct_media_group2_cd(), book.getBqct_media_group2_name(),
                book.getBqgm_sales_date(), book.getBqgm_publisher_cd(), book.getBqio_trn_date(),
                Constants.FLAG_0, String.valueOf(book.getYear_rank()), String.valueOf(book.getJoubi()),
                String.valueOf(book.getSts_total_sales()), String.valueOf(book.getSts_total_supply()),
                String.valueOf(book.getSts_total_return()), book.getLocation_id()};
        checkIndexList(barcode);

        int intSalesDate = common.ConvertStringDateToInt(common.FormatDateTime(flagSettingNew.getFlagReleaseDate()));
        int intUndisturbed = common.ConvertStringDateToInt(common.FormatDateTime(flagSettingNew.getFlagUndisturbed()));
        int intNumberOfStock = Integer.parseInt(flagSettingNew.getFlagNumberOfStocks());
        float intStockPercent = Float.parseFloat(common.FormatPercent(flagSettingNew.getFlagStockPercent()));

        if (!barcode.equals(book.getJan_cd())) {
            noReturnSound.start();
            LogManager.i(TAG, String.format(Message.TAG_SCANNER_ACTIVITY_OUT_LIST, barcode));
        } else {
            //If select all group1cd
            if ((Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup1Cd().get(0))) &&
                    (Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagPublisher().get(0)))) {
                if (common.ConvertStringDateToInt(book.getBqgm_sales_date()) < intSalesDate
                        && common.ConvertStringDateToInt(book.getBqio_trn_date()) <= intUndisturbed
                        && common.ConvertStringDateToInt(book.getBqtse_last_supply_date()) <= intUndisturbed
                        && common.ConvertStringDateToInt(book.getBqtse_last_sale_date()) <= intUndisturbed
                        && book.getBqsc_stock_count() >= intNumberOfStock
                        && book.getPercent() >= intStockPercent) {
                    if (Constants.VALUE_YES_STANDING.equals(flagSettingNew.getFlagJoubi())) {
                        if (book.getJoubi() != 5) {
                            normalSound.start();
                            item[17] = Constants.FLAG_1;
                            LogManager.i(TAG, String.format(Message.TAG_SCANNER_ACTIVITY_INLIST, book.getJan_cd(),
                                    book.getBqgm_goods_name()));
                        } else {
                            noReturnSound.start();
                            item[17] = Constants.FLAG_0;
                            LogManager.i(TAG, String.format(Message.TAG_SCANNER_ACTIVITY_INLIST, book.getJan_cd(),
                                    book.getBqgm_goods_name()));
                        }
                    } else {
                        normalSound.start();
                        item[17] = Constants.FLAG_1;
                        LogManager.i(TAG, String.format(Message.TAG_SCANNER_ACTIVITY_INLIST, book.getJan_cd(),
                                book.getBqgm_goods_name()));
                    }
                } else {
                    noReturnSound.start();
                    item[17] = Constants.FLAG_0;
                    LogManager.i(TAG, String.format(Message.TAG_SCANNER_ACTIVITY_INLIST, book.getJan_cd(),
                            book.getBqgm_goods_name()));
                }
            } else {
                if (Constants.FLAG_DEFAULT.equals(book.getFlag_sales())) {
                    normalSound.start();
                    item[17] = Constants.FLAG_1;
                    LogManager.i(TAG, String.format(Message.TAG_SCANNER_ACTIVITY_INLIST, book.getJan_cd(),
                            book.getBqgm_goods_name()));
                } else {
                    noReturnSound.start();
                    item[17] = Constants.FLAG_0;
                    LogManager.i(TAG, String.format(Message.TAG_SCANNER_ACTIVITY_INLIST, book.getJan_cd(),
                            book.getBqgm_goods_name()));
                }
            }
        }

        arrBookInlist.add(0, item);

        // Update list view
        ListViewScanAdapter adapterBook = new ListViewScanAdapter(this, arrBookInlist);
        lvBook.setAdapter(adapterBook);

    }


    private String formatOCRJAN(String jancdOCR) {
        String result;
        result = jancdOCR.replaceAll("ISBN|-| ", "");
        if (!result.startsWith(Constants.PREFIX_JAN_CODE_978)) {
            result = common.GenerateJAN(Constants.PREFIX_JAN_CODE_978 + result);
        }
        return result;
    }

    /**
     * Remove duplicate data in scanned listYên
     *
     * @param janCode {@link String}
     */

    public void checkIndexList(String janCode) {

        for (String[] item : arrBookInlist) {
            if (item[0].equals(janCode)) {
                arrBookInlist.remove(item);
                return;
            }
        }
    }

    /**
     * Show warning dialog when 取消 is clicked
     */
    private void showDialog() {

        AlertDialog.Builder dialog =
                new AlertDialog.Builder(this);
        dialog.setMessage(getString(R.string.cancel_msg)).setCancelable(false)
                .setPositiveButton(getString(R.string.logout_yes),
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                cancelItem();
                                btnCancel.setEnabled(true);
                            }
                        })
                .setNegativeButton(getString(R.string.logout_no),
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                dialog.dismiss();
                                btnCancel.setEnabled(true);
                            }
                        });
        AlertDialog alert = dialog.show();
        TextView messageText = (TextView) alert.findViewById(android.R.id.message);
        assert messageText != null;
        messageText.setGravity(Gravity.CENTER);
    }

    /**
     * Remove latest scanned item in list
     */
    private void cancelItem() {
        if (arrBookInlist.size() > 0) {
            LogManager.i(TAG,
                    String.format(Message.TAG_SCANNER_ACTIVITY_CANCEL, arrBookInlist.get(0)[0]));
            arrBookInlist.remove(0);
            // Adapter init
            // Set data adapter to list view
            ListViewScanAdapter adapterBook = new ListViewScanAdapter(this, arrBookInlist);
            lvBook.setAdapter(adapterBook);
        }
    }

    /**
     * Back event handler
     */
    @Override
    public void onBackPressed() {
        super.onBackPressed();

        LogManager.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.MESSAGE_ACTIVITY_END);

        finishAffinity();
    }

    @Override
    public boolean onKeyDown(int keyCode, KeyEvent event) {
        if (keyCode == KeyEvent.KEYCODE_BACK) {
            //Enable button setting
            imbFilter.setClickable(true);
            //Enable button switch
            aSwitchOCR.setClickable(true);
            return true;
        }
        return super.onKeyDown(keyCode, event);
    }

    /**
     * Show logout warning dialog
     */
    private void showDialog(boolean isLogout) {
        progress.show();
        android.support.v7.app.AlertDialog.Builder dialog =
                new android.support.v7.app.AlertDialog.Builder(this);
        if (isLogout) {
            dialog
                    .setMessage(getString(R.string.logout_msg))
                    .setCancelable(false)
                    .setPositiveButton(getString(R.string.logout_yes),
                            new DialogInterface.OnClickListener() {
                                @Override
                                public void onClick(DialogInterface dialog, int which) {
                                    LogManager.i(TAG,
                                            Message.TAG_SCANNER_ACTIVITY + Message.MESSAGE_LOGOUT);
                                    // print log end process
                                    LogManager.i(TAG, Message.TAG_SCANNER_ACTIVITY
                                            + Message.MESSAGE_ACTIVITY_END);
                                    // print log move screen
                                    LogManager.i(TAG,
                                            String.format(Message.MESSAGE_ACTIVITY_MOVE,
                                                    Message.SCANNER_ACTIVITY_NAME,
                                                    Message.LOGIN_ACTIVITY_NAME));
                                    compressFile();
                                    sendFileLog();
                                }
                            })
                    .setNegativeButton(getString(R.string.logout_no),
                            new DialogInterface.OnClickListener() {
                                @Override
                                public void onClick(DialogInterface dialog, int which) {
                                    dialog.dismiss();
                                    progress.dismiss();
                                    btnLogout.setEnabled(true);
                                }
                            });
        } else {
            dialog
                    .setMessage(Message.MESSAGE_NETWORK_ERR)
                    .setCancelable(false)
                    .setPositiveButton(getString(R.string.retry),
                            new DialogInterface.OnClickListener() {
                                @Override
                                public void onClick(DialogInterface dialog, int which) {
                                    sendFileLog();
                                }
                            })
                    .setNegativeButton(getString(R.string.cancel),
                            new DialogInterface.OnClickListener() {
                                @Override
                                public void onClick(DialogInterface dialog, int which) {
                                    dialog.dismiss();
                                    progress.show();
                                    delete(Constants.STRING_EMPTY);
                                    clearAndLogout();
                                    progress.dismiss();
                                }
                            });
        }
        android.support.v7.app.AlertDialog alert = dialog.show();
        TextView messageText = (TextView) alert.findViewById(android.R.id.message);
        assert messageText != null;
        messageText.setGravity(Gravity.CENTER);
    }

    /**
     * Response output
     *
     * @param output       {@link String}
     * @param typeLocation int
     * @param fileName     {@link String}
     */
    @Override
    public void progressFinish(String output, int typeLocation, String fileName) {
        if (output.contains(Message.CODE_200)) {
            countFile++;
            delete(fileName);
        } else {
            countFile++;
        }
        if (countFile == files.length) {
            clearAndLogout();
        }
    }

    /**
     * Wipe data and log out
     */
    private void clearAndLogout() {
        // Process logout
        // Wipe database
        DatabaseManager.initializeInstance(
                new DatabaseHelper(getApplicationContext()));
        DatabaseHelper ds = new DatabaseHelper(
                SdmScannerActivity.this);
        ds.clearTables();

        // Stop process loading screen
        progress.dismiss();
        finishAffinity();

        // Transition to login screen
        Intent intent = new Intent(SdmScannerActivity.this,
                LoginActivity.class);
        startActivity(intent);
    }

    /**
     * Remove old logs
     *
     * @param fileName {@link String}
     */
    public void delete(String fileName) {

        // Remove old log file for new session user
        File file = new File(Log4JHelper.fileName);
        if (file.exists()) {
            file.delete();
        }
        if (!fileName.isEmpty()) {
            File fileGz = new File(fileName);
            if (fileGz.exists()) {
                fileGz.delete();
            }
        }
    }

    /**
     * Send log file to GCS
     */
    public void sendFileLog() {

        if (checkNetwork()) {
            File root = new File(Environment.getExternalStorageDirectory(), "returncandidate_log");
            files = root.listFiles();
            for (File file : files) {
                String[] params =
                        new String[]{Config.API_POST_FILE, file.toString()};
                new HttpPostFile(this).execute(params);
            }
        } else {
            // Stop process loading screen
            progress.dismiss();
            showDialog(false);
        }
    }

    /**
     * Check network
     *
     * @return This returns true if network connectivity is okay
     */
    public boolean checkNetwork() {
        ConnectivityManager connectivityManager = (ConnectivityManager) this.getSystemService(
                Context.CONNECTIVITY_SERVICE);
        NetworkInfo networkInfo = connectivityManager.getActiveNetworkInfo();
        return !(networkInfo == null || !networkInfo.isConnected() || !networkInfo.isAvailable());
    }

    /**
     * onClick event handler
     *
     * @param v {@link View}
     */
    @Override
    public void onClick(View v) {

        switch (v.getId()) {
            case R.id.btn_logout:
                btnLogout.setEnabled(false);
                showDialog(true);
                break;
            case R.id.btn_cancel:
                if (arrBookInlist.size() > 0) {
                    btnCancel.setEnabled(false);
                    showDialog();
                }
                break;
            // Call filter
            case R.id.imb_filter:
                if (!flagFilterSubmit) {
                    putFlagOldToFlagNew();
                }
                //Disable scan
                hsmDecoder.disableSymbology(Symbology.OCR);
                hsmDecoder.disableSymbology(Symbology.EAN13);
                isEnableScan = false;
                //Disable button setting
                imbFilter.setClickable(false);
                //Enable button switch
                aSwitchOCR.setClickable(false);

                //Call filter setting when click button setting
                DFilterSettingFragment dSettingFragment = new DFilterSettingFragment();
                FragmentManager fm = getSupportFragmentManager();
                Bundle bundle = common.DataPutActivity(flagSettingNew, flagSettingOld);
                //put flag click setting
                bundle.putString(Constants.FLAG_CLICK_SETTING, Constants.VALUE_CHECK_ONCLICK_SETTING);
                //put flag switch ocr
                bundle.putString(Constants.FLAG_SWITCH_OCR, flagSwitchOCR);
                dSettingFragment.setArguments(bundle);
                dSettingFragment.show(fm, null);
                break;
        }
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

    /**
     * On Pause app when go use another app
     */
    @Override
    protected void onPause() {
        super.onPause();
        timeout = System.currentTimeMillis();

        HSMDecoder.disposeInstance();
    }

    @Override
    protected void onStop() {
        super.onStop();
    }

    /**
     * onRestart event handler
     */
    @Override
    protected void onRestart() {

        // Print log end process
        LogManager.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.MESSAGE_ACTIVITY_END);
        // Print log move screen
        LogManager.i(TAG, String.format(Message.MESSAGE_ACTIVITY_MOVE,
                Message.SCANNER_ACTIVITY_NAME,
                Message.UNLOCK_ACTIVITY_NAME));
        super.onRestart();
        if (timeout < (System.currentTimeMillis() - Constants.TIME_OUT)) {
            Intent intent = new Intent(this, UnlockScreenActivity.class);
            Bundle bundle = new Bundle();
            HashMap<String, LinkedList<String[]>> hashMapArrBook = new HashMap<>();
            hashMapArrBook.put(Constants.COLUMN_INFOR_LIST_SCAN, arrBookInlist);
            bundle.putString(Constants.COLUMN_USER_ID, userID);
            bundle.putString(Constants.COLUMN_SHOP_ID, shopID);
            bundle.putString(Constants.COLUMN_SERVER_NAME, serverName);
            bundle.putString(Constants.COLUMN_LICENSE, license);
            bundle.putSerializable(Constants.COLUMN_INFOR_LIST_SCAN, hashMapArrBook);
            intent.putExtras(bundle);
            startActivity(intent);
            finish();
            Log.d("LINCONGPVUNCLOCK", license);
        }

//         Activate HSM license
        Log.d("LINCONGPVUNCLOCK1", license);
        ActivationResult activationResult = ActivationManager.activate(this,
                license);
        Toast.makeText(this, "Activation result: " + activationResult, Toast.LENGTH_LONG).show();

        decCom = (HSMDecodeComponent) findViewById(R.id.hsm_decodeComponent);

        // HSM init
        hsmDecoder = HSMDecoder.getInstance(this);

        // Declare symbology
        if (aSwitchOCR.isChecked()) {
            flagSwitchOCR = Constants.FLAG_1;
            hsmDecoder.enableSymbology(Symbology.OCR);
            hsmDecoder.setOCRActiveTemplate(OCRActiveTemplate.ISBN);
            hsmDecoder.disableSymbology(Symbology.EAN13);
        } else {
            flagSwitchOCR = Constants.FLAG_0;
            hsmDecoder.enableSymbology(Symbology.EAN13);
            //hsmDecoder.setOCRActiveTemplate(OCRActiveTemplate.ISBN);
            hsmDecoder.disableSymbology(Symbology.OCR);

        }

        if (!isEnableScan){
            hsmDecoder.disableSymbology(Symbology.EAN13);
            hsmDecoder.disableSymbology(Symbology.OCR);
        }
        //hsmDecoder.enableSymbology(Symbology.EAN13);

        // Declare HSM component UI
        hsmDecoder.enableFlashOnDecode(false);
        hsmDecoder.enableSound(false);
        hsmDecoder.enableAimer(false);
        hsmDecoder.setWindowMode(WindowMode.CENTERING);
        hsmDecoder.setWindow(18, 42, 0, 100);

        // Assign listener
        hsmDecoder.addResultListener(this);
    }

    /**
     * onDestroy event handler<br>
     * Unregister all HSM instance
     */
    @Override
    public void onDestroy() {

        super.onDestroy();
    }

    /**
     * Log file compression
     */
    private void compressFile() {
        @SuppressLint("SimpleDateFormat") SimpleDateFormat dateFormat = new SimpleDateFormat(
                "yyyyMMddHHmmss");
        Calendar cal = Calendar.getInstance();
        String strDate = dateFormat.format(cal.getTime());

        // Compress file log
        String filepath = Log4JHelper.fileName;
        File root = new File(Environment.getExternalStorageDirectory(), "returncandidate_log");
        if (!root.exists()) {
            root.mkdirs();
        }

        userID = userID.replaceAll("__", "_");
        serverName = serverName.replaceAll("__", "_");
        shopID = shopID.replaceAll("__", "_");

        //Gzip file name
        String gzipFileName = root + "/" + serverName + "__" + shopID + "__" + userID + "__"
                + strDate + ".log.gz";
        GzipFile gzipFile = new GzipFile();
        gzipFile.compressGzipFile(filepath, gzipFileName);
    }

    @Override
    protected void onResume() {
        super.onResume();

        //Event click item detail
        lvBook.setOnItemClickListener(new AdapterView.OnItemClickListener() {
            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
                DProductDetailFragment dProductDetailFragment = new DProductDetailFragment();
                FragmentManager fm = getSupportFragmentManager();
                Bundle bundleItems = new Bundle();
                bundleItems.putString(Constants.COLUMN_JAN_CD, arrBookInlist.get(position)[0]);
                bundleItems.putString(Constants.COLUMN_STOCK_COUNT, arrBookInlist.get(position)[1]);
                bundleItems.putString(Constants.COLUMN_GOODS_NAME, arrBookInlist.get(position)[2]);
                bundleItems.putString(Constants.COLUMN_WRITER_NAME, arrBookInlist.get(position)[3]);
                bundleItems.putString(Constants.COLUMN_PUBLISHER_NAME, arrBookInlist.get(position)[4]);
                bundleItems.putString(Constants.COLUMN_PRICE, arrBookInlist.get(position)[5]); // TODO
                bundleItems.putString(Constants.COLUMN_FIRST_SUPPLY_DATE, arrBookInlist.get(position)[6]);
                bundleItems.putString(Constants.COLUMN_LAST_SUPPLY_DATE, arrBookInlist.get(position)[7]);
                bundleItems.putString(Constants.COLUMN_LAST_SALES_DATE, arrBookInlist.get(position)[8]);
                bundleItems.putString(Constants.COLUMN_LAST_ORDER_DATE, arrBookInlist.get(position)[9]);
                bundleItems.putString(Constants.COLUMN_MEDIA_GROUP1_CD, arrBookInlist.get(position)[10]);
                bundleItems.putString(Constants.COLUMN_MEDIA_GROUP1_NAME, arrBookInlist.get(position)[11]);
                bundleItems.putString(Constants.COLUMN_MEDIA_GROUP2_CD, arrBookInlist.get(position)[12]);
                bundleItems.putString(Constants.COLUMN_MEDIA_GROUP2_NAME, arrBookInlist.get(position)[13]);
                bundleItems.putString(Constants.COLUMN_SALES_DATE, arrBookInlist.get(position)[14]);
                bundleItems.putString(Constants.COLUMN_YEAR_RANK, arrBookInlist.get(position)[18]);
                bundleItems.putString(Constants.COLUMN_JOUBI, arrBookInlist.get(position)[19]);
                bundleItems.putString(Constants.COLUMN_TOTAL_SALES, arrBookInlist.get(position)[20]);
                bundleItems.putString(Constants.COLUMN_TOTAL_SUPPLY, arrBookInlist.get(position)[21]);
                bundleItems.putString(Constants.COLUMN_TOTAL_RETURN, arrBookInlist.get(position)[22]);
                bundleItems.putString(Constants.COLUMN_LOCATION_ID, arrBookInlist.get(position)[23]);
                bundleItems.putInt(Constants.COLUMN_MAX_YEAR_RANK, maxYearRank.getMaxYearRank());
                dProductDetailFragment.setArguments(bundleItems);
                dProductDetailFragment.show(fm, null);
            }
        });
    }

    @Override
    public void onLitSelectedFilterSettingDialog(FlagSettingNew _flagSettingNew, FlagSettingOld
            _flagSettingOld, Boolean _flagFilterSubmit, String flagSwitchOCR) {
        //Enable button filter
        imbFilter.setClickable(true);
        //Enable button switch
        aSwitchOCR.setClickable(true);
        //Enable scan
        isEnableScan = true;

        flagSettingNew = _flagSettingNew;
        flagSettingOld = _flagSettingOld;
        flagFilterSubmit = _flagFilterSubmit;
        if (!flagFilterSubmit) {
            if (Constants.FLAG_1.equals(flagSwitchOCR)) {
                hsmDecoder.enableSymbology(Symbology.OCR);
            } else {
                hsmDecoder.enableSymbology(Symbology.EAN13);
            }
        }
    }
}
