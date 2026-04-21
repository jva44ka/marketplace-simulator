# marketplace-simulator

Оркестрирующий репозиторий учебного проекта «Симулятор маркетплейса».

Запускает всю инфраструктуру через `docker-compose`: микросервисы, базы данных, генератор нагрузки и observability-стек.

## Состав системы

| Сервис                  | Репозиторий                                                                                     | Описание                                                            |
|-------------------------|-------------------------------------------------------------------------------------------------|---------------------------------------------------------------------|
| **product**             | [marketplace-simulator-product](https://github.com/jva44ka/marketplace-simulator-product)       | Управление товарами (gRPC + REST, PostgreSQL, Kafka outbox)         |
| **cart**                | [marketplace-simulator-cart](https://github.com/jva44ka/marketplace-simulator-cart)             | Корзина покупок (REST, PostgreSQL, Kafka outbox)                    |
| **loadgen**             | [marketplace-simulator-loadgen](https://github.com/jva44ka/marketplace-simulator-loadgen)       | Генератор нагрузки (replenisher, order flow, cart viewer)           |
| **product-db**          | postgres:17.7                                                                                   | БД сервиса товаров                                                  |
| **cart-db**             | postgres:17.7                                                                                   | БД сервиса корзины                                                  |
| **product-migrations**  | migrator из [product](https://github.com/jva44ka/marketplace-simulator-product)                 | Накатывает миграции в product-db при старте                         |
| **cart-migrations**     | migrator из [cart](https://github.com/jva44ka/marketplace-simulator-cart)                       | Накатывает миграции в cart-db при старте                            |
| **kafka**               | confluentinc/cp-kafka:7.9.0                                                                     | Брокер сообщений (события об изменении товаров)                     |
| **kafka-ui**            | provectuslabs/kafka-ui                                                                          | Веб-интерфейс для Kafka                                             |
| **prometheus**          | prom/prometheus                                                                                 | Сбор метрик с сервисов                                              |
| **tempo**               | grafana/tempo:2.6.1                                                                             | Хранилище трейсов (OTLP)                                            |
| **loki**                | grafana/loki:3.4.2                                                                              | Хранилище логов                                                     |
| **promtail**            | grafana/promtail:3.4.2                                                                          | Агент сбора логов из Docker-контейнеров                             |
| **grafana**             | grafana/grafana                                                                                 | Дашборды, метрики, трейсы, логи                                     |

## Быстрый старт

```bash
docker-compose up
```

Все сервисы запустятся автоматически. Миграции применятся при первом старте.

## Порты

| Сервис      | Хост-порт | Описание                        |
|-------------|-----------|---------------------------------|
| [product](https://github.com/jva44ka/marketplace-simulator-product)     | 5001      | HTTP (grpc-gateway + REST)      |
| [cart](https://github.com/jva44ka/marketplace-simulator-cart)           | 5002      | HTTP REST                       |
| product-db  | 5433      | PostgreSQL                      |
| cart-db     | 5434      | PostgreSQL                      |
| kafka       | 9092      | Kafka broker                    |
| kafka-ui    | 8090      | Kafka UI                        |
| prometheus  | 9090      | Prometheus UI                   |
| tempo       | 4317      | OTLP gRPC receiver              |
| loki        | 3100      | Loki HTTP API                   |
| grafana     | 3000      | Grafana (admin / admin)         |

## Конфигурация

Файлы конфигурации находятся в `configs/`:

| Файл               | Назначение                               |
|--------------------|------------------------------------------|
| `product.yaml`     | Конфиг сервиса товаров                   |
| `cart.yaml`        | Конфиг сервиса корзины                   |
| `loadgen.yaml`     | Конфиг генератора нагрузки               |
| `prometheus.yml`   | Конфиг Prometheus (scrape jobs)          |
| `tempo.yaml`       | Конфиг Tempo                             |
| `loki.yaml`        | Конфиг Loki                              |
| `promtail.yaml`    | Конфиг Promtail (сбор логов)             |
| `grafana/`         | Provisioning и дашборды Grafana          |

## Observability

| Инструмент | Что собирает           | Адрес                                    |
|------------|------------------------|------------------------------------------|
| Prometheus | Метрики сервисов       | [http://localhost:9090](http://localhost:9090) |
| Tempo      | Распределённые трейсы  | через Grafana                            |
| Loki       | Логи всех контейнеров  | через Grafana                            |
| Grafana    | Единый UI              | [http://localhost:3000](http://localhost:3000) (admin / admin) |

В Grafana предустановлены дашборды:
- **products** — gRPC-запросы, latency, ошибки БД, optimistic lock failures, outbox метрики
- **cart** — HTTP-запросы, latency, ошибки БД, outbox метрики
- **outbox-overview** — сводный дашборд по outbox обоих сервисов

Datasource-связки:
- Из трейса (Tempo) → переход в логи (Loki) по `traceId`
- Из лога (Loki) → переход в трейс (Tempo) по `traceId`
- Из трейса (Tempo) → переход в метрики (Prometheus)

## Архитектура взаимодействия сервисов

```
              HTTP REST
  Client ──────────────► cart :5002
                          │
       gRPC (Docker net.) │
                          ▼
                       product :8002
                          │
                    Kafka outbox
                          ▼
                        kafka :9092
```

При оформлении заказа (`POST /user/{user_id}/cart/checkout`) сервис cart:
1. Вызывает `ReserveProduct` на product — резервирует каждый товар из корзины
2. Сохраняет задачи подтверждения в outbox (в одной транзакции с очисткой корзины)
3. Outbox job асинхронно вызывает `ConfirmReservation` — списывает товары со склада

При пополнении товаров loadgen слушает Kafka-топик `product.events` и вызывает `IncreaseProductCount`, когда остаток падает ниже порога.

## Документация сервисов

- [marketplace-simulator-product](https://github.com/jva44ka/marketplace-simulator-product) — Swagger UI: [http://localhost:5001/swagger/](http://localhost:5001/swagger/), метрики: [http://localhost:5001/metrics](http://localhost:5001/metrics)
- [marketplace-simulator-cart](https://github.com/jva44ka/marketplace-simulator-cart) — Swagger UI: [http://localhost:5002/swagger/](http://localhost:5002/swagger/), метрики: [http://localhost:5002/metrics](http://localhost:5002/metrics)
- [marketplace-simulator-loadgen](https://github.com/jva44ka/marketplace-simulator-loadgen) — генератор нагрузки
