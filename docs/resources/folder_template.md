# jenkins_folder_template Resource

Manages a folder by template within Jenkins.

~> The Jenkins installation that uses this resource is expected to have the [Cloudbees Folders Plugin](https://plugins.jenkins.io/cloudbees-folder) installed in their system.

## Example Usage

```hcl
resource "jenkins_folder_template" "example" {
  name = "folder-name"
  template = templatefile("${path.module}/folder-template.xml", {
    description = "A top-level folder"
  })
}

resource "jenkins_folder_template" "example_child" {
  name   = "child-name"
  folder = jenkins_folder_template.example.id
  template = templatefile("${path.module}/folder-template.xml", {
    description = "A nested subfolder"
  })
}
```

And in `folder-template.xml`:

```xml
<?xml version="1.1" encoding="UTF-8"?><com.cloudbees.hudson.plugins.folder.Folder plugin="cloudbees-folder@6.729.v2b_9d1a_74d673">
  <description>${description}</description>
  <properties/>
  <folderViews class="com.cloudbees.hudson.plugins.folder.views.DefaultFolderViewHolder">
    <views>
      <hudson.model.AllView>
        <owner class="com.cloudbees.hudson.plugins.folder.Folder" reference="../../../.."/>
        <name>All</name>
        <filterExecutors>false</filterExecutors>
        <filterQueue>false</filterQueue>
        <properties class="hudson.model.View$PropertyList"/>
      </hudson.model.AllView>
    </views>
    <tabBar class="hudson.views.DefaultViewsTabBar"/>
  </folderViews>
  <healthMetrics/>
  <icon class="com.cloudbees.hudson.plugins.folder.icons.StockFolderIcon"/>
</com.cloudbees.hudson.plugins.folder.Folder>
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the folder being created.
* `folder` - (Optional) The folder namespace to store the subfolder in. If creating in a nested folder structure you may separate folder names with `/`, such as `parent/child`. This name cannot be changed once the folder has been created, and all parent folders must be created in advance.
* `template` - (Required) A Jenkins-compatible XML template to describe the folder. You can retrieve an existing folder's XML by appending `/config.xml` to its URL and viewing the source in your browser. The `template` property is rendered using a Golang template that takes the other resource arguments as variables. Do not include the XML prolog in the definition.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The full canonical folder path, E.G. `/job/parent`.

## Import

Folders may be imported by their canonical name, e.g.

```sh
$ terraform import jenkins_folder_template.example /job/folder-name
```
