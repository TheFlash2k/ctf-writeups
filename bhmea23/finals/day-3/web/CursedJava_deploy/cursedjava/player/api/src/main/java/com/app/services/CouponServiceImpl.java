package com.app.services;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Repository;

import com.app.caching.Cached;
import com.app.caching.Coupon;
import com.app.entities.User;

@Repository
public class CouponServiceImpl implements CouponService {

    @Autowired
    private CachedService cachedService;

    @Autowired
    private UserService userService;

    @Override
    public Coupon getCouponById(String id) {
        Cached cached = cachedService.one("coupon:" + id);
        if (cached == null) {
            return null;
        }
        return (Coupon) cached.getValue();
    }

    @Override
    public Coupon createCoupons(String code) throws Exception {

        Coupon existing = getCouponById(code);

        if (existing != null) {
            throw new Exception("Coupon already exists");
        }

        Coupon coupon = new Coupon(code, code);
        cachedService.save(new Cached("coupon:" + coupon.getId(), coupon));
        return coupon;
    }

    @Override
    public void useCoupon(String id, String userId) throws Exception {
        User u = userService.getUserById(userId);
        if (u == null) {
            throw new Exception("User not found");
        }

        Coupon coupon = getCouponById(id);

        if (coupon == null) {
            throw new Exception("Coupon not found");
        }

        if (!coupon.getIsValid()) {
            throw new Exception("Coupon is not valid");
        }

        u.setSubscribed(true);

        userService.updateUser(u);
    }
}