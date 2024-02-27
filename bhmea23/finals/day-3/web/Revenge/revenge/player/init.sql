CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    hashed_password VARCHAR(255) NOT NULL,
    profile_image_data TEXT CHECK (length(profile_image_data) <= 200720), 
    created_at TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS tracking (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    http_referer TEXT
);

GRANT SELECT, INSERT ON TABLE users, tracking TO flaskuser;
GRANT CONNECT ON DATABASE db TO flaskuser;
GRANT pg_write_server_files TO flaskuser;
grant usage, select on all sequences in schema public to flaskuser;

