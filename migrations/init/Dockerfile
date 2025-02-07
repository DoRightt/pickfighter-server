FROM postgres

COPY ./migrations/init/create_database.sql /init/
COPY ./migrations/init/create_schemas.sql /init/
COPY ./migrations/init/auth_service_init.sql /init/
COPY ./migrations/init/fighters_service_init.sql /init/
COPY ./migrations/init/events_service_init.sql /init/

COPY ./migrations/init/init.sh /docker-entrypoint-initdb.d/



