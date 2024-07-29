# Metadata Extractor

**Metadata Extractor** is a powerful command-line tool designed for extracting and managing metadata from various types of files. It supports image files, YAML, JSON, XML files, and even ZIP archives containing multiple files.

## Features

- **EXIF Metadata Extraction**: Extracts detailed EXIF metadata from image files (e.g., JPG, PNG).
- **YAML Parsing**: Reads and parses YAML metadata.
- **JSON Parsing**: Reads and parses JSON metadata.
- **XML Parsing**: Reads and parses XML metadata.
- **ZIP Archive Support**: Extracts metadata from files within ZIP archives.
- **Flexible Output**: Saves extracted metadata to a formatted JSON file.
- **Metadata Filtering**: Selectively extract specific types of metadata or all available metadata.

## Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/Dem0n4ik/metadata-extractor.git
   ```

2. **Navigate to the Project Directory**:
   ```bash
   cd metadata-extractor
   cd metadata_extractor
   ```
   
3. **Initialize Go Modules**

   Initialize a Go module in the root directory of the project:

   ```bash
   go mod init metadata_extractor
   ```

4. **Add Dependencies**

   Fetch the necessary dependencies:

   ```bash
   go get github.com/rwcarlsen/goexif/exif
   go get github.com/rwcarlsen/goexif/tiff
   go get gopkg.in/yaml.v2
   ```

5. **Build the Project**:
   ```bash
   go build -o metadata_extractor
   ```

   This will create an executable named `metadata_extractor` in the project directory.

## Usage

The tool accepts several command-line flags to control its behavior:

### Command-Line Flags

- `-files` **(required)**: A comma-separated list of file paths. This can include individual files or ZIP archives.
- `-output` **(optional)**: The path to the output JSON file where metadata will be saved. If omitted, metadata will not be saved to a file.
- `-log` **(optional)**: The path to the log file. The default is `app.log`.
- `-type` **(optional)**: Specifies the type of metadata to extract. Options include `exif`, `yaml`, `json`, `xml`, or `all`. The default is `all`.

### Examples

1. **Extract Metadata from Individual Files**:
   ```bash
   ./metadata_extractor -files image1.jpg,config.yaml -output metadata.json -type all
   ```
   This command extracts metadata from `image1.jpg` and `config.yaml`, saving the results to `metadata.json`.

2. **Extract Metadata from ZIP Archive**:
   ```bash
   ./metadata_extractor -files archive.zip -output archive_metadata.json -type all
   ```
   This command extracts metadata from all files within `archive.zip` and saves it to `archive_metadata.json`.

3. **Extract Only JSON Metadata**:
   ```bash
   ./metadata_extractor -files data.json,info.xml -output json_metadata.json -type json
   ```
   This command extracts metadata from `data.json` and `info.xml`, but only the JSON metadata, saving it to `json_metadata.json`.

## Metadata File Format

The output JSON file is formatted as a list of metadata objects. Each object contains the following fields:

- **`Filename`**: The name of the file from which metadata was extracted.
- **`Type`**: The type of metadata (e.g., EXIF, YAML, JSON, XML).
- **`Data`**: The extracted metadata in a key-value pair format.

### JSON Metadata File Structure

Example of JSON output for an image file with EXIF metadata:

```json
[
  {
    "Filename": "image1.jpg",
    "Type": "EXIF",
    "Data": {
      "ImageWidth": 1920,
      "ImageHeight": 1080,
      "Make": "Canon",
      "Model": "EOS 5D Mark IV",
      "DateTime": "2023:07:25 14:22:30",
      "ExposureTime": "1/125",
      "FNumber": "2.8",
      ...
    }
  }
]
```

Example of JSON output for a YAML file:

```json
[
  {
    "Filename": "config.yaml",
    "Type": "YAML",
    "Data": {
      "version": "1.0",
      "settings": {
        "theme": "dark",
        "language": "en"
      },
      "features": [
        "notifications",
        "auto-updates"
      ]
    }
  }
]
```

## Contribution

Contributions to the project are welcome. You can help by:

- Reporting bugs or issues.
- Suggesting features or improvements.
- Submitting pull requests with enhancements or fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
