This dir contains migration scripts

# Naming convention
Migrations (should) come in pairs, one up script for upgrade and one down for downgrade.
```
<major><sprint><minor>__<description>.<up|down>.sql
```
# Run migrations
Follow instructions here: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
```
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.11.0/migrate.darwin-amd64.tar.gz | tar xvz
./migrate.darwin-amd64 -verbose -database mysql://$DB_MIGRATE_CONN_STRING -path $MIGRATION_PATH up
```
