---
subcategory: "Auto Scaling"
---

# g42cloud_as_configuration

Manages an AS Configuration resource within G42Cloud.

## Example Usage

### Basic AS Configuration

```hcl
resource "g42cloud_as_configuration" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config {
    flavor = var.flavor
    image  = var.image_id

    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }
    key_name  = var.keyname
    user_data = file("userdata.txt")
  }
}
```

### AS Configuration With Encrypted Data Disk

```hcl
resource "g42cloud_as_configuration" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config {
    flavor = var.flavor
    image  = var.image_id

    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }

    disk {
      size        = 100
      volume_type = "SSD"
      disk_type   = "DATA"
      kms_id      = var.kms_id
    }
    key_name  = var.keyname
    user_data = file("userdata.txt")
  }
}
```

### AS Configuration With User Data and Metadata

```hcl
resource "g42cloud_as_configuration" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config {
    flavor = var.flavor
    image  = var.image_id

    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }

    key_name  = var.keyname
    user_data = file("userdata.txt")
    metadata  = {
      some_key = "some_value"
    }
  }
}
```

### AS Configuration uses the existing instance specifications as the template

```hcl
resource "g42cloud_as_configuration" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config {
    instance_id = "4579f2f5-cbe8-425a-8f32-53dcb9d9053a"
    key_name    = var.keyname
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the AS configuration.
  If omitted, the `region` argument of the provider is used. Changing this will create a new resource.

* `scaling_configuration_name` - (Required, String, ForceNew) Specifies the AS configuration name.
  The name contains only letters, digits, underscores (_), and hyphens (-), and cannot exceed 64 characters.
  Changing this will create a new resource.

* `instance_config` - (Required, List, ForceNew) Specifies the information about instance configuration.
  The object structure is documented below. Changing this will create a new resource.

The `instance_config` block supports:

* `instance_id` - (Optional, String, ForceNew) Specifies the ECS instance ID when using its specification
  as the template to create AS configurations. In this case, `flavor`, `image`, and `disk` arguments do not take effect.
  If this argument is not specified, `flavor`, `image`, and `disk` arguments are mandatory.
  Changing this will create a new resource.

* `flavor` - (Optional, String, ForceNew) Specifies the ECS flavor name. A maximum of 10 flavors can be selected.
  Use a comma (,) to separate multiple flavor names. Changing this will create a new resource.

* `image` - (Optional, String, ForceNew) Specifies the ECS image ID. Changing this will create a new resource.

* `disk` - (Optional, List, ForceNew) Specifies the disk group information. System disks are mandatory and
  data disks are optional. The [object](#instance_config_disk_object) structure is documented below.
  Changing this will create a new resource.

* `key_name` - (Required, String, ForceNew) Specifies the name of the SSH key pair used to log in to the instance.
  Changing this will create a new resource.

* `user_data` - (Optional, String, ForceNew) Specifies the user data to provide when launching the instance.
  The file content must be encoded with Base64. Changing this will create a new resource.

* `public_ip` - (Optional, List, ForceNew) Specifies the EIP of the ECS instance.
  The [object](#instance_config_public_ip_object) structure is documented below.
  Changing this will create a new resource.

* `metadata` - (Optional, Map, ForceNew) Specifies the key/value pairs to make available from within the instance.
  Changing this will create a new resource.

* `personality` - (Optional, List, ForceNew) Specifies the customize personality of an instance by defining one or
  more files and their contents. The [object](#instance_config_personality_object) structure is documented below.
  Changing this will create a new resource.

<a name="instance_config_disk_object"></a>
The `disk` block supports:

* `size` - (Required, Int, ForceNew) Specifies the disk size. The unit is GB.
  The system disk size ranges from 40 to 32768, and the data disk size ranges from 10 to 32768.
  Changing this will create a new resource.

* `volume_type` - (Required, String, ForceNew) Specifies the disk type. Changing this will create a new resource.
  Available options are:
  + `SAS`: high I/O disk type.
  + `SSD`: ultra-high I/O disk type.
  + `GPSSD`: general purpose SSD disk type.

* `disk_type` - (Required, String, ForceNew) Specifies whether the disk is a system disk or a data disk.
  Option **DATA** indicates a data disk, option **SYS** indicates a system disk.
  Changing this will create a new resource.

* `kms_id` - (Optional, String, ForceNew) Specifies the encryption KMS ID of the **DATA** disk.
  Changing this will create a new resource.

<a name="instance_config_public_ip_object"></a>
The `public_ip` block supports:

* `eip` - (Required, List, ForceNew) Specifies the EIP configuration that will be automatically assigned to the instance.
  The object structure is documented below. Changing this will create a new resource.

The `eip` block supports:

* `ip_type` - (Required, String, ForceNew) Specifies the EIP type. Possible values are **5_bgp** (dynamic BGP)
  and **5_sbgp** (static BGP). Changing this will create a new resource.

* `bandwidth` - (Required, List, ForceNew) Specifies the bandwidth information. The object structure is documented below.
  Changing this will create a new resource.

The `bandwidth` block supports:

* `share_type` - (Required, String, ForceNew) Specifies the bandwidth sharing type. The system only supports
  **PER** (indicates exclusive bandwidth). Changing this will create a new resource.

* `charging_mode` - (Required, String, ForceNew) Specifies whether the bandwidth is billed by traffic or by bandwidth
  size. The value can be **traffic** or **bandwidth**. Changing this creates a new resource.

* `size` - (Required, Int, ForceNew) Specifies the bandwidth (Mbit/s). The value range for bandwidth billed by bandwidth
  is 1 to 2000 and that for bandwidth billed by traffic is 1 to 300.
  Changing this creates a new resource.

<a name="instance_config_personality_object"></a>
The `personality` block supports:

* `path` - (Required, String, ForceNew) Specifies the path of the injected file. Changing this creates a new resource.

* `contents` - (Required, String, ForceNew) Specifies the content of the injected file, which must be encoded with base64.
  Changing this creates a new resource.
