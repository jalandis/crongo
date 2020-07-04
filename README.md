# crongo
Simple GO cron system

## Purpose

The goal of this project is offer an alternative to the Heroku scheduler addon for GO apps.  

I hope to implement a few improvements over Heroku scheduler or other GO alternatives but
this will largely be a learning experience with some of GO's concurrency features.

## Usage 

Project has an example Heroku setup. [Heroku GO instructions](https://devcenter.heroku.com/articles/getting-started-with-go?singlepage=true) should be suffiecent but remember to scale worker up.

Ex. `heroku ps:scale worker=1 -a crongo`

