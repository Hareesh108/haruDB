# HaruDB 🚀

**HaruDB** is a modern, fully custom database built from scratch in Go, inspired by PostgreSQL and SQLite.
It's designed to be **client-server, TCP-based, and feature-rich**, supporting SQL-like commands, persistence, crash recovery, and ACID compliance.

---

## ✨ Current Features (v0.0.5)

### 🏗️ **Core Architecture**

- **TCP-based client-server** architecture with interactive REPL (like `psql`)
- **Write-Ahead Logging (WAL)** for crash recovery and data durability
- **Atomic operations** with proper error handling and rollback
- **Persistent storage** with JSON-based table files
- **Memory-first design** with disk persistence

### 📊 **SQL Operations**

- **Data Definition Language (DDL)**:
  - `CREATE TABLE` - Create tables with custom schemas
  - `DROP TABLE` - Remove tables and associated data
- **Data Manipulation Language (DML)**:
  - `INSERT` - Add new rows to tables
  - `SELECT` - Query and display table data
  - `UPDATE` - Modify existing rows by index
  - `DELETE` - Remove rows by index

### 🚀 **Indexes & Query Optimization**

- `CREATE INDEX ON <table> (<column>)` - Build an in-memory hash index
- `SELECT ... WHERE <column> = 'value'` - Equality filters use indexes when available
- Index metadata persisted; indexes are rebuilt on startup

### 🔒 **Transactions & ACID Compliance**

- **ACID Transactions** – Full atomicity, consistency, isolation, and durability demonstrated in HaruDB.
- `BEGIN TRANSACTION` – Start a new transaction:

  ```sql
  BEGIN TRANSACTION;
  ```

- `COMMIT` – Commit all changes in the current transaction:

  ```sql
  COMMIT;
  ```

- `ROLLBACK` – Rollback all changes in the current transaction:

  ```sql
  ROLLBACK;
  ```

- `SAVEPOINT name` – Create savepoints within transactions:

  ```sql
  SAVEPOINT sp1;
  ```

- `ROLLBACK TO SAVEPOINT name` – Rollback to a specific savepoint:

  ```sql
  ROLLBACK TO SAVEPOINT sp1;
  ```

- **Isolation Levels** – Control how transactions interact:

  ```sql
  BEGIN TRANSACTION ISOLATION LEVEL READ COMMITTED;
  BEGIN TRANSACTION ISOLATION LEVEL REPEATABLE READ;
  BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE;
  ```

- **Multi-table Operations** – Transactions spanning multiple tables:

  ```sql
  INSERT INTO orders VALUES (1, 1, 'Laptop', '1', '999.99');
  UPDATE inventory SET stock = '9' ROW 0;
  COMMIT;
  ```

- **Error Handling in Transactions** – Safely handle invalid operations:

  ```sql
  BEGIN TRANSACTION;
  UPDATE inventory SET stock = '5' ROW 10; -- invalid row
  ROLLBACK;
  ```

- **Large Transaction Performance** – Efficiently handle bulk inserts:

  ```sql
  BEGIN TRANSACTION;
  INSERT INTO logs VALUES (1, '2024-01-01', 'Log entry 1');
  INSERT INTO logs VALUES (2, '2024-01-01', 'Log entry 2');
  COMMIT;
  ```

- **WAL Integration** – All transaction operations logged for crash recovery.

---

### 🔒 **Data Integrity & Recovery**

- **Write-Ahead Logging (WAL)** ensures all changes are logged before being applied
- **Crash recovery** - Automatic WAL replay on startup
- **Atomic writes** - Changes are either fully applied or not at all
- **Data consistency** - WAL ensures database state integrity

### 🔐 **Authentication & Security**

- **User Management** - Create, delete, and manage database users
- **Role-based Access Control** - Admin, User, and ReadOnly roles
- **Session Management** - Secure session handling with automatic cleanup
- **TLS Encryption** - Optional TLS support for secure connections
- **Password Security** - SHA-256 password hashing

### 💾 **Backup & Restore**

- **Full Database Backup** - Compressed tar.gz backups with metadata
- **Point-in-time Recovery** - Restore from any backup
- **Backup Information** - Detailed backup metadata and statistics
- **Backup Management** - List and manage backup files

