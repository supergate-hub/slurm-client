package slurmclient

// slurmdbProxy implements SlurmdbProxy with pre-created service instances.
type slurmdbProxy struct {
	accounts     AccountsService
	users        UsersService
	associations AssociationsService
	clusters     ClustersService
	qos          QOSService
	jobs         SlurmdbJobsService
	wckeys       WCKeysService
	tres         TRESService
	instances    InstancesService
	config       ConfigService
	diag         SlurmdbDiagService
}

func newSlurmdbProxy(t *Transport) SlurmdbProxy {
	return &slurmdbProxy{
		accounts:     &accountsServiceImpl{t: t},
		users:        &usersServiceImpl{t: t},
		associations: &associationsServiceImpl{t: t},
		clusters:     &clustersServiceImpl{t: t},
		qos:          &qosServiceImpl{t: t},
		jobs:         &slurmdbJobsServiceImpl{t: t},
		wckeys:       &wckeysServiceImpl{t: t},
		tres:         &tresServiceImpl{t: t},
		instances:    &instancesServiceImpl{t: t},
		config:       &configServiceImpl{t: t},
		diag:         &slurmdbDiagServiceImpl{t: t},
	}
}

func (p *slurmdbProxy) Accounts() AccountsService       { return p.accounts }
func (p *slurmdbProxy) Users() UsersService             { return p.users }
func (p *slurmdbProxy) Associations() AssociationsService { return p.associations }
func (p *slurmdbProxy) Clusters() ClustersService       { return p.clusters }
func (p *slurmdbProxy) QOS() QOSService                 { return p.qos }
func (p *slurmdbProxy) Jobs() SlurmdbJobsService        { return p.jobs }
func (p *slurmdbProxy) WCKeys() WCKeysService           { return p.wckeys }
func (p *slurmdbProxy) TRES() TRESService               { return p.tres }
func (p *slurmdbProxy) Instances() InstancesService     { return p.instances }
func (p *slurmdbProxy) Config() ConfigService           { return p.config }
func (p *slurmdbProxy) SlurmdbDiag() SlurmdbDiagService { return p.diag }
