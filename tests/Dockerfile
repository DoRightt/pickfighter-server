FROM postgres:alpine

COPY init.sql /docker-entrypoint-initdb.d/

ENV POSTGRES_DB=fighters_db
ENV POSTGRES_USER=postgres