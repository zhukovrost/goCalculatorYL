# Распределенный вычислитель арифметических выражений

Финальное задание по второму спринту Яндекс Лицея (GoLang)

## Установка

1. Клонируйте репозиторий:
    ```sh
    git clone https://github.com/zhukovrost/goCalculatorYL.git
    ```

2. Перейдите в директорию проекта:
    ```sh
    cd goCalculatorYL
    ```

3. Установите зависимости:
    ```sh
    go mod tidy
    ```

## Запуск

Для запуска сервера выполните:

```sh
go run cmd/calculator/main.go
```

## Инструкция по использованию



## Структура проекта

```
mywebsite/
├── cmd/
│   └── calculator/
│       └── main.go
├── internal/
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
├── web/
│   ├── static/
│   └── templates/
├── .gitignore
├── go.mod
└── README.md
```
