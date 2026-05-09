# marketplace-simulator

[Русский](README.md) · en English

Pet project "Marketplace Simulator" — two Go microservices with a load generator and a full observability stack.

## About

The system consists of two microservices and one load generator:
1. **[cart](https://github.com/jva44ka/marketplace-simulator-cart)** — http-service, master system for user shopping carts. Endpoints for:
   1. Adding / removing items from the cart
   2. Clearing the cart
   3. Viewing the current cart contents
   4. Checking out an order
2. **[product](https://github.com/jva44ka/marketplace-simulator-product)** — grpc-service, master system for products / inventory. Endpoints for increasing / decreasing stock counts and querying product data.
3. **[loadgen](https://github.com/jva44ka/marketplace-simulator-loadgen)** — load generator for cart and product.

Supporting infrastructure:
1. **PostgreSQL** for persistent storage
2. **Redis** for caching in the product service
3. **Kafka** for async product events
4. **etcd** for live configuration changes without restarts (e.g. adjusting the rate limiter in product)
5. **Grafana** for dashboards, alerts, logs, and traces
6. **Prometheus, Tempo, Loki** for metrics storage, distributed traces, and logs

### Main flows

**Add item to cart**

The client sends `POST /user/{id}/cart/{sku}` to cart. Cart calls product via gRPC (`GetProduct`) to fetch price and name. Product serves data from Redis (Cache-Aside); on a cache miss it reads from PostgreSQL. Cart saves the item to its own database and returns a response.
<img width="944" height="364" alt="image" src="https://github.com/user-attachments/assets/6b027bd0-640d-4f21-9eab-bcda9e89b959" />

**Checkout**

The client sends `POST /user/{id}/cart/checkout` to cart. Cart fetches the cart items from the database, then calls `ReserveProduct` on product (product creates reservation records — stock counts are not changed yet). Cart then atomically clears the cart and writes confirmation tasks to an outbox table in a single transaction. A background job asynchronously calls `ConfirmReservation` on product — product deducts stock, deletes the reservations, and publishes events to Kafka. The loadgen replenisher reads these events and restocks inventory when the count drops below a threshold.
<img width="1564" height="857" alt="image" src="https://github.com/user-attachments/assets/6e5c83db-60e7-4ee3-8ef0-697e1363172f" />

## Quick start

```bash
git clone https://github.com/jva44ka/marketplace-simulator.git
cd marketplace-simulator
docker-compose up
```

After startup:

|                    |                                                                  |
|--------------------|------------------------------------------------------------------|
| Grafana            | [http://localhost:3000](http://localhost:3000) (admin / admin)   |
| Prometheus         | [http://localhost:9090](http://localhost:9090)                   |
| Kafka UI           | [http://localhost:8090](http://localhost:8090)                   |
| etcd UI            | [http://localhost:8091](http://localhost:8091)                   |
| Swagger — product  | [http://localhost:5001/swagger/](http://localhost:5001/swagger/) |
| Swagger — cart     | [http://localhost:5002/swagger/](http://localhost:5002/swagger/) |

## Architectural patterns

**Transactional Outbox** — cart atomically clears the cart and writes reservation confirmation tasks in a single transaction. Product similarly atomically updates stock and writes Kafka events. No message loss on service crash.

**Circuit Breaker** — cart's gRPC client to product is protected by a circuit breaker: once the error threshold is exceeded, the circuit opens and requests are rejected immediately without waiting for a timeout. Configurable: threshold, window, minimum traffic to trip.

**Retry + Exponential Backoff + Jitter** — transient gRPC failures to product are retried with exponential delay and random jitter to avoid thundering herd on recovery.

**Dead Letter Queue** — outbox records that fail to deliver after N attempts are moved to dead letter and do not block the queue. The counter is monitored in Grafana with an alert.

**Optimistic Locking** — concurrent stock updates in product are protected by optimistic locking on the record version. Version conflicts are tracked with a dedicated metric.

**Cache-Aside (Redis)** — `GetProduct` reads are served from Redis; on a miss the service reads from the database and warms the cache asynchronously. After stock updates, product writes to `outbox.cache_updates` and a background job invalidates the cache. Stale entries are detected by comparing PostgreSQL xmin. When Redis is unavailable the service continues without cache (graceful degradation).

**Read-Your-Writes (RYW)** — product publishes messages to the `product.events` topic with a `transaction_id` field. Subsequent requests can pass this field to `GET /v1/products/{sku}` to ensure they receive the up-to-date version. If the cached entry is stale, the service falls back to the database. This guarantees a client using `transactionId` will never read a product state older than the event it already received.

**Rate Limiter** — incoming traffic to product is limited by a token bucket, protecting the database from overload during traffic spikes.

**Drain Mode** — outbox jobs switch between two intervals: a minimal pause when the queue is empty and an immediate next tick while there are records pending, avoiding wasteful idle polling of the database.

---

→ [Detailed documentation](docs/README.en.md)
