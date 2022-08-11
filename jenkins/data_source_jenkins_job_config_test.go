package jenkins

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJenkinsJobConfigDataSource_basic(t *testing.T) {
	testDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(testDir, "test.xml"), testXML, 0644)
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource jenkins_job foo {
	name = "tf-acc-test-%s"
	template = templatefile("%s/test.xml", {
		description = "Acceptance testing Jenkins provider"
	})
}

data jenkins_job_config foo {
	name = jenkins_job.foo.name
	node = "properties"
}`, randString, testDir),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_job.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_job_config.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_job_config.foo", "name", "tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_job_config.foo", "config", "<properties/>"),
				),
			},
		},
	})
}

func TestAccJenkinsJobConfigDataSource_nested(t *testing.T) {
	testDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(testDir, "test.xml"), testXML, 0644)
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource jenkins_folder foo {
	name = "tf-acc-test-%s"
}

resource jenkins_job sub {
	name = "subfolder"
	folder = jenkins_folder.foo.id
	template = templatefile("%s/test.xml", {
		description = "Acceptance testing Jenkins provider"
	})
}

data jenkins_job_config sub {
	name = jenkins_job.sub.name
	folder = jenkins_job.sub.folder
	regex = "<properties>[^\\0]*</properties>|<properties/>"
}`, randString, testDir),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_folder.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_job.sub", "id", "/job/tf-acc-test-"+randString+"/job/subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_job_config.sub", "name", "subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_job_config.sub", "folder", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_job_config.sub", "config", "<properties/>"),
				),
			},
		},
	})
}
