
## Описание

На основе проекта из L0 была проведена попытка профилировки:

- Профилирование CPU и памяти
- Нагрузочное тестирование
- Benchmarking с анализом производительности

---

## Структура проекта

Структура проекта аналогична проекту из L0, однако была добавлена папка с профилями, которые были получена в рамках выполнения задания(pprof файлы где показаны метрики по аллокациям на куче и нагрущке на CPU)

---




##  Нагрузочное тестирование

Используется `wrk` для тестирования производительности эндпоинтов(тестил только один эндпоинт по получению заказа по айди):

```bash
wrk -t8 -c200 -d30s http://localhost:8080/order/test
```

Результаты показывают количество запросов в секунду, среднюю задержку и ошибки.

---

## 3. Benchmarking

В папке `i/benchmarks` реализованы Go-бенчмарки:

```bash
go test -bench=. -benchmem ./internal/benchmarks > before.txt
# проведена оптимизация и запускалась 
go test -bench=. -benchmem ./internal/benchmarks > after.txt
```

Для сравнения результатов используется `benchstat`:

```bash
benchstat before.txt after.txt > benchstat.txt
```

Пример результата:

| Metric    | Before  | After   | Change  |
| --------- | ------- | ------- | ------- |
| sec/op    | 1.473µ  | 1.454µ  | -1.29%  |
| B/op      | 1.644Ki | 1.456Ki | -11.41% |
| allocs/op | 14      | 13      | -7.14%  |

> Это показывает улучшение производительности и снижение аллокаций после оптимизации.

---

## Как воспроизвести

1. Запуск контейнеров:

```bash
docker-compose up --build
```

2. Снятие профилей:

```bash
go tool pprof "http://localhost:6060/debug/pprof/profile?seconds=30"
go tool pprof "http://localhost:6060/debug/pprof/heap"
```

3. Нагрузочное тестирование:

```bash
wrk -t8 -c200 -d30s http://localhost:8080/order/test
```

4. Benchmarking:

```bash
go test -bench=. -benchmem ./benchmarks
benchstat before.txt after.txt
```


---

## Как результат

* Снизилось время выполнения эндпоинта `GetOrder` на ~1.3%
* Аллокация памяти снизилась на 11.4%