---

## 🌟 Real-World Examples

HaruDB is designed to handle complex, production-ready applications. We provide comprehensive examples that demonstrate the full capabilities of the database:

### 🏦 Banking System Example
- **Complete financial system** with customers, accounts, transactions, and loans
- **Multi-table relationships** with foreign keys and referential integrity
- **Financial calculations** with decimal precision and audit trails
- **Performance optimization** with strategic indexing for high-volume transactions
- **Sample data**: 5,000+ customers, 10,000+ accounts, 100,000+ transactions

### 🍕 Food Ordering App Example (Swiggy/Zomato Style)
- **Full e-commerce platform** with restaurants, menus, orders, and delivery tracking
- **Complex relationships** between users, restaurants, orders, and delivery persons
- **Real-time tracking** of delivery status and order management
- **Review and rating system** for quality control and customer feedback
- **Sample data**: 50,000+ users, 1,000+ restaurants, 100,000+ orders

### 📁 Example Files
- `examples/banking_system.sql` - Complete banking system schema and data
- `examples/food_ordering_app.sql` - Full food delivery platform schema and data
- `examples/load_examples.sh` - Interactive script to load examples
- `docs/guides/examples/` - Comprehensive documentation with architecture diagrams

### 🚀 Quick Start with Examples
```bash
# Start HaruDB server
./harudb --data-dir ./data

# Load examples (in another terminal)
cd examples
./load_examples.sh

# Connect and explore
./haru-cli
LOGIN admin admin123
SELECT * FROM customers;      # Banking example
SELECT * FROM restaurants;    # Food app example
```

---

## 📦 Planned Full Features (Roadmap)

| Feature                          | Status        |
|---------------------------------|---------------|
| Disk-based persistence           | ✅ **Implemented** |
| Write-Ahead Logging (WAL)        | ✅ **Implemented** |
| Crash recovery                   | ✅ **Implemented** |
| Basic SQL operations (CRUD)      | ✅ **Implemented** |
| Indexes & query optimization     | ✅ **Implemented** |
| Advanced WHERE clauses           | ✅ **Implemented**    |
| Transactions & ACID compliance   | ✅ **Implemented**    |
| Concurrency & locking            | ✅ **Implemented** |
| Custom wire protocol             | ✅ **Implemented** |
| CLI client (`haru-cli`)          | ✅ **Implemented** |
| Authentication & TLS             | ✅ **Implemented** |
| Multi-user support               | ✅ **Implemented** |
| Backup & restore                 | ✅ **Implemented** |
| Docker & Kubernetes deployment   | ✅ Ready      |

---

## 🐧 Linux / 🍎 macOS Installation (Native)

### 1️⃣ Install HaruDB

Run the following command:

```bash
curl -sSL https://raw.githubusercontent.com/Hareesh108/haruDB/main/scripts/install-harudb.sh | bash
```

### ❌ Uninstall HaruDB (Native)

To fully remove HaruDB, including **active server processes, binary, data, logs, and temp files**, run:

```bash
curl -sSL https://raw.githubusercontent.com/Hareesh108/haruDB/main/scripts/uninstall-harudb.sh | bash
```

## 🐳 Run HaruDB via Docker

Pull the image:

```bash
docker pull hareesh108/harudb:latest
````

Run the container:

```bash
docker run -p 54321:54321 hareesh108/harudb:latest
```

Server output:

```
🚀 HaruDB server started on port 54321
```

---

## 🔌 Connect to HaruDB

Use **Telnet** (basic) or later, the HaruDB CLI client:

```bash
telnet localhost 54321
```

Use **haru-cli** the HaruDB CLI client:

```bash
haru-cli
```

Example session:

```
Welcome to HaruDB v0.0.5 🎉
Type 'exit' to quit.

haruDB>
haruDB> CREATE TABLE users (id, name, email);
Table users created

haruDB> INSERT INTO users VALUES (1, 'Hareesh', 'hareesh@example.com');
1 row inserted

haruDB> INSERT INTO users VALUES (2, 'Bhittam', 'bhittam@example.com');
1 row inserted

