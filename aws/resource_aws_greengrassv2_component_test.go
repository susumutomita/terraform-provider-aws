package aws

import (
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go/service/greengrassv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAwsGreengrassv2Component_JsonFormat(t *testing.T) {

	var component greengrassv2.DescribeComponentOutput
	resourceName := "aws_greengrassv2_component.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		ErrorCheck:   testAccErrorCheck(t, greengrassv2.EndpointsID),
		CheckDestroy: testAccCheckComponentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGreengrassv2ComponentConfigJson(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGreengrassv2ComponentExists(resourceName, &component),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr("aws_greengrassv2_component.test", "tags.Name", "tagValue"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"inline_recipe"}, // This doesn't work because the API doesn't provide attachments info directly
			},
		},
	})
}

func TestAccAwsGreengrassv2Component_YamlFormat(t *testing.T) {

	var component greengrassv2.DescribeComponentOutput
	resourceName := "aws_greengrassv2_component.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		ErrorCheck:   testAccErrorCheck(t, greengrassv2.EndpointsID),
		CheckDestroy: testAccCheckComponentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGreengrassv2ComponentConfigYaml(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGreengrassv2ComponentExists(resourceName, &component),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr("aws_greengrassv2_component.test", "tags.Name", "tagValue"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"inline_recipe"}, // This doesn't work because the API doesn't provide attachments info directly
			},
		},
	})
}

func TestAccAwsGreengrassv2Component_Lambda(t *testing.T) {

	var component greengrassv2.DescribeComponentOutput
	resourceName := "aws_greengrassv2_component.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		ErrorCheck:   testAccErrorCheck(t, greengrassv2.EndpointsID),
		CheckDestroy: testAccCheckComponentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGreengrassv2ComponentConfigLambda(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGreengrassv2ComponentExists(resourceName, &component),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr("aws_greengrassv2_component.test", "tags.Name", "tagValue"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"inline_recipe", "lambda_function"}, // This doesn't work because the API doesn't provide attachments info directly
			},
		},
	})
}

func testAccCheckGreengrassv2ComponentExists(n string, component *greengrassv2.DescribeComponentOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).greengrassv2conn
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		describeComponentOpts := &greengrassv2.DescribeComponentInput{
			Arn: aws.String(rs.Primary.ID),
		}
		resp, err := conn.DescribeComponent(describeComponentOpts)
		if err != nil {
			return err
		}

		*component = *resp

		return nil
	}
}

func testAccGreengrassv2ComponentConfigJson(rName string) string {
	return fmt.Sprintf(`
resource "aws_greengrassv2_component" "test" {
  tags = {
    Name  = "tagValue"
  }
  inline_recipe = jsonencode(
    {
      "RecipeFormatVersion" : "2020-01-25",
      "ComponentName" : "com.example.test.json.%s",
      "ComponentVersion" : "1.0.0",
      "ComponentType" : "aws.greengrass.generic",
      "ComponentDescription" : "sample",
      "ComponentConfiguration" : {
        "DefaultConfiguration" : {
          "Message" : "sample"
        }
      },
      "Manifests" : [
        {
          "Platform" : {
            "os" : "linux"
          },
          "Name" : "Linux",
          "Lifecycle" : {
            "Install" : {
              "Script" : "ls"
            },
            "Run" : {
              "Script" : "ls -l"
            }
          }
        }
      ],
    }
	)
}
`, rName)
}

func testAccGreengrassv2ComponentConfigYaml(rName string) string {
	return fmt.Sprintf(`
resource "aws_greengrassv2_component" "test" {
  tags = {
    Name = "tagValue"
  }
  inline_recipe          = <<EOF
---
RecipeFormatVersion: '2020-01-25'
ComponentName: "com.example.test.yaml.%s"
ComponentVersion: 1.0.0
ComponentType: aws.greengrass.generic
ComponentDescription: sample
ComponentConfiguration:
  DefaultConfiguration:
    Message: sample
Manifests:
- Platform:
    os: linux
  Name: Linux
  Lifecycle:
    Install:
      Script: ls
    Run:
      Script: ls -l
Lifecycle: {}
EOF
}
`, rName)
}

