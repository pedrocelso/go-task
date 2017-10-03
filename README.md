# go-rest-service
(Working on a better name)
[![CircleCI](https://circleci.com/gh/pedrocelso/go-rest-service/tree/master.svg?style=shield)](https://circleci.com/gh/pedrocelso/go-rest-service/tree/master)

An basic REST service (tasks) written in GO.

There's an basic UI for this service at: https://github.com/pedrocelso/go-rest-service-ui

## About the project
An generic task abstraction. Users can have unlimited tasks, each task can have unlimited incidents.

An incident can be any kind of interation between the user and its task. It can be observations from a phone call, an email body, even an numeric workflow (not implemented).

## Milestones
1. ~~Create basic routes and handlers;~~
2. ~~Add basic CRUD operations for `Users`;~~
3. ~~Do whatever is needed to make it GAE-compatible;~~
4. ~~Move the database from MySQL to Google Datastore?;~~
4. [Add JWT authentication](https://github.com/pedrocelso/go-rest-service/issues/9);
5. Add an new Entity `Task` that will be assigned to a given `User` See [Task Endpoints](https://github.com/pedrocelso/go-rest-service/issues/8)
6. Add an new Entity `Incident` that will be assigned to a given `Task`

## Trivia
https://cloud.google.com/free/docs/always-free-usage-limits 

### Local Unit Testing for Go
https://cloud.google.com/appengine/docs/standard/go/tools/localunittesting/

### AppEngine Test Package (aetest)
https://cloud.google.com/appengine/docs/standard/go/tools/localunittesting/reference
