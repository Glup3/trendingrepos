# GitHub Data Loader

## Generating GraphQL Go Code

```sh
go run github.com/Khan/genqlient
```

## 100 Concurrent Requests (Secondary rate limits)

The [GitHub GraphQL API docs](https://docs.github.com/en/graphql/overview/rate-limits-and-node-limits-for-the-graphql-api#secondary-rate-limits)
state that it's possible to make at maximum 100 concurrent requests. After some try
and error I figured out that it's only possible if you actually fire 100 concurrent
requests AT ONCE! My previous approach of staggered requests failed after around
50 requests (I did 9 parallel, 1 blocking --> and repeat). The blocking one was
necessary in order to minimize the star range (--> repository counts).

The workaround (solution) is to aggregate the star ranges beforehand and then
firing 100 requests at once.
