# SocketIO Mini Service

[![en](https://img.shields.io/badge/lang-en-yellow.svg)](https://github.com/bydanovm/socketio-mini-service/blob/main/README.md)

SocketIO Mini Service - это минималистичный сервис для работы с Socket.IO, созданный для упрощения разработки приложений реального времени. Этот сервис предоставляет обобщенный интерфейс для управления подключениями Socket.IO с поддержкой пользовательской авторизации и обработки событий.

Основа сервиса построена на [Socket.IO от zishang520](https://github.com/zishang520/socket.io).

## Основные возможности

- Управление подключениями Socket.IO с поддержкой обобщенных типов
- Встроенная система авторизации через токены
- Поддержка пользовательских middleware
- Потокобезопасное управление соединениями
- Гибкая система обработки событий
- Поддержка CORS
- Настраиваемые параметры таймаутов и буферизации

## Установка

```bash
go get github.com/bydanovm/socketio-mini-service
```

## Использование

Пример создания базового сервиса:

```go
package main

import (
    "context"
    "net/http"
    "log"
    
    "github.com/bydanovm/socketio-mini-service"
)

type User struct {
    ID   string
    Name string
}

func (u *User) GetId() string {
    return u.ID
}

func main() {
    // Создание нового сервиса
    service := socketiominiservice.NewSocketIOService[string]()
    
    // Добавление обработчика событий
    service.AddEvent(func(ctx context.Context) (string, func(...any)) {
        return "message", func(data ...any) {
            // Обработка сообщения
        }
    })
    
    // Запуск сервиса
    if err := service.Run(); err != nil {
        log.Fatal(err)
    }
    
    // Запуск HTTP сервера
    http.Handle("/socket.io/", service.Handle())
    log.Fatal(http.ListenAndServe(":3000", nil))
}
```

## Лицензия

Этот проект лицензирован в соответствии с Apache License 2.0 - подробности в файле [LICENSE](LICENSE).

## Главный разработчик

**bydanovm** - основной разработчик минималистичного сервиса.

## Контрибьюторы

Проект открыт для вкладов от сообщества разработчиков.