exports.apps = [
  {
    name: "kalimbu-api",
    script: "/home/connor/kalimbu/cmd/api/main.go",
    interpreter: "go",
    interpreter_args: "run",
  },
  {
    name: "kalimbu-autotask-wh",
    script: "/home/connor/kalimbu/cmd/autotaskWarehouse/main.go",
    interpreter: "go",
    interpreter_args: "run",
  },
  {
    name: "kalimbu-cybercns-wh",
    script: "/home/connor/kalimbu/cmd/cybercnsWarehouse/main.go",
    interpreter: "go",
    interpreter_args: "run",
  },
  {
    name: "kalimbu-rocketcyber-wh",
    script: "/home/connor/kalimbu/cmd/rocketcyberWarehouse/main.go",
    interpreter: "go",
    interpreter_args: "run",
  },
  {
    name: "kalimbu-vsa-wh",
    script: "/home/connor/kalimbu/cmd/vsaWarehouse/main.go",
    interpreter: "go",
    interpreter_args: "run",
  },
  {
    name: "kalimbu-app",
    script: "serve",
    env: {
      PM2_SERVE_PATH: "/home/connor/kalimbu/frontend/dist",
      PM2_SERVE_PORT: 5177,
      PM2_SERVE_SPA: "true",
      PM2_SERVE_HOMEPAGE: "/index.html",
    },
  },
];
