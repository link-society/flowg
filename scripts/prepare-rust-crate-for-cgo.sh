#!/bin/sh
set -eu

crate_dir=${1:?Usage: prepare-rust-crate-for-cgo.sh <crate-directory>}

host=$(rustc -vV | sed -n 's/^host: //p')
objcopy="$(rustc --print sysroot)/lib/rustlib/$host/bin/llvm-objcopy"
if [ ! -x "$objcopy" ]; then
    echo "Missing llvm-objcopy; install it with: rustup component add llvm-tools" >&2
    exit 1
fi

# Rust has already performed LTO. Remove embedded bitcode from the final
# archive, keeping the native code and symbols needed by cgo. Use a separate
# directory so Go does not scan Cargo's intermediate archives for LTO sections.
mkdir -p "$crate_dir/target/native"
"$objcopy" \
    --remove-section=.llvmbc \
    --remove-section=.llvmcmd \
    --remove-section='__LLVM,*' \
    "$crate_dir/target/release/libflowg_vrl.a" "$crate_dir/target/native/libflowg_vrl.a"
