-- Testing sql script to check the state of the database
-- Use python test_db_state.py to run this script
-- With: python test_db_state.py

-- Show all tables and their counts in each database
USE ednevnik_workspace;
SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'ednevnik_workspace';
SHOW TABLES;

-- USE ednevnik_primary;
-- SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'ednevnik_primary';
-- SHOW TABLES;

-- USE ednevnik_secondary;
-- SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'ednevnik_secondary';
-- SHOW TABLES;
