# marketplace-simulator

Учебный проект «Симулятор маркетплейса» — два микросервиса на Go с генератором нагрузки и полным observability-стеком.

```bash
git clone https://github.com/jva44ka/marketplace-simulator.git
cd marketplace-simulator
docker-compose up
```

После запуска:

| | |
|-|-|
| Grafana | [http://localhost:3000](http://localhost:3000) (admin / admin) |
| Prometheus | [http://localhost:9090](http://localhost:9090) |
| Kafka UI | [http://localhost:8090](http://localhost:8090) |
| Swagger — product | [http://localhost:5001/swagger/](http://localhost:5001/swagger/) |
| Swagger — cart | [http://localhost:5002/swagger/](http://localhost:5002/swagger/) |

## Использованные паттерны для улучшения stability

**Transactional Outbox** — cart атомарно очищает корзину и записывает задачи подтверждения резервирований в одной транзакции. Product так же атомарно обновляет остатки и пишет события в Kafka. Никаких потерь сообщений при падении сервиса.

**Circuit Breaker** — gRPC-клиент cart к product защищён автоматом: при превышении порога ошибок цепь размыкается и запросы отклоняются сразу, не тратя время на таймаут. Настраивается: порог, окно, минимальный трафик для срабатывания.

**Retry + Exponential Backoff + Jitter** — временные сбои gRPC-вызовов к product переповторяются с экспоненциальной паузой и случайным отклонением, чтобы не создавать thundering herd при восстановлении.

**Dead Letter Queue** — outbox-записи которые не удалось доставить после N попыток уходят в dead letter и не блокируют очередь. Счётчик мониторится в Grafana с алертом.

**Optimistic Locking** — конкурентное обновление остатков товаров в product защищено оптимистичной блокировкой по версии записи. Конфликты версий отслеживаются отдельной метрикой.

**Rate Limiter** — входящий трафик на product-сервис ограничен token bucket'ом, защищая БД от перегрузки при всплесках нагрузки.

**Drain Mode** — outbox jobs переключаются между двумя интервалами: минимальная пауза в idle и немедленный следующий тик пока в очереди есть записи, не нагружая базу холостыми запросами.

---

→ [Подробная документация](docs/README.md)
