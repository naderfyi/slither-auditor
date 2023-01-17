#!/bin/bash

# Check if a file was passed as an argument
if [ -z "$1" ]; then
    echo "Error: No file passed to script."
    exit 1
fi

# Check if a compiler version was passed as an argument
if [ -z "$2" ]; then
    echo "Error: No compiler version passed to script."
    exit 1
fi

# Check if the file exists
if [ ! -f "$1" ]; then
    echo "Error: File $1 not found."
    exit 1
fi

# Check if solc-select is installed
if ! command -v solc-select > /dev/null; then
    echo "Error: solc-select is not installed."
    exit 1
fi

# Check if slither is installed
if ! command -v slither > /dev/null; then
    echo "Error: slither is not installed."
    exit 1
fi

# Use solc-select to set the correct version of the compiler
solc-select use "$2"

# Run the analysis using the correct version of Slither
slither "$1"
