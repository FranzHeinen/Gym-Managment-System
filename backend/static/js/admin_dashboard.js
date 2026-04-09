document.addEventListener('DOMContentLoaded', function() {
    // 1. Verifica autenticación y rol de Admin
    if (!requireAuth()) return;
    if (localStorage.getItem('userRole') !== 'ADMIN') {
        window.location.href = '/rutinas'; 
        return;
    }

    // 2. Carga las estadísticas
    cargarEstadisticas();
});

/**
 * Llama a la API de estadísticas globales y actualiza la UI.
 */
async function cargarEstadisticas() {
    try {
        const response = await fetch('/api/admin/stats', {
            headers: getApiHeaders(false)
        });

        if (!response.ok) {
            if (response.status === 401) logout();
            throw new Error('Error al cargar las estadísticas');
        }

        const stats = await response.json();

        // Actualizar tarjetas
        actualizarTotalUsuarios(stats.total_usuarios);
        actualizarLista(stats.ejercicios_populares, 'ejerciciosPopulares', 'No hay ejercicios populares');

    } catch (error) {
        console.error('Error en cargarEstadisticas:', error);
        document.getElementById('totalUsuarios').textContent = 'Error';
    }
}

/**
 * Actualiza el contador de usuarios.
 */
function actualizarTotalUsuarios(total) {
    const el = document.getElementById('totalUsuarios');
    if (el) {
        el.textContent = total !== undefined ? total : '0';
    }
}

/**
 * Rellena una lista <ul> con los datos de popularidad.
 */
function actualizarLista(items, elementId, mensajeVacio) {
    const ul = document.getElementById(elementId);
    if (!ul) return;

    if (!items || items.length === 0) {
        ul.innerHTML = `<li class="list-group-item text-muted">${mensajeVacio}</li>`;
        return;
    }

    ul.innerHTML = ''; // Limpiar spinner
    items.forEach(item => {
        const li = document.createElement('li');
        li.className = 'list-group-item d-flex justify-content-between align-items-center';
        
        // Gracias a la modificación del backend, item.nombre ya existe
        li.textContent = item.nombre || 'Nombre no disponible';
        
        const badge = document.createElement('span');
        badge.className = 'badge bg-primary rounded-pill';
        badge.textContent = item.count;
        
        li.appendChild(badge);
        ul.appendChild(li);
    });
}