
# mygrep

Распределённая CLI-утилита, аналог `grep`, с поддержкой кворума и конкурентности.

## Особенности

- Поддержка флагов: `-i`, `-v`, `-n`, `-c`, `-A`, `-B`, `-C`, `-F`.
- Работа в режиме **координатор-воркер**.
- **Кворум**: результат считается валидным, если `N/2 + 1` воркеров ответили.
- **Параллельная обработка** данных между воркерами.
- Совместимость с поведением оригинальной утилиты `grep`.

## Установка

Клонируйте репозиторий:

```bash
git clone https://github.com/v1adis1av28/level4
cd mygrep
```

## Запуск

### Запуск воркеров

Запустите воркеров на разных портах:

```bash
go run cmd/worker/main.go --mode=worker --port=8080
go run cmd/worker/main.go --mode=worker --port=8081
go run cmd/worker/main.go --mode=worker --port=8082
```

### Запуск координатора

Координатор принимает паттерн первым аргументом, затем флаги:

```bash
echo -e "line1\nerror line2\nline3\nerror line4" | \
go run cmd/coordinator/main.go \
  error \
  -n \
  -workers="localhost:8080,localhost:8081,localhost:8082" \
  --quorum=2
```

### Скрипт запуска

Для быстрого запуска с тремя воркерами:

```bash
./examples/run.sh
```

## Тестирование

Запуск всех тестов:

```bash
go test ./...
```

Запуск тестов с подробным выводом:

```bash
go test -v ./...
```


## Сборка бинарников

Собрать воркер и координатор:

```bash
go build -o bin/worker cmd/worker/main.go
go build -o bin/coordinator cmd/coordinator/main.go
```

## Структура проекта

```
mygrep/
├── cmd/
│   ├── coordinator/ — точка входа для координатора
│   └── worker/      — точка входа для воркера
├── internal/
│   ├── config/      — парсинг флагов
│   ├── grep/        — логика фильтрации строк
│   └── network/     — HTTP-взаимодействие
├── examples/
│   └── run.sh       — скрипт запуска
├── go.mod
└── README.md
```