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
ascot find-access-key-owner --access-key-id AKIA123...
```

If an owner is found, the ARN of the IAM user is displayed.

**It is recommended to have unrestricted `iam:ListUser` privileges in the
AWS account in order to ensure all users are searched.**

## License

[MIT](LICENSE)
