# F1DB

A distributed database interface based on Mixin Network and IPFS.

## Features

- [x] Write and read record from distributed DAG network.
- [x] Large content supports
- [x] RESTful API supports
- [x] User register and authorization
- [ ] More storage implements
- [ ] Better API authentication
- [ ] Cache
- [ ] Index interface

## Projects using F1DB 

- Whisper.news: https://github.com/fox.one/whisper-news
- Love seal: https://github.com/fox.one/love-seal

## Usage

### Register an App at Open.Fox.ONE

Email open@fox.one for application. Inquire a key/secret pair and a quota. 

### Setup and run an IPFS daemon

FYI: https://docs.ipfs.io/introduction/usage/

### Modify config.yml

- Use config.yml.example as template.
- collector_user_id: a user that use to collector quota.
- generate a new `pin` and keep it safe.

### Register App users

Run 

```bash
$ ./f1db -m register
```

copy the user info and write them into config.yml

### Run http server

```bash
$ ./f1db
```

## License

MIT License
