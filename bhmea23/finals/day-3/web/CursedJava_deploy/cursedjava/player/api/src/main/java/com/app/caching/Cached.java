package com.app.caching;

import lombok.Data;
import java.io.Serializable;

@Data
public class Cached implements Serializable {
    private String id;
    private Object value;

    public Cached(String id, Object value) {
        this.id = id;
        this.value = value;
    }

    public String getId() {
        return id;
    }

    public Object getValue() {
        return value;
    }
}