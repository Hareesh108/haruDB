-- HaruDB Transaction Examples SQL Script
-- This file contains SQL commands demonstrating transaction functionality

-- Example 1: Basic Transaction Workflow
-- ======================================

-- Create tables and insert initial data
CREATE TABLE accounts (id, name, balance);
INSERT INTO accounts VALUES (1, 'Alice', '1000');
INSERT INTO accounts VALUES (2, 'Bob', '500');
SELECT * FROM accounts;

-- Example 2: Money Transfer with Transaction
-- ==========================================

-- Begin transaction
BEGIN TRANSACTION;

-- Transfer $100 from Alice to Bob
UPDATE accounts SET balance = '900' ROW 0;  -- Alice: 1000 -> 900
UPDATE accounts SET balance = '600' ROW 1;  -- Bob: 500 -> 600

-- Check intermediate state
SELECT * FROM accounts;

-- Commit the transaction
COMMIT;

-- Verify final state
SELECT * FROM accounts;

-- Example 3: Transaction Rollback
-- ===============================

-- Begin transaction
BEGIN TRANSACTION;

-- Attempt to transfer more money than Alice has
UPDATE accounts SET balance = '800' ROW 0;  -- Alice: 900 -> 800
UPDATE accounts SET balance = '700' ROW 1;  -- Bob: 600 -> 700

-- Check intermediate state
SELECT * FROM accounts;

-- Rollback the transaction (simulating error condition)
ROLLBACK;

-- Verify original state is restored
SELECT * FROM accounts;

-- Example 4: Savepoints
-- ====================

-- Begin transaction
BEGIN TRANSACTION;

-- Add new account
INSERT INTO accounts VALUES (3, 'Charlie', '200');

-- Create savepoint
SAVEPOINT sp1;

-- Add more accounts
INSERT INTO accounts VALUES (4, 'David', '300');
SAVEPOINT sp2;
INSERT INTO accounts VALUES (5, 'Eve', '400');

-- Check current state
SELECT * FROM accounts;

-- Rollback to first savepoint
ROLLBACK TO SAVEPOINT sp1;

-- Check state after rollback
SELECT * FROM accounts;

-- Commit transaction
COMMIT;

-- Verify final state
SELECT * FROM accounts;

-- Example 5: Different Isolation Levels
-- ====================================

-- Read Committed isolation
BEGIN TRANSACTION ISOLATION LEVEL READ COMMITTED;
SELECT * FROM accounts;
COMMIT;

-- Repeatable Read isolation
BEGIN TRANSACTION ISOLATION LEVEL REPEATABLE READ;
SELECT * FROM accounts;
COMMIT;

-- Serializable isolation
BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE;
SELECT * FROM accounts;
COMMIT;

-- Example 6: Complex Transaction with Multiple Operations
-- =======================================================

-- Create additional tables
CREATE TABLE orders (id, customer_id, product, quantity, price);
CREATE TABLE inventory (id, product, stock);

-- Begin complex transaction
BEGIN TRANSACTION;

-- Create order
INSERT INTO orders VALUES (1, 1, 'Laptop', '1', '999.99');
INSERT INTO inventory VALUES (1, 'Laptop', '10');

-- Update inventory
UPDATE inventory SET stock = '9' ROW 0;

-- Create savepoint
SAVEPOINT order_created;

-- Add another order
INSERT INTO orders VALUES (2, 2, 'Mouse', '2', '29.99');
UPDATE inventory SET stock = '8' ROW 0;

-- Rollback to savepoint (simulating order cancellation)
ROLLBACK TO SAVEPOINT order_created;

-- Add different order
INSERT INTO orders VALUES (2, 2, 'Keyboard', '1', '79.99');
UPDATE inventory SET stock = '7' ROW 0;

-- Commit transaction
COMMIT;

-- Verify final state
SELECT * FROM orders;
SELECT * FROM inventory;

-- Example 7: Error Handling in Transactions
-- ==========================================

-- Begin transaction
BEGIN TRANSACTION;

-- Add order
INSERT INTO orders VALUES (3, 3, 'Monitor', '1', '299.99');

-- This will fail due to invalid row index (simulating error)
UPDATE inventory SET stock = '5' ROW 10;

-- Rollback due to error
ROLLBACK;

-- Verify no changes were made
SELECT * FROM orders;

-- Example 8: Large Transaction Performance
-- ========================================

-- Create logs table
CREATE TABLE logs (id, timestamp, message);

-- Begin large transaction
BEGIN TRANSACTION;

-- Insert multiple log entries
INSERT INTO logs VALUES (1, '2024-01-01', 'Log entry 1');
INSERT INTO logs VALUES (2, '2024-01-01', 'Log entry 2');
INSERT INTO logs VALUES (3, '2024-01-01', 'Log entry 3');
INSERT INTO logs VALUES (4, '2024-01-01', 'Log entry 4');
INSERT INTO logs VALUES (5, '2024-01-01', 'Log entry 5');

-- Commit transaction
COMMIT;

-- Verify all entries were inserted
SELECT COUNT(*) FROM logs;

-- Example 9: ACID Compliance Demonstration
-- ==========================================

-- Atomicity test
BEGIN TRANSACTION;
INSERT INTO acid_test VALUES (1, 'A');
INSERT INTO acid_test VALUES (2, 'B');
ROLLBACK;
-- Should show no rows

-- Consistency test
BEGIN TRANSACTION;
INSERT INTO acid_test VALUES (1, 'A');
INSERT INTO acid_test VALUES (2, 'B');
COMMIT;
-- Should show 2 rows

-- Isolation test
BEGIN TRANSACTION;
UPDATE acid_test SET value = 'C' ROW 0;
-- Check intermediate state
SELECT * FROM acid_test;
ROLLBACK;
-- Should show original values

-- Durability test (implicit - data persists after commit)
SELECT * FROM acid_test;

-- Example 10: Transaction Cleanup
-- ================================

-- Begin cleanup transaction
BEGIN TRANSACTION;

-- Drop all tables
DROP TABLE accounts;
DROP TABLE orders;
DROP TABLE inventory;
DROP TABLE logs;
DROP TABLE acid_test;

-- Commit cleanup
COMMIT;

-- Verify cleanup
SELECT * FROM accounts;  -- Should show "Table accounts not found"
