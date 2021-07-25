# rfrl-be

rfrl Backend

To delete volume for docker db - `docker-compose rm -fv db`

sq : `sq.Select("*").From("auth").Where(sq.And{sq.Eq{"auth.token": token}, sq.Eq{"auth.auth_type": authType}})`

psql: `psql "postgres://rfrl:secretpassword1@localhost:5432/rfrl?sslmode=disable"`

## Env

`ASSETS_FOLDER` : outline where the assets folder is.

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

## Dockerfile

<https://levelup.gitconnected.com/complete-guide-to-create-docker-container-for-your-golang-application-80f3fb59a15e>

## Production

`psql "sslmode=verify-ca sslrootcert=server-ca.pem sslcert=client-cert.pem sslkey=client-key.pem hostaddr=34.69.205.182 port=5432 user=rfrl dbname=rfrl"`
`psql "postgres://rfrl:__PASSWORD__@34.69.205.182:5432/rfrl?sslcert=client-cert.pem&sslkey=client-key.pem&sslrootcert=server-ca.pem&sslmode=verify-ca"`

`migrate -database "postgres://rfrl:__PASSWORD__@34.69.205.182:5432/rfrl?sslcert=client-cert.pem&sslkey=client-key.pem&sslrootcert=server-ca.pem&sslmode=verify-ca" -path ./migration up`