haruDB> SELECT * FROM users;
id | name    | email
1  | Hareesh | hareesh@example.com
2  | Bhittam | bhittam@example.com

haruDB> UPDATE users SET name = 'Hareesh Updated' ROW 0;
1 row updated

haruDB> SELECT * FROM users;
id | name            | email
1  | Hareesh Updated | hareesh@example.com
2  | Bhittam         | bhittam@example.com

haruDB> DELETE FROM users ROW 1;
1 row deleted

haruDB> SELECT * FROM users;
id | name            | email
1  | Hareesh Updated | hareesh@example.com

haruDB> DROP TABLE users;
Table users dropped

haruDB> SELECT * FROM users;
Table users not found
```

### 🔒 **Transaction Example**

```sql
-- Create accounts table
CREATE TABLE accounts (id, name, balance);
INSERT INTO accounts VALUES (1, 'Alice', '1000');
INSERT INTO accounts VALUES (2, 'Bob', '500');

-- Begin transaction
BEGIN TRANSACTION;

-- Transfer money
UPDATE accounts SET balance = '900' ROW 0;  -- Alice: 1000 -> 900
UPDATE accounts SET balance = '600' ROW 1;  -- Bob: 500 -> 600

-- Commit transaction
COMMIT;

-- Verify changes
SELECT * FROM accounts;
```

---

## 📊 **Feature Comparison**

| Feature | HaruDB | SQLite | PostgreSQL |
|---------|--------|--------|------------|
| **Write-Ahead Logging** | ✅ | ✅ | ✅ |
| **Crash Recovery** | ✅ | ✅ | ✅ |
| **TCP Server** | ✅ | ❌ | ✅ |
| **JSON Storage** | ✅ | ❌ | ❌ |
| **Memory-First** | ✅ | ❌ | ❌ |
| **Go Native** | ✅ | ❌ | ❌ |
| **Docker Ready** | ✅ | ✅ | ✅ |
| **ACID Transactions** | ✅ | ✅ | ✅ |
| **Advanced Indexing** | ✅ | ✅ | ✅ |
| **Concurrent Access** | 🔜 | Limited | ✅ |

---

## 🏛️ **Technical Architecture**

### **Write-Ahead Logging (WAL)**

HaruDB implements a robust WAL system that ensures data durability and crash recovery:

- **Binary WAL Format**: Efficient storage with timestamps and operation metadata
- **Atomic Operations**: All changes are logged before being applied to data files
- **Crash Recovery**: Automatic WAL replay on startup restores database state
- **Checkpointing**: Periodic checkpoints mark successful data persistence
- **Thread-Safe**: Concurrent access protection with mutex locks

### **Storage Engine**

- **Memory-First Design**: Fast in-memory operations with disk persistence
- **JSON Persistence**: Human-readable table files (`.harudb` format)
- **Atomic Writes**: Temp file + rename pattern ensures data integrity
- **File System Sync**: Proper `fsync()` calls ensure data reaches disk

### **SQL Parser**

- **String-Based Parser**: Simple but effective command parsing
- **Error Handling**: Comprehensive validation and user-friendly error messages
- **Extensible Design**: Easy to add new SQL operations
- **Indexing & WHERE**: Supports `CREATE INDEX` and `SELECT ... WHERE col = 'value'`

---

## 🚀 **Quick Start Guide**

### **1. Start the Server**

```bash
# Using the binary (without TLS)
./harudb --data-dir ./data

# Using the binary (with TLS encryption)
./harudb --data-dir ./data --tls

# Using Docker
docker run -p 54321:54321 hareesh108/harudb:latest
```

### **2. Connect and Use**

```bash
# Connect via telnet
telnet localhost 54321

# Or use netcat
nc localhost 54321
```

### **3. Basic Operations**

```sql
-- Create a table
CREATE TABLE products (id, name, price);

-- Insert data
INSERT INTO products VALUES (1, 'Laptop', '999.99');
INSERT INTO products VALUES (2, 'Mouse', '29.99');

-- Query data
SELECT * FROM products;

-- Create an index and run indexed queries
CREATE INDEX ON products (name);
SELECT * FROM products WHERE name = 'Laptop';

-- Update data
UPDATE products SET price = '1099.99' ROW 0;

