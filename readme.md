Simple rss aggregator that follows the [boot.dev](https://boot.dev) course assignments.

The reference video is [there](https://www.youtube.com/watch?v=dpXhDzgUSe4&t=1s).


> We're going to build an RSS feed aggregator in Go! It's a web server that allows clients to:
>* Add [RSS](https://en.wikipedia.org/wiki/RSS) feeds to be collected
>* Follow and unfollow RSS feeds that other users have added
>* Fetch all the latest posts from the RSS feeds they follow

# Usage

This projet launches a web server that exposes a REST API,
and store data in a PostgreSQL database.

A makefile is provided to build and run the project.

```shell
$make help
help:          Show this help.
build:         clean Build the application.
clean:         Clean the application.
run:           Run the application.
test:          Test the application.
migrate-up:    Apply all up migrations.
migrate-down:  Apply all down migrations.
sqlc:          Generate the database code.
mock:          Generate a store mock.
```

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
I use mock to generate the [mock](https://github.com/uber-go/mock) of the querier interface.
It is an exercise to use the mock package and to understand how it works.
We could launch a db in a container for the tests and avoid to mock the database.

## Tests

I integrate [testcontainers](https://testcontainers.com/) to launch a postgresql container for the tests.
It will be a good exercise to use testcontainers and to understand how it works.
For the installation, I followed the [quickstart](https://golang.testcontainers.org/quickstart/) of testcontainers.
Then I used [postgres container](https://golang.testcontainers.org/modules/postgres/) to launch a postgres container for the tests.
I add a containers file to centralize the containers configuration.
In the same way, I add a migrations file to centralize the goose migrations configuration.
Then I use both containers and migrations in the tests.
To launch only one container for all the tests,
I use the test main function to launch the container and to close it at the end of the tests.