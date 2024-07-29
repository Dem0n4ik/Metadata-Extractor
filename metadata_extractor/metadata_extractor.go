package main

import (
    "archive/zip"
    "encoding/json"
    "encoding/xml"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"

    "github.com/rwcarlsen/goexif/exif"
    "github.com/rwcarlsen/goexif/tiff"
    "gopkg.in/yaml.v2"
)

// Metadata structure
type Metadata struct {
    Filename string
    Type     string
    Data     interface{}
}

// ExifData structure
type ExifData struct {
    Fields map[string]interface{}
}

// YamlData structure
type YamlData struct {
    Fields map[string]interface{}
}

// JsonData structure
type JsonData struct {
    Fields map[string]interface{}
}

// XmlData structure
type XmlData struct {
    Fields map[string]interface{}
}

// Extract EXIF data from an image
func extractExif(filename string) (ExifData, error) {
    log.Printf("Extracting EXIF data from %s\n", filename)
    file, err := os.Open(filename)
    if err != nil {
        return ExifData{}, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    x, err := exif.Decode(file)
    if err != nil {
        return ExifData{}, fmt.Errorf("failed to decode exif data: %w", err)
    }

    data := make(map[string]interface{})
    x.Walk(&exifWalker{data})

    return ExifData{Fields: data}, nil
}

// ExifWalker implements the exif.Walker interface
type exifWalker struct {
    data map[string]interface{}
}

func (w *exifWalker) Walk(name exif.FieldName, tag *tiff.Tag) error {
    w.data[string(name)] = tag.String()
    return nil
}

// Extract YAML metadata from a file
func extractYaml(filename string) (YamlData, error) {
    log.Printf("Extracting YAML data from %s\n", filename)
    file, err := os.Open(filename)
    if err != nil {
        return YamlData{}, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    var data map[string]interface{}
    decoder := yaml.NewDecoder(file)
    err = decoder.Decode(&data)
    if err != nil {
        return YamlData{}, fmt.Errorf("failed to decode yaml data: %w", err)
    }

    return YamlData{Fields: data}, nil
}

// Extract JSON metadata from a file
func extractJson(filename string) (JsonData, error) {
    log.Printf("Extracting JSON data from %s\n", filename)
    file, err := os.Open(filename)
    if err != nil {
        return JsonData{}, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    var data map[string]interface{}
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&data)
    if err != nil {
        return JsonData{}, fmt.Errorf("failed to decode json data: %w", err)
    }

    return JsonData{Fields: data}, nil
}

// Extract XML metadata from a file
func extractXml(filename string) (XmlData, error) {
    log.Printf("Extracting XML data from %s\n", filename)
    file, err := os.Open(filename)
    if err != nil {
        return XmlData{}, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    var data map[string]interface{}
    decoder := xml.NewDecoder(file)
    err = decoder.Decode(&data)
    if err != nil {
        return XmlData{}, fmt.Errorf("failed to decode xml data: %w", err)
    }

    return XmlData{Fields: data}, nil
}

// Extract metadata from a file based on its extension
func extractMetadata(filename, ext string) (Metadata, error) {
    log.Printf("Extracting metadata from %s\n", filename)
    var metadata Metadata
    metadata.Filename = filename

    switch ext {
    case "jpg", "jpeg", "png", "tiff", "bmp":
        exifData, err := extractExif(filename)
        if err != nil {
            return metadata, fmt.Errorf("error extracting EXIF data: %w", err)
        }
        metadata.Type = "EXIF"
        metadata.Data = exifData

    case "yaml", "yml":
        yamlData, err := extractYaml(filename)
        if err != nil {
            return metadata, fmt.Errorf("error extracting YAML metadata: %w", err)
        }
        metadata.Type = "YAML"
        metadata.Data = yamlData

    case "json":
        jsonData, err := extractJson(filename)
        if err != nil {
            return metadata, fmt.Errorf("error extracting JSON metadata: %w", err)
        }
        metadata.Type = "JSON"
        metadata.Data = jsonData

    case "xml":
        xmlData, err := extractXml(filename)
        if err != nil {
            return metadata, fmt.Errorf("error extracting XML metadata: %w", err)
        }
        metadata.Type = "XML"
        metadata.Data = xmlData

    default:
        return metadata, fmt.Errorf("unsupported file type: %s", ext)
    }

    return metadata, nil
}

// Save metadata to a JSON file
func saveMetadata(metadata []Metadata, outputFile string) error {
    log.Printf("Saving metadata to %s\n", outputFile)
    file, err := os.Create(outputFile)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(metadata); err != nil {
        return fmt.Errorf("failed to encode metadata: %w", err)
    }

    return nil
}

// Print metadata
func printMetadata(metadata Metadata) {
    fmt.Printf("Metadata for %s (%s):\n", metadata.Filename, metadata.Type)
    switch data := metadata.Data.(type) {
    case ExifData:
        for key, value := range data.Fields {
            fmt.Printf("%s: %v\n", key, value)
        }
    case YamlData:
        for key, value := range data.Fields {
            fmt.Printf("%s: %v\n", key, value)
        }
    case JsonData:
        for key, value := range data.Fields {
            fmt.Printf("%s: %v\n", key, value)
        }
    case XmlData:
        for key, value := range data.Fields {
            fmt.Printf("%s: %v\n", key, value)
        }
    }
}

// Extract metadata from files within a ZIP archive
func processZipFile(zipFilename string, metadataType string) ([]Metadata, error) {
    log.Printf("Processing ZIP archive %s\n", zipFilename)
    var allMetadata []Metadata

    zipFile, err := zip.OpenReader(zipFilename)
    if err != nil {
        return nil, fmt.Errorf("failed to open ZIP archive: %w", err)
    }
    defer zipFile.Close()

    for _, file := range zipFile.File {
        if file.FileInfo().IsDir() {
            continue
        }

        ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Name), "."))
        if metadataType != "all" && metadataType != ext {
            continue
        }

        f, err := file.Open()
        if err != nil {
            log.Printf("Error opening file %s in ZIP: %v\n", file.Name, err)
            continue
        }
        defer f.Close()

        tmpFile, err := os.CreateTemp("", "tmpfile-")
        if err != nil {
            log.Printf("Error creating temp file for %s: %v\n", file.Name, err)
            continue
        }
        defer os.Remove(tmpFile.Name())

        _, err = io.Copy(tmpFile, f)
        if err != nil {
            log.Printf("Error copying file %s to temp file: %v\n", file.Name, err)
            continue
        }
        tmpFile.Close()

        metadata, err := extractMetadata(tmpFile.Name(), ext)
        if err != nil {
            log.Printf("Error extracting metadata from %s: %v\n", file.Name, err)
            continue
        }

        allMetadata = append(allMetadata, metadata)
    }

    return allMetadata, nil
}

