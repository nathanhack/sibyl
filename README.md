![sibylwithtextgrad2](https://user-images.githubusercontent.com/9204400/49330993-96c7a080-f564-11e8-9285-069ebb91d06c.png)

Sibyl is WIP stock trading platform.  Sibyl is setup to use a couple of data sources namely: Alpaca and Polygon.io.  For executing trades it Alpaca.  

# Table of Contents

- [Installing](#installing)
- [Running](#running)
  - [SibylServer](#server)
    * [Sqlite3](#sqlite3)
    * [MySQL](#mysql)
  - [SibylCli](#cli)
    * [Adding](#adding)
    * [Showing](#showing)
  - [Roadmap](#Roadmap)

# Running
There is a server and a client. The server run in the background and downloads data and stores it in a database (either Sqlite3 or MySql configured via `sibyl.json` config).

## server
For help and usage:
`go run cmd/server/main.go -h`

For the quick and easy just run:
`go run cmd/server/main.go`

This will create the `sibyl.json` config if it doesn't exists. Then you can update it for you configuration. Then rerun to start the server. 


### Sqlite3
Sibyl can use Sqlite3 just provide it configuration info in the sibl.json config. ex:
```
"databaseold": {
    "dialect": "sqlite3",
    "dsn": "./db.sqlite?_fk=1"
  },
```

### MySQL
Sibyl expects to use a database calls `sibyl` so for MySql it must already exists.
To connect to a MySql database add the following to the configuration sibyl.json updated for your configuration:
```
 "database": {
    "dialect": "mysql",
    "dsn": "root:my-secret-pw@tcp(localhost:3306)/sibyl?parseTime=True"
  },
```
The above sibyl.json configuration works with the example local insecure MySql docker example below.

The following is a fast way to setup a insecure docker MySql for testing:
```
docker run --name mysql -p 3306:3306 -v <DIRECTORY_FOR_MYSQL>:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:latest
```

If your running MySql but would like Sibyl to create the sibyl database for you. Just include `"root"` in the database section of the sibyl.json ex:
```
 "database": {
    "dialect": "mysql",
    "dsn": "root:my-secret-pw@tcp(localhost:3307)/sibyl?parseTime=True"
    "root": "root:my-secret-pw@tcp(localhost:3307)/?parseTime=True"
  },
```


## cli

To interact with the server there is a cli tool. The tool supports adding and deleting stocks, showing stocks, history and intraday values.  Once a stock is added it will validate the symbol was a valid stock and depending on which actions are enabled it will being downloading the appropriate information.

### Adding
To add a stock for the SibylServer track you execute the following command. Note you must [enable](#enabling) downloading for the stock, because any new stocks added default to disabled for all actions.

    `go run cmd/cli/main.go add STOCK`

### Showing
Information is all stored in the database and is accessed via the server.
To show a list of stocks currently being tracked:

```
go run cmd/cli/main.go show stocks
```

or
    
```
go run cmd/cli/main.go show stocks --details
```

To show the history for a particular stock:

```
go run cmd/cli/main.go show bars daily STOCK
```

## Roadmap
client
    
- [ ] Add a Load command (reverse of dump)

server

- [ ] Add a web GUI
- [ ] Get stock financials from Polygon.io
- [X] Add support trades history
- [ ] Add support for realtime data (websockets)
- [ ] Add external web api to trade through Alpaca
