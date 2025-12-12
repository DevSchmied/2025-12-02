# URL Status & PDF Report Service — 2025-12-02

## Project Overview
The **web service** accepts URLs, checks their availability, stores URL sets,
and generates PDF reports upon request.

---

## Technologies and Stack
- **Go 1.23+**
- **Gin** — web framework
- **gopdf** — PDF generation
- **JSON storage** for data
- **sync.WaitGroup**, **channels** — concurrency
- **Graceful Shutdown** support

---

## Core Functionality
### 1. URL Availability Check

Request example:

```json
{"links": ["google.com", "malformed.gg"]}
```


The service:
- checks each URL using a worker pool,
- determines the status (available / not available),
- stores the set under a unique number (links_num),
- returns the result:

```json
{
  "links": {
    "google.com": "available",
    "malformed.gg": "not available"
  },
  "links_num": 1
}
```

### 2. PDF Generation by Set Numbers

Request example:

```json
{"links_list": [1, 2]}
```

The service:
- retrieves stored data from JSON storage,
- generates a PDF report with statuses,
- sends the PDF as a binary file.
A TTF font **DejaVu Sans** is used.

### 3. Graceful Shutdown

Upon receiving SIGINT/SIGTERM, the service:
- continues accepting tasks for another 3 seconds,
- closes the task queue,
- waits for workers to finish (WaitGroup),
- saves the storage to disk,
- shuts down gracefully without data loss.

---

## Design Patterns Used
- **WorkerPool** — parallel URL checking
- **Handler Factory** — CheckURLs(...) gin.HandlerFunc, dynamic handler generation
- **Dependency Injection** — for testing the storage module
- **Delayed Graceful Shutdown** — correct service termination

---

## Unit Testing

Unit testing is demonstrated on the **storage** module.

The tests follow:
- **AAA structure** (Arrange–Act–Assert)
- **FIRST principles**
- usage of **mock objects**
- **table-driven tests**
- **Dependency Injection**

**Covered methods:**
- LoadFromDisk()
- SaveToDisk()

---

## Architecture and Project Structure

The project structure is detailed in the **Package Diagram**, located at:

```
docs/diagrams/package_diagram.png
docs/diagrams/package_diagram.dia
```

The diagram reflects the key project components and their responsibilities.

### Main Packages
- **cmd/**  
  Application entry point; server startup and graceful shutdown handling.
- **internal/check/**  
  URL availability checking function.
- **internal/handlers/**  
  HTTP handlers for URL checking and PDF generation.
- **internal/pdf/**  
  PDF report generation.
- **internal/server/**  
  HTTP server, routing, initialization.
- **internal/service/**  
  Worker Pool for parallel task processing.
- **internal/storage/**  
  In-memory storage, JSON file, Dependency Injection.

---

## UML Diagrams

As part of the assignment, the following diagrams were created:
- **Use Case Diagram**
- **Activity Diagram**
- **Package Diagram**

All materials are located in the folder:

**docs/diagrams/**

Formats:
- source files: .dia (created with **Dia**)
- diagram screenshots

Each diagram is provided in **two versions:**
1. **.dia** — editable format
2. **.png** — screenshot for quick review

---

## Running the Project

```bash
go run ./cmd/main.go
```

The service is available at:

http://localhost:8080

---

## HTTP Request Examples
1. URL availability check

**POST /check-links**

```json
{"links": ["google.com", "ya.ru"]}
```

2. PDF generation

**POST /make-pdf**

```json
{"links_list": [1, 2]}
```

---

## Manual Testing (Postman)

End-to-end testing was performed:
1. URL status checking
2. PDF generation from saved sets

Screenshots are located in:

docs/screenshots/


Files:
- check_links — successful check
- make_pdf_01 — PDF for a single set
- make_pdf_02 — PDF for multiple sets

All tests were performed on localhost:8080.