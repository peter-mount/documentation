#!/bin/sh

newPath=$1
fontCss=font.css

src="https://cdn.jsdelivr.net/gh/rastikerdar/vazir-font@v27.0.1/dist/"
FILES=$(grep "url(" ${fontCss}|cut -f2 -d'('|cut -f1 -d')')
for file in $FILES
do
  fileName=$(basename $file | sed "s/'//g")
  dirName=$(dirname $file)

  echo "$file -> $fileName"
  wget --quiet -O $fileName "${src}$(echo -n $file | sed "s/'//g")"

  #sed -i "s|${dirName}|${newPath}|g" ${fontCss}
done
