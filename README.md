# go-task
[![CircleCI](https://circleci.com/gh/pedrocelso/go-rest-service/tree/master.svg?style=shield)](https://circleci.com/gh/pedrocelso/go-rest-service/tree/master)

An basic REST service (tasks) written in GO.

~~There's an basic UI for this service at: https://github.com/pedrocelso/go-rest-service-ui~~ A basic UI for this will be made using a create-react-app boilerplate.

## About the project
_My main goal for this project is to learn how to write Go code to be used on AppEngine, using the cloud datastore. I'm trying to add new useful things to the project, and learning them, for example, how to use JWT middleware on gin. You shouldn't expect anything advanced, as I'm struggling to keep this as simple as possible. I hope this helps new go programmers to ger started to AppEngine, and to write basic REST services._

__An generic task abstraction. Users can have unlimited tasks, each task can have unlimited incidents:
An incident can be any kind of interation between the user and its task. It can be observations from a phone call, an email body, even an numeric workflow (not implemented). This should be the first step for complex projects, for instance a Follow-up backend, an BPMS backend, etc.__

## Milestones
1. ~~Create basic routes and handlers;~~
2. ~~Add basic CRUD operations for `Users`;~~
3. ~~Do whatever is needed to make it GAE-compatible;~~
4. ~~Move the database from MySQL to Google Datastore?;~~
4. ~~[Add JWT authentication](https://github.com/pedrocelso/go-rest-service/issues/9);~~
5. ~~Add an new Entity `Task` that will be assigned to a given `User` See [Task Endpoints](https://github.com/pedrocelso/go-rest-service/issues/8)~~
6. Add an new Entity `Incident` that will be assigned to a given `Task` See [Incidents](https://github.com/pedrocelso/go-rest-service/issues/120)

## Trivia
https://cloud.google.com/free/docs/always-free-usage-limits 

### Local Unit Testing for Go
https://cloud.google.com/appengine/docs/standard/go/tools/localunittesting/

### AppEngine Test Package (aetest)
https://cloud.google.com/appengine/docs/standard/go/tools/localunittesting/reference
