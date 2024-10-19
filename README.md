# gator

Blog aggregator project for boot.dev.

## requirements

- Postgres
- Go

## installation

To install, use 'go install github.com/ffb6c1/bootdev_aggregator@latest' from the command line.

## usage and commands

Before use, you'll need to register your username.

If there are spaces in any part of a command (for example, the name of a feed), wrap that part in quotes.

**the 'agg' command runs indefinitely**. It is intended to run in the background in a separate command window. To end the process manually, use ctrl+c or cmd+c.

### User Management

- bootdev_aggregator register [ username ] - adds a new user and logs them in.
- bootdev_aggregator login [ username] - logs in using the specified username.
- bootdev_aggregator users - view a list of users.

### Feed Management

- bootdev_aggregator addfeed [ feed name ] [ url ] - adds a new feed to the feed list and subscribes the active user to it.
- bootdev_aggregator feeds - displays a list of the feeds added by all users.
- bootdev_aggregator follow [ url ] - subscribes to a feed previously added by another user.
- bootdev_aggregator unfollow [ url ] - unsubscribes from a feed.
- bootdev_aggregator following - displays a list of feeds followed by the active user.

### Feed Viewing

- bootdev_aggregator agg [ duration between requests (s/m/h) ] - commands the program to automatically gather new posts from the followed feeds every time the duration passes.
- bootdev_aggregator browse [ optional - number ] - view the current n latest posts from your followed feeds. Number defaults to 2. 
