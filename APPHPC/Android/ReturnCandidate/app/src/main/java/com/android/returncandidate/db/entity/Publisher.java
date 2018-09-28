package com.android.returncandidate.db.entity;

import java.io.Serializable;

/**
 * Large_classifications, Location, Publishers entity
 *
 * @author cong-pv
 * @since 2018-07-09
 */

public class Publisher implements Serializable {

    /**
     * id
     */
    private String id;

    /**
     * name
     */
    private String name;

    /**
     * Constructor CLP
     */
    public Publisher() {

    }

    /**
     * Constructor CLP
     */
    public Publisher(String id, String name) {

        this.id = id;
        this.name = name.trim();
    }

    /**
     * Get id
     */
    public String getId() {
        return id;
    }

    /**
     * Set id
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Get name
     */
    public String getName() {
        return name;
    }

    /**
     * Set name
     */
    public void setName(String name) {
        this.name = name.trim();
    }
}
