CREATE TABLE orders (
                        id TEXT PRIMARY KEY,
                        user_id TEXT NOT NULL,
                        total_price REAL NOT NULL,
                        status TEXT NOT NULL
);
CREATE TABLE order_items (
                             order_id TEXT,
                             product_id TEXT,
                             quantity INTEGER,
                             price REAL,
                             FOREIGN KEY(order_id) REFERENCES orders(id)
);
