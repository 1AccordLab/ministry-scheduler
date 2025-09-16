# 部署指南

## 技術需求

### 基礎環境

* **作業系統**: Linux (Ubuntu 20.04+ 或 CentOS 8+ 推薦)
* **記憶體**: 最少 2GB，建議 4GB+
* **儲存空間**: 最少 20GB，建議 50GB+
* **網路**: 外網連線 (Line Bot, Google Calendar 整合需要)

### 必需軟體

* **Podman** 3.0+
* **PostgreSQL** 14+
* **Nginx** (反向代理)
* **域名與 SSL 憑證**

## 部署架構

### 容器化部署 (推薦)

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    Nginx    │────│   Go App    │────│  PostgreSQL │
│ (SSL/Proxy) │    │ (Container) │    │ (Container) │
└─────────────┘    └─────────────┘    └─────────────┘
```

### 檔案結構

```
/opt/ministry-scheduler/
├── docker-compose.yml
├── nginx/
│   └── nginx.conf
├── app/
│   ├── config.yaml
│   └── migrations/
├── data/
│   └── postgres/
└── logs/
    ├── nginx/
    ├── app/
    └── postgres/
```

## 部署步驟

### 1. 環境準備

**安裝 Podman**：

```bash
# Ubuntu/Debian
sudo apt update && sudo apt install -y podman

# CentOS/RHEL
sudo dnf install -y podman
```

**安裝 PostgreSQL**：

```bash
# 使用容器部署
podman run -d \
    --name ministry-db \
    -e POSTGRES_DB=ministry \
    -e POSTGRES_USER=ministry_user \
    -e POSTGRES_PASSWORD=your_secure_password \
    -v /opt/ministry-scheduler/data/postgres:/var/lib/postgresql/data \
    -p 5432:5432 \
    postgres:14
```

### 2. 應用程式部署

**建立設定檔 `/opt/ministry-scheduler/app/config.yaml`**：

```yaml
server:
  port: 8080
  host: "0.0.0.0"

database:
  host: "ministry-db"
  port: 5432
  user: "ministry_user"
  password: "your_secure_password"
  dbname: "ministry"
  sslmode: "require"

line_bot:
  channel_secret: "your_line_channel_secret"
  channel_token: "your_line_channel_token"
  webhook_url: "https://yourdomain.com/webhook/line"

google_calendar:
  client_id: "your_google_client_id"
  client_secret: "your_google_client_secret"
  redirect_url: "https://yourdomain.com/auth/google/callback"

security:
  jwt_secret: "your_jwt_secret_key"
  session_key: "your_session_secret"
```

**建立 Docker Compose 檔案**：

```yaml
version: '3.8'

services:
  ministry-db:
    image: postgres:14
    container_name: ministry-db
    environment:
      POSTGRES_DB: ministry
      POSTGRES_USER: ministry_user
      POSTGRES_PASSWORD: your_secure_password
    volumes:
      - /opt/ministry-scheduler/data/postgres:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    restart: unless-stopped

  ministry-app:
    image: ministry-scheduler:latest
    container_name: ministry-app
    depends_on:
      - ministry-db
    volumes:
      - /opt/ministry-scheduler/app/config.yaml:/app/config.yaml
      - /opt/ministry-scheduler/logs/app:/app/logs
    ports:
      - "8080:8080"
    restart: unless-stopped
    environment:
      - CONFIG_FILE=/app/config.yaml
```

### 3. 應用程式建置

**建立 Dockerfile**：

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ministry-scheduler ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/ministry-scheduler .
COPY --from=builder /app/web ./web
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./ministry-scheduler"]
```

**建置映像檔**：

```bash
cd /path/to/source
podman build -t ministry-scheduler:latest .
```

### 4. 資料庫初始化

**執行遷移**：

```bash
podman exec ministry-app ./ministry-scheduler migrate up
```

**建立初始管理員**：

```bash
podman exec ministry-app ./ministry-scheduler create-admin \
    --email admin@church.org \
    --password admin123 \
    --name "系統管理員"
```

### 5. Nginx 設定

**建立 `/etc/nginx/sites-available/ministry-scheduler`**：

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /path/to/ssl/cert.pem;
    ssl_certificate_key /path/to/ssl/private.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;

    client_max_body_size 10M;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket support
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    location /static/ {
        alias /opt/ministry-scheduler/app/web/static/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    access_log /var/log/nginx/ministry-scheduler-access.log;
    error_log /var/log/nginx/ministry-scheduler-error.log;
}
```

**啟用站台**：

```bash
sudo ln -s /etc/nginx/sites-available/ministry-scheduler /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 外部服務設定

### Line Bot 設定

1. **申請 Line Developer 帳號**
2. **建立 Messaging API Channel**
3. **設定 Webhook URL**: `https://yourdomain.com/webhook/line`
4. **取得 Channel Secret 和 Channel Token**

### Google Calendar API 設定

1. **建立 Google Cloud Project**
2. **啟用 Calendar API**
3. **建立 OAuth 2.0 憑證**
4. **設定授權重新導向 URI**: `https://yourdomain.com/auth/google/callback`

