package com.app.controllers;

import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.net.InetAddress;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

import com.app.caching.Coupon;
import com.app.services.CouponService;
import com.app.utils.ExceptionHandlers;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.validation.Valid;

@RestController
@CrossOrigin()
@RequestMapping("/api/coupon")
public class CouponController extends ExceptionHandlers {
    private CouponService couponService;

    public CouponController(CouponService couponService) {
        super();
        this.couponService = couponService;
    }

    @RequestMapping("/{id}")
    public ResponseEntity<Coupon> get(
            @PathVariable("id") String id,
            HttpServletRequest request) {
        validateOrigin(request);
        Coupon todo = couponService.getCouponById(id);
        return new ResponseEntity<>(todo, HttpStatus.OK);
    }

    @RequestMapping("/generate/{code}")
    public ResponseEntity<Coupon> create(
            @Valid @PathVariable String code,
            HttpServletRequest request) throws Exception {
        validateOrigin(request);
        Coupon newCoupon = couponService.createCoupons(code);
        return new ResponseEntity<>(newCoupon, HttpStatus.CREATED);
    }

    @RequestMapping("/use/{id}")
    public void create(
            @Valid @PathVariable String id,
            @RequestParam(value = "userId", defaultValue = "0") String userId,
            HttpServletRequest request) throws Exception {
        validateOrigin(request);
        couponService.useCoupon(id, userId);
    }

    private void validateOrigin(HttpServletRequest request) {
        try {

            InetAddress requestAddress = InetAddress.getByName(request.getRemoteAddr());
            String uri = requestAddress.toString().split("/")[1];

            boolean isLocalhost = uri.equals("127.0.0.1");

            if (isLocalhost) {
                return;
            } else {
                throw new Exception("Not allowed");
            }
        } catch (Exception e) {
            System.out.println(e);
            throw new RuntimeException(e);
        }

    }
}