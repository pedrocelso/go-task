# go-rest-service
[![CircleCI](https://circleci.com/gh/pedrocelso/go-rest-service/tree/master.svg?style=shield)](https://circleci.com/gh/pedrocelso/go-rest-service/tree/master)

An basic REST service (todo app) written in GO.

There's an basic UI for this service at: https://github.com/pedrocelso/go-rest-service-ui

## About the project
I'm working on this project to get to know more about REST services in GO. Specifically, my final goal is to have an fully functional REST service running on Google AppEngine (GAE).

It should support Task creation per User, and each Task can have multiple incidents (or observations).

## Milestones
1. ~~Create an basic MVP saving and retrieving data on MySQL;~~
2. ~~Do whatever is needed to get it up and running on GAE.~~
3. ~~Move the database from MySQL to a noSQL (Google Datastore?)~~
4. Add authentication;
5. Add an new Entity `Task` that will be assigned to a given `User`
6. Add an new Entity `Incident` that will be assigned to a given `Task`

## Trivia
https://cloud.google.com/free/docs/always-free-usage-limits 

### Local Unit Testing for Go
https://cloud.google.com/appengine/docs/standard/go/tools/localunittesting/

### AppEngine Test Package (aetest)
https://cloud.google.com/appengine/docs/standard/go/tools/localunittesting/reference
