Ambrosia<!-- omit from toc -->
========
[![tests](https://github.com/reactorsh/ambrosia/actions/workflows/tests.yml/badge.svg)](https://github.com/reactorsh/ambrosia/actions/workflows/tests.yml)

Ambrosia is a cross-platform command line tool for improving the text datasets you use for machine learning.

It has fast versions of the usual dataset tasks: dedupe, filtering, checking lengths, etc.

It also has a unique LLM-based filtering option, `psort`, which sends entries in your dataset to an LLM for sorting and filtering.

## Table of Contents<!-- omit from toc -->
- [Release Status](#release-status)
- [Installing](#installing)
- [Getting started](#getting-started)
- [Commands and Options](#commands-and-options)
  - [Global Options](#global-options)
  - [whitespace](#whitespace)
  - [dedupe](#dedupe)
  - [length](#length)
  - [filter](#filter)
  - [psort](#psort)

## Release Status

Ambrosia is pre-1.0 software.  You may encounter bugs.  

Please report (or even better, send a PR) them.

## Installing

Ambrosia is available as a single binary.  There are no dependencies.  

Simply download the latest version from the [release page](https://github.com/reactorsh/ambrosia/releases) and run it from your terminal.

## Getting started

A complete feature walkthrough of cleaning a dataset is available here.

## Commands and Options

Note: this is only intended to be a reference.  Refer to the post linked in the [Getting Started](#getting-started) section for a detailed walkthrough.

All commands have a `--help` option that will print out the available options.

All flags can also be specified via environment variables.  The environment variable names will be output when you use the `--help` option.

### Global Options

`--debug, -d`  
Setting this flag will enable debug logging.  You should definitely use this if you're experiencing an issue or want more details.  Specifying it multiple times (`-d -d`) will increase verbosity for commands that support it.

`--json, -j`
The expectation is that most people will be executing commands interactively.  You can set this if you aren't and would prefer JSON-structured outputs.

### whitespace
`whitespace` trims Unicode-defined whitespace at the beginning and end of a field.

`--fields, -f`<br>
Specifies the fields to trim.  Multiple fields can be selected by passing them as a comma-separated list.  E.g., `--fields input,output`.

### dedupe

`dedupe` is used to deduplicate a dataset.  It has one required value, `field`, and a few options.  Empty fields are *not* treated as duplicates.

`--fields, -f`<br>
Specifies the field(s) in the data to compare on.  E.g. for a dataset of objects like this:

```json
[
  {
    "instruction": "This is an instruction."
  },
  {
    "input": "This is an input."
  },
  {
    "output": "This is an output."
  }
]
```

You would use `--fields instruction` to dedupe on the `instruction` field.  You can specify multiple fields by passing a comma-separated list, like so: `--fields instruction,input`.  

When multiple fields are specified, the data is compared with each field being newline-delimited and prefixed with the field name.  Continuing with the examples provided, the ROUGE-L score would be calculated on the following text:

```
instruction: this is an instruction.
input: this is an input.
```

Note that empty fields are not considered to be duplicates of each other.

`--ignore-case, -i`<br>
If this flag is set, the comparison will be case-insensitive.  All Unicode characters will be lower-cased before comparison.

`--rougel, -rl`<br>
Enable [ROUGE-L](https://aclanthology.org/P04-1077.pdf) deduplication.  Words are not stemmed, and Unicode characters are always lower-cased before comparison.

`--rl-threshold, -rlt`<br>
If `--rougel` is set, this option can be used to specify the threshold for ROUGE-L deduplication.  

`--progress, -p`<br>
If set, `--progress` will display a progress bar for ROUGE-L deduplication.

### length

`length` filters data by *byte count*.  Tokenizers vary, so filtering on the number of tokens output by a particular tokenizer is not consistently meaningful.  Filtering by byte count is a more reliable way to filter on length.

`--fields, -f`<br>
Specifies the fields to use for determining the length.  Multiple fields can be selected by passing them as a comma-separated list.  E.g., `--fields input,output`.

`--min`<br>
Specifies the minimum length to filter on.  E.g., `--min 10` will filter out all entries with a length less than or equal to 10 bytes.

`--max`<br>
Specifies the maximum length to filter on.  E.g., `--max 100` will filter out all entries with a length greater than or equal to 100 bytes.

### filter

The `filter` command is used to filter data containing particular strings.  These can be 'simple' strings where, if the string is present in the data, it will be filtered out, or 'regex' strings, where the provided string is treated as a regular expression that will be matched against the data.

`--field, -f`<br>
Which field to check against.  This is required.

`--regex, -r`<br>
If this flag is set, the provided string or wordlist will be treated as a regular expression.

`--string, -s`<br>
If this flag is set, the provided string will be used for filtering.  E.g., `-s foo` will filter out entries that contain 'foo'.

`--wordlist, -w`<br>
If this flag is set, the value will be treated as a path to a file containing a newline-delimited list of strings.  Each string will be used for filtering.  E.g., `-w ./wordlist.txt` will filter out entries that contain any of the strings in `./wordlist.txt`.

### psort

The `psort` command sorts data using a provided prompt and an LLM.  The response from the LLM is used to sort the data by taking the first Unicode character of the response and writing the data to a file suffixed with the character.  As an example, if the LLM responded with:

> `foo bar baz...`

and

> 'baz bar foo...'

The data would be sorted into two different files named:

`<INPUTFILE>_prompt_f.<EXT>`<br>
and<br>
`<INPUTFILE>_prompt_b.<EXT>`.

The prompt used to obtain that response is a combination of your data, `--sysprompt` (if supported), and `--instruction`.  For each piece of data, the sysprompt will be set, and the request for a response will be structured as such:

```
<INSTRUCTION>

field1: <FIELD1>
field2: <FIELD2>
field3: <FIELD3>
```

`--model, -m`<br>
Specify the LLM model to use.  At this time, ambrosia supports:

* `gpt-3.5-turbo`
* `gpt-4`

`--token, -t`<br>
The authentication token to use for the specified model, if required.  In the case of OpenAI models, this would be your OAI platform token, e.g., `sk-***`.

`--baseurl, -b`<br>
The base URL to use for the specified model (if required), or if you wish to override the default.

`--instruction, -i`<br>
This is the text that will be passed to the LLM *before* the data.  It should be a prompt that will result in the LLM generating a response that can be used to filter the data.  E.g., you might use `--instruction "Is this an intelligible sentence?"`.  

Tip: try to force the LLM to restrict its responses to something easy to sort on.  Ask the LLM to respond with `yes` or `no`, or `0`/`1`/`2`.

`--end-instruction, -e`<br>
This is text that will be added *after* the data.  This is helpful for certain use cases where the dataset unintentionally acts as a prompt injection.

`--sysprompt, -sp`<br>
The system prompt to use for models that support it.

`--fields, -f`<br>
The fields to pass along in the prompt for each piece of data.  Multiple fields can be specified by passing them as a comma-separated list.  E.g., `--fields input,output`.  By default, all fields are sent, **but the order is not guaranteed**.  If you need a specific order (and you probably do), specify the fields explicitly.

`--include-resp, -ir`<br>
If set, `--include-resp` will add the complete response from the LLM as a new field named `ambrosia` in the output file(s).  Very helpful for prompt debugging.

`--dry-run, -d`<br>
If set, ambrosia will not actually make any inference requests to the LLM; it will just print the generated prompt to the console.  

Doing this before spending compute/money on prompt filtering is a good idea.

`--concurrency, -c`<br>
This is the number of simultaneous inference requests to make.  The default is 1, but you can increase this to speed up the filtering process.  Be careful, though, as some models limit the number of simultaneous requests you can make.

`--rpm`<br>
This sets the maximum number of inference requests per minute.

`--tpm`<br>
This sets the maximum number of tokens per minute.  

We can't determine the number of tokens an arbitrary model will use.  Instead, when we check the token limiter for additional capacity, we use the number of bytes in the prompt.  This should always be less than the token count.

When we receive a response from the LLM that includes the total tokens consumed, the token limiter is decremented by the correct amount.

`--max-tokens, -mt`<br>
This is included in the request (for models that support it) and specifies the maximum number of tokens that should be returned.  

Setting this to a small value (e.g., the default of `5`) will save on costs but reduces interpretability.

`--timeout, -to`<br>
This is the timeout for each inference request in seconds.  The default is 15 seconds.  Timed-out requests will be retried indefinitely, once per second.

`--progress, -p`<br>
If set, `--progress` will display a progress bar for the filtering process.
