#!/bin/bash
FILE="buildnr.txt"
BUILD=`cat $FILE`
NEWBUILD=`expr $BUILD + 1`
echo "$NEWBUILD" > $FILE
echo "$NEWBUILD"
