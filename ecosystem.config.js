exports.apps = [
  {
    name: "three-api",
    script: "/home/connor/threereco-nextgen/cmd/api/main.go",
    interpreter: "go",
    interpreter_args: "run",
  },
  {
    name: "three-app",
    script: "serve",
    env: {
      PM2_SERVE_PATH: "/home/connor/threereco-nextgen/frontend/dist",
      PM2_SERVE_PORT: 5177,
      PM2_SERVE_SPA: "true",
      PM2_SERVE_HOMEPAGE: "/index.html",
    },
  },
];
