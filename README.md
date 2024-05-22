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
│   ├── utils/
│   │   └── utils.go
│   └── middleware/
│       └── middleware.go
├── .gitignore
├── go.mod
└── README.md
```
