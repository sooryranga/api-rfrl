# tutorme-be

TutorMe Backend

To delete volume for docker db - `docker-compose rm -fv db`

sq : `sq.Select("*").From("auth").Where(sq.And{sq.Eq{"auth.token": token}, sq.Eq{"auth.auth_type": authType}})`
