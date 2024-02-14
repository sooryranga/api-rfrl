# rfrl-be

rfrl Backend

To delete volume for docker db - `docker-compose rm -fv db`

sq : `sq.Select("*").From("auth").Where(sq.And{sq.Eq{"auth.token": token}, sq.Eq{"auth.auth_type": authType}})`

psql: `psql "postgres://rfrl:secretpassword1@localhost:5432/rfrl?sslmode=disable"`

## Logging

```
 log.Errorj(log.JSON{
  "sql":  sql,
  "args": args,
 })
```

To create migration:
`migrate create -ext sql  -dir ./migration <file_name>`

To migrate up:
`migrate -database "postgres://rfrl:secretpassword1@localhost:5432/rfrl?sslmode=disable" -path ./migration up`

## Create id_rsa

`openssl genrsa -out ./id_rsa 4096`
`openssl rsa -in ./id_rsa -pubout -out ./id_rsa.pub`

## Delete all images

`docker rmi -f $(docker images -a -q)`
