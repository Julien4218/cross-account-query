#!/usr/bin/env bash

date_from=2023-02-01
date_to_max=2023-03-01

while [ "$date_from" != "$date_to_max" ]; do 
  date_to=$(date -j -v +1d -f "%Y-%m-%d" $date_from +%Y-%m-%d)
  echo "Querying $1 from:$date_from to:$date_to..."

  QUERY_FROM="$date_from 00:00:00 PT" QUERY_TO="$date_to 00:00:00 PT" go run . $1 --json > "output_$date_from.json"

  # mac option for d decl (the +1d is equivalent to + 1 day)
  date_from=$(date -j -v +1d -f "%Y-%m-%d" $date_from +%Y-%m-%d)
done
