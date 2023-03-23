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

# Installing
(To be updated)
# Running
(To be updated)

## server
(To be updated)

### Sqlite3
(To be updated)

### MySQL
(To be updated)


## cli
(To be updated)

To interact with the server there is a cli tool. The tool supports adding and deleting stocks, showing stocks, history and intraday values.  Once a stock is added it will validate the symbol was a valid stock and depending on which actions are enabled it will being downloading the appropriate information.

### Adding
(To be updated)

To add a stock for the SibylServer track you execute the following command. Note you must [enable](#enabling) downloading for the stock, because any new stocks added default to disabled for all actions.

    cli add STOCK

### Showing
(To be updated)

Information is all stored in the database and is accessed via the server.
To show a list of stocks currently being tracked:

    cli show stocks
or
    
    cli show stocks --details

To show the history for a particular stock:

    cli show bars daily STOCK

## Roadmap
* Add a web GUI
* cli
    * Add a Load command (reverse of dump)
* server
    * Get stock financials from Polygon.io
    * Add support trades history
    * Add support for realtime data (websockets)
    * Add external web api to trade through Alpaca
