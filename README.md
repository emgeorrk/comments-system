# commentsSystem
### Запуск
1. Выбрать режим хранилища в файле `docker-compose.yam`:
```yaml
STORAGE_TYPE: postgres # postgres | in-memory
```
2. Запустить docker-compose:
```bash
docker-compose up -d
```
3. Сервер будет запущен по адресу [http://localhost:8084/graphql](http://localhost:8084/graphql)
