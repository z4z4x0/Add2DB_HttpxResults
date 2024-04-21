# Data Integration Tool 🔄

## Overview 📄

This Go-based tool is designed to efficiently process and store web scan results from HTTPx into a SQLite database. It handles large volumes of JSON data, extracting key details and populating a database for further analysis. Perfect for cybersecurity professionals and researchers.

## Features 🌟

- **JSON Processing**: Automatically reads JSON files and extracts essential data.
- **Database Management**: Uses SQLite to manage data, with a schema designed to handle detailed records efficiently.
- **Error Handling**: Robust error handling to ensure stability even with corrupt or malformed data.
- **Performance Optimized**: Fast processing of large datasets.

## Getting Started 🚀

To use this tool, follow these steps:

1. **Clone the repository**:
   ```
   git clone https://github.com/z4z4x0/Add2DB_HttpxResults.git
   ```
2. **Navigate to the project directory**:
   ```
   cd Add2DB_HttpxResults
   ```
3. **Build the project** (Go must be installed):
   ```
   go build
   ```
4. **Run the executable**:
   ```
   ./Add2DB_HttpxResults
   ```

Ensure that the `data` directory contains the JSON files you intend to process and that `data.db` is properly set up in the root directory.

## Schema 📐

The SQLite database is set up with a table `records` that includes numerous fields such as `timestamp`, `url`, `status_code`, etc., designed to capture a comprehensive range of data from each scan result.

## Contributing 🤝

Interested in contributing? Great! Please send pull requests, or file issues for bugs you've noticed or features you'd like to see.

## License ⚖️

Distributed under the GPL-3.0 license.

## Contact 📧

@z4z4_h1

