document.addEventListener('DOMContentLoaded', function() {
    if (!requireAuth()) return; // De common.js

    const userRole = localStorage.getItem('userRole');
    const userName = localStorage.getItem('userName') || "Usuario";

    document.querySelector('h1.h3').textContent = `Bienvenido, ${userName}`;

    // --- Definición de Tarjetas ---

    const userCards = [
        {
            title: "Mis Rutinas",
            text: "Crea, edita y gestiona tus rutinas de entrenamiento.",
            link: "/rutinas",
            icon: "bi-list-ul"
        },
        {
            title: "Mi Progreso",
            text: "Visualiza tus estadísticas y el historial de entrenamientos.",
            link: "/dashboard",
            icon: "bi-graph-up"
        },
        {
            title: "Mi Perfil",
            text: "Actualiza tu información personal y tus objetivos.",
            link: "/profile",
            icon: "bi-person-fill"
        }
    ];

    const adminCards = [
        {
            title: "Dashboard Admin",
            text: "Ver estadísticas globales de la plataforma.",
            link: "/admin/dashboard",
            icon: "bi-bar-chart-line-fill"
        },
        {
            title: "Gestionar Ejercicios",
            text: "Añadir, editar o eliminar ejercicios del catálogo.",
            link: "/admin/ejercicios",
            icon: "bi-tools"
        },
        {
            title: "Gestionar Usuarios",
            text: "Ver la lista de usuarios registrados en el sistema.",
            link: "/admin/users",
            icon: "bi-people-fill"
        },
        {
            title: "Ver Logs",
            text: "Revisar los registros de actividad del sistema.",
            link: "/admin/logs",
            icon: "bi-clipboard-data"
        }
    ];

    // --- Renderizado ---

    const userMenu = document.getElementById('userMenuCards');
    renderCards(userCards, userMenu);

    if (userRole === 'ADMIN') {
        const adminMenuContainer = document.getElementById('adminMenuContainer');
        const adminMenu = document.getElementById('adminMenuCards');
        
        adminMenuContainer.style.display = 'block'; // Muestra la sección de admin
        renderCards(adminCards, adminMenu);
    }
});

/**
 * Función para renderizar las tarjetas en un contenedor
 */
function renderCards(cards, container) {
    container.innerHTML = cards.map(card => `
        <div class="col-md-6 col-lg-3">
            <div class="card card-menu h-100">
                <a href="${card.link}" class="card-body text-decoration-none d-flex flex-column">
                    <div class="d-flex align-items-center mb-3">
                        <i class="bi ${card.icon} text-primary" style="font-size: 2rem;"></i>
                        <h5 class="card-title mb-0 ms-3">${card.title}</h5>
                    </div>
                    <p class="card-text text-muted small flex-grow-1">${card.text}</p>
                    <span class="text-primary fw-bold small">Ir ahora <i class="bi bi-arrow-right-short"></i></span>
                </a>
            </div>
        </div>
    `).join('');
}