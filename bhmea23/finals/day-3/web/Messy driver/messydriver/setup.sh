#!/bin/bash
#!/bin/bash
mkdir plugins

USERNAME=admin
PASSWORD=$(openssl rand -hex 16)
DIRECTORY_NAME=$(openssl rand -hex 16)


sqlite3 idp.db << EOF
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);
INSERT INTO users (username, password) VALUES ('$USERNAME', '$PASSWORD');
EOF

sqlite3 fileserver.db << EOF

CREATE TABLE IF NOT EXISTS profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nickname TEXT NOT NULL,
    nationality TEXT
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_id INTEGER,
    username TEXT UNIQUE NOT NULL,
    activated_plugin TEXT, -- This will store comma-separated plugin names. Consider using a separate table for a normalized approach.
    directory TEXT,
    FOREIGN KEY (profile_id) REFERENCES profiles(id)
);

INSERT INTO profiles (nickname, nationality)
VALUES ('adminos', 'UNKNOWN');

INSERT INTO users (profile_id, username, activated_plugin, directory)
VALUES (last_insert_rowid(), '$USERNAME', '', '$DIRECTORY_NAME');
EOF


mkdir -p uploads/$DIRECTORY_NAME uploads/tmp