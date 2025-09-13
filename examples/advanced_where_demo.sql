-- Advanced WHERE Clauses Demo for HaruDB
-- This script demonstrates all supported WHERE clause features

-- Create test table
CREATE TABLE employees (id, name, age, salary, department, status, email);

-- Insert sample data
INSERT INTO employees VALUES (1, 'John Doe', 25, 50000, 'Engineering', 'active', 'john@company.com');
INSERT INTO employees VALUES (2, 'Jane Smith', 30, 60000, 'Marketing', 'active', 'jane@company.com');
INSERT INTO employees VALUES (3, 'Bob Johnson', 35, 70000, 'Engineering', 'inactive', 'bob@company.com');
INSERT INTO employees VALUES (4, 'Alice Brown', 28, 55000, 'Sales', 'active', 'alice@company.com');
INSERT INTO employees VALUES (5, 'Charlie Wilson', 45, 80000, 'Engineering', 'active', 'charlie@company.com');
INSERT INTO employees VALUES (6, 'Diana Prince', 22, 45000, 'Marketing', 'active', 'diana@company.com');
INSERT INTO employees VALUES (7, 'Eve Davis', 32, 65000, 'Sales', 'active', 'eve@company.com');
INSERT INTO employees VALUES (8, 'Frank Miller', 29, 52000, 'Engineering', 'inactive', 'frank@company.com');

-- Create indexes for better performance
CREATE INDEX ON employees (age);
CREATE INDEX ON employees (department);
CREATE INDEX ON employees (status);
CREATE INDEX ON employees (email);

-- ===========================================
-- 1. BASIC COMPARISON OPERATORS
-- ===========================================

-- Equality (=)
SELECT * FROM employees WHERE age = 25;
-- Expected: John Doe (25 years old)

-- Not Equals (!= or <>)
SELECT * FROM employees WHERE department != 'Engineering';
-- Expected: Jane, Alice, Diana, Eve (non-Engineering employees)

-- Less Than (<)
SELECT * FROM employees WHERE age < 30;
-- Expected: John (25), Alice (28), Diana (22), Frank (29)

-- Greater Than (>)
SELECT * FROM employees WHERE salary > 60000;
-- Expected: Jane (60000), Bob (70000), Charlie (80000), Eve (65000)

-- Less Than or Equal (<=)
SELECT * FROM employees WHERE age <= 28;
-- Expected: John (25), Alice (28), Diana (22)

-- Greater Than or Equal (>=)
SELECT * FROM employees WHERE salary >= 60000;
-- Expected: Jane (60000), Bob (70000), Charlie (80000), Eve (65000)

-- ===========================================
-- 2. LIKE PATTERN MATCHING
-- ===========================================

-- Names starting with 'J'
SELECT * FROM employees WHERE name LIKE 'J%';
-- Expected: John Doe, Jane Smith

-- Names containing 'ohn'
SELECT * FROM employees WHERE name LIKE '%ohn%';
-- Expected: John Doe, Bob Johnson

-- Names with exactly 4 characters (using underscore)
SELECT * FROM employees WHERE name LIKE '____';
-- Expected: None (all names are longer)

-- Email addresses ending with '@company.com'
SELECT * FROM employees WHERE email LIKE '%@company.com';
-- Expected: All employees

-- Names with 'a' as second character
SELECT * FROM employees WHERE name LIKE '_a%';
-- Expected: Jane Smith, Diana Prince

-- ===========================================
-- 3. LOGICAL OPERATORS (AND/OR)
-- ===========================================

-- AND condition
SELECT * FROM employees WHERE age > 25 AND department = 'Engineering';
-- Expected: Bob Johnson (35), Charlie Wilson (45), Frank Miller (29)

-- OR condition
SELECT * FROM employees WHERE age < 25 OR age > 40;
-- Expected: Diana Prince (22), Charlie Wilson (45)

-- Multiple AND conditions
SELECT * FROM employees WHERE age >= 25 AND age <= 35 AND status = 'active';
-- Expected: Jane Smith (30), Alice Brown (28), Eve Davis (32)

-- Multiple OR conditions
SELECT * FROM employees WHERE department = 'Marketing' OR department = 'Sales';
-- Expected: Jane Smith, Alice Brown, Diana Prince, Eve Davis

-- ===========================================
-- 4. COMPLEX COMBINATIONS
-- ===========================================

-- AND with OR
SELECT * FROM employees WHERE department = 'Engineering' AND (age > 30 OR salary > 60000);
-- Expected: Bob Johnson (35, 70000), Charlie Wilson (45, 80000)

-- OR with AND
SELECT * FROM employees WHERE (age < 25 OR age > 40) AND status = 'active';
-- Expected: Diana Prince (22), Charlie Wilson (45)

-- Complex salary and department filtering
SELECT * FROM employees WHERE salary >= 55000 AND (department = 'Engineering' OR department = 'Sales');
-- Expected: Bob Johnson (70000, Engineering), Alice Brown (55000, Sales), Charlie Wilson (80000, Engineering), Eve Davis (65000, Sales)

-- ===========================================
-- 5. EDGE CASES AND SPECIAL SCENARIOS
-- ===========================================

-- No matches (should return empty)
SELECT * FROM employees WHERE age > 100;
-- Expected: (no rows)

-- All matches
SELECT * FROM employees WHERE age > 0;
-- Expected: All employees

-- String comparison (lexicographic)
SELECT * FROM employees WHERE name > 'M';
-- Expected: Names starting with M-Z: Bob Johnson, Charlie Wilson, Diana Prince, Eve Davis, Frank Miller

-- Numeric comparison with string values
SELECT * FROM employees WHERE id > 5;
-- Expected: Diana Prince (6), Eve Davis (7), Frank Miller (8)

-- ===========================================
-- 6. PERFORMANCE TESTING WITH INDEXES
-- ===========================================

-- These queries should use indexes for better performance
SELECT * FROM employees WHERE age = 30;
-- Uses age index

SELECT * FROM employees WHERE department = 'Engineering';
-- Uses department index

SELECT * FROM employees WHERE status = 'active';
-- Uses status index

-- ===========================================
-- 7. QUOTED STRINGS AND SPECIAL CHARACTERS
-- ===========================================

-- Single quotes
SELECT * FROM employees WHERE name = 'John Doe';

-- Double quotes
SELECT * FROM employees WHERE name = "Jane Smith";

-- Names with spaces
SELECT * FROM employees WHERE name LIKE '% %';
-- Expected: All employees (all have spaces in names)

-- ===========================================
-- 8. VERIFICATION QUERIES
-- ===========================================

-- Count total employees
SELECT * FROM employees;

-- Count by department
SELECT * FROM employees WHERE department = 'Engineering';
SELECT * FROM employees WHERE department = 'Marketing';
SELECT * FROM employees WHERE department = 'Sales';

-- Count by status
SELECT * FROM employees WHERE status = 'active';
SELECT * FROM employees WHERE status = 'inactive';

-- Age ranges
SELECT * FROM employees WHERE age BETWEEN 25 AND 35;
-- Note: BETWEEN not implemented yet, but equivalent to:
SELECT * FROM employees WHERE age >= 25 AND age <= 35;

-- ===========================================
-- CLEANUP (optional)
-- ===========================================

-- DROP TABLE employees;
