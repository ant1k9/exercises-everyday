#!/bin/bash
set -euo pipefail

cd /tmp

HEROKU_APP="exercises-everyday"
DB_BACKUP_PATH="exercises.dump"
DB_BACKUP_ARCHIVE="$DB_BACKUP_PATH.zip"
DROPBOX_PATH="/exercises/$(date +'%Y%m%d_%H%M%S').dump.zip"

heroku="/usr/local/bin/heroku"
CURRENT_BACKUP_VERSION="$(
    $heroku pg:backups \
        -a $HEROKU_APP \
        | grep -Eo 'b[0-9]+' \
        | tail -1\
)"

# Delete last backup
$heroku pg:backups:delete \
    -a $HEROKU_APP \
    --confirm "$HEROKU_APP" \
    "$CURRENT_BACKUP_VERSION"

# Create new backup
$heroku pg:backups:capture \
    -a $HEROKU_APP

# Download backup
$heroku pg:backups:download \
    -a $HEROKU_APP \
    -o "$DB_BACKUP_PATH"

zip "$DB_BACKUP_ARCHIVE" "$DB_BACKUP_PATH"
rm "$DB_BACKUP_PATH"

curl -X POST https://content.dropboxapi.com/2/files/upload \
    --header "Authorization: Bearer $DB_BACKUP_TOKEN" \
    --header "Content-Type: application/octet-stream" \
    --header "Dropbox-API-Arg: {\"path\": \"$DROPBOX_PATH\",\"mode\": \"add\",\"autorename\": true,\"mute\": false,\"strict_conflict\": false}" \
    --data-binary @"$DB_BACKUP_ARCHIVE"

rm "$DB_BACKUP_ARCHIVE"
