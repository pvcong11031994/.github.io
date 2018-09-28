package com.android.returncandidate.common.utils;

import android.os.Bundle;

import com.android.returncandidate.common.constants.Constants;
import com.android.returncandidate.db.entity.Books;
import com.android.returncandidate.db.entity.CLP;
import com.android.returncandidate.db.models.ReturnbookModel;

import java.util.ArrayList;
import java.util.Calendar;
import java.util.HashMap;
import java.util.List;

/**
 * Created by cong-pv on
 * 2018/07/11.
 */

public class Common {


    //Put data into activity
    public Bundle DataPutActivity(FlagSettingNew flagSettingNew, FlagSettingOld flagSettingOld) {

        Bundle bundle = new Bundle();
        //Flag New
        //put flag classification 分類
        bundle.putStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP1_CD, flagSettingNew.getFlagClassificationGroup1Cd());
        bundle.putStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP1_NAME, flagSettingNew.getFlagClassificationGroup1Name());
        //bundle.put
        //bundle.putSerializable(Constants.FLAG_CLASSIFICATION_GROUP2_CD, flagSettingNew.getFlagClassificationGroup2Cd());
        //bundle.putSerializable(Constants.FLAG_CLASSIFICATION_GROUP2_NAME, flagSettingNew.getFlagClassificationGroup2Name());
        bundle.putStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP2_CD, flagSettingNew.getFlagClassificationGroup2Cd());
        bundle.putStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP2_NAME, flagSettingNew.getFlagClassificationGroup2Name());
        //put flag publisher 出版社
        bundle.putStringArrayList(Constants.FLAG_PUBLISHER, flagSettingNew.getFlagPublisher());
        bundle.putStringArrayList(Constants.FLAG_PUBLISHER_SHOW_SCREEN, flagSettingNew.getFlagPublisherShowScreen());
        //put flag release date 発売日
        bundle.putString(Constants.FLAG_RELEASE_DATE, flagSettingNew.getFlagReleaseDate());
        bundle.putString(Constants.FLAG_RELEASE_DATE_SHOW_SCREEN, flagSettingNew.getFlagReleaseDateShowScreen());
        //put flag undisturbed 未動期間
        bundle.putString(Constants.FLAG_UNDISTURBED, flagSettingNew.getFlagUndisturbed());
        bundle.putString(Constants.FLAG_UNDISTURBED_SHOW_SCREEN, flagSettingNew.getFlagUndisturbedShowScreen());
        //put flag number of stocks 在庫数
        bundle.putString(Constants.FLAG_NUMBER_OF_STOCKS, flagSettingNew.getFlagNumberOfStocks());
        bundle.putString(Constants.FLAG_NUMBER_OF_STOCKS_SHOW_SCREEN, flagSettingNew.getFlagNumberOfStocksShowScreen());
        //put flag stocks percent 在庫％
        bundle.putString(Constants.FLAG_STOCKS_PERCENT, flagSettingNew.getFlagStockPercent());
        bundle.putString(Constants.FLAG_STOCKS_PERCENT_SHOW_SCREEN, flagSettingNew.getFlagStockPercentShowScreen());
        //put flag old back group 1 cd
        bundle.putStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP1_CD_BACK, flagSettingNew.getFlagClassificationGroup1Cd());
        bundle.putStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP1_NAME_BACK, flagSettingNew.getFlagClassificationGroup1Name());
        //put flag joubi 対象外
        bundle.putString(Constants.FLAG_JOUBI, flagSettingNew.getFlagJoubi());

        //Flag Old
        bundle.putStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP1_CD_OLD, flagSettingOld.getFlagClassificationGroup1Cd());
        bundle.putStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP1_NAME_OLD, flagSettingOld.getFlagClassificationGroup1Name());
        bundle.putStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP2_CD_OLD, flagSettingOld.getFlagClassificationGroup2Cd());
        bundle.putStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP2_NAME_OLD, flagSettingOld.getFlagClassificationGroup2Name());
        //put flag publisher 出版社
        bundle.putStringArrayList(Constants.FLAG_PUBLISHER_OLD, flagSettingOld.getFlagPublisher());
        bundle.putStringArrayList(Constants.FLAG_PUBLISHER_SHOW_SCREEN_OLD, flagSettingOld.getFlagPublisherShowScreen());
        //put flag release date 発売日
        bundle.putString(Constants.FLAG_RELEASE_DATE_OLD, flagSettingOld.getFlagReleaseDate());
        bundle.putString(Constants.FLAG_RELEASE_DATE_SHOW_SCREEN_OLD, flagSettingOld.getFlagReleaseDateShowScreen());
        //put flag undisturbed 未動期間
        bundle.putString(Constants.FLAG_UNDISTURBED_OLD, flagSettingOld.getFlagUndisturbed());
        bundle.putString(Constants.FLAG_UNDISTURBED_SHOW_SCREEN_OLD, flagSettingOld.getFlagUndisturbedShowScreen());
        //put flag number of stocks 在庫数
        bundle.putString(Constants.FLAG_NUMBER_OF_STOCKS_OLD, flagSettingOld.getFlagNumberOfStocks());
        bundle.putString(Constants.FLAG_NUMBER_OF_STOCKS_SHOW_SCREEN_OLD, flagSettingOld.getFlagNumberOfStocksShowScreen());
        //put flag stocks percent 在庫％
        bundle.putString(Constants.FLAG_STOCKS_PERCENT_OLD, flagSettingOld.getFlagStockPercent());
        bundle.putString(Constants.FLAG_STOCKS_PERCENT_SHOW_SCREEN_OLD, flagSettingOld.getFlagStockPercentShowScreen());

        //put flag joubi 対象外
        bundle.putString(Constants.FLAG_JOUBI_OLD, flagSettingOld.getFlagJoubi());
        //return bundle
        return bundle;
    }

    //Set flag to arguments
    public void SetArgumentsFlagData(FlagSettingNew flagSettingNew, FlagSettingOld flagSettingOld, Bundle bundleArguments) {

        //New
        //save flag classification
        flagSettingNew.setFlagClassificationGroup1Cd(bundleArguments.getStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP1_CD));
        flagSettingNew.setFlagClassificationGroup1Name(bundleArguments.getStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP1_NAME));
        flagSettingNew.setFlagClassificationGroup2Cd(bundleArguments.getStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP2_CD));
        flagSettingNew.setFlagClassificationGroup2Name(bundleArguments.getStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP2_NAME));
        //flagSettingNew.setFlagClassificationGroup2Cd((HashMap<String,ArrayList<String>>) bundleArguments.getSerializable(Constants.FLAG_CLASSIFICATION_GROUP2_CD));
        //flagSettingNew.setFlagClassificationGroup2Name((HashMap<String,ArrayList<String>>) bundleArguments.getSerializable(Constants.FLAG_CLASSIFICATION_GROUP2_NAME));
        //save flag publisher
        flagSettingNew.setFlagPublisher(bundleArguments.getStringArrayList(Constants.FLAG_PUBLISHER));
        flagSettingNew.setFlagPublisherShowScreen(bundleArguments.getStringArrayList(Constants.FLAG_PUBLISHER_SHOW_SCREEN));
        //save flag release date
        flagSettingNew.setFlagReleaseDate(bundleArguments.getString(Constants.FLAG_RELEASE_DATE));
        flagSettingNew.setFlagReleaseDateShowScreen(bundleArguments.getString(Constants.FLAG_RELEASE_DATE_SHOW_SCREEN));
        //save flag undisturbed
        flagSettingNew.setFlagUndisturbed(bundleArguments.getString(Constants.FLAG_UNDISTURBED));
        flagSettingNew.setFlagUndisturbedShowScreen(bundleArguments.getString(Constants.FLAG_UNDISTURBED_SHOW_SCREEN));
        //save flag number of stocks
        flagSettingNew.setFlagNumberOfStocks(bundleArguments.getString(Constants.FLAG_NUMBER_OF_STOCKS));
        flagSettingNew.setFlagNumberOfStocksShowScreen(bundleArguments.getString(Constants.FLAG_NUMBER_OF_STOCKS_SHOW_SCREEN));
        //save flag stocks percent
        flagSettingNew.setFlagStockPercent(bundleArguments.getString(Constants.FLAG_STOCKS_PERCENT));
        flagSettingNew.setFlagStockPercentShowScreen(bundleArguments.getString(Constants.FLAG_STOCKS_PERCENT_SHOW_SCREEN));
        //save flag joubi
        flagSettingNew.setFlagJoubi(bundleArguments.getString(Constants.FLAG_JOUBI));
        //save flag click setting
        flagSettingNew.setFlagClickSetting(bundleArguments.getString(Constants.FLAG_CLICK_SETTING));

        //save flag old back group 1
        flagSettingNew.setFlagClassificationGroup1CdOld(bundleArguments.getStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP1_CD_BACK));
        flagSettingNew.setFlagClassificationGroup1NameOld(bundleArguments.getStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP1_NAME_BACK));

        //Old
        //save flag classification
        flagSettingOld.setFlagClassificationGroup1Cd(bundleArguments.getStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP1_CD_OLD));
        flagSettingOld.setFlagClassificationGroup1Name(bundleArguments.getStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP1_NAME_OLD));
        flagSettingOld.setFlagClassificationGroup2Cd(bundleArguments.getStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP2_CD_OLD));
        flagSettingOld.setFlagClassificationGroup2Name(bundleArguments.getStringArrayList(Constants.FLAG_CLASSIFICATION_GROUP2_NAME_OLD));
        //save flag publisher
        flagSettingOld.setFlagPublisher(bundleArguments.getStringArrayList(Constants.FLAG_PUBLISHER_OLD));
        flagSettingOld.setFlagPublisherShowScreen(bundleArguments.getStringArrayList(Constants.FLAG_PUBLISHER_SHOW_SCREEN_OLD));
        //save flag release date
        flagSettingOld.setFlagReleaseDate(bundleArguments.getString(Constants.FLAG_RELEASE_DATE_OLD));
        flagSettingOld.setFlagReleaseDateShowScreen(bundleArguments.getString(Constants.FLAG_RELEASE_DATE_SHOW_SCREEN_OLD));
        //save flag undisturbed
        flagSettingOld.setFlagUndisturbed(bundleArguments.getString(Constants.FLAG_UNDISTURBED_OLD));
        flagSettingOld.setFlagUndisturbedShowScreen(bundleArguments.getString(Constants.FLAG_UNDISTURBED_SHOW_SCREEN_OLD));
        //save flag number of stocks
        flagSettingOld.setFlagNumberOfStocks(bundleArguments.getString(Constants.FLAG_NUMBER_OF_STOCKS_OLD));
        flagSettingOld.setFlagNumberOfStocksShowScreen(bundleArguments.getString(Constants.FLAG_NUMBER_OF_STOCKS_SHOW_SCREEN_OLD));
        //save flag stocks percent
        flagSettingOld.setFlagStockPercent(bundleArguments.getString(Constants.FLAG_STOCKS_PERCENT_OLD));
        flagSettingOld.setFlagStockPercentShowScreen(bundleArguments.getString(Constants.FLAG_STOCKS_PERCENT_SHOW_SCREEN_OLD));
        //save flag joubi
        flagSettingOld.setFlagJoubi(bundleArguments.getString(Constants.FLAG_JOUBI_OLD));
    }


    public String FormatDateTime(String strReleaseDateInput) {

        int intReleaseDateInput = Integer.parseInt(strReleaseDateInput);

        String day, month, year;
        Calendar calendarNow = Calendar.getInstance();

        //Get Month Current
        int monthCurrent = (calendarNow.get(Calendar.MONTH)) + 1;

        //Get Day Current
        if ((calendarNow.get(Calendar.DATE)) < 10) {
            day = "0" + String.valueOf(calendarNow.get(Calendar.DATE));
        } else {
            day = String.valueOf(calendarNow.get(Calendar.DATE));
        }

        //Check month input > month current
        if (intReleaseDateInput == monthCurrent) {
            year = String.valueOf(calendarNow.get(Calendar.YEAR) - 1);
            month = "12";
        } else if (intReleaseDateInput < monthCurrent) {
            year = String.valueOf(calendarNow.get(Calendar.YEAR));
            if ((monthCurrent - intReleaseDateInput) < 10) {
                month = "0" + String.valueOf(monthCurrent - intReleaseDateInput);
            } else {
                month = String.valueOf(monthCurrent - intReleaseDateInput);
            }
        } else {
            int a = intReleaseDateInput / monthCurrent;
            int b = intReleaseDateInput % monthCurrent;
            year = String.valueOf(calendarNow.get(Calendar.YEAR) - a);
            if ((12 - b) < 10) {
                month = "0" + String.valueOf(12 - b);
            } else {
                month = String.valueOf(12 - b);
            }
        }
        return year + month + day;
    }

    public String FormatPercent(String strPercentInput) {

        int intPercentInput = Integer.parseInt(strPercentInput);
        return String.valueOf(1 - intPercentInput * 0.05);
    }

    public String FormatPercentLocal(String strPercentInput) {

        int intPercentInput = Integer.parseInt(strPercentInput);
        return String.valueOf(intPercentInput * 0.05);
    }

    public int ConvertStringDateToInt(String strDateTime) {

        if (Constants.BLANK.equals(strDateTime) || Constants.NULL.equals(strDateTime)) {
            return Constants.VALUE_DEFAULT_DATE_INT;
        } else {
            return Integer.parseInt(strDateTime);
        }
    }

    public String GenerateJAN(String codeISBN) {
        String result = codeISBN.substring(0, codeISBN.length() - 1);
        return result + GetCheckDigit(codeISBN);
    }


    private String GetCheckDigit(String code) {

        int odd = 0;
        int even = 0;
        for (int i = 0; i < code.length() - 1; i++) {
            int num = Integer.parseInt(code.substring(i, i + 1));
            if (i % 2 == 0) {
                even += num;
            } else {
                odd += num;
            }
        }
        return String.valueOf(((10 - (odd * 3 + even) % 10) % 10));
    }

    //Check list group2 cd
    public HashMap<String, ArrayList<String>> ListGroup2Cd(ArrayList<String> listGroup2CdNew, ArrayList<String> listGroup2CdOld,
                                                           ArrayList<String> listGroup2NameOld) {

        Boolean check;
        ArrayList<String> resultCd = new ArrayList<>();
        ArrayList<String> resultName = new ArrayList<>();
        for (int i = 0; i < listGroup2CdOld.size(); i++) {
            check = false;
            for (int j = 0; j < listGroup2CdNew.size(); j++) {
                if (listGroup2CdNew.get(j).equals(listGroup2CdOld.get(i))) {
                    check = true;
                    break;
                }
            }
            if (check) {
                resultCd.add(listGroup2CdOld.get(i));
                resultName.add(listGroup2NameOld.get(i));
            }
        }
        HashMap<String, ArrayList<String>> hashMap = new HashMap<>();
        hashMap.put(Constants.COLUMN_MEDIA_GROUP2_CD, resultCd);
        hashMap.put(Constants.COLUMN_MEDIA_GROUP2_NAME, resultName);
        return hashMap;
    }

    public HashMap<String, ArrayList<String>> ListGroup2CdNew(FlagSettingNew flagSettingNew, List<CLP> rlist) {

        Boolean check;
        ArrayList<String> resultCd = new ArrayList<>();
        ArrayList<String> resultName = new ArrayList<>();
        for (int i = 0; i < flagSettingNew.getFlagClassificationGroup2Cd().size(); i++) {
            check = false;
            for (int j = 0; j < rlist.size(); j++) {
                if (rlist.get(j).getId().equals(flagSettingNew.getFlagClassificationGroup2Cd().get(i))) {
                    check = true;
                    break;
                }
            }
            if (!check) {
                resultCd.add(flagSettingNew.getFlagClassificationGroup2Cd().get(i));
                resultName.add(flagSettingNew.getFlagClassificationGroup2Name().get(i));
            }
        }
        HashMap<String, ArrayList<String>> hashMap = new HashMap<>();
        hashMap.put(Constants.COLUMN_MEDIA_GROUP2_CD, resultCd);
        hashMap.put(Constants.COLUMN_MEDIA_GROUP2_NAME, resultName);
        return hashMap;
    }

}
