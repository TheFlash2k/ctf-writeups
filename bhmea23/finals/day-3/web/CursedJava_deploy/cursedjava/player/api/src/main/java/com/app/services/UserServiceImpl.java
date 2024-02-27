package com.app.services;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.InputStream;
import java.io.ObjectInputStream;
import java.io.ObjectOutputStream;
import java.util.Base64;
import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import com.app.caching.Cached;
import com.app.entities.User;
import com.app.exceptions.ResourceNotFoundException;
import com.app.repositories.UserRepository;
import com.app.services.Security;

import jakarta.servlet.http.Cookie;
import jakarta.servlet.http.HttpServletResponse;

import java.net.HttpURLConnection;
import java.net.URL;

@Service
public class UserServiceImpl implements UserService {

    private UserRepository userRepository;

    @Autowired
    private CachedService cachedService;

    @Bean
    public PasswordEncoder encoder() {
        return new BCryptPasswordEncoder();
    }

    public UserServiceImpl(UserRepository userRepository) {
        super();
        this.userRepository = userRepository;
    }

    @Override
    public User getUserById(String id) {
        Long longId = Long.parseLong(id);
        User existingUser = userRepository
                .findById(longId)
                .orElseThrow(() -> new ResourceNotFoundException(
                        "user not found with id: " + id));
        return existingUser;
    }

    @Override
    public User getUserBySession(String session) throws Exception {
        if (session == null) {
            return null;
        }
        try {
            byte[] decodedBytes = Base64.getDecoder().decode(session.getBytes());
            ByteArrayInputStream bis = new ByteArrayInputStream(decodedBytes);
            ObjectInputStream ois = new Security(bis);
            Object o = ois.readObject();
            com.app.caching.User obj = (com.app.caching.User) o;
            return getUserById(obj.getId());
        } catch (Exception e) {
            return null;
        }
    }

    @Override
    public String getAvatar(String provider, String username) throws Exception {
        String avatarUrl = "http://%s/api/%s";
        avatarUrl = String.format(avatarUrl, provider, username);
        try {
            URL url = new URL(avatarUrl);
            HttpURLConnection con = (HttpURLConnection) url.openConnection();
            con.setRequestMethod("GET");
            InputStream inputStream = con.getInputStream();
            byte[] bytes = inputStream.readAllBytes();
            byte[] encoded = Base64.getEncoder().encode(bytes);
            String encodedString = new String(encoded);

            if (encodedString.length() == 0)
                return avatarUrl;
            return encodedString;
        } catch (Exception e) {
            e.printStackTrace();
            throw new Exception(e.getMessage());
        }
    }

    @Override
    public void register(User user) throws Exception {
        user.setPassword(encoder().encode(user.getPassword()));
        List<User> existing = userRepository.findByUsername(user.getUsername());

        if (existing.size() > 0) {
            throw new Exception("User already exists");
        }

        User created = new User();
        created.setUsername(user.getUsername());
        created.setPassword(user.getPassword());
        created.setAvatar(user.getAvatar());
        userRepository.save(created);
    }

    @Override
    public User login(User user, HttpServletResponse response) throws Exception {
        List<User> existing = userRepository.findByUsername(user.getUsername());

        if (existing.size() == 0) {
            throw new Exception("User does not exist");
        }

        User created = existing.get(0);

        if (!encoder().matches(user.getPassword(), created.getPassword())) {
            throw new Exception("Password is incorrect");
        }

        com.app.caching.User object = new com.app.caching.User(
                created.getId().toString(),
                created.getUsername(),
                created.getAvatar());

        Cached cached = new Cached(
                "user:" + created.getId().toString(),
                object);

        ByteArrayOutputStream baos = new ByteArrayOutputStream();
        ObjectOutputStream oos = new ObjectOutputStream(baos);
        oos.writeObject(object);
        oos.close();
        byte[] bytes = baos.toByteArray();
        byte[] encoded = Base64.getEncoder().encode(bytes);
        String encodedString = new String(encoded);
        Cookie cookie = new Cookie("user", encodedString);
        cookie.setPath("/");
        response.addCookie(cookie);
        cachedService.save(cached);
        return created;
    }

    @Override
    public void refresh(String session) throws Exception {
        byte[] decodedBytes = Base64.getDecoder().decode(session.getBytes());
        ByteArrayInputStream bis = new ByteArrayInputStream(decodedBytes);
        ObjectInputStream ois = new Security(bis);
        Object o = ois.readObject();
        com.app.caching.Blueprint obj = (com.app.caching.Blueprint) o;

        Cached cached = new Cached(
                obj.getId().toString(),
                o);

        cachedService.save(cached);
    }

    @Override
    public User updateUser(User user) {
        User existingUser = userRepository
                .findById(user.getId())
                .orElseThrow(() -> new ResourceNotFoundException(
                        "user not found with id: " + user.getId()));

        if (user.getUsername() != null)
            existingUser.setUsername(user.getUsername());
        if (user.getPassword() != null)
            existingUser.setPassword(encoder().encode(user.getPassword()));

        return userRepository.save(existingUser);
    }

    @Override
    public User deleteUser(Long id) {
        User existingUser = userRepository
                .findById(id)
                .orElseThrow(() -> new ResourceNotFoundException(
                        "user not found with id: " + id));
        userRepository.delete(existingUser);
        return existingUser;
    }
}
