# Pantry

_A simple pantry app using Golang, HTTP, and PostgreSQL._

## Requirements

* Golang: `1.18`
* Go-Bindata
* PostgreSQL Server (may be remote)

## Quickstart

1. `chmod +x run.sh build.sh update_assets.sh`
    * This will allow your system to run the shell scripts
    * All are required as the previous relies on the next
2. `psql -h {host} -U {username} -d {databaseName} -a -f pantry.sql`
    * This will ensure you have the right table in your database
3. `./run.sh`
    * This will build the assets directory
      * You may include a `.env` file if you wish as `/assets/.env`
      * You may also include environment variables in the standard way
    * This will build the binary for your system
    * This will run the go binary built

## Environment Variables

* Database Variables
  * `DB_HOST` (localhost): The location of the Postgres server
  * `DB_PORT` (5432): The port on which the server runs
  * `DB_USER` (postgres): The service account user which will execute the queries
  * `DB_PASS` (postgres): The password of the service account user which will execute the queries
  * `DB_NAME` (candahome): The name of the database to be used by this system
* Server Variables
  * `SERVER_PORT` (8080): The port on which the HTTP server of this app runs
