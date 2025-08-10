# VIDEO - https://drive.google.com/file/d/15RVY8UpEPOUJb8T6gIf0QGBVUsDSs8Yq/view?usp=sharing

# web-l0

Golang service with web-interface and API

# Architecture

```mermaid
graph LR
    %% Стили узлов
    classDef app fill:#4CAF50,stroke:#388E3C,color:white,stroke-width:2px
    classDef server fill:#2196F3,stroke:#0b7dda,color:white
    classDef storage fill:#FF9800,stroke:#F57C00,color:white
    classDef queue fill:#607D8B,stroke:#455A64,color:white
    classDef db fill:#f44336,stroke:#d32f2f,color:white
    classDef api fill:#9C27B0,stroke:#7B1FA2,color:white
    classDef internal fill:#0ABAB5,stroke:#7B1FA2,color:white
    
    %% Узлы
    A[Client]:::app
    B[Web App]:::internal
    C[Storage]:::internal
    L[Service]:::internal
    Q[API]:::api
    W[Pages]:::api
    G[Kafka]:::queue
    E[(Redis)]:::db
    F[(PostgreSQL)]:::db
    
 
    
    subgraph "Data Layer"
        E
        F
    end

    subgraph "HTTP Layer"
        Q
        W
    end

    subgraph "Service Layer"
        G
    end
    
    %% Четкие связи с пояснениями
    A -->|Get new orders| L
    L -->|Read new orders| G
    A -->|Listen And Serve| B
    B -->|API Calls| Q
    B -->|Browser Calls| W
    C -->|Cache Access| E
    C -->|DB Persistence| F
    A -->|Get saved orders| C
```

# Data

### Api Request Order

```mermaid
sequenceDiagram

 autonumber
    title Поиск заказа по order_uid

    participant Client
    participant API as "API Layer"
    participant Storage as "Storage Service"
    participant Cache as "Cache (Redis)"
    participant Database as "Database (PostgreSQL)"

    Note over Client: Запрос информации о заказе

    Client->>API: GET /order/{order_uid}
    API->>Storage: FindOrder(order_uid)

    rect rgb(240,240,255)
        Note right of Storage: Поиск в кеше
        Storage->>Cache: Find(order_uid)
        
        alt Найдено в кеше
            Cache-->>Storage: order.Order object
            Storage-->>API: order.Order object
            API-->>Client: 200 OK (with data, json)
        else Не найдено в кеше
            rect rgb(255,240,240)
                Note right of Storage: Поиск в БД
                Storage->>DB: Find(order_uid)
                
                alt Найдено в БД
                    DB-->>Storage: Order data
                    Storage->>Cache: Save(order_uid, data)
                    Storage-->>API: order.Order object
                    API-->>Client: 200 OK (with data json)
                else Не найдено
                    DB-->>Storage: Error (Not Found)
                    Storage-->>API: Error
                    API-->>Client: 404 Not Found
                end
            end
        end
    end

```

# API Documentation

### localhost:8080/swagger/index.html


# Local startup

- docker-compose up

  ```
  postgres: localhost:5432 (user=dev, password=qqq, dbname=mydb) | load dump with test data
  pgadmin: localhost:8081 (admin@example.com:admin)
  kafka: localhost:9092 (topic=orders)
  redis: localhost:6379
  ```

- config loading
  You should use -c "..." with config file path, config file is .yml

```
type Config struct {

WebConfig  `yaml:"web_config" env-required:"true"`

PostgresConfig  `yaml:"postgres_config"`

RedisConfig  `yaml:"redis"`

KafkaOrdersConfig  `yaml:"kafka"`



InitialDataSize  int  `yaml:"initial_data_size" env-default:"100"`

}



type WebConfig struct {

Host  string  `yaml:"host" env-required:"true"`

Port  string  `yaml:"port" env-required:"true"`

ReadTimeout  time.Duration  `yaml:"read_timeout" env-default:"10s"`

WriteTimeout  time.Duration  `yaml:"write_timeout" env-default:"10s"`

}



type PostgresConfig struct {

Host  string  `yaml:"host" env-required:"true"`

Port  string  `yaml:"port" env-required:"true"`

User  string  `yaml:"user" env-required:"true"`

Password  string  `yaml:"password" env-required:"true"`

DBName  string  `yaml:"db_name" env-required:"true"`

SSLMode  bool  `yaml:"sslmode" env-default:"false"`

}



type RedisConfig struct {

Host  string  `yaml:"host" env-required:"true"`

Port  string  `yaml:"port" env-required:"true"`

Password  string  `yaml:"password"`

DBName  int  `yaml:"db_name"`

}



type KafkaOrdersConfig struct {

Brokers  []string  `yaml:"brokers" env-required:"true"`

Topic  string  `yaml:"topic" env-required:"true"`

MinBytes  int  `yaml:"min_bytes" env-default:"1"`

MaxBytes  int  `yaml:"max_bytes" env-default:"10e6"`

}
```

# logs
  Logs saved in ./log/app.log
