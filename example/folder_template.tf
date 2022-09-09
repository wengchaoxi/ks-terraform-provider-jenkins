resource "jenkins_folder_template" "example" {
  name = "folder-name-example"

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
