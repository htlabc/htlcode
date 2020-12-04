package memchange

const (
	FAIL    = 0
	SUCCESS = 1
)

type ClusterMemberChageResult struct {
	status     int
	leaderHint string
}

func NewClusterMemberChageResult() *ClusterMemberChageResult {
	return &ClusterMemberChageResult{}
}

func (cmcr *ClusterMemberChageResult) STATUS(val int) *ClusterMemberChageResult {
	cmcr.status = val
	return cmcr
}

func (cmcr *ClusterMemberChageResult) LEADERHINT(val string) *ClusterMemberChageResult {
	cmcr.leaderHint = val
	return cmcr
}
