# HaruDB 🚀

**HaruDB** is a fully custom database built from scratch in Go, inspired by PostgreSQL and MySQL.  
It’s designed to be **client-server, TCP-based, and feature-rich**, supporting SQL-like commands, persistence, indexing, transactions, and more.

---

## ✨ Current Features (v0.0.2)

- TCP server with interactive REPL (like `psql`)
- In-memory storage engine
- `CREATE TABLE`, `INSERT`, `SELECT` support

---

## 📦 Planned Full Features (Roadmap)

| Feature                          | Status        |
|---------------------------------|---------------|
| Disk-based persistence           | 🔜 Planned    |
| Indexes & query optimization     | 🔜 Planned    |
| Transactions & ACID compliance  | 🔜 Planned    |
| SQL parser & more commands       | 🔜 Planned    |
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
Welcome to HaruDB v0.1 🎉
Type 'exit' to quit.

haruDB> CREATE TABLE users (id, name);
Table users created

haruDB> INSERT INTO users VALUES (1, 'Hareesh');
1 row inserted

haruDB> INSERT INTO users VALUES (2, 'Bhittam');
1 row inserted

haruDB> SELECT * FROM users;
id | name
1  | Hareesh
2  | Bhittam
```

---

## ⚡ Vision for HaruDB

HaruDB aims to become a **production-capable database** with features like:

- SQL compliance with rich query support
- Persistent storage with crash recovery
- ACID transactions & MVCC
- Indexing & fast query execution
- Concurrent access and connection pooling
- Cross-platform and container-ready deployment

Think of it as your **own open-source Postgres/MySQL clone in Go**.

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

HaruDB is currently **experimental** and intended for learning and experimentation.
Do **not use in production** until full persistence and transaction support are implemented.
