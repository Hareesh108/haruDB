# HaruDB 🚀

**HaruDB** is a modern, fully custom database built from scratch in Go, inspired by PostgreSQL and SQLite.
It's designed to be **client-server, TCP-based, and feature-rich**, supporting SQL-like commands, persistence, crash recovery, and ACID compliance.

---

## ✨ Current Features (v0.0.3)

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

### 🔒 **Data Integrity & Recovery**

- **Write-Ahead Logging (WAL)** ensures all changes are logged before being applied
- **Crash recovery** - Automatic WAL replay on startup
- **Atomic writes** - Changes are either fully applied or not at all
- **Data consistency** - WAL ensures database state integrity

---

## 📦 Planned Full Features (Roadmap)

| Feature                          | Status        |
|---------------------------------|---------------|
| Disk-based persistence           | ✅ **Implemented** |
| Write-Ahead Logging (WAL)        | ✅ **Implemented** |
| Crash recovery                   | ✅ **Implemented** |
| Basic SQL operations (CRUD)      | ✅ **Implemented** |
| Indexes & query optimization     | 🔜 Planned    |
| Advanced WHERE clauses           | 🔜 Planned    |
| Transactions & ACID compliance  | 🔜 Planned    |
| Concurrency & locking            | 🔜 Planned    |
| Custom wire protocol             | 🔜 Planned    |
| CLI client (`haru-cli`)          | 🔜 Planned    |
| Authentication & TLS             | 🔜 Planned    |
| Multi-user support               | 🔜 Planned    |
| Backup & restore                 | 🔜 Planned    |
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

Example session:

```
Welcome to HaruDB v0.0.3 🎉
Type 'exit' to quit.

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
| **ACID Transactions** | 🔜 | ✅ | ✅ |
| **Advanced Indexing** | 🔜 | ✅ | ✅ |
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

---

## 🚀 **Quick Start Guide**

### **1. Start the Server**

```bash
# Using the binary
./harudb --data-dir ./data

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

-- Update data
UPDATE products SET price = '1099.99' ROW 0;

-- Delete data
DELETE FROM products ROW 1;

-- Drop table
DROP TABLE products;
```

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

**Current Status**: ✅ Persistence & WAL | 🔜 Full ACID Transactions | 🔜 Advanced Querying
