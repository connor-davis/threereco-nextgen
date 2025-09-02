exports.apps = [
  {
    name: "three-api",
    script: "/home/3reco/cmd/api/main.go",
    interpreter: "go",
    interpreter_args: "run",
  },
  {
    name: "three-app",
    script: "serve",
    env: {
      PM2_SERVE_PATH: "/home/3reco/frontend/dist",
      PM2_SERVE_PORT: 5173,
      PM2_SERVE_SPA: "true",
      PM2_SERVE_HOMEPAGE: "/index.html",
    },
  },
];
