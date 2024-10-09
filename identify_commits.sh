#!/bin/bash

# Ensure we're up to date
git fetch origin

# Update master
git checkout master
git pull origin master

# Function to display commits in a more readable format
display_commits() {
    git log --graph --oneline --decorate -n 20
}

echo "Recent commits on master branch:"
display_commits

echo -e "\nLet's identify the merge commit and the START_COMMIT."

# Find the most recent merge commit
MERGE_COMMIT=$(git log --merges -n 1 --pretty=format:"%H")
MERGE_COMMIT_MSG=$(git log --merges -n 1 --pretty=format:"%s")

echo -e "\nMost recent merge commit:"
echo "Hash: $MERGE_COMMIT"
echo "Message: $MERGE_COMMIT_MSG"

# Display commits before the merge commit
echo -e "\nCommits leading up to the merge commit:"
git log $MERGE_COMMIT^^ --oneline -n 10

echo -e "\nTo identify the START_COMMIT, look for the commit just before your PR's first commit."
echo "This is typically the commit where your feature branch was created from master."

# Prompt user for START_COMMIT
read -p "Enter the hash of the START_COMMIT (the commit just before your PR's first commit): " START_COMMIT

# Verify START_COMMIT
if git cat-file -e $START_COMMIT^{commit} 2>/dev/null; then
    echo "START_COMMIT is valid."
    echo "Commits in your PR:"
    git log --oneline $START_COMMIT..$MERGE_COMMIT
else
    echo "Error: Invalid commit hash. Please run the script again with a valid commit hash."
    exit 1
fi

echo -e "\nSummary:"
echo "Merge Commit: $MERGE_COMMIT"
echo "START_COMMIT: $START_COMMIT"
echo -e "\nYou can use these commit hashes in the revert and squash process."
