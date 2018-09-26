package com.android.productchange.db.entity;

import com.google.gson.annotations.Expose;
import com.google.gson.annotations.SerializedName;

import java.io.Serializable;

/**
 * <h1> CLP </h1>
 * Entity of Large_classifications, Location, Publishers
 *
 * @author tien-lv
 * @since 2017-12-25
 */
public class CLP implements Serializable {

    /**
     * id
     */
    @SerializedName("id")
    @Expose
    private String id;

    /**
     * name
     */
    @SerializedName("name")
    @Expose
    private String name;

    /**
     * Constructor CLP
     */
    public CLP() {

    }

    /**
     * Get id
     *
     * @return id
     */
    public String getId() {
        return id;
    }

    /**
     * Set id
     *
     * @param id is id
     */
    public void setId(String id) {
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
        this.name = name.trim();
    }
}
