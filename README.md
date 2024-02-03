# Project Small Links

Проект по уменьшению длины ссылок. Хранение данных производится в БД PostgreSQL.

```sql
    CREATE TABLE links (
        slot_id SERIAL PRIMARY KEY,
        origlink VARCHAR(300) UNIQUE NOT NULL,
        customlink VARCHAR(30) NOT NULL
    )
```

Перед запуском **необходимо**:
* ввести данные в ```config.yaml```
* создать ```таблицу links```
* [Опционально] добавить задний фон для интерфейса