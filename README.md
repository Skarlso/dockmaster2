Dockmaster
==========

Dockmaster who oversees the workings of containers throughout a network.

Dockmaster is a overseer for all the containers running on all the servers. The architecture of Dockmaster is simple.

Server
------

Will sit on a central server which is reachable by all the nodes on which an Agent is operating.

Agent
-----

An Agent sits on a node where there are docker containers running. The Agent reports to the Server about statuses of the containers next-to it.


Currently the Agent only reports about running containers. Other containers are not regarded because it would create a very convoluted view.

The frontend first needs to handle proper showing of the data that it gets.

Current Look
------------

The current look is a very simple table view:

![Dockmaster](dockmaster.png)
