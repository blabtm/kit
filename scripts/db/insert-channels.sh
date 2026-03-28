#!/bin/sh

psql -h 0.0.0.0 -p 5433 -U postgres -c "\\copy channels FROM '$1' DELIMITER ',' CSV HEADER"
