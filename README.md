# Go Project

## To-Do:
- API
- CLI
- Scheduled job (either gocron or just cron)

## Development
- Run `go build -o main .`
- Run `./main`

## Deployment

`set _VERSION $(cat VERSION)`

### Locally (or in build pipeline)
- Run `docker build --tag silasbrack/go-project .`
- Run `docker save silasbrack/go-project -o go-project.tar`
- Run `aws s3 cp go-project.tar s3://silas-s3-bucket/container-registry/go-project.tar`
- Run `rm go-project.tar`

### On Digital Ocean Droplet
- Run `aws s3 cp s3://silas-s3-bucket/container-registry/go-project.tar .`
- Run `cat go-project.tar | sudo docker load && rm go-project.tar`
- Run `sudo docker container run -e host="db-postgresql-ams3-31352-do-user-17154797-0.a.db.ondigitalocean.com" -e user="doadmin" -e password="XXX" -e dbname="defaultdb" -e port="25060" -e sslmode="require" -e TimeZone="UTC" silasbrack/go-project:latest`

## Environment Variables
- `PG_HOST` - host of postgres database
- `PG_USER` - username of postgres database
- `PG_PASSWORD` - password of postgres database
- `PG_DBNAME` - database name of postgres database
- `PG_PORT` - port of postgres database
- `PG_SSLMODE` - sslmode of postgres database
- `PG_TIMEZONE` - timezone of postgres database
- `PORT` - port to listen on
- `LOG_DIR` - directory to store logs. If not set, logs will be written to stdout.
