# SSO Service

SSO Микросервис команда Go19

---

## 📁 Структура проекта
``` text
.
├── cmd/
│   └── app/
│       └── main.go             
├── config/
│   └── config.yaml             
├── google/
│   └── api/
│       ├── annotations.proto    
│       ├── http.proto
│       └── httpbody.proto
├── internal/
│   ├── auth/
│   │   └── auth.go             
│   ├── security/
│   │   └── security.go         
│   ├── config/
│   │   └── config.go            
│   ├── interceptor/
│   │   └── interceptor.go       
│   ├── logger/
│   │   └── logger.go            
│   ├── models/
│   │   └── models.go         
│   ├── repository/
│   │   └── repository.go       
│   ├── server/
│   │   ├── handlers.go          
│   │   └── server.go            
│   └── service/
│       └── service.go           
├── migrations/
│   ├── 001init.up.sql          
│   └── 001init.down.sql         
├── pkg/
│   └── api/
│       ├── sso_service.proto   
│       ├── sso_service.swagger.json 
│       └── test/
│           ├── sso_service.pb.go    
│           ├── sso_service_grpc.pb.go
│           └── sso_service.pb.gw.go  
├── .gitignore
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## 🚀 Быстрый старт

### Запуск всех сервисов

```bash
docker-compose up -d
```

### Остановка

```bash
docker-compose down
```

### Остановка с удалением данных

```bash
docker-compose run migrate-down
docker-compose down
```

---

## 🔌 Порты

| Сервис | Порт | Назначение |
|--------|------|------------|
| SSO (gRPC) | 50051 | gRPC API |
| SSO (HTTP) | 8080 | REST API |
| PostgreSQL | 5432 | База данных |
| Redis | 6379 | Кэш/сессии |

---

## ⚙️ Конфигурация

Конфиги монтируются из `./config` в `/app/config`.

### Переменные окружения БД

| Переменная | Значение |
|------------|----------|
| `POSTGRES_USER` | postgres |
| `POSTGRES_PASSWORD` | postgres |
| `POSTGRES_DB` | sso |

---

## 📡 Подключение

### PostgreSQL

```
postgres://postgres:postgres@localhost:5432/sso?sslmode=disable
```

### Redis

```
localhost:6379
```

### gRPC

```
localhost:50051
```

### HTTP

```
http://localhost:8080
```

---

## 🔧 Полезные команды

```bash
# Логи SSO
docker logs -f sso_service

# Логи БД
docker logs -f sso_db

# Перезапуск SSO
docker-compose restart sso

# Пересборка SSO
docker-compose up -d --build sso

# Зайти в контейнер БД
docker exec -it sso_db psql -U postgres -d sso
```

---

## 📋 Требования

- Docker 20.10+
- Docker Compose 2.0+
