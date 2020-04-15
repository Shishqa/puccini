// This file was auto-generated from a YAML file

package v1_3

func init() {
	Profile["/tosca/simple/1.3/js/functions/get_artifact.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.8.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.8.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.8.1
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.8.1

clout.exec('tosca.lib.utils');

function evaluate(entity, artifactName, location, remove) {
	if (arguments.length < 2)
		throw 'must have at least 2 arguments';
	var nodeTemplate = tosca.getModelableEntity(entity);
	if (!nodeTemplate.artifacts || !(artifactName in nodeTemplate.artifacts))
		throw puccini.sprintf('artifact "%s" not found in "%s"', artifactName, nodeTemplate.name);
	var artifact = nodeTemplate.artifacts[artifactName];
	return artifact.$artifact;
}
`
}
