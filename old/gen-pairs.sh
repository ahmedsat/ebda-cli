#!/bin/bash

# todo: embed this script into the main application

# Read from stdin into an array
mapfile -t list

# Generate pairs and output as TSV
for ((i=0; i<${#list[@]}; i++)); do
    for ((j=i+1; j<${#list[@]}; j++)); do
        echo -e "${list[i]}\t${list[j]}"
    done
done

