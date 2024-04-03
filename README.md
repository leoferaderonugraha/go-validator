# Validator for Typed Struct in Golang (WIP)

### Example
```go
package main

import (
    "fmt"
    "github.com/leoferaderonugraha/go-validator"
)

type TemporaryFileUploadRequest struct {
    Email     string  `json:"email" validate:"required|email"`
    FileData  string  `json:"file_data" validate:"required|ext:png,jpg"`
    FileName  string  `json:"file_name" validate:"required"`
}

func main() {
    request := TemporaryFileUploadRequest{
        "leo@github.com",
        // an encoded 1x1 pixel png data
        "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVQIW2P4v5ThPwAG7wKklwQ/bwAAAABJRU5ErkJggg==",
        "avatar.png",
    }

    errs, ok := validator.Validate(request)

    if !ok {
        for i := 0; i < len(errs); i++ {
            fieldName := errs[i][0]
            fmt.Printf("field %s contains some errors\n", fieldName)
            for errIdx := 1; errIdx < len(errs[i]); errIdx++ {
                fmt.Printf("\t- %s\n", errs[i][errIdx])
            }
        }

        return
    }

    fmt.Println("There's no problem with the input")
}
```
