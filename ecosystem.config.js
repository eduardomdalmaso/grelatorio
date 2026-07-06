const path = require('path');

const wailsPath = path.join(process.env.USERPROFILE || 'C:\\Users\\eduar', 'go', 'bin', 'wails.exe');

module.exports = {
  apps: [
    {
      name: "relatorio-wails",
      script: wailsPath,
      args: "dev",
      cwd: __dirname,
      watch: false,
      autorestart: false,
      max_memory_restart: "1G",
      env: {
        PATH: process.env.PATH + ";C:\\Program Files\\Go\\bin;" + path.dirname(wailsPath)
      }
    },
    {
      name: "relatorio-frontend",
      script: "npm",
      args: "run dev",
      cwd: path.join(__dirname, "frontend"),
      watch: false,
      autorestart: true,
      max_memory_restart: "1G"
    }
  ]
};


