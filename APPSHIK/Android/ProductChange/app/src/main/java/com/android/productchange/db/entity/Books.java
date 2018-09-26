package com.android.productchange.db.entity;

import com.google.gson.annotations.Expose;
import com.google.gson.annotations.SerializedName;

import java.io.Serializable;

/**
 * <h1>Entity Books</h1>
 *
 * @author tien-lv
 * @since 2017-11-30
 */
public class Books implements Serializable {

    /**
     * id
     */
    @SerializedName("id")
    @Expose
    private int id;
    /**
     * location id
     */
    @SerializedName("location_id")
    @Expose
    private String location_id;
    /**
     * location name
     */
    @SerializedName("location_name")
    @Expose
    private String location_name;
    /**
     * large classifications id
     */
    @SerializedName("large_classification_id")
    @Expose
    private String large_classifications_id;
    /**
     * large classifications name
     */
    @SerializedName("large_classification_name")
    @Expose
    private String large_classifications_name;
    /**
     * jan code
     */
    @SerializedName("jan_code")
    @Expose
    private String jan_code;
    /**
     * name
     */
    @SerializedName("name")
    @Expose
    private String name;
    /**
     * author
     */
    @SerializedName("author")
    @Expose
    private String author;
    /**
     * publisher id
     */
    @SerializedName("publisher_id")
    @Expose
    private String publisher_id;
    /**
     * publisher name
     */
    @SerializedName("publisher_name")
    @Expose
    private String publisher_name;
    /**
     * publish date
     */
    @SerializedName("publish_date")
    @Expose
    private String publish_date;
    /**
     * price
     */
    @SerializedName("price")
    @Expose
    private float price;
    /**
     * shop id
     */
    @SerializedName("shop_id")
    @Expose
    private String shop_id;
    /**
     * final purchase date
     */
    @SerializedName("final_purchase_date")
    @Expose
    private String final_purchase_date;
    /**
     * final sale date
     */
    @SerializedName("final_sale_date")
    @Expose
    private String final_sale_date;
    /**
     * sale number
     */
    @SerializedName("sale_number")
    @Expose
    private int sale_number;
    /**
     * inventory number
     */
    @SerializedName("inventory_number")
    @Expose
    private int inventory_number;
    /**
     * national sale number
     */
    @SerializedName("national_sale_number")
    @Expose
    private int national_sale_number;
    /**
     * scan status
     */
    @SerializedName("scan_status")
    @Expose
    private int scan_status;
    /**
     * old catagory rank
     */
    @SerializedName("old_category_rank")
    @Expose
    private int old_catagory_rank;
    /**
     * new catagory rank
     */
    @SerializedName("new_category_rank")
    @Expose
    private int new_catagory_rank;
    /**
     * flag order now
     */
    @SerializedName("flag_order_now")
    @Expose
    private int flag_order_now;
    /**
     * ranking
     */
    @SerializedName("Ranking")
    @Expose
    private int ranking;

    /**
     * Constructor Books
     */
    public Books() {

    }

    /**
     * Get id
     *
     * @return id
     */
    public int getId() {
        return id;
    }

    /**
     * Set id
     *
     * @param id is id
     */
    public void setId(int id) {
        this.id = id;
    }

    /**
     * Get name
     *
     * @return name
     */
    public String getName() {
        return name;
    }

    /**
     * Set name
     *
     * @param name is name
     */
    public void setName(String name) {
        this.name = name;
    }

    /**
     * Get location id
     *
     * @return location id
     */
    public String getLocation_id() {
        return location_id;
    }

    /**
     * Set location id
     *
     * @param location_id is location id
     */
    public void setLocation_id(String location_id) {
        this.location_id = location_id;
    }

    /**
     * Get large classifications
     *
     * @return large classifications
     */
    public String getLarge_classifications_id() {
        return large_classifications_id;
    }

    /**
     * Set large classifications
     *
     * @param large_classifications_id is large classifications id
     */
    public void setLarge_classifications_id(String large_classifications_id) {
        this.large_classifications_id = large_classifications_id;
    }

    /**
     * Get jan code
     *
     * @return jan code
     */
    public String getJan_code() {
        return jan_code;
    }

    /**
     * Set jan code
     *
     * @param jan_code is jan code
     */
    public void setJan_code(String jan_code) {
        this.jan_code = jan_code;
    }

    /**
     * Get author
     *
     * @return author
     */
    public String getAuthor() {
        return author;
    }

    /**
     * Set author
     *
     * @param author is author
     */
    public void setAuthor(String author) {
        this.author = author;
    }

    /**
     * Get price
     *
     * @return price
     */
    public float getPrice() {
        return price;
    }

    /**
     * Set price
     *
     * @param price is price
     */
    public void setPrice(float price) {
        this.price = price;
    }

    /**
     * Get final purchase date
     *
     * @return final purchase date
     */
    public String getFinal_purchase_date() {
        return final_purchase_date;
    }

