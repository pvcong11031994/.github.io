package com.android.returncandidate.db.entity;

import java.io.*;

/**
 * Users entity
 *
 * @author tien-lv
 * @since 2017-12-05
 */

public class Users implements Serializable {

    /**
     * user id
     */
    private String userID;

    /**
     * password
     */
    private String password;

    /**
     * user name
     */
    private String name;

    /**
     * uid
     */
    private String uid;

    /**
     * shop id
     */
    private String shop_id;

    /**
     * shop name
     */
    private String shop_name;

    /**
     * login key
     */
    private String login_key;

    /**
     * server name
     */
    private String server_name;

    /**
     * user role
     */
    private int role;

    /**
     * license HSM
     */
    private String license;

    /**
     * Constructor Users
     */
    public Users() {
    }

    /**
     * Constructor Users
     */
    public Users(String userID, String password, String name, String uid, String shop_id,
                 String shop_name, String login_key, String server_name, int role) {

        this.userID = userID;
        this.password = password;
        this.name = name;
        this.uid = uid;
        this.shop_id = shop_id;
        this.shop_name = shop_name;
        this.login_key = login_key;
        this.server_name = server_name;
        this.role = role;
    }

    /**
     * Get userID
     */
    public String getUserid() {
        return userID;
    }

    /**
     * Get user name
     */
    public String getName() {
        return name;
    }

    /**
     * Get uid
     */
    public String getUid() {
        return uid;
    }

    /**
     * Get shop id
     */
    public String getShop_id() {
        return shop_id;
    }

    /**
     * Get shop name
     */
    public String getShop_name() {
        return shop_name;
    }

    /**
     * Set user id
     */
    public void setUserID(String userID) {
        this.userID = userID;
    }

    /**
     * Set user name
     */
    public void setName(String name) {
        this.name = name;
    }

    /**
     * Set uid
     */
    public void setUid(String Uid) {
        this.uid = Uid;
    }

    /**
     * Set shop id
     */
    public void setShop_id(String shop_id) {
        this.shop_id = shop_id;
    }

    /**
     * Set shop name
     */
    public void setShop_name(String shop_name) {
        this.shop_name = shop_name;
    }

    /**
     * Get login key
     */
    public String getLogin_key() {
        return login_key;
    }

    /**
     * Set login key
     */
    public void setLogin_key(String login_key) {
        this.login_key = login_key;
    }

    /**
     * Get user role
     */
    public int getRole() {
        return role;
    }

    /**
     * Set user role
     */
    public void setRole(int role) {
        this.role = role;
    }

    /**
     * Get server name
     */
    public String getServer_name() {
        return server_name;
    }

    /**
     * Set server name
     */
    public void setServer_name(String server_name) {
        this.server_name = server_name;
    }

    /**
     * Get password
     */
    public String getPassword() {
        return password;
    }

    /**
     * Set password
     */
    public void setPassword(String password) {
        this.password = password;
    }

    /**
     * Get license
     */
    public String getLicense() {
        return license;
    }

    /**
     * Set license
     */
    public void setLicense(String license) {
        this.license = license;
    }
}
