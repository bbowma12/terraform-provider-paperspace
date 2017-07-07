package main

import (
  "encoding/json"
  "fmt"
  "github.com/hashicorp/terraform/helper/schema"
  "log"
  "reflect"
)

type MapIf map[string]interface{}

func (m *MapIf) Append(d *schema.ResourceData, k string) {
  v := d.Get(k)
  (*m)[k] = v
}

func (m *MapIf) AppendIfSet(d *schema.ResourceData, k string) {
  v := d.Get(k)
  if reflect.ValueOf(v).Interface() != reflect.Zero(reflect.TypeOf(v)).Interface() {
    (*m)[k] = v
  }
}

func SetResDataFrom(d *schema.ResourceData, m map[string]interface{}, dn, n string) {
  v, ok := m[n]
  //log.Printf("%v %v\n", n, v)
  if ok {
    d.Set(dn, v)
  }
}

func SetResData(d *schema.ResourceData, m map[string]interface{}, n string) {
  SetResDataFrom(d, m, n, n)
}

func resourceMachineCreate(d *schema.ResourceData, m interface{}) error {
  client := m.(PaperspaceClient).RestyClient

  log.Printf("[INFO] paperspace resourceMachineCreate Client ready")

  body := make(MapIf)
  body.Append(d, "region")
  body.Append(d, "machineType")
  body.Append(d, "size")
  body.Append(d, "billingType")
  body.Append(d, "machineName")
  body.Append(d, "templateId")
  body.AppendIfSet(d, "networkId")
  body.AppendIfSet(d, "teamId")
  body.AppendIfSet(d, "userId")
  body.AppendIfSet(d, "email")
  body.AppendIfSet(d, "password")
  body.AppendIfSet(d, "firstName")
  body.AppendIfSet(d, "lastName")
  body.AppendIfSet(d, "notificationEmail")

  data, _ := json.MarshalIndent(body, "", "  ")
  log.Println(string(data))

  resp, err := client.R().
  SetBody(body).
  Post("/machines/createSingleMachinePublic")

  if err != nil {
    return fmt.Errorf("Error creating paperspace machine: %s", err)
  }

  statusCode := resp.StatusCode()
  log.Printf("[INFO] paperspace resourceMachineCreate StatusCode: %v", statusCode)
  LogResponse("paperspace resourceMachineCreate", resp, err)
  if statusCode != 200 {
    return fmt.Errorf("Error creating paperspace machine: Response: %s", resp.Body())
  }

  var f interface{}
  err = json.Unmarshal(resp.Body(), &f)

  /*fake := []byte(`{"id":"psmfffm3","name":"Tom Terraform Test 4","os":null,"ram":null,
    "cpus":1,"gpu":null,"storageTotal":null,"storageUsed":null,"usageRate":"C1 Hourly",
    "shutdownTimeoutInHours":null,"shutdownTimeoutForces":false,"performAutoSnapshot":false,
    "autoSnapshotFrequency":null,"autoSnapshotSaveCount":null,"agentType":"LinuxHeadless",
    "dtCreated":"2017-06-22T04:29:59.501Z","state":"provisioning","networkId":null,
    "privateIpAddress":null,"publicIpAddress":null,"region":null,"userId":"uijn3il",
    "teamId":null}`)
  err := json.Unmarshal(fake, &f)*/

  if err != nil {
    return fmt.Errorf("Error unmarshalling paperspace machine create response: %s", err)
  }

  mp := f.(map[string]interface{})
  id, _ := mp["id"].(string)

  if id == "" {
    return fmt.Errorf("Error in paperspace machine create data: id not found")
  }

  log.Printf("[INFO] paperspace resourceMachineCreate returned id: %v", id)

  SetResDataFrom(d, mp, "machineName", "name") //overlays, but should be the same
  SetResData(d, mp, "name") //duplicate of the above
  SetResData(d, mp, "os")
  SetResData(d, mp, "ram")
  SetResData(d, mp, "cpus")
  SetResData(d, mp, "gpu")
  SetResData(d, mp, "storageTotal")
  SetResData(d, mp, "storageUsed")
  SetResData(d, mp, "usageRate")
  SetResData(d, mp, "shutdownTimeoutInHours")
  SetResData(d, mp, "shutdownTimeoutForces")
  SetResData(d, mp, "performAutoSnapshot")
  SetResData(d, mp, "autoSnapshotFrequency")
  SetResData(d, mp, "autoSnapshotSaveCount")
  SetResData(d, mp, "agentType")
  SetResData(d, mp, "dtCreated")
  SetResData(d, mp, "state")
  SetResData(d, mp, "networkId") //overlays with null initially
  SetResData(d, mp, "privateIpAddress")
  SetResData(d, mp, "publicIpAddress")
  SetResData(d, mp, "region") //overlays with null initially
  SetResData(d, mp, "userId")
  SetResData(d, mp, "teamId")

  d.SetId(id);

  return nil
}

func resourceMachineRead(d *schema.ResourceData, m interface{}) error {
  return nil
}

func resourceMachineUpdate(d *schema.ResourceData, m interface{}) error {
  return nil
}

func resourceMachineDelete(d *schema.ResourceData, m interface{}) error {
  return nil
}

func resourceMachine() *schema.Resource {
  return &schema.Resource{
    Create: resourceMachineCreate,
    Read:   resourceMachineRead,
    Update: resourceMachineUpdate,
    Delete: resourceMachineDelete,

    Schema: map[string]*schema.Schema{
      "region": &schema.Schema{
          Type:     schema.TypeString,
          Required: true,
      },
      "machineType": &schema.Schema{
          Type:     schema.TypeString,
          Required: true,
      },
      "size": &schema.Schema{
          Type:     schema.TypeInt,
          Required: true,
      },
      "billingType": &schema.Schema{
          Type:     schema.TypeString,
          Required: true,
      },
      "machineName": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
      },
      "templateId": &schema.Schema{
          Type:     schema.TypeString,
          Required: true,
      },
      "networkId": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "teamId": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "userId": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "email": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "password": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "firstName": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "lastName": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "notificationEmail": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "name": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "os": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "ram": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "cpus": &schema.Schema{
          Type:     schema.TypeInt,
          Computed: true,
      },
      "gpu": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "storageTotal": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "storageUsed": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "usageRate": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "shutdownTimeoutInHours": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "shutdownTimeoutForces": &schema.Schema{
          Type:     schema.TypeBool,
          Computed: true,
      },
      "performAutoSnapshot": &schema.Schema{
          Type:     schema.TypeBool,
          Computed: true,
      },
      "autoSnapshotFrequency": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "autoSnapshotSaveCount": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "agentType": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "dtCreated": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "state": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "privateIpAddress": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "publicIpAddress": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
    },
  }
}
