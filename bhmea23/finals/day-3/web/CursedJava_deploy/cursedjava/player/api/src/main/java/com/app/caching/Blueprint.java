package com.app.caching;

import lombok.Data;
import java.io.Serializable;

@Data
public class Blueprint implements Serializable {

    private static final long serialVersionUID = 1L;

    private String id;

    public Blueprint(String id) {
        this.id = id;
    }

    public String getId() {
        return id;
    }

}