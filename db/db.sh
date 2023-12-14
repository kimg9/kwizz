#!/bin/bash

psql -c "drop database kwizz"
psql -c "create database kwizz"
psql kwizz -f db.sql