[chat](https://www.notion.so/chat-519ef4e7f8f747899ec33825d8368a86)

## Go言語勉強のためにnet/httpパッケージを使った簡単なチャットアプリ。

### Install cmd

```
go get -u github.com/tom1236908745/goChat
```

### You should install gorilla/websocket, if you haven't installed yet. (https://github.com/gorilla/websocket)

```bash
go get github.com/gorilla/websocket
```

### You have to check GOPATH, and then type command bellow.

### install

```
go get github.com/tom1236908745/goChat
```
### build

```bash
go build -o chat # fileName can be named whatever you like!
```

### run

```bash
./chat -addr=":portNumber"  # portNumber can be named whatever you like! ;)
```

(or 

```bash
./chat # default port number :8080)
```

Start up localhost (default [http://localhost:8080](http://localhost:8080) ) to view it in the browser.