-- Delete data
DELETE FROM products ROW 1;

-- Drop table
DROP TABLE products;
```

### **4. End-to-End Example (Indexes & WHERE)**

```sql
CREATE TABLE users (id, name, email);
INSERT INTO users VALUES (1, 'Hareesh', 'hareesh@example.com');
INSERT INTO users VALUES (2, 'Bhittam', 'bhittam@example.com');
INSERT INTO users VALUES (3, 'Alice', 'alice@example.com');

-- Create an index on email
CREATE INDEX ON users (email);

-- Uses index
SELECT * FROM users WHERE email = 'bhittam@example.com';

-- Falls back to full scan (no index on name yet)
SELECT * FROM users WHERE name = 'Alice';

```

## Advanced WHERE Clauses

Your Bash demo shows that HaruDB can now handle **rich conditional queries** well beyond simple equality.
Key capabilities—illustrated by the commands in your script:

| Capability               | Example from your script                                                                   | What it means                                                                                   |
| ------------------------ | ------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------- |
| **Comparison operators** | `SELECT * FROM employees WHERE age > 25;`                                                  | `<`, `>`, `<=`, `>=`, `!=` work on numbers and strings.                                         |
| **Pattern matching**     | `SELECT * FROM employees WHERE name LIKE 'J%';`                                            | `LIKE`, `%` (wildcard), `_` (single char) for flexible text search.                             |
| **Logical operators**    | `SELECT * FROM employees WHERE age > 25 AND department = 'Engineering';`                   | Combine multiple conditions with `AND`, `OR`, and parentheses for grouping.                     |
| **Complex combinations** | `SELECT * FROM employees WHERE department = 'Engineering' AND (age > 30 OR salary > 60000);` | Mix nested logic for precise filtering.                                                         |
| **Edge cases**           | `SELECT * FROM employees WHERE age > 100;`                                                 | Returns empty sets gracefully, supports lexicographic string

---

### **4. End-to-End Example (Transactions)**

```sql
CREATE TABLE accounts (id, name, balance);
INSERT INTO accounts VALUES (1, 'Alice', '1000');
INSERT INTO accounts VALUES (2, 'Bob', '500');
SELECT * FROM accounts;

BEGIN TRANSACTION;
UPDATE accounts SET balance = '95500' ROW 0;
UPDATE accounts SET balance = '69900' ROW 1;
SELECT * FROM accounts;
COMMIT;
SELECT * FROM accounts;

BEGIN TRANSACTION;
UPDATE accounts SET balance = '800' ROW 0;
UPDATE accounts SET balance = '700' ROW 1;
SELECT * FROM accounts;
ROLLBACK;
SELECT * FROM accounts;

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

```

### **5. Authentication & User Management**

```sql
-- Login with default admin user
LOGIN admin admin123

-- Create new users
CREATE USER john mypassword USER
CREATE USER jane herpassword ADMIN
CREATE USER readonly readpass READONLY

-- List all users (admin only)
LIST USERS

-- Logout
LOGOUT
```

### **6. Backup & Restore**

```sql
-- Create a backup
BACKUP TO ./backups/my_backup.backup DESCRIPTION "Daily backup"

-- List available backups
LIST BACKUPS ./backups

-- Get backup information
BACKUP INFO ./backups/my_backup.backup

