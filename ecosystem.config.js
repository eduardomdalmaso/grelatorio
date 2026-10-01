const path = require('path');

const userProfile = process.env.USERPROFILE || 'C:\\Users\\hades';
const wailsPath = path.join(userProfile, 'go', 'bin', 'wails.exe');

module.exports = {
  apps: [
    {
      name: "relatorio-wails",
      script: wailsPath,
      args: "dev",
      cwd: __dirname,
      interpreter: "none",
      watch: false,
      autorestart: false,
      max_memory_restart: "1G",
      env: {
        PATH: process.env.PATH + ";C:\\Program Files\\Go\\bin;" + path.dirname(wailsPath)
      }
    },
    {
      name: "relatorio-frontend",
      script: "node_modules/vite/bin/vite.js",
      cwd: path.join(__dirname, "frontend"),
      watch: false,
      autorestart: true,
      max_memory_restart: "1G"
    }
  ]
};


