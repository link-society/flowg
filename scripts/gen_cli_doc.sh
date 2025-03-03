#!/bin/sh

cat > ./docs/cli.md <<EOF
# Command Line Interface

\`\`\`
$(./bin/flowg-server --help)
\`\`\`
EOF
