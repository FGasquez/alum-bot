# Alum Discord Bot

## Build Docker Image
```bash
$ docker build -t discord-bot .

```

## Build go binary
```bash
$ go build -o bin/main cmd/alum-bot/*

# Run the binary
$ ./bin/main -t <DISCORD_BOT_TOKEN>
```

## Configurtions
- `--messages-file` Path to file with custom messages in yaml format
- `--test-guilds` List of test guild IDs, separated by commas, where bot register commands. this can be configured with the environment variable `TEST_GUILD_ID`
- `--token` Discord token, this can be configured with environment variable `DISCORD_TOKEN`

## Custom messages
Go templates are used to configure custom responses.

yaml keys for messages:
- `nextHoliday`: response for nex-holiday command
- `daysLeft`: response for days-left command
- `nextLargeHoliday`: response for next-large-holiday command
- `holidaysOfMonth`: response for holidays-of-month command

The keys passed for commands `nextHoliday`, `daysLeft` and `nextLargeHoliday` is:

- `HolidayName`: Name of holiday
- `DaysLeft`: Days left to holiday
- `FormattedDate`: Date formated to spanish
- `NamedDate`:
    - `Day`: Day name
    - `Month`: Month name
- `RawDate`: 
    - `Day`: Day number
    - `Month`: Month number
    - `Year`: Year
- `FullDate`: Date in format `yyyy-mm-dd`
- `IsToday`: Boolean, true if the holiday is today
- `Adjacents`: Adjacents holidays

The keys passed for command `holidaysOfMonth` is:

- `Month`: Month name.
- `Count`: Number of holidays in for this month.
- `HolidaysList`: The list of holidays with the full information.
- `Adjacents`: List of list of adjacents holidays, to determine the large holidays, weekends are consider holidays in this lists. Eg, if the holiday is in friday, the follow `saturda` and `sunday` are added as adjacents. 
