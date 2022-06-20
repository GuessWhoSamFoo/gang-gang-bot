## Getting Started

Install the [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli).

Set the environment variables with:

```
heroku config:set DISCORD_GUILD_ID="$GUILD_ID" DISCORD_TOKEN="$TOKEN"
```

Instead of using the local `config.yml` file for secrets, create a `.env` file instead:

```
DISCORD_GUILD_ID=$GUILD_ID
DISCORD_TOKEN=$TOKEN
```

Start the local worker by:

```
heroku local
```

It will use secrets from `.env` and allow for multiple development environments.

The `Procfile` can be used have finer control over what happens during runtime on a Heroku dyno.
