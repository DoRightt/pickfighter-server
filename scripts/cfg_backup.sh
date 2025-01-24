#!/bin/bash

# Backup directory
BACKUP_DIR="$HOME/backups/pickfighter"
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")

# Project root
PROJECT_DIR="."

# Backup list
FILES_TO_BACKUP=(
    ".env"
    "auth/configs/config.yaml"
    "events/configs/config.yaml"
    "fighters/configs/config.yaml"
    "pickfighter/configs/config.yaml"
    "scraper/configs/proxy.yaml"
    "auth/logs/log.json"
    "events/logs/log.json"
    "fighters/logs/log.json"
    "pickfighter/logs/log.json"
    "scraper/logs/log.json"
    "scraper/collection/fighters.json"
)

# Create backup directory if not exists
mkdir -p $BACKUP_DIR

# Copy files and directories
for FILE in "${FILES_TO_BACKUP[@]}"; do
    if [ -e "$PROJECT_DIR/$FILE" ]; then
        # Creating file directory in backup directory
        BACKUP_PATH="$BACKUP_DIR/$FILE"
        BACKUP_DIR_PATH=$(dirname "$BACKUP_PATH")
        mkdir -p "$BACKUP_DIR_PATH"
        
        # Copy a file to appropriate backup directory
        cp -r "$PROJECT_DIR/$FILE" "$BACKUP_PATH"
        echo "Backup completed for: $PROJECT_DIR/$FILE"
    else
        echo "Warning: $PROJECT_DIR/$FILE not found"
    fi
done

echo "Backup process finished."