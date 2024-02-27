package com.app.utils;

public class Flag {
    public String getFlag() {
        if (System.getenv().containsKey("FLAG")) {
            return System.getenv("FLAG");
        } else {
            throw new RuntimeException("FLAG not found");
        }
    }
}
