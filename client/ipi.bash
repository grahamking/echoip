#!/bin/sh
# OSX USERS: REMOVE '-q 1'

# I put this in my .bashrc:
#alias ipi='echo "" | nc -q 1 -u plebis.net 7777'

echo "" | nc -q 1 -u plebis.net 7777
