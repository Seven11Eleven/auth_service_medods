# Test task BackDev

Тестовое задание на позицию Junior Backend Developer

**Используемые технологии:**
**• Go** 
**• JWT**
**• PostgreSQL**
**Задание:**

Написать часть сервиса аутентификации.

Два REST маршрута:
Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов

Требования:

Access токен тип JWT, алгоритм SHA512, хранить в базе строго запрещено.

Refresh токен тип произвольный, формат передачи base64, хранится в базе исключительно в виде bcrypt хеша, должен быть защищен от изменения на стороне клиента и попыток повторного использования.

Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена можно выполнить только тем Refresh токеном который был выдан вместе с ним.

Payload токенов должен содержать сведения об ip адресе клиента, которому он был выдан. В случае, если ip адрес изменился, при рефреш операции нужно послать email warning на почту юзера (для упрощения можно использовать моковые данные).



## API Эндпоинты
# 1. Регистрация пользователя

    Эндпоинт: POST /signup

    Тело запроса:

```json
{
    "Username": "vasya",
    "Email": "vasiliy@gmail.com",
    "Password": "vasiliy123"
}
```

# 2. Логин пользователя

    Эндпоинт: POST /login

    Тело запроса:

```json

{
    "Username": "vasya",
    "Email": "vasiliy@gmail.com",
    "Password": "vasiliy123"
}
```
# 3. "Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса"
    Эндпоинт: GET /tokenByGUID?userID=<вашАйди>

# 4. Рефреш токенов 

    Эндпоинт: POST /refresh

    Тело запроса:

```json

{
    "refresh_token": "токен"
}
```