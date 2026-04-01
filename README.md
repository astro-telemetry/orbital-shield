# Orbital Shield
**Real-time Network Integrity & Telemetry for Orbital Infrastructure**

## Vision
Orbital Shield provides high-fidelity observability for satellite network constellations. By deploying lightweight, Go-native sensors at the edge, the system maintains a real-time integrity map, feeding structured telemetry into an ML-driven inference engine to detect anomalies before they escalate into kinetic failures.

## The Problem: "The Rogue Network Paradigm"
Satellite networks offer global connectivity but introduce a massive, decentralized attack surface. Conventional monitoring (like FIRE or ARTEMIS) focuses on aggregate, long-term behavior. However, bad actors utilize:
* **Ephemeral Infrastructure:** Rapidly deployed cloud nodes that are "burned" after an attack.
* **BGP Hijacking:** Exploiting the "Prefer Customer" routing logic to divert sensitive orbital data.
* **Interface Ghosting:** Unauthorized virtual tunnels used for Command & Control (C2) exfiltration.

## The Solution
Orbital Shield flips the script from *reactive* to *proactive*.
1. **Low-Level Sensing:** A Go-native agent performs bitwise hardware status checks to identify unauthorized interface changes.
2. **Structured Telemetry:** Real-time data is marshaled into JSON feature vectors, optimized for machine learning ingestion.
3. **ML Pipeline:** Anomaly detection identifies "Rogue" patterns in interface behavior and traffic routing.

## Architecture
The system is designed for modularity, high-performance execution, and immediate deployment in distributed Linux environments.
* **Sensor Core:** System-level Go utilizing the `net` package for kernel-level interface interrogation.
* **Containerized Daemon:** Packaged via multi-stage Docker builds into a minimal Alpine Linux image for a zero-dependency, secure footprint.
* **Structured Logging:** Utilizes Go 1.21+ `log/slog` to output strictly typed JSON operational logs to `stdout` for SIEM ingestion.
* **Data Protocol:** Standardized JSON schemas for the telemetry payloads, piped securely to the host machine via volume mounts for ML model integration.
  
## Security Goals
* **Interface Integrity:** Detect and alert on unauthorized status changes or new virtual bridges.
* **C2 Suppression:** Identify signatures of Command and Control traffic hidden in standard telemetry.
* **Route Validation:** Monitor local gateway behavior to flag potential BGP prefix hijacks.

## Getting Started

### Prerequisites
* Docker Engine
* Go 1.21+ (For local development only)

### Deployment (Production Container)
The agent is designed to run as a detached, containerized daemon that mounts a secure volume to the host to deposit telemetry payloads.

**1. Build the image:**
```bash
docker build -t orbital-shield .
```

**2. Run the agent and mount the payload directory:**
```bash
docker run -d \
  --name orbshield_agent \
  -v "$(pwd)/telemetry_logs:/root/telemetry_logs" \
  orbital-shield
```
**3. Monitor the structured operational logs:**
```bash
docker logs -f orbshield_agent
```

### Telemetry Payload Example
```json
{
  "hostname": "satellite-gateway-01",
  "timestamp": 1711982345,
  "interface_name": "eth0",
  "is_active": true,
  "mtu": 1500,
  "ip_list": ["192.168.1.50/24"]
}
```
