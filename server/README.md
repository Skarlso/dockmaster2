Directory for the Server
========================

Dockmaster server. This server will store data received from the Agent in a mongodb collection.

The form of this data is as follows:

```json
{
	"agentid": "007",
	"expireAfterSeconds": 300,
	"containers": [{
		"id": "cont1",
		"name": "bla",
		"command": "catalina.sh",
		"port": "3321"
	}, {
		"id": "cont2",
		"name": "bla2",
		"command": "catalina2.sh",
		"port": "3333"
	}]
}
```

This will create a record for agent 007 and with a list of containers defined in ```containers```. Whenever the agent resends data, at the current version, all previous data for this agent is completely overwritten. This is so that avoid duplications and not have to deal with individual updates to containers.

In a next version perhaps the change will only be applied if there is actual change to one for the containers or information regarding the agent.

For now, the server provides these end-points:

```bash
[GIN-debug] GET    /api/1/list               --> main.listContainers (3 handlers)
[GIN-debug] POST   /api/1/add                --> main.addContainers (3 handlers)
[GIN-debug] POST   /api/1/delete             --> main.deleteContainers (3 handlers)
[GIN-debug] GET    /api/1/inspect/:agentID/:containerID --> main.inspectContainer (3 handlers)
[GIN-debug] POST   /api/1/stopAll            --> main.stopAll (3 handlers)
[GIN-debug] POST   /api/1/stop/:agentID      --> main.stop (3 handlers)
```

The ```delete``` however, is not used at the moment.
