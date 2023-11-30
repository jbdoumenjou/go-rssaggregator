Simple rss aggregator that follows the [boot.dev](https://boot.dev) course assignments.

The reference video is [there](https://www.youtube.com/watch?v=dpXhDzgUSe4&t=1s).


> We're going to build an RSS feed aggregator in Go! It's a web server that allows clients to:
>* Add [RSS](https://en.wikipedia.org/wiki/RSS) feeds to be collected
>* Follow and unfollow RSS feeds that other users have added
>* Fetch all of the latest posts from the RSS feeds they follow


# Bootstrap

## Web Server
We will use the following stack:

* [chi](https://github.com/go-chi/chi)
* [cors](https://github.com/go-chi/cors)
* [godotenv](https://github.com/joho/godotenv)

Like for the [web server project](https://github.com/jbdoumenjou/mygoserver),
we will use a .env file to store the configuration.
Don't forget to add the .env file to your .gitignore file.

Let's start with something like this:

```bash
PORT="8080"
```

## Database

We will use [PostgreSQL](https://www.postgresql.org/) as a database.

>We'll be using a couple of tools to help us out:
>* [database/sql](https://pkg.go.dev/database/sql): This is part of Go's standard library.
>It provides a way to connect to a SQL database, execute queries, and scan the results into Go types.
>* [sqlc](https://sqlc.dev/): SQLC is an amazing Go program that generates Go code from SQL queries.
>It's not exactly an ORM, but rather a tool that makes working with raw SQL almost as easy as using an ORM.
>* [Goose](https://github.com/pressly/goose): Goose is a database migration tool written in Go.
>It runs migrations from the same SQL files that SQLC uses, making the pair of tools a perfect fit.

We will generate go code from SQL query and use a migration tool to manage the database schema.

The course proposes to install postgres and pgadmin locally, I prefer to use docker images.
I used to use [dbeaver](https://dbeaver.io/) as DB client, but I will try to use pgadmin this time.

The [Getting Started with PostgreSQL](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html) of sqlc is a good start.

I configure sqlc to generate a querier interface to easily mock the database in the tests.
