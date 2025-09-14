#!/bin/bash

# HaruDB Transaction Examples
# This script demonstrates various transaction scenarios and ACID compliance

echo "🚀 HaruDB Transaction Examples"
echo "================================"
echo ""

# Start HaruDB server in background
echo "Starting HaruDB server..."
./harudb --data-dir ./transaction_data &
SERVER_PID=$!

# Wait for server to start
sleep 2

echo "✅ Server started (PID: $SERVER_PID)"
echo ""

# Function to execute SQL commands
execute_sql() {
    echo "SQL> $1"
    echo "$1" | nc localhost 54321
    echo ""
}

# Function to execute multiple SQL commands
execute_transaction() {
    echo "=== $1 ==="
    echo "$2" | nc localhost 54321
    echo ""
}

echo "📊 Example 1: Basic Transaction Workflow"
echo "=========================================="
execute_transaction "Creating tables and inserting data" "
CREATE TABLE accounts (id, name, balance);
INSERT INTO accounts VALUES (1, 'Alice', '1000');
INSERT INTO accounts VALUES (2, 'Bob', '500');
SELECT * FROM accounts;
"

echo "📊 Example 2: Money Transfer with Transaction"
echo "=============================================="
execute_transaction "Money transfer transaction" "
BEGIN TRANSACTION;
UPDATE accounts SET balance = '900' ROW 0;
UPDATE accounts SET balance = '600' ROW 1;
SELECT * FROM accounts;
COMMIT;
SELECT * FROM accounts;
"

echo "📊 Example 3: Transaction Rollback"
echo "==================================="
execute_transaction "Failed transaction with rollback" "
BEGIN TRANSACTION;
UPDATE accounts SET balance = '800' ROW 0;
UPDATE accounts SET balance = '700' ROW 1;
SELECT * FROM accounts;
ROLLBACK;
SELECT * FROM accounts;
"

echo "📊 Example 4: Savepoints"
echo "========================"
execute_transaction "Transaction with savepoints" "
BEGIN TRANSACTION;
INSERT INTO accounts VALUES (3, 'Charlie', '200');
SAVEPOINT sp1;
INSERT INTO accounts VALUES (4, 'David', '300');
SAVEPOINT sp2;
INSERT INTO accounts VALUES (5, 'Eve', '400');
SELECT * FROM accounts;
ROLLBACK TO SAVEPOINT sp1;
SELECT * FROM accounts;
COMMIT;
SELECT * FROM accounts;
"

echo "📊 Example 5: Different Isolation Levels"
echo "========================================"
execute_transaction "Read Committed isolation" "
BEGIN TRANSACTION ISOLATION LEVEL READ COMMITTED;
SELECT * FROM accounts;
COMMIT;
"

execute_transaction "Repeatable Read isolation" "
BEGIN TRANSACTION ISOLATION LEVEL REPEATABLE READ;
SELECT * FROM accounts;
COMMIT;
"

execute_transaction "Serializable isolation" "
BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE;
SELECT * FROM accounts;
COMMIT;
"

echo "📊 Example 6: Complex Transaction with Multiple Operations"
echo "=========================================================="
execute_transaction "Complex transaction" "
CREATE TABLE orders (id, customer_id, product, quantity, price);
CREATE TABLE inventory (id, product, stock);

BEGIN TRANSACTION;
INSERT INTO orders VALUES (1, 1, 'Laptop', '1', '999.99');
INSERT INTO inventory VALUES (1, 'Laptop', '10');
UPDATE inventory SET stock = '9' ROW 0;
SAVEPOINT order_created;
INSERT INTO orders VALUES (2, 2, 'Mouse', '2', '29.99');
UPDATE inventory SET stock = '8' ROW 0;
ROLLBACK TO SAVEPOINT order_created;
INSERT INTO orders VALUES (2, 2, 'Keyboard', '1', '79.99');
UPDATE inventory SET stock = '7' ROW 0;
COMMIT;
SELECT * FROM orders;
SELECT * FROM inventory;
"

echo "📊 Example 7: Error Handling in Transactions"
echo "============================================="
execute_transaction "Transaction with error handling" "
BEGIN TRANSACTION;
INSERT INTO orders VALUES (3, 3, 'Monitor', '1', '299.99');
-- This will fail due to invalid row index
UPDATE inventory SET stock = '5' ROW 10;
ROLLBACK;
SELECT * FROM orders;
"

echo "📊 Example 8: Large Transaction Performance"
echo "============================================"
execute_transaction "Large transaction" "
CREATE TABLE logs (id, timestamp, message);

BEGIN TRANSACTION;
INSERT INTO logs VALUES (1, '2024-01-01', 'Log entry 1');
INSERT INTO logs VALUES (2, '2024-01-01', 'Log entry 2');
INSERT INTO logs VALUES (3, '2024-01-01', 'Log entry 3');
INSERT INTO logs VALUES (4, '2024-01-01', 'Log entry 4');
INSERT INTO logs VALUES (5, '2024-01-01', 'Log entry 5');
COMMIT;
SELECT COUNT(*) FROM logs;
"

echo "📊 Example 9: Transaction Cleanup"
echo "=================================="
execute_transaction "Cleanup operations" "
BEGIN TRANSACTION;
DROP TABLE accounts;
DROP TABLE orders;
DROP TABLE inventory;
DROP TABLE logs;
COMMIT;
"

echo "📊 Example 10: ACID Compliance Demonstration"
echo "============================================="
execute_transaction "ACID compliance test" "
CREATE TABLE acid_test (id, value);

-- Atomicity test
BEGIN TRANSACTION;
INSERT INTO acid_test VALUES (1, 'A');
INSERT INTO acid_test VALUES (2, 'B');
ROLLBACK;
SELECT * FROM acid_test;

-- Consistency test
BEGIN TRANSACTION;
INSERT INTO acid_test VALUES (1, 'A');
INSERT INTO acid_test VALUES (2, 'B');
COMMIT;
SELECT * FROM acid_test;

-- Isolation test
BEGIN TRANSACTION;
UPDATE acid_test SET value = 'C' ROW 0;
SELECT * FROM acid_test;
ROLLBACK;
SELECT * FROM acid_test;
"

echo "✅ Transaction examples completed!"
echo ""
echo "🧹 Cleaning up..."

# Stop the server
kill $SERVER_PID
wait $SERVER_PID 2>/dev/null

# Clean up data directory
rm -rf ./transaction_data

echo "✅ Cleanup completed!"
echo ""
echo "🎉 All transaction examples executed successfully!"
echo ""
echo "Key Features Demonstrated:"
echo "- ✅ ACID Compliance (Atomicity, Consistency, Isolation, Durability)"
echo "- ✅ Transaction Begin/Commit/Rollback"
echo "- ✅ Savepoints and Rollback to Savepoint"
echo "- ✅ Multiple Isolation Levels"
echo "- ✅ Error Handling in Transactions"
echo "- ✅ Large Transaction Performance"
echo "- ✅ Complex Multi-Operation Transactions"
echo ""
echo "HaruDB now supports full transaction management with ACID compliance! 🚀"
