# Data Loader

## Concurrent Requests (Secondary rate limits)

The [GitHub GraphQL API docs](https://docs.github.com/en/graphql/overview/rate-limits-and-node-limits-for-the-graphql-api#secondary-rate-limits)
state that it's possible to make at maximum 100 concurrent requests. After some try
and error I figured out that it's only possible if you actually fire 100 concurrent
requests AT ONCE!

The workaround (solution) is to aggregate the star ranges beforehand and then
firing 100 requests at once.

UPDATE: I don't know why but I can consistently make 200 requests at once with
a 90 seconds cooldown. I hope they don't patch it 😂

## Development

Load environment variables. I am using [direnv](https://direnv.net), take a look
at `.envrc.example`.

```sh
direnv allow .
```

Run TimescaleDB locally with docker.

```sh
docker run -d --name timescaledb -p 5432:5432 -e POSTGRES_PASSWORD=password timescale/timescaledb-ha:pg16
```

Init Go modules.

```sh
go mod tidy
```

Apply database migrations using [goose](https://github.com/pressly/goose).

```sh
goose status
goose up
```

Generate GraphQL code.

```sh
go run github.com/Khan/genqlient
```

Generate [sqlc](https://docs.sqlc.dev/en/stable/overview/install.html) code.

```sh
sqlc generate
```

Start the data loader.

```sh
go run cmd/loader/main.go
```
