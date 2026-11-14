# vrl

The package at `internal/utils/langs/vrl` binds FlowG to
[Vector Remap Language (VRL)](https://vector.dev/docs/reference/vrl/), the
language transformer nodes use to reshape log events.

It exists because VRL is implemented in Rust. This package wraps that
implementation behind a small Go API through cgo, so the pipeline engine can
compile and run VRL programs without knowing anything about the foreign function
interface underneath.

## Responsibilities

- **Compilation** — `NewScriptRunner` compiles VRL source into a reusable
  `ScriptRunner`, returning a `*CompileError` when the source is invalid.
- **Transformation** — `ScriptRunner.TransformLog` runs a compiled program
  against a log event and returns the resulting event(s), surfacing a
  `*EvalError` on runtime failures.
- **Resource management** — `ScriptRunner.Close` releases the native resources
  held by a runner.

## Usage notes

- A `ScriptRunner` owns native resources; always call `Close` when finished.
- A runner is **not** safe for concurrent use — serialize calls or use one
  runner per goroutine.
- A program may emit zero, one or several events, so `TransformLog` always
  returns a slice. Nested objects and arrays are flattened into dotted field
  names.

## Layout

- **ffi.go / ffi.h** — the cgo binding and the C header it links against.
- **errors.go** — the `CompileError` and `EvalError` types.
- **rust-crate/** — the Rust implementation, built into the static library the
  binding links against.

## Building

Install Rust's object tools once with `rustup component add llvm-tools`, then
run `task build:backend:vrl`. After Cargo's release build,
`scripts/prepare-rust-crate-for-cgo.sh <crate-directory>` copies the library to
`rust-crate/target/native/libflowg_vrl.a`, removing
embedded LLVM bitcode while preserving the optimized native code. Both the
Task and Docker builds use this step because Go's cgo checks reject bitcode
sections and recursively inspect library search directories. Keeping the
prepared library separate leaves Cargo's intermediate LTO archives intact.
