package gsuite

import (
  "fmt"
  "strings"

  "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
  directory "google.golang.org/api/admin/directory/v1"
)

func dataUserValues() *schema.Resource {
  return &schema.Resource{
    Read: dataUserValuesRead,
    Schema: map[string]*schema.Schema{
      "primary_email": {
        Type:     schema.TypeString,
        Required: true,
        StateFunc: func(val interface{}) string {
          return strings.ToLower(val.(string))
        },
      },

      "custom_schemas": {
        Type:     schema.TypeList,
        Computed: true,
        Elem: &schema.Resource{
          Schema: map[string]*schema.Schema{
            "name": {
              Computed: true,
              Type: schema.TypeString,
            },
            "value": {
              Computed: true,
              Type: schema.TypeString,
            },
          },
        },
      },
    },
  }
}

func dataUserValuesRead(d *schema.ResourceData, meta interface{}) error {
  config := meta.(*Config)

  var user *directory.User
  var err error
  err = retry(func() error {
    user, err = config.directory.Users.Get(d.Get("primary_email").(string)).Projection("full").Do()
    return err
  }, config.TimeoutMinutes)

  if err != nil {
    return handleNotFoundError(err, d, fmt.Sprintf("User %q", d.Id()))
  }

  u_err, u_msg := flattenCustomSchema(user.CustomSchemas)
  if u_err != nil {
    return handleNotFoundError(u_err, d, fmt.Sprintf("Flatten schema error with user %q", d.Id()))
  }

  d.SetId(user.Id)
  d.Set("custom_schemas", u_msg)

  return nil
}
