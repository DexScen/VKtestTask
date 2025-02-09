# Система мониторинга IP контейнеров

## Описание проекта
**Система мониторинга IP контейнеров** — это веб-приложение для мониторинга доступности IP-адресов контейнеров. Система пингует ip-адреса запущенных контейнеров, отображает время пингов и время последнего успешного пинга, по этим метрикам можно понять доступность ip контейнеров в виде таблицы.

## Функциональность
- **Мониторинг IP-адресов** с сохранением данных в PostgreSQL
- **Отображение данных** в виде таблицы с динамическим обновлением
- **Frontend на React** с использованием Ant Design
- **Docker-контейнеризация** всех сервисов

## Стек технологий
### Backend:
- Go, Mux, PostgreSQL

### Frontend:
- React + Ant, Fetch API, nginx

А также Dockerизация проекта

## Установка и запуск
### 1. Клонирование репозитория
```sh
git clone https://github.com/DexScen/VKtestTask.git
cd VKtestTask
```
### 2. Запуск с Docker Compose
```sh
docker-compose up --build -d
```
После этого сервисы будут запущены:
- **Backend**: `http://localhost:8080`
- **Frontend**: `http://localhost:3000`

Увидеть таблицу можно вбив http://localhost:3000 в адрес браузера =)

## API Endpoints
### Получить список IP-адресов
```
GET /containers
```
**Пример ответа:**
```json
[
  {"ip":"172.19.0.4", "pingtime":"2025-02-08T13:19:36.859Z", "successdate":"2025-02-08T13:19:36.859Z"},
  {"ip":"172.19.0.3", "pingtime":"2025-02-08T13:19:36.860Z", "successdate":"2025-02-08T13:19:36.860Z"}
]
```

### Отправить список IP-адресов
```
POST /containers
```
```json
[
  {"ip":"172.19.0.4", "pingtime":"2025-02-08T13:19:36.859Z", "successdate":"2025-02-08T13:19:36.859Z"},
  {"ip":"172.19.0.3", "pingtime":"2025-02-08T13:19:36.860Z", "successdate":"2025-02-08T13:19:36.860Z"}
]
```
**Пример ответа:**
```
OK
```

### Остановка и удаление контейнеров
```sh
docker-compose down -v
```
### Запуск контейнеров с пересборкой
```sh
docker-compose up --build -d
```

## Автор
**Александр Самарцев** – [GitHub - DexScen](https://github.com/DexScen)
