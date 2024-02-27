package com.app.caching;

public class Coupon extends Blueprint {
    private static final long serialVersionUID = 1L;

    private String code;
    private boolean isValid = false;

    public Coupon(String id, String code) {
        super(id);
        this.code = code;
    }

    public String getCode() {
        return code;
    }

    public boolean getIsValid() {
        return isValid;
    }
}