# jenkins_job_config Data Source

Get the config of a job within Jenkins.

## Example Usage

```hcl
data "jenkins_job_config" "example" {
  name        = "job-name"
  node        = "node-name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the job being read.
* `folder` - (Optional) The folder namespace containing this job.
* `node` - (Optional) The node used to filter the job configuration.
* `regex` - (Optional) The regex used to filter the job configuration.


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The full canonical job path, E.G. `/job/job-name`.
* `config` - The configuration of the job.
