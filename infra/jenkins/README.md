# Jenkins (Controller + Docker Agent) Setup (Docker Compose)

This directory contains a minimal Jenkins setup using a dedicated **controller** and a separate inbound **agent** intended for Docker-based workloads.

The goal is to keep the controller responsible for orchestration and configuration, while all builds run on the agent labeled `docker`.

## Prerequisites
- Docker Engine installed on the host
- Docker Compose available as `docker compose`
- A host directory `/opt/shagram` (used by pipelines as a shared location for source code and deployment files)

## Start Jenkins
From the repository root:

```bash
cd infra/jenkins
docker compose up -d
docker compose ps
```

Jenkins UI will be available at:

- http://<server-ip>:8080

### Initial admin password
Retrieve the initial password with:

```bash
docker exec -it jenkins-controller cat /var/jenkins_home/secrets/initialAdminPassword
```

## Disable builds on the built-in node (Executors = 0)
To ensure builds do not run on the controller:

1. Open Jenkins UI
2. Go to: **Manage Jenkins → Manage Nodes and Clouds**
3. Open **Built-In Node → Configure**
4. Set **Number of executors** to `0`
5. Save

## Create an inbound agent node (docker-agent)
Create a dedicated node for running pipelines:

1. Go to: **Manage Jenkins → Manage Nodes and Clouds → New Node**
2. Set:
   - **Node name**: `docker-agent`
   - **Type**: Permanent Agent
   - **Remote root directory**: `/home/jenkins/agent`
   - **Labels**: `docker`
   - **Usage**: Only build jobs with label expressions matching this node
3. Save

After saving, open the agent page:

- **Manage Nodes and Clouds → docker-agent**

On that page, Jenkins provides the inbound connection details, including the **secret** required by the inbound agent container.

## Configure the agent secret in Compose
Edit `infra/jenkins/compose.yaml` and replace:

- `JENKINS_SECRET=__PASTE_ME__`

with the real secret value shown on the `docker-agent` node page.

Restart only the agent container:

```bash
docker compose up -d --force-recreate jenkins-agent-docker
docker logs -f jenkins-agent-docker
```

## Verification
- In Jenkins UI: **Manage Nodes and Clouds**, the node `docker-agent` should be **Online**
- Any pipeline using `agent { label 'docker' }` should execute on this agent

## Security note
The Docker agent container mounts `/var/run/docker.sock`, which effectively grants high-level control over the Docker host. Use this setup only in trusted environments and limit access to Jenkins accordingly.
