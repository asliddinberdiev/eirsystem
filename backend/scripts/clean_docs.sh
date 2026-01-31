#!/bin/bash

# Navigate to docs directory
cd ./docs

# Use sed to replace the prefix in all files
for file in *; do
    if [ -f "$file" ]; then
        sed -i 's/github_com_asliddinberdiev_eirsystem_internal_domain.//g' "$file"
        sed -i 's/github_com_asliddinberdiev_eirsystem_pkg_response.//g' "$file"
    fi
done

echo "Prefix removal completed"