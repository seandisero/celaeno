# celaeno
a terminal-based real-time encripted messaging app

spend a lot of time in the terminal, well now you don't have to leave it to chat with friends, or coworkers *caugh caugh*. messages are symetrically encrypted so only you and your fellow user with the same cipher will ever be able to read what you talk about.

## Motivation
I didn't like switching to my browser just to respond to a chat, I wanted a way to stay in the same window (provided one uses something like tmux), and not deviate too much from what I was doing. So I built this terminal app that lets me do just that. 

## Quick Start Client

### 1. install using go

```bash
# clone repo

cd celaeno
go instal cmd/celaeno-cli/main.go
```

### 2. setting up a user and chat
start by registering a user.
```bash
# run client
celaeno-cli

# after startup
/register <username> <password>
```
once you've created a user, you can log in and create a chat.
```bash
/login <username> <password>

# after login confirmation
/create-chat
```
the chat created is stored under your username so others can join by using the connect command
```bash
/connect <username>
```
you'll need to set a cipher for each client chatting 
```bash
/set cipher <cipher>
```
this must be the same for everyone chatting
### 3. start chatting
anyone can set a display name if they don't want to use their username
```bash
/set displayname <new_displayname>
```
then just type any messages you wish into the window and start chatting.

## Quick Start Server
use any sqlite database with a single table for users
```bash
CREATE TABLE users(
	id BLOB UNIQUE PRIMARY KEY,
	username TEXT NOT NULL,
	displayname TEXT,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
    hashed_password TEXT NOT NULL
);
```
there are only a few environment variables that need to be set for the server
- DB_URL=<the_url_for_your_database>
- JWT_SECRET=<your_secret>
and the build command is 
```bash
go build -tags netgo -ldflags '-s -w' -o celaeno-server ./cmd/celaeno-server
```
the start command is then
```bash
./celaeno-server
```
## Usage
available commands:
```bash
# this command tries to start up the server if using a service like Render
# and it needs to wake back up. the program will respond with a 'good to go'
# message when done.
/startup 

# register a new user
/register <username> <password>

# log in as a user
/login <username> <password>

/logout

# print current user information
/whoami

# delete current user, must be logged in as user, celaeno will prompt for password
/deleteme 

/create-chat

# must know the username of who you wish to connect to and chat must 
# already be created.
/connect <username>

# leave the current chat
/leave

#list currently available chats
/chats
```

## Contributing 
if you'd like to contribute to the project please fork the repository and open a pull request to the `main` branch.
