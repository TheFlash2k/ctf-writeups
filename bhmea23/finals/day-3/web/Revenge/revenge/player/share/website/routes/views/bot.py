from .route import web
from flask import Flask, request, jsonify
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
import time
import threading
import queue
from website.config import SECRET_AGENT, MAX_CONCURRENT_THREADS, FLAG
import re

thread_local = threading.local()
url_queue = queue.Queue()


def get_webdriver():
    if not hasattr(thread_local, "webdriver"):
        options = webdriver.ChromeOptions()
        options.add_argument("--no-sandbox")
        options.add_argument("--headless")
        options.add_argument("--ignore-certificate-errors")
        options.add_argument("--disable-dev-shm-usage")
        options.add_argument("--disable-http2")
        options.add_argument("--disable-infobars")
        options.add_argument("--disable-default-apps")
        options.add_argument("--disable-gpu")
        options.add_argument("--disable-sync")
        options.add_argument("--hide-scrollbars")
        options.add_argument("--metrics-recording-only")
        options.add_argument("--no-first-run")
        options.add_argument("--safebrowsing-disable-auto-update")
        options.add_argument("--no-default-browser-check")
        options.add_argument("--media-cache-size=1")
        options.add_argument("--disk-cache-size=1")
        options.add_argument(f"--user-agent={SECRET_AGENT}")

        service = Service("/usr/bin/chromedriver")
        driver = webdriver.Chrome(options=options, service=service)

        thread_local.webdriver = driver

    return thread_local.webdriver


def release_webdriver():
    if hasattr(thread_local, "webdriver"):
        thread_local.webdriver.quit()
        del thread_local.webdriver


def process_url(url):
    try:
        driver = get_webdriver()

        driver.get("http://proxy.local/")
        driver.add_cookie({"name": "flag", "value": FLAG, "domain": "proxy.local"})

        driver.get(url)
        time.sleep(50)
        return "done"
    except Exception as e:
        return f"Error: {str(e)}"


def process_queue():
    while True:
        try:
            if threading.active_count() <= MAX_CONCURRENT_THREADS:
                # Get the next URL from the waiting list
                url = url_queue.get(timeout=1)
                threading.Thread(target=process_and_release, args=(url,)).start()
            else:
                time.sleep(1)
        except queue.Empty:
            time.sleep(1)


def process_and_release(url):
    try:
        result = process_url(url)
        print(result)
    finally:
        url_queue.task_done()
        release_webdriver()


# Start the URL processing thread
url_processing_thread = threading.Thread(target=process_queue)
url_processing_thread.daemon = True
url_processing_thread.start()


def use_regex(input_text):
    return input_text.startswith(
        ("http://secret.local/", "http://flask.local/"),
    )


@web.route("/bot", methods=["GET", "POST"])
def index():
    if request.method == "POST":
        url = request.form["url"]
        if use_regex(url):
            url_queue.put(url)
            return "URL added to the waiting list"
    return "none"
