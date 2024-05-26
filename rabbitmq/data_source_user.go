package rabbitmq

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

func dataSourcesUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"max_connections": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_channels": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	rmqc := meta.(*rabbithole.Client)

	user, err := rmqc.GetUser(name)
	if err != nil {
		return diag.Errorf("user '%s' is not found: %#v", name, err)
	}
	d.Set("name", user.Name)

	if len(user.Tags) > 0 {
		var tagList []string
		for _, v := range user.Tags {
			if v != "" {
				tagList = append(tagList, v)
			}
		}
		if len(tagList) > 0 {
			d.Set("tags", tagList)
		}
	}

	myUserLimits, err := rmqc.GetUserLimits(name)
	if err != nil {
		return diag.Errorf("error to get user limits for '%s': %#v", name, err)
	}

	if len(myUserLimits) > 0 {
		if val, ok := myUserLimits[0].Value["max-connections"]; ok {
			d.Set("max_connections", strconv.Itoa(val))
		}

		if val, ok := myUserLimits[0].Value["max-channels"]; ok {
			d.Set("max_channels", strconv.Itoa(val))
		}
	}

	d.SetId(user.Name)
	return nil
}
