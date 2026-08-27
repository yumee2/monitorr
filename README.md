# monitorr

Uptime checker. Reads a list of services (URLs) and pings each one on an
interval to see if it's up or down. Built as a portfolio project.

## Idea

- config: list of `{name, url, interval}`
- worker per service makes HTTP requests, checks status code / timeout
- persist check history to compute uptime %
- notify via Telegram bot on state change (up -> down, down -> up)
- HTTP API exposing status + uptime %, for a frontend later

## Quick Description

Monitorr is a lightweight uptime monitoring tool written in Go that continuously monitors the availability of your services. It tracks historical data to calculate uptime percentages and sends notifications via Telegram when services go up or down.

## How to Run the App

1. **Clone the repository**
   ```bash
   git clone https://github.com/yumee2/monitorr.git
   cd monitorr
   ```

2. **Prerequisites**
   - Go 1.16 or higher installed
   - Telegram bot token (for notifications)

3. **Configure services**
   - Create a configuration file with your services (see config format above)
   - Set up Telegram bot credentials in the config file (bot token and chat ID)
   - Set environment variables or configuration file path as needed

4. **Run the application**
   ```bash
   go run main.go
   ```

   Or build and run:
   ```bash
   go build -o monitorr
   ./monitorr
   ```

5. **Access the API**
   - The HTTP API will expose service status and uptime metrics
   - Configure your frontend to consume the API endpoints

## Status

- Config and persist check history to compute uptime % are implemented
- Worker goroutines are spawned to check each service

# TODO
- Delete old check history
- Add logger
