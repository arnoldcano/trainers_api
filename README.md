# trainers_api

This is an api for trainers

## Setup

### Install Docker Desktop
### Start database

`docker-compose start db`

### Create database

`docker-compose run -v $PWD/schema.sql:/schema.sql db psql -h db -U postgres -a -q -f schema.sql`

### Run app

`docker-compose up`

## Endpoints

### GET /trainers/{id}

Get a trainer by id

```
curl http://localhost:8080/trainers/1
```

### POST /trainers

Create a new trainer

```
curl -H "Content-Type: application/json" \
     -d '{"email":"foo@bar.com","phone":"11111111111","first_name":"foo","last_name":"bar"}' \
     http://localhost:8080/trainers
```

## Testing

The project uses standard `go test` for testing. Seed data is loaded from the `testdata/seeds.sql` file. 

### Run tests

`docker-compose run web go test`

## Next Steps

 - [ ] Add more tests
 - [ ] Database migrations
 - [ ] Refactor separation of concerns
 - [ ] Multi-stage Dockerfile for production container
 - [ ] Log aggregation
 - [ ] Status endpoint for healthchecks
 - [ ] Monitoring
 - [ ] Error reporting
