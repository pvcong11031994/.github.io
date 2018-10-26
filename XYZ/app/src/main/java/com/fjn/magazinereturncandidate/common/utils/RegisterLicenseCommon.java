package com.fjn.magazinereturncandidate.common.utils;

import com.fjn.magazinereturncandidate.common.constants.Constants;
import com.honeywell.barcode.HSMDecoder;
import com.honeywell.barcode.OCRActiveTemplate;
import com.honeywell.barcode.Symbology;

/**
 * Function enable and disable key
 * Created by cong-pv on 2018/10/22.
 */

public class RegisterLicenseCommon {

    public String EnableOCRDisableJanCode(HSMDecoder hsmDecoder) {

        String flagSwitchOCR = Constants.FLAG_1;
        hsmDecoder.enableSymbology(Symbology.OCR);
        hsmDecoder.setOCRActiveTemplate(OCRActiveTemplate.ISBN);
        hsmDecoder.disableSymbology(Symbology.EAN13);
        hsmDecoder.disableSymbology(Symbology.CODE128);
        hsmDecoder.disableSymbology(Symbology.EAN13_5CHAR_ADDENDA);
        return flagSwitchOCR;
    }

    public String EnableJanCodeDisableOCR(HSMDecoder hsmDecoder) {

        String flagSwitchOCR = Constants.FLAG_0;
        hsmDecoder.enableSymbology(Symbology.EAN13);
        hsmDecoder.enableSymbology(Symbology.CODE128);
        hsmDecoder.enableSymbology(Symbology.EAN13_5CHAR_ADDENDA);
        hsmDecoder.disableSymbology(Symbology.OCR);
        return flagSwitchOCR;
    }

    public void DisableScan(HSMDecoder hsmDecoder) {

        hsmDecoder.disableSymbology(Symbology.EAN13);
        hsmDecoder.disableSymbology(Symbology.CODE128);
        hsmDecoder.disableSymbology(Symbology.OCR);
        hsmDecoder.disableSymbology(Symbology.EAN13_5CHAR_ADDENDA);
    }

    public void EnableOCROrJanCode(String flagSwitchOCR, HSMDecoder hsmDecoder) {

        if (Constants.FLAG_1.equals(flagSwitchOCR)) {
            hsmDecoder.enableSymbology(Symbology.OCR);
            hsmDecoder.setOCRActiveTemplate(OCRActiveTemplate.ISBN);
        } else {
            hsmDecoder.enableSymbology(Symbology.EAN13);
            hsmDecoder.enableSymbology(Symbology.CODE128);
            hsmDecoder.enableSymbology(Symbology.EAN13_5CHAR_ADDENDA);
        }
    }
}
