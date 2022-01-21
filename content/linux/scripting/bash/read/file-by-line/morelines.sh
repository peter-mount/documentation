#!/bin/bash
input="$1"
(
  lc=1
  while IFS= read -r line
  do
    printf "%03d %s\n" ${lc} "${line}"
    lc=$((lc+1))
  done
) < "${input}"
