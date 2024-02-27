package com.app.controllers;

import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.time.Duration;

import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseCookie;
import org.springframework.http.ResponseEntity;

import com.app.entities.User;
import com.app.services.UserService;
import com.app.utils.ExceptionHandlers;
import com.app.utils.Flag;

import jakarta.servlet.http.Cookie;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import jakarta.validation.Valid;

@RestController
@CrossOrigin(originPatterns = "*", allowCredentials = "true", allowedHeaders = "*", methods = {
		RequestMethod.GET, RequestMethod.POST, RequestMethod.PUT, RequestMethod.DELETE, RequestMethod.OPTIONS })
@RequestMapping("/api/user")
public class UserController extends ExceptionHandlers {
	private UserService userService;

	public UserController(UserService userService) {
		super();
		this.userService = userService;
	}

	@GetMapping("")
	public ResponseEntity<User> get(
			HttpServletRequest request) throws Exception {
		String session = getSession(request);
		User user = userService.getUserBySession(session);
		if (user == null) {
			return null;
		}
		user.setPassword(null);
		return new ResponseEntity<>(user, HttpStatus.OK);
	}

	@GetMapping("/flag")
	public ResponseEntity<String> getFlag(
			HttpServletRequest request) throws Exception {
		String flag = "You need to subscribe to get the flag.";
		String session = getSession(request);
		User user = userService.getUserBySession(session);
		if (user == null) {
			return new ResponseEntity<>(flag, HttpStatus.OK);
		}
		if (user.getSubscribed()) {
			flag = new Flag().getFlag();
		}
		return new ResponseEntity<>(flag, HttpStatus.OK);
	}

	@PostMapping("/logout")
	public void logout(
			HttpServletResponse response) throws Exception {
		ResponseCookie cookie = ResponseCookie.from("user", "")
				.maxAge(Duration.ofSeconds(0))
				.path("/")
				.build();
		response.addHeader(HttpHeaders.SET_COOKIE, cookie.toString());
	}

	@PostMapping("register")
	public void register(
			@Valid @RequestBody User user,
			@RequestParam(value = "avatar", defaultValue = "ui-avatars.com") String avatarProvider) throws Exception {

		try {
			String avatar = userService.getAvatar(avatarProvider, user.getUsername());
			user.setAvatar(avatar);
		} catch (Exception e) {
		}
		userService.register(user);
	}

	@PostMapping("login")
	public ResponseEntity<User> register(
			@Valid @RequestBody User user,
			HttpServletResponse response) throws Exception {
		User newUser = userService.login(user, response);
		newUser.setPassword(null);
		return new ResponseEntity<>(newUser, HttpStatus.CREATED);
	}

	@PostMapping("refresh")
	public void refresh(HttpServletRequest request) throws Exception {
		Cookie[] cookies = request.getCookies();

		Cookie userCookie = null;
		for (Cookie cookie : cookies) {
			if (cookie.getName().equals("user")) {
				userCookie = cookie;
				break;
			}
		}
		userService.refresh(userCookie.getValue());
	}

	private String getSession(HttpServletRequest request) {
		Cookie[] cookies = request.getCookies();

		Cookie userCookie = null;
		if (cookies != null) {
			for (Cookie cookie : cookies) {
				if (cookie.getName().equals("user")) {
					userCookie = cookie;
					break;
				}
			}
		}

		if (userCookie == null)
			return null;

		return userCookie.getValue();
	}
}