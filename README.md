# Spotify

![screenshot](./img/image.png)

## Project work for the last assignment of Alif Academy Course.

Spotify - is backend part of the music streaming service application written on Go (Golang).

## Features

- Authentication. Users can register and sign in.
- Protected endpoints. Only signed-in users can create actions.
- RESTful routing.
- Middleware.
- PostgreSQL database.


#
First thing that you should do, is install [Golang](https://golang.org/dl/) and any code editor. For example i use [VSCode](https://code.visualstudio.com/) or you can install Goland IDE.

If you have a Windows installed, need to download [docker tool-box](https://github.com/docker/toolbox/releases) for Win7 and [Docker desktop](https://www.docker.com/products/docker-desktop) for Win10 and for Mac [this](https://hub.docker.com/editions/community/docker-ce-desktop-mac).
If you have a Linux, please use this [guide](https://hub.docker.com/search?offering=community&operating_system=linux&q=&type=edition) for installation.

Then, if you have completed any of this points above, we will begin.

##
# Installation

1. Create a directory and clone this repository: 

```sh
$ git clone https://github.com/SardorMS/project.git .
```

2. If you have VScode, you will need to install additional extensions for it:
- RestClient, SQLTools, SQLTools Postgre SQL and Docker.

3. Then you will need to customize the dockep-compose.yml file by specifying your desired values in enviromets  settings and the same values will need to set in main.go file:
```go
//userlogin@host:port/db
dsn := "postgres://app:123@localhost:5432/db"
```

4. And now you configure SQL Tools extantion to run PostgreSQL image:
- Add new connection (use PostgreSQL);
- Set name of connection;
- Set Server Addres - on Win7 it's a 192.168.99.100. On Win10, Linux or Mac you can use localhost address;
- Set database name, username and password and check test connection below. If it's ok you can create and save connection.
All setting will be avialable in the settings.json file.


5. After the connection you can start docker-compose file to run PostgreSQL Database:

```sh
$ sudo docker-compose up
```

6. Run main.go file, pre-setting the address and port to start server itself. Use localhost or 127.0.0.1 addres with 9999 port.
```go
host := "localhost"
port := "9999"
```

```sh
$ go run cmd/main.go
```


## References

To get started, please read the [quick start](https://developer.spotify.com/documentation/web-api/).

I'll hope that you were able to configure and run everything, and also read the manual.

Let's start :)

## Usage

1. Next step is a tables:
- schema.sql  - provides necessary table for further work. Of course you can make your own changes. To add tables, you need to select each table and press Ctrl+E+E. After that, all tables will appear in the SQLTools extension itself.

2. List of all methods:
- server.go - file with lists of all possible endpoint methods;

3. And here is the RestClient extansion files:
- test.http - file with examples of using these methods for application, users and playlists. 
- test2.http - file with examples of methods for tracks.

4. Before using the methods, we need to generate some variables with random length. For this you need to run main.go in the pkce directory:

```sh
$ go run cmd/pkce/main.go
```
And we got a pkce.json file, which has all data we need. Then replace that data into the first method in the test.http file and make the request. You should receive the answer in the form of json.

Then follow the instructions below and perform the methods one by one.

PDF file (on russian) for reading will be available at [here](https://drive.google.com/file/d/1K3MT1Mt5JOnjv9CqdRJ2gASAgn7AuoQl/view?usp=sharing).

English version wil be soon :)
