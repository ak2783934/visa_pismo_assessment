CREATE TABLE IF NOT EXISTS accounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_number TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    account_id INTEGER NOT NULL,
    operation_type_id INTEGER NOT NULL,
    amount REAL NOT NULL,
    event_date DATETIME NOT NULL,
    FOREIGN KEY(account_id) REFERENCES accounts(id)
);
