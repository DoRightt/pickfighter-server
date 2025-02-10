#!/bin/bash
psql -U $POSTGRES_USER -d $POSTGRES_DB -f /init/create_database.sql
psql -U $POSTGRES_USER -d $POSTGRES_DB -f /init/create_schemas.sql
psql -U $POSTGRES_USER -d $POSTGRES_DB -f /init/auth_service_init.sql
psql -U $POSTGRES_USER -d $POSTGRES_DB -f /init/fighters_service_init.sql
psql -U $POSTGRES_USER -d $POSTGRES_DB -f /init/events_service_init.sql