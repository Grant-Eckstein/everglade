# Everglade
*A file encryption framework*

#### Overview
So, this framework has a few main features, (tbh this should be 3 projects), namely:
1. File encryption
2. File management
3. Key management

###### File Encryption
File encryption in Everglade is the main concern. The idea here is to make secure management of files super easy. This is done by file discovery and support encryption functions. 

###### Example Usage
```go
//  Discovers all files in current directory and then encrypts and decrypts them!
files := DiscoverFilesInDirectory(".")
for _, f := range files.files {
    f.encrypt()
}
for _, f := range files.files {
    f.decrypt()
}
```







