package com.app.services;

import com.app.caching.Coupon;

public interface CouponService {
    Coupon getCouponById(String id);

    Coupon createCoupons(String code) throws Exception;

    void useCoupon(String id, String userId) throws Exception;
}
