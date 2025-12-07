-- Create table
CREATE TABLE IF NOT EXISTS books (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  author TEXT NOT NULL,
  year INT
);

-- Sample seed (id auto)
INSERT INTO books (title, author, year) VALUES
('The Hobbit', 'J.R.R. Tolkien', 1937),
('Dune', 'Frank Herbert', 1965)
ON CONFLICT DO NOTHING;
