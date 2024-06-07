# TFModuleIndex

This tool injects a index of terraform modules in a `README.md` file.
It builds this index based on `module.yml` files in all terraform modules.

This module.yml currently contains the following fields:
```
short_description: "Short description of what the modules does"
```

It reads by default the `README.md` and injects the generated table between two markers:
```
<!-- BEGIN_MODULE_INDEX -->
<!-- END_MODULE_INDEX -->
```

Usage:
```
Usage of tfmoduleindex:
  -dir string
    	Directory to scan recursively (default ".")
  -inject string
    	File to inject module index (default "README.md")
  -template string
    	(Optional) Template file to use, e.g. .tfmoduleindex.tmpl
```

By default it uses the embedded [`template.tpl`](./template.tpl) but you can also override the template.

Example output:
```
<!-- BEGIN_MODULE_INDEX -->
| Module | Description |
|------|------|
| **Acm** | |
| [public/cert_dns](modules/acm/public/cert_dns) | Public ACM Certificate with DNS Validation |
| **Cloudwatch** | |
| [eventbridge/lambda_schedule](modules/cloudwatch/eventbridge/lambda_schedule) | Lambda schedule (cron) |
| **Ecr** | |
| [repository/lifecycle_policy/generator](modules/ecr/repository/lifecycle_policy/generator) | ECR Lifecycle policy generator |
| [repository/lifecycle_policy/generic](modules/ecr/repository/lifecycle_policy/generic) | ECR Lifecycle policy generic |
| [repository/policy/pull](modules/ecr/repository/policy/pull) | ECR policy to pull images |
<!-- END_MODULE_INDEX -->
```
