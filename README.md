# Real-Time Log Monitoring System

This project is a Real-Time Log Monitoring System built with Golang and Kafka. It consists of the following services:

- **Producer:** Simulates generating logs and sends them to Kafka.
- **Consumer:** Consumes logs from Kafka, processes them, and optionally sends alerts.
- **Dashboard:** Displays logs in real-time through a web interface.

## Prerequisites

- Docker and Docker Compose
- Golang (for local development)

## Setup

1. Clone this repository:

   ```bash
   git clone <repository_url>
   cd real-time-log-monitoring
