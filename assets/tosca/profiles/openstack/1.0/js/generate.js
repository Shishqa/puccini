
clout.exec('tosca.utils');

tosca.coerce();

puccini.writeString('\
[defaults]\n\
ansible_managed=Modified by Ansible on %Y-%m-%d %H:%M:%S %Z\n\
inventory=./inventory.yaml\n\
transport=ssh\n\
command_warnings=false\n\
\n\
# Better concurrency\n\
forks=25\n\
\n\
# Required because in the cloud IP addresses can be reused\n\
host_key_checking=false\n\
\n\
[ssh_connection]\n\
# Faster SSH\n\
pipelining=true\n\
', 'ansible.cfg');

puccini.write({
	all: {
		hosts: {
			localhost: {
				ansible_connection: 'local'
			}
		}
	}
}, 'inventory.' + puccini.format);

puccini.write({
	clouds: {}
}, 'clouds.' + puccini.format);

puccini.write([[{
	name: 'Provision keypair',
	register: 'keypair',
	os_keypair: {
		state: 'present',
		name: 'topology'
	}
}, {
	name: 'Write keys',
	when: 'keypair.key.private_key is not none', // will only be available when keypair is created
	block: [{
		file: {
			path: '{{ playbook_dir }}/keys',
			state: 'directory'
		}
	}, {
		copy: {
			content: '{{ keypair.key.private_key }}',
			dest: '{{ playbook_dir }}/keys/{{ keypair.key.name }}'
		}
	}, {
		file: {
			mode: 0600, // required by ssh
			dest: '{{ playbook_dir }}/keys/{{ keypair.key.name }}'
		}
	}, {
		copy: {
			content: '{{ keypair.key.public_key }}',
			dest: '{{ playbook_dir }}/keys/{{ keypair.key.name }}.pub'
		}
	}]
}]], 'roles/openstack-keypair/tasks.main.' + puccini.format);

puccini.write([[{
	name: 'Provision servers',
	async: 300, // 5 minutes
	register: 'servers_async',
	with_items: '{{ topology.servers }}',
	os_server: {
	    state: 'present',
	    name: '{{ item.type }}-{{ item.index }}.{{ topology.site_name }}.{{ topology.zone }}',
	    image: '{{ topology.image }}',
	    flavor: '{{ item.flavor }}',
	    key_name: '{{ keypair.key.name }}'
	}
}, {
	name: 'Wait for servers to become active',
	retries: 300, // delay is 5 seconds
	register: 'servers',
	until: 'servers.finished',
	with_items: '{{ servers_async.results }}',
	async_status: {
		jid: '{{ item.ansible_job_id }}'
	}
}, {
	name: 'Add servers to group',
	with_items: '{{ servers.results }}',
	add_host: {
		name: '{{ item.server.name }}',
		groups: 'servers',
		// Custom attributes:
		server: '{{ item.server }}',
		// Ansible attributes:
	    ansible_ssh_host: '{{ item.server.public_v4 }}',
	    ansible_ssh_user: 'root',
	    ansible_ssh_private_key_file: '{{ playbook_dir }}/keys/{{ keypair.key.name }}',
	}
}]], 'roles/openstack-servers/tasks.main.' + puccini.format);


playbook = [];

provision = {
	hosts: 'localhost',
	gather_facts: false,
	tasks: [{
		name: 'Provision servers',
		async: 300, // 5 minutes
		register: 'servers_async',
		with_items: '{{ topology.servers }}',
		os_server: {
		    state: 'present',
		    name: '{{ item.type }}-{{ item.index }}.{{ topology.site_name }}.{{ topology.zone }}',
		    image: '{{ topology.image }}',
		    flavor: '{{ item.flavor }}',
		    key_name: '{{ keypair.key.name }}'
		}
	}]
};

playbook.push(provision);

for (v in clout.vertexes) {
	vertex = clout.vertexes[v];
	if (!tosca.isNodeTemplate(vertex, 'openstack.Nova.Server'))
		continue;
	nodeTemplate = vertex.properties;

	// provision.tasks.push();
}

puccini.write([[{
	hosts: 'localhost',
	gather_facts: false,
	tasks: [{
		name: 'Configure OpenStack',
		os_client_config: null,
	}, {
		include_role: {
			name: 'openstack-keypair'
		}
	}, {
		include_role: {
			name: 'openstack-servers'
		}
	}]
}]], 'install.' + puccini.format);
