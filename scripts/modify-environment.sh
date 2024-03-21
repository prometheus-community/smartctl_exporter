#!/usr/bin/env bash

# This script modifies the environment, namely ~/.bashrc, to preserve the bash history.

# The workspace in the devcontainer is preserved across rebuilds.
# Use .cache to keep history and local scripts.
# .cache is excluded from git (i.e. in .gitignore)
mkdir -p "${PWD}/.cache/"

# Preserve history
[[ ! -L "${HOME}/.bash_history" ]] && ln -sf "${PWD}/.cache/bash_history" "${HOME}/.bash_history"
[[ ! -f "${PWD}/.cache/bash_history" ]] && touch "${PWD}/.cache/bash_history"

# Write history after every command to preserve it across rebuilds.
if ! grep -q '^### CUSTOM: Preserve Bash History ###$' "${HOME}/.bashrc"; then
    cat >> "${HOME}/.bashrc" <<'EOT'
### CUSTOM: Preserve Bash History ###
PROMPT_COMMAND="history -a; ${PROMPT_COMMAND}"
EOT
fi
