#!/bin/bash

# Advanced WHERE Clauses Demo Script for HaruDB
# This script demonstrates all supported WHERE clause features

echo "🚀 HaruDB Advanced WHERE Clauses Demo"
echo "======================================"
echo ""

# Start HaruDB server in background
echo "Starting HaruDB server..."
./harudb --data-dir ./data &
SERVER_PID=$!

# Wait for server to start
sleep 2

echo "Setting up test data..."
echo ""

# Create test table and insert data
{
    echo "CREATE TABLE employees (id, name, age, salary, department, status, email);"
    echo "INSERT INTO employees VALUES (1, 'John Doe', 25, 50000, 'Engineering', 'active', 'john@company.com');"
    echo "INSERT INTO employees VALUES (2, 'Jane Smith', 30, 60000, 'Marketing', 'active', 'jane@company.com');"
    echo "INSERT INTO employees VALUES (3, 'Bob Johnson', 35, 70000, 'Engineering', 'inactive', 'bob@company.com');"
    echo "INSERT INTO employees VALUES (4, 'Alice Brown', 28, 55000, 'Sales', 'active', 'alice@company.com');"
    echo "INSERT INTO employees VALUES (5, 'Charlie Wilson', 45, 80000, 'Engineering', 'active', 'charlie@company.com');"
    echo "INSERT INTO employees VALUES (6, 'Diana Prince', 22, 45000, 'Marketing', 'active', 'diana@company.com');"
    echo "INSERT INTO employees VALUES (7, 'Eve Davis', 32, 65000, 'Sales', 'active', 'eve@company.com');"
    echo "INSERT INTO employees VALUES (8, 'Frank Miller', 29, 52000, 'Engineering', 'inactive', 'frank@company.com');"
    echo ""
    echo "CREATE INDEX ON employees (age);"
    echo "CREATE INDEX ON employees (department);"
    echo "CREATE INDEX ON employees (status);"
    echo "CREATE INDEX ON employees (email);"
    echo ""
} | nc localhost 54321

echo ""
echo "==========================================="
echo "1. BASIC COMPARISON OPERATORS"
echo "==========================================="
echo ""

echo "🔍 Equality (=): Find employees aged exactly 25"
echo "SELECT * FROM employees WHERE age = 25;" | nc localhost 54321
echo ""

echo "🔍 Not Equals (!=): Find employees not in Engineering"
echo "SELECT * FROM employees WHERE department != 'Engineering';" | nc localhost 54321
echo ""

echo "🔍 Less Than (<): Find employees younger than 30"
echo "SELECT * FROM employees WHERE age < 30;" | nc localhost 54321
echo ""

echo "🔍 Greater Than (>): Find employees with salary > 60000"
echo "SELECT * FROM employees WHERE salary > 60000;" | nc localhost 54321
echo ""

echo "🔍 Less Than or Equal (<=): Find employees aged 28 or younger"
echo "SELECT * FROM employees WHERE age <= 28;" | nc localhost 54321
echo ""

echo "🔍 Greater Than or Equal (>=): Find employees with salary >= 60000"
echo "SELECT * FROM employees WHERE salary >= 60000;" | nc localhost 54321
echo ""

echo "==========================================="
echo "2. LIKE PATTERN MATCHING"
echo "==========================================="
echo ""

echo "🔍 Names starting with 'J'"
echo "SELECT * FROM employees WHERE name LIKE 'J%';" | nc localhost 54321
echo ""

echo "🔍 Names containing 'ohn'"
echo "SELECT * FROM employees WHERE name LIKE '%ohn%';" | nc localhost 54321
echo ""

echo "🔍 Email addresses ending with '@company.com'"
echo "SELECT * FROM employees WHERE email LIKE '%@company.com';" | nc localhost 54321
echo ""

echo "🔍 Names with 'a' as second character"
echo "SELECT * FROM employees WHERE name LIKE '_a%';" | nc localhost 54321
echo ""

echo "==========================================="
echo "3. LOGICAL OPERATORS (AND/OR)"
echo "==========================================="
echo ""

echo "🔍 AND condition: Engineering employees older than 25"
echo "SELECT * FROM employees WHERE age > 25 AND department = 'Engineering';" | nc localhost 54321
echo ""

echo "🔍 OR condition: Employees younger than 25 or older than 40"
echo "SELECT * FROM employees WHERE age < 25 OR age > 40;" | nc localhost 54321
echo ""

echo "🔍 Multiple AND conditions: Active employees aged 25-35"
echo "SELECT * FROM employees WHERE age >= 25 AND age <= 35 AND status = 'active';" | nc localhost 54321
echo ""

echo "🔍 Multiple OR conditions: Marketing or Sales employees"
echo "SELECT * FROM employees WHERE department = 'Marketing' OR department = 'Sales';" | nc localhost 54321
echo ""

echo "==========================================="
echo "4. COMPLEX COMBINATIONS"
echo "==========================================="
echo ""

echo "🔍 AND with OR: Engineering employees who are either over 30 or have salary > 60000"
echo "SELECT * FROM employees WHERE department = 'Engineering' AND (age > 30 OR salary > 60000);" | nc localhost 54321
echo ""

echo "🔍 OR with AND: Young or old employees who are active"
echo "SELECT * FROM employees WHERE (age < 25 OR age > 40) AND status = 'active';" | nc localhost 54321
echo ""

echo "🔍 Complex salary and department filtering"
echo "SELECT * FROM employees WHERE salary >= 55000 AND (department = 'Engineering' OR department = 'Sales');" | nc localhost 54321
echo ""

echo "==========================================="
echo "5. EDGE CASES"
echo "==========================================="
echo ""

echo "🔍 No matches (should return empty)"
echo "SELECT * FROM employees WHERE age > 100;" | nc localhost 54321
echo ""

echo "🔍 All matches"
echo "SELECT * FROM employees WHERE age > 0;" | nc localhost 54321
echo ""

echo "🔍 String comparison (lexicographic)"
echo "SELECT * FROM employees WHERE name > 'M';" | nc localhost 54321
echo ""

echo "==========================================="
echo "6. PERFORMANCE TESTING WITH INDEXES"
echo "==========================================="
echo ""

echo "🔍 These queries should use indexes for better performance"
echo "SELECT * FROM employees WHERE age = 30;" | nc localhost 54321
echo ""

echo "SELECT * FROM employees WHERE department = 'Engineering';" | nc localhost 54321
echo ""

echo "SELECT * FROM employees WHERE status = 'active';" | nc localhost 54321
echo ""

echo "==========================================="
echo "7. VERIFICATION QUERIES"
echo "==========================================="
echo ""

echo "🔍 Count total employees"
echo "SELECT * FROM employees;" | nc localhost 54321
echo ""

echo "🔍 Count by department"
echo "SELECT * FROM employees WHERE department = 'Engineering';" | nc localhost 54321
echo ""

echo "SELECT * FROM employees WHERE department = 'Marketing';" | nc localhost 54321
echo ""

echo "SELECT * FROM employees WHERE department = 'Sales';" | nc localhost 54321
echo ""

echo "🔍 Count by status"
echo "SELECT * FROM employees WHERE status = 'active';" | nc localhost 54321
echo ""

echo "SELECT * FROM employees WHERE status = 'inactive';" | nc localhost 54321
echo ""

echo "==========================================="
echo "Demo completed! 🎉"
echo "==========================================="

# Clean up
echo "Cleaning up..."
kill $SERVER_PID 2>/dev/null
rm -rf ./demo_data
echo "Done!"
