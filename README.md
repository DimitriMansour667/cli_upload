# cmdim
cmdim is a CLI tool that allows you to manage your PocketBase files.

## Installation

```bash
go install github.com/dimi-mansour/cmdim@latest
```

## Notes

- You need to have a PocketBase Instance running to use this tool
- You need to have a valid record in your PocketBase Instance to use this tool
- You need to have the fields "name", "file", "link", "created" in your record

## Usage

### Config
Set the url of the PocketBase Instance
1. Your URL should point to a PocketBase collection's API endpoint
2. The URL format should be:
   https://your-pocketbase-instance.com/api/collections/YOUR_COLLECTION/records

3. Required collection fields:
   - name
   - file
   - link
   - created
```bash
cmdim config --set <url>
```

Get the url of the PocketBase Instance
```bash
cmdim config --get
```

Get the path of the config file
```bash
cmdim config --path
```

### Check
Check if the PocketBase Instance is running
```bash
cmdim check
```

### Upload
Upload a file to the PocketBase Instance
```bash
cmdim upload <file>
```

### List
List all files in the PocketBase Instance
```bash
cmdim list
```

### Delete
Delete a file from the PocketBase Instance
```bash
cmdim delete <file>
```










