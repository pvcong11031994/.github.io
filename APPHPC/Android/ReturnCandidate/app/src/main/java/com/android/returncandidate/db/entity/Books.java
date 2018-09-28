package com.android.returncandidate.db.entity;

import com.google.gson.annotations.*;

import java.io.*;

/**
 * Books entity
 *
 * @author cong-pv
 * @since 2018-06-20
 */

public class Books implements Serializable {

    /**
     * location id
     */
    @SerializedName("location_id")
    @Expose
    private String location_id;
    /**
     * jan cd
     */
    @SerializedName("jan_cd")
    @Expose
    private String jan_cd;

    /**
     * bqsc_stock_count
     */
    @SerializedName("bqsc_stock_count")
    @Expose
    private int bqsc_stock_count;

    /**
     * bqgm_goods_name
     */
    @SerializedName("bqgm_goods_name")
    @Expose
    private String bqgm_goods_name;

    /**
     * bqgm_writer_name
     */
    @SerializedName("bqgm_writer_name")
    @Expose
    private String bqgm_writer_name;

    /**
     * bqgm_publisher_cd
     */
    @SerializedName("bqgm_publisher_cd")
    @Expose
    private String bqgm_publisher_cd;
    /**
     * bqgm_publisher_name
     */
    @SerializedName("bqgm_publisher_name")
    @Expose
    private String bqgm_publisher_name;

    /**
     * bqgm_price
     */
    @SerializedName("bqgm_price")
    @Expose
    private Float bqgm_price;


    /**
     * bqtse_first_supply_date
     */
    @SerializedName("bqtse_first_supply_date")
    @Expose
    private String bqtse_first_supply_date;


    /**
     * bqtse_last_supply_date
     */
    @SerializedName("bqtse_last_supply_date")
    @Expose
    private String bqtse_last_supply_date;


    /**
     * bqtse_last_sale_date
     */
    @SerializedName("bqtse_last_sale_date")
    @Expose
    private String bqtse_last_sale_date;


    /**
     * bqtse_last_order_date
     */
    @SerializedName("bqtse_last_order_date")
    @Expose
    private String bqtse_last_order_date;

    /**
     * bqct_media_group1_cd
     */
    @SerializedName("bqct_media_group1_cd")
    @Expose
    private String bqct_media_group1_cd;

    /**
     * bqct_media_group1_name
     */
    @SerializedName("bqct_media_group1_name")
    @Expose
    private String bqct_media_group1_name;

    /**
     * bqct_media_group2_cd
     */
    @SerializedName("bqct_media_group2_cd")
    @Expose
    private String bqct_media_group2_cd;

    /**
     * bqct_media_group2_name
     */
    @SerializedName("bqct_media_group2_name")
    @Expose
    private String bqct_media_group2_name;

    /**
     * bqgm_sales_date
     */
    @SerializedName("bqgm_sales_date")
    @Expose
    private String bqgm_sales_date;

    /**
     * bqio_trn_date
     */
    @SerializedName("bqio_trn_date")
    @Expose
    private String bqio_trn_date;

    /**
     * percent
     */
    @SerializedName("percent")
    @Expose
    private Float percent;


    /**
     * flag_sales_return
     */
    @SerializedName("flag_sales")
    @Expose
    private String flag_sales;

    /**
     * sumStocks
     */
    @SerializedName("sumStocks")
    @Expose
    private int sumStocks;

    /**
     * countJan_Cd
     */
    @SerializedName("countJan_Cd")
    @Expose
    private int countJan_Cd;

    /**
     * joubi
     */
    @SerializedName("joubi")
    @Expose
    private int joubi;

    /**
     * sts_total_sales
     */
    @SerializedName("sts_total_sales")
    @Expose
    private int sts_total_sales;


    /**
     * sts_total_supply
     */
    @SerializedName("sts_total_supply")
    @Expose
    private int sts_total_supply;


    /**
     * sts_total_return
     */
    @SerializedName("sts_total_return")
    @Expose
    private int sts_total_return;


    /**
     * year_rank
     */
    @SerializedName("year_rank")
    @Expose
    private int year_rank;
    /**
     * percent_local
     */
    @SerializedName("percent_local")
    @Expose
    private Float percent_local;

    /**
     * Constructor Books
     */
    public Books() {

    }

    /**
     * Constructor Books
     */
    public Books(String jan_cd, int bqsc_stock_count, String bqgm_publisher_cd, String bqtse_last_supply_date,
                 String bqtse_last_sale_date, String bqct_media_group1_cd, String bqct_media_group2_cd,
                 String bqgm_sales_date, String bqio_trn_date, Float percent_local, String flag_sales) {

        this.jan_cd = jan_cd;
        this.bqsc_stock_count = bqsc_stock_count;
        this.bqgm_publisher_cd = bqgm_publisher_cd;
        this.bqtse_last_supply_date = bqtse_last_supply_date;
        this.bqtse_last_sale_date = bqtse_last_sale_date;
        this.bqct_media_group1_cd = bqct_media_group1_cd;
        this.bqct_media_group2_cd = bqct_media_group2_cd;
        this.bqgm_sales_date = bqgm_sales_date;
        this.bqio_trn_date = bqio_trn_date;
        this.percent_local = percent_local;
        this.flag_sales = flag_sales;
    }

    /**
     * Constructor Books
     */
    public Books(String jan_cd, Float percent_local, String flag_sales) {

        this.jan_cd = jan_cd;
        this.percent_local = percent_local;
        this.flag_sales = flag_sales;
    }

