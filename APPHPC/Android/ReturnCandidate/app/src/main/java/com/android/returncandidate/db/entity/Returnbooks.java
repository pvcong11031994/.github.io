package com.android.returncandidate.db.entity;

import java.io.*;

/**
 * Returnbooks entity
 *
 * @author tien-lv
 * @since 2017-12-05
 */

public class Returnbooks implements Serializable {

    /**
     * id
     */
    private int id;
    /**
     * shop id
     */
    private String shop_id;
    /**
     * userid
     */
    private String userid;
    /**
     * return date
     */
    private String return_date;
    /**
     * jan code
     */
    private String jan_code;
    /**
     * number
     */
    private int number;
    /**
     * list status
     */
    private int list_status;

    /**
     * Constructor Returnbooks
     */
    public Returnbooks() {

    }

    /**
     * Constructor Returnbooks
     */
    public Returnbooks(int id, String shop_id, String userid, String return_date, String jan_code,
            int number, int list_status) {

        this.id = id;
        this.shop_id = shop_id;
        this.userid = userid;
        this.return_date = return_date;
        this.jan_code = jan_code;
        this.number = number;
        this.list_status = list_status;
    }

    /**
     * Get id
     */
    public int getId() {
        return id;
    }

    /**
     * Set id
     */
    public void setId(int id) {
        this.id = id;
    }

    /**
     * Get shop id
     */
    public String getShop_id() {
        return shop_id;
    }

    /**
     * Set shop id
     */
    public void setShop_id(String shop_id) {
        this.shop_id = shop_id;
    }

    /**
     * Get user id
     */
    public String getUserid() {
        return userid;
    }

    /**
     * Set user id
     */
    public void setUserid(String userid) {
        this.userid = userid;
    }

    /**
     * Get return date
     */
    public String getReturn_date() {
        return return_date;
    }

    /**
     * Set return date
     */
    public void setReturn_date(String return_date) {
        this.return_date = return_date;
    }

    /**
     * Get jan code
     */
    public String getJan_code() {
        return jan_code;
    }

    /**
     * Set jan code
     */
    public void setJan_code(String jan_code) {
        this.jan_code = jan_code;
    }

    /**
     * Get number
     */
    public int getNumber() {
        return number;
    }

    /**
     * Set number
     */
    public void setNumber(int number) {
        this.number = number;
    }

    /**
     * Get list status
     */
    public int getList_status() {
        return list_status;
    }

    /**
     * Set list status
     */
    public void setList_status(int list_status) {
        this.list_status = list_status;
    }
}
