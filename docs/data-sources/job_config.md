# jenkins_job_config Data Source

Get the XML configuration of a job within Jenkins.

## Example Usage

```hcl
data "jenkins_job_config" "example" {
  name        = "job-name"
  xml_node    = "xml-node-name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the job being read.
* `folder` - (Optional) The folder namespace containing this job.
* `xml_node` - (Optional) The xml_node used to filter the XML configuration of the job.
* `regex` - (Optional) The regex used to filter the XML configuration of the job.

> Note: If both `xml_node` and `regex` are empty, the complete XML configuration of the job will be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The full canonical job path, E.G. `/job/job-name`.
* `config` - The XML configuration of the job.
