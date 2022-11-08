# ascot - AWS Security Compliance and Operations Toolkit

[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=scottbrown_ascot&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=scottbrown_ascot)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=scottbrown_ascot&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=scottbrown_ascot)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=scottbrown_ascot&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=scottbrown_ascot)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=scottbrown_ascot&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=scottbrown_ascot)

This is a single-binary tool that can be easily used by GRC (Compliance)
and Security Operations teams if they need to support AWS infrastructure
but don't have good experience yet with the AWS CLI.  This tool is intended
to complement the use of AWS CLI, not replace it.

This tool is also intended as a learning tool, for those that lack the
experience in AWS.  You can find what IAM permissions are needed for a
command by running `ascot [COMMAND] --show-required-permissions`.  You
can also run `ascot [COMMAND] --how-it-works` in order to understand the
business logic being run, which can aid you when you want to use the
AWS CLI directly instead of this tool.

Most, if not all, of the commands in this tool can also be done with
the AWS CLI and some `bash`, but not everyone is comfortable in this
environment.  This tool solves for that.

The tool is built to have multiple sub-commands, one for each type of
investigation being performed.

# Features

A list of all supported commands can be found by runnning `ascot -h`.

### Figure out who you are

Running `ascot` with no arguments will print out the ARN of the user you
have authenticated with (or assumed, when using a role).

```bash
$ ascot
AWS login was successful.
You are currently logged in as arn:aws:iam::012345678901:user/johndoe
```

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

### Audit Default VPCs

Allows you to see which regions in the AWS account still allow their
default VPCs to exist.

```bash
ascot audit-default-vpcs
```

If a default VPC exists in the region, the result is `FAIL`.  Otherwise,
the region receives a `PASS`.

### Missing Images

Lists any EC2 instances that are using AMIs that no longer exist.  This
means if the EC2 instance needs to be rebuilt, it will fail because it
depends on an AMI that cannot be found.

This searches all regions.

```bash
ascot missing-images
```

The output is a list of instance IDs that are affected by the missing
AMI, and the AMI ID that is missing.

## Contributing

Any help is appreciated.  Please put your changes in a branch and then
create a Pull Request (PR).

## License

[MIT](LICENSE)
