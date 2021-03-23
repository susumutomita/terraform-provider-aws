---
subcategory: "Greengrassv2"
layout: "aws"
page_title: "AWS: aws_greengrassv2_component"
description: |-
  Creates and manages an AWS IoT Greengrassv2 Component Definition
---

# Resource: aws_greengrassv2_component

Provides a Greengrassv2 component Job resource.

## Example Usage

## Example Usage for Json inline_recipe

```terraform
resource "aws_greengrassv2_component" "default" {
  tags = {
    tag-key = "tag-value"
  }
  inline_recipe = <<EOF
{
  "RecipeFormatVersion": "2020-01-25",
  "ComponentName": "com.example.test",
  "ComponentVersion": "1.0.0",
  "ComponentDescription": "sample",
  "ComponentConfiguration": {
    "DefaultConfiguration": {
      "Message": "sample"
    }
  },
  "Manifests": [
    {
      "Platform": {
        "os": "linux"
      },
      "Name": "Linux",
      "Lifecycle": {
        "Install": {
          "Script": "ls"
        },
        "Run": {
          "Script": "ls -l"
        }
      }
    }
  ]
}
EOF
}
```

## Argument Reference

The following arguments are supported:

* `inlineRecipe` – (Optional) The recipe to use to create the component. The recipe defines the component's metadata, parameters, dependencies, lifecycle, artifacts, and platform compatibility.You must specify either inlineRecipe or lambdaFunction.
* `lambdaFunction` – (Optional) The parameters to create a component from a Lambda function.You must specify either inlineRecipe or lambdaFunction.
* `tags` - (Optional) Key-value map of resource tags

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the Greengrassv2 Component.
* `id` - The ID of the Greengrassv2 Component.

## Import

IoT Greengrassv2 Component can be imported using the `arn`, e.g.

```
$ terraform import aws_greengrassv2_component.default <arn>
```
