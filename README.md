## gobank

JSON API app with user create, deposit to user, transfer between users.

### Run

Environment variables load from .env file

Up Postgres

```
docker-compose up -d
```

Create tables

```
make migrate
```

Start application

```
make run
```

### Tests

```
make test
```
