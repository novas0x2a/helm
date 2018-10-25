/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"io"
	"testing"

	"github.com/spf13/cobra"

	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	rpb "k8s.io/helm/pkg/proto/hapi/release"
)

func TestRollbackCmd(t *testing.T) {
	mk := func(name string, vers int32, code rpb.Status_Code, appVersion string) *rpb.Release {
		ch := &chart.Chart{
			Metadata: &chart.Metadata{
				Name:       name,
				Version:    "0.1.0-beta.1",
				AppVersion: appVersion,
			},
		}

		return helm.ReleaseMock(&helm.MockReleaseOptions{
			Name:       name,
			Chart:      ch,
			Version:    vers,
			StatusCode: code,
		})
	}

	mkRels := func() []*rpb.Release {
		return []*rpb.Release{
			mk("funny-honey", 1, rpb.Status_DEPLOYED, "1.1"),
		}
	}

	tests := []releaseCase{
		{
			name:     "rollback a release",
			args:     []string{"funny-honey", "1"},
			rels:     mkRels(),
			expected: "Rollback was a success.",
		},
		{
			name:     "rollback a release with timeout",
			args:     []string{"funny-honey", "1"},
			flags:    []string{"--timeout", "120"},
			rels:     mkRels(),
			expected: "Rollback was a success.",
		},
		{
			name:     "rollback a release with wait",
			args:     []string{"funny-honey", "1"},
			flags:    []string{"--wait"},
			rels:     mkRels(),
			expected: "Rollback was a success.",
		},
		{
			name:     "rollback a release with description",
			args:     []string{"funny-honey", "1"},
			flags:    []string{"--description", "foo"},
			rels:     mkRels(),
			expected: "Rollback was a success.",
		},
		{
			name: "rollback a release without revision",
			args: []string{"funny-honey"},
			rels: mkRels(),
			err:  true,
		},
	}

	cmd := func(c *helm.FakeClient, out io.Writer) *cobra.Command {
		return newRollbackCmd(c, out)
	}

	runReleaseCases(t, tests, cmd)

}
