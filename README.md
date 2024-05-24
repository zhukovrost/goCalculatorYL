# Распределенный вычислитель арифметических выражений (оркестратор)

Финальное задание по второму спринту Яндекс Лицея (GoLang)

## Требуется установка агента

Подробнее тут: https://github.com/zhukovrost/agentYL.git

## Установка

1. Клонируйте репозиторий:
    ```sh
    git clone https://github.com/zhukovrost/orchestratorYL.git
    ```

2. Перейдите в директорию проекта:
    ```sh
    cd orchestratorYL
    ```

3. Установите зависимости:
    ```sh
    go mod tidy
    ```

## Запуск

Для запуска сервера выполните:

```sh
go run cmd/orchestrator/main.go
```

## Инструкция по использованию

### 1. Добавление арифметического выражения

```sh
curl http://localhost:8080/api/v1/calculate \
--header 'Content-Type: application/json' \
--data '{
  "id": "12342",
  "expression": "2 + 2"
}'
```

### 2. Получение всех алгоритмических выражений

```sh 
curl http://localhost:8080/api/v1/expressions
```

### 3. Получение алгоритмических выражений по ID

Дан пример получения выражения с ID 12342

```sh
curl http://localhost:8080/api/v1/expressions/12342
```

### 4. Получение задачи

```sh
curl http://localhost:8080/internal/task
```

### 5. Установить результат задачи

```sh
curl http://localhost:8080/internal/task \
--header 'Content-Type: application/json' \
--data '{
  "id": 0,
  "result": 4
}'
```

## Структура проекта

```
orchestratorYL/
├── cmd/
│   └── orchestrator/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   └── handler.go
│   ├── router/
│   │   └── router.go
│   └── service/
│       └── service.go
├── pkg/
│   └── utils/
│       └── utils.go
├── .gitignore
├── go.mod
└── README.md
```
