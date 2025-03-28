# Alum Discord Bot

## Build Docker Image
```bash
$ docker build -t discord-bot .

```

## Build go binary
```bash
$ go build -o bin/main cmd/alum-bot/main.go

# Run the binary
$ ./bin/main -t <DISCORD_BOT_TOKEN>
```