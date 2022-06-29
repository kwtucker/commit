# Commit

Commit will help you craft a formatted commit message by reading staged files for statements between the defined delimiters. The output of the commit command will be a formatted message that can be used when actually committing to source control.

## Install

```
$ go get -u github.com/kwtucker/commit
```

## Configuration

Defaults:

```
COMMIT_PREFIX="*"
COMMIT_TITLE_PREFIX=""
```

Set in your environment.

```
export COMMIT_PREFIX="-"
export COMMIT_TITLE_PREFIX="Ticket-123:"
```

Or Envionment file ".envrc"

```
COMMIT_PREFIX="-"
COMMIT_TITLE_PREFIX="Ticket-123:"
```

## Delimiters

Delimiters are use as markers to retrieve the start and end of a commit message entry.
The commit message can be a single line or span multiple lines both wrapped in the start and end delimiters.

| Delimiter | Description                 |
| --------- | --------------------------- |
| `(:`      | Start of a new commit entry |
| `:)`      | End of a commit entry       |

## How to write commits

- In your source code. The idea here is to write the commit as you code.

  ```
  // (:Adds a commit formatter:) <----- This will be read/removed by commit

  // String returns a formatted commit message.
  func String() string {
   return "Great commit message"
  }
  ```

  Run

  ```
  $ commit -t "Commit file"
  Commit file

  * Adds a commit formatter
  ```

- In a ".commit" file.
  The ".commit" file will allow you to write commits in one place if you do not want to add the commit to the source code files. The ".commit" file can be added to your git ignore and still be read. The ".commit" file will be read and parsed where the commit command is executed. This means if there is ".commit" files in other sub directories they will be ignored.

  ```
  (:Adds a commit formatter:)

  (:
  This will be concatenated
  into one line for the final output.
  :)
  ```

  Run

  ```
  $ commit -t "Commit file"
  Commit file

  * Adds a commit formatter

  * This will be concatenated into one line for the final output.
  ```

## Run

1. Stage the files you want to commit.

2. Run the commit command with an optional title message. This will not actually commit to version control. Commit will print a formatted message to use when actually committing.

   ```
   $ commit --copy --rm-text -title-msg "Added an exciting feature"

   Added an exciting feature

   * Feature
   ```

## Help

```
$ commit --help
Commit will help construct a commit message.

Usage:
  commit [flags]

Flags:
  -c, --copy               copy commit message to clipboard
  -d, --dry-run            dry run to inspect the result
  -h, --help               help for commit
  -r, --rm-text            remove text within the delimiters from the file after reading message
  -t, --title-msg string   quoted title of the commit message

```
