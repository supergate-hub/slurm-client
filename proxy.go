package slurmclient

// SlurmProxy provides access to Slurm API resources.
// Each method returns a service interface for a specific resource type.
//
// Usage: client.Slurm.Jobs().List(ctx)
type SlurmProxy interface {
	Jobs() JobsService
	Nodes() NodesService
	Partitions() PartitionsService
	Reservations() ReservationsService
	Licenses() LicensesService
	Shares() SharesService
	Diag() DiagService
	Ping() PingService
	Reconfigure() ReconfigureService
}

// SlurmdbProxy provides access to SlurmDB API resources.
//
// Usage: client.Slurmdb.Accounts().List(ctx)
type SlurmdbProxy interface {
	Accounts() AccountsService
	Users() UsersService
	Associations() AssociationsService
	Clusters() ClustersService
	QOS() QOSService
	Jobs() SlurmdbJobsService
	WCKeys() WCKeysService
	TRES() TRESService
	Instances() InstancesService
	Config() ConfigService
	SlurmdbDiag() SlurmdbDiagService
}
