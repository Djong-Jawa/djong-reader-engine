#!/bin/bash

# Script to organize resolver files with naming prefixes
# Renames files to include schema prefix (djong-jukung or djong-phinisi)

RESOLVERS_DIR="graph/resolvers"

echo "Organizing resolver files with naming prefixes..."

# Rename djong-jukung resolvers
if [ -f "$RESOLVERS_DIR/mstpricingtiers.resolvers.go" ]; then
    if [ ! -f "$RESOLVERS_DIR/djong-jukung.mstpricingtiers.resolvers.go" ]; then
        echo "Renaming mstpricingtiers.resolvers.go -> djong-jukung.mstpricingtiers.resolvers.go"
        mv "$RESOLVERS_DIR/mstpricingtiers.resolvers.go" "$RESOLVERS_DIR/djong-jukung.mstpricingtiers.resolvers.go"
    else
        echo "Removing duplicate mstpricingtiers.resolvers.go (prefixed version exists)"
        rm "$RESOLVERS_DIR/mstpricingtiers.resolvers.go"
    fi
fi

# Rename djong-phinisi resolvers
for base in mstbooking mstlead mstsalespipeline; do
    if [ -f "$RESOLVERS_DIR/${base}.resolvers.go" ]; then
        if [ ! -f "$RESOLVERS_DIR/djong-phinisi.${base}.resolvers.go" ]; then
            echo "Renaming ${base}.resolvers.go -> djong-phinisi.${base}.resolvers.go"
            mv "$RESOLVERS_DIR/${base}.resolvers.go" "$RESOLVERS_DIR/djong-phinisi.${base}.resolvers.go"
        else
            echo "Removing duplicate ${base}.resolvers.go (prefixed version exists)"
            rm "$RESOLVERS_DIR/${base}.resolvers.go"
        fi
    fi
done

echo "Done! Resolver files organized with naming prefixes."
