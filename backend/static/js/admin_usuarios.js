document.addEventListener('DOMContentLoaded', function() {
    // 1. Verifica autenticación y rol de Admin
    if (!requireAuth()) return;
    if (localStorage.getItem('userRole') !== 'ADMIN') {
        window.location.href = '/rutinas'; 
        return;
    }

    // 2. Carga la lista de usuarios
    cargarUsuarios();
});

/**
 * Llama a la API de admin/users y renderiza la tabla.
 */
async function cargarUsuarios() {
    const tbody = document.getElementById('listaUsuarios');
    
    try {
        const response = await fetch('/api/admin/users', {
            headers: getApiHeaders(false)
        });

        if (!response.ok) {
            if (response.status === 401) logout();
            throw new Error('Error al cargar los usuarios');
        }

        const usuarios = await response.json() || [];
        
        if (usuarios.length === 0) {
            tbody.innerHTML = '<tr><td colspan="4" class="text-center text-muted">No se encontraron usuarios.</td></tr>';
            return;
        }

        tbody.innerHTML = ''; // Limpiar spinner
        
        usuarios.forEach(user => {
            const tr = document.createElement('tr');
            
            // Formatear la fecha (createdAt)
            const fechaRegistro = new Date(user.created_at).toLocaleDateString('es-ES', {
                year: 'numeric',
                month: 'long',
                day: 'numeric'
            });

            tr.innerHTML = `
                <td>${escapeHtml(user.nombre)}</td>
                <td>${escapeHtml(user.email)}</td>
                <td>
                    <span class="badge ${user.rol === 'ADMIN' ? 'bg-danger' : 'bg-secondary'}">
                        ${escapeHtml(user.rol)}
                    </span>
                </td>
                <td>${fechaRegistro}</td>
            `;
            tbody.appendChild(tr);
        });

    } catch (error) {
        console.error('Error en cargarUsuarios:', error);
        tbody.innerHTML = '<tr><td colspan="4" class="text-center text-danger">No se pudieron cargar los usuarios.</td></tr>';
    }
}

/**
 * Helper para escapar HTML (debe estar en common.js o aquí).
 */
function escapeHtml(text) {
    if (typeof text !== 'string') return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}