![sibylwithtextgrad2](https://user-images.githubusercontent.com/9204400/49330993-96c7a080-f564-11e8-9285-069ebb91d06c.png)

Sibyl is a **Display Only** stock trading platform. For most discount brokers that provide an API, any access to that API would be required to be Display Only program.  Sibyl was created to provide a platform that would satisfy that requirement and provide a UI for trading regardless of your discount broker.

# Table of Contents

- [Installing](#installing)
- [Running](#running)
  - [SibylServer](#sibylserver)
    * [MySQL](#mysql)
  - [SibylCli](#sibylcli)
    * [Adding](#adding)
    * [Showing](#showing)
    * [Enabling](#enabling)
    * [Deleting](#deleting)
    * [Setting](#setting)
  - [Web GUI](#web-gui)

# Installing
Using Sibyl is easy. First, use `go get` to install the latest version.  Note Sibyl requires Golang version 1.11+.

    go get -u github.com/nathanhack/sibyl/cmd/...

# Running
After installing Sibyl there are two programs that are installed: SibylServer and SibylCli.

## SibylServer
The server requires Mysql server be running (some where - defaults to localhost) with the user:'sibly' and password:'pa$$word'.
By running the following at the commandline will run SibylServer using all defaults.

    SibylServer http

### MySQL
As example the following could be executed from commandline to add the user.

    CREATE USER 'sibyl'@'localhost' IDENTIFIED BY 'pa$$word';
    GRANT CREATE ON *.* TO `sibyl`@'localhost';
    GRANT DELETE ON *.* TO `sibyl`@'localhost';
    GRANT DROP ON *.* TO `sibyl`@'localhost';
    GRANT INDEX ON *.* TO `sibyl`@'localhost';
    GRANT INSERT ON *.* TO `sibyl`@'localhost';
    GRANT SELECT ON *.* TO `sibyl`@'localhost';
    GRANT UPDATE ON *.* TO `sibyl`@'localhost';
    FLUSH PRIVILEGES;

Additionally, this is a system value that must be set:

    SET GLOBAL local_infile = 'ON';

## SibylCli
To interact with the SibylServer there is a cli tool called SibylCli. The tool supports adding and deleting stocks, showing stocks, history and intraday values, and setting up creds for the discount broker.  Once a stock is added it will validate the symbol was a valid stock and depending on which actions are enabled it will being downloading the appropriate information.

### Adding
To add a stock for the SibylServer track you execute the following command. Note you must [enable](#enabling) downloading for the stock, because any new stocks added default to disabled for all actions.

    SibylCli add STOCK

### Showing
Information is all stored in MySQL database and is accessed via SibylServer.
To show a list of stocks currently being tracked by SibylServer:

    SibylCli show stocks

To show the history for a particular stock:

    SibylCli show history STOCK

### Enabling
SibylServer will take perform several actions for all the stocks that has been added to it.  The actions are determined by internal state to the server for each particular stock.  There are lots of options but the easiest one to get started is to run the following command after your done adding stocks (or after adding new ones later).

    SibylCli enable all

### Deleting
To remove a stock from SibylServer use the following command. Note that it will delete all data associated with that stock so use carefully.

    SibylCli delete STOCK

### Setting
To set the creds for the discount broker use the following command:

    SibylCli set ally key secret token tokenSecret
    
Note currently only Ally is supported.

## Web GUI
To be added in a future release.
