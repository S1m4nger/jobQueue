# Job Queue + Worker API

This project implements a simplified job queue and worker API, a common backend pattern for asynchronous task processing.

## Idea

The service accepts tasks via an API, places them into a queue, processes them with workers, stores execution status, and allows checking the results.

Examples of tasks:
- Sending emails (mock)
- File processing
- Webhook requests (without retries initially)

## Stack

- Go 1.21+
- net/http
- SQLite
- In-memory queue (channel)
- context

## Project Structure

```
job-queue/
├── cmd/
│   └── app/
│       └── main.go
│
├── internal/
│   ├── http/
│   │   ├── handlers/
│   │   │   └── jobs.go
│   │   └── router.go
│   │
│   ├── service/
│   │   └── jobs.go
│   │
│   ├── repository/
│   │   └── sqlite/
│   │       └── jobs.go
│   │
│   ├── domain/
│   │   └── job.go
│   │
│   └── queue/
│       └── memory.go
│
├── migrations/
│   └── 000001_create_jobs_table.up.sql
├── go.mod
├── go.sum
└── README.md
```

## Minimal Functionality (MVP)

### API Endpoints

- `POST /jobs`: Creates a new job and adds it to the queue.
- `GET /jobs/{id}`: Retrieves the status and result of a specific job.

### Job Statuses

- `pending`
- `running`
- `done`
- `failed`

## How to Run

1.  **Build the application (with CGO disabled for SQLite):**
    ```bash
    CGO_ENABLED=0 go build -o app ./cmd/app/main.go
    ```

2.  **Run the application:**
    ```bash
    ./app
    ```
    The server will start on `http://localhost:8080`.

### Example API Usage

**1. Create a new job:**

```bash
curl -X POST -H "Content-Type: application/json" -d '{"type":"send_email","payload":{"to":"test@example.com","subject":"Hello from Job Queue"}}' http://localhost:8080/jobs
```

Example response:
```json
{
  "id": "some-uuid",
  "type": "send_email",
  "payload": "eyAidG8iOiAidGVzdHVAZXhhbXBsZS5jb20iLCAic3ViamVjdCI6ICJIZWxsbyBmcm9tIEpvYiBRdWV1ZSIgfQ==",
  "status": "pending",
  "created_at": "2023-10-27T10:00:00Z",
  "updated_at": "2023-10-27T10:00:00Z"
}
```

**2. Get job status and result:**

Replace `some-uuid` with the actual job ID from the `POST` response.

```bash
curl http://localhost:8080/jobs/some-uuid
```

Example response (while pending):
```json
{
  "id": "some-uuid",
  "type": "send_email",
  "payload": "eyAidG8iOiAidGVzdHVAZXhhbXBsZS5jb20iLCAic3ViamVjdCI6ICJIZWxsbyBmcm9tIEpvYiBRdWV1ZSIgfQ==",
  "status": "pending",
  "created_at": "2023-10-27T10:00:00Z",
  "updated_at": "2023-10-27T10:00:00Z"
}
```

Example response (after processing):
```json
{
  "id": "some-uuid",
  "type": "send_email",
  "payload": "eyAidG8iOiAidGVzdHVAZXhhbXBsZS5jb20iLCAic3ViamVjdCI6ICJIZWxsbyBmcm9tIEpvYiBRdWV1ZSIgfQ==",
  "status": "done",
  "result": "UHJvY2Vzc2VkOiBzZW5kX2VtYWls",
  "created_at": "2023-10-27T10:00:00Z",
  "updated_at": "2023-10-27T10:00:02Z"
}
```