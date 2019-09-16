package nexus

import (
	"bytes"
	"fmt"
	"path"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNexusComponentRaw() *schema.Resource {
	return &schema.Resource{
		Create: resourceNexusComponentRawCreate,
		Read:   resourceNexusComponentRawRead,
		Update: resourceNexusComponentRawUpdate,
		Delete: resourceNexusComponentRawDelete,

		Schema: map[string]*schema.Schema{
			"repository": {
				Type:        schema.TypeString,
				Description: "The repository to store the component in",
				Required:    true,
			},
			"src": {
				Type:        schema.TypeString,
				Description: "Path or URI of source file to upload",
				Required:    true,
			},
			"filename": {
				Type:        schema.TypeString,
				Description: "The filename to assign to the file. If empty, will use the base filename of the uploaded src file",
				Optional:    true,
			},

			"dest": {
				Type:        schema.TypeString,
				Description: "Destination for uploaded files (e.g. /path/to/files/)",
				Required:    true,
			},
		},
	}
}

func resourceNexusComponentRawCreate(d *schema.ResourceData, meta interface{}) error {
	fmt.Println("resourceNexusComponentRawCreate")
	client := meta.(*client)

	repository := d.Get("repository").(string)
	src := d.Get("src").(string)
	filename := d.Get("filename").(string)
	dest := d.Get("dest").(string)

	if filename == "" {
		filename = path.Base(src)
	}

	assetPath := path.Join(dest, filename)

	exists, err := client.FileExists(repository, assetPath)
	if err != nil {
		return err
	}

	if !exists {
		fmt.Println("Artifact does not exist. Uploading")
		b, err := getFileContents(src)
		if err != nil {
			return err
		}
		if err := client.Put(repository, assetPath, bytes.NewReader(b)); err != nil {
			return err
		}
	}

	fullPath := path.Join(repository, assetPath)

	d.SetId(fullPath)
	return resourceNexusComponentRawRead(d, meta)
}

func resourceNexusComponentRawRead(d *schema.ResourceData, meta interface{}) error {
	fmt.Println("resourceNexusComponentRawRead")
	return nil
}

func resourceNexusComponentRawUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
	//d.SetId(config.ID)
	//return resourceNexusComponentRawRead(d, meta)
}

func resourceNexusComponentRawDelete(d *schema.ResourceData, meta interface{}) error {
	fmt.Println("resourceNexusComponentRawDelete")
	client := meta.(*client)
	repository := d.Get("repository").(string)
	src := d.Get("src").(string)
	filename := d.Get("filename").(string)
	dest := d.Get("dest").(string)

	if filename == "" {
		filename = path.Base(src)
	}

	assetPath := path.Join(dest, filename)

	if err := client.Delete(repository, assetPath); err != nil {
		fmt.Println("Delete Error: ", err)
		return err
	}
	d.SetId("")
	return nil
}
