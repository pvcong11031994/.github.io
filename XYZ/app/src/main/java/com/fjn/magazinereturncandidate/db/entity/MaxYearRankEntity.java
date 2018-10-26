package com.fjn.magazinereturncandidate.db.entity;

import com.google.gson.annotations.Expose;
import com.google.gson.annotations.SerializedName;

import java.io.Serializable;

public class MaxYearRankEntity implements Serializable {

    public int getMaxYearRank() {
        return maxYearRank;
    }

    public void setMaxYearRank(int maxYearRank) {
        this.maxYearRank = maxYearRank;
    }

    /**
     * max_year_rank
     */
    @SerializedName("max_year_rank")
    @Expose
    private int maxYearRank;

    public MaxYearRankEntity() {

    }
}
