package manager

import (
	"loadbalancer/APIRequest"
	"reflect"
	"testing"

	"gotest.tools/v3/assert"
)

func TestLoadBalancerManagerApplyConfig(t *testing.T) {
	t.Run("Test apply config", func(t *testing.T) {
		m := NewManager("ROUND_ROBIN", "TOKEN_BUCKET")
		m.ApplyConfigs("some endpoint", []string{"3.3.3.3", "4.4.4.4"})

		output := m.serverIPMapping.mapping
		expected := map[string][]string{
			"some endpoint": []string{"3.3.3.3", "4.4.4.4"},
		}

		if !reflect.DeepEqual(output, expected) {
			t.Fatalf("expected %v, got %v", expected, output)
		}

	})

}

func TestLoadBalancerMatcher(t *testing.T) {
	t.Run("Test server matching", func(t *testing.T) {
		m := NewManager("ROUND_ROBIN", "TOKEN_BUCKET")
		m.ApplyConfigs("some endpoint", []string{"3.3.3.3", "4.4.4.4"})
		req := APIRequest.NewAPIRequest("1.1.1.1", "2.2.2.2", "some endpoint")
		err := m.matcher.AssignServerIP(req)
		assert.Equal(t, err, nil)

		// verify destination IP of request changed to 3.3.3.3
		if req.DestinationIP != "3.3.3.3" {
			t.Fatalf("expected destination server ip"+
				" to be 3.3.3.3, got %v", req.DestinationIP)
		}

		err = m.matcher.AssignServerIP(req)
		assert.Equal(t, err, nil)
		// this time destIP should be 4.4.4.4
		if req.DestinationIP != "4.4.4.4" {
			t.Fatalf("expected destination server ip"+
				" to be 4.4.4.4, got %v", req.DestinationIP)
		}

		// reconfigure load balancer again
		m.ApplyConfigs("some endpoint", []string{})
		err = m.matcher.AssignServerIP(req)
		assert.Equal(t, err.Error(), "no servers available for APIName some endpoint")

	})

}
