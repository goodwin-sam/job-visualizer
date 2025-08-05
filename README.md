# Job Visualizer

A desktop application for visualizing and analyzing job data from Excel files. The application runs with a GUI or through command line with a --headless argument.  It has processing capabilities for job data analysis, mapping, and filtering.

## Features

- **Excel Processing**: Import and analyze job data from `.xlsx` and `.xls` files
- **Interactive GUI**: Desktop interface with job filtering and mapping
- **Headless Mode**: Command-line processing
- **Job Mapping**: Visualize job locations on interactive maps
- **Database Storage**: SQLite database for efficient data management

## Download and Installation

### Option 1: AppImage (Recommended)

1. **Download the AppImage**:
   - Go to the [Releases](https://github.com/samg111/job-visualizer/releases) page
   - Download the latest `Job_Visualizer-x86_64.AppImage` file

2. **Make it executable**:
   ```bash
   chmod +x Job_Visualizer-x86_64.AppImage
   ```

3. **Run the application**:
   ```bash
   # If appimage file saved in current directory
   ./Job_Visualizer-x86_64.AppImage
   
   # If saved in PATH (e.g., /usr/local/bin/)
   Job_Visualizer-x86_64.AppImage
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

### GUI Mode (Default)

1. Launch the application
2. Use the file selection button to choose one or more Excel files (`.xlsx` or `.xls`)
3. Select an output directory for processed data
4. Click "Start Application" to begin processing
5. Once processing is complete, you'll see the main interface with job data
6. Use the filtering options to search and filter jobs by:
   - Keywords
   - Location
   - Minimum salary
   - Work-from-home options

**Output Files**: The application creates SQLite database files in your selected output directory
- These database files can be opened with any SQLite browser for further analysis

### Headless Mode

Run the application with argument --headless for headless processing:

```bash
# Process all Excel files in current directory
./Job_Visualizer-x86_64.AppImage --headless
```

The headless mode will:
- Process all `.xlsx` and `.xls` files in the current directory
- Display job information in a sqlite database table
- Output results in a formatted table to the command line

## Data Format

The application expects Excel files with the following columns structure and a header line. Please see demo excel file demoData.xlsx:
- A1: Company Name
- B1: Posting Age
- C1: JobId
- D1: Country
- E1: Location
- F1: Publication Date
- G1: Salary Max
- H1: Salary Min
- I1: Salary Type
- J1: Job Title

## System Requirements

- **Operating System**: Linux (AppImage)
- **Dependencies**: None required for AppImage (self-contained)

## Troubleshooting

### AppImage Issues

If the AppImage doesn't run:

1. **Check permissions**:
   ```bash
   chmod +x Job_Visualizer-x86_64.AppImage
   ```

2. **Run from terminal to see errors**:
   ```bash
   ./Job_Visualizer-x86_64.AppImage
   ```

3. **Check if your system supports AppImages**:
   ```bash
   # Install AppImageLauncher (Ubuntu/Debian)
   sudo apt install appimagelauncher
   ```

### Build Issues

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
