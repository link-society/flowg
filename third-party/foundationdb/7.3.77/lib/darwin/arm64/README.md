# FoundationDB client library for macOS (arm64)

Apple does not publish a prebuilt macOS client library for FoundationDB, so the
`libfdb_c.dylib` used to build the macOS (`darwin/arm64`) release of FlowG must
be built from source and placed in this directory:

```
third-party/foundationdb/7.3.77/lib/darwin/arm64/libfdb_c.dylib
```

## Building it

On an Apple Silicon Mac (requires `cmake`, `ninja` and `boost`, e.g. via
Homebrew, plus the Xcode command line tools):

```bash
git clone https://github.com/apple/foundationdb.git
cd foundationdb
git checkout 7.3.77          # must match the pinned version (API version 730)
mkdir build && cd build
cmake -G Ninja -DCMAKE_BUILD_TYPE=Release ..
ninja fdb_c                  # build only the C client target
```

## Normalizing the install name

The `flowg-server` binary records the dylib's install name at link time, so it
must match where the library is installed at runtime (see the cluster setup
guide, which installs it to `/usr/local/lib/libfdb_c.dylib`):

```bash
install_name_tool -id /usr/local/lib/libfdb_c.dylib libfdb_c.dylib
otool -D libfdb_c.dylib   # verify the install name
```

Then copy the resulting `libfdb_c.dylib` into this directory.

## Making the release binary relocatable

Keep the install name as `/usr/local/lib/libfdb_c.dylib` above: the backend build
(`scripts/build.taskfile.yml`) relies on it to re-point `flowg-server` at
`@rpath/libfdb_c.dylib` with `install_name_tool -change`, and links the binary
with an `LC_RPATH` of `@loader_path/../lib`. Because the library sits at `../lib`
relative to the binary in both the release tarball (`bin/` + `lib/`) and a system
install (`/usr/local/bin` + `/usr/local/lib`), the same binary finds `libfdb_c`
in either layout without `DYLD_LIBRARY_PATH`.
