## TgPoster - CLI утилита для постинга в Telegram канал

## Возможности:

### Быстрый постинг
```sh 
tg-poster post --file example-post/direct.md --channel @channel_name
```
### Отложенный постинг
```sh
tg-poster schedule -f example-post/deferred.md -c @channel_name -T "2026-02-18 14:00"
```
### API для постинга
```sh
tg-poster serve --port 8080
