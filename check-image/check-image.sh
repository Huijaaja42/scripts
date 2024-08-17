#!/bin/bash

if [ -z "$1" ]; then
	exit 1
fi

p=$(find "$1" -type f \( -name '*.png' -o -name '*.PNG' -o -name '*.jpg' -o -name '*.JPG' -o -name '*.nef' -o -name '*.NEF' \) -printf '%p,')
c=$(echo "$p" | tr ' ,' '- ' | wc -w)
echo "Found $c files"

IFS=',' read -r -a a <<< "$p" 
i=1
for f in "${a[@]}"
do
	echo -e "\nChecking $i/$c $f"
	ffmpeg -nostdin -i "$f" -loglevel error -f null - 2>&1
	i=$((i+1)) 
done
