// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/bpmn/1.0/policies.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_1

policy_types:

  bpmn.Root:
    description: >-
      Root for policies implemented by BPM software.
    derived_from: Root

  bpmn.Process:
    description: >-
      Policy implemented by a process defined in BPMN.
    derived_from: bpmn.Root
    properties:
      bpmn_process_id:
        description: >-
          Execute this BPMN process when triggered.
        type: string
    targets:
    - Root
`
}
