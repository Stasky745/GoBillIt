# GoBillIt
GoBillIt allows you to generate an invoice from a config file using [Invoice-Generator.com](https://invoice-generator.com/) API. It also can send notifications and add some additional items to it using [NTFY](https://github.com/binwiederhier/ntfy), and automatically send the generated PDFs to emails.

It also has implemented conversion using [ApiLayer's exchange rate API](https://apilayer.com/marketplace/exchangerates_data-api).

- [GoBillIt](#gobillit)
  - [Usage](#usage)
    - [Docker](#docker)
    - [Binary](#binary)
  - [Configuration](#configuration)
    - [Config File](#config-file)
    - [Config Dir](#config-dir)
    - [Variables](#variables)
      - [Date Format](#date-format)
    - [Environment Variables](#environment-variables)
  - [NTFY](#ntfy)
  - [Possible Improvements (not being worked on)](#possible-improvements-not-being-worked-on)


## Usage
### Docker
The preferred method is by using Docker:

```sh
docker run \
    -v /path/to/config.yaml:/etc/gbi/config.yaml \ 
    -v /path/to/config/dir:/etc/gbi/config.d \
    ghcr.io/stasky745/gbi:v1.0.2
```

### Binary
This repository also makes executables available for use. Feel free to download those and use them from the command line.

```sh
/path/to/gbi --help
```

## Configuration

GBI uses configuration files in any of these extensions: `.yaml`, `.yml`, `.json`.

### Config File
This is the base config file. By default, GBI searches this file on `/etc/gbi/config.yaml` unless you specify it through the `GBI_CONFIG` environment variable or the `-config <path>` flag.

You can see an example of this file in `config.example.yaml`

### Config Dir
If you want to split the configuration in different files (eg: `ntfy.yaml`, `email.yaml`, etc.) you can add any number of files with the accepted extensions in the `GBI_CONFIG_DIR` (or `-config_dir` flag), which defaults to `/etc/gbi/config.d`.

These files must have the format shown in `config.example.yaml`, so if you wanted to add just the invoice's item list in a file called `/etc/gbi/config.d/items.yaml`, it would contain this:

```yaml
inv:
  items:
  list:
    - label: main
      name: "{{date[m]}} Services"
      description: "1 {{ date[m] }} - {{ date[D] }} {{ date[m] }}"
      quantity: 1
      unit_cost: 4800
```

### Variables
Items' names and descriptions, extras and many other informational fields allow the use of variables.

The format is:
```yaml
{{<key>[<option>]}}
```

The possible values:

|Key|Options|Description|
|---|-------|-----------|
|date|date format (shown below)|Current date when running the program in the format of the option. Default format is `YYYY-MM-DD`. Same as `inv.date` key.|
|conversion|---|Conversion value used for this run of the program. Same as `inv.conversion` and `inv.conversion.value` keys|

Any other key from the config file in the format of `key1.key2.key3` which would return `value1` in the following setting.

```yaml
key1:
  key2:
    key3: value1
```

This allows you to add more configuration to the file. For example a `vars:` section to use for this.

#### Date Format
This is what each key value would return for the date `2002-07-09`

|Key|Value|
|---|-----|
|YYYY|2002|
|YY|02|
|MM|07|
|M|7|
|DD|09|
|D|9|
|m|July|
|d|Tuesday|

### Environment Variables
If you prefer using environment variables rather than a config file, you can do so by appending the `GBI_` prefix and using underscores (`_`) to split the path.

For example:

```yaml
inv:
  conversion:
    value: 0.93583
```

can be set with `GBI_INV_CONVERSION_VALUE=0.93583`.

## NTFY
You can setup a cron so that this happens once every so often, but even then you might not be home or want to configure a bit further. GBI can send a notification through NTFY which allows you to add a few pre-configured items in exceptional cases.

If configured to do so, GBI will send a notification with the invoice in PDF format to an NTFY channel, where you can preview it and decide whether you want to add optional items (set in the path `ntfy.extras.items` of `config.yaml`). **Because of NTFY limitations, a maximum of 2 items can be configured**. 

It will keep asking to add items until done, and then preview the newly created invoice to be able to continue with the process (sending email or not).

## Possible Improvements (not being worked on)
- Move away from NTFY and create a phone app that will handle this better with webhooks (would have to open the app to the internet).
  - Open to implement better options instead.
- Use a *LaTeX* template instead of relying on [Invoice-Generator.com](https://invoice-generator.com/) API.
