# OpsUser

Модуль предназначен для получения информации о пользователях в Unix-системах. 
Он предоставляет функции для получения данных о текущем пользователе или пользователе по имени, включая UID, имя пользователя и домашнюю директорию.

## Возможности
- Получение информации о текущем пользователе (`GetCurrent`).
- Получение информации о пользователе по username (`GetByUsername`). 
- Проверка, является ли текущий пользователь суперпользователем (`IsRoot`).
- Проверка наличия пользователя по username (`CheckExists`).

## Требования
- Go 1.20 или выше.
- Unix-подобная ОС (Linux, macOS и т.д.).

## Установка
Склонируйте репозиторий или добавьте модуль в ваш проект:
```shell
go get github.com/mosregdata/ops-user
```

## Использование
Пример использования модуля:
```go
package main

import (
    "fmt"
    opsuser "github.com/mosregdata/ops-user"
)

func main() {
    // Получить информацию о текущем пользователе
    current, err := opsuser.GetCurrent()
    if err != nil {
        fmt.Printf("Ошибка: %v\n", err)
        return
    }
    fmt.Printf("Текущий пользователь: %s (UID: %s, Home: %s)\n",
        current.Username, current.UID, current.HomeDir)

    // Получить информацию о пользователе по имени
    info, err := opsuser.GetByUsername("root")
    if err != nil {
        fmt.Printf("Ошибка: %v\n", err)
        return
    }
    fmt.Printf("Пользователь: %s (UID: %s, Home: %s)\n",
        info.Username, info.UID, info.HomeDir)

    // Проверка, является ли текущий пользователь суперпользователем
    msg := "Программа запущена от имени обычного пользователя"
    if opsuser.IsRoot() {
        msg = "Программа запущена от имени root"
    }
    fmt.Println(msg)

    // Проверка наличия пользователя
    username := "some"
    mess := fmt.Sprintf("Пользователь %s существует\n", username)
    _, err = opsuser.CheckExists(username)
    if err != nil {
        mess = fmt.Sprintf("Пользователь %s не существует\n", username)
    }
    fmt.Println(mess)
}
```

## Результат
Код примера представленный выше вернет такой результат:
```
Текущий пользователь: admin (UID: 1000, Home: /home/admin)
Пользователь: root (UID: 0, Home: /root)
Программа запущена от имени обычного пользователя
Пользователь some не существует
```