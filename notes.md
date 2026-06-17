---
  Struktur gold-gym-be-v2 sekarang terdiri dari dua layer:

  Layer 1 — Kode existing (sudah ada sebelumnya)
  - Single Go service (cmd/http/main.go) pakai Gin framework
  - Semua package sudah public → tidak ada yang perlu diganti di sini

  Layer 2 — Struktur courier-style yang baru dibuat

  gold-gym-be-v2/
  ├── docker-compose.yml          ← infra courier-style (MySQL, Redis, RabbitMQ, nginx)
  ├── app.env / .env.example      ← env config
  ├── logs.sh                     ← tail log container (./logs.sh users)
  ├── up.sh / halt.sh / ssh.sh    ← start/stop/shell ke container
  ├── scripts/functions.sh        ← helper functions (mirip courier)
  ├── servers/
  │   ├── gobase/                 ← Base Dockerfile untuk semua Go service
  │   └── reverse-proxy/          ← Nginx route /users/* dan /members/*
  ├── pkg/gymkit/                 ← ★ Pengganti sange (SEMUA PUBLIC PACKAGES)
  │   ├── gymkit.go               ← konstanta
  │   ├── utility.go              ← GetEnv, InArray, IsJSON, IsEmptyString
  │   ├── response.go             ← NewResponse (JSON standard response)
  │   ├── db.go                   ← DatabaseBase, GetMysqlConnection (gorm)
  │   ├── config.go               ← GetRabbitConfig (amqp URL builder)
  │   └── middleware/auth/        ← JWT middleware (pakai golang-jwt/jwt, bukan sange)
  └── workspace/
      ├── users/                  ← Service user management (CRUD + Login + JWT)
      │   ├── main.go
      │   ├── models/user.go
      │   ├── controllers/        ← user_controller.go, auth_controller.go
      │   ├── routers/base.go
      │   ├── helpers/apps.go     ← IsUserAppID, IsAdminAppID (mirip courier)
      │   └── middleware/
      └── members/                ← Service gym member (CRUD membership)
          ├── main.go
          ├── models/member.go
          ├── controllers/member_controller.go
          └── ...

  Private packages yang diganti:

  ┌─────────────────────────────────────────────────┬────────────────────────────────────────────────┐
  │                Courier (Private)                │            gold-gym-be-v2 (Public)             │
  ├─────────────────────────────────────────────────┼────────────────────────────────────────────────┤
  │ github.com/mataharibiz/sange/v2                 │ github.com/gold-gym/gymkit (local /pkg/gymkit) │
  ├─────────────────────────────────────────────────┼────────────────────────────────────────────────┤
  │ github.com/mataharibiz/sange/v2/middleware/auth │ github.com/gold-gym/gymkit/middleware/auth     │
  ├─────────────────────────────────────────────────┼────────────────────────────────────────────────┤
  │ github.com/mataharibiz/cony                     │ github.com/rabbitmq/amqp091-go (langsung)      │
  └─────────────────────────────────────────────────┴────────────────────────────────────────────────┘

  Cara pakai:
  cp .env.example app.env   # isi konfigurasi
  ./up.sh                   # start semua service
  ./logs.sh users           # lihat log service users
  ./logs.sh members         # lihat log service members
  ./ssh.sh users            # masuk ke shell container