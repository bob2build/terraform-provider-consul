package consul

import (
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceConsulPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulPolicyCreate,
		Read:   resourceConsulPolicyRead,
		Update: resourceConsulPolicyUpdate,
		Delete: resourceConsulPolicyDelete,

		Schema: map[string]*schema.Schema{
			"policy": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"token": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "terraform_managed",
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "client",
			},
		},
	}
}

func resourceConsulPolicyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*consulapi.Client)

	var dc string
	if v, ok := d.GetOk("datacenter"); ok {
		dc = v.(string)
	} else {
		var err error
		if dc, err = getDC(d, client); err != nil {
			return err
		}
	}

	var token string
	if v, ok := d.GetOk("token"); ok {
		token = v.(string)
	}

	// Setup the operations using the datacenter
	wOpts := consulapi.WriteOptions{Datacenter: dc, Token: token}

	name := d.Get("name").(string)
	policy := d.Get("policy").(string)
	acltype := d.Get("type").(string)
	id := d.Get("token").(string)

	aclClient := client.ACL()
	acl_entry := &consulapi.ACLEntry{Name: name, Rules: policy, Type: acltype, ID: id}
	id, _, e := aclClient.Create(acl_entry, &wOpts)
	if e == nil {
		d.SetId(id)
		d.Set("token", id)
	} else {
		return e
	}
	return nil
}

func resourceConsulPolicyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*consulapi.Client)

	var dc string
	if v, ok := d.GetOk("datacenter"); ok {
		dc = v.(string)
	} else {
		var err error
		if dc, err = getDC(d, client); err != nil {
			return err
		}
	}

	var token string
	if v, ok := d.GetOk("token"); ok {
		token = v.(string)
	}

	// Setup the operations using the datacenter
	qOpts := consulapi.QueryOptions{Datacenter: dc, Token: token}

	aclClient := client.ACL()
	acl, _, e := aclClient.Info(d.Get("token").(string), &qOpts)

	if e == nil && acl != nil {
		d.Set("type", acl.Type)
		d.Set("policy", acl.Rules)
		d.Set("token", acl.ID)
		d.Set("name", acl.Name)
	} else {
		return e
	}
	return nil
}

func resourceConsulPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*consulapi.Client)

	var dc string
	if v, ok := d.GetOk("datacenter"); ok {
		dc = v.(string)
	} else {
		var err error
		if dc, err = getDC(d, client); err != nil {
			return err
		}
	}

	var token string
	if v, ok := d.GetOk("token"); ok {
		token = v.(string)
	}

	// Setup the operations using the datacenter
	wOpts := consulapi.WriteOptions{Datacenter: dc, Token: token}

	if d.HasChange("type") || d.HasChange("policy") || d.HasChange("name") {
		name := d.Get("name").(string)
		policy := d.Get("policy").(string)
		acltype := d.Get("type").(string)
		id := d.Get("token").(string)

		aclClient := client.ACL()
		acl_entry := &consulapi.ACLEntry{Name: name, Rules: policy, Type: acltype, ID: id}
		id, _, e := aclClient.Create(acl_entry, &wOpts)
		if e == nil {
			d.SetId(id)
			d.Set("token", id)
		} else {
			return e
		}
	}
	return nil

}

func resourceConsulPolicyDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
