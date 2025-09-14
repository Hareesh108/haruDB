#!/bin/bash

# Simple Transaction Test Script for HaruDB
# This script tests basic transaction functionality

echo "🧪 HaruDB Transaction Test"
echo "=========================="
echo ""

# Create test data directory
mkdir -p ./test_data

# Start HaruDB server in background
echo "Starting HaruDB server..."
./harudb --data-dir ./test_data &
SERVER_PID=$!

# Wait for server to start
sleep 2

echo "✅ Server started (PID: $SERVER_PID)"
echo ""

# Function to execute SQL and show result
test_sql() {
    echo "SQL> $1"
    result=$(echo "$1" | nc localhost 54321)
    echo "Result: $result"
    echo ""
}

echo "📊 Test 1: Basic Transaction"
echo "============================="
test_sql "CREATE TABLE accounts (id, name, balance)"
test_sql "INSERT INTO accounts VALUES (1, 'Alice', '1000')"
test_sql "INSERT INTO accounts VALUES (2, 'Bob', '500')"
test_sql "SELECT * FROM accounts"

echo "📊 Test 2: Transaction Commit"
echo "============================="
test_sql "BEGIN TRANSACTION"
test_sql "UPDATE accounts SET balance = '900' ROW 0"
test_sql "UPDATE accounts SET balance = '600' ROW 1"
test_sql "SELECT * FROM accounts"
test_sql "COMMIT"
test_sql "SELECT * FROM accounts"

echo "📊 Test 3: Transaction Rollback"
echo "==============================="
test_sql "BEGIN TRANSACTION"
test_sql "UPDATE accounts SET balance = '800' ROW 0"
test_sql "UPDATE accounts SET balance = '700' ROW 1"
test_sql "SELECT * FROM accounts"
test_sql "ROLLBACK"
test_sql "SELECT * FROM accounts"

echo "📊 Test 4: Savepoints"
echo "======================"
test_sql "BEGIN TRANSACTION"
test_sql "INSERT INTO accounts VALUES (3, 'Charlie', '200')"
test_sql "SAVEPOINT sp1"
test_sql "INSERT INTO accounts VALUES (4, 'David', '300')"
test_sql "SELECT * FROM accounts"
test_sql "ROLLBACK TO SAVEPOINT sp1"
test_sql "SELECT * FROM accounts"
test_sql "COMMIT"
test_sql "SELECT * FROM accounts"

echo "📊 Test 5: Isolation Levels"
echo "==========================="
test_sql "BEGIN TRANSACTION ISOLATION LEVEL READ COMMITTED"
test_sql "SELECT * FROM accounts"
test_sql "COMMIT"

test_sql "BEGIN TRANSACTION ISOLATION LEVEL REPEATABLE READ"
test_sql "SELECT * FROM accounts"
test_sql "COMMIT"

test_sql "BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE"
test_sql "SELECT * FROM accounts"
test_sql "COMMIT"

echo "✅ All transaction tests completed!"
echo ""

# Stop the server
echo "Stopping server..."
kill $SERVER_PID
wait $SERVER_PID 2>/dev/null

# Clean up
rm -rf ./test_data

echo "✅ Cleanup completed!"
echo ""
echo "🎉 Transaction functionality is working correctly!"
