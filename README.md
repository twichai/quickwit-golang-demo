Below is a clean `README.md` for your logging stack:

---

# Logging Stack (Go + Fluent Bit + Quickwit + Grafana)

This project demonstrates the real-time data recording process for Golang project:

Go Application → Fluent Bit → Quickwit → Grafana

---

## Architecture

- **Go App**
  Writes structured JSON logs to `app.log`

- **Fluent Bit**
  Tails log file and forwards logs to Quickwit

- **Quickwit**
  Stores and indexes logs (Elasticsearch-compatible)

- **Grafana**
  Visualizes logs and metrics

---

## Project Structure

```
.
├── docker-compose.yml
├── fluent-bit/
│   └── fluent-bit.conf
├── quickwit/
│   └── index-config.yaml
├── logs/
│   └── app.log
└── go-app/
```

---

## Run the Stack

### 1. Start services

```bash
docker compose up -d
```

Services:

- Quickwit → [http://localhost:7280](http://localhost:7280)
- Grafana → [http://localhost:3000](http://localhost:3000)

---

## Quickwit Index Setup

Example `index-config.yaml`:

```yaml
index_id: stackoverflow

doc_mapping:
  timestamp_field: creationDate
  field_mappings:
    - name: creationDate
      type: datetime
    - name: level
      type: text
    - name: msg
      type: text
    - name: username
      type: text

indexing_settings:
  commit_timeout_secs: 5
```

Create index:

```bash
curl -X POST http://localhost:7280/api/v1/indexes \
  -H "Content-Type: application/yaml" \
  --data-binary @index-config.yaml
```

---

## Fluent Bit Configuration

Example:

```
[INPUT]
    Name              tail
    Path              /app/logs/app.log
    Parser            json
    Tag               myapp
    Refresh_Interval  1
    Read_from_Head    true
    DB                /fluent-bit/tail.db

[OUTPUT]
    Name              http
    Match             *
    Host              quickwit
    Port              7280
    URI               /api/v1/stackoverflow/ingest
    Format            json
```

---

## Access Grafana

1. Open:
   [http://localhost:3000](http://localhost:3000)

2. Login:
   user: `admin`
   pass: `admin`

3. Add Data Source:
   - Type: Elasticsearch
   - URL: `http://quickwit:7280`
   - Index: `stackoverflow`
   - Time field: `creationDate`

---

## Example Queries (Grafana)

Show all logs:

```
*
```

Filter by level:

```
level:ERROR
```

Filter by user:

```
username:john
```

---

## Common Issues

### Logs not appearing

- Ensure Quickwit index exists
- Check `commit_timeout_secs`
- Confirm `creationDate` is mapped as datetime
- Verify time picker range in Grafana

### Logs delayed in Docker (macOS)

Use:

```
./logs:/app/logs:delegated
```

Or run Fluent Bit locally.

---

## Development Tips

For stable real-time logs on macOS:

Option A:

- Run Go locally
- Run Fluent Bit locally
- Quickwit + Grafana in Docker

Option B:

- Run everything inside Docker

Avoid mixing local file writes with Docker bind mounts.

---

This stack provides near real-time log ingestion (~5 seconds commit window) and searchable dashboards.
