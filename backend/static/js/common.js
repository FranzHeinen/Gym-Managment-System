function updateAuthSection() {
    const authSection = document.getElementById('authSection');
    if (!authSection) return;
    
    const token = localStorage.getItem('token');
    const userRole = localStorage.getItem('userRole'); // Obtenemos el rol
    const userName = localStorage.getItem('userName') || "Usuario";

    if (token) {
        let menuHTML = `
            <li class="nav-item dropdown">
                <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown">
                    <i class="bi bi-person-circle me-1"></i> ${userName}
                </a>
                <ul class="dropdown-menu dropdown-menu-end">
                    <li><a class="dropdown-item" href="/profile"><i class="bi bi-person-fill me-2"></i>Mi Perfil</a></li>
                    <li><a class="dropdown-item" href="/rutinas"><i class="bi bi-list-ul me-2"></i>Mis Rutinas</a></li>
                    <li><a class="dropdown-item" href="/dashboard"><i class="bi bi-graph-up me-2"></i>Mi Progreso</a></li>
                    <li class="dropdown-divider"></li>`;

        // Si el rol es 'ADMIN', añadimos el enlace de gestión
        if (userRole === 'ADMIN') {
            menuHTML += `<li><a class="dropdown-item" href="/admin/dashboard"><i class="bi bi-bar-chart-line-fill me-2"></i>Dashboard</a></li>
                         <li><a class="dropdown-item" href="/admin/ejercicios"><i class="bi bi-tools me-2"></i>Gestionar Ejercicios</a></li>
                         
                         <li><a class="dropdown-item" href="/admin/users"><i class="bi bi-people-fill me-2"></i>Gestionar Usuarios</a></li>
                         <li><a class="dropdown-item" href="/admin/logs"><i class="bi bi-clipboard-data me-2"></i>Ver Logs</a></li>

                         <li class="dropdown-divider"></li>`;
        }

        menuHTML += `
                    <li><a class="dropdown-item text-danger" href="#" onclick="logout()"><i class="bi bi-box-arrow-right me-2"></i>Cerrar Sesión</a></li>
                </ul>
            </li>`;
        
        authSection.innerHTML = menuHTML;

    } else {
        // Menú para usuario no logueado
        authSection.innerHTML = `
            <li class="nav-item">
                <a class="nav-link" href="/login"><i class="bi bi-box-arrow-in-right me-1"></i>Iniciar Sesión</a>
            </li>
            <li class="nav-item">
                <a class="nav-link" href="/register"><i class="bi bi-person-plus me-1"></i>Registrarse</a>
            </li>
        `;
    }
}

async function logout() { 
    if (confirm('¿Estás seguro de que quieres cerrar sesión?')) {
        
        try {
            await fetch('/api/auth/logout', {
                method: 'POST',
                headers: getApiHeaders(false) 
            });

        } catch (error) {
            console.error("Error al cerrar sesión en el servidor:", error);
        } finally {
            localStorage.removeItem('token');
            localStorage.removeItem('userRole'); 
            localStorage.removeItem('userName'); 
            window.location.href = '/login';
        }
    }
}

function isAuthenticated() {
    const token = localStorage.getItem('token');
    return !!token;
}

function requireAuth() {
    if (!isAuthenticated()) {
        window.location.href = '/login';
        return false;
    }
    return true;
}

function getApiHeaders(needsContentType = true) { 
    const token = localStorage.getItem('token');
    const headers = {
        'Authorization': `Bearer ${token}`
    };

    if (needsContentType) { 
        headers['Content-Type'] = 'application/json';
    }

    return headers;
}

document.addEventListener('DOMContentLoaded', function() {
    const isIndexPage = window.location.pathname === '/';
    
    if (!isIndexPage) {
        updateAuthSection();
        
        window.addEventListener('storage', function(e) {
            if (e.key === 'token' || e.key === 'userRole') {
                updateAuthSection();
            }
        });
    }
});