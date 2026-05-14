#!/usr/bin/env node
const fs = require('fs');
const path = require('path');

const husky_dir = path.join(__dirname, '.husky');
const pre_commit = path.join(husky_dir, 'pre-commit');

// Create .husky directory if it doesn't exist
if (!fs.existsSync(husky_dir)) {
  fs.mkdirSync(husky_dir, { recursive: true });
}

// Create pre-commit hook
const pre_commit_content = `#!/bin/sh
. "$(dirname "$0")/_/husky.sh"

npx lint-staged
npx commitlint --edit $1
`;

if (!fs.existsSync(pre_commit)) {
  fs.writeFileSync(pre_commit, pre_commit_content);
  fs.chmodSync(pre_commit, '0755');
}

// Create commit-msg hook
const commit_msg = path.join(husky_dir, 'commit-msg');
const commit_msg_content = `#!/bin/sh
. "$(dirname "$0")/_/husky.sh"

npx --no -- commitlint --edit $1
`;

if (!fs.existsSync(commit_msg)) {
  fs.writeFileSync(commit_msg, commit_msg_content);
  fs.chmodSync(commit_msg, '0755');
}

console.log('✅ Husky hooks installed successfully');
