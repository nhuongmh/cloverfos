
### Postgres SQL
- In case of upgrade postgresql version fail: consider migrate data in /var/lib/pgsql
- Start postgresql service: sudo systemctl start postgresql
- postgresql admin user pass: clover.fox
- create user: sudo -u postgres psql
- CREATE USER clover PASSWORD 'foxie' CREATEDB;


NORMAL dev use:
- psql -U clover -h localhost -d postgres