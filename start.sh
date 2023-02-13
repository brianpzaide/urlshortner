#!/bin/sh
while ! nc -z urlshortner_db 5432; do sleep 3; done
./urlshortner
