
# greload

greload is a reverse-proxy for HTML live-reloading. The program intercepts your http request and injects a script to the HTML response. When local files is modified, the page automatically reloads.

## usage

```
greload [options...] url
```

**options**

    --port              port to listen
    --watch             path to watch
    --exclude           path to exclude from watch list     




