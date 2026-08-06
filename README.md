# monitorr

Uptime checker. Reads a list of services (URLs) and pings each one on an
interval to see if it's up or down. Built as a portfolio project.

## Idea

- config: list of `{name, url, interval}`
- worker per service makes HTTP requests, checks status code / timeout
- persist check history to compute uptime %
- notify via Telegram bot on state change (up -> down, down -> up)
- HTTP API exposing status + uptime %, for a frontend later

## Status

- Config yaml reading is implemented
- Worker goroutines are spawned to check each service

# TODO
- Correct context cancellation
- Add waitgroups to wait for all workers to finish
- Set and use checkedAt field
- Write simple tests
