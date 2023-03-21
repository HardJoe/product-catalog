CREATE TABLE product (
  sku VARCHAR(20) PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  category VARCHAR(50),
  display_case VARCHAR(50),
  image_id INT REFERENCES image(id),
  weight DECIMAL(10, 2), --in kilograms
  price DECIMAL(10, 2), --in rupiahs
  review_id INT REFERENCES product_review(id)
);

CREATE TABLE image (
  id SERIAL PRIMARY KEY,
  url TEXT NOT NULL,
  description TEXT
);

CREATE TABLE product_review (
  id SERIAL PRIMARY KEY,
  rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
  comment TEXT
);
