# marketplace-simulator

Учебный проект «Симулятор маркетплейса» — два микросервиса на Go с генератором нагрузки и полным observability-стеком.

## О проекте

Система состоит из двух микросервисов (**cart** и **product**), генератора нагрузки (**loadgen**) и инфраструктуры: PostgreSQL, Redis, Kafka, etcd, Prometheus, Grafana, Tempo, Loki.

### Основные флоу

**Добавление товара в корзину**

Клиент отправляет `POST /user/{id}/cart/{sku}` на cart. Cart обращается к product по gRPC (`GetProduct`) — получает цену и название товара. Product отдаёт данные из Redis (Cache-Aside); при промахе идёт в PostgreSQL. Cart сохраняет позицию в своей БД и возвращает ответ.

**Оформление заказа (checkout)**

Клиент отправляет `POST /user/{id}/cart/checkout` на cart. Cart получает позиции корзины из БД, вызывает `ReserveProduct` на product (product создаёт записи резервирований — остатки пока не меняются). Затем в одной транзакции cart очищает корзину и записывает задачи подтверждения в outbox-таблицу. Фоновая job асинхронно вызывает `ConfirmReservation` на product — product списывает товары со склада, удаляет резервирования и публикует события в Kafka. Loadgen-replenisher читает эти события и пополняет склад, когда остаток падает ниже порога.

## Технологии

| Область | Инструменты |
|---------|------------|
| **Транспорт** | HTTP (`net/http`) — cart; gRPC + grpc-gateway (REST-обёртка) — product; Protobuf — сериализация gRPC |
| **База данных** | PostgreSQL (pgx/v5, pgxpool); миграции — goose |
| **Кеширование** | Redis 7 — Cache-Aside в product; декоратор `CachedProductRepository`; graceful degradation при недоступности |
| **Брокер сообщений** | Apache Kafka — product публикует события `product.events`; loadgen читает и пополняет склад |
| **Трейсинг** | OpenTelemetry SDK → OTLP gRPC → Grafana Tempo; инструментированы HTTP, gRPC, PostgreSQL (pgx tracer), Redis |
| **Логирование** | `log/slog` (stdlib); Promtail собирает логи Docker-контейнеров → Loki; просмотр в Grafana с drill-down из трейсов |
| **Метрики** | Prometheus; дашборды в Grafana: Cart, Products, Business Metrics, Outbox Overview, Postgres Overview |
| **Динамическая конфигурация** | etcd — hot-reload без рестарта сервисов; первый старт сидирует конфиг из YAML |

---

## Быстрый старт

```bash
git clone https://github.com/jva44ka/marketplace-simulator.git
cd marketplace-simulator
docker-compose up
```

После запуска:

|                   |                                                                  |
|-------------------|------------------------------------------------------------------|
| Grafana           | [http://localhost:3000](http://localhost:3000) (admin / admin)   |
| Prometheus        | [http://localhost:9090](http://localhost:9090)                   |
| Kafka UI          | [http://localhost:8090](http://localhost:8090)                   |
| ETCD UI           | [http://localhost:8091](http://localhost:8091)                   |
| Swagger — product | [http://localhost:5001/swagger/](http://localhost:5001/swagger/) |
| Swagger — cart    | [http://localhost:5002/swagger/](http://localhost:5002/swagger/) |

## Архитектурные паттерны

**Transactional Outbox** — cart атомарно очищает корзину и записывает задачи подтверждения резервирований в одной транзакции. Product так же атомарно обновляет остатки и пишет события в Kafka. Никаких потерь сообщений при падении сервиса.

**Circuit Breaker** — gRPC-клиент cart к product защищён автоматом: при превышении порога ошибок цепь размыкается и запросы отклоняются сразу, не тратя время на таймаут. Настраивается: порог, окно, минимальный трафик для срабатывания.

**Retry + Exponential Backoff + Jitter** — временные сбои gRPC-вызовов к product переповторяются с экспоненциальной паузой и случайным отклонением, чтобы не создавать thundering herd при восстановлении.

**Dead Letter Queue** — outbox-записи которые не удалось доставить после N попыток уходят в dead letter и не блокируют очередь. Счётчик мониторится в Grafana с алертом.

**Optimistic Locking** — конкурентное обновление остатков товаров в product защищено оптимистичной блокировкой по версии записи. Конфликты версий отслеживаются отдельной метрикой.

**Cache-Aside (Redis)** — чтения товаров через `GetProduct` обслуживаются из Redis; промах → БД с async-прогревом. После обновления остатков product пишет в `outbox.cache_updates` и фоновая джоба инвалидирует кеш. Stale-записи отлавливаются сравнением PostgreSQL xmin. При недоступности Redis сервис продолжает работу без кеша (graceful degradation).

**Read-Your-Writes (RYW)** — product публикует сообщения в топик `product.events` с полем `transaction_id`. При последующих запросах можно передать это поле в ручку `get: "/v1/products/{sku}"` чтобы получить актуальную версию продукта. И если например в кеше будет неактуальная запись - сервис пойдет в базу. Это гарантирует, что клиент, использующий поле `transactionId`, никогда не прочитает состояние товара старее того, о котором уже получил событие.

**Rate Limiter** — входящий трафик на product-сервис ограничен token bucket'ом, защищая БД от перегрузки при всплесках нагрузки.

**Drain Mode** — outbox jobs переключаются между двумя интервалами: минимальная пауза в idle и немедленный следующий тик пока в очереди есть записи, не нагружая базу холостыми запросами.

---

→ [Подробная документация](docs/README.md)
