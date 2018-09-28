package com.android.returncandidate.common.utils;

import android.os.Bundle;

import com.android.returncandidate.common.constants.Constants;

import java.util.ArrayList;

/**
 * Created by cong-pv
 * on 2018/07/10.
 */

public class FlagSettingOld {

    public FlagSettingOld() {
    }

    public String getFlagReleaseDate() {
        return flagReleaseDate;
    }

    public void setFlagReleaseDate(String flagReleaseDate) {
        this.flagReleaseDate = flagReleaseDate;
    }

    public String getFlagUndisturbed() {
        return flagUndisturbed;
    }

    public void setFlagUndisturbed(String flagUndisturbed) {
        this.flagUndisturbed = flagUndisturbed;
    }

    public String getFlagNumberOfStocks() {
        return flagNumberOfStocks;
    }

    public void setFlagNumberOfStocks(String flagNumberOfStocks) {
        this.flagNumberOfStocks = flagNumberOfStocks;
    }

    public String getFlagStockPercent() {
        return flagStockPercent;
    }

    public void setFlagStockPercent(String flagStockPercent) {
        this.flagStockPercent = flagStockPercent;
    }

    public String getFlagReleaseDateShowScreen() {
        return flagReleaseDateShowScreen;
    }

    public void setFlagReleaseDateShowScreen(String flagReleaseDateShowScreen) {
        this.flagReleaseDateShowScreen = flagReleaseDateShowScreen;
    }

    public String getFlagUndisturbedShowScreen() {
        return flagUndisturbedShowScreen;
    }

    public void setFlagUndisturbedShowScreen(String flagUndisturbedShowScreen) {
        this.flagUndisturbedShowScreen = flagUndisturbedShowScreen;
    }

    public String getFlagNumberOfStocksShowScreen() {
        return flagNumberOfStocksShowScreen;
    }

    public void setFlagNumberOfStocksShowScreen(String flagNumberOfStocksShowScreen) {
        this.flagNumberOfStocksShowScreen = flagNumberOfStocksShowScreen;
    }

    public String getFlagStockPercentShowScreen() {
        return flagStockPercentShowScreen;
    }

    public void setFlagStockPercentShowScreen(String flagStockPercentShowScreen) {
        this.flagStockPercentShowScreen = flagStockPercentShowScreen;
    }

    public String getFlagJoubi() {
        return flagJoubi;
    }

    public void setFlagJoubi(String flagJoubi) {
        this.flagJoubi = flagJoubi;
    }

    public ArrayList<String> getFlagClassificationGroup1Cd() {
        return flagClassificationGroup1Cd;
    }

    public void setFlagClassificationGroup1Cd(ArrayList<String> flagClassificationGroup1Cd) {
        this.flagClassificationGroup1Cd = flagClassificationGroup1Cd;
    }

    public ArrayList<String> getFlagClassificationGroup1Name() {
        return flagClassificationGroup1Name;
    }

    public void setFlagClassificationGroup1Name(ArrayList<String> flagClassificationGroup1Name) {
        this.flagClassificationGroup1Name = flagClassificationGroup1Name;
    }

    public ArrayList<String> getFlagClassificationGroup2Cd() {
        return flagClassificationGroup2Cd;
    }

    public void setFlagClassificationGroup2Cd(ArrayList<String> flagClassificationGroup2Cd) {
        this.flagClassificationGroup2Cd = flagClassificationGroup2Cd;
    }

    public ArrayList<String> getFlagClassificationGroup2Name() {
        return flagClassificationGroup2Name;
    }

    public void setFlagClassificationGroup2Name(ArrayList<String> flagClassificationGroup2Name) {
        this.flagClassificationGroup2Name = flagClassificationGroup2Name;
    }

    public ArrayList<String> getFlagPublisher() {
        return flagPublisher;
    }

    public void setFlagPublisher(ArrayList<String> flagPublisher) {
        this.flagPublisher = flagPublisher;
    }

    public ArrayList<String> getFlagPublisherShowScreen() {
        return flagPublisherShowScreen;
    }

    public void setFlagPublisherShowScreen(ArrayList<String> flagPublisherShowScreen) {
        this.flagPublisherShowScreen = flagPublisherShowScreen;
    }

    //Flag selected classification
    private ArrayList<String> flagClassificationGroup1Cd;
    private ArrayList<String> flagClassificationGroup1Name;

    private ArrayList<String> flagClassificationGroup2Cd;
    private ArrayList<String> flagClassificationGroup2Name;

    //Flag selected publisher
    private ArrayList<String> flagPublisher;
    private ArrayList<String> flagPublisherShowScreen;

    //Flag selected release date
    private String flagReleaseDate;
    private String flagReleaseDateShowScreen;

    //Flag Undisturbed
    private String flagUndisturbed;
    private String flagUndisturbedShowScreen;

    //Flag Number of stocks
    private String flagNumberOfStocks;
    private String flagNumberOfStocksShowScreen;

    //Flag Stocks Percent
    private String flagStockPercent;
    private String flagStockPercentShowScreen;

    //Flag joubi
    private String flagJoubi;
}
