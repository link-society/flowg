"use strict";

const fs = require("node:fs");


function getInput(name, required = false) {
  const value = process.env[`INPUT_${name.toUpperCase()}`]?.trim() ?? "";

  if (required && value === "") {
    throw new Error(`Missing required input: ${name}`);
  }

  return value;
}

function setOutput(name, value) {
  const outputFile = process.env.GITHUB_OUTPUT;

  if (!outputFile) {
    throw new Error("Missing GITHUB_OUTPUT environment variable");
  }

  fs.appendFileSync(outputFile, `${name}=${value}\n`);
}


async function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function truncate(str, maxLength = 1000) {
  if (str.length <= maxLength) {
    return str;
  }

  return `${str.slice(0, maxLength)}... [truncated]`;
}


function renderBuildUrl(endpoint) {
  return `https://api.render.com/v1/${endpoint}`;
}

function renderBuildServiceUrl(serviceId) {
  return renderBuildUrl(`services/${encodeURIComponent(serviceId)}`);
}

function renderBuildServiceDeploysUrl(serviceId) {
  return `${renderBuildServiceUrl(serviceId)}/deploys`;
}

function renderBuildDeployUrl(serviceId, deployId) {
  return `${renderBuildServiceDeploysUrl(serviceId)}/${encodeURIComponent(deployId)}`;
}

async function renderApiRequest({ url, method = "GET", body = null, apiKey }) {
  const attempts = 4;

  const TRANSIENT_HTTP_STATUS_CODES = new Set([
    429, // Too Many Requests
    500, // Internal Server Error
    502, // Bad Gateway
    503, // Service Unavailable
    504, // Gateway Timeout
  ]);

  for (let attempt = 1; attempt <= attempts; attempt++) {
    let response;

    try {
      response = await fetch(url, {
        method,
        signal: AbortSignal.timeout(60_000),
        headers: {
          "Authorization": `Bearer ${apiKey}`,
          "Accept": "application/json",
          "Content-Type": "application/json",
          "User-Agent": "flowg-github-action-render-deploy",
        },
        body: body ? JSON.stringify(body) : null,
      });
    }
    catch (error) {
      if (attempt === attempts) {
        throw new Error(
          `${method} ${url} failed: ${error.message}`,
          { cause: error },
        );
      }

      await sleep(attempt * 2000);
      continue;
    }

    const responseText = await response.text();

    if (response.ok) {
      if (responseText === "") {
        return null;
      }

      try {
        return JSON.parse(responseText);
      }
      catch (error) {
        throw new Error(
          `${method} ${url} returned invalid JSON: ${truncate(responseText)}`,
          { cause: error },
        );
      }
    }

    if (
      TRANSIENT_HTTP_STATUS_CODES.has(response.status) &&
      attempt < attempts
    ) {
      const retryAfter = Number.parseInt(
        response.headers.get("Retry-After") ?? "",
        10,
      );

      const delayMs = Number.isInteger(retryAfter)
        ? retryAfter * 1000
        : attempt * 3000;

      console.log(`${method} ${url} returned ${response.status}. Retrying in ${delayMs}ms...`);
      await sleep(delayMs);
      continue;
    }

    throw new Error(`${method} ${url} returned ${response.status}: ${truncate(responseText)}`);
  }

  throw new Error(`${method} ${url} failed after ${attempts} attempts`);
}


async function updateService({ apiKey, serviceId, imageName }) {
  console.log(`⬆️ Update Render Service image to: ${imageName}`);

  await renderApiRequest({
    url: renderBuildServiceUrl(serviceId),
    method: "PATCH",
    body: {
      image: {
        name: imageName,
      },
    },
    apiKey,
  });
}

async function triggerDeploy({ apiKey, serviceId, imageName }) {
  console.log(`🚀 Trigger Render Service deploy for image: ${imageName}`);

  const deployment = await renderApiRequest({
    url: renderBuildServiceDeploysUrl(serviceId),
    method: "POST",
    body: {
      imageUrl: imageName,
    },
    apiKey,
  });

  if (!deployment?.id) {
    throw new Error(`Render did not return a deployment ID`);
  }

  console.log(`Render Deployment ID: ${deployment.id}`);
  return deployment.id;
}

async function waitForDeploy({ apiKey, serviceId, deployId }) {
  console.log(`⏳ Waiting for Render Service deploy to complete...`);

  const TIMEOUT = 1800 /* seconds */;
  const POLL_INTERVAL = 10 /* seconds */;

  const FAILURE_STATUSES = new Set([
    "build_failed",
    "update_failed",
    "pre_deploy_failed",
    "canceled",
    "deactivated",
  ]);

  const deadline = Date.now() + TIMEOUT * 1000;
  let previousStatus = null;

  while (Date.now() < deadline) {
    const deployment = await renderApiRequest({
      url: renderBuildDeployUrl(serviceId, deployId),
      apiKey,
    });

    const status = deployment?.status;

    if (!status) {
      throw new Error(`Render did not return a deployment status`);
    }

    if (status !== previousStatus) {
      console.log(`Render Deployment Status: ${status}`);
      previousStatus = status;
    }

    if (status === "live") {
      console.log(`✅ Render Service deploy completed successfully`);
      return;
    }

    if (FAILURE_STATUSES.has(status)) {
      throw new Error(`Render Service deploy failed with status: ${status}`);
    }

    await sleep(POLL_INTERVAL * 1000);
  }

  throw new Error(`Render Service deploy timed out after ${TIMEOUT} seconds`);
}


async function main() {
  const apiKey = getInput("render_api_key", true);
  const serviceId = getInput("render_service_id", true);
  const imageName = getInput("docker_image_name", true);

  await updateService({ apiKey, serviceId, imageName });
  const deployId = await triggerDeploy({ apiKey, serviceId, imageName });
  setOutput("render_deploy_id", deployId);
  await waitForDeploy({ apiKey, serviceId, deployId });
}

main().catch((error) => {
  console.error("❌ Render deploy failed:", error.message);
  console.error(error.stack);
  process.exitCode = 1;
});
