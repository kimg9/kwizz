#!/bin/bash

psql -c "drop database kwizz"
psql -c "create database kwizz"
cat db/db.sql | psql kwizz