#!/usr/bin/env bash
set -e

cd "$(dirname "$0")"

if [ ! -f .env ]; then
    echo "No .env file found in $(pwd)" >&2
    exit 1
fi

# Execute each example
examples=(
    create_database
    list_databases
    get_database
    update_database
    delete_database
    create_collection
    list_collections
    get_collection
    update_collection
    delete_collection
    create_document
    list_documents
    get_document
    update_document
    delete_document
    count_documents
    create_attribute
    get_attribute
    delete_attribute
)

for ex in "${examples[@]}"; do
    if [ -d "$ex" ]; then
        echo "Running $ex..."
        go run "./$ex" || exit 1
        echo
    fi
done
