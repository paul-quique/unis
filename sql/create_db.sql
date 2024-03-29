-- Créations des tables de la base de données --

CREATE TABLE user_info (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(45) NOT NULL,
    last_name VARCHAR(45) NOT NULL,
    email VARCHAR(75) UNIQUE NOT NULL,
    salted_hash VARCHAR(100) NOT NULL,
    points INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(45) NOT NULL
);

CREATE TABLE product (
    id SERIAL PRIMARY KEY,
    name VARCHAR(45) NOT NULL,
    user_id INT NOT NULL,
    category_id INTEGER NOT NULL,
    price INTEGER NOT NULL,
    FOREIGN KEY (category_id) REFERENCES category (id),
    FOREIGN KEY (user_id) REFERENCES user_info (id)
);

CREATE TABLE offer (
    id SERIAL PRIMARY KEY,
    borrower_id INTEGER NOT NULL,
    lender_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (borrower_id) REFERENCES user_info (id),
    FOREIGN KEY (lender_id) REFERENCES user_info (id),
    FOREIGN KEY (product_id) REFERENCES product (id)
);

CREATE TABLE session (
    id VARCHAR(100) PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user_info (id)
);

CREATE TABLE transaction (
    id SERIAL PRIMARY KEY,
    offer_id INTEGER NOT NULL,
    return_date TIMESTAMP NOT NULL,
    FOREIGN KEY (offer_id) REFERENCES offer (id)
);

CREATE TABLE message (
    id SERIAL PRIMARY KEY,
    sent_from INTEGER NOT NULL,
    sent_to INTEGER NOT NULL,
    sent_at TIMESTAMP NOT NULL,
    content VARCHAR(1000) NOT NULL,
    FOREIGN KEY (sent_from) REFERENCES user_info (id),
    FOREIGN KEY (sent_to) REFERENCES user_info (id)
);
