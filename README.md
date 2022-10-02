# ascot - AWS Security Compliance and Operations Toolkit

This is a single-binary tool that can be easily used by GRC (Compliance)
and Security Operations teams if they need to support AWS infrastructure
but don't have good experience yet with the AWS CLI.

Most, if not all, of the commands in this tool can also be done with
the AWS CLI and some `bash`, but not everyone is comfortable in this
environment.  This tool solves for that.

The tool is built to have multiple sub-commands, one for each type of
investigation being performed.

## Features

### Find the Owner of a given AWS Access Key ID

Given an AWS Access Key ID, this command searches all IAM users in your
AWS account (either given a `--profile` or using the default) and attempts
to find a matching owner.

```bash
ascot access-key-owner --access-key-id AKIA123...
```

If an owner is found, the ARN of the IAM user is displayed.  If the key
is active, which is presents a risk during a key exposure incident, the
output displays an alert.

**It is recommended to have unrestricted `iam:ListUser` privileges in the
AWS account in order to ensure all users are searched.**

### List all Active AWS Regions

Returns an alphabetical list of the regions that are active within the
AWS account.  This is helpful when needing to iterate over all regions
looking for a particular resource (e.g. listing all EC2 instances in all
regions).

```bash
ascot active-regions
```

## Contributing

Any help is appreciated.  Please put your changes in a branch and then
create a Pull Request (PR).

## License

[MIT](LICENSE)
