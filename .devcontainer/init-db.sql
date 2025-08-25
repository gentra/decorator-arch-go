-- Initialize database with required extensions and setup
-- This script runs when the PostgreSQL container starts

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create a test database for running tests
CREATE DATABASE taskdb_test;

-- Grant permissions on test database
GRANT ALL PRIVILEGES ON DATABASE taskdb_test TO "user";