    /**
     * Set final purchase date
     *
     * @param final_purchase_date is final purchase date
     */
    public void setFinal_purchase_date(String final_purchase_date) {
        this.final_purchase_date = final_purchase_date;
    }

    /**
     * Get sale number
     *
     * @return sale number
     */
    public int getSale_number() {
        return sale_number;
    }

    /**
     * Set sale number
     *
     * @param sale_number is sale number
     */
    public void setSale_number(int sale_number) {
        this.sale_number = sale_number;
    }

    /**
     * Get inventory number
     *
     * @return inventory number
     */
    public int getInventory_number() {
        return inventory_number;
    }

    /**
     * Set inventory number
     *
     * @param inventory_number is inventory number
     */
    public void setInventory_number(int inventory_number) {
        this.inventory_number = inventory_number;
    }

    /**
     * Get national sale number
     *
     * @return national sale number
     */
    public int getNational_sale_number() {
        return national_sale_number;
    }

    /**
     * Set national sale number
     *
     * @param national_sale_number is national sale number
     */
    public void setNational_sale_number(int national_sale_number) {
        this.national_sale_number = national_sale_number;
    }

    /**
     * Get scan status
     *
     * @return scan status
     */
    public int getScan_status() {
        return scan_status;
    }

    /**
     * Set scan status
     *
     * @param scan_status is scan status
     */
    public void setScan_status(int scan_status) {
        this.scan_status = scan_status;
    }

    /**
     * Get publisher id
     *
     * @return publisher id
     */
    public String getPublisher_id() {
        return publisher_id;
    }

    /**
     * Set publisher id
     *
     * @param publisher_id is publisher id
     */
    public void setPublisher_id(String publisher_id) {
        this.publisher_id = publisher_id;
    }

    /**
     * Get publish date
     *
     * @return publish date
     */
    public String getPublish_date() {
        return publish_date;
    }

    /**
     * Set publish date
     *
     * @param publish_date is publish date
     */
    public void setPublish_date(String publish_date) {
        this.publish_date = publish_date;
    }

    /**
     * Get shop id
     *
     * @return shop id
     */
    public String getShop_id() {
        return shop_id;
    }

    /**
     * Set shop id
     *
     * @param shop_id is shop id
     */
    public void setShop_id(String shop_id) {
        this.shop_id = shop_id;
    }

    /**
     * Get final sale date
     *
     * @return final sale date
     */
    public String getFinal_sale_date() {
        return final_sale_date;
    }

    /**
     * Set final sale date
     *
     * @param final_sale_date is final sale date
     */
    public void setFinal_sale_date(String final_sale_date) {
        this.final_sale_date = final_sale_date;
    }

    /**
     * Get location name
     *
     * @return location name
     */
    public String getLocation_name() {
        return location_name;
    }

    /**
     * Set location name
     *
     * @param location_name is location name
     */
    public void setLocation_name(String location_name) {
        this.location_name = location_name;
    }

    /**
     * Get large classifications name
     *
     * @return large classifications name
     */
    public String getLarge_classifications_name() {
        return large_classifications_name;
    }

    /**
     * Set large classifications name
     *
     * @param large_classifications_name is large classifications name
     */
    public void setLarge_classifications_name(String large_classifications_name) {
        this.large_classifications_name = large_classifications_name;
    }

    /**
     * Get publisher name
     *
     * @return publisher name
     */
    public String getPublisher_name() {
        return publisher_name;
    }

    /**
     * Set publisher name
     *
     * @param publisher_name is publisher name
     */
    public void setPublisher_name(String publisher_name) {
        this.publisher_name = publisher_name;
    }

    /**
     * Get old catagory rank
     *
     * @return old catagory rank
     */
    public int getOld_catagory_rank() {
        return old_catagory_rank;
    }

    /**
     * Set old catagory rank
     *
     * @param old_catagory_rank is old catagory rank
     */
    public void setOld_catagory_rank(int old_catagory_rank) {
        this.old_catagory_rank = old_catagory_rank;
    }

    /**
     * Get new catagory rank
     *
     * @return new catagory rank
     */
    public int getNew_catagory_rank() {
        return new_catagory_rank;
    }

    /**
     * Set new catagory rank
     *
     * @param new_catagory_rank is new catagory rank
     */
    public void setNew_catagory_rank(int new_catagory_rank) {
        this.new_catagory_rank = new_catagory_rank;
    }

    /**
     * Get flag order now
     *
     * @return flag order now
     */
    public int getFlag_order_now() {
        return flag_order_now;
    }

    /**
     * Set flag order now
     *
     * @param flag_order_now is flag order now
     */
    public void setFlag_order_now(int flag_order_now) {
        this.flag_order_now = flag_order_now;
    }

    /**
     * Get ranking
     *
     * @return ranking
     */
    public int getRanking() {
        return ranking;
    }

    /**
     * Set ranking
     *
     * @param ranking is ranking
     */
    public void setRanking(int ranking) {
        this.ranking = ranking;
    }
}