-- Restore from backup (admin only)
RESTORE FROM ./backups/my_backup.backup
```

## Advanced Transaction Features

HaruDB can now handle **full-fledged transactional operations** with ACID compliance, covering a wide range of scenarios from simple inserts to complex multi-table workflows.
Key capabilities—illustrated by the commands in your script:

| Capability                         | Example from your script                                                                                               | What it means                                                                                     |
| ---------------------------------- | ---------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| **Basic transaction workflow**     | `CREATE TABLE accounts (id, name, balance); INSERT INTO accounts VALUES (1, 'Alice', '1000'); SELECT * FROM accounts;` | Create tables, insert rows, and query data within a transactional context.                        |
| **Begin/Commit/Rollback**          | `BEGIN TRANSACTION; UPDATE accounts SET balance = '900' ROW 0; COMMIT;`                                                | Start a transaction, apply multiple operations, and commit to make changes permanent.             |
| **Transaction rollback**           | `BEGIN TRANSACTION; UPDATE accounts SET balance = '800' ROW 0; ROLLBACK;`                                              | Undo all operations within a transaction to maintain consistency.                                 |
| **Savepoints**                     | `SAVEPOINT sp1; INSERT INTO accounts VALUES (4, 'David', '300'); ROLLBACK TO SAVEPOINT sp1;`                           | Create intermediate checkpoints to selectively undo operations without rolling back everything.   |
| **Isolation levels**               | `BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE; SELECT * FROM accounts; COMMIT;`                                      | Control how concurrent transactions interact (Read Committed, Repeatable Read, Serializable).     |
| **Complex multi-table operations** | `BEGIN TRANSACTION; INSERT INTO orders VALUES (...); UPDATE inventory SET stock = '9' ROW 0; COMMIT;`                  | Support transactions spanning multiple tables with dependent operations.                          |
| **Error handling**                 | `BEGIN TRANSACTION; UPDATE inventory SET stock = '5' ROW 10; ROLLBACK;`                                                | Gracefully handle errors by rolling back failed transactions.                                     |
| **Large transaction performance**  | `BEGIN TRANSACTION; INSERT INTO logs VALUES (...); COMMIT;`                                                            | Efficiently handle bulk inserts or updates in a single transaction.                               |
| **ACID compliance demonstration**  | `-- Atomicity, Consistency, Isolation, Durability`                                                                     | Demonstrates that all transactional properties are supported, ensuring reliability and integrity. |
| **Cleanup & teardown**             | `BEGIN TRANSACTION; DROP TABLE accounts; COMMIT;`                                                                      | Safely remove tables and data at the end of a session.                                            |

---

## ⚡ Vision for HaruDB

HaruDB aims to become a **production-capable database** with features like:

- **Full SQL Compliance** with rich query support and advanced WHERE clauses
- **ACID Transactions** with Multi-Version Concurrency Control (MVCC)
- **Advanced Indexing** with B-tree and hash indexes for fast query execution
- **Concurrent Access** with proper locking and connection pooling
- **High Availability** with replication and clustering support
- **Performance Optimization** with query planning and execution optimization
- **Security Features** with authentication, authorization, and encryption
- **Cross-Platform** deployment with Docker, Kubernetes, and cloud-native support

Think of it as your **own open-source Postgres/MySQL clone in Go** - a modern, performant, and reliable database system.

---

## 🔧 **Development & Building**

### **Prerequisites**

- Go 1.24.0 or later
- Git

### **Build from Source**

```bash
# Clone the repository
git clone https://github.com/Hareesh108/haruDB.git
cd haruDB

# Build the binary
go build -o harudb ./cmd/server

# Run the server
./harudb --data-dir ./data
```

### **Project Structure**

```
haruDB/
├── cmd/server/          # Main server application
├── internal/
│   ├── parser/          # SQL parser and query engine
│   └── storage/         # Storage engine and WAL implementation
├── data/               # Database files directory
├── scripts/            # Installation and utility scripts
└── dist/               # Pre-built binaries
```

---

## 🐛 **Troubleshooting**

### **Common Issues**

**Port already in use:**

```bash
# Kill existing process
pkill harudb
# Or use a different port
./harudb --port 54322

# HaruDB automatically detects port conflicts and suggests solutions
./harudb --data-dir ./data
# Output: ⚠️  Port 54321 is already in use by: python3
#         Please stop the other service or use a different port
```

**WAL file corruption:**

```bash
# Remove WAL file to start fresh (data loss warning!)
rm data/wal.log
```

**Permission issues:**

```bash
# Ensure proper permissions
chmod +x harudb
chmod 755 data/
```

---

## 📖 Contributing

Contributions are welcome!

- Report issues
- Submit PRs for new features
- Help with testing, documentation, or CLI tools

---

## 👨‍💻 Author

**Hareesh Bhittam** – [GitHub](https://github.com/Hareesh108)

---

## ⚠️ Disclaimer

HaruDB is currently in **active development** and includes robust persistence and crash recovery features.
While the core functionality is stable, it's recommended for **development and testing environments**.
For production use, ensure thorough testing and consider the current feature limitations.
