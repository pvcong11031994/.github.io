package com.fjn.magazinereturncandidate.common.utils;

import android.widget.TextView;

import com.fjn.magazinereturncandidate.common.constants.Constants;
import com.fjn.magazinereturncandidate.common.constants.Message;

import static com.fjn.magazinereturncandidate.common.constants.Constants.JAN_12_CHAR;
import static com.fjn.magazinereturncandidate.common.constants.Constants.JAN_13_CHAR;
import static com.fjn.magazinereturncandidate.common.constants.Constants.JAN_18_CHAR;

/**
 * Common check input ......
 * Created by cong-pv on 2018/10/20.
 */

public class CheckDataCommon {

    /**
     * Function check validate fields null when input janCode
     *
     * @param tv textView show message error
     * @return
     */
    public Boolean validateFieldsNotNull(TextView tv) {

        if (tv.getText().length() == 0) {
            tv.setError(Message.MESSAGE_ERROR_BLANK_INPUT);
            return false;
        }
        return true;
    }

    /**
     * Function check validate fields input janCode
     *
     * @param tv textView show message error
     * @return
     */
    public boolean validateFields(TextView tv) {

        if (tv.getText().length() != JAN_12_CHAR &&
                tv.getText().length() != JAN_13_CHAR &&
                tv.getText().length() != JAN_18_CHAR) {
            tv.setError(Message.MESSAGE_ERROR_INPUT_JANCODE);
            return false;
        }
        return true;
    }

    /**
     * Function check check digit true/false
     *
     * @param strJanInput is janCode input
     * @param tv          is textView show message error
     * @return
     */
    public String validateCheckDigit(String strJanInput, TextView tv) {

        String strResult;

        int lenStrJanInput = strJanInput.length();

        if (lenStrJanInput == Constants.JAN_12_CHAR) {
            return strJanInput + GetCheckDigit12Character(strJanInput);
        } else if (lenStrJanInput == Constants.JAN_13_CHAR) {
            strResult = GenerateJAN(strJanInput);
            if (strJanInput.equals(strResult)) {
                return strJanInput;
            } else {
                tv.setError(Message.MESSAGE_ERROR_CHECK_DIGIT_JANCODE);
                return null;
            }
        } else {
            String valueJan13Character = strJanInput.substring(0, lenStrJanInput - 5);
            strResult = GenerateJAN(valueJan13Character);
            if (valueJan13Character.equals(strResult)) {
                return strJanInput;
            } else {
                tv.setError(Message.MESSAGE_ERROR_CHECK_DIGIT_JANCODE);
                return null;
            }
        }
    }


    /**
     * Function Format OCR
     *
     * @param janCodeOCR is OCR scan
     * @return
     */
    public String validateOCR(String janCodeOCR) {

        String result;
        result = janCodeOCR.replaceAll("ISBN|-| ", "");
        if (!result.startsWith(Constants.PREFIX_JAN_CODE_978)) {
            result = GenerateJAN(Constants.PREFIX_JAN_CODE_978 + result);
        }
        return result;
    }

    /**
     * Function check digit JanCode 13/18 character
     *
     * @param codeISBN is param JanCode 13/18 character
     * @return
     */
    private String GenerateJAN(String codeISBN) {
        String result = codeISBN.substring(0, codeISBN.length() - 1);
        return result + GetCheckDigit(codeISBN);
    }

    /**
     * Function check digit JanCode 13/18 character
     *
     * @param code is param JanCode 13/18 character
     * @return
     */
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


    /**
     * Function check digit JanCode 12 character
     *
     * @param code is param JanCode 12 character
     * @return
     */
    private String GetCheckDigit12Character(String code) {

        int odd = 0;
        int even = 0;
        for (int i = 0; i < code.length(); i++) {
            int num = Integer.parseInt(code.substring(i, i + 1));
            if (i % 2 == 0) {
                even += num;
            } else {
                odd += num;
            }
        }
        return String.valueOf(((10 - (odd * 3 + even) % 10) % 10));
    }
}

