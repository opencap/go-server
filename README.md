# go-opencap

Basic OpenCAP server written in Go. Supports "Address Query" and "Alias Management"

## Download source code

```bash
go get github.com/opencap/go-server
```

## Quickstart

For this guide I will assume you will be running your server using Ubuntu 16.04 LTS.

### Download or build the "go-server" program and place its own folder

### Add a file to the same folder called ".env"

.env must contain the following variables:

```bash
PORT=443
DB_URL=database.db
DB_TYPE=sqlite3
JWT_EXPIRATION_MINUTES=30
JWT_SECRET=FSHDFKSDHFSDLFABSDJFLSDFSDFJSLDFBASDFBDSJFSHDFASDKFJSLDKFALSDFKJADFKASDFD
PLATFORM_ENV=prod
TIMEOUT_SECONDS=30
CREATE_USER_PASSWORD=somepassword
```

The above configuration will work if you are planning on using sqlite as a database. However be sure to change CREATE_USER_PASSWORD and JWT_SECRET for security puposes. PORT must be set to 443 for running in producion.

### Install sqlite3

```bash
sudo apt-get update.
sudo apt-get install sqlite3 libsqlite3-dev
```

### HTTPS/TLS

Upon startup the server will check for two files, "cert.pem" and "key.pem" in the same directory as the executable. If they are not found it will generate new certificates so that it can serve using HTTPS/TLS.

### Port Forwarding

If you are running the server on a computer that is connected to the internet via a home router, the router must forward the PORT to the computer. On startup the server tries to do this automatically using uPnP (Universal plug n play). If it is unable to do so it will print an error message. Make sure that your router is upgraded to the latest version of its firmware, that the go-server application has proper permissions from your firewall, and that uPnP is enabled on your router.

Port forwarding can also be done manually if that is easier for you.

### Run the server

Make sure you are in the directory containing the executable and you can start the server using ./go-server

### Add DNS Records

You will need to setup your server with a domain name that you own. There are many tutorials online of how to do this step, but basically you log into your DNS providers website and add an A record to your domain name that points to your computer's IP address.

You will also need to add a SRV record that points to the domain of your server. The SRV record should have the following information:

```bash
_opencap._tcp.DOMAIN_NAME_HERE. 86400 IN SRV 5 12 443 DOMAIN_NAME_HERE.
```

For example:

```bash
_opencap._tcp.mysite.com. 86400 IN SRV 5 12 443 mysite.com.
```

If your server is running at a different location that the "domain" section of the aliases you are hosting then it would look like this:

```bash
_opencap._tcp.aliasdomain.com. 86400 IN SRV 5 12 443 serverdomain.com.
```

## Using the server

Once it is up and running, you can create users by making requests of the following format:

```json
POST https://myserver.com/v1/users
content-type: application/json

{
    "alias": "username$myserver.com",
    "password": "myNewUserPassword",
    "create_user_password": "somepassword" // same password used in .env file
}
```

All other requests follow the OpenCAP protocol. Examples can be seen in the test.rest file

## Additional database options

If you want to use something other tha sqlite some other valid DB_TYPE, and DB_URL variables are:

- "sqlite3" "/tmp/gorm.db"
- "mssql" "sqlserver://username:password@host:1433?database=dbname"
- "mysql" "user:password@/dbname?charset=utf8&parseTime=True&loc=Local"
- "postgres" "postgres://username:password@host:5432/dbname?sslmode=disable"

## Testing

docker-compose is used for testing:

```bash
docker-compose build go-server
docker-compose run go-server
cd src/github.com/opencap/go-server
go test ./...
```

## Contribute

Feel free to make pull requests and open issues!
