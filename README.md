# Лабораторная работа № 11
## Контейнеризация мультиязычных приложений

**Студент:** Бойченко Даниэль Дмитриевич  
**Группа:** 220032-11  
**Вариант:** 2

---


### Задания

| № | Описание |
|---|----------|
| 1 | Dockerfile для Go-приложения с многоэтапной сборкой |
| 2 | Собрать образы и сравнить их размеры. |
| 3 | Настроить сеть между контейнерами. |
| 4 | Собрать Rust-приложение с поддержкой musl для полностью статической сборки. |
| 5 | Использовать docker buildx для кросс-платформенной сборки (arm64/amd64). |

---

## Технологии

- **Языки:** Go, Rust, Python
- **Контейнеризация:** Docker, Docker Compose
- **Подходы:** многоэтапная сборка (multi-stage build), scratch-образы, статическая компиляция

---

## Инструкция по сборке и запуску

### Предварительные требования

- [Docker](https://docs.docker.com/get-docker/) установлен и запущен
- Для локальной сборки (без Docker): Go 1.22+, Rust 1.75+, Python 3.12+

---

### Задание 1 — Go-приложение с многоэтапной сборкой

Перейти в папку задания:
```bash
cd task1
```

Собрать Docker-образ:
```bash
docker build -t lr11-task1 .
```

Запустить контейнер:
```bash
docker run -d -p 8080:8080 --name task1-app lr11-task1
```

Проверить работу (доступные эндпоинты):
```bash
curl http://localhost:8080/health
curl "http://localhost:8080/hello?name=World"
```

Остановить и удалить контейнер:
```bash
docker stop task1-app
docker rm task1-app
```

---

### Задание 2 — Собрать образы и сравнить их размеры

Три идентичных по функционалу приложения (эндпоинты `/health` и `/hello`) собраны в Docker-образы с использованием многоэтапной сборки и минимальных базовых образов.

#### Go

```bash
cd task2
docker build -f Dockerfile.go -t lr11-task2-go .
docker run -d -p 8081:8080 --name task2-go lr11-task2-go
curl http://localhost:8081/health
curl "http://localhost:8081/hello?name=Go"
docker stop task2-go && docker rm task2-go
```

#### Python

```bash
cd task2
docker build -f Dockerfile.python -t lr11-task2-python .
docker run -d -p 8082:8080 --name task2-python lr11-task2-python
curl http://localhost:8082/health
curl "http://localhost:8082/hello?name=Python"
docker stop task2-python && docker rm task2-python
```

#### Rust

```bash
cd task2
docker build -f Dockerfile.rust -t lr11-task2-rust .
docker run -d -p 8083:8080 --name task2-rust lr11-task2-rust
curl http://localhost:8083/health
curl "http://localhost:8083/hello?name=Rust"
docker stop task2-rust && docker rm task2-rust
```

#### Сравнение размеров

```bash
docker images --format "table {{.Repository}}\t{{.Size}}" | findstr lr11-task2
```

| Язык | Базовый образ (финальный) | Размер | Многоэтапная сборка |
|------|--------------------------|--------|---------------------|
| Go | `scratch` (0 B) | ~7 MB | ✅ builder → scratch |
| Rust | `alpine:3.19` (~7 MB) | ~17 MB | ✅ rust → alpine |
| Python | `python:3.12-slim` | ~177 MB | ❌ нет |

#### Вывод

Самый маленький образ даёт **Go**, потому что:

1. Статическая компиляция (`CGO_ENABLED=0`) — один бинарник без зависимостей.
2. Финальный образ `scratch` — пустой, 0 байт оверхеда.
3. Флаги `-ldflags="-s -w"` вырезают отладочную информацию.

**Rust** на втором месте — статический бинарник (musl), но финальный образ `alpine` добавляет ~7 MB.

**Python** самый большой — интерпретатор + стандартная библиотека занимают ~170 MB, многоэтапная сборка неприменима (нужен интерпретатор в runtime).

---

### Задание 3 — Настроить сеть между контейнерами

Три сервиса (Go, Python, Rust) объединены в общую Docker-сеть через `docker-compose.yml`. Go-сервер выступает шлюзом: эндпоинт `/status` опрашивает Python и Rust по внутренней сети и возвращает общий статус.

```
┌─────────────┐
│  Go (:8080)  │  ← доступен на хосте
│  /health     │
│  /hello      │
│  /status     │  → опрашивает python:8080 и rust:8080
└──┬───────┬───┘
   │       │  docker network: app-net
   ▼       ▼
┌──────┐ ┌──────┐
│Python│ │ Rust │
└──────┘ └──────┘
```

Запустить все сервисы:
```bash
cd task3
docker compose up -d --build
```

Проверить работу:
```bash
curl http://localhost:8080/health
curl "http://localhost:8080/hello?name=Network"
curl http://localhost:8080/status
```

Проверить сеть:
```bash
docker network ls
docker network inspect task3_app-net
```

Остановить:
```bash
docker compose down
```

---

## Тесты

| Язык | Команда | Кол-во |
|------|---------|--------|
| Go | `cd task3/src/go && go test -v .` | 13 (10 unit + 3 integration) |
| Python | `cd task3 && python -m pytest src/python/tests/test_python.py -v` | 15 |
| Rust | `cd task3/src/rust && cargo test` | 11 |

