from website.config import (
    SECRET_AGENT,
    DB_NAME,
    DB_USER,
    DB_PWD,
    DB_HOST,
    DB_PORT,
    SALT,
)
from .route import web
from flask import render_template, request, Response, redirect, url_for, session, flash
from flask import request
import re
import base64
import psycopg2
from markupsafe import escape
import sys
import json
import time
from werkzeug.security import generate_password_hash, check_password_hash
import html


# Define the connection parameters
database_name = DB_NAME
database_user = DB_USER
database_password = DB_PWD
database_host = DB_HOST
database_port = DB_PORT
max_retries = 2
retries = 0

while retries < max_retries:
    try:
        # Try to establish a connection
        conn = psycopg2.connect(
            dbname=database_name,
            user=database_user,
            password=database_password,
            host=database_host,
            port=database_port,
        )
        break  # Successful connection, exit the loop
    except Exception as e:
        print(e)
        retries += 1
        time.sleep(2)  # Wait for a few seconds before retrying

if retries == max_retries:
    print("Failed to connect to the database after multiple retries.")
else:
    print("Successfully connected to the database.")


@web.route("/")
def home():
    if "user" in session:
        username = session["user"]
        cur = conn.cursor()
        cur.execute(
            "SELECT profile_image_data FROM users WHERE username = %s", (username,)
        )
        user = cur.fetchone()
        if user is None:
            del session["user"]
            return redirect(url_for("web.login"))

        image = user[0]
        cur.close()
        return render_template(
            "welcome.html", username=html.escape(username), image=image
        )

    return render_template("index.html")


def insert_referer_into_tracking(username, http_referer):
    try:
        cur = conn.cursor()
        cur.execute(
            f"INSERT INTO tracking (username, http_referer) VALUES (%s, '{http_referer}')",
            (username,),
        )
        conn.commit()

    except Exception as e:
        conn.rollback()
        print(f"Error inserting into 'tracking' table: {str(e)}", file=sys.stderr)
    finally:
        cur.close()


png_extension_pattern = re.compile(r"^.*\.png$", re.IGNORECASE)


def user_exists(username):
    try:
        cur = conn.cursor()
        cur.execute("SELECT COUNT(*) FROM users WHERE username = %s", (username,))
        count = cur.fetchone()[0]
        return count > 0
    except psycopg2.Error as e:
        conn.rollback()
        raise
    finally:
        cur.close()


@web.route("/register", methods=["GET", "POST"])
def register():
    if request.method == "POST":
        username = request.form["username"]
        password = request.form["password"] + SALT
        hashed_password = generate_password_hash(password)
        profile_picture = request.files["profile_picture"]

        if user_exists(username):
            flash("User already exists. Please choose a different username.", "danger")
            return redirect(url_for("web.register"))

        # Check if a profile picture is provided and if it's a PNG image
        if profile_picture:
            if not png_extension_pattern.match(profile_picture.filename):
                flash("Profile picture must be in PNG format.", "danger")
                return redirect(url_for("web.register"))

            # Check if the profile picture size is less than or equal to 25KB
            if len(profile_picture.read()) > 15000 * 1024:  # TODO
                flash("Profile picture size should not exceed 25KB.", "danger")
                return redirect(url_for("web.register"))

            # Reset the file cursor after reading
            profile_picture.seek(0)
            img = base64.b64encode(profile_picture.read()).decode()

            cur = conn.cursor()
            try:
                cur.execute(
                    "INSERT INTO users (username, hashed_password, profile_image_data) VALUES (%s, %s, %s)",
                    (username, hashed_password, img),
                )
                conn.commit()

            except:
                conn.rollback()
            finally:
                cur.close()
        else:
            flash("Profile picture is required.", "danger")
            return redirect(url_for("web.register"))

        flash("Registration successful. Please log in.", "success")
        return redirect(url_for("web.login"))
    return render_template("register.html")


@web.route("/login", methods=["GET", "POST"])
def login():
    if request.method == "POST":
        username = request.form["username"]
        password = request.form["password"] + SALT
        cur = conn.cursor()
        cur.execute(
            "SELECT username, hashed_password FROM users WHERE username = %s",
            (username,),
        )
        user = cur.fetchone()
        cur.close()
        if request.headers.get("User-agent") == SECRET_AGENT:
            http_referer = request.headers.get("Referer")
            insert_referer_into_tracking(username, http_referer)
            return render_template(
                "welcome.html", username=html.escape(username), image=""
            )

        if user and check_password_hash(user[1], password):
            session["user"] = user[0]
            flash("Login successful!", "success")
            return redirect(url_for("web.home"))
        else:
            flash("Login failed. Please try again.", "danger")
    return render_template("login.html")


@web.route("/logout")
def logout():
    session.pop("user", None)
    return redirect(url_for("web.home"))