func main() {
    filenames := flag.String("files", "", "Comma-separated list of filenames to extract metadata from")
    outputFile := flag.String("output", "", "Output file to save metadata (optional)")
    logFile := flag.String("log", "app.log", "Log file to save logs")
    metadataType := flag.String("type", "all", "Type of metadata to extract (exif, yaml, json, xml, all)")
    flag.Parse()

    if *filenames == "" {
        log.Fatal("Please specify at least one filename using the -files flag")
    }

    // Set up logging
    logOutput, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }
    defer logOutput.Close()
    log.SetOutput(logOutput)

    files := strings.Split(*filenames, ",")
    var allMetadata []Metadata

    for _, filename := range files {
        if _, err := os.Stat(filename); os.IsNotExist(err) {
            log.Printf("File does not exist: %s\n", filename)
            continue
        }

        ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
        if ext == "zip" {
            metadata, err := processZipFile(filename, *metadataType)
            if err != nil {
                log.Printf("Error processing ZIP file %s: %v\n", filename, err)
                continue
            }
            allMetadata = append(allMetadata, metadata...)
        } else {
            metadata, err := extractMetadata(filename, ext)
            if err != nil {
                log.Printf("Error processing file %s: %v\n", filename, err)
                continue
            }
            if *metadataType == "all" || strings.ToLower(metadata.Type) == strings.ToLower(*metadataType) {
                printMetadata(metadata)
                allMetadata = append(allMetadata, metadata)
            }
        }
    }

    if *outputFile != "" {
        err := saveMetadata(allMetadata, *outputFile)
        if err != nil {
            log.Fatalf("Error saving metadata to file: %v", err)
        }
        fmt.Printf("All metadata saved to %s\n", *outputFile)
    }

    log.Println("Process completed successfully")
}