    public int getJoubi() {
        return joubi;
    }

    public void setJoubi(int joubi) {
        this.joubi = joubi;
    }

    public int getBqsc_stock_count() {
        return bqsc_stock_count;
    }

    public void setBqsc_stock_count(int bqsc_stock_count) {
        this.bqsc_stock_count = bqsc_stock_count;
    }

    public int getYear_rank() {
        return year_rank;
    }

    public void setYear_rank(int year_rank) {
        this.year_rank = year_rank;
    }

    public Float getPercent() {
        return percent;
    }

    public void setPercent(Float percent) {
        this.percent = percent;
    }

    public Float getPercent_local() {
        return percent_local;
    }

    public void setPercent_local(Float percent_local) {
        this.percent_local = percent_local;
    }

    public String getFlag_sales() {
        return flag_sales;
    }

    public void setFlag_sales(String flag_sales) {
        this.flag_sales = flag_sales;
    }

    public int getSumStocks() {
        return sumStocks;
    }

    public void setSumStocks(int sumStocks) {
        this.sumStocks = sumStocks;
    }

    public int getCountJan_Cd() {
        return countJan_Cd;
    }

    public void setCountJan_Cd(int countJan_Cd) {
        this.countJan_Cd = countJan_Cd;
    }

    public String getJan_cd() {
        return jan_cd;
    }

    public void setJan_cd(String jan_cd) {
        this.jan_cd = jan_cd;
    }


    public String getBqgm_writer_name() {
        return bqgm_writer_name;
    }

    public void setBqgm_writer_name(String bqgm_writer_name) {
        this.bqgm_writer_name = bqgm_writer_name;
    }

    public String getBqgm_publisher_name() {
        return bqgm_publisher_name;
    }

    public void setBqgm_publisher_name(String bqgm_publisher_name) {
        this.bqgm_publisher_name = bqgm_publisher_name;
    }

    public Float getBqgm_price() {
        return bqgm_price;
    }

    public void setBqgm_price(Float bqgm_price) {
        this.bqgm_price = bqgm_price;
    }

    public String getBqtse_first_supply_date() {
        return bqtse_first_supply_date;
    }

    public void setBqtse_first_supply_date(String bqtse_first_supply_date) {
        this.bqtse_first_supply_date = bqtse_first_supply_date;
    }

    public String getBqtse_last_supply_date() {
        return bqtse_last_supply_date;
    }

    public void setBqtse_last_supply_date(String bqtse_last_supply_date) {
        this.bqtse_last_supply_date = bqtse_last_supply_date;
    }

    public String getBqtse_last_sale_date() {
        return bqtse_last_sale_date;
    }

    public void setBqtse_last_sale_date(String bqtse_last_sale_date) {
        this.bqtse_last_sale_date = bqtse_last_sale_date;
    }

    public String getBqtse_last_order_date() {
        return bqtse_last_order_date;
    }

    public void setBqtse_last_order_date(String bqtse_last_order_date) {
        this.bqtse_last_order_date = bqtse_last_order_date;
    }

    public String getBqgm_sales_date() {
        return bqgm_sales_date;
    }

    public void setBqgm_sales_date(String bqgm_sales_date) {
        this.bqgm_sales_date = bqgm_sales_date;
    }

    public String getBqct_media_group1_cd() {
        return bqct_media_group1_cd;
    }

    public String getBqct_media_group1_name() {
        return bqct_media_group1_name;
    }

    public void setBqct_media_group1_name(String bqct_media_group1_name) {
        this.bqct_media_group1_name = bqct_media_group1_name;
    }

    public String getBqct_media_group2_cd() {
        return bqct_media_group2_cd;
    }

    public void setBqct_media_group2_cd(String bqct_media_group2_cd) {
        this.bqct_media_group2_cd = bqct_media_group2_cd;
    }

    public String getBqct_media_group2_name() {
        return bqct_media_group2_name;
    }

    public void setBqct_media_group2_name(String bqct_media_group2_name) {
        this.bqct_media_group2_name = bqct_media_group2_name;
    }

    public void setBqct_media_group1_cd(String bqct_media_group1_cd) {
        this.bqct_media_group1_cd = bqct_media_group1_cd;
    }

    public String getBqgm_goods_name() {
        return bqgm_goods_name;
    }

    public void setBqgm_goods_name(String bqgm_goods_name) {
        this.bqgm_goods_name = bqgm_goods_name;
    }

    public String getBqgm_publisher_cd() {
        return bqgm_publisher_cd;
    }

    public void setBqgm_publisher_cd(String bqgm_publisher_cd) {
        this.bqgm_publisher_cd = bqgm_publisher_cd;
    }

    public String getBqio_trn_date() {
        return bqio_trn_date;
    }

    public void setBqio_trn_date(String bqio_trn_date) {
        this.bqio_trn_date = bqio_trn_date;
    }

    public int getSts_total_sales() {
        return sts_total_sales;
    }

    public void setSts_total_sales(int sts_total_sales) {
        this.sts_total_sales = sts_total_sales;
    }

    public int getSts_total_supply() {
        return sts_total_supply;
    }

    public void setSts_total_supply(int sts_total_supply) {
        this.sts_total_supply = sts_total_supply;
    }

    public int getSts_total_return() {
        return sts_total_return;
    }

    public void setSts_total_return(int sts_total_return) {
        this.sts_total_return = sts_total_return;
    }

    public String getLocation_id() {
        return location_id;
    }

    public void setLocation_id(String location_id) {
        this.location_id = location_id;
    }
}
