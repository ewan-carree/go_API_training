-- Create the extension if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the people table
CREATE TABLE IF NOT EXISTS people (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100)
);
