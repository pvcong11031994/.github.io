package com.android.productchange.common.utils;

import com.android.productchange.common.constants.Constants;

/**
 * Class common condition filter
 * Created by cong-pv on 2018/08/31.
 */

public class ConditionQueryCommon {

    private Common common = new Common();

    /**
     * Condition select filter other
     *
     * @param flagSettingNew
     * @return
     */
    public String conditionFilterSetting(FlagSettingNew flagSettingNew) {

        String queryCondition = "";

        //Format Date time
        String strReleaseDate = common.FormatDateTime(flagSettingNew.getFlagReleaseDate());
        String strUndisturbed = common.FormatDateTime(flagSettingNew.getFlagUndisturbed());

        queryCondition += String.format(" WHERE bqgm_sales_date < '%s' ", strReleaseDate);
        //Condition undisturbed
        queryCondition += String.format(" AND bqio_trn_date <= '%s' AND bqtse_last_sale_date <= '%s' " +
                        "AND bqtse_last_supply_date <= '%s' ", strUndisturbed,
                strUndisturbed, strUndisturbed);
        queryCondition += String.format(" AND bqsc_stock_count >= %s ", Integer.parseInt(flagSettingNew.getFlagNumberOfStocks()));
        //Condition

        if (flagSettingNew.getFlagPublisherShowScreen().size() > 0 &&
                !Constants.ROW_ALL.equals(flagSettingNew.getFlagPublisherShowScreen().get(0))) {
            String strPublisher = "";
            for (int i = 0; i < flagSettingNew.getFlagPublisherShowScreen().size(); i++) {
                if (i > 0) {
                    strPublisher += ",";
                }
                strPublisher += "'";
                strPublisher += flagSettingNew.getFlagPublisherShowScreen().get(i);
                strPublisher += "'";
            }
            queryCondition += String.format(" AND publisher_name IN (%s) ", strPublisher);
        }
        if (Constants.VALUE_YES_STANDING.equals(flagSettingNew.getFlagJoubi())) {
            queryCondition += String.format(" AND joubi != %s", Constants.VALUE_JOUBI);
        }

        return queryCondition;
    }


    /**
     * Return condition select group1cd/group2cd
     *
     * @param flagSettingNew
     * @return
     */
    public String conditionFilterSettingGroupCd(FlagSettingNew flagSettingNew) {

        String queryConditionGroupCd = "";
        if (flagSettingNew.getFlagClassificationGroup1Cd().size() >= 1) {
            if (Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup1Cd().get(0))) {
                return queryConditionGroupCd;
            } else {
                if (flagSettingNew.getFlagClassificationGroup2Cd().size() == 1) {
                    if (Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup2Cd().get(0))) {
                        return queryConditionGroupCd;
                    } else {
                        queryConditionGroupCd += String.format(" AND bqct_media_group2_cd = '%s' ", flagSettingNew.getFlagClassificationGroup2Cd().get(0));
                    }
                } else if (flagSettingNew.getFlagClassificationGroup2Cd().size() > 1) {
                    String strCondition = "";
                    for (int i = 0; i < flagSettingNew.getFlagClassificationGroup2Cd().size(); i++) {
                        if (i > 0) {
                            strCondition += ",";
                        }
                        strCondition += "'";
                        strCondition += flagSettingNew.getFlagClassificationGroup2Cd().get(i);
                        strCondition += "'";

                    }
                    queryConditionGroupCd += String.format(" AND bqct_media_group2_cd IN (%s) ", strCondition);
                }
            }
        }
        return queryConditionGroupCd;
    }
}
