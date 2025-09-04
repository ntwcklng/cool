#!/bin/bash

# Pfad zum Quellcode
SRC="$HOME/dev/cool/main.go"

# Zielordner
TARGET="$HOME/go/bin"

# Executable Name
NAME="cool"

# Ordner anlegen, falls nicht vorhanden
mkdir -p "$TARGET"

# Build
echo "Building $NAME..."
go build -o "$TARGET/$NAME" "$SRC"

if [ $? -eq 0 ]; then
  echo "Finished building $NAME!"
else
  echo "Failed building $NAME!"
fi