## 監控與維護

### 系統監控

**建立監控腳本 `/opt/ministry-scheduler/monitor.sh`**：

```bash
#!/bin/bash

# 檢查容器狀態
if ! podman ps | grep -q ministry-app; then
    echo "應用程式容器未執行，嘗試重啟..."
    podman restart ministry-app
fi

if ! podman ps | grep -q ministry-db; then
    echo "資料庫容器未執行，嘗試重啟..."
    podman restart ministry-db
fi

# 檢查磁碟空間
DISK_USAGE=$(df /opt/ministry-scheduler | tail -1 | awk '{print $5}' | sed 's/%//')
if [ $DISK_USAGE -gt 80 ]; then
    echo "磁碟空間不足：使用了 ${DISK_USAGE}%"
fi

# 檢查記憶體使用率
MEMORY_USAGE=$(free | grep Mem | awk '{printf("%.0f", $3/$2*100)}')
if [ $MEMORY_USAGE -gt 85 ]; then
    echo "記憶體使用率過高：${MEMORY_USAGE}%"
fi
```

**設定 Cron 自動監控**：

```bash
# 每 5 分鐘檢查一次
*/5 * * * * /opt/ministry-scheduler/monitor.sh >> /var/log/ministry-monitor.log 2>&1
```

### 日誌管理

**設定 Logrotate `/etc/logrotate.d/ministry-scheduler`**：

```
/opt/ministry-scheduler/logs/app/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 root root
    postrotate
        podman exec ministry-app killall -USR1 ministry-scheduler
    endscript
}
```

### 資料備份

**建立備份腳本 `/opt/ministry-scheduler/backup.sh`**：

```bash
#!/bin/bash

BACKUP_DIR="/opt/ministry-scheduler/backups"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# 資料庫備份
podman exec ministry-db pg_dump -U ministry_user ministry > \
    $BACKUP_DIR/ministry_db_$DATE.sql

# 設定與日誌備份
tar -czf $BACKUP_DIR/config_$DATE.tar.gz \
    /opt/ministry-scheduler/app/config.yaml \
    /opt/ministry-scheduler/logs

# 清理 30 天前的備份
find $BACKUP_DIR -name "*.sql" -mtime +30 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +30 -delete
```

**設定每日備份**：

```bash
# 每天晚上 2 點執行備份
0 2 * * * /opt/ministry-scheduler/backup.sh
```

### 更新流程

**應用程式更新**：

```bash
# 1. 備份資料
/opt/ministry-scheduler/backup.sh

# 2. 停止舊容器
podman stop ministry-app

# 3. 拉取新映像檔
podman pull ministry-scheduler:latest

# 4. 執行資料庫遷移 (如有需要)
podman run --rm -v /opt/ministry-scheduler/app/config.yaml:/app/config.yaml \
    ministry-scheduler:latest ./ministry-scheduler migrate up

# 5. 啟動新容器
podman start ministry-app
```

## 效能調校

### PostgreSQL 最佳化

**編輯 PostgreSQL 設定**：

```bash
podman exec -it ministry-db psql -U ministry_user -d ministry
```

```sql
-- 調整記憶體設定
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';

-- 調整連線設定
ALTER SYSTEM SET max_connections = '200';

-- 重新載入設定
SELECT pg_reload_conf();
```

### 應用程式最佳化

**調整 Go 應用程式設定**：

```yaml
# config.yaml
performance:
  max_connections: 100
  idle_timeout: 30s
  read_timeout: 10s
  write_timeout: 10s
  max_request_size: 10485760  # 10MB
```

## 故障排除

### 常見問題

**1. 應用程式無法連接資料庫**

```bash
# 檢查資料庫狀態
podman logs ministry-db

# 檢查網路連線
podman exec ministry-app ping ministry-db
```

**2. Line Bot 無法接收訊息**

```bash
# 檢查 Webhook 設定
curl -X POST https://yourdomain.com/webhook/line

# 檢查應用程式日誌
podman logs ministry-app | grep -i line
```

**3. SSL 憑證問題**

```bash
# 檢查憑證有效性
openssl x509 -in /path/to/cert.pem -text -noout

# 更新憑證 (Let's Encrypt)
certbot renew
systemctl reload nginx
```

### 啟動腳本

**建立系統服務 `/etc/systemd/system/ministry-scheduler.service`**：

```ini
[Unit]
Description=Ministry Scheduler Application
Requires=multi-user.target
After=multi-user.target

[Service]
Type=forking
ExecStart=/opt/ministry-scheduler/start.sh
ExecStop=/opt/ministry-scheduler/stop.sh
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**啟動腳本 `/opt/ministry-scheduler/start.sh`**：

```bash
#!/bin/bash
cd /opt/ministry-scheduler
podman-compose up -d
```

**停止腳本 `/opt/ministry-scheduler/stop.sh`**：

```bash
#!/bin/bash
cd /opt/ministry-scheduler
podman-compose down
```
