package slurmclient

// slurmProxy implements SlurmProxy with pre-created service instances.
type slurmProxy struct {
	jobs         JobsService
	nodes        NodesService
	partitions   PartitionsService
	reservations ReservationsService
	licenses     LicensesService
	shares       SharesService
	diag         DiagService
	ping         PingService
	reconfigure  ReconfigureService
}

func newSlurmProxy(t *Transport) SlurmProxy {
	return &slurmProxy{
		jobs:         &jobsServiceImpl{t: t},
		nodes:        &nodesServiceImpl{t: t},
		partitions:   &partitionsServiceImpl{t: t},
		reservations: &reservationsServiceImpl{t: t},
		licenses:     &licensesServiceImpl{t: t},
		shares:       &sharesServiceImpl{t: t},
		diag:         &diagServiceImpl{t: t},
		ping:         &pingServiceImpl{t: t},
		reconfigure:  &reconfigureServiceImpl{t: t},
	}
}

func (p *slurmProxy) Jobs() JobsService               { return p.jobs }
func (p *slurmProxy) Nodes() NodesService             { return p.nodes }
func (p *slurmProxy) Partitions() PartitionsService    { return p.partitions }
func (p *slurmProxy) Reservations() ReservationsService { return p.reservations }
func (p *slurmProxy) Licenses() LicensesService        { return p.licenses }
func (p *slurmProxy) Shares() SharesService            { return p.shares }
func (p *slurmProxy) Diag() DiagService                { return p.diag }
func (p *slurmProxy) Ping() PingService                { return p.ping }
func (p *slurmProxy) Reconfigure() ReconfigureService  { return p.reconfigure }
