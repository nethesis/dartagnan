# foil

Foil populates Redis cache with systems having a valid subscription:

1. it flushes all Redis keys
2. read the database and set a key for each valid system

Foil can be scheduled to run at regular intervals (eg. once a day).

