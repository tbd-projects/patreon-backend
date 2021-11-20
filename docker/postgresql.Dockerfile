FROM postgres:14

RUN apt-get update \
      && apt-cache showpkg postgresql-$PG_MAJOR-rum \
      && apt-get install -y --no-install-recommends \
           postgresql-$PG_MAJOR-rum \
      && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /docker-entrypoint-initdb.d