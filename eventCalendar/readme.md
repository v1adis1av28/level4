## Развернуть проект через docker compose

Клонируйте репозиторий:

```bash
git clone https://github.com/v1adis1av28/level4
cd EventCalendar
```
    #Чтобы запустить проект выполените:
    docker-compose up --build
```
Приложение будет доступно по адресу `http://localhost:8080` (или другому порту, указанному в компоузе или конфиг файле).
Для остановки нажми `Ctrl+C` в терминале или выполните `docker-compose down`.



## Структура проекта

```
eventCalendar/
├── README.md                   
├── go.mod, go.sum              
├── main.go                     # Точка входа приложения
├── config/
│   └── dev.yml                 # Конфигурация здесь можно поменять на то что надо
├── migrations/                 # Файлы миграции
├── internal/
│   ├── workers/                # воркеры(отслеживание архивирования и уведомлений)
│   ├── config/                 
│   ├── server/                 
│   ├── handlers/               
│   ├── storage/                
│   ├── models/                 
│   ├── utils/                  
│   └── logger/                 #  ассинхронный логгер
├── docker-compose.yml          
│__ Dockerfile                  


## API 


### Создание события

*   **URL:** `/create_event`
*   **Метод:** `POST`
*   **Тело запроса:**
    ```json
    {
        "name": "Название события",
        "description": "Описание события",
        "date": "2025-12-25", // Формат: YYYY-MM-DD
        "status": "scheduled", // Допустимые значения: "scheduled"
        "have_notification": true // true или false для указания необходимости напоминания
    }
    ```
*   **Ответ:**
    *   `200 OK`:
        ```json
        {
            "code": 200,
            "message": "Event created successfully",
            "result": {
                "id": 1,
                "name": "Название события",
                "description": "Описание события",
                "date": "2025-12-25",
                "status": "scheduled",
                "have_notification": true
            }
        }
        ```
    *   `400 Bad Request`: При неверных данных или формате даты.
    *   `500 Internal Server Error`: При ошибке сервера.

### Обновление события

*   **URL:** `/update_event`
*   **Метод:** `POST`
*   **Тело запроса:**
    ```json
    {
        "id": 1,
        "new_date": "2025-12-26", // (опционально)
        "new_name": "Новое название" // (опционально)
    }
    ```
*   **Ответ:**
    *   `200 OK`:
        ```json
        {
            "code": 200,
            "message": "succesfully update event",
            "result": {
                "id": 1,
                "new_date": "2025-12-26",
                "new_name": "Новое название"
            }
        }
        ```
    *   `400 Bad Request`: При неверных данных.
    *   `500 Internal Server Error`: При ошибке сервера.

### Удаление события

*   **URL:** `/delete_event`
*   **Метод:** `POST`
*   **Тело запроса:**
    ```json
    {
        "id": 1
    }
    ```
*   **Ответ:**
    *   `200 OK`:
        ```json
        {
            "code": 200,
            "message": "Event deleted successfully"
        }
        ```
    *   `400 Bad Request`: При отсутствии ID или неверном формате.
    *   `404 Not Found`: Если событие не найдено.
    *   `500 Internal Server Error`: При ошибке сервера.

### Получение событий за день

*   **URL:** `/events_for_day?date=2025-12-25`
*   **Метод:** `GET`
*   **Параметры строки запроса:**
    *   `date` (обязательный): Дата в формате `YYYY-MM-DD`.
*   **Ответ:**
    *   `200 OK`:
        ```json
        {
            "code": 200,
            "message": "Events retrieved successfully",
            "result": [
                {
                    "id": 1,
                    "name": "Событие 1",
                    "description": "Описание 1",
                    "date": "2025-12-25T00:00:00Z",
                    "status": "scheduled",
                    "have_notification": false
                }
            ]
        }
        ```
    *   `400 Bad Request`: При отсутствии параметра `date` или неверном формате.
    *   `500 Internal Server Error`: При ошибке сервера.

### Получение событий за неделю

*   **URL:** `/events_for_week?date=2025-12-25`
*   **Метод:** `GET`
*   **Параметры строки запроса:**
    *   `date` (обязательный): Дата в формате `YYYY-MM-DD`, по которой определяется неделя (с понедельника).
*   **Ответ:**
    *   `200 OK`:
        ```json
        {
            "code": 200,
            "message": "Events retrieved successfully",
            "result": [...]
        }
        ```
    *   `400 Bad Request`: При отсутствии параметра `date` или неверном формате.
    *   `500 Internal Server Error`: При ошибке сервера.

### Получение событий за месяц

*   **URL:** `/events_for_month?date=2025-12-25`
*   **Метод:** `GET`
*   **Параметры строки запроса:**
    *   `date` (обязательный): Дата в формате `YYYY-MM-DD`, по которой определяется месяц.
*   **Ответ:**
    *   `200 OK`:
        ```json
        {
            "code": 200,
            "message": "Events retrieved successfully",
            "result": [...]
        }
        ```
    *   `400 Bad Request`: При отсутствии параметра `date` или неверном формате.
    *   `500 Internal Server Error`: При ошибке сервера.

## Запуск тестов

### Юнит-тесты

Для запуска всех юнит-тестов (без `testcontainers`):

```bash
go test ./...
```


### Интеграционные тесты

Для запуска интеграционных тестов, которые используют `testcontainers`:

1.  Нужно чтоьы был запущен докер.
2.  Запустить:

    ```bash
    go test -tags=integration ./...
    ```
