package com.app.services;

import com.app.entities.User;

import jakarta.servlet.http.HttpServletResponse;

public interface UserService {
    User getUserById(String id);

    User getUserBySession(String session) throws Exception;

    void refresh(String session) throws Exception;

    String getAvatar(String provider, String username) throws Exception;

    void register(User user) throws Exception;

    User login(User user, HttpServletResponse response) throws Exception;

    User updateUser(User user);
    User deleteUser(Long id);
}
