# marketplace-simulator

Оркестрирующий репозиторий учебного проекта «Симулятор маркетплейса».

Запускает всю инфраструктуру через `docker-compose`: микросервисы, базы данных, генератор нагрузки и observability-стек (Prometheus, Grafana, Tempo, Loki).

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

→ [Подробная документация](docs/README.md)
