# Homebrew Installation - Step by Step Guide

## Overview

This guide will walk you through setting up netmon for Homebrew installation. Homebrew uses "taps" (custom repositories) to distribute software that isn't in the main Homebrew repository.

## Prerequisites

- ✅ GitHub account
- ✅ Homebrew installed on your Mac
- ✅ Git installed

## Step 1: Create GitHub Repositories

You need **two** repositories:

### Repository 1: Main Source Code Repository
- **Name**: `netmon` (or your preferred name)
- **Purpose**: Your source code
- **Visibility**: Public (required for Homebrew)

### Repository 2: Homebrew Tap Repository
- **Name**: `homebrew-netmon` (must start with `homebrew-`)
- **Purpose**: Contains the Homebrew formula
- **Visibility**: Public (required for Homebrew)

**Action**: Create both repositories on GitHub now.

---

## Step 2: Push Your Source Code

If you haven't already:

```bash
cd /Users/titus/projects/veTech/usage

# Initialize git if needed
git init
git add .
git commit -m "Initial commit"

# Add your GitHub remote
git remote add origin https://github.com/YOUR_USERNAME/netmon.git
git branch -M main
git push -u origin main

```

**Replace `YOUR_USERNAME` with your GitHub username.**

---

## Step 3: Update the Formula File

Edit `Formula/netmon.rb` and replace all instances of `YOUR_USERNAME`:

```bash
# Open the formula
open Formula/netmon.rb

# Or edit manually:
# Replace: https://github.com/YOUR_USERNAME/netmon
# With:    https://github.com/your-actual-username/netmon
```

**Key fields to update:**
- `homepage`: Your repository URL
- `url`: Your repository URL (for source downloads)
- `head`: Your repository URL (for development builds)

---

## Step 4: Create the Homebrew Tap Repository

```bash
# Create a new directory
mkdir ~/homebrew-netmon
cd ~/homebrew-netmon

# Initialize git
git init
git branch -M main

# Create Formula directory
mkdir -p Formula

# Copy the formula
cp /Users/titus/projects/veTech/usage/Formula/netmon.rb Formula/

# Commit
git add Formula/netmon.rb
git commit -m "Add netmon formula"

# Add remote (replace YOUR_USERNAME)
git remote add origin https://github.com/YOUR_USERNAME/homebrew-netmon.git

# Push
git push -u origin main
```

---

## Step 5: Test the Formula Locally

Before sharing, test it works:

```bash
# Install from your local tap
brew install --build-from-source --verbose ~/homebrew-netmon/Formula/netmon.rb
```

Or test from the remote tap:

```bash
# Add your tap
brew tap YOUR_USERNAME/netmon

# Install
brew install netmon

# Test it works
netmon --help
netmon-service --help
```

---

## Step 6: Create a Versioned Release (Optional but Recommended)

For stable releases, use version tags:

### 6a. Tag a Release

```bash
cd /Users/titus/projects/veTech/usage

# Create a version tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### 6b. Update Formula for Versioned Release

Edit `Formula/netmon.rb` in your tap repository:

```ruby
class Netmon < Formula
  desc "macOS Network Usage Monitor - Track network usage by interface and application"
  homepage "https://github.com/YOUR_USERNAME/netmon"
  url "https://github.com/YOUR_USERNAME/netmon/archive/refs/tags/v1.0.0.tar.gz"
  version "1.0.0"
  sha256 "CALCULATE_THIS_BELOW"
  license "MIT"
  # ... rest of formula
end
```

### 6c. Calculate SHA256 Hash

```bash
curl -L https://github.com/YOUR_USERNAME/netmon/archive/refs/tags/v1.0.0.tar.gz | shasum -a 256
```

Copy the hash and paste it into the formula.

### 6d. Update Tap Repository

```bash
cd ~/homebrew-netmon
# Edit Formula/netmon.rb with the new version and SHA256
git add Formula/netmon.rb
git commit -m "Update netmon to v1.0.0"
git push origin main
```

---

## Step 7: Share Installation Instructions

Users can now install with:

```bash
# Add your tap
brew tap YOUR_USERNAME/netmon

# Install netmon
brew install netmon

# Run setup
netmon setup
```

Add this to your README.md:

```markdown
## Installation via Homebrew

```bash
brew tap YOUR_USERNAME/netmon
brew install netmon
netmon setup
```
```

---

## Step 8: Updating for New Versions

When you release a new version:

1. **Tag the release** in your main repository:
   ```bash
   git tag -a v1.1.0 -m "Release v1.1.0"
   git push origin v1.1.0
   ```

2. **Calculate new SHA256**:
   ```bash
   curl -L https://github.com/YOUR_USERNAME/netmon/archive/refs/tags/v1.1.0.tar.gz | shasum -a 256
   ```

3. **Update the formula** in your tap repository:
   - Update `version`
   - Update `url` to point to new tag
   - Update `sha256` with new hash

4. **Commit and push**:
   ```bash
   cd ~/homebrew-netmon
   git add Formula/netmon.rb
   git commit -m "Update netmon to v1.1.0"
   git push origin main
   ```

5. **Users update** with:
   ```bash
   brew upgrade netmon
   ```

---

## Troubleshooting

### Formula doesn't install
- ✅ Check repository is public
- ✅ Verify formula syntax
- ✅ Ensure tap name is correct: `YOUR_USERNAME/netmon`

### Build fails
- ✅ Check Go version (needs 1.22+)
- ✅ Verify all dependencies in go.mod
- ✅ Check source URL is accessible

### "No available formula"
- ✅ Ensure tap repository name starts with `homebrew-`
- ✅ Make sure repository is public
- ✅ Verify you've pushed the formula file

### Test formula syntax
```bash
# Check for syntax errors
ruby -c Formula/netmon.rb

# Check Homebrew style (if brew audit works)
brew style Formula/netmon.rb
```

---

## Quick Reference

**Your repositories:**
- Main: `https://github.com/YOUR_USERNAME/netmon`
- Tap: `https://github.com/YOUR_USERNAME/homebrew-netmon`

**User installation:**
```bash
brew tap YOUR_USERNAME/netmon
brew install netmon
```

**Update formula:**
```bash
cd ~/homebrew-netmon
# Edit Formula/netmon.rb
git add Formula/netmon.rb
git commit -m "Update to vX.X.X"
git push
```

---

## Next Steps

1. ✅ Create both GitHub repositories
2. ✅ Push your source code
3. ✅ Update formula with your username
4. ✅ Create tap repository with formula
5. ✅ Test installation locally
6. ✅ Create first versioned release
7. ✅ Share installation instructions

You're all set! 🎉

