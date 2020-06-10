package configuration

const (
	gwDeployEnv              = "GW_DEPLOY_ENV"
	gwDbPasswordDev          = "GW_DATABASE_SECRET_DEV"
	gwDbPasswordQA           = "GW_DATABASE_SECRET_QA"
	gwDbPasswordLocal        = "GW_DATABASE_SECRET_LOCAL"
	devEnvironment           = "DEV"
	qaEnvironment            = "QA"
	configDev                = "config_dev"
	configQA                 = "config_qa"
	configLocal              = "config_local"
	configTypeYaml           = "yaml"
	configPathFromRootDir    = "src/configuration"
	configPathFromSubModules = "../configuration"
	cypherPassPhaseLocal     = "CYPHER_PASS_PHASE_LOCAL"
	cypherPassPhaseDev       = "CYPHER_PASS_PHASE_DEV"
	cypherPassPhaseQA        = "CYPHER_PASS_PHASE_QA"
)
