package com.fjn.magazinereturncandidate.activities;

import android.annotation.SuppressLint;
import android.app.ProgressDialog;
import android.content.DialogInterface;
import android.content.Intent;
import android.graphics.Bitmap;
import android.media.MediaPlayer;
import android.os.Bundle;
import android.os.Environment;
import android.support.v4.app.FragmentManager;
import android.support.v7.app.AlertDialog;
import android.support.v7.app.AppCompatActivity;
import android.view.Gravity;
import android.view.KeyEvent;
import android.view.View;
import android.widget.AdapterView;
import android.widget.Button;
import android.widget.CompoundButton;
import android.widget.ListView;
import android.widget.Switch;
import android.widget.TextView;
import android.widget.Toast;

import com.fjn.magazinereturncandidate.R;
import com.fjn.magazinereturncandidate.adapters.ListViewScanAdapter;
import com.fjn.magazinereturncandidate.api.Config;
import com.fjn.magazinereturncandidate.api.HttpPostFile;
import com.fjn.magazinereturncandidate.api.HttpPostFileDataReturnMagazine;
import com.fjn.magazinereturncandidate.api.HttpResponse;
import com.fjn.magazinereturncandidate.common.constants.Constants;
import com.fjn.magazinereturncandidate.common.constants.Message;
import com.fjn.magazinereturncandidate.common.helpers.DatabaseHelper;
import com.fjn.magazinereturncandidate.common.helpers.Log4JHelper;
import com.fjn.magazinereturncandidate.common.utils.CheckDataCommon;
import com.fjn.magazinereturncandidate.common.utils.DatabaseManagerCommon;
import com.fjn.magazinereturncandidate.common.utils.FormatCommon;
import com.fjn.magazinereturncandidate.common.utils.GzipFileCommon;
import com.fjn.magazinereturncandidate.common.utils.LogManagerCommon;
import com.fjn.magazinereturncandidate.common.utils.CSVFileCommon;
import com.fjn.magazinereturncandidate.common.utils.MyCustomPlugin;
import com.fjn.magazinereturncandidate.common.utils.MyCustomPluginResultListener;
import com.fjn.magazinereturncandidate.common.utils.RegisterLicenseCommon;
import com.fjn.magazinereturncandidate.db.entity.MaxYearRankEntity;
import com.fjn.magazinereturncandidate.db.entity.ReturnMagazineEntity;
import com.fjn.magazinereturncandidate.db.entity.UsersEntity;
import com.fjn.magazinereturncandidate.db.models.MaxYearRankModel;
import com.fjn.magazinereturncandidate.db.models.ReturnMagazineModel;
import com.fjn.magazinereturncandidate.db.models.UserModel;
import com.fjn.magazinereturncandidate.fragments.InputJanCodeFragment;
import com.fjn.magazinereturncandidate.fragments.ProductDetailFragment;
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
import java.util.Calendar;
import java.util.HashMap;
import java.util.LinkedList;
import java.util.List;

import static com.fjn.magazinereturncandidate.common.constants.Message.MESSAGE_CANCEL;
import static com.fjn.magazinereturncandidate.common.constants.Message.MESSAGE_CONFIRM_LOGOUT;
import static com.fjn.magazinereturncandidate.common.constants.Message.MESSAGE_CONFIRM_SEND_DATA;
import static com.fjn.magazinereturncandidate.common.constants.Message.MESSAGE_NO;
import static com.fjn.magazinereturncandidate.common.constants.Message.MESSAGE_NO_JP;
import static com.fjn.magazinereturncandidate.common.constants.Message.MESSAGE_RETRY;
import static com.fjn.magazinereturncandidate.common.constants.Message.MESSAGE_SEND_LIST_FILE_CSV;
import static com.fjn.magazinereturncandidate.common.constants.Message.MESSAGE_YES;
import static com.fjn.magazinereturncandidate.common.constants.Message.MESSAGE_YES_JP;
import static com.fjn.magazinereturncandidate.common.utils.NetworkCommon.isNetworkConnection;

/**
 * <h1>Main barcode scanner activity.</h1>
 * Begin scanning barcode, if scanned barcode is matched with database will result with a beep
 * sound, otherwise will result with a buzz sound. User can remove recently scanned object via 取消
 * button
 * and logout via ログアウト button
 *
 * @author cong-pv
 * @version 1.0
 * @since 2018-10-15
 */
