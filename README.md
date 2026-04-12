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

## Configurations
- `--messages-file` Path to file with custom messages in yaml format
- `--test-guilds` List of test guild IDs, separated by commas, where bot register commands. This can be configured with the environment variable `TEST_GUILD_ID`
- `--token` Discord token, this can be configured with the environment variable `DISCORD_TOKEN`
- `--birthdays-file` Path to the file with birthdays data in yaml format. This can be configured with the environment variable `BIRTHDAYS_FILE`
- `--advertise-today-birthday` Flag to run the job that advertises today's birthday and exit
- `--advertise-channel` Channel ID used to advertise birthdays. This can be configured with the environment variable `ADVERTISE_CHANNEL`
- `--timezone` Timezone for the bot (e.g. `America/Argentina/Buenos_Aires`). This can be configured with the environment variable `TZ`

## Custom messages
Go templates are used to configure custom responses.

yaml keys for messages:
- `nextHoliday`: response for next-holiday command
- `daysLeft`: response for days-left command
- `nextLargeHoliday`: response for next-large-holiday command
- `holidaysOfMonth`: response for holidays-of-month command
- `activityStatus`: response for the bot's custom activity status
- `nextBday`: response for the next-birthday advertisement
- `todayBday`: response for the today-birthday advertisement
- `bdayLeft`: response for the bday-left command

The keys passed for commands `nextHoliday`, `daysLeft`, `nextLargeHoliday` and `activityStatus` are:

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

The keys passed for command `holidaysOfMonth` are:

- `Month`: Month name.
- `Count`: Number of holidays in for this month.
- `HolidaysList`: The list of holidays with the full information.
- `Adjacents`: List of list of adjacents holidays, to determine the large holidays, weekends are consider holidays in this lists. Eg, if the holiday is in friday, the follow `saturda` and `sunday` are added as adjacents. 

The keys passed for commands `nextBday` and `todayBday` are:

- `Name`: Name of the person
- `Username`: Username or ping representation
- `Date`: Birthday date string
- `DaysLeft`: Days left to birthday

The keys passed for command `bdayLeft` are:

- `Username`: Username or ping representation
- `DaysLeft`: Days left to birthday
