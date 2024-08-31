CREATE SCHEMA wisdom;
CREATE TABLE wisdom.services (
  service_id UUID PRIMARY KEY,
  name VARCHAR UNIQUE,
  price NUMERIC(12,2)
);

CREATE TABLE wisdom.customers (
   customer_id UUID PRIMARY KEY,
   first_name VARCHAR,
   last_name VARCHAR,
   email VARCHAR,
   phone VARCHAR,
   address VARCHAR
);

CREATE TABLE wisdom.vendors (
     vendor_id UUID PRIMARY KEY,
     name VARCHAR NOT NULL,
     contact VARCHAR,
     phone VARCHAR,
     email VARCHAR,
     address VARCHAR
);

CREATE TABLE wisdom.products (
      product_id UUID PRIMARY KEY,
      name VARCHAR UNIQUE,
      price NUMERIC (12,2),
      vendor_id UUID NOT NULL,
      FOREIGN KEY (VENDOR_ID) references wisdom.vendors(VENDOR_ID)
);