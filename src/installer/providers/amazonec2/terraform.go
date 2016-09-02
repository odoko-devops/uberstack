package amazonec2

import (
	"installer/model"
	"log"
	"utils"
)

func (aws Amazonec2) applyTerraform(config model.Config, provider model.Provider) {
	log.Println("Create AWS VPC Environment")

	cwd := "terraform/aws"
	env := utils.Environment{
		"TF_VAR_aws_access_key": aws.accessKey,
		"TF_VAR_aws_secret_key": aws.secretKey,
	}

	utils.Execute("terraform apply -state=/state/terraform.tfstate", env, cwd)
	aws.vpcId = utils.ExecuteAndRetrieve("terraform output -state=/state/terraform.tfstate vpc_id", env, cwd)
	aws.subnetId = utils.ExecuteAndRetrieve("terraform output -state=/state/terraform.tfstate subnet_id", env, cwd)
}

func (aws Amazonec2) destroyTerraform(config model.Config, provider model.Provider) {
	log.Println("Destroy AWS VPC Environment")

	cwd := "terraform/aws"
	env := utils.Environment{
		"TF_VAR_aws_access_key": aws.accessKey,
		"TF_VAR_aws_secret_key": aws.secretKey,
	}

	utils.ExecuteAndRetrieve("terraform destroy -state=/state/terraform.tfstate -force", env, cwd)
}