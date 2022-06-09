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
    check 25 "main(){3; {2; 25}}"
    check 14 "main(){3; {2; 25}; 14}"
    check 1 "main(){ 1==1 }"
    check 0 "main(){ 0==1 }"
    check 0 "main(){ 1!=1 }"
    check 1 "main(){ 0!=1 }"
    check 1 "main(){ 1<=1 }"
    check 1 "main(){ 0<=1 }"
    check 1 "main(){ 1>=1 }"
    check 0 "main(){ 0>=1 }"
    check 0 "main(){ 1<1 }"
    check 1 "main(){ 0<1 }"
    check 0 "main(){ 1>1 }"
    check 0 "main(){ 0>1 }"
    check 25 "main(){ if 1==1 then 5*5 }"
    check 25 "main(){ if 1==1 then 5*5 else 5*5-5 }"
    check 20 "main(){ if 1==0 then 5*5 else 5*5-5 }"
    check 0 "main(){ if 1==0 then 5*5 }"
    check 25 "main() {if {if 1 then 1 else 0} then {if 1 then 25 else 0} else {if 1 then 0 else 0}}"
    echo ok
}

main
