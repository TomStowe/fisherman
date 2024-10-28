#!/bin/sh
# Example pre-commit hook script
echo "Running tests..."
go test ./...
if [ $? -ne 0 ]; then
    echo "Tests failed. Commit aborted."
    exit 1
fi

echo "All checks passed. Proceeding with commit."
