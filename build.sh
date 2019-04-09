#!/bin/bash

target=$1
if [ -z "$target" ]; then
	echo "usage: $0 <build-target> eg: windows, linux etc."
	exit 1
fi
    output="hunterD-$target-amd64"
    if [ $target = "windows" ]; then
	 output+=".exe"
    fi

   env GOOS=$target GOARCH=amd64 go build -o $output .

if [ $? -ne 0 ]; then 
	   echo "Failed to build $target target."
	   exit 1
fi

