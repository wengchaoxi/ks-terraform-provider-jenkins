package jenkins

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceJenkinsJobConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJenkinsJobConfigRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The unique name of the JenkinsCI job.",
				Required:         true,
				ValidateDiagFunc: validateJobName,
			},
			"folder": {
				Type:             schema.TypeString,
				Description:      "The folder namespace that the job exists in.",
				Optional:         true,
				ValidateDiagFunc: validateFolderName,
			},
			"xml_node": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The xml_node used to filter the job configuration.",
			},
			"regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The regex used to filter the job configuration.",
			},
			"config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The XML configuration of the job.",
			},
		},
	}
}

func dataSourceJenkinsJobConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jenkinsAdapter)

	name := d.Get("name").(string)
	folderName := d.Get("folder").(string)

	id := formatFolderName(folderName + "/" + name)
	job, err := client.GetJob(ctx, id)
	if err != nil {
		log.Printf("[DEBUG] jenkins::job::config - Could not find job %q: %s", id, err.Error())
		return nil
	}
	config, err := job.GetConfig(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::job::config - Error get config in job %q: %w", id, err))
	}

	xmlNode := d.Get("xml_node").(string)
	if xmlNode != "" {
		config, err = filterJobConfigByXMLNode(config, xmlNode)
		if err != nil {
			log.Printf("[DEBUG] jenkins::job::config - Job %q: %s", id, err.Error())
		}
	}

	regex := d.Get("regex").(string)
	if regex != "" && config != "" {
		config, err = filterJobConfigByRegex(config, regex)
		if err != nil {
			log.Printf("[DEBUG] jenkins::job::config - Job %q: %s", id, err.Error())
		}
	}

	d.Set("config", config)
	d.SetId(job.Base)
	return nil
}

func filterJobConfigByRegex(config string, regex string) (string, error) {
	re := regexp.MustCompile(fmt.Sprintf("%s", regex))
	result := re.FindStringSubmatch(config)
	if len(result) == 0 {
		return "", errors.New("no match result")
	}
	return result[0], nil
}

func filterJobConfigByXMLNode(config string, node string) (string, error) {
	result, err := filterJobConfigByRegex(config, fmt.Sprintf(`<%s ?.*>[^\0]*</%s>|<%s ?.*/>`, node, node, node))
	if err != nil {
		return "", fmt.Errorf("not found XML node: %q", node)
	}
	return result, nil
}
