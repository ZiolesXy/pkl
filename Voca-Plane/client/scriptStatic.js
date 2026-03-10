const API_URL = "http://172.16.17.123:8000/api/v1/user/device-info";
// const API_URL = "https://undeliberatingly-decemviral-petronila.ngrok-free.dev/api/v1/user/device-info";
const ctx = document.getElementById("chart").getContext("2d");

// Kita gunakan interval 3 detik agar server tidak lelah, tapi visual tetap smooth
const UPDATE_INTERVAL = 3000; 

const chart = new Chart(ctx, {
    type: "line",
    data: {
        labels: Array(10).fill(""),
        datasets: [
            {
                label: "CPU %",
                data: Array(10).fill(0),
                borderColor: "#22c55e",
                backgroundColor: "rgba(34, 197, 94, 0.1)",
                tension: 0.3, // Sedikit lengkungan agar tidak terlalu kaku
                fill: true,
                pointRadius: 2
            },
            {
                label: "RAM %",
                data: Array(10).fill(0),
                borderColor: "#facc15",
                backgroundColor: "rgba(250, 204, 21, 0.1)",
                tension: 0.3,
                fill: true,
                pointRadius: 2
            },
            {
                label: "Disk %",
                data: Array(10).fill(0),
                borderColor: "#a855f7",
                backgroundColor: "rgba(168, 85, 247, 0.1)",
                tension: 0.3,
                fill: true,
                pointRadius: 2
            }
        ]
    },
    options: {
        responsive: true,
        maintainAspectRatio: false,
        // Kuncinya di sini: Matikan animasi global agar tidak "melompat"
        animation: false, 
        // Tapi berikan transisi halus pada elemen garis saat data berubah
        transitions: {
            active: {
                animation: {
                    duration: 1000 // Transisi pergeseran titik selama 1 detik saja
                }
            }
        },
        plugins: {
            legend: { labels: { color: "#9ca3af" } }
        },
        scales: {
            y: { 
                min: 0, 
                max: 100,
                grid: { color: "rgba(255, 255, 255, 0.05)" },
                ticks: { color: "#9ca3af" }
            },
            x: {
                grid: { display: false },
                ticks: { display: false }
            }
        }
    }
});

async function loadData() {
    try {
        const res = await fetch(API_URL, {
            method: "GET",
            headers: {
                "ngrok-skip-browser-warning": "true",
                "Accept": "application/json"
            }
        });
        const json = await res.json();
        const server = json.data.server;
        const client = json.data.client;

        // --- UPDATE TEXT CARDS (ATAS) ---
        if(document.getElementById("status")) document.getElementById("status").textContent = "Live Monitoring";
        if(document.getElementById("cpuUsage")) document.getElementById("cpuUsage").textContent = server.server_cpu_usage;
        if(document.getElementById("ramUsage")) document.getElementById("ramUsage").textContent = server.server_ram_usage;
        if(document.getElementById("diskUsage")) document.getElementById("diskUsage").textContent = server.server_disk_usage;
        if(document.getElementById("healthScore")) document.getElementById("healthScore").textContent = server.health_score + "/100";

        // --- UPDATE SERVER INFORMATION (TENGAH) ---
        if(document.getElementById("serverName")) document.getElementById("serverName").textContent = server.server_host_name;
        if(document.getElementById("serverOS")) document.getElementById("serverOS").textContent = server.server_os;
        if(document.getElementById("cpuModel")) document.getElementById("cpuModel").textContent = server.server_cpu_model;
        if(document.getElementById("uptime")) document.getElementById("uptime").textContent = server.system_uptime;
        if(document.getElementById("appUptime")) document.getElementById("appUptime").textContent = server.app_uptime;
        if(document.getElementById("publicIP")) document.getElementById("publicIP").textContent = server.server_public_ip;
        if(document.getElementById("internetStatus")) document.getElementById("internetStatus").textContent = server.internet_status;
        
        // Perhatikan ID 'env' -> mengambil dari server.environment
        if(document.getElementById("environment")) document.getElementById("environment").textContent = server.environment;
        if(document.getElementById("suggestion")) document.getElementById("suggestion").textContent = server.suggestion;

        // --- UPDATE CLIENT INFO (BAWAH) ---
        if(document.getElementById("clientIP")) document.getElementById("clientIP").textContent = client.user_ip;
        if(document.getElementById("device")) document.getElementById("device").textContent = (client.device_brand || "") + " " + (client.device_model || "");
        if(document.getElementById("browser")) document.getElementById("browser").textContent = client.browser;
        if(document.getElementById("os")) document.getElementById("os").textContent = client.os;
        if(document.getElementById("fingerprint")) document.getElementById("fingerprint").textContent = client.fingerprint;
        if(document.getElementById("country")) document.getElementById("country").textContent = client.country;
        if(document.getElementById("city")) document.getElementById("city").textContent = client.city;
        if(document.getElementById("isp")) document.getElementById("isp").textContent = client.isp;
        
        // Perhatikan ID 'isBot' -> logika boolean
        if(document.getElementById("bot")) {
            document.getElementById("bot").textContent = client.is_bot ? "Yes" : "No";
        }

        // --- DATA UPDATE UNTUK CHART ---
        // Menghapus simbol '%' agar bisa diproses sebagai angka
        const cpu = parseFloat(server.server_cpu_usage.replace('%', '')) || 0;
        
        // Regex untuk mengambil angka di dalam tanda kurung (e.g., "Used 88.0%")
        const ramMatch = server.server_ram_usage.match(/(\d+(\.\d+)?)%/);
        const ram = ramMatch ? parseFloat(ramMatch[1]) : 0;
        
        // Regex untuk mengambil angka di dalam tanda kurung (e.g., "62.5% used")
        const diskMatch = server.server_disk_usage.match(/(\d+(\.\d+)?)%/);
        const disk = diskMatch ? parseFloat(diskMatch[1]) : 0;

        chart.data.datasets[0].data.shift();
        chart.data.datasets[0].data.push(cpu);
        chart.data.datasets[1].data.shift();
        chart.data.datasets[1].data.push(ram);
        chart.data.datasets[2].data.shift();
        chart.data.datasets[2].data.push(disk);

        chart.update('active');

    } catch (e) {
        if(document.getElementById("status")) document.getElementById("status").textContent = "Disconnected";
        console.error("Fetch Error:", e);
    }
}

loadData();
setInterval(loadData, UPDATE_INTERVAL);