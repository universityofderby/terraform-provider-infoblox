package infoblox

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/fanatic/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceInfobloxRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfobloxRecordCreate,
		Read:   resourceInfobloxRecordRead,
		Update: resourceInfobloxRecordUpdate,
		Delete: resourceInfobloxRecordDelete,

		Schema: map[string]*schema.Schema{
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fqdn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4addr": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6addr": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nextavailableip": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"view": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceInfobloxRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	var comment, name, recID, value, view string
	var err error
	var nextavailableip bool
	var ttl int
	values := url.Values{}

	if attr, ok := d.GetOk("comment"); ok {
		comment = attr.(string)
	}
	if attr, ok := d.GetOk("name"); ok {
		name = attr.(string)
	}
	if attr, ok := d.GetOk("domain"); ok {
		name = strings.Join([]string{name, attr.(string)}, ".")
	}
	if attr, ok := d.GetOk("nextavailableip"); ok {
		nextavailableip = attr.(bool)
	}
	if attr, ok := d.GetOk("ttl"); ok {
		ttl = attr.(int)
	}
	if attr, ok := d.GetOk("value"); ok {
		value = attr.(string)
	}
	if attr, ok := d.GetOk("view"); ok {
		view = attr.(string)
	}

	switch strings.ToUpper(d.Get("type").(string)) {
	case "A":
		if nextavailableip {
			value = "func:nextavailableip:" + value
		}
		opts := &infoblox.Options{
			ReturnFields: []string{"name", "ipv4addr", "ttl", "view"},
		}
		body := &infoblox.RecordAObject{
			Comment:  comment,
			Ipv4Addr: value,
			Name:     name,
			Ttl:      ttl,
			View:     view,
		}
		recID, err = client.RecordA().Create(values, opts, body)
	case "AAAA":
		if nextavailableip {
			value = "func:nextavailableip:" + value
		}
		opts := &infoblox.Options{
			ReturnFields: []string{"name", "ipv6addr", "ttl", "view"},
		}
		body := &infoblox.RecordAAAAObject{
			Comment:  comment,
			Ipv6Addr: value,
			Name:     name,
			Ttl:      ttl,
			View:     view,
		}
		recID, err = client.RecordAAAA().Create(values, opts, body)
	case "CNAME":
		opts := &infoblox.Options{
			ReturnFields: []string{"name", "canonical", "ttl", "view"},
		}
		body := &infoblox.RecordCnameObject{
			Canonical: value,
			Comment:   comment,
			Name:      name,
			Ttl:       ttl,
			View:      view,
		}
		recID, err = client.RecordCname().Create(values, opts, body)
	case "HOST":
		if nextavailableip {
			value = "func:nextavailableip:" + value
		}
		opts := &infoblox.Options{
			ReturnFields: []string{"name", "ipv4addrs", "ttl", "view"},
		}
		body := &infoblox.RecordHostObject{
			Comment: comment,
			Ipv4Addrs: []infoblox.HostIpv4Addr{
				infoblox.HostIpv4Addr{
					Ipv4Addr: value,
				},
			},
			Name: name,
			Ttl:  ttl,
			View: view,
		}
		recID, err = client.RecordHost().Create(values, opts, body)
	default:
		return fmt.Errorf("resourceInfobloxRecordCreate: unknown type")
	}

	if err != nil {
		return fmt.Errorf("Failed to create Infblox Record: %s", err)
	}

	d.SetId(recID)

	log.Printf("[INFO] record ID: %s", d.Id())

	return resourceInfobloxRecordRead(d, meta)
}

func resourceInfobloxRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	var nextavailableip bool
	if attr, ok := d.GetOk("nextavailableip"); ok {
		nextavailableip = attr.(bool)
	}

	switch strings.ToUpper(d.Get("type").(string)) {
	case "A":
		opts := &infoblox.Options{
			ReturnFields: []string{"name", "ipv4addr", "ttl", "view"},
		}
		rec, err := client.GetRecordA(d.Id(), opts)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox A record: %s", err)
		}
		d.Set("fqdn", rec.Name)
		fqdn := strings.Split(rec.Name, ".")
		d.Set("name", fqdn[0])
		d.Set("domain", strings.Join(fqdn[1:], "."))
		d.Set("ipv4addr", rec.Ipv4Addr)
		d.Set("ttl", rec.Ttl)
		d.Set("type", "A")
		if !nextavailableip {
			d.Set("value", rec.Ipv4Addr)
		}
		d.Set("view", rec.View)
	case "AAAA":
		opts := &infoblox.Options{
			ReturnFields: []string{"name", "ipv6addr", "ttl", "view"},
		}
		rec, err := client.GetRecordAAAA(d.Id(), opts)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox AAAA record: %s", err)
		}
		d.Set("fqdn", rec.Name)
		fqdn := strings.Split(rec.Name, ".")
		d.Set("name", fqdn[0])
		d.Set("domain", strings.Join(fqdn[1:], "."))
		d.Set("ipv6addr", rec.Ipv6Addr)
		d.Set("ttl", rec.Ttl)
		d.Set("type", "AAAA")
		if !nextavailableip {
			d.Set("value", rec.Ipv6Addr)
		}
		d.Set("view", rec.View)
	case "CNAME":
		opts := &infoblox.Options{
			ReturnFields: []string{"name", "canonical", "ttl", "view"},
		}
		rec, err := client.GetRecordCname(d.Id(), opts)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox CNAME record: %s", err)
		}
		d.Set("fqdn", rec.Name)
		fqdn := strings.Split(rec.Name, ".")
		d.Set("name", fqdn[0])
		d.Set("domain", strings.Join(fqdn[1:], "."))
		d.Set("ttl", rec.Ttl)
		d.Set("type", "CNAME")
		d.Set("value", rec.Canonical)
		d.Set("view", rec.View)
	case "HOST":
		opts := &infoblox.Options{
			ReturnFields: []string{"name", "ipv4addrs", "ttl", "view"},
		}
		rec, err := client.GetRecordHost(d.Id(), opts)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox Host record: %s", err)
		}
		d.Set("fqdn", rec.Name)
		fqdn := strings.Split(rec.Name, ".")
		d.Set("name", fqdn[0])
		d.Set("domain", strings.Join(fqdn[1:], "."))
		if len(rec.Ipv4Addrs) > 0 {
			d.Set("ipv4addr", rec.Ipv4Addrs[0].Ipv4Addr)
			if !nextavailableip {
				d.Set("value", rec.Ipv4Addrs[0].Ipv4Addr)
			}
		}
		d.Set("ttl", rec.Ttl)
		d.Set("type", "Host")
		d.Set("view", rec.View)
	default:
		return fmt.Errorf("resourceInfobloxRecordRead: unknown type")
	}

	return nil
}

func resourceInfobloxRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	var comment, name, recID, value string
	var err, updateErr error
	var nextavailableip bool
	var ttl int
	values := url.Values{}

	switch strings.ToUpper(d.Get("type").(string)) {
	case "A":
		_, err = client.GetRecordA(d.Id(), nil)
	case "AAAA":
		_, err = client.GetRecordAAAA(d.Id(), nil)
	case "CNAME":
		_, err = client.GetRecordCname(d.Id(), nil)
	case "HOST":
		_, err = client.GetRecordHost(d.Id(), nil)
	default:
		return fmt.Errorf("resourceInfobloxRecordUpdate: unknown type")
	}

	if err != nil {
		return fmt.Errorf("Couldn't find Infoblox record: %s", err)
	}

	if attr, ok := d.GetOk("comment"); ok {
		comment = attr.(string)
	}
	if attr, ok := d.GetOk("name"); ok {
		name = attr.(string)
	}
	if attr, ok := d.GetOk("domain"); ok {
		name = strings.Join([]string{name, attr.(string)}, ".")
	}
	if attr, ok := d.GetOk("nextavailableip"); ok {
		nextavailableip = attr.(bool)
	}
	if attr, ok := d.GetOk("ttl"); ok {
		ttl = attr.(int)
	}
	if attr, ok := d.GetOk("value"); ok {
		value = attr.(string)
	}

	//log.Printf("[DEBUG] Infoblox Record update configuration: %#v", record)

	switch strings.ToUpper(d.Get("type").(string)) {
	case "A":
		if nextavailableip {
			if d.HasChange("value") {
				value = "func:nextavailableip:" + value
			} else if ipv4addr, ok := d.GetOk("ipv4addr"); ok {
				value = ipv4addr.(string)
			}
		}
		opts := &infoblox.Options{
			ReturnFields: []string{"name", "ipv4addr", "ttl", "view"},
		}
		body := &infoblox.RecordAObject{
			Comment:  comment,
			Ipv4Addr: value,
			Name:     name,
			Ttl:      ttl,
		}
		recID, updateErr = client.RecordAObject(d.Id()).Update(values, opts, body)
	case "AAAA":
		if nextavailableip {
			if d.HasChange("value") {
				value = "func:nextavailableip:" + value
			} else if ipv6addr, ok := d.GetOk("ipv6addr"); ok {
				value = ipv6addr.(string)
			}
		}
		opts := &infoblox.Options{
			ReturnFields: []string{"name", "ipv6addr", "ttl", "view"},
		}
		body := &infoblox.RecordAAAAObject{
			Comment:  comment,
			Ipv6Addr: value,
			Name:     name,
			Ttl:      ttl,
		}
		recID, updateErr = client.RecordAAAAObject(d.Id()).Update(values, opts, body)
	case "CNAME":
		opts := &infoblox.Options{
			ReturnFields: []string{"name", "canonical", "ttl", "view"},
		}
		body := &infoblox.RecordCnameObject{
			Canonical: value,
			Comment:   comment,
			Name:      name,
			Ttl:       ttl,
		}
		recID, updateErr = client.RecordCnameObject(d.Id()).Update(values, opts, body)
	case "HOST":
		if nextavailableip {
			if d.HasChange("value") {
				value = "func:nextavailableip:" + value
			} else if ipv4addr, ok := d.GetOk("ipv4addr"); ok {
				value = ipv4addr.(string)
			}
		}
		opts := &infoblox.Options{
			ReturnFields: []string{"name", "ipv4addrs", "ttl", "view"},
		}
		body := &infoblox.RecordHostObject{
			Comment: comment,
			Ipv4Addrs: []infoblox.HostIpv4Addr{
				infoblox.HostIpv4Addr{
					Ipv4Addr: value,
				},
			},
			Name: name,
			Ttl:  ttl,
		}
		recID, updateErr = client.RecordHostObject(d.Id()).Update(values, opts, body)
	default:
		return fmt.Errorf("resourceInfobloxRecordUpdate: unknown type")
	}

	if updateErr != nil {
		return fmt.Errorf("Failed to update Infblox Record: %s", err)
	}

	d.SetId(recID)

	return resourceInfobloxRecordRead(d, meta)
}

func resourceInfobloxRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	log.Printf("[INFO] Deleting Infoblox Record: %s, %s", d.Get("name").(string), d.Id())
	switch strings.ToUpper(d.Get("type").(string)) {
	case "A":
		_, err := client.GetRecordA(d.Id(), nil)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox A record: %s", err)
		}

		deleteErr := client.RecordAObject(d.Id()).Delete(nil)
		if deleteErr != nil {
			return fmt.Errorf("Error deleting Infoblox A Record: %s", err)
		}
	case "AAAA":
		_, err := client.GetRecordAAAA(d.Id(), nil)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox AAAA record: %s", err)
		}

		deleteErr := client.RecordAAAAObject(d.Id()).Delete(nil)
		if deleteErr != nil {
			return fmt.Errorf("Error deleting Infoblox AAAA Record: %s", err)
		}
	case "CNAME":
		_, err := client.GetRecordCname(d.Id(), nil)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox CNAME record: %s", err)
		}

		deleteErr := client.RecordCnameObject(d.Id()).Delete(nil)
		if deleteErr != nil {
			return fmt.Errorf("Error deleting Infoblox CNAME Record: %s", err)
		}
	case "HOST":
		_, err := client.GetRecordHost(d.Id(), nil)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox Host record: %s", err)
		}

		deleteErr := client.RecordHostObject(d.Id()).Delete(nil)
		if deleteErr != nil {
			return fmt.Errorf("Error deleting Infoblox Host Record: %s", err)
		}
	default:
		return fmt.Errorf("resourceInfobloxRecordDelete: unknown type")
	}
	return nil
}
