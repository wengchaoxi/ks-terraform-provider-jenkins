package jenkins

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceJenkinsPlugins() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluginsRead,
		Schema: map[string]*schema.Schema{
			"plugins": {
				Type:        schema.TypeList,
				Description: "The list of the Jenkins plugins.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePluginsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jenkinsAdapter)

	p, err := client.GetPlugins(ctx, 1)
	if err != nil {
		return diag.FromErr(err)
	}

	var plugins []map[string]string
	for _, v := range p.Raw.Plugins {
		p := make(map[string]string)
		p["name"] = v.ShortName
		p["version"] = v.Version
		plugins = append(plugins, p)
	}

	if err := d.Set("plugins", plugins); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("jenkins-data-source-plugins-id")
	return nil
}
