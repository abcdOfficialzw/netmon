# Homebrew Installation Setup Guide

This guide will help you set up netmon for installation via Homebrew.

## Overview

Homebrew uses "taps" (custom repositories) to install software. We'll create a tap that contains a formula (Ruby file) describing how to build and install netmon.

## Prerequisites

1. A GitHub account
2. A GitHub repository for your netmon releases
3. Homebrew installed on your Mac

## Step-by-Step Instructions

### Step 1: Create a GitHub Release Repository

1. Create a new GitHub repository (e.g., `netmon-releases` or use your existing repo)
2. This will host your release binaries and the Homebrew formula

### Step 2: Create Release Binaries

For Homebrew, you have two options:

#### Option A: Source-based Installation (Recommended)
- Homebrew builds from source
- No need to host binaries
- Works for all architectures automatically

#### Option B: Binary Distribution
- Pre-build binaries for each architecture
- Faster installation
- Requires maintaining multiple binaries

We'll use **Option A** (source-based) as it's simpler and more maintainable.

### Step 3: Create the Homebrew Formula

The formula file is in this repository: `Formula/netmon.rb`

### Step 4: Create a Homebrew Tap Repository

1. Create a new GitHub repository named `homebrew-netmon` (or `homebrew-<yourname>`)
   - The `homebrew-` prefix is required
   - Make it public (Homebrew needs to access it)

2. Clone it locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/homebrew-netmon.git
   cd homebrew-netmon
   ```

3. Copy the formula:
   ```bash
   mkdir -p Formula
   cp /path/to/netmon/Formula/netmon.rb Formula/
   ```

4. Commit and push:
   ```bash
   git add Formula/netmon.rb
   git commit -m "Add netmon formula"
   git push origin main
   ```

### Step 5: Test the Formula Locally

```bash
# Install from your tap
brew install YOUR_USERNAME/netmon/netmon

# Or if using the default tap name:
brew tap YOUR_USERNAME/netmon
brew install netmon
```

### Step 6: Update the Formula for New Versions

When you release a new version:

1. Update the version and SHA256 in `Formula/netmon.rb`
2. Commit and push to your tap repository
3. Users can update with: `brew upgrade netmon`

## Installation Instructions for Users

Once your tap is set up, users can install with:

```bash
# Add your tap
brew tap YOUR_USERNAME/netmon

# Install netmon
brew install netmon

# Run setup
netmon setup
```

## Formula File Location

The formula is located at: `Formula/netmon.rb` in this repository.

## Version Management

To get the SHA256 hash for a new version:

```bash
# For source releases
curl -L https://github.com/YOUR_USERNAME/netmon/archive/v1.0.0.tar.gz | shasum -a 256

# Or use Homebrew's built-in tool
brew fetch --build-from-source YOUR_USERNAME/netmon/netmon
```

## Troubleshooting

### Formula doesn't install
- Check that the repository is public
- Verify the formula syntax: `brew audit --strict Formula/netmon.rb`
- Test locally: `brew install --build-from-source --verbose Formula/netmon.rb`

### Build fails
- Ensure all dependencies are listed in the formula
- Check that Go version requirements are met
- Verify the source URL is accessible

### Users can't find the tap
- Ensure the repository name starts with `homebrew-`
- Make sure the repository is public
- Verify the tap name matches: `YOUR_USERNAME/netmon`

## Next Steps

1. Create your GitHub repositories
2. Copy the formula file to your tap repository
3. Test installation locally
4. Share the installation instructions with users
5. Tag releases in your main repository for version tracking

