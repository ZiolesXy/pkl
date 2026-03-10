document.addEventListener('DOMContentLoaded', () => {
    // --- 1. Inisialisasi Chart.js ---
    const ctx = document.getElementById('performanceChart').getContext('2d');
    const performanceChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: Array(10).fill(''),
            datasets: [{
                label: 'RAM Usage (%)',
                data: Array(10).fill(0),
                borderColor: '#fbbf24', // Yellow-400
                backgroundColor: 'rgba(251, 191, 36, 0.1)',
                fill: true,
                tension: 0.4
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            animation: { duration: 800 },
            scales: {
                y: { beginAtZero: true, max: 100, grid: { color: '#374151' }, ticks: { color: '#9ca3af' } },
                x: { grid: { display: false } }
            },
            plugins: {
                legend: { labels: { color: '#9ca3af', usePointStyle: true } }
            }
        }
    });

    // --- 2. Fungsi Utama Fetch Data ---
    async function fetchDashboardData() {
        // const API_URL = 'https://undeliberatingly-decemviral-petronila.ngrok-free.dev/api/v1/user/device-info';
        const API_URL = 'http://10.192.132.101:8000/api/v1/user/device-info';
        const statusEl = document.getElementById('status-message');

        try {
            const response = await fetch(API_URL, {
                method: 'GET',
                headers: {
                    'ngrok-skip-browser-warning': 'true',
                    'Accept': 'application/json'
                }
            });

            // Proteksi jika ngrok mengirim HTML (halaman peringatan)
            const contentType = response.headers.get("content-type");
            if (!contentType || !contentType.includes("application/json")) {
                throw new TypeError("Server tidak mengirim JSON. Periksa koneksi ngrok.");
            }

            const result = await response.json();

            if (result.success) {
                const { server, client } = result.data;

                // --- Update UI Server ---
                document.getElementById('server-host').textContent = server.server_host_name || '-';
                document.getElementById('server-os').textContent = server.server_os || '-';
                document.getElementById('server-cpu').textContent = server.server_cpu_model || '-';
                document.getElementById('server-suggestion').textContent = server.suggestion || '-';
                
                // Menangani Internet Status (Green/Red badge)
                const netStatus = document.getElementById('net-status');
                netStatus.textContent = server.internet_status || '-';
                netStatus.className = server.internet_status?.toLowerCase() === 'online' 
                    ? "px-2 py-0.5 bg-green-900/50 text-green-400 text-[10px] font-bold rounded uppercase"
                    : "px-2 py-0.5 bg-red-900/50 text-red-400 text-[10px] font-bold rounded uppercase";

                // --- Logic RAM & Chart ---
                // Regex untuk mengambil angka dari string "Total (terpakai: 87.0%)"
                const ramMatch = server.server_ram_usage.match(/(\d+(\.\d+)?)%/);
                let ramPercent = 0;

                if (ramMatch) {
                    ramPercent = parseFloat(ramMatch[1]);
                    
                    // Update Text & Progress Bar
                    document.getElementById('server-ram-val').textContent = `${ramPercent}%`;
                    const ramBar = document.getElementById('ram-progress-bar');
                    ramBar.style.width = `${ramPercent}%`;

                    // Update Warna Bar (Kritis > 90)
                    if (ramPercent > 90) {
                        ramBar.classList.replace('bg-yellow-500', 'bg-red-500');
                    } else {
                        ramBar.classList.replace('bg-red-500', 'bg-yellow-500');
                    }

                    // Update Chart Data
                    performanceChart.data.datasets[0].data.shift();
                    performanceChart.data.datasets[0].data.push(ramPercent);
                    performanceChart.update();
                }

                // --- Update UI Client ---
                document.getElementById('client-ip').textContent = client.user_ip || '-';
                document.getElementById('client-browser').textContent = `${client.browser} v${client.browser_version}`;
                document.getElementById('client-device').textContent = client.device_type || '-';
                document.getElementById('client-ua').textContent = client.user_agent || '-';

                // Update Status Bar
                statusEl.innerHTML = `<span class="flex items-center gap-2"><span class="w-2 h-2 bg-green-500 rounded-full"></span> Data Tersinkronisasi</span>`;
                statusEl.className = "text-green-400 text-sm";
            }
        } catch (error) {
            console.error("Detail Error:", error);
            statusEl.textContent = "Koneksi Gagal: " + error.message;
            statusEl.className = "text-red-500 text-sm font-bold";
        }
    }

    // Jalankan pertama kali
    fetchDashboardData();

    // Auto-refresh setiap 5 detik agar diagram bergerak
    setInterval(fetchDashboardData, 5000);
});