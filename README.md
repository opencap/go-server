# go-opencap

Basic OpenCAP server written in Go. Supports "Address Query" and "Alias Management"

## Full Installation Guide

This guide will show you how to setup and run an OpenCAP server that you can host yourself. It requires that you have a computer that is always on and connected to the internet, and that you own a domain name. E.g. "mywebsite.<span></span>com". You should be aware that you will be running a full webserver on your home internet connection and bandwidth restrictions may apply. The developers of this open-source software are in no way responsible for any security issues you encounter.

### Download

Download the latest folder for your operating system from the releases tab above and unzip it.

### Terminal Disclaimer

In some of the steps below you will be required to use the command terminal. The way to open a command terminal in the proper location depends on your operating system.

#### Windows 10

Use the file explorer to navigate into the folder you just downloaded. In the address bar, type "cmd" and hit enter. This should open a terminal for you.

#### Linux

If you're using linux you probably know how to use a terminal.

### Find your IP address and setup your domain records

You can lookup your computer's IP address by googling "My IP address" or by using the tool built into the program:

#### Linux

```bash
./go-server --getip
```

#### Windows Command Line

```cmd
go-server.exe --getip
```

### Add an A record to point to your IP address

You will need to setup your server with a domain name that you own. There are many tutorials online of how to do this step, but basically you log into your DNS provider's website and add an A record for your domain name that points to your computer's IP address.

The domain name should NOT include "www". It should be the equivalent of "example.com".

Verify that the record is live using google's tool: https://dns.google.com/

Type in your domain name to the search bar. The result should contain the IP address of your computer. If it doesn't, double check your A record. Also try waiting a few minutes, it usually takes a bit for the DNS record to propogate. While you are waiting you can setup your SRV record.

### Add a SRV record to point to your webserver

You will also need to add a SRV record that points to the domain of your server. This is done in the same place where you updated the A record. The SRV record should have the following information:

```bash
_opencap._tcp.DOMAIN_NAME_HERE. 300 IN SRV 5 12 443 DOMAIN_NAME_HERE.
```

For example:

```bash
_opencap._tcp.example.com. 300 IN SRV 5 12 443 example.com.
```

Broken down:

```bash
Record Name = _opencap._tcp.example.com.
TTL = 300
Class = IN
Record Type = SRV
Priority = 5
Weight = 12
Port = 443
Target = example.com
```

Fill in the information that your provider asks for, it may not all be required.

Verify that the record is live using google's tool: https://dns.google.com/

Type in your record name (e.g. \_opencap.\_tcp.example.com.) to the search bar. The result should contain the domain name you set. If it doesn't, double check your SRV record. Also try waiting a few minutes, it usually takes a bit for the DNS record to propogate.

### Forward your computers web ports through your router

Your router needs to open its webserver ports to the public and communicate traffic to your computer. The program has a helper tool:

#### Linux

```bash
./go-server --openport 80
```

and

```bash
./go-server --openport 443
```

#### Windows Command Line

```cmd
go-server.exe --openport 80
```

and

```cmd
go-server.exe --openport 443
```

If it prints an error message, you will probably need to do this step manually:

https://kb.netgear.com/24290/How-do-I-add-a-custom-port-forwarding-service-on-my-Nighthawk-router

### Modify the file named ".env"

Set JWT_SECRET equal to some random text (no spaces, only letters) at least 50 characters long.

Set CREATE_USER_PASSWORD to a password that you can use to add new aliases to your server.

Set the DOMAIN_NAME equal to the domain name you are using.

Everything else can be left alone.

### Setup the database

#### Install sqlite3

##### Windows

https://www.sqlite.org/download.html

In "Precompiled Binaries for Windows" download the x64 version if you are on a 64 bit machine (newer machine)

##### Linux

```bash
sudo apt-get update
sudo apt-get install sqlite3 libsqlite3-dev
```

##### Setup the OpenCAP database configuration

##### Windows Command Line

```bash
go-server.exe --setupdatabase
```

##### Linux

```bash
./go-server --setupdatabase
```

### Run the server

#### Windows

```bash
go-server.exe
```

#### Linux

```bash
./go-server
```

If all went well, your server is up and running!

### Troubleshoot

One thing that can go wrong here is that the go-server program was unable to correctly setup HTTPS (secure encryption). If this is the case, double check that your DNS records are correct and that your ports are being forwarded correctly.

You may also need to make sure that windows (or another operating system's) firewall isn't blocking it from doing its job. You can go to your firewall settings and allow the go-server program by browsing for it.

## Using the server

Once it is up and running, you can create users by making requests of the following format:

```json
POST https://example.com/v1/users
content-type: application/json

{
    "alias": "username$myserver.com",
    "password": "myNewUserPassword",
    "create_user_password": "somepassword"
}
```

All other requests follow the OpenCAP protocol.

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
