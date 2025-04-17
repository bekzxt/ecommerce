

CREATE TABLE products (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL,
                          description TEXT NOT NULL,
                          price NUMERIC(10, 2) NOT NULL,
                          stock INT NOT NULL,
                          category_id INT NOT NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
