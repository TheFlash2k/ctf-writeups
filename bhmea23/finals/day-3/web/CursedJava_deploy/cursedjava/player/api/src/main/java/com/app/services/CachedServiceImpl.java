package com.app.services;

import java.util.ArrayList;
import java.util.List;

import org.springframework.data.redis.core.HashOperations;
import org.springframework.stereotype.Repository;

import com.app.caching.Cached;

import jakarta.annotation.Resource;

@Repository
public class CachedServiceImpl implements CachedService {

    private final String hashReference = "Cached";

    @Resource(name = "redisTemplate")
    private HashOperations<String, String, Cached> hashOperations;

    @Override
    public void save(Cached cached) {
        hashOperations.put(cached.getId(), hashReference, cached);
    }

    @Override
    public Cached one(String id) {
        return hashOperations.get(id, hashReference);
    }
}