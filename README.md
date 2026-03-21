# ozon-simulator-go

Оркестрирующий репозиторий учебного проекта «Симулятор Ozon».

Запускает всю инфраструктуру через `docker-compose`: оба микросервиса, базы данных, мигратор, Prometheus и Grafana.

## Состав системы

| Сервис                | Репозиторий                        | Описание                                        |
|-----------------------|------------------------------------|-------------------------------------------------|
| **products**          | [ozon-simulator-go-products](https://github.com/jva44ka/ozon-simulator-go-products) | Управление товарами (gRPC + REST, PostgreSQL)   |
| **cart**              | [ozon-simulator-go-cart](https://github.com/jva44ka/ozon-simulator-go-cart)         | Корзина покупок (REST, PostgreSQL)              |
| **products-db**       | postgres:17.7                      | БД сервиса товаров                              |
| **cart-db**           | postgres:17.7                      | БД сервиса корзины                              |
| **products-migrations** | migrator из [products](https://github.com/jva44ka/ozon-simulator-go-products)    | Накатывает миграции в products-db при старте    |
| **cart-migrations**   | migrator из [cart](https://github.com/jva44ka/ozon-simulator-go-cart)              | Накатывает миграции в cart-db при старте        |
| **prometheus**        | prom/prometheus                    | Сбор метрик с обоих сервисов                    |
| **grafana**           | grafana/grafana                    | Дашборды для visualize метрик                   |

## Быстрый старт

```bash
docker-compose up
```

Все сервисы запустятся автоматически. Миграции применятся при первом старте.

## Порты

| Сервис     | Хост-порт | Описание                   |
|------------|-----------|----------------------------|
| [products](https://github.com/jva44ka/ozon-simulator-go-products) | 5001 | HTTP (grpc-gateway + REST) |
| [cart](https://github.com/jva44ka/ozon-simulator-go-cart)         | 5002 | HTTP REST                  |
| products-db| 5433      | PostgreSQL                 |
| cart-db    | 5434      | PostgreSQL                 |
| prometheus | 9090      | Prometheus UI              |
| grafana    | 3000      | Grafana (admin / admin)    |

## Конфигурация сервисов

Файлы конфигурации находятся в `configs/`:

| Файл                  | Назначение                        |
|-----------------------|-----------------------------------|
| `products.yaml`       | Конфиг сервиса товаров            |
| `cart.yaml`           | Конфиг сервиса корзины            |
| `prometheus.yml`      | Конфиг Prometheus (scrape jobs)   |
| `grafana/`            | Provisioning и дашборды Grafana   |

## Мониторинг

Prometheus автоматически собирает метрики с обоих сервисов каждые 15 секунд.

В Grafana предустановлены два дашборда:
- **[products](https://github.com/jva44ka/ozon-simulator-go-products)** — метрики сервиса товаров (gRPC-запросы, время ответа, ошибки БД, optimistic lock failures)
- **[cart](https://github.com/jva44ka/ozon-simulator-go-cart)** — метрики сервиса корзины (HTTP-запросы, время ответа, ошибки БД)

Grafana: [http://localhost:3000](http://localhost:3000) — логин `admin`, пароль `admin`.

## Архитектура взаимодействия сервисов

```
           HTTP REST
  Client ──────────────► cart :5002
                           │
                           │ gRPC (внутренняя сеть Docker)
                           ▼
                      products :8002
                           │
                           │ SQL
                           ▼
                      products-db :5432
  cart ──────────────► cart-db :5432
```

При оформлении заказа (`POST /user/{user_id}/cart/checkout`) сервис [cart](https://github.com/jva44ka/ozon-simulator-go-cart):
1. Запрашивает данные товаров из `products` по gRPC
2. Вызывает `DecreaseProductCount` для списания со склада
3. Очищает корзину пользователя

## Документация сервисов

- [Products](https://github.com/jva44ka/ozon-simulator-go-products) Swagger UI: [http://localhost:5001/swagger/](http://localhost:5001/swagger/)
- [Cart](https://github.com/jva44ka/ozon-simulator-go-cart) Swagger UI: [http://localhost:5002/swagger/](http://localhost:5002/swagger/)
- [Products](https://github.com/jva44ka/ozon-simulator-go-products) метрики: [http://localhost:5001/metrics](http://localhost:5001/metrics)
- [Cart](https://github.com/jva44ka/ozon-simulator-go-cart) метрики: [http://localhost:5002/metrics](http://localhost:5002/metrics)
