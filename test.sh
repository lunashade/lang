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
    check 0 "main(){0}"
    check 1 "main(){1}"
    check 255 "main(){255}"
    check 2 "main(){1+1}"
    check 6 "main(){2*3}"
    check 7 "main(){ 1+2*3 }"
    check 8 "main(){ 280 / 20 - 2 * 3 }"
    check 14 "main(){ 2 * (3 + 4) }"
    check 0 "main(){10;}"
    check 0 "main(){}"
    check 25 "main(){1+1; 5*5}"
    check 0 "main(){1+1; 5*5;}"
    echo ok
}

main
