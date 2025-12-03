# How to Update netmon

## For End Users

### Updating via Homebrew

If you installed netmon via Homebrew, updating is simple:

```bash
# Update netmon to the latest version
brew upgrade netmon
```

This will:
- Fetch the latest formula from the tap
- Download and install the new version
- Keep your existing database and configuration

### Check Current Version

```bash
# See what version you have
netmon version

# Or check via Homebrew
brew list --versions netmon
```

### Update Homebrew Tap First (if needed)

If you're not getting the latest version, update the tap:

```bash
# Update the tap repository
brew update

# Then upgrade netmon
brew upgrade netmon
```

### Troubleshooting Updates

**If upgrade says "already installed":**
```bash
# Force reinstall
brew reinstall netmon
```

**If you get "No available formula":**
```bash
# Make sure tap is added
brew tap abcdofficialzw/netmon

# Update tap
brew update

# Try upgrade again
brew upgrade netmon
```

**If you want to install a specific version:**
```bash
# Uninstall current version
brew uninstall netmon

# Install specific version (if available)
brew install abcdofficialzw/netmon/netmon@0.1.0
```

---

## For Maintainers

### How to Release a New Version

When you want to release a new version:

#### 1. Update Version in Code

Edit `cmd/netmon/main.go`:
```go
const Version = "0.2.0"  // Update this
```

#### 2. Create a Git Tag

```bash
# Commit your changes
git add .
git commit -m "Release v0.2.0"

# Create and push tag
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

#### 3. Update Homebrew Formula

Edit `Formula/netmon.rb` in your tap repository:

```ruby
class Netmon < Formula
  # ... other fields ...
  url "https://github.com/abcdOfficialzw/netmon/archive/refs/tags/v0.2.0.tar.gz"
  version "0.2.0"
  sha256 "CALCULATE_THIS"
  # ... rest of formula ...
end
```

#### 4. Calculate SHA256

```bash
curl -sL https://github.com/abcdOfficialzw/netmon/archive/refs/tags/v0.2.0.tar.gz | shasum -a 256
```

Copy the hash and update it in the formula.

#### 5. Update Tap Repository

```bash
cd ~/homebrew-netmon  # or wherever your tap is
# Edit Formula/netmon.rb with new version and SHA256
git add Formula/netmon.rb
git commit -m "Update netmon to v0.2.0"
git push origin main
```

#### 6. Test the Update

```bash
# Test locally
brew install --build-from-source Formula/netmon.rb

# Or test from tap
brew upgrade netmon
```

### Automated Update Script

You can create a script to automate this process. Here's a template:

```bash
#!/bin/bash
# update-version.sh

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: ./update-version.sh <version>"
    echo "Example: ./update-version.sh 0.2.0"
    exit 1
fi

# Update version in code
sed -i '' "s/const Version = \".*\"/const Version = \"$VERSION\"/" cmd/netmon/main.go

# Create tag
git add cmd/netmon/main.go
git commit -m "Release v$VERSION"
git tag -a "v$VERSION" -m "Release v$VERSION"
git push origin main
git push origin "v$VERSION"

# Calculate SHA256
SHA256=$(curl -sL "https://github.com/abcdOfficialzw/netmon/archive/refs/tags/v$VERSION.tar.gz" | shasum -a 256 | awk '{print $1}')

echo ""
echo "Update Formula/netmon.rb in your tap repository:"
echo "  url \"https://github.com/abcdOfficialzw/netmon/archive/refs/tags/v$VERSION.tar.gz\""
echo "  version \"$VERSION\""
echo "  sha256 \"$SHA256\""
echo ""
echo "Then commit and push to your tap repository."
```

---

## Version Numbering

Follow semantic versioning (semver):
- **MAJOR.MINOR.PATCH** (e.g., 1.2.3)
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

Examples:
- `0.1.0` → `0.1.1` (bug fix)
- `0.1.0` → `0.2.0` (new feature)
- `0.1.0` → `1.0.0` (major release/breaking changes)

---

## Best Practices

1. **Always tag releases** - Makes it easy to track versions
2. **Update formula promptly** - Users expect timely updates
3. **Test before releasing** - Verify the formula works
4. **Changelog** - Consider maintaining a CHANGELOG.md
5. **Announce updates** - Let users know about new versions

---

## Quick Reference

**Users:**
```bash
brew upgrade netmon          # Update to latest
netmon version               # Check current version
```

**Maintainers:**
```bash
# 1. Update version in code
# 2. git tag -a vX.X.X -m "Release vX.X.X"
# 3. git push origin vX.X.X
# 4. Update Formula/netmon.rb
# 5. Push to tap repository
```