func testAccGreengrassv2ComponentConfigLambda(rName string) string {
	return fmt.Sprintf(`
resource "aws_greengrassv2_component" "test" {
  tags = {
    Name = "tagValue"
  }
  lambda_function {
    component_dependencies {
      component_name      = "aws.greengrass.test"
      dependency_type     = "SOFT"
      version_requirement = "1.0.0"
    }
		component_dependencies {
      component_name      = "aws.greengrass.test2"
      dependency_type     = "HARD"
      version_requirement = ">1.0.0"
    }
		component_platforms {
      attributes {
        os           = "Linux"
        architecture = "arm"
      }
      name = "test Linux platform"
    }
    component_platforms {
      attributes {
        os           = "Windows"
        architecture = "x86"
      }
      name = "test Windows platform"
    }
		component_lambda_parameters {
      max_idle_time_in_seconds    = 2147483647
      max_instances_count         = 1
      max_queue_size              = 1
      status_timeout_in_seconds   = 2147483647
      timeout_in_seconds          = 10
      input_payload_encoding_type = "binary"
      exec_args                   = ["hoge", "fuga"]
      environment_variables = {
        hoge   = "hoge"
        number = 1
      }
      event_sources {
        topic = aws_sns_topic.test.arn
        type  = "IOT_CORE"
      }
      event_sources {
        topic = aws_sns_topic.test2.arn
        type  = "PUB_SUB"
      }
			linux_process_params {
        container_params {
          devices {
            add_group_owner = true
            path            = "/dev/stdout"
            permission      = "ro"
          }
          memory_size_in_kb = 2048
          mount_ro_sysfs    = true
          volumes {
            add_group_owner  = true
            destination_path = "/tmp"
            permission       = "ro"
            source_path      = "/tmp"
          }
        }
        isolation_mode = "GreengrassContainer"
      }
    }

    component_name    = "com.example.test.lambda.%s"
    component_version = "1.0.0"
    lambda_arn        = aws_lambda_function.test.qualified_arn
  }
}

resource "aws_iam_role_policy" "iam_policy_for_lambda" {
  name = "hoge"
  role = aws_iam_role.iam_for_lambda.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "hoge"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}
resource "aws_lambda_function" "test" {
  filename      = "test-fixtures/lambdatest.zip"
  function_name = "hoge"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "exports.example"
  runtime       = "nodejs12.x"
	publish       = "true"
}

resource "aws_sns_topic" "test" {
  name = "lambda-test-topic"
}

resource "aws_sns_topic" "test2" {
  name = "lambda-test-topic2"
}
`, rName)
}

func testAccCheckComponentDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).greengrassv2conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_greengrassv2_component" {
			continue
		}

		describeComponentOpts := &greengrassv2.DescribeComponentInput{
			Arn: aws.String(rs.Primary.ID),
		}

		resp, err := conn.DescribeComponent(describeComponentOpts)
		if err == nil {
			if len(aws.StringValue(resp.Arn)) > 0 {
				return fmt.Errorf("Greengrassv2 component still exists.")
			}
			return nil
		}
	}

	return nil
}

func testAccCheckComponentExists(n string, app *elasticbeanstalk.ApplicationVersionDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		fmt.Errorf("Greengrassv2 component is not set")

		if rs.Primary.ID == "" {
			return fmt.Errorf("Greengrassv2 component is not set")
		}

		conn := testAccProvider.Meta().(*AWSClient).greengrassv2conn
		describeComponentOpts := greengrassv2.DescribeComponentInput{
			Arn: aws.String(rs.Primary.ID),
		}

		log.Printf("[DEBUG] Greengrassv2 component TEST describe opts: %s", describeComponentOpts)

		resp, err := conn.DescribeComponent(&describeComponentOpts)
		if err != nil {
			return err
		}
		if len(aws.StringValue(resp.Arn)) == 0 {
			return fmt.Errorf("Greengrassv2 component not found.")
		}

		return nil
	}
}
