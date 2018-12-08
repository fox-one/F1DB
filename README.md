# F1DB

A distributed database interface based on Mixin Network and IPFS.

## Features

- [x] Write and read record from distributed DAG network.
- [x] Large content supports
- [x] RESTful API supports
- [x] User register and authorization
- [x] Better API authentication
- [x] JSON RPC service
- [ ] More storage implements
- [ ] Cache
- [ ] Storage node communication
- [ ] Index interface

## Documents

- JSON RPC Document: https://documenter.getpostman.com/view/140588/RzfiG8ZN#intro

## Projects using F1DB 

- Whisper.news: https://github.com/fox.one/whisper-news
- Love seal: https://github.com/lyricat/love-seal

## Usage

### 1 Register an App at Open.Fox.ONE

Email `open@fox.one` for application. Inquire a key/secret pair and a quota. 

### 2 Setup and run an IPFS daemon

FYI: https://docs.ipfs.io/introduction/usage/

### 3 Modify config.yml

- Use config.yml.example as template.
- collector_user_id: a user that use to collector quota.
- generate a new `pin` and keep it safe.

### 4 Run as JSON RPC service or use in your code

**JSON RPC service**

```bash
$ ./f1db
```
**use F1DB in your code**

check out the [source](https://github.com/fox-one/F1DB/tree/master/controller).

## License

MIT License
