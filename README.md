# sc-grad-2025

The technical take home for SC graduate program of 2025.

## Getting started

Requires `Go` >= `1.20`

follow the official install instruction: [Golang Installation](https://go.dev/doc/install)

To run the code on your local machine
```
  go run main.go
```

## Folder structure

```
| go.mod
| README.md
| sample.json
| main.go
| folders
    | folders.go
    | folders_test.go
    | static.go
```

## Instructions

- This technical assessment consists of 2 components:
- Component 1:
  - within `folders.go`.
    - We would like you to read through, and run, the code.
    - Write some comments on what you think the code does.
    - Suggest some improvements that can be made to the code.
    - Implement any suggested improvements.
    - Write up some unit tests in `folders_test.go` for your new `GetAllFolders` method

- Component 2:
  - Extend your improved code to now facilitate pagination.
  - You can copy over the existing methods into `folders_pagination.go` to get started.
  - Write a short explanation of your chosen solution.

## What is pagination?
  - Pagination helps break down a large dataset into smaller chunks.
  - Those smaller chunks can then be served to the client, and are usually accompanied by a token pointing to the next chunk.
  - The end result could potentially look like this:
```
  original data: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

  This might result in the following payload to the client:
  { data: [1, 2, 3, ..., 10] }

  However, with pagination implemented, the payload might instead look like this:
  request() -> { data: [1, 2], token: "nQsjz" }

  The token could then be used to fetch more results:

  request("nQsjz") -> { data : [3, 4], token: "uJsnQ" }

  .
  .
  .

  And more results until there's no data left:

  { data: [9, 10], token: "" }
```

## Submission

Create a repo in your chosen git repository (make sure it is public so we can access it) and reply with the link to your code. We recommend using GitHub.


## Contact

If you have any questions feel free to contact us at: graduates@safetyculture.io