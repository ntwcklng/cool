package types

type Deployment struct {
	ID              int    `json:"id"`
	ApplicationName string `json:"name"`
	DeploymentUUID  string `json:"uuid"`
	FQDN            string `json:"fqdn"`
}
