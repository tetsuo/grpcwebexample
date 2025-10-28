import { GreeterClient } from "@gen/helloworld/helloworld_grpc_web_pb.js";
import { HelloRequest } from "@gen/helloworld/helloworld_pb.js";

const endpoint = import.meta.env.VITE_GRPC_ENDPOINT;
const client = new GreeterClient(endpoint, null, {
  withCredentials: true,
});

const nameInput = document.getElementById("nameInput");
const submitBtn = document.getElementById("callBtn");
const output = document.getElementById("responseArea");

function setOutput(text) {
  output.textContent = text;
}

function callHello(name) {
  const request = new HelloRequest();
  request.setName(name);

  return new Promise((resolve, reject) => {
    client.sayHello(request, {}, (err, response) => {
      if (err) {
        reject(err);
        return;
      }
      resolve(response.getMessage());
    });
  });
}

submitBtn?.addEventListener("click", async () => {
  const name = (nameInput?.value || "").trim();
  setOutput("Calling...");
  try {
    const message = await callHello(name);
    setOutput(JSON.stringify({ message }, null, 2));
  } catch (err) {
    const msg = err?.message || err?.details || String(err);
    setOutput(`Error: ${msg}`);
  }
});
