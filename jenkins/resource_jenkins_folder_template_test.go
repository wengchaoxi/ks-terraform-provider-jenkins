package jenkins

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	//go:embed "resource_jenkins_folder_template_test.xml"
	testFloderTemplateXML []byte
	//go:embed "resource_jenkins_folder_template_want.xml"
	testFloderTemplateXMLWant string
)

func TestAccJenkinsFolderTemplate_basic(t *testing.T) {
	testDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(testDir, "test.xml"), testFloderTemplateXML, 0644)
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJenkinsFolderTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_folder_template foo {
					name = "%s"
 					template = templatefile("%s/test.xml", {
 						description = "Acceptance testing Jenkins provider"
 					})
				}`, randString, testDir),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_folder_template.foo", "id", "/job/"+randString),
					resource.TestCheckResourceAttr("jenkins_folder_template.foo", "name", randString),
					resource.TestCheckResourceAttr("jenkins_folder_template.foo", "template", testFloderTemplateXMLWant),
				),
			},
		},
	})
}

func TestAccJenkinsFolderTemplate_nested(t *testing.T) {
	testDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(testDir, "test.xml"), testFloderTemplateXML, 0644)
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJenkinsFolderTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_folder_template foo {
					name = "tf-acc-test-%s"
					template = templatefile("%s/test.xml", {
						description = "Acceptance testing Jenkins provider"
					})
				}

				resource jenkins_folder_template sub {
					name = "subfolder"
					folder = jenkins_folder_template.foo.id
					template = templatefile("%s/test.xml", {
						description = "Acceptance testing Jenkins provider"
					})
				}`, randString, testDir, testDir),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_folder_template.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_folder_template.foo", "name", "tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_folder_template.sub", "id", "/job/tf-acc-test-"+randString+"/job/subfolder"),
					resource.TestCheckResourceAttr("jenkins_folder_template.sub", "name", "subfolder"),
					resource.TestCheckResourceAttr("jenkins_folder_template.sub", "folder", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_folder_template.sub", "template", testFloderTemplateXMLWant),
				),
			},
		},
	})
}

func testAccCheckJenkinsFolderTemplateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(jenkinsClient)
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jenkins_folder_template" {
			continue
		}

		_, err := client.GetJob(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Folder %s still exists", rs.Primary.ID)
		}
	}

	return nil
}
