#!/bin/bash

APIMODELS=api/models

echo "stripping unneeded code from generated models"
for FILE in ${APIMODELS}/*; do
  echo "${FILE}"

  TMPFILE=$(mktemp -d "/tmp/$(basename $0).XXXXXX")
  TMPFILE="${TMPFILE}/temp.go"

  GOOD=1
  IMPORT=0
  MODEL=0

  cat ${FILE} | while IFS="" read LINE; do
    if [[ $LINE == 'import (' ]]; then
      IMPORT=1
      GOOD=0
    fi

    if [[ $LINE == type* ]]; then
      MODEL=1
      GOOD=1
    fi

    if [[ "$GOOD" == '1' ]]; then
      echo "${LINE}" >> $TMPFILE
    fi

    if [[ $LINE == ')' ]]; then
      if [[ "$IMPORT" == '1' ]]; then
        IMPORT=0
        GOOD=1
      fi
    fi

    if [[ $LINE == '}' ]]; then
      if [[ "$MODEL" == '1' ]]; then
        MODEL=0
        GOOD=0
      fi
    fi
  done

  cp $TMPFILE $FILE
  echo ""
done
