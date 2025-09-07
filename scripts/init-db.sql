-- Grant schema privileges to fitbyte user
GRANT ALL ON SCHEMA public TO fitbyte;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO fitbyte;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO fitbyte;

-- Set default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO fitbyte;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO fitbyte;
