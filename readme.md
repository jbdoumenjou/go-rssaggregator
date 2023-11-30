Simple rss aggregator that follows the [boot.dev](https://boot.dev) course assignments.

The reference video is [there](https://www.youtube.com/watch?v=dpXhDzgUSe4&t=1s).


> We're going to build an RSS feed aggregator in Go! It's a web server that allows clients to:
>* Add [RSS](https://en.wikipedia.org/wiki/RSS) feeds to be collected
>* Follow and unfollow RSS feeds that other users have added
>* Fetch all of the latest posts from the RSS feeds they follow


# Bootstrap

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
