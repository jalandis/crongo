# crongo
Simple GO cron system

## Purpose

The goal of this project is offer a simple alternative to the Heroku scheduler addon for GO apps.  

Minimal features but a stable example that is flexible enough for most use cases.

## Usage 

The project has an example Heroku setup in `cmd/cron` that illustrates usage. 

For setup of heroku apps, the [Heroku GO instructions](https://devcenter.heroku.com/articles/getting-started-with-go?singlepage=true) should be sufficient but remember to scale worker up.

`heroku ps:scale worker=1 -a crongo`
