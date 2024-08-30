# Adding data to database

There are several "easier" ways to do this, but this is the most repeatable for different levels of experience as well as different operating systems.

1. Start database - ./postgres_start.sh
   - note: you can make this persistent by specifying volumes in the script such as adding:
   ```
   -e PGDATA=/var/lib/postgresql/data/pgdata \
   -v /Users/frankmoley/.local/docker/data:/var/lib/postgresql/data \
   ```
   
2. Exec into docker container
   ```
   docker exec -it local-pg /bin/bash
   ```

3. Launch psql from inside the docker container
   ```
   psql -U postgres
   ```
4. Copy/paste schema file and then data file from this directory into psql


