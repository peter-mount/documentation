#!/bin/sh

newPath=$1
fontCss=font.css

FILES=$(grep "url(" ${fontCss}|cut -f2 -d'('|cut -f1 -d')')
for file in $FILES
do
  fileName=$(basename $file)
  dirName=$(dirname $file)

  echo "$file -> $fileName"
  wget --quiet -O $fileName $file

  sed -i "s|${dirName}|${newPath}|g" ${fontCss}
done
