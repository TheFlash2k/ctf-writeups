package com.app.caching;

public class User extends Blueprint {
    private static final long serialVersionUID = 1L;

    private String username, avatar;

    public User(String id, String username, String avatar) {
        super(id);
        this.username = username;
        this.avatar = avatar;
    }

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getAvatar() {
        return avatar;
    }

    public void setAvatar(String avatar) {
        this.avatar = avatar;
    }
}