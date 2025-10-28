import { defineConfig } from "vite";
import path from "node:path";

export default defineConfig({
  root: ".",
  publicDir: "public",
  base: "/static/",
  resolve: {
    alias: {
      "@gen": path.resolve(__dirname, "gen"),
    },
  },
  optimizeDeps: {
    include: [
      "google-protobuf",
      "grpc-web",
      "@gen/helloworld/helloworld_pb.js",
      "@gen/helloworld/helloworld_grpc_web_pb.js",
    ],
  },
  build: {
    outDir: "dist",
    sourcemap: true,
    commonjsOptions: {
      include: [/node_modules/, /web\/gen/, /gen\/helloworld/],
    },
  },
});