public class SdmScannerActivity extends AppCompatActivity implements View.OnClickListener,
        HttpResponse, MyCustomPluginResultListener, ProductDetailFragment.SubmitEditDataNumberReturnMagazine,
        InputJanCodeFragment.SubmitInputJanDataReturnMagazine {

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
    private Button btnLogout;
    private Button btnSend;
    private Button btnInputJan;
    private ListView lvBook;
    private MediaPlayer normalSound;
    private MediaPlayer noReturnSound;

    long timeout;

    private Switch aSwitchOCR;
    private String flagSwitchOCR;
    private Boolean isEnableScan = true;

    private ReturnMagazineModel returnMagazineModel;
    private MaxYearRankEntity maxYearRank;

    private CheckDataCommon checkDataCommon;
    private CSVFileCommon csvFileCommon;
    private RegisterLicenseCommon registerLicenseCommon;
    private FormatCommon formatCommon;
    private long endScan = 0;
    private String strBarcodeOld;

    private String janCodeNew18Character;


    private MyCustomPlugin customPlugin;


    /**
     * Initialize screen layout
     *
     * @param state {@link Bundle}
     */
    @Override
    public void onCreate(Bundle state) {

        LogManagerCommon.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.MESSAGE_ACTIVITY_START);
        super.onCreate(state);

        setContentView(R.layout.activity_sdm_scanner);

        returnMagazineModel = new ReturnMagazineModel();
        csvFileCommon = new CSVFileCommon();
        registerLicenseCommon = new RegisterLicenseCommon();
        maxYearRank = new MaxYearRankEntity();
        checkDataCommon = new CheckDataCommon();
        formatCommon = new FormatCommon();
        strBarcodeOld = "";

        // Load user info
        Bundle bundle = getIntent().getExtras();
        if (bundle != null) {
            //get flag switch OCR
            flagSwitchOCR = bundle.getString(Constants.FLAG_SWITCH_OCR);
            userID = bundle.getString(Constants.COLUMN_USER_ID);
            shopID = bundle.getString(Constants.COLUMN_SHOP_ID);
            serverName = bundle.getString(Constants.COLUMN_SERVER_NAME);
            license = bundle.getString(Constants.COLUMN_LICENSE);
            hashMapArrBook = (HashMap<String, LinkedList<String[]>>) bundle.getSerializable(Constants.COLUMN_INFOR_LIST_SCAN);
        }
        // UI init
        UsersEntity user = new UserModel().getUserInfo();
        TextView tvUserName = (TextView) findViewById(R.id.txv_user_name);
        tvUserName.setText(user.getName());

        lvBook = (ListView) findViewById(R.id.list_book);
        arrBookInlist = new LinkedList<>();
        if (hashMapArrBook != null) {
            arrBookInlist = new LinkedList<>(hashMapArrBook.get(Constants.COLUMN_INFOR_LIST_SCAN));
        }


        //Check OCR null
        if (flagSwitchOCR == null) {
            //Flag disable OCR (default)
            flagSwitchOCR = Constants.FLAG_0;
        }
        MaxYearRankModel maxYearRankModel = new MaxYearRankModel();
        DatabaseManagerCommon.initializeInstance(
                new DatabaseHelper(getApplicationContext()));

        maxYearRank = maxYearRankModel.getMaxYearRank();

        // Init process loading screen
        progress = new ProgressDialog(this);
        progress.setMessage(Message.MESSAGE_UPLOAD_LOG_SCREEN);
        progress.setCancelable(false);

        // Activate HSM license
        ActivationResult activationResult = ActivationManager.activate(this,
                license);
        Toast.makeText(this, "Activation result: " + activationResult, Toast.LENGTH_LONG).show();

        // HSM init
        hsmDecoder = HSMDecoder.getInstance(this);

        // Declare symbology
        hsmDecoder.enableSymbology(Symbology.EAN13);
        hsmDecoder.enableSymbology(Symbology.EAN13_5CHAR_ADDENDA);
        hsmDecoder.enableSymbology(Symbology.CODE128);
        //hsmDecoder.enableSymbology(Symbology.OCR);
        //hsmDecoder.setOCRActiveTemplate(OCRActiveTemplate.ISBN);

        //create plug-in instance and add a result listener
        customPlugin = new MyCustomPlugin(getApplicationContext(), 100);
        customPlugin.addResultListener(this);
        //register the plug-in with the system
        hsmDecoder.registerPlugin(customPlugin);


        // Declare HSM component UI
        hsmDecoder.enableFlashOnDecode(false);
        hsmDecoder.enableSound(false);
        hsmDecoder.enableAimer(false);
        hsmDecoder.setWindowMode(WindowMode.CENTERING);
        hsmDecoder.setWindow(18, 42, 0, 100);

        // Assign listener
        //hsmDecoder.addResultListener(this);

        // Sound init
        normalSound = MediaPlayer.create(this,
                R.raw.pingpong_main); // sound is inside res/raw/mysound
        //totalTimeNormalAudio = normalSound.getDuration();

        noReturnSound = MediaPlayer.create(this,
                R.raw.wrong_main); // sound is inside res/raw/mysound
        //totalTimeNoReturnAudio = noReturnSound.getDuration();

        // Button init
        btnLogout = (Button) findViewById(R.id.btn_logout);
        btnInputJan = (Button) findViewById(R.id.btn_input_jancode);
        btnSend = (Button) findViewById(R.id.btn_send_data);

        aSwitchOCR = (Switch) findViewById(R.id.switch_OCR);
        btnLogout.setOnClickListener(this);
        btnSend.setOnClickListener(this);
        btnInputJan.setOnClickListener(this);
/*
        aSwitchOCR.setOnCheckedChangeListener(new CompoundButton.OnCheckedChangeListener() {
            @Override
            public void onCheckedChanged(CompoundButton buttonView, boolean isChecked) {
                if (aSwitchOCR.isChecked()) {
                    flagSwitchOCR = registerLicenseCommon.EnableOCRDisableJanCode(hsmDecoder);
                } else {
                    flagSwitchOCR = registerLicenseCommon.EnableJanCodeDisableOCR(hsmDecoder);
                }
            }
        });*/

        //Check if arr list books not null
        if (arrBookInlist != null) {
            // Set data adapter to list view
            ListViewScanAdapter adapterBook = new ListViewScanAdapter(this, arrBookInlist);
            lvBook.setAdapter(adapterBook);
        }

        //Check network send file and check csvFileCommon.isExistFileCSV()
        if (csvFileCommon.isExistFileCSV() != 0) {
            //Disable scan
            registerLicenseCommon.DisableScan(hsmDecoder);
            //Show dialog choose send data
            showDialogNotifyListCSVWhenLogin();
        }

    }
    //callback method that returns the plug-in result

    @Override
    public void onMyCustomPluginResult(HSMDecodeResult[] hsmDecodeResults) {

       //int abc =  customPlugin.getDecodeResultRequirement();

       /* if (hsmDecodeResults.length > 0) {

            long startScan = formatCommon.getCurrentTimeLong();
            String strBarcodeNew;
            StringBuilder stringBuilder = new StringBuilder();
            for (HSMDecodeResult hsmDecodeResult : hsmDecodeResults) {
                String strJanCode = hsmDecodeResult.getBarcodeData();
                stringBuilder.append(strJanCode);
            }
            strBarcodeNew = stringBuilder.toString();

            if (endScan == 0) {
                for (HSMDecodeResult hsmDecodeResult : hsmDecodeResults) {
                    barcodeValidation(hsmDecodeResult.getBarcodeData());
                }
            } else {
                if (!strBarcodeOld.equals(strBarcodeNew) || (startScan - endScan) >= 1000) {
                    for (HSMDecodeResult hsmDecodeResult : hsmDecodeResults) {
                        barcodeValidation(hsmDecodeResult.getBarcodeData());
                    }
                }
            }

            strBarcodeOld = strBarcodeNew;
            endScan = formatCommon.getCurrentTimeLong();

        }*/
    }

    /**
     * Decoded results handler
     *
     * @param hsmDecodeResults {@link HSMDecodeResult}
     *//*
    @Override
    public void onHSMDecodeResult(HSMDecodeResult[] hsmDecodeResults) {

        if (hsmDecodeResults.length > 0) {

            long startScan = formatCommon.getCurrentTimeLong();
            String strBarcodeNew;
            StringBuilder stringBuilder = new StringBuilder();
            for (HSMDecodeResult hsmDecodeResult : hsmDecodeResults) {
                String strJanCode = hsmDecodeResult.getBarcodeData();
                stringBuilder.append(strJanCode);
            }
            strBarcodeNew = stringBuilder.toString();

            if (endScan == 0) {
                for (HSMDecodeResult hsmDecodeResult : hsmDecodeResults) {
                    barcodeValidation(hsmDecodeResult.getBarcodeData());
                }
            } else {
                if (!strBarcodeOld.equals(strBarcodeNew) || (startScan - endScan) >= 1000) {
                    for (HSMDecodeResult hsmDecodeResult : hsmDecodeResults) {
                        barcodeValidation(hsmDecodeResult.getBarcodeData());
                    }
                }
            }

            strBarcodeOld = strBarcodeNew;
            endScan = formatCommon.getCurrentTimeLong();

            //Delay scan
           *//* try {
                TimeUnit.MILLISECONDS.sleep(700);
            } catch (InterruptedException ignored) {
            }*//*
        }
    }*/

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
        //if (aSwitchOCR.isChecked()) {
        if (barcode.substring(0, 4).equals("ISBN")) {
            barcode = checkDataCommon.validateOCR(barcode);
        }

        //}

        int lenJanCodeInput = barcode.length();
        String janCode18Character = barcode.substring(0, lenJanCodeInput - 5);

        //Get info
        ReturnMagazineEntity magazineEntity;
        if (lenJanCodeInput == Constants.JAN_18_CHAR) {
            magazineEntity = returnMagazineModel.getItemReturnMagazine(janCode18Character);
        } else {
            magazineEntity = returnMagazineModel.getItemReturnMagazine(barcode);
        }

        //Set default
        String strPrice = Constants.BLANK;
        if (magazineEntity.getBqgm_price() != null) {
            strPrice = magazineEntity.getBqgm_price().toString();
        }

        int checkValueRecordNumber = checkIndexList(barcode);

        //Check jan_code in list and remove
        String JanCodeResult;
        if (lenJanCodeInput == Constants.JAN_18_CHAR) {
            JanCodeResult = barcode;
        } else {
            if (janCodeNew18Character != null) {
                JanCodeResult = janCodeNew18Character;
            } else {
                JanCodeResult = barcode;
            }
        }
        //Reset janCode
        janCodeNew18Character = null;

        int valueRecordNumber;
        if (checkValueRecordNumber != -1) {
            valueRecordNumber = checkValueRecordNumber + 1;
        } else {
            if (magazineEntity.getBqsc_stock_count() == 0) {
                valueRecordNumber = 1;
            } else {
                valueRecordNumber = magazineEntity.getBqsc_stock_count();
            }
        }

        String[] item = new String[]{JanCodeResult, String.valueOf(valueRecordNumber), magazineEntity.getBqgm_goods_name(),
                magazineEntity.getBqgm_writer_name(), magazineEntity.getBqgm_publisher_name(), strPrice,
                magazineEntity.getBqtse_first_supply_date(), magazineEntity.getBqtse_last_supply_date(),
                magazineEntity.getBqtse_last_sale_date(), magazineEntity.getBqtse_last_order_date(),
                magazineEntity.getBqct_media_group1_cd(), magazineEntity.getBqct_media_group1_name(),
                magazineEntity.getBqct_media_group2_cd(), magazineEntity.getBqct_media_group2_name(),
                magazineEntity.getBqgm_sales_date(), magazineEntity.getBqgm_publisher_cd(), magazineEntity.getBqio_trn_date(),
                Constants.FLAG_0, String.valueOf(magazineEntity.getYear_rank()),
                String.valueOf(magazineEntity.getSts_total_sales()), String.valueOf(magazineEntity.getSts_total_supply()),
                String.valueOf(magazineEntity.getSts_total_return()), magazineEntity.getLocation_id()};


        if (!JanCodeResult.equals(magazineEntity.getJan_cd())) {
            noReturnSound.start();
            LogManagerCommon.i(TAG, String.format(Message.TAG_SCANNER_ACTIVITY_OUT_LIST, JanCodeResult));
        } else {
            normalSound.start();
            item[17] = Constants.FLAG_1;
            LogManagerCommon.i(TAG, String.format(Message.TAG_SCANNER_ACTIVITY_INLIST, JanCodeResult,
                    magazineEntity.getBqgm_goods_name()));
        }

        arrBookInlist.add(0, item);

        // Update list view
        ListViewScanAdapter adapterBook = new ListViewScanAdapter(this, arrBookInlist);
        lvBook.setAdapter(adapterBook);

    }

    /**
     * Remove duplicate data in scanned listYên
     *
     * @param janCode {@link String}
     */

    public int checkIndexList(String janCode) {

        int saveItemRemove = -1;
        for (String[] item : arrBookInlist) {
            if (item[0].substring(0, 13).equals(janCode.substring(0, 13))) {
                janCodeNew18Character = item[0];
                saveItemRemove = Integer.parseInt(item[1]);
                arrBookInlist.remove(item);
                break;
            }

        }
        return saveItemRemove;
    }

    /**
     * Remove duplicate data in scanned when input jan
     *
     * @param janCode {@link String}
     */

    public void checkIndexListInputJan(String janCode) {

        for (String[] item : arrBookInlist) {
            if (item[0].substring(0, 13).equals(janCode.substring(0, 13))) {
                janCodeNew18Character = item[0];
                arrBookInlist.remove(item);
                break;
            }

        }
    }

    /**
     * Back event handler
     */
    @Override
    public void onBackPressed() {

        super.onBackPressed();

        LogManagerCommon.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.MESSAGE_ACTIVITY_END);

        finishAffinity();
    }

    @Override
    public boolean onKeyDown(int keyCode, KeyEvent event) {

        if (keyCode == KeyEvent.KEYCODE_BACK) {
            //Enable button setting
            btnInputJan.setClickable(true);
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
        AlertDialog.Builder dialog =
                new AlertDialog.Builder(this);
        if (isLogout) {
            dialog
                    .setMessage(MESSAGE_CONFIRM_LOGOUT)
                    .setCancelable(false)
                    .setPositiveButton(MESSAGE_YES_JP,
                            new DialogInterface.OnClickListener() {
                                @Override
                                public void onClick(DialogInterface dialog, int which) {
                                    LogManagerCommon.i(TAG,
                                            Message.TAG_SCANNER_ACTIVITY + Message.MESSAGE_LOGOUT);
                                    // print log end process
                                    LogManagerCommon.i(TAG, Message.TAG_SCANNER_ACTIVITY
                                            + Message.MESSAGE_ACTIVITY_END);
                                    // print log move screen
                                    LogManagerCommon.i(TAG,
                                            String.format(Message.MESSAGE_ACTIVITY_MOVE,
                                                    Message.SCANNER_ACTIVITY_NAME,
                                                    Message.LOGIN_ACTIVITY_NAME));
                                    compressFile();
                                    sendFileLog();
                                }
                            })
                    .setNegativeButton(MESSAGE_NO_JP,
                            new DialogInterface.OnClickListener() {
                                @Override
                                public void onClick(DialogInterface dialog, int which) {
                                    dialog.dismiss();
                                    progress.dismiss();
                                    btnLogout.setEnabled(true);
                                    isEnableScan = true;
                                    registerLicenseCommon.EnableOCROrJanCode(flagSwitchOCR, hsmDecoder);
                                }
                            });
        } else {
            dialog
                    .setMessage(Message.MESSAGE_NETWORK_ERR)
                    .setCancelable(false)
                    .setPositiveButton(MESSAGE_RETRY,
                            new DialogInterface.OnClickListener() {
                                @Override
                                public void onClick(DialogInterface dialog, int which) {
                                    sendFileLog();
                                }
                            })
                    .setNegativeButton(MESSAGE_CANCEL,
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
        AlertDialog alert = dialog.show();
        TextView messageText = (TextView) alert.findViewById(android.R.id.message);
        assert messageText != null;
        messageText.setGravity(Gravity.CENTER);
    }

    /**
     * Show warning dialog when send data
     */
    private void showDialogWarningSendData() {

        progress.show();
        AlertDialog.Builder dialog =
                new AlertDialog.Builder(this);
        dialog
                .setMessage(MESSAGE_CONFIRM_SEND_DATA)
                .setCancelable(false)
                .setPositiveButton(MESSAGE_NO,
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                dialog.dismiss();
                                progress.dismiss();
                                btnSend.setEnabled(true);
                                isEnableScan = true;
                                registerLicenseCommon.EnableOCROrJanCode(flagSwitchOCR, hsmDecoder);
                            }
                        })
                .setNegativeButton(MESSAGE_YES,
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                LogManagerCommon.i(TAG,
                                        Message.TAG_SCANNER_ACTIVITY + Message.MESSAGE_SEND_DATA);
                                // print log end process
                                LogManagerCommon.i(TAG, Message.TAG_SCANNER_ACTIVITY
                                        + Message.MESSAGE_ACTIVITY_END);
                                // print log send file data
                                LogManagerCommon.i(TAG, Message.MESSAGE_SEND_DATA);

                                //Save csv and send list file to Bigquery
                                sendData();
                            }
                        });

        AlertDialog alert = dialog.show();
        TextView messageText = (TextView) alert.findViewById(android.R.id.message);
        assert messageText != null;
        messageText.setGravity(Gravity.CENTER);
    }

    /**
     * Show message info list file csv
     */
    private void showDialogNotifyListCSVWhenLogin() {

        progress.show();
        AlertDialog.Builder dialog =
                new AlertDialog.Builder(this);
        dialog
                .setMessage(MESSAGE_SEND_LIST_FILE_CSV)
                .setCancelable(false)
                .setPositiveButton(MESSAGE_NO,
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                dialog.dismiss();
                                progress.dismiss();
                                btnSend.setEnabled(true);
                                //Enable scan
                                isEnableScan = true;
                                registerLicenseCommon.EnableOCROrJanCode(flagSwitchOCR, hsmDecoder);
                            }
                        })
                .setNegativeButton(MESSAGE_YES,
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                LogManagerCommon.i(TAG,
                                        Message.TAG_SCANNER_ACTIVITY + Message.MESSAGE_SEND_DATA);
                                // print log end process
                                LogManagerCommon.i(TAG, Message.TAG_SCANNER_ACTIVITY
                                        + Message.MESSAGE_ACTIVITY_END);
                                // print log send file data
                                LogManagerCommon.i(TAG, Message.MESSAGE_SEND_DATA);

                                //Send list file csv to bigQuery
                                connectSendDataWhenLogin();
                            }
                        });
        AlertDialog alert = dialog.show();
        TextView messageText = (TextView) alert.findViewById(android.R.id.message);
        assert messageText != null;
        messageText.setGravity(Gravity.CENTER);
    }

    /**
     * Function check connect send data when login success
     */
    private void connectSendDataWhenLogin() {

        //Check network
        if (isNetworkConnection(SdmScannerActivity.this)) {
            File root = new File(Environment.getExternalStorageDirectory(), Config.FOLDER_ANDROID_SAVE_DATA);
            files = root.listFiles();
            for (File file : files) {
                String[] params =
                        new String[]{Config.API_POST_FILE_DATA_RETURN_MAGAZINE, file.toString()};
                new HttpPostFileDataReturnMagazine(this).execute(params);
            }
        } else {
            //Not connect
            // stop process loading screen
            progress.dismiss();
            showDialogConnectFailedWhenLogin();
        }
    }

    /**
     * Function show dialog when connect fail (not internet - retry)
     */
    private void showDialogConnectFailedWhenLogin() {

        progress.show();
        android.support.v7.app.AlertDialog.Builder dialog =
                new android.support.v7.app.AlertDialog.Builder(this);
        dialog
                .setMessage(Message.MESSAGE_NETWORK_ERR)
                .setCancelable(false)
                .setPositiveButton(MESSAGE_RETRY,
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                connectSendDataWhenLogin();
                            }
                        })
                .setNegativeButton(MESSAGE_CANCEL,
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                dialog.dismiss();
                                progress.dismiss();
                                //Enable scan
                                isEnableScan = true;
                                registerLicenseCommon.EnableOCROrJanCode(flagSwitchOCR, hsmDecoder);
                            }
                        });
        android.support.v7.app.AlertDialog alert = dialog.show();
        TextView messageText = (TextView) alert.findViewById(android.R.id.message);
        assert messageText != null;
        messageText.setGravity(Gravity.CENTER);
    }

    /**
     * Wipe data and log out
     */
    private void clearAndLogout() {
        // Process logout
        // Wipe database
        DatabaseManagerCommon.initializeInstance(
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

        if (isNetworkConnection(SdmScannerActivity.this)) {
            File root = new File(Environment.getExternalStorageDirectory(), "magazinereturncandidate_log");
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
     * Function send data when click send data in screen scan barcode
     */
    private void sendData() {

        //Save file csv
        csvFileCommon.saveCSVFile(userID, shopID, serverName, arrBookInlist);

        //Delete all data listView scan
        arrBookInlist = new LinkedList<>();

        // Update list view
        ListViewScanAdapter adapterBook = new ListViewScanAdapter(this, arrBookInlist);
        lvBook.setAdapter(adapterBook);

        //Send data
        connectSendData();

    }

    /**
     * Check connect when send data
     */
    private void connectSendData() {

        //Check network
        if (isNetworkConnection(SdmScannerActivity.this)) {
            File root = new File(Environment.getExternalStorageDirectory(), Config.FOLDER_ANDROID_SAVE_DATA);
            files = root.listFiles();
            for (File file : files) {
                String[] params =
                        new String[]{Config.API_POST_FILE_DATA_RETURN_MAGAZINE, file.toString()};
                new HttpPostFileDataReturnMagazine(this).execute(params);
            }
        } else {
            //Not connect
            // stop process loading screen
            progress.dismiss();
            showDialog();
        }
    }

    /**
     * Warning dialog when failed to connect ) (retry function connectSendData())
     */
    private void showDialog() {

        progress.show();
        android.support.v7.app.AlertDialog.Builder dialog =
                new android.support.v7.app.AlertDialog.Builder(this);
        dialog
                .setMessage(Message.MESSAGE_NETWORK_ERR)
                .setCancelable(false)
                .setPositiveButton(MESSAGE_RETRY,
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                connectSendData();
                            }
                        })
                .setNegativeButton(MESSAGE_CANCEL,
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                dialog.dismiss();
                                progress.dismiss();
                                showDialogInfoPathSaveCSV();
                                //Enable scan
                                isEnableScan = true;
                                registerLicenseCommon.EnableOCROrJanCode(flagSwitchOCR, hsmDecoder);
                            }
                        });
        android.support.v7.app.AlertDialog alert = dialog.show();
        TextView messageText = (TextView) alert.findViewById(android.R.id.message);
        assert messageText != null;
        messageText.setGravity(Gravity.CENTER);
    }

    /**
     * Function show path file save CSV when not connect internet
     */
    private void showDialogInfoPathSaveCSV() {

        progress.show();
        android.support.v7.app.AlertDialog.Builder dialog =
                new android.support.v7.app.AlertDialog.Builder(this);
        dialog
                .setMessage(Message.MESSAGE_INFO_PATH_CSV)
                .setCancelable(false)
                .setPositiveButton(MESSAGE_YES,
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
     * onClick event handler
     *
     * @param v {@link View}
     */
    @Override
    public void onClick(View v) {

        switch (v.getId()) {
            case R.id.btn_logout:
                btnLogout.setEnabled(false);
                registerLicenseCommon.DisableScan(hsmDecoder);
                isEnableScan = false;
                showDialog(true);
                break;
            case R.id.btn_input_jancode:
                registerLicenseCommon.DisableScan(hsmDecoder);
                isEnableScan = false;
                showFragmentInputJan();
                break;
            case R.id.btn_send_data:
                registerLicenseCommon.DisableScan(hsmDecoder);
                isEnableScan = false;
                if (arrBookInlist.size() > 0) {
                    showDialogWarningSendData();
                }
                break;
        }
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

    /**
     * onRestart event handler
     */
    @Override
    protected void onRestart() {

        // Print log end process
        LogManagerCommon.i(TAG, Message.TAG_SCANNER_ACTIVITY + Message.MESSAGE_ACTIVITY_END);
        // Print log move screen
        LogManagerCommon.i(TAG, String.format(Message.MESSAGE_ACTIVITY_MOVE,
                Message.SCANNER_ACTIVITY_NAME,
                Message.UNLOCK_ACTIVITY_NAME));
        super.onRestart();
        if (timeout < (System.currentTimeMillis() - Constants.TIME_OUT)) {
            Intent intent = new Intent(this, UnlockScreenActivity.class);
            HashMap<String, LinkedList<String[]>> hashMapArrBook = new HashMap<>();
            hashMapArrBook.put(Constants.COLUMN_INFOR_LIST_SCAN, arrBookInlist);

            Bundle bundle = new Bundle();
            //put flag switch OCR
            bundle.putString(Constants.FLAG_SWITCH_OCR, flagSwitchOCR);
            bundle.putString(Constants.COLUMN_USER_ID, userID);
            bundle.putString(Constants.COLUMN_SHOP_ID, shopID);
            bundle.putString(Constants.COLUMN_SERVER_NAME, serverName);
            bundle.putString(Constants.COLUMN_LICENSE, license);
            bundle.putSerializable(Constants.COLUMN_INFOR_LIST_SCAN, hashMapArrBook);
            intent.putExtras(bundle);
            startActivity(intent);
            finish();
        }

//         Activate HSM license
        ActivationResult activationResult = ActivationManager.activate(this,
                license);
        Toast.makeText(this, "Activation result: " + activationResult, Toast.LENGTH_LONG).show();

        // HSM init
        hsmDecoder = HSMDecoder.getInstance(this);

        // Declare symbology
        if (aSwitchOCR.isChecked()) {
            registerLicenseCommon.EnableOCRDisableJanCode(hsmDecoder);
        } else {
            registerLicenseCommon.EnableJanCodeDisableOCR(hsmDecoder);
        }

        //Check enable scan
        if (!isEnableScan) {
            registerLicenseCommon.DisableScan(hsmDecoder);
        }

        // Declare HSM component UI
        hsmDecoder.enableFlashOnDecode(false);
        hsmDecoder.enableSound(false);
        hsmDecoder.enableAimer(false);
        hsmDecoder.setWindowMode(WindowMode.CENTERING);
        hsmDecoder.setWindow(18, 42, 0, 100);

        // Assign listener
        hsmDecoder.registerPlugin(customPlugin);
    }

    /**
     * onDestroy event handler<br>
     * Unregister all HSM instance
     */
    @Override
    public void onDestroy() {

        super.onDestroy();

        HSMDecoder.disposeInstance();

        //dispose of plug-in
        customPlugin.dispose();
    }

    @Override
    protected void onStop() {

        super.onStop();

        HSMDecoder.disposeInstance();
    }

    /**
     * Show fragment input janCode
     */
    private void showFragmentInputJan() {

        InputJanCodeFragment inputJanCodeFragment = new InputJanCodeFragment();
        FragmentManager fm = getSupportFragmentManager();

        Bundle bundle = new Bundle();
        //put flag switch OCR
        bundle.putString(Constants.FLAG_SWITCH_OCR, flagSwitchOCR);
        inputJanCodeFragment.setArguments(bundle);

        inputJanCodeFragment.show(fm, null);
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
        File root = new File(Environment.getExternalStorageDirectory(), "magazinereturncandidate_log");
        if (!root.exists()) {
            root.mkdirs();
        }

        userID = userID.replaceAll("__", "_");
        serverName = serverName.replaceAll("__", "_");
        shopID = shopID.replaceAll("__", "_");

        //Gzip file name
        String gzipFileName = root + "/" + serverName + "__" + shopID + "__" + userID + "__"
                + strDate + ".log.gz";
        GzipFileCommon gzipFileCommon = new GzipFileCommon();
        gzipFileCommon.compressGzipFile(filepath, gzipFileName);
    }

    @Override
    protected void onResume() {

        super.onResume();

        //Event click item detail
        lvBook.setOnItemClickListener(new AdapterView.OnItemClickListener() {
            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {

                //Disable scan
                registerLicenseCommon.DisableScan(hsmDecoder);

                ProductDetailFragment dProductDetailFragment = new ProductDetailFragment();
                FragmentManager fm = getSupportFragmentManager();
                Bundle bundleItems = new Bundle();
                bundleItems.putString(Constants.COLUMN_JAN_CD, arrBookInlist.get(position)[0]);
                bundleItems.putString(Constants.COLUMN_STOCK_COUNT, arrBookInlist.get(position)[1]);
                bundleItems.putString(Constants.COLUMN_GOODS_NAME, arrBookInlist.get(position)[2]);
                bundleItems.putString(Constants.COLUMN_WRITER_NAME, arrBookInlist.get(position)[3]);
                bundleItems.putString(Constants.COLUMN_PUBLISHER_NAME, arrBookInlist.get(position)[4]);
                bundleItems.putString(Constants.COLUMN_PRICE, arrBookInlist.get(position)[5]);
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
                //put position edit
                bundleItems.putInt(Constants.POSITION_EDIT_PRODUCT, position);
                //put flag switch OCR
                bundleItems.putString(Constants.FLAG_SWITCH_OCR, flagSwitchOCR);

                dProductDetailFragment.setArguments(bundleItems);
                dProductDetailFragment.show(fm, null);
            }
        });

        //Event long press item detail
        lvBook.setOnItemLongClickListener(new AdapterView.OnItemLongClickListener() {

            @Override
            public boolean onItemLongClick(AdapterView<?> parent, View view, final int position, long id) {

                //Disable scan
                registerLicenseCommon.DisableScan(hsmDecoder);

                AlertDialog.Builder alertDialogBuilder = new AlertDialog.Builder(SdmScannerActivity.this);
                alertDialogBuilder.setMessage(String.format(Message.MESSAGE_CONFIRM_DELETE_JAN_CD, arrBookInlist.get(position)[0]));

                alertDialogBuilder.setCancelable(false);
                // Configure alert dialog button
                alertDialogBuilder.setPositiveButton(Message.MESSAGE_SELECT_YES,
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                //Remove items long press
                                if (arrBookInlist.size() > 0) {
                                    LogManagerCommon.i(TAG,
                                            String.format(Message.TAG_SCANNER_ACTIVITY_CANCEL, arrBookInlist.get(position)[0]));
                                    arrBookInlist.remove(position);

                                    // Set data adapter to list view
                                    ListViewScanAdapter adapterBook = new ListViewScanAdapter(SdmScannerActivity.this, arrBookInlist);
                                    lvBook.setAdapter(adapterBook);
                                }
                            }
                        });
                alertDialogBuilder.setNegativeButton(Message.MESSAGE_SELECT_NO,
                        new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                dialog.dismiss();
                            }
                        });
                AlertDialog alert = alertDialogBuilder.show();
                TextView messageText = (TextView) alert.findViewById(android.R.id.message);
                assert messageText != null;
                messageText.setGravity(Gravity.CENTER);

                //Enable scan
                isEnableScan = true;
                registerLicenseCommon.EnableOCROrJanCode(flagSwitchOCR, hsmDecoder);
                return true;
            }
        });
    }

    /**
     * Function call when click Product detail and edit number return magazine
     */
    @Override
    public void onDataEdit(int positionEdit, String valueEdit) {

        LogManagerCommon.i(TAG,
                String.format(Message.TAG_SCANNER_ACTIVITY_EDIT, arrBookInlist.get(positionEdit)[positionEdit]));
        String[] item = new String[]{arrBookInlist.get(positionEdit)[0], valueEdit, arrBookInlist.get(positionEdit)[2],
                arrBookInlist.get(positionEdit)[3], arrBookInlist.get(positionEdit)[4], arrBookInlist.get(positionEdit)[5],
                arrBookInlist.get(positionEdit)[6], arrBookInlist.get(positionEdit)[7],
                arrBookInlist.get(positionEdit)[8], arrBookInlist.get(positionEdit)[9],
                arrBookInlist.get(positionEdit)[10], arrBookInlist.get(positionEdit)[11],
                arrBookInlist.get(positionEdit)[12], arrBookInlist.get(positionEdit)[13],
                arrBookInlist.get(positionEdit)[14], arrBookInlist.get(positionEdit)[15],
                arrBookInlist.get(positionEdit)[16], arrBookInlist.get(positionEdit)[17],
                arrBookInlist.get(positionEdit)[18], arrBookInlist.get(positionEdit)[19],
                arrBookInlist.get(positionEdit)[20], arrBookInlist.get(positionEdit)[21],
                arrBookInlist.get(positionEdit)[22], arrBookInlist.get(positionEdit)[23]};
        arrBookInlist.set(positionEdit, item);

        // Set data adapter to list view
        ListViewScanAdapter adapterBook = new ListViewScanAdapter(SdmScannerActivity.this, arrBookInlist);
        lvBook.setAdapter(adapterBook);
    }

    /**
     * Function call when click submit input JanCode
     */
    @Override
    public void onDataInput(String valueJanCode, String valueNumberReturn) {

        int lenJanCodeInput = valueJanCode.length();
        String janCode18Character = valueJanCode.substring(0, lenJanCodeInput - 5);
        //Get info
        ReturnMagazineEntity magazineEntity;
        if (lenJanCodeInput == Constants.JAN_18_CHAR) {
            magazineEntity = returnMagazineModel.getItemReturnMagazine(janCode18Character);
        } else {
            magazineEntity = returnMagazineModel.getItemReturnMagazine(valueJanCode);
        }

        //Set default
        String strPrice = Constants.BLANK;
        if (magazineEntity.getBqgm_price() != null) {
            strPrice = magazineEntity.getBqgm_price().toString();
        }

        checkIndexListInputJan(valueJanCode);

        //Check jan_code in list and remove
        String JanCodeResult;
        if (lenJanCodeInput == Constants.JAN_18_CHAR) {
            JanCodeResult = valueJanCode;
        } else {
            if (janCodeNew18Character != null) {
                JanCodeResult = janCodeNew18Character;
            } else {
                JanCodeResult = valueJanCode;
            }
        }
        //Reset janCode
        janCodeNew18Character = null;


        String[] item = new String[]{JanCodeResult, valueNumberReturn, magazineEntity.getBqgm_goods_name(),
                magazineEntity.getBqgm_writer_name(), magazineEntity.getBqgm_publisher_name(), strPrice,
                magazineEntity.getBqtse_first_supply_date(), magazineEntity.getBqtse_last_supply_date(),
                magazineEntity.getBqtse_last_sale_date(), magazineEntity.getBqtse_last_order_date(),
                magazineEntity.getBqct_media_group1_cd(), magazineEntity.getBqct_media_group1_name(),
                magazineEntity.getBqct_media_group2_cd(), magazineEntity.getBqct_media_group2_name(),
                magazineEntity.getBqgm_sales_date(), magazineEntity.getBqgm_publisher_cd(), magazineEntity.getBqio_trn_date(),
                Constants.FLAG_0, String.valueOf(magazineEntity.getYear_rank()), String.valueOf(magazineEntity.getJoubi()),
                String.valueOf(magazineEntity.getSts_total_sales()), String.valueOf(magazineEntity.getSts_total_supply()),
                String.valueOf(magazineEntity.getSts_total_return()), magazineEntity.getLocation_id()};


        if (!JanCodeResult.equals(magazineEntity.getJan_cd())) {
            noReturnSound.start();
            LogManagerCommon.i(TAG, String.format(Message.TAG_SCANNER_ACTIVITY_OUT_LIST, valueJanCode));
        } else {
            normalSound.start();
            item[17] = Constants.FLAG_1;
            LogManagerCommon.i(TAG, String.format(Message.TAG_SCANNER_ACTIVITY_INLIST, valueJanCode,
                    magazineEntity.getBqgm_goods_name()));
        }

        arrBookInlist.add(0, item);

        // Update list view
        ListViewScanAdapter adapterBook = new ListViewScanAdapter(this, arrBookInlist);
        lvBook.setAdapter(adapterBook);

        //Enable scan
        isEnableScan = true;
        registerLicenseCommon.EnableOCROrJanCode(flagSwitchOCR, hsmDecoder);
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
            //Type location = 10 is send file
            if (typeLocation == 10) {
                csvFileCommon.deleteCSVFile(fileName);
            } else {
                delete(fileName);
            }
        } else {
            countFile++;
        }
        //typeLocation == 0 is button logout
        if (countFile == files.length && typeLocation == 0) {
            clearAndLogout();
        } else {
            //Event when click send data ok
            csvFileCommon.deleteCSVFile(fileName);
            // Stop process loading screen
            progress.dismiss();
        }
        //Enable scan
        isEnableScan = true;
        registerLicenseCommon.EnableOCROrJanCode(flagSwitchOCR, hsmDecoder);
    }
}
