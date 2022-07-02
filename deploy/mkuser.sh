#!/usr/bin/bash

if id nfarl  &>/dev/null ; then
    echo 'User nfarl exists'
else
    echo 'Creating user nfarl'
    useradd -p nfdfhyMCZWk6w -c "NFARL" -m nfarl
fi
