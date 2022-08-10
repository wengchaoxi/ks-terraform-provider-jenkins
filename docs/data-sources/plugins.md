# jenkins_plugin Data Source

Get the information of plugins within Jenkins.

## Example Usage

```hcl
data "jenkins_plugins" "example" {}
```

## Argument Reference

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Fixed value: `jenkins-data-source-plugins-id`.
* `plugins` - The list of the Jenkins plugins.
    * `name` - The name of plugin.
    * `version` - The version of plugin.
