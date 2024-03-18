# vk-film-library
Тестовое задание VK | Бэкенд приложения "Фильмотека"

Для запуска сервиса - команда make compose-up    
Документацию после запуска можно посмотреть по адресу http://localhost:8080/swagger/index.html

## Некоторые примеры запросов

### Регистрация
Регистрация администратора:
```curl
curl -X 'POST' \
  'http://localhost:8080/admin/signup' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "password": "admin2",
  "username": "admin2"
}'
```
Пример ответа:
```json
{"id":2}
```

### Аутентификация
Аутентификация для получения токена доступа:
```curl
curl -X 'POST' \
  'http://localhost:8080/signin' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "password": "admin1",
  "username": "admin1"
}'
```
Пример ответа:
```json
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTA3OTM2MjUsImlhdCI6MTcxMDc4NjQyNSwiVXNlclJvbGUiOiJhZG1pbiJ9.KaUTvr85RM1807pS4ScurgdaHZS9tIeqC1x2YEHbFRk"}
```

### Получение фильмов, отсортированных по рейтингу
```curl
curl -X 'POST' \
  'http://localhost:8080/api/v1/films/sorted' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTA3OTM2MjUsImlhdCI6MTcxMDc4NjQyNSwiVXNlclJvbGUiOiJhZG1pbiJ9.KaUTvr85RM1807pS4ScurgdaHZS9tIeqC1x2YEHbFRk' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "rating"
}'
```
Пример ответа:
```json
{"films":[
  {"Id":3,"name":"string","description":"string","created_at":"2000-01-01","rating":5,"Actors":["asher"]},
  {"Id":1,"name":"murder","description":"string","created_at":"2010-01-01","rating":7,"Actors":["asher"]},
  {"Id":2,"name":"murder2","description":"string","created_at":"2015-01-01","rating":8,"Actors":["asher"]}
]}
```