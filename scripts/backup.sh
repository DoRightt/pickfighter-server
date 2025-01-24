#!/bin/bash

TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
BACKUP_DIR="$HOME/backups/pickfighter/database"
DB_NAME="pickfighter_db"
DB_USER="postgres"
DB_HOST="localhost"

# Creation DB dump
pg_dump -U $DB_USER -h $DB_HOST $DB_NAME > "$BACKUP_DIR/db_backup_$TIMESTAMP.sql"

# Removing old backups (older than 7 days)
find $BACKUP_DIR -type f -name "*.sql" -mtime +7 -exec rm {} \;

echo "Backup completed: $BACKUP_DIR/db_backup_$TIMESTAMP.sql"