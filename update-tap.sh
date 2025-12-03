#!/bin/bash
# Script to update the Homebrew tap repository with the latest formula

set -e

echo "🔍 Finding Homebrew tap repository..."
TAP_PATH=$(brew --repository abcdofficialzw/netmon 2>/dev/null || echo "")

if [ -z "$TAP_PATH" ]; then
    echo "❌ Tap not found. Make sure you've run: brew tap abcdofficialzw/netmon"
    exit 1
fi

echo "📍 Tap location: $TAP_PATH"
echo ""

echo "📋 Copying updated formula..."
cp "$(dirname "$0")/Formula/netmon.rb" "$TAP_PATH/Formula/netmon.rb"

echo "✅ Formula copied!"
echo ""

cd "$TAP_PATH"
echo "📝 Changes:"
git diff Formula/netmon.rb || echo "No changes detected (formula already up to date)"
echo ""

read -p "Commit and push these changes? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    git add Formula/netmon.rb
    git commit -m "Fix: Add version and SHA256 attributes"
    git push origin main
    echo ""
    echo "✅ Tap repository updated!"
    echo ""
    echo "Now try: brew install abcdofficialzw/netmon/netmon"
else
    echo "Changes saved but not committed. You can commit manually:"
    echo "  cd $TAP_PATH"
    echo "  git add Formula/netmon.rb"
    echo "  git commit -m 'Fix: Add version and SHA256 attributes'"
    echo "  git push origin main"
fi
