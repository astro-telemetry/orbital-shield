# Orbital Shield  
Lightweight network monitoring agent for real-time interface and routing visibility

---

## Overview

Orbital Shield is a containerized Go-based agent that collects network interface and routing state from Linux systems and outputs structured logs for downstream monitoring and detection systems.

The goal is to provide **low-level visibility into network changes** that are often missed by higher-level monitoring tools, enabling faster identification of anomalous behavior.

---

## What It Does

- Monitors network interfaces and captures:
  - interface status (up/down)
  - IP address assignments
  - MTU and configuration changes
- Outputs structured JSON logs for ingestion into SIEM or other monitoring pipelines
- Runs as a lightweight containerized daemon with minimal dependencies

---

## Current Capabilities

- System-level network inspection using Go (`net` package)
- Structured logging using `log/slog` (JSON output)
- Containerized deployment via multi-stage Docker builds
- Secure log output via mounted volumes for downstream processing

---

## Example Output

```json
{
  "hostname": "host-01",
  "timestamp": 1711982345,
  "interface_name": "eth0",
  "is_active": true,
  "mtu": 1500,
  "ip_list": ["192.168.1.50/24"]
}
```

---

## Use Cases

- Detecting unexpected interface changes or configuration drift
- Monitoring for anomalous routing or network behavior
- Providing structured input for downstream detection or alerting systems
- Supporting investigation workflows with consistent network state data

---

## Architecture

- **Agent:** Go-based process collecting network state
- **Container:** Minimal Docker image for portable deployment
- **Logging:** JSON output to stdout for integration with logging pipelines
- **Data Flow:** Agent → structured logs → external processing / detection systems

---

## Future Work

- Integration with anomaly detection models
- Policy enforcement (e.g., OPA) for validating network state
- Deployment via infrastructure-as-code (Terraform/Kubernetes)
- Expanded signal collection (routing tables, connection metadata)

---

## Getting Started

### Build

```bash
docker build -t orbital-shield .
```

### Run

```bash
docker run -d \
  --name orbshield_agent \
  -v "$(pwd)/telemetry_logs:/root/telemetry_logs" \
  orbital-shield
```

### View Logs

```bash
docker logs -f orbshield_agent
```
