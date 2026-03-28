#!/bin/bash

awk -v side="$2" '
BEGIN {
  system("rm -f v2k.*")
}
{
  sub(/\./, "", $1)
  printf "%s,%f\n", $1, $4 >> "v2k.beam.cur"

  if (side == "l") {
    printf "%s,%f\n", $1, $23 >> "v2k.ccd.1m1l.sx"
    printf "%s,%f\n", $1, $24 >> "v2k.ccd.1m1l.sz"
    printf "%s,%f\n", $1, $31 >> "v2k.ccd.1m2l.sx"
    printf "%s,%f\n", $1, $32 >> "v2k.ccd.1m2l.sz"
    printf "%s,%f\n", $1, $47 >> "v2k.ccd.2m1l.sx"
    printf "%s,%f\n", $1, $48 >> "v2k.ccd.2m1l.sz"
    printf "%s,%f\n", $1, $39 >> "v2k.ccd.2m2l.sx"
    printf "%s,%f\n", $1, $40 >> "v2k.ccd.2m2l.sz"
    printf "%s,%f\n", $1, $55 >> "v2k.ccd.3m1l.sx"
    printf "%s,%f\n", $1, $56 >> "v2k.ccd.3m1l.sz"
    printf "%s,%f\n", $1, $63 >> "v2k.ccd.3m2l.sx"
    printf "%s,%f\n", $1, $64 >> "v2k.ccd.3m2l.sz"
    printf "%s,%f\n", $1, $79 >> "v2k.ccd.4m1l.sx"
    printf "%s,%f\n", $1, $80 >> "v2k.ccd.4m1l.sz"
    printf "%s,%f\n", $1, $71 >> "v2k.ccd.4m2l.sx"
    printf "%s,%f\n", $1, $72 >> "v2k.ccd.4m2l.sz"
  } else {
    printf "%s,%f\n", $1, $27 >> "v2k.ccd.1m1r.sx"
    printf "%s,%f\n", $1, $28 >> "v2k.ccd.1m1r.sz"
    printf "%s,%f\n", $1, $35 >> "v2k.ccd.1m2r.sx"
    printf "%s,%f\n", $1, $36 >> "v2k.ccd.1m2r.sz"
    printf "%s,%f\n", $1, $51 >> "v2k.ccd.2m1r.sx"
    printf "%s,%f\n", $1, $52 >> "v2k.ccd.2m1r.sz"
    printf "%s,%f\n", $1, $43 >> "v2k.ccd.2m2r.sx"
    printf "%s,%f\n", $1, $44 >> "v2k.ccd.2m2r.sz"
    printf "%s,%f\n", $1, $59 >> "v2k.ccd.3m1r.sx"
    printf "%s,%f\n", $1, $60 >> "v2k.ccd.3m1r.sz"
    printf "%s,%f\n", $1, $67 >> "v2k.ccd.3m2r.sx"
    printf "%s,%f\n", $1, $68 >> "v2k.ccd.3m2r.sz"
    printf "%s,%f\n", $1, $83 >> "v2k.ccd.4m1r.sx"
    printf "%s,%f\n", $1, $84 >> "v2k.ccd.4m1r.sz"
    printf "%s,%f\n", $1, $75 >> "v2k.ccd.4m2r.sx"
    printf "%s,%f\n", $1, $76 >> "v2k.ccd.4m2r.sz"
  }
}' "$1"
