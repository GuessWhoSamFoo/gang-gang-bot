# gang-gang bot

This Discord bot is intended for use by the [SALSA Discord Group](https://discord.gg/jmKXruqvz4).

## Roadmap

 * Recurring events
 * Integration with Google/Outlook/etc Calendar
 * Accessibility
 * Viewing, sorting, filtering events
 * Localization
 * Custom signup rules

## Development

To obtain a Guild ID, open the Discord client on Desktop then go to `User Settings` > `Advanced`
and enable developer mode. Right-click a server icon on the left to see a `Copy ID` option.

To setup a bot for development purposes, go to [Applications](https://discord.com/developers/applications)
then create a new application. Under `Settings` > `Bot`, there is an option to generate a token.

It is also recommended to create a new server for testing.

1. Install Go 1.18+

2. Create a `config.yaml` file with the Guild ID and Bot Token in the repo root. This file contains secrets and
should not be committed to version control or uploaded online without encryption.

```
discord:
  guild_id: {{ GUILD_ID }}
secret:
  token: {{ DISCORD_BOT_TOKEN }}
```

If using Heroku, see [docs/](/docs/heroku.md).

3. Add the bot to a server for testing. See [this guide](https://discordjs.guide/preparations/adding-your-bot-to-servers.html#adding-your-bot-to-servers)
for detailed instructions. It is possible to invite a users to app test or create a link through the OAuth2 URL Generator.

Minimally, `bot` and `application.commands` scopes must be enabled with the ability to read/send messages.

```
https://discord.com/api/oauth2/authorize?client_id=920479421478076468&permissions=2147568640&scope=bot%20applications.commands
```

4. Start the bot by running `go run cmd/gang-gang-bot/main.go` from the root directory of the git repo.

5. A new Slash command `/event` should be available which still start a DM sequence for creating a new event.

## Resources

 - [Visualizer](https://autocode.com/tools/discord/embed-builder/) for embed and components
 - [Discord API Documentation](https://discord.com/developers/docs/intro)
 - [Discord.js Documentation](https://discord.js.org/#/docs/discord.js/stable/general/welcome)
 - [discordgo](https://github.com/bwmarrin/discordgo)
