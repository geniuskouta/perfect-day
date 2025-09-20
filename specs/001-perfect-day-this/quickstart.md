# Perfect Day Terminal App - Quickstart Guide

## Installation

1. Download the binary for your platform from GitHub releases
2. Make it executable: `chmod +x perfectday`
3. Move to PATH: `sudo mv perfectday /usr/local/bin/`
4. Verify installation: `perfectday --version`

## Initial Setup

### 1. Initialize Your Profile
```bash
perfectday init kouta
```

This creates:
- Configuration at `~/.perfectday/config.json`
- Data directory at `~/.perfectday/days/kouta/`
- User profile in `~/.perfectday/users.json`

### 2. Configure Google Places API (Optional)
Edit `~/.perfectday/config.json` to add your API key:
```json
{
  "google_places_api_key": "your-api-key-here"
}
```

Without an API key, you can still use custom text for locations.

## Creating Your First Perfect Day

### 1. Start the Creation Process
```bash
perfectday create
```

### 2. Follow the Interactive Prompts

**Day Title**: Enter a memorable title (5-100 characters)
```
What was this perfect day? Morning in Shibuya
```

**Day Description**: Describe the overall experience (10-1000 characters)
```
A perfect morning exploring the bustling heart of Tokyo, from quiet coffee to vibrant streets.
```

**Date**: Enter the date (defaults to today)
```
What date was this? (YYYY-MM-DD) [2025-01-15]: 2025-01-15
```

**Main Areas**: List the key areas you spent time in
```
What areas did you spend time in? (comma-separated): Shibuya, Harajuku
```

### 3. Add Activities

For each activity, you'll be prompted for:

**Location** (with Google Places search if API key configured):
```
Location name: Blue Bottle Coffee
🔍 Searching Google Places...

1. Blue Bottle Coffee Shibuya - Tokyo, Shibuya City, Shibuya, 2 Chome−24−12
2. Blue Bottle Coffee Roppongi - Tokyo, Minato City, Roppongi, 6 Chome−10−1

Select (1-2) or press Enter for custom text: 1
```

**Timing**:
```
Start time (HH:MM): 08:00
End time (HH:MM): 09:30
```

**Personal Commentary**:
```
Your thoughts about this activity: Perfect quiet spot to start the day. The flat white was exceptional and the morning light through the windows was beautiful.
```

**Continue or Finish**:
```
Add another activity? (y/n): y
```

### 4. Save and Review
After adding all activities, your perfect day is automatically saved and you'll see a summary.

## Viewing and Exploring Perfect Days

### List All Perfect Days
```bash
perfectday list
```

### Search by Area or Activity
```bash
perfectday search "coffee"
perfectday search --area "Shibuya" "morning"
```

### View Detailed Timeline
```bash
perfectday view "Morning in Shibuya"
```

### Export Data
```bash
perfectday list --format=json > my-perfect-days.json
```

## Managing Your Perfect Days

### Edit an Existing Day
```bash
perfectday edit "Morning in Shibuya"
```

### Delete a Day (Soft Delete)
```bash
perfectday delete "Morning in Shibuya"
```

## Sharing with Others

All perfect days are publicly viewable by design. Share your data file or export specific days:

```bash
# Export all days as JSON
perfectday list --format=json

# Export specific user's days
perfectday list --user=kouta --format=json

# Share your data directory
tar -czf my-perfect-days.tar.gz ~/.perfectday/days/
```

## Example Perfect Day Workflow

Here's the complete flow for creating Kouta's perfect day:

```bash
# 1. Initialize (one time setup)
perfectday init kouta

# 2. Create the perfect day
perfectday create
# Follow prompts to create "Morning in Shibuya" with 2 activities

# 3. View the result
perfectday view "Morning in Shibuya"

# 4. List all your days
perfectday list --user=kouta

# 5. Search for coffee-related days
perfectday search "coffee"

# 6. Export for sharing
perfectday list --format=json > share-with-friends.json
```

## Troubleshooting

### Common Issues

**"No perfect days found"**:
- Check if you're in the right user directory
- Verify files exist in `~/.perfectday/days/your-username/`

**"Google Places search failed"**:
- Verify your API key in `~/.perfectday/config.json`
- Check internet connection
- Use custom text entry as fallback

**"Permission denied"**:
- Ensure `~/.perfectday/` directory is writable
- Check file permissions: `ls -la ~/.perfectday/`

**"Invalid time format"**:
- Use 24-hour format: HH:MM (e.g., 09:30, 14:15)
- Ensure end time is after start time

### Getting Help

```bash
# General help
perfectday --help

# Command-specific help
perfectday create --help
perfectday search --help
```

### File Locations

- **Config**: `~/.perfectday/config.json`
- **User data**: `~/.perfectday/users.json`
- **Perfect days**: `~/.perfectday/days/username/YYYY-MM-DD_title.json`
- **Logs**: `~/.perfectday/logs/` (if debugging enabled)

## Validation Test Scenarios

This quickstart serves as the basis for integration tests. Each section above should be executable as automated test scenarios to verify the complete user workflow functions correctly.