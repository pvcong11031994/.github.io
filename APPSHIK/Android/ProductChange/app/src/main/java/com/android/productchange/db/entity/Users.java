package com.android.productchange.db.entity;

import java.io.Serializable;

/**
 * <h1>Users</h1>
 *
 * Entity Users
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
     * create date
     */
    private String create_date;

    /**
     * Constructor Users
     */
    public Users() {
    }

    /**
     * Get userID
     *
     * @return userID
     */
    public String getUserid() {
        return userID;
    }

    /**
     * Get user name
     *
     * @return user name
     */
    public String getName() {
        return name;
    }

    /**
     * Get uid
     *
     * @return uid
     */
    public String getUid() {
        return uid;
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
     * Get shop name
     *
     * @return shop name
     */
    public String getShop_name() {
        return shop_name;
    }

    /**
     * Set user id
     *
     * @param userID is user id
     */
    public void setUserID(String userID) {
        this.userID = userID;
    }

    /**
     * Set user name
     *
     * @param name is user name
     */
    public void setName(String name) {
        this.name = name;
    }

    /**
     * Set uid
     *
     * @param Uid is uid
     */
    public void setUid(String Uid) {
        this.uid = Uid;
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
     * Set shop name
     *
     * @param shop_name is shop name
     */
    public void setShop_name(String shop_name) {
        this.shop_name = shop_name;
    }

    /**
     * Get login key
     *
     * @return login key
     */
    public String getLogin_key() {
        return login_key;
    }

    /**
     * Set login key
     *
     * @param login_key is login key
     */
    public void setLogin_key(String login_key) {
        this.login_key = login_key;
    }

    /**
     * Get user role
     *
     * @return user role
     */
    public int getRole() {
        return role;
    }

    /**
     * Set user role
     *
     * @param role is user role
     */
    public void setRole(int role) {
        this.role = role;
    }

    /**
     * Get server name
     *
     * @return server name
     */
    public String getServer_name() {
        return server_name;
    }

    /**
     * Set server name
     *
     * @param server_name is server name
     */
    public void setServer_name(String server_name) {
        this.server_name = server_name;
    }

    /**
     * Get password
     *
     * @return password
     */
    public String getPassword() {
        return password;
    }

    /**
     * Set password
     *
     * @param password is password
     */
    public void setPassword(String password) {
        this.password = password;
    }

    /**
     * Get create date
     *
     * @return String
     */
    public String getCreate_date() {
        return create_date;
    }

    /**
     * Set create date
     *
     * @param create_date {@link String}
     */
    public void setCreate_date(String create_date) {
        this.create_date = create_date;
    }
}
