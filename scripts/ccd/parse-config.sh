#!/bin/bash

awk '
BEGIN {
  printf "ccd:\n"
}
{
  printf "  - name: %s\n", tolower($1)
  printf "    axes:\n"
  printf "      x:\n"
  printf "        beta: %.7e\n", $3
  printf "        dispersion: %.7e\n", $5
  printf "      z:\n"
  printf "        beta: %.7e\n", $4
  printf "        dispersion: %.7e\n", $6
}' "$1"