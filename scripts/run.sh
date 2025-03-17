#!/bin/sh
set -e

echo "Running database migrations..."
for migration in /app/migrations/*.sql; do
    echo "Applying migration: $(basename "$migration")"
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f "$migration"
done
echo "Migrations completed successfully!"

echo "Starting order service..."
/app/order-service-test