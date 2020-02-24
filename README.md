# Turbo-Pancake

## Table of Contents
* [Getting Started](#getting-started)
* [About](#about)
* [Usage](#usage)

## Getting Started

### Requirements
* GNU Make
* Golang `>=1.12`

### Compiling
```
~# make
```

### Install:
```
~# make install
```

### Test:
```
~# make test
```

## About

Inspired and Originates from [@voodooEntity's Threader](https://github.com/voodooEntity/threader)

However, this is a Non-Compatible Fork and barely shares any of the original code. You could call it a complete rewrite.

In fact, the original Threader project and this one share a total of 40 lines of code, most of them are the golang imports.

Because of this the tool has been renamed to `Turbo Pancake` with permission of the original author to avoid ambiguity.

## Usage

### Command Line Interface

Using the CLI Utility is very straight forward.

```
~# turbo-pancake -h
Usage of turbo-pancake:
  -command string
        Command to execute
  -d    Debug Output
  -delimiter string
        Input Delimiter, as fmt.Sprintf (default "\\n")
  -out-format string
        Output Format, as fmt.Printf (default "%s\\n")
  -threads int
        Amount of threads to use (default 1)
```

|Option|Description|
|------|-----------|
|`-command`|Almost arbitrary shell code to execute. Heavily depends on the workers' environments|
|`-d`|Flag to enable Debug Mode. This will spam a lot.|
|`-delimiter`|String or sequence to delimit the STDIN Input to distribute work|
|`-out-format`|fmt.Printf Parsable string to format Outputs by|
|`-threads`|Amount of Workers to keep running|

#### Variables Provided to `Command`

|Variable|Content|
|--------|-------|
|`$INPUT`|Input passed to the Worker|
|`$BASEINPUT`|Base64 Encoded `$INPUT`|
|`$JOBID`|UUIDv4 unique JobID used to identify and trace the current Job|
|`$WORKERID`|UUIDv4 unique WokerID used to identify and trace the Worker executing the Job|
|`$COMMAND`|Original Command that will be executed by the Worker|

### Examples

```
~# echo -n '1,2,3' | turbo-pancake -delimiter ',' -threads 2 -command 'echo -n $INPUT;sleep $((3-$INPUT))' -out-format '%s'
231
~#
```

### Networked Interface

TBA