#!/bin/sh

cat > ./docs/cli.md <<EOF
# Command Line Interface

\`\`\`
$(./bin/flowg --help)
\`\`\`

## 1. \`flowg serve\`

\`\`\`
$(./bin/flowg serve --help)
\`\`\`

## 2. \`flowg admin\`

\`\`\`
$(./bin/flowg admin --help)
\`\`\`

### 2.1. \`flowg admin role\`

\`\`\`
$(./bin/flowg admin role --help)
\`\`\`

#### 2.1.1. \`flowg admin role create\`

\`\`\`
$(./bin/flowg admin role create --help)
\`\`\`

#### 2.1.2. \`flowg admin role delete\`

\`\`\`
$(./bin/flowg admin role delete --help)
\`\`\`

#### 2.1.3. \`flowg admin role list\`

\`\`\`
$(./bin/flowg admin role list --help)
\`\`\`

### 2.2. \`flowg admin user\`

\`\`\`
$(./bin/flowg admin user --help)
\`\`\`

#### 2.2.1. \`flowg admin user create\`

\`\`\`
$(./bin/flowg admin user create --help)
\`\`\`

#### 2.2.2. \`flowg admin user delete\`

\`\`\`
$(./bin/flowg admin user delete --help)
\`\`\`

#### 2.2.3. \`flowg admin user list\`

\`\`\`
$(./bin/flowg admin user list --help)
\`\`\`

### 2.3. \`flowg admin token\`

\`\`\`
$(./bin/flowg admin token --help)
\`\`\`

#### 2.3.1. \`flowg admin token create\`

\`\`\`
$(./bin/flowg admin token create --help)
\`\`\`
EOF
