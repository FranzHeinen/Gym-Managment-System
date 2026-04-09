document.addEventListener('DOMContentLoaded', function() {
    if (!requireAuth()) return; // De common.js

    loadStats();
    loadHistory(); // <-- Esta es la nueva función
});

/**
 * Carga las estadísticas personales del usuario.
 */
async function loadStats() {
    try {
        const response = await fetch(`/api/profile/stats`, {
            headers: getApiHeaders(false) // De common.js
        });

        if (!response.ok) {
            if (response.status === 401) logout(); // De common.js
            throw new Error('Error al cargar estadísticas');
        }

        const stats = await response.json();

        // 1. Actualizar Frecuencia
        const elFrecuencia = document.getElementById('statFrecuencia');
        if (elFrecuencia) {
            elFrecuencia.textContent = stats.frecuencia_entrenamiento || 0;
        }

        // 2. Actualizar Rutinas Más Usadas
        const ulRutinas = document.getElementById('statRutinas');
        if (ulRutinas) {
            ulRutinas.innerHTML = ''; // Limpiar spinner
            if (stats.rutinas_mas_utilizadas && stats.rutinas_mas_utilizadas.length > 0) {
                stats.rutinas_mas_utilizadas.forEach(item => {
                    const li = document.createElement('li');
                    li.className = 'list-group-item d-flex justify-content-between align-items-center';
                    li.textContent = item.nombre || 'Rutina desconocida';
                    li.innerHTML += `<span class="badge bg-secondary rounded-pill">${item.count}</span>`;
                    ulRutinas.appendChild(li);
                });
            } else {
                ulRutinas.innerHTML = '<li class="list-group-item text-muted text-center">No hay datos</li>';
            }
        }

        // 3. Renderizar Gráfico de Progreso
        const chartLoading = document.getElementById('chartLoading');
        if (chartLoading) chartLoading.style.display = 'none'; // Ocultar spinner
        renderProgresoChart(stats.progreso_en_el_tiempo || []);

    } catch (error) {
        console.error('Error en loadStats:', error);
        const statsContainer = document.getElementById('statsContainer');
        if (statsContainer) {
            statsContainer.innerHTML = '<p class="text-danger text-center">No se pudieron cargar las estadísticas.</p>';
        }
    }
}

/**
 * Renderiza el gráfico de "Progreso en el tiempo" usando Chart.js
 */
function renderProgresoChart(data) {
    const ctx = document.getElementById('progresoChart');
    if (!ctx) return;

    if (data.length === 0) {
        const chartLoading = document.getElementById('chartLoading');
        if (chartLoading) {
            chartLoading.textContent = 'No hay suficientes datos para mostrar un gráfico.';
            chartLoading.style.display = 'block';
        }
        return;
    }

    const labels = data.map(item => `Semana ${item._id.week} (${item._id.year})`);
    const values = data.map(item => item.total_workouts);

    new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Entrenamientos por Semana',
                data: values,
                fill: false,
                borderColor: 'rgb(75, 192, 192)',
                tension: 0.1
            }]
        },
        options: {
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        stepSize: 1 
                    }
                }
            },
            plugins: {
                legend: {
                    display: false
                }
            }
        }
    });
}

 
async function loadHistory() {
    const tBodyHistorial = document.getElementById('statHistorial'); 
    if (!tBodyHistorial) return;

    try {
        const response = await fetch(`/api/workouts`, {
            headers: getApiHeaders(false)
        });
        if (!response.ok) throw new Error('Error al cargar historial');

        const history = await response.json() || [];
        tBodyHistorial.innerHTML = ''; // Limpiar spinner

        if (history.length === 0) {
            tBodyHistorial.innerHTML = '<tr><td class="text-muted text-center">No hay entrenamientos registrados.</td></tr>';
            return;
        }
        
        history.forEach(item => {
            const tr = document.createElement('tr');
            
            const fecha = new Date(item.fecha).toLocaleString('es-ES', {
                day: '2-digit', month: 'short', year: 'numeric',
                hour: '2-digit', minute: '2-digit'
            });
            
            const nombreRutina = item.rutina_nombre || 'Rutina Desconocida';
            
            const rutinaId = item.rutina_id; 
            
            let botonHTML = '';
            if (nombreRutina === 'Rutina Eliminada' || nombreRutina === 'Rutina Desconocida') {
                // Si la rutina no existe, muestra un botón deshabilitado
                tr.classList.add('text-muted'); 
                botonHTML = `
                    <button class="btn btn-sm btn-outline-secondary" disabled>
                        <i class="bi bi-eye-slash"></i> Ver
                    </button>`;
            } else {
                // Si la rutina SÍ existe, muestra el enlace normal
                botonHTML = `
                    <a href="/rutina-detalle?id=${rutinaId}" class="btn btn-sm btn-outline-primary">
                        <i class="bi bi-eye"></i> Ver
                    </a>`;
            }

        tr.innerHTML = `
                <td>
                    <div class="fw-bold">${escapeHtml(nombreRutina)}</div>
                    <small class_name="text-muted">${fecha} hs</small>
                </td>
                <td class="text-end align-middle">
                    ${botonHTML}
                </td>
            `;
            tBodyHistorial.appendChild(tr);
        });

    } catch (error) {
        console.error('Error en loadHistory:', error);
        if (tBodyHistorial) {
            tBodyHistorial.innerHTML = '<tr><td class="text-danger text-center">No se pudo cargar el historial.</td></tr>';
        }
    }
}

function escapeHtml(text) {
    if (typeof text !== 'string') return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}