CREATE TABLE MessageDatabase (
    slackID TEXT PRIMARY KEY,
    discordID TEXT,
    ChannelName TEXT,
    slackChannelID TEXT,
    discordChannelID TEXT
);

CREATE TABLE ThreadDatabase (
    slackThreadID TEXT PRIMARY KEY,
    discordThreadID TEXT,
    ChannelName TEXT,
    slackChannelID TEXT,
    discordChannelID TEXT
);