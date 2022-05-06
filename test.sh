#!/bin/bash
# test script for lang
TARGET=$1  # target executable

function check {
    want=$1
    input=$2

    echo "$input" | ${TARGET} | lli
    got=$?
    if [[ "$want" != "$got" ]]; then
        echo "input: {$input}, want: ${want} but got: ${got}";
        exit 1;
    fi
}


function main {
    echo "target: ${TARGET}"
    check 0 "0"
    check 1 "1"
    check 255 "255"
    echo ok
}

main
