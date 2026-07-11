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
