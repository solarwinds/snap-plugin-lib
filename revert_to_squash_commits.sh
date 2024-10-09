#!/bin/bash

# Ensure we're up to date
git fetch origin

# Update master
git checkout master
git pull origin master

# Get the merge commit hash
echo "Enter the hash of the merge commit you want to revert:"
read MERGE_COMMIT

# Get the commit hash to start from (the commit just before your PR commits started)
echo "Enter the commit hash just before your PR commits started:"
read START_COMMIT

# Create a new branch
NEW_BRANCH="squashed-pr-$(date +%Y%m%d)"
git checkout -b $NEW_BRANCH master
echo "Created new branch: $NEW_BRANCH based on current master"

# Revert the merge commit
git revert -m 1 $MERGE_COMMIT

# Check if there are any commits to revert
COMMIT_COUNT=$(git rev-list --count $START_COMMIT..$MERGE_COMMIT)
if [ $COMMIT_COUNT -eq 0 ]; then
    echo "Error: No commits found between $START_COMMIT and $MERGE_COMMIT."
    echo "Here are the last 5 commits before the merge:"
    git log $MERGE_COMMIT^^ -n 5 --oneline
    echo "Please check the START_COMMIT and try again."
    exit 1
fi

echo "Found $COMMIT_COUNT commits to revert."

# Perform the revert of PR commits
git revert --no-commit $START_COMMIT..$MERGE_COMMIT

# Check for conflicts
if [ -n "$(git status --porcelain)" ]; then
    echo "Conflicts detected. Please resolve them manually, then run 'git revert --continue'."
    echo "After resolving all conflicts, stage the changes with 'git add' and run this script again."
    exit 1
fi

# Commit the squashed changes
echo "Enter a commit message for the squashed commit (press Enter, then Ctrl+D when finished):"
COMMIT_MSG=$(cat)

git commit -m "$COMMIT_MSG"

# Show the log to confirm
git log -n 3
echo "Above are the latest commits. Please verify they look correct."

# Offer to push the changes
read -p "Do you want to push this new branch to origin? (y/n) " PUSH_CONFIRM
if [[ $PUSH_CONFIRM =~ ^[Yy]$ ]]
then
    git push -u origin $NEW_BRANCH
    echo "Pushed $NEW_BRANCH to origin"
else
    echo "Branch not pushed. You can push later with: git push -u origin $NEW_BRANCH"
fi

echo "Process completed. New branch '$NEW_BRANCH' contains two commits:"
echo "1. A revert of the merge commit"
echo "2. A single squashed commit with all changes from the original PR reverted"
echo "You can now create a new PR from this branch to master on your Git hosting platform."