# Quick Start: Homebrew Installation

## For End Users

### Installation

```bash
# Add the tap
brew tap YOUR_USERNAME/netmon

# Install netmon
brew install netmon

# Run setup wizard
netmon setup
```

### Usage

```bash
# View network usage (default: apps view)
netmon

# View today's totals
netmon stats today

# View this month
netmon stats month

# View all-time stats
netmon stats all
```

### Update

```bash
brew upgrade netmon
```

### Uninstall

```bash
brew uninstall netmon
```

---

## For Maintainers

### Initial Setup

1. **Create GitHub repositories:**
   - Main repo: `netmon` (your source code)
   - Tap repo: `homebrew-netmon` (the formula)

2. **Set up the tap repository:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/homebrew-netmon.git
   cd homebrew-netmon
   mkdir -p Formula
   # Copy Formula/netmon.rb from this repo
   git add Formula/netmon.rb
   git commit -m "Add netmon formula"
   git push origin main
   ```

3. **Update the formula:**
   - Edit `Formula/netmon.rb`
   - Replace `YOUR_USERNAME` with your GitHub username
   - Update URLs to point to your repository

### Releasing a New Version

1. **Tag a release in your main repository:**
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **Update the formula:**
   ```bash
   cd homebrew-netmon
   # Edit Formula/netmon.rb:
   # - Update version
   # - Update url to point to the tag
   # - Calculate and update sha256
   ```

3. **Get SHA256 hash:**
   ```bash
   curl -L https://github.com/YOUR_USERNAME/netmon/archive/refs/tags/v1.0.0.tar.gz | shasum -a 256
   ```

4. **Test locally:**
   ```bash
   brew install --build-from-source --verbose Formula/netmon.rb
   ```

5. **Commit and push:**
   ```bash
   git add Formula/netmon.rb
   git commit -m "Update netmon to v1.0.0"
   git push origin main
   ```

### Testing the Formula

```bash
# Install from local formula
brew install --build-from-source --verbose Formula/netmon.rb

# Or test from tap
brew tap YOUR_USERNAME/netmon
brew install netmon

# Test the binaries
netmon --help
netmon-service --help
netmon setup
```

### Formula Validation

```bash
# Check formula syntax
brew audit --strict Formula/netmon.rb

# Check for style issues
brew style Formula/netmon.rb
```

---

## Formula Template Customization

Before using the formula, update these values in `Formula/netmon.rb`:

1. **homepage**: Your project homepage URL
2. **url**: Your GitHub repository URL
3. **YOUR_USERNAME**: Replace with your GitHub username
4. **license**: Update if using a different license

Example:
```ruby
homepage "https://github.com/johndoe/netmon"
url "https://github.com/johndoe/netmon/archive/refs/heads/main.tar.gz"
```

---

## Common Issues

### "No available formula"
- Ensure tap repository name starts with `homebrew-`
- Make sure repository is public
- Verify tap name: `brew tap YOUR_USERNAME/netmon`

### Build fails
- Check Go version: `go version` (needs 1.22+)
- Verify all dependencies are available
- Check formula syntax: `brew audit Formula/netmon.rb`

### Binary not found after install
- Check installation: `brew list netmon`
- Verify PATH: `which netmon`
- Reinstall: `brew reinstall netmon`

