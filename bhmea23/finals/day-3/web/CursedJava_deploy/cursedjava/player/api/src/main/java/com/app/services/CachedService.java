package com.app.services;

import com.app.caching.Cached;

public interface CachedService {
    void save(Cached cached);

    Cached one(String id);
}