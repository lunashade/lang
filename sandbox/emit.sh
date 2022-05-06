#!/bin/bash
set -e
code=${1:-sandbox/code.c}
llcode=${code%.*}.ll
clang -c -S -emit-llvm -o $llcode $code
cat $llcode
