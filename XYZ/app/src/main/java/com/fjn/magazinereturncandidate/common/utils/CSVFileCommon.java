package com.fjn.magazinereturncandidate.common.utils;


import android.annotation.SuppressLint;
import android.os.Build;
import android.os.Environment;

import com.fjn.magazinereturncandidate.api.Config;
import com.fjn.magazinereturncandidate.common.constants.Constants;
import com.opencsv.CSVWriter;

import java.io.File;
import java.io.FileWriter;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.List;


/**
 * Save file to folder android
 * Created by cong-pv on 2018/10/18.
 */

public class CSVFileCommon {

    /**
     * Function save file CSV to storage android
     *
     * @param user_id     id login
     * @param shop_id     shop login
     * @param server_name server name login
     * @param dataScan    list data scan
     */
    public void saveCSVFile(String user_id, String shop_id, String server_name, List<String[]> dataScan) {

        try {
            //Create folder save data
            File root = new File(Environment.getExternalStorageDirectory(), "/MagazineReturnCandidate/datasend");
            if (!root.exists()) {
                root.mkdirs();
            }

            //Get Info devices
            String myDeviceModel = Build.ID;

            //Filter data import
            List<String[]> dataImport = new ArrayList<>();
            int index = 1;
            //Create file .csv
            @SuppressLint("SimpleDateFormat") SimpleDateFormat dateFormat = new SimpleDateFormat(
                    "yyyyMMddHHmmss");
            Calendar cal = Calendar.getInstance();
            String strDate = dateFormat.format(cal.getTime());

            //Add header
            String[] itemsHeader = new String[]{"Device Name", "User Login", "Shop Cd", "Server Name",
                    "Jan Code", "Product Name", "Number Return", "Date Send"};
            dataImport.add(itemsHeader);

            //Add list items scan
            for (int i = 0; i < dataScan.size(); i++) {
                if (Constants.FLAG_1.equals(dataScan.get(i)[17])) {
                    String[] item = new String[]{myDeviceModel, user_id, shop_id, server_name,
                            dataScan.get(i)[0], dataScan.get(i)[2], dataScan.get(i)[1], strDate};
                    dataImport.add(index, item);
                    index++;
                }
            }

            //Create file csv
            String fileName = server_name + "_" + shop_id + "_" + user_id + "_" + strDate + ".csv";
            String csv = Environment.getExternalStorageDirectory().getAbsolutePath() + Config.FOLDER_ANDROID_SAVE_DATA + fileName;
            CSVWriter writer = new CSVWriter(new FileWriter(csv));

            writer.writeAll(dataImport);

            writer.close();
        } catch (Exception e) {
            e.printStackTrace();
        }

    }

    /**
     * Function delete file CSV in storage android
     *
     * @param fileName fileName delete
     */
    public void deleteCSVFile(String fileName) {

        // Remove old data file for new session user
        File file = new File(fileName);
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
     * Check file CSV is folder storage android
     *
     * @return
     */
    public int isExistFileCSV() {

        File root = new File(Environment.getExternalStorageDirectory(), Config.FOLDER_ANDROID_SAVE_DATA);
        if (!root.exists()) {
            root.mkdirs();
        }
        return root.listFiles().length;
    }

}

