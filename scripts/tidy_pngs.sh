#!/usr/bin/env bash

# This script:
#  1) Takes all .png files in the current directory.
#  2) Extracts only the text after the LAST hyphen ("-") in the filename.
#  3) Removes any leading spaces from that portion.
#  4) Sanitizes any remaining special characters (only alphanumeric, underscore, dash, space).
#  5) Renames them to "XX - extracted_part.png", where XX is a 2-digit number.
#  6) Moves them to "etc/gui/images".
#
# Usage (in Git Bash on Windows):
#   bash tidy_pngs.sh

# Create the target directory if it doesn't exist
mkdir -p etc/gui/images

index=1
for file in *.png; do
    # If there are no .png files, skip
    [ -e "$file" ] || continue

    # Strip the .png extension
    baseName="$(basename "$file" .png)"

    # Extract the portion AFTER the last hyphen (if any), then remove leading spaces
    tailPart="$(echo "$baseName" | sed -E 's/^.*-\s*//')"

    # Replace unwanted characters (anything not alphanumeric, underscore, dash, or space) with underscores
    cleanedName="$(echo "$tailPart" | sed -E 's/[^[:alnum:]_ -]/_/g')"
    
    # Build the new name: 01 - something.png, 02 - somethingElse.png, etc.
    newName="$(printf "%02d - %s.png" "$index" "$cleanedName")"

    # Rename and move to etc/gui/images
    echo "Renaming '$file' to '$newName' and moving to etc/gui/images..."
    mv -- "$file" "etc/gui/images/$newName"

    ((index++))
done