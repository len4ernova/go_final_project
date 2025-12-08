# Планировщик задач
Реализован базовый и дополнительный (*) функционал.

Основной функционал:
- создание;
- редактирование;
- удаление;
- повторение;
- отмана задачи.

Дополнительный:
- считывание данных из переменных окружения (TODO_PASSWORD, TODO_PORT, TODO_DBFILE);
- аутентификация;
- Dockerfile.

Запуск:
 - локально:
 ```
 TODO_PORT=7540 TODO_DBFILE="scheduler.db" TODO_PASSWORD="12345" ...
 ```
 - docker:
   ```
    docker run -p 8080:7540 --name planner:v1
   ```

Шаг 2
Шаг 5

docker build -p planner:latest .
docker run -p 7540:7540 -e TODO..=.. planner:latest
docker run -p 7540:7540 --env-file .env  planner:latest