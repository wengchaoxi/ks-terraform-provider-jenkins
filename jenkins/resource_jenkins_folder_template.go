package jenkins

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceJenkinsFolderTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJenkinsFolderTemplateCreate,
		ReadContext:   resourceJenkinsFolderTemplateRead,
		UpdateContext: resourceJenkinsFolderTemplateUpdate,
		DeleteContext: resourceJenkinsJobDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The unique name of the JenkinsCI folder.",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateJobName,
			},
			"folder": {
				Type:             schema.TypeString,
				Description:      "The folder namespace that the folder will be added to as a subfolder.",
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateFolderName,
			},
			"template": {
				Type:        schema.TypeString,
				Description: "The configuration file template, used to communicate with Jenkins.",
				Required:    true,
			},
		},
	}
}

func resourceJenkinsFolderTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	name := d.Get("name").(string)
	folderName := d.Get("folder").(string)

	// Validate that the folder exists
	if err := folderExists(ctx, client, folderName); err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::create - Could not find folder '%s': %w", folderName, err))
	}

	xml := d.Get("template").(string)

	folders := extractFolders(folderName)
	_, err := client.CreateJobInFolder(ctx, xml, name, folders...)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::create - Error creating job for %q in folder %s: %w", name, folderName, err))
	}

	log.Printf("[DEBUG] jenkins::create - job %q created in folder %s", name, folderName)
	d.SetId(formatFolderName(folderName + "/" + name))

	return resourceJenkinsFolderTemplateRead(ctx, d, meta)
}

func resourceJenkinsFolderTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	name, folders := parseCanonicalJobID(d.Id())

	log.Printf("[DEBUG] jenkins::read - Looking for job %q", name)

	job, err := client.GetJob(ctx, name, folders...)
	if err != nil {
		if strings.HasPrefix(err.Error(), "404") {
			// Job does not exist
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("jenkins::read - Job %q does not exist: %w", name, err))
	}

	// Extract the raw XML configuration
	config, err := job.GetConfig(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::read - Job %q could not extract configuration: %v", job.Base, err))
	}

	log.Printf("[DEBUG] jenkins::read - Job %q exists", job.Base)
	d.SetId(job.Base)

	// The config content returned by Jenkins removes the '\n' at the end of the text
	if err := d.Set("template", config+"\n"); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] jenkins::read - \n%s", config)

	if err := d.Set("name", name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("folder", formatFolderID(folders)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceJenkinsFolderTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	name, folders := parseCanonicalJobID(d.Id())

	// grab job by current name
	job, err := client.GetJob(ctx, name, folders...)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::update - Could not find job %q: %w", name, err))
	}

	xml := d.Get("template").(string)

	err = job.UpdateConfig(ctx, xml)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::update - Error updating job %q configuration: %w", name, err))
	}

	return resourceJenkinsFolderTemplateRead(ctx, d, meta)
}
