CREATE TABLE Users (
  id SERIAL PRIMARY KEY,
  username TEXT,
  password_hash TEXT,
  email TEXT ,
  phone TEXT,
  last_login TIMESTAMP,
  registered_at TIMESTAMP,
  session JSONB
);

ALTER TABLE Users ADD CONSTRAINT unique_email UNIQUE (email);

CREATE TABLE Admins (
  id SERIAL PRIMARY KEY,
  username TEXT,
  password TEXT,
  email TEXT,
  last_login TIMESTAMP,
  session JSONB
);


CREATE TABLE Products (
  id SERIAL  PRIMARY KEY,
  name TEXT,
  description TEXT,
  price REAL,
  quantity INTEGER
);


CREATE TABLE Carts (
  id SERIAL  PRIMARY KEY,
  user_id INTEGER,
  FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

CREATE TABLE CartItems (
  id SERIAL  PRIMARY KEY,
  cart_id INTEGER,
  product_id INTEGER,
  quantity INTEGER,
  FOREIGN KEY (cart_id) REFERENCES Carts(id) ON DELETE CASCADE,
  FOREIGN KEY (product_id) REFERENCES Products(id)
);

CREATE TYPE delivery_method_enum AS ENUM ('method1', 'method2', 'method3');
CREATE TYPE payment_method_enum AS ENUM ('method1', 'method2', 'method3');

CREATE TABLE Orders (
  id SERIAL PRIMARY KEY,
  user_id INTEGER,
  status TEXT,
  delivery_method delivery_method_enum,
  payment_method payment_method_enum,
  FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

CREATE TABLE OrderItems (
  id SERIAL PRIMARY KEY,
  order_id INTEGER,
  product_id INTEGER,
  quantity INTEGER,
  FOREIGN KEY (order_id) REFERENCES Orders(id) ON DELETE CASCADE,
  FOREIGN KEY (product_id) REFERENCES Products(id) ON DELETE CASCADE
);