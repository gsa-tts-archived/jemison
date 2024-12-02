#!/bin/bash

cat $1 | while read line 
do
  http --ignore-stdin put http://localhost:10001/api/pack scheme=https path="/" api-key=lego host=${line}
done
