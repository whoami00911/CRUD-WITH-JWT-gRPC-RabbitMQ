## CRUD Приложение предоставляющее web API к данным
### Стэк
- **Go**: Версия 1.22.5

### Начало работы

Для запуска программы используйте следующую команду:

```bash
go run cmd/main.go
```

### Использование
**Для легкого взаимодействия с API вы можете перейти на:**
```http://localhost:8080/swagger/index.html```

**Пример:**
- **Зарегистрировать пользователя:** localhost:8080/auth/signUp [POST]
```
{
    "name": "user",
    "username": "User123",
    "password": "password"
}
```
Ответ:
```
{
    "id": 16
}
```
- **Залогиниться:** localhost:8080/auth/signIn [POST]
```
{
    "username": "User123",
    "password": "password"
}
```
Ответ:
```
{
    "RefreshToken": "z3ZFlUPhhOO34naL3gflp-P59iIRsD7Cel3B5bYdLvU=",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzg3NzgwNTIsImlhdCI6MTczODc3Nzk5MiwiVXNlcklkIjoxNn0.yXgqDK6hwsB8tvAHlED4DLt7qJo27p7_NFNJlyBR0FA"
}
```
- **Обновить токены:** localhost:8080/auth/refresh [GET]
Ответ:
```
{
    "RefreshToken": "-RQy4qOIXwqrQHMmhV7bPQWsJQlR_PLHQt9g8s1wi44=",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzg3ODE5ODYsImlhdCI6MTczODc4MTkyNiwiVXNlcklkIjoxNn0.amzL5v7blv1XPvxqb3H94gT37uigsmTFJVc0POoVvlg"
}
```
- **Добавить сущность:**
Добавить сущность с уже существующим IP нельзя ни в кеш ни в БД
```
curl -v --request POST "http://localhost:8080/Abuseip/" -H "Content-Type: application/json" -d "{\"ipAddress\": \"192.25.6.39\", \"isPublic\": true, \"ipVersion\": 4, \"isWhitelisted\": false, \"abuseConfidenceScore\": 100, \"countryCode\": \"CN\", \"countryName\": \"China\", \"usageType\": \"WebHosting\", \"isp\": \"ChinChinarem\"}" + ВАШ ТОКЕН
```

- **Получить все сущности из БД:**
```
curl -v --request GET http://localhost:8080/Abuseip/ + ВАШ ТОКЕН
```

- **Получить одну сущность из кеша(Хранится 15 секунд после добавления) или БД:**
- Комманда: ```curl -v --request GET http://localhost:8080/Abuseip/192.25.6.39 + ВАШ ТОКЕН```
- Вывод:
Из БД: "isFromDB":true | "IsCache":false
```
{"ipAddress":"192.25.6.39","isPublic":true,"ipVersion":4,"isWhitelisted":false,"abuseConfidenceScore":100, <...>  "IsCache":false,"isFromDB":true}*
```
Из кеша: "IsCache":true | "isFromDB":false
```
{"ipAddress":"192.25.6.39","isPublic":true,"ipVersion":4,"isWhitelisted":false,"abuseConfidenceScore":100, <...> "IsCache":true,"isFromDB":false}*
```

- **Удалить все сущности в БД и кеше:**
```
curl -v --request DELETE http://localhost:8080/Abuseip/ + ВАШ ТОКЕН
```
- **Удалить сущность по IP:**
```
curl -v --request GET http://localhost:8080/Abuseip/192.25.6.39 + ВАШ ТОКЕН
```

- **Обновить сущность по IP**
```
curl -v --request PUT "http://localhost:8080/Abuseip/" -H "Content-Type: application/json" -d "{\"ipAddress\": \"192.25.6.39\", \"isPublic\": true, \"ipVersion\": 4, \"isWhitelisted\": false, \"abuseConfidenceScore\": 100, \"countryCode\": \"CN\", \"countryName\": \"Chinazes\", \"usageType\": \"WebHosting\", \"isp\": \"ChinChinarem\"}" + ВАШ ТОКЕН
```
