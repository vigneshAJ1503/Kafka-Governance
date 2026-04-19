# Kafka Governance Control Plane Architecture

This document describes the high-level architecture of the Kafka Governance proxy.

## High Level Diagram

```mermaid
flowchart TD
    Client[Client Browser / CLI]

    subgraph "Kafka Governance Service"
        Router[Chi/Gin Router & Middleware]
        API[API Handlers]
        Service[Service Layer (Business Logic)]
        Auth[JWT Middleware Validation]
        
        Router --> Auth
        Auth --> API
        API --> Service
    end

    subgraph "Data & Policy Engine"
        MongoDB[(MongoDB)]
        CedarAgent[Cedar CLI Agent (Docker)]
    end

    subgraph "Kafka Cluster (Data Plane)"
        KafkaBrokers[Kafka Brokers]
    end

    Client -- HTTP Requests --> Router
    Service -- Topic / Policy State --> MongoDB
    Service -- Request Authorization --> CedarAgent

    CedarAgent -. Decision .-> Service
    
    %% Note to indicate scope
    Service -. "Manual Admin setup needed for actual topic" .-> KafkaBrokers
```

### Components overview

1. **API Router**: Exposes `/topics`, `/policies`, and `/auth` routes. Processes incoming HTTP requests.
2. **Auth Middleware**: Parses and validates the JWT, ensuring specific routes are protected. It injects the `userId` and user `role` onto the request context for further authorization assertions.
3. **Service Layer**: Maps business workflows. When a user requests to `Create Topic`, it triggers the application to query the defined policies and push them into the `CedarAgent`.
4. **Data Plane**: Uses MongoDB instance for persisting `Topic` structs and `Policy` structs.
5. **Policy Engine (AWS Cedar)**: Running as an independent process (Docker Container). Analyzes `Principal` (User), `Action` (e.g. `CreateTopic`), `Resource` (e.g., `orders.topic`), and evaluates them against the policies we define.
