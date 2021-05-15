package plan

import (
	"github.com/exepirit/cf-ddns/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPlan_AddNewDomain_Eval(t *testing.T) {
	var current []*domain.Endpoint
	desired := []*domain.Endpoint{
		{
			DNSName: "sub.example.com",
			Target:  domain.Target{"127.0.0.1"},
		},
	}
	plan := &Plan{Current: current, Desired: desired}

	plan.Eval()

	require.NotNil(t, plan.Changes)
	require.Len(t, plan.Changes.Create, 1)
	require.Len(t, plan.Changes.Delete, 0)
	require.Len(t, plan.Changes.Update, 0)
	require.Equal(t, desired[0], plan.Changes.Create[0])
}

func TestPlan_RemoveDomain_Eval(t *testing.T) {
	current := []*domain.Endpoint{
		{
			DNSName: "sub.example.com",
			Target:  domain.Target{"127.0.0.1"},
		},
	}
	var desired []*domain.Endpoint
	plan := &Plan{Current: current, Desired: desired}

	plan.Eval()

	require.NotNil(t, plan.Changes)
	require.Len(t, plan.Changes.Create, 0)
	require.Len(t, plan.Changes.Delete, 1)
	require.Len(t, plan.Changes.Update, 0)
	require.Equal(t, current[0], plan.Changes.Delete[0])
}

func TestPlan_ChangeDomainTarget_Eval(t *testing.T) {
	current := []*domain.Endpoint{
		{
			DNSName: "sub.example.com",
			Target:  domain.Target{"127.0.0.1"},
		},
	}
	desired := []*domain.Endpoint{
		{
			DNSName: "sub.example.com",
			Target:  domain.Target{"10.0.8.1"},
		},
	}
	plan := &Plan{Current: current, Desired: desired}

	plan.Eval()

	require.NotNil(t, plan.Changes)
	require.Len(t, plan.Changes.Create, 0, "unexpected create in diff")
	require.Len(t, plan.Changes.Delete, 0, "unexpected delete in diff")
	require.Len(t, plan.Changes.Update, 1, "unexpected update in diff")
	require.Equal(t, desired[0].Target, plan.Changes.Update[0].Target)
}

func TestPlan_EmptyChanges_Eval(t *testing.T) {
	var current, desired []*domain.Endpoint
	plan := &Plan{Current: current, Desired: desired}

	plan.Eval()

	require.NotNil(t, plan.Changes)
	require.Len(t, plan.Changes.Create, 0)
	require.Len(t, plan.Changes.Delete, 0)
	require.Len(t, plan.Changes.Update, 0)
}
