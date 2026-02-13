# Commit

## Install

```
$ go install github.com/kwtucker/commit
```

## Configuration

Defaults:

```
COMMIT_BODY_PREFIX="*"
```

Set in your environment.

```
export COMMIT_BODY_PREFIX="-"
```

Or Envionment file ".envrc"

```
COMMIT_BODY_PREFIX="-"
```

## Run

1. Stage the files you want to commit.

2. Run the commit command with an optional copy flag.

   ```
   $ commit --copy

   feat(core): Added an exciting feature

   * Feature
   ```

## Help

```
$ commit --help
Constructs formatted commit messages

Usage:
  commit [flags]

Flags:
  -c, --copy   copy commit message to clipboard
  -h, --help   help for commit
```
