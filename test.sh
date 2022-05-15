#!/bin/bash
# test script for lang
TARGET=$1  # target executable
TMPDIR=tmp/
mkdir $TMPDIR

function check {
    want=$1
    input=$2

    echo "$input" | ${TARGET} > $TMPDIR/tmp.ll
    cat $TMPDIR/tmp.ll | lli
    got=$?
    if [[ "$want" == "$got" ]]; then
        echo "[SUCCESS] ${input} => ${got}"
    else
        echo "[FAIL] ${input} => ${want} but ${got}";
        exit 1;
    fi
}


function main {
    echo "target: ${TARGET}"
    check 0 "0"
    check 1 "1"
    check 255 "255"
    check 255 "  255  "
    check 2 "1+1"
    check 6 "2*3"
    check 7 "1+2*3"
    check 8 "280 / 20 - 2 * 3"
    echo ok
}

main
