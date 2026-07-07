# Kafka Governance Control Plane Architecture

This document describes the high-level architecture of the Kafka Governance proxy including the newly added caching infrastructure components.

## High-Level Design (HLD)

```mermaid
flowchart TD
    Client[Client Browser / CLI]

    subgraph "Kafka Governance Platform"
        subgraph "Frontend Layer"
            UI[Next.js Application]
        end

        subgraph "Backend Services"
            Router[Gin Router & Middleware]
            API[API Handlers]
            Service[Business Logic Layer]
            Auth[JWT Middleware Validation]
            
            Router --> Auth
            Auth --> API
            API --> Service
        end
    end

    subgraph "Infrastructure Data & Policy Engine"
        MongoDB[(MongoDB)]
        Redis[(Redis Cache)]
        CedarAgent[Cedar CLI Agent Container]
    end

    subgraph "Kafka Cluster (Data Plane)"
        KafkaBrokers[Kafka Brokers]
        KafkaUI[Kafka UI Web Manager]
    end

    Client -- HTTP UI/Proxy Requests --> UI
    UI -- '/api/v1/*' Proxied --> Router
    Service -- Topic / Policy State --> MongoDB
    
    %% Caching Layer Mapping
    Service -- 1. Query Cached Authorization State --> Redis
    Service -- 2. Evaluate Policy Payload (Cache Miss) --> CedarAgent
    CedarAgent -. Decision .-> Service
    Service -- 3. Store Computed Policies Async --> Redis
    
    %% Boundaries
    Client -. Produce/Consume TCP Stream Data .-> KafkaBrokers
    KafkaUI -- Monitor Broker Info --> KafkaBrokers
```

### Components overview

1. **Frontend**: Next.js 14 client mapping React architectures and dynamically proxying authorization bounds natively securely.
2. **Auth Middleware**: Parses and validates the JWT decoupled locally. It injects the `userId` and user `role` onto the request context targeting API route maps.
3. **Service Layer**: Maps business workflows. When a user requests to `Create Topic`, it triggers the application to query the defined policies through the authentication topology.
4. **Data Plane**: Uses MongoDB instance for persisting structured elements like `Topic` & `Policy` logic natively decoupled.
5. **Caching Layer (Redis)**: Deployed to significantly accelerate latency-intensive validation hits targeted dynamically against AWS Cedar. By evaluating policies mapped across identity and resources rapidly and locally caching them securely, the payload execution minimizes external container networking blocks dramatically over heavily loaded infrastructures.
6. **Policy Engine (AWS Cedar)**: Running securely decoupled. Analyzes `Principal` (User), `Action` (e.g. `CreateTopic`), `Resource` (e.g., `orders.topic`), and evaluates them.

---

## Low-Level Design (LLD): Policy Evaluation Sequence

This illustrates the exact sequence during an authorization enforcement phase where Redis reduces latency intelligently.

```mermaid
sequenceDiagram
    actor Client
    participant API as API Handler 
    participant Service as Topic Service
    participant Redis as Redis Cache 
    participant Cedar as Cedar Agent

    Client->>API: POST /api/v1/topics (Create Request)
    API->>Service: CreateTopic(ctx, TopicStruct)
    
    Service->>Redis: GET cedar_policy:{user}:{action}:{topic} (1ms)
    
    alt Cache Hit
        Redis-->>Service: return "true"
        Service->>API: Create MongoDB Document (Approved)
    else Cache Miss
        Redis-->>Service: return nil (Key Not Found)
        Service->>+Cedar: Evaluate Policy (HTTP Execution) -> 15ms+
        Cedar-->>-Service: return ALLOW (true)
        
        par Write-Behind Cache
            Service->>Redis: SET EX cedar_policy:{user}... "true" 300s
        end
        
        Service->>API: Create MongoDB Document (Approved)
    end
    
    API-->>Client: 201 Created Status
```
