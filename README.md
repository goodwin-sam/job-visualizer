# Job Visualizer 📊

[![Build Status](https://github.com/goodwin-sam/job-visualizer/workflows/Build,%20test%20and%20lint%20job-visualizer/badge.svg)](https://github.com/goodwin-sam/job-visualizer/actions)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20Windows-lightgrey.svg)]()

A desktop application for visualizing and analyzing job data from Excel files. The application runs with a GUI or through command line with a --headless argument.  It has processing capabilities for job data analysis, mapping, and filtering.

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Download](#download)
- [Usage](#usage)
- [Data Format](#data-format)
- [System Requirements](#system-requirements)
- [Troubleshooting](#troubleshooting)



## Features

- **Excel Processing**: Import and analyze job data from `.xlsx` and `.xls` files
- **Interactive GUI**: Desktop interface with job filtering and mapping
- **Headless Mode**: Command-line processing
- **Job Mapping**: Visualize job locations on interactive maps with web browser integration
- **Geocoding & Caching**: Automatic location geocoding with local caching for faster processing
- **Database Storage**: SQLite database for efficient data management

## 🛠 Tech Stack

**Core**
- **Go**
- **Fyne** — Cross-platform GUI framework
- **SQLite** — Lightweight local database

**Data & Integration**
- **Excelize** — Excel file import/export
- **Nominatim** — Address geocoding using OpenStreetMap

**DevOps**
- **GitHub Actions** — Continuous Integration
- **Cross-Platform** — Multi-platform releases (AppImage, .exe, etc.)

## Screenshots

### Geocoded Map View
![Map View w/ Pin Description](screenshots/map-w-pin-desc.png)

### Data Filtering & Analysis
![Filtering Location](screenshots/data-filtering-location.png)

### CLI Mode
![CLI](screenshots/cli.png)

## Quick Start

### 🐧 Linux

1. **Download** the latest AppImage from [Releases](https://github.com/goodwin-sam/job-visualizer/releases)
2. **Make it executable**: `chmod +x job-visualizer-vX.X.X-x86_64.AppImage`
3. **Run**:
   - Double-click the `.AppImage` file
   - Or run from terminal: `./job-visualizer-vX.X.X-x86_64.AppImage` or `./job-visualizer-vX.X.X-x86_64.AppImage --headless`
4. **Process your data**:
   - **If in GUI**: Select your Excel files and output directory, then click "Start Application" to begin processing
   - **If in headless**: The application will automatically process all Excel files in the current directory without any intervention

### 🪟 Windows

1. **Download** the latest `.exe` file from [Releases](https://github.com/goodwin-sam/job-visualizer/releases)
2. **Run**:
   - Double-click the `.exe` file
   - Or run from PowerShell/command prompt: `./job-visualizer-vX.X.X-windows-amd64.exe` or `./job-visualizer-vX.X.X-windows-amd64.exe --headless`
3. **Process your data**:
   - **If in GUI**: Select your Excel files and output directory, then click "Start Application" to begin processing
   - **If in headless**: The application will automatically process all Excel files in the current directory without any intervention

## Download and Run

> **💡 Note**: No installation required - these are portable executables that run directly.

### Option 1: Pre-built Binaries (Recommended)

#### 🐧 Linux - AppImage

1. **Download the AppImage**:
   - Go to the [Releases](https://github.com/goodwin-sam/job-visualizer/releases) page
   - Download the latest `job-visualizer-vX.X.X-x86_64.AppImage` file

2. **Make it executable**:
   ```bash
   chmod +x job-visualizer-vX.X.X-x86_64.AppImage
   ```

3. **Run the application**:
   
   - **Double-click** the AppImage file to launch the GUI
   - **Or run from terminal**:
     ```bash
     # If appimage file saved in current directory
     ./job-visualizer-vX.X.X-x86_64.AppImage
     
     # If saved in PATH (e.g., /usr/local/bin/)
     job-visualizer-vX.X.X-x86_64.AppImage
     ```
   - **Or run headless** (command-line only):
     ```bash
     ./job-visualizer-vX.X.X-x86_64.AppImage --headless
     ```

#### 🪟 Windows - Executable

1. **Download the executable**:
   - Go to the [Releases](https://github.com/goodwin-sam/job-visualizer/releases) page
   - Download the latest `job-visualizer-vX.X.X-windows-amd64.exe` file

2. **Run the application**:
   
   - **Double-click** the `.exe` file to launch the GUI
   - **Or run from PowerShell/command prompt**:
     ```powershell
     ./job-visualizer-vX.X.X-windows-amd64.exe
     ```
   - **Or run headless** (command-line only):
     ```powershell
     ./job-visualizer-vX.X.X-windows-amd64.exe --headless
     ```

### Option 2: Build from Source

If you prefer to build from source, you'll need the following dependencies:

```bash
   # Ubuntu/Debian
   sudo apt install libgl1-mesa-dev xorg-dev
   
   # Fedora
   sudo dnf install mesa-libGL-devel libX11-devel
   
   # Arch Linux
   sudo pacman -S mesa libx11
   ```

Then build and run:

```bash
# GUI
go run ./cmd/app

# Headless
go run ./cmd/app --headless
```

## Usage

### 🖥️ GUI Mode (Default)

1. **Launch** the application
2. **Select files** using the file selection button to choose one or more Excel files (`.xlsx` or `.xls`)
3. **Choose output directory** for the SQLite database with processed data
4. **Click** "Start Application" to begin processing
5. **Explore** the main interface with filtering options, job list, and detailed job information
6. **Filter jobs** using the search options:
   - **Keywords** - Search by job title or company
   - **Location** - Filter by city, state, or country
   - **Minimum salary** - Set salary thresholds
   - **Work-from-home** - Filter remote work options

> **📁 Output Files**: The application creates SQLite database files in your selected output directory that can be opened with any SQLite browser for further analysis.

### ⚡ Headless Mode

Run the application with the `--headless` argument for command-line processing:

```bash
# Process all Excel files in current directory
./job-visualizer-vX.X.X-x86_64.AppImage --headless
```

**What headless mode does:**
- Processes all `.xlsx` and `.xls` files in the current working directory
- Creates a SQLite database table in the current working directory
- Outputs results in a formatted table to the command line

## 📋 Data Format

The application expects Excel files with the following requirements:

### 📊 Worksheet Requirements
- **Worksheet Name**: Must be named "jobs" (case-insensitive - "jobs", "Jobs", "JOBS" all work)
- **Header Row**: First row must contain column headers
- **Data Rows**: All subsequent rows contain job data

### 📐 Column Structure
> **💡 Reference**: See the demo Excel file `demoData.xlsx` for a complete example.

| Column | Header | Description |
|--------|--------|-------------|
| A | Company Name | Name of the hiring company |
| B | Posting Age | How long the job has been posted |
| C | JobId | Unique identifier for the job |
| D | Country | Country where the job is located |
| E | Location | City/State/Region of the job |
| F | Publication Date | When the job was posted |
| G | Salary Max | Maximum salary offered |
| H | Salary Min | Minimum salary offered |
| I | Salary Type | Type of salary (hourly, annual, etc.) |
| J | Job Title | Title of the position |

### 🌍 Geocoding & Caching
- **📍 Automatic Geocoding**: Job locations are automatically converted to coordinates using OpenStreetMap's Nominatim API
- **💾 Intelligent Caching**: Geocoded locations are cached in `~/.job-visualizer/cached_locations.json` to avoid repeated API calls
- **Location Standardization**: Location strings are cleaned and standardized to reduce repeated API calls on the same cities/towns
- **⚡ Performance**: Cached locations load instantly, significantly speeding up subsequent processing runs

## 💻 System Requirements

- **Operating System**: Linux (AppImage), Windows (.exe)
- **Dependencies**: None required for AppImage or .exe (self-contained)

## 🔧 Troubleshooting

### AppImage Issues

If the AppImage doesn't run:

1. **Check permissions**:
   ```bash
   chmod +x job-visualizer-vX.X.X-x86_64.AppImage
   ```

2. **Run from terminal to see errors**:
   ```bash
   ./job-visualizer-vX.X.X-x86_64.AppImage
   ```

3. **Check if your system supports AppImages**:
   ```bash
   # Install AppImageLauncher (Ubuntu/Debian)
   sudo apt install appimagelauncher
   ```

### 🔨 Build Issues

If building from source fails:

1. **Ensure Go is installed** (version 1.19 or later)
2. **Install OpenGL development packages**:
   ```bash
   # Ubuntu/Debian
   sudo apt install libgl1-mesa-dev xorg-dev
   
   # Fedora
   sudo dnf install mesa-libGL-devel libX11-devel
   
   # Arch Linux
   sudo pacman -S mesa libx11
   ```
