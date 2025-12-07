# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

Реализованные *:
Шаг 1
Шаг 2
Шаг 5

docker build -p planner:latest .
docker run -p 7540:7540 -e TODO..=.. planner:latest
docker run -p 7540:7540 --env-file .env  planner:latest