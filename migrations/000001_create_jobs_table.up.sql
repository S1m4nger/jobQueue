CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    payload BLOB,
    status TEXT NOT NULL,
    result BLOB,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
