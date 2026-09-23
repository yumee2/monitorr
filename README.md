# monitorr
[![CI](https://github.com/yumee2/monitorr/actions/workflows/ci.yml/badge.svg)](https://github.com/yumee2/monitorr/actions/workflows/ci.yml) </br>
Инструмент для мониторинга доступности сервисов. Читает список сервисов (URL) и с заданным интервалом проверяет, доступны ли они. Реализован как pet-проект для портфолио.


<img width="900" alt="Screenshot from 2026-08-27 09-29-13" src="https://github.com/user-attachments/assets/f9961988-9d1d-4026-b7e6-f9853fb59c39" />
<img width="900" alt="Screenshot from 2026-08-27 09-31-30" src="https://github.com/user-attachments/assets/a3ee4a88-b8c2-475e-9da4-8c4b33c67257" />



## Идея

- конфиг: список `{name, url, interval}`
- отдельный воркер на каждый сервис делает HTTP-запросы, проверяет статус-код / таймаут
- история проверок сохраняется для расчёта uptime %
- уведомления в Telegram-бот при смене состояния (up -> down, down -> up)
- HTTP API отдаёт статус + uptime % для фронтенда


## Краткое описание

Monitorr — лёгкий инструмент мониторинга доступности сервисов, написанный на Go. Непрерывно отслеживает доступность сервисов, сохраняет историю для расчёта uptime % и отправляет уведомления в Telegram при изменении статуса сервиса.

## Как запустить

1. **Клонировать репозиторий**
   ```bash
   git clone https://github.com/yumee2/monitorr.git
   cd monitorr
   ```

2. **Требования**
   - Go 1.16 или новее
   - Telegram bot token (для уведомлений)

3. **Настройка сервисов**
   - Создать конфигурационный файл со списком сервисов (формат см. выше)
   - Указать данные Telegram-бота в конфиге (bot token и chat ID)
   - При необходимости задать переменные окружения или путь до конфига

4. **Запуск приложенияn**
   ```bash
   go run .\cmd\main.go
   ```

   Или build and run:
   ```bash
   go build -o monitorr.exe ./cmd
   ./monitorr
   ```

5. **Работа с API**
   - HTTP API отдаёт статус сервисов и метрики uptime
   - Фронтенд может обращаться к этим эндпоинтам напрямую

