#!/bin/bash

# Ensure we're up to date
git fetch origin

# Get the original PR's base commit
echo "Enter the commit hash that the original PR was based on:"
read BASE_COMMIT

# Create a new branch from this base commit
NEW_BRANCH="squashed-pr-$(date +%Y%m%d)"
git checkout -b $NEW_BRANCH $BASE_COMMIT
echo "Created new branch: $NEW_BRANCH based on commit $BASE_COMMIT"

# Get the range of commits to squash
echo "Enter the commit hash of the first commit in the original PR:"
read FIRST_COMMIT
echo "Enter the commit hash of the last commit in the original PR:"
read LAST_COMMIT

# Cherry-pick the range of commits without committing
git cherry-pick -n $FIRST_COMMIT^..$LAST_COMMIT

# Add all changes
git add .

# Commit with a new squash message
echo "Enter a commit message for the squashed commit (press Enter, then Ctrl+D when finished):"
COMMIT_MSG=$(cat)

git commit -m "$COMMIT_MSG"

# Show the log to confirm
git log -n 2
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

echo "Process completed. New branch '$NEW_BRANCH' contains a single squashed commit with all changes from the original PR."
echo "You can now create a new PR from this branch to master on your Git hosting platform."
echo "Note: This branch is based on the original PR's base. If master has progressed, the PR may show as behind master."
