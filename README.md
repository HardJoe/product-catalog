# Product Catalog

How to use:
1. Clone this repo
```bash
git clone https://github.com/HardJoe/product-catalog
```
2. Build docker containers
```bash
docker compose up --build
```
3. Enter db container
```bash
docker exec -it product-catalog-db-1 /bin/sh
```
4. Enter PostgreSQL shell
```bash
psql -U divrhinotrivia
```
5. Insert sql queries in test.sql to create tables

