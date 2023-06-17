/*
Copyright 2023 The Radius Authors.

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

package create

import (
	"context"

	corerp "github.com/project-radius/radius/pkg/corerp/api/v20220315privatepreview"
	"github.com/project-radius/radius/pkg/to"
	"github.com/spf13/cobra"

	v1 "github.com/project-radius/radius/pkg/armrpc/api/v1"
	"github.com/project-radius/radius/pkg/cli"
	"github.com/project-radius/radius/pkg/cli/clients"
	"github.com/project-radius/radius/pkg/cli/clierrors"
	"github.com/project-radius/radius/pkg/cli/cmd/commonflags"
	"github.com/project-radius/radius/pkg/cli/cmd/env/namespace"
	"github.com/project-radius/radius/pkg/cli/connections"
	"github.com/project-radius/radius/pkg/cli/framework"
	"github.com/project-radius/radius/pkg/cli/kubernetes"
	"github.com/project-radius/radius/pkg/cli/output"
	"github.com/project-radius/radius/pkg/cli/workspaces"
	"github.com/project-radius/radius/pkg/ucp/resources"
)

// NewCommand creates an instance of the command and runner for the `rad env create` command.
func NewCommand(factory framework.Factory) (*cobra.Command, framework.Runner) {
	runner := NewRunner(factory)

	cmd := &cobra.Command{
		Use:   "create [envName]",
		Short: "Create a new Radius environment",
		Long: `Create a new Radius environment
Radius environments are prepared "landing zones" for Radius applications.
Applications deployed to an environment will inherit the container runtime, configuration, and other settings from the environment.`,
		Args:    cobra.MinimumNArgs(1),
		Example: `rad env create myenv`,
		RunE:    framework.RunCommand(runner),
	}

	commonflags.AddEnvironmentNameFlag(cmd)
	commonflags.AddWorkspaceFlag(cmd)
	commonflags.AddResourceGroupFlag(cmd)
	commonflags.AddNamespaceFlag(cmd)

	return cmd, runner
}

// Runner is the runner implementation for the `rad env create` command.
type Runner struct {
	ConfigHolder        *framework.ConfigHolder
	Output              output.Interface
	Workspace           *workspaces.Workspace
	EnvironmentName     string
	ResourceGroupName   string
	Namespace           string
	ConnectionFactory   connections.Factory
	ConfigFileInterface framework.ConfigFileInterface
	KubernetesInterface kubernetes.Interface
	NamespaceInterface  namespace.Interface
}

// NewRunner creates a new instance of the `rad env create` runner.
func NewRunner(factory framework.Factory) *Runner {
	return &Runner{
		ConfigHolder:        factory.GetConfigHolder(),
		Output:              factory.GetOutput(),
		ConnectionFactory:   factory.GetConnectionFactory(),
		ConfigFileInterface: factory.GetConfigFileInterface(),
		KubernetesInterface: factory.GetKubernetesInterface(),
		NamespaceInterface:  factory.GetNamespaceInterface(),
	}
}

// Validate runs validation for the `rad env create` command.
func (r *Runner) Validate(cmd *cobra.Command, args []string) error {
	workspace, err := cli.RequireWorkspace(cmd, r.ConfigHolder.Config, r.ConfigHolder.DirectoryConfig)
	if err != nil {
		return err
	}
	r.Workspace = workspace

	r.EnvironmentName, err = cli.RequireEnvironmentNameArgs(cmd, args, *workspace)
	if err != nil {
		return err
	}

	r.Workspace.Scope, err = cli.RequireScope(cmd, *r.Workspace)
	if err != nil {
		return err
	}

	r.Namespace, err = cmd.Flags().GetString("namespace")
	if err != nil {
		return err
	} else if r.Namespace == "" {
		r.Namespace = r.EnvironmentName
	}

	client, err := r.ConnectionFactory.CreateApplicationsManagementClient(cmd.Context(), *r.Workspace)
	if err != nil {
		return err
	}

	// Parse the resource group name so we can use it later. DO NOT use the
	// --group argument, because we want to find the right group when the user
	// didn't pass it.
	scopeId, err := resources.Parse(r.Workspace.Scope)
	if err != nil {
		return err
	}
	r.ResourceGroupName = scopeId.FindScope(resources.ResourceGroupsSegment)

	_, err = client.ShowUCPGroup(cmd.Context(), "radius", "local", r.ResourceGroupName)
	if clients.Is404Error(err) {
		return clierrors.Message("Resource group %q could not be found.", r.ResourceGroupName)
	} else if err != nil {
		return err
	}

	err = r.NamespaceInterface.ValidateNamespace(cmd.Context(), r.Namespace)
	if err != nil {
		return err
	}

	return nil
}

// Run runs the `rad env create` command.
func (r *Runner) Run(ctx context.Context) error {
	r.Output.LogInfo("Creating Environment...")

	client, err := r.ConnectionFactory.CreateApplicationsManagementClient(ctx, *r.Workspace)
	if err != nil {
		return err
	}

	envProperties := &corerp.EnvironmentProperties{
		Compute: &corerp.KubernetesCompute{
			Namespace: to.Ptr(r.Namespace),
		},
	}

	err = client.CreateEnvironment(ctx, r.EnvironmentName, v1.LocationGlobal, envProperties)
	if err != nil {
		return err
	}
	r.Output.LogInfo("Successfully created environment %q in resource group %q", r.EnvironmentName, r.ResourceGroupName)

	return nil
}
