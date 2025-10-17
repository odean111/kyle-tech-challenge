-- Deploy lothrop-backend:companies to pg

BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    jurisdiction VARCHAR(20) NOT NULL CHECK (jurisdiction IN ('UK', 'Singapore', 'Caymens')),
    company_name VARCHAR(255) NOT NULL,
    company_address TEXT NOT NULL,
    nature_of_business TEXT,
    number_of_directors INTEGER CHECK (number_of_directors >= 0),
    number_of_shareholders INTEGER CHECK (number_of_shareholders >= 0),
    sec_code VARCHAR(50),
    date_created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    date_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on jurisdiction for faster queries
CREATE INDEX idx_companies_jurisdiction ON companies(jurisdiction);

-- Create index on company_name for searching
CREATE INDEX idx_companies_name ON companies(company_name);

-- Create index on date_created for sorting
CREATE INDEX idx_companies_date_created ON companies(date_created);

COMMIT;
