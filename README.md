# Database Query Analyzer

Read-only MCP server for MySQL. ChatGPT can inspect the database, indexes, and `SELECT` plans through tools exposed at `http://localhost:9091/mcp`.

## Requirements

- Go 1.27 or newer
- A running MySQL server

## Install

```bash
git clone <repository-url>
cd database-query-analyzer
go mod download
```

Edit `config.yaml` and set the MySQL connection. Do not commit real passwords.

```yaml
server:
  port: 9091

db:
  mysql:
    host: localhost
    port: 3306
    user: root
    password: your-password
    database: your_database
```

Build and start:

```bash
go build -o database-query-analyzer ./cmd
./database-query-analyzer
```

Or run without building:

```bash
go run ./cmd
```

Check that the process is up:

```bash
curl http://localhost:9091/healthz
curl http://localhost:9091/readyz
```

The MCP endpoint is `http://localhost:9091/mcp` (streamable HTTP, JSON responses).

## Use it in ChatGPT

This app is an MCP server, not a classic ChatGPT plugin manifest. Connect it as a custom MCP app in ChatGPT.

ChatGPT’s servers cannot reach `localhost`. Expose the process with a public HTTPS URL (for example [ngrok](https://ngrok.com) or Cloudflare Tunnel) and use that URL as the connector endpoint.

1. Start this server and confirm `/healthz` returns success.
2. Publish port `9091` over HTTPS so the public path ends in `/mcp`. Example local tunnel:

   ```bash
   ngrok http 9091
   ```

   Connector URL: `https://<your-tunnel-host>/mcp`
3. In ChatGPT, open **Settings → Apps → Advanced settings** and turn on **Developer mode**.
4. Go to **Settings → Apps**, choose **Create app** (or **Add MCP server**), and enter:
   - **Name:** Database Query Analyzer
   - **MCP server URL:** `https://<your-tunnel-host>/mcp`
   - **Authentication:** None (the server does not check a bearer token; protect the tunnel and the MySQL user)
5. Create the app and enable it in a new chat.

Only `SELECT` statements are accepted by the explain and stats tools. The MySQL user should have read access to the target database and permission to run `EXPLAIN`.

## How to use

In a chat with the app enabled, ask in plain language. ChatGPT calls the tools below.

| Tool | What it does | Arguments |
| --- | --- | --- |
| `get_database_info` | Engine, version, and schemas. Optional table list. | `include_tables` (bool, optional) |
| `get_schema` | Columns for one table, or every table if omitted. | `table` (string, optional) |
| `get_indexes` | Index name, uniqueness, type, columns, cardinality. | `table` (string, optional) |
| `explain_query` | Optimizer plan (`EXPLAIN FORMAT=JSON`). Does not run the query. | `sql` (single `SELECT`) |
| `explain_analyze` | Runs the `SELECT` with `EXPLAIN ANALYZE` and returns real timings. | `sql` (single `SELECT`) |
| `get_query_stats` | Runs the `SELECT` and returns time, row count, both explain forms, and warnings. | `sql` (single `SELECT`) |

Example prompts:

- “What version is this database, and which tables does it have?”
- “Show the schema and indexes for `orders`.”
- “Explain this query without running it: `SELECT * FROM orders WHERE customer_id = 42`.”
- “Run EXPLAIN ANALYZE on that query and tell me what is slow.”

Prefer `explain_query` before `explain_analyze` or `get_query_stats` on large tables. Those two tools execute the statement.
