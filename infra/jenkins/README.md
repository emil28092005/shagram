# Jenkins Controller + Docker Agent (Docker Compose)

This directory contains a minimal Jenkins setup running as Docker containers:
- `jenkins-controller`: Jenkins UI + configuration (no builds should run here).
- `jenkins-agent-docker`: inbound agent used for Docker-based pipelines (label: `docker-agent`).

The goal is to keep the controller responsible for orchestration/configuration while all CI jobs run on the dedicated agent.

## Prerequisites
- Docker Engine installed on the host.
- Docker Compose available as `docker compose`.
- A host directory `/opt/shagram` (used by pipelines as a shared location for source code and deployment files).

## Start Jenkins
From the repository root:

```bash
cd infra/jenkins
docker compose up -d
docker compose ps
```

Jenkins UI will be available at:

```text
http://<server-ip>:8080
```

## Initial admin password
Retrieve the initial password:

```bash
docker exec -it jenkins-controller \
  cat /var/jenkins_home/secrets/initialAdminPassword
```

## Disable builds on the controller
To ensure builds do not run on the controller:
1. Open Jenkins UI.
2. Go to: `Manage Jenkins` → `Manage Nodes and Clouds`.
3. Open: `Built-In Node` → `Configure`.
4. Set **Number of executors** to `0`.
5. Save.

## Configure the inbound agent node
Create a dedicated node for running pipelines:
1. `Manage Jenkins` → `Manage Nodes and Clouds` → `New Node`.
2. Set:
   - Node name: `docker-agent`
   - Type: `Permanent Agent`
   - Remote root directory: `/home/jenkins/agent`
   - Labels: `docker-agent` (must match your Jenkinsfile `agent { label ... }`)
   - Usage: “Only build jobs with label expressions matching this node”
3. Save.

After saving, open the agent page:
- `Manage Nodes and Clouds` → `docker-agent`

On that page Jenkins will show the inbound connection details, including the **secret** required by the inbound agent container.

## Set agent secret in Compose
Edit `infra/jenkins/compose.yaml` and replace the placeholder with the real secret value shown on the `docker-agent` node page:

- `JENKINS_SECRET=PASTE_ME`

Then restart only the agent container:

```bash
docker compose up -d --force-recreate jenkins-agent-docker
docker logs -f jenkins-agent-docker
```

Verification:
- In Jenkins UI, the node `docker-agent` should become **Online**.
- Any pipeline using `agent { label 'docker-agent' }` should execute on this agent.

## Security note
The Docker agent mounts `/var/run/docker.sock`, which effectively grants high-level control over the Docker host.

Use this setup only in trusted environments and limit access to Jenkins accordingly.
