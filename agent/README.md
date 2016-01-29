Agent
=====

The agent, who talks with the docker socket on a given machine. The purpose of the agent is to be started on a given server where there are running docker containers. The Agent will talk to the Docker daemon through his socket. For now, the agent will only talk with a docker daemon who is located on the same server where the agent is.

The agent's startup flags are the following:

* server default: "http://localhost:8989" explanation: "The server uri where dockmaster is located."
* agent default: "localhost", explanation: "The name of an Agent. Example: TestQA1."
* refresh default: 60, explanation: "The rate at which this agent should check for changes in seconds."
* expireAfterSeconds default: 60, explanation: "The rate at which data sent by this agent should expire in seconds."

The agent will go into an endless loop and continuously talk to the docker daemon and request running containers. Once received, it will post them to the dockmaster server and moves on to sleep the amount of time defined in refresh.

The agent will only quit if it cannot talk to the docker socket. Any other error it will just log and then restart the cycle and try again.
