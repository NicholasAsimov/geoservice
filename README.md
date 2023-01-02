# geoservice

## Description
geoservice is a microservice for ingesting and serving geolocation data from a CSV file.

CSV files can be ingested either by using an API endpoint `/api/v1/ingest` or
by using a CLI utility `importcsv` that imports a local file.


## Running

### Dependencies and API server

`docker-compose up` will start a postgres DB container (mapped to local port 5432) and the API
server (mapped to local port 5000).

### CSV importer
`go run ./cmd/importcsv` will start the import process (by default using csv file `data_dump.csv` in the current directory).

### Fake data generator
`go run ./cmd/generator -file fakedata.csv -records 10000000` will generate a fake CSV file with 10,000,000 records.

Generator writes CSV records directly to disk in a streaming manner so it doesn't allocate any memory
to hold the records, thus allowing to generate arbitrary file sizes independent of available RAM.


## Example usage
```
> docker-compose up -d
[+] Running 4/4
 ⠿ Network geoservice_default       Created
 ⠿ Volume "geoservice_data"         Created
 ⠿ Container geoservice-postgres-1  Started
 ⠿ Container geoservice-api-1       Started

> go run ./cmd/generator -file data_dump.csv -records 1000000
INF cmd/generator/main.go:76 > csv generated file=data_dump.csv records=1000000 took=2.468177767s

> go run ./cmd/importcsv
INF cmd/importcsv/main.go:62 > parsing file file=./data_dump.csv
INF cmd/importcsv/main.go:73 > persisting to db records=999892
INF cmd/importcsv/main.go:90 > import finished accepted=999892 db_took=27.75483288s parse_took=2.001588949s records=1000000 skipped=108

> curl -X POST -F file=@data_dump.csv -D- 'http://localhost:5000/api/v1/ingest'
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 49

{"records":1000000,"accepted":999892,"skipped":108}
```


## Notable things
1. csv parsing logic is separated from validation (business) logic by passing in "validationFunc" closure

2. csv parsing is streaming and doesn't allocate extra memory by doing all the work in a single loop

3. migrations are done properly and are embedded into the binary using Go 1.16 embed library so it's
   still a single binary distribution

4. very fast insertion into the database by using postgres' native COPY FROM command. The rows are
   being quickly inserted into a temporary table and then the upsert is performed from the temporary
   table to the actual one, of course in a single transaction (see store/store.go: UpsertRecords
   function).
   This database insertion also doesn't allocate additional memory, as it would be in case of using
   GORM v2, for example, because it would prepare a huge INSERT/UPDATE statement for all 1M rows.
   This is achieved by implementing pgx CopyFromSource interface for GeoRecord type to transform
   application type to db type on the fly (see store/pgxutil.go: CopyFromRecords function).

5. configuration is done via env variables and all the default values are correctly set for local
   development, so there's no requirement to specify any configuration, only as needed. You can
   check all available configuration parameters by checking the help:
```
> go run ./cmd/api -h

This application is configured via the environment. The following environment
variables can be used:

KEY                            TYPE             DEFAULT            REQUIRED    DESCRIPTION
GEOSERVICE_SERVER_ADDR          String           localhost
GEOSERVICE_SERVER_PORT          String           5000
GEOSERVICE_DB_HOST              String           localhost
GEOSERVICE_DB_PORT              Integer          5432
GEOSERVICE_DB_NAME              String           geoservice
GEOSERVICE_DB_USER              String           geoservice
GEOSERVICE_DB_PASSWORD          String           geoservice
GEOSERVICE_DB_SSL               String           disable
GEOSERVICE_IMPORTER_FILEPATH    String           ./data_dump.csv
GEOSERVICE_LOGLEVEL             String           debug
GEOSERVICE_PRETTYLOG            True or False    true
```

## File structure
```
├── cmd
│   ├── api
│   │   └── main.go
│   ├── generator
│   │   └── main.go
│   └── importcsv
│       └── main.go
├── config
│   └── config.go
├── csvparse
│   ├── csvparse.go
│   └── csvparse_test.go
├── model
│   └── model.go
├── resources
│   ├── sql
│   │   ├── 000001_create_georecords_table.down.sql
│   │   └── 000001_create_georecords_table.up.sql
│   └── embed.go
├── server
│   ├── routes.go
│   └── server.go
├── store
│   ├── migrate.go
│   ├── pgxutil.go
│   ├── store.go
│   └── store_test.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```
