# tutorme-be

TutorMe Backend

To delete volume for docker db - `docker-compose rm -fv db`

sq : `sq.Select("*").From("auth").Where(sq.And{sq.Eq{"auth.token": token}, sq.Eq{"auth.auth_type": authType}})`

psql: `psql "postgres://tutorme:secretpassword1@localhost:5432/tutorme?sslmode=disable"`

Logging :

```
 log.Errorj(log.JSON{
  "sql":  sql,
  "args": args,
 })
```
