## Objective: Multi-Node Key-Value Store with Replication


## Replication
Replication is the process of keeping multiple nodes alive with the same data as the current main node.

### Two Types of Replication:

- **Synchronous**:
  Our elected leader node does not confirm or accept writes until `(N // 2) + 1` nodes have verified they received the write and sent an OK back to the elected leader.

- **Asynchronous**:
  Our elected leader confirms writes upon receiving but will asynchronously send out the updates to followers. In this, our follower nodes must agree on a quorum read, meaning out of `K` nodes that are read, the majority answer is returned.

## Failover
Failover is the process of our current leader failing and being replaced with a replica.

To elect a leader, we can use any well-defined leader election algorithm. Any writes should be redirected to a durable queue during the failover process, then be asynchronously pushed to the new leader.

### Types of Failover:

- **Active-Active**:
  Leaderless system where writes and reads are processed by any other replica node, keeping availability.

- **Active-Passive**:
  Leader-follower system where, upon failure of a leader, writes are blocked until a replica replaces it.

This project will follow **Active-Passive failover** as we are using a leader-follower architecture.

## Health Check

Health checks are the process of pinging our nodes to check if they are responding, if a node does not respond by N intervals we mark it as dead. If the leader is dead start the failover process otherwise
