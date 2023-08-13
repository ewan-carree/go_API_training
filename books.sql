-- Create the extension if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the books table
CREATE TABLE IF NOT EXISTS books (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    title VARCHAR(100),
    author VARCHAR(100),
    person_id UUID
);
