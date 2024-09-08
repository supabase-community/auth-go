# Testing auth-go

## Postgres set up
The postgres server requires some set up beyond the default postgres docker image. For example, it needs to set up the `supabase_auth_admin` user and the `auth` namespace. This is done using `testing/init_postgres.sh`, which is copied into the postgres image's entrypoint directory so that it runs on startup.

This script in run by the `postgres` user, as configured for the postgres image in `docker-compose`. This is why the docker image is configured with `postgres:root`, but the Auth server connects to postgres using the newly created `supabase_auth_admin` user, which also uses `root` as its password.

This is based on the same set up in the [supabase/auth](https://github.com/supabase/auth) repo.

### Start up

As postgres can take some time to startup and run the entrypoint on it's first run, the Auth server will often fail to connect to postgres on the first attempt. This is normal, and it should be restarted by docker compose and eventually succeed in connecting and running migrations.

Our tests can also sometimes fail by running too early. `TestMain(m *testing.M)` in `client_test.go` should ensure the Auth server is ready for requests before running any tests.

### Volume

To ensure a clean start, the postgres container does not use a volume, meaning all of its data is lost every time the container stops.

## Multiple servers

We actually spin up 3 Auth instances, each with different configurations, to facilitate testing the client where the server config affects what gets returned. For example, if autoconfirm (of the user's email/phone) is disabled, signup returns a `types.User`. If autoconfirm is enabled, a `types.Session` is returned instead. Autoconfirm is a configuration on the server, so there's no easy way to test this without running multiple servers.

They all seem to be happy enough interacting with the same database, though.

## Rate limiting

Auth implements rate limits. We can configure high limits for HTTP requests, but in the case of signup requests, for example, there are further tests that a confirmation email wasn't sent too recently (typically within the last 60 seconds). Instead, we set the email rate limit to 1 request per nanosecond.
