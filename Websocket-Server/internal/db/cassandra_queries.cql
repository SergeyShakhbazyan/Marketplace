CREATE KEYSPACE "messenger_keyspace" WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

CREATE TABLE messenger_keyspace.chatroom (
                                             chatID UUID,
                                             user1 UUID,
                                             user2 UUID,
                                             PRIMARY KEY ((user1, user2), chatID)
);
CREATE INDEX ON messenger_keyspace.chatroom (user1);
CREATE INDEX ON messenger_keyspace.chatroom (user2);

CREATE MATERIALIZED VIEW messenger_keyspace.chatroom_by_id AS
SELECT chatID, user1, user2
FROM messenger_keyspace.chatroom
WHERE chatID IS NOT NULL
  AND user1 IS NOT NULL
  AND user2 IS NOT NULL
PRIMARY KEY (chatID, user1, user2);



CREATE TABLE messenger_keyspace.message(
                                           id UUID,
                                           content TEXT,
                                           senderID UUID,
                                           chatID UUID,
                                           timeStamp Timestamp,
                                           primary key (chatID, timeStamp, ID, senderID)
)WITH CLUSTERING ORDER BY (timeStamp DESC);
