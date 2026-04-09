// URL base de la API real
const BASE_API_URL = "/api";

let ejercicios = [];
let ejercicioEditando = null;
let gruposMusculares = new Set(); // Para llenar el filtro

/**
 * Se ejecuta cuando el DOM está completamente cargado.
 */
document.addEventListener('DOMContentLoaded', function() {
    // 1. Verifica autenticación y rol de Admin
    if (!requireAuth()) return;
    if (localStorage.getItem('userRole') !== 'ADMIN') {
        // Si no es admin, redirige (o muestra un mensaje)
        window.location.href = '/rutinas'; 
        return;
    }

    // 2. Configura los event listeners
    document.getElementById('formEjercicio').addEventListener('submit', guardarEjercicio);
    document.getElementById('modalEjercicio').addEventListener('hidden.bs.modal', limpiarModal);
    document.getElementById('filtroNombre').addEventListener('input', filtrarEjercicios);
    document.getElementById('filtroGrupo').addEventListener('change', filtrarEjercicios);
    document.getElementById('limpiarFiltros').addEventListener('click', limpiarYFiltrar);

    // 3. Carga los ejercicios iniciales
    cargarEjercicios();
});

/**
 * Carga todos los ejercicios desde la API.
 */
async function cargarEjercicios() {
    const container = document.getElementById('listaEjercicios');
    container.innerHTML = '<div class="col-12 text-center"><span class="spinner-border"></span></div>';

    try {
        // Llama a la API (endpoint público para lectura)
        const response = await fetch(`${BASE_API_URL}/ejercicios`, { 
            headers: getApiHeaders() // Necesario por si la ruta GET se protege en el futuro
        });
        
        if (response.status === 401) return logout();
        if (!response.ok) throw new Error('Error al cargar ejercicios');

        ejercicios = await response.json() || [];
        
        // Extraer grupos musculares únicos para el filtro
        gruposMusculares = new Set(ejercicios.map(e => e.grupo_muscular).filter(Boolean));
        llenarFiltroGrupos();
        
        renderizarEjercicios(ejercicios); // Muestra todos al inicio

    } catch (error) {
        console.error(error);
        container.innerHTML = `<p class="text-danger text-center">No se pudieron cargar los ejercicios.</p>`;
    }
}

/**
 * Rellena el <select> de filtro con los grupos musculares.
 */
function llenarFiltroGrupos() {
    const select = document.getElementById('filtroGrupo');
    select.innerHTML = '<option value="">Filtrar por grupo muscular...</option>';
    gruposMusculares.forEach(grupo => {
        const option = document.createElement('option');
        option.value = grupo;
        option.textContent = grupo;
        select.appendChild(option);
    });
}

/**
 * Filtra los ejercicios mostrados según los inputs de nombre y grupo.
 */
function filtrarEjercicios() {
    const nombreFiltro = document.getElementById('filtroNombre').value.toLowerCase();
    const grupoFiltro = document.getElementById('filtroGrupo').value;

    const ejerciciosFiltrados = ejercicios.filter(ej => {
        const nombreCoincide = ej.nombre.toLowerCase().includes(nombreFiltro);
        const grupoCoincide = !grupoFiltro || ej.grupo_muscular === grupoFiltro;
        return nombreCoincide && grupoCoincide;
    });

    renderizarEjercicios(ejerciciosFiltrados);
}

/**
 * Limpia los filtros y vuelve a mostrar todos los ejercicios.
 */
function limpiarYFiltrar() {
    document.getElementById('filtroNombre').value = '';
    document.getElementById('filtroGrupo').value = '';
    renderizarEjercicios(ejercicios);
}

/**
 * Dibuja las tarjetas (cards) de los ejercicios en el HTML.
 */
function renderizarEjercicios(ejerciciosAMostrar) {
    const container = document.getElementById('listaEjercicios');
    if (!ejerciciosAMostrar || ejerciciosAMostrar.length === 0) {
        container.innerHTML = `<p class="text-muted text-center col-12">No se encontraron ejercicios.</p>`;
        return;
    }

    container.innerHTML = ejerciciosAMostrar.map(ej => `
    <div class="col-12 col-md-6 col-lg-4">
        <div class="card h-100 shadow-sm">
            <div class="card-body d-flex flex-column">
                <h5 class="card-title">${escapeHtml(ej.nombre)}</h5>
                <h6 class="card-subtitle mb-2 text-muted">${escapeHtml(ej.grupo_muscular)} - ${escapeHtml(ej.dificultad)}</h6>
                <p class="card-text small flex-grow-1">${escapeHtml(ej.descripcion)}</p>
                <div class="mt-auto d-flex justify-content-end gap-2">
                    <button class="btn btn-sm btn-outline-secondary" onclick="editarEjercicio('${ej.id}')">
                        <i class="bi bi-pencil"></i> Editar
                    </button>
                    <button class="btn btn-sm btn-outline-danger" onclick="eliminarEjercicio('${ej.id}')">
                        <i class="bi bi-trash"></i> Eliminar
                    </button>
                </div>
            </div>
        </div>
    </div>
    `).join('');
}

/**
 * Abre el modal para editar un ejercicio existente.
 */
function editarEjercicio(id) {
    ejercicioEditando = ejercicios.find(e => e.id === id);
    if (!ejercicioEditando) return;

    // Rellenar el formulario del modal con los datos
    document.getElementById('ejercicioId').value = ejercicioEditando.id;
    document.querySelector('#formEjercicio [name="nombre"]').value = ejercicioEditando.nombre;
    document.querySelector('#formEjercicio [name="categoria"]').value = ejercicioEditando.categoria;
    document.querySelector('#formEjercicio [name="grupo_muscular"]').value = ejercicioEditando.grupo_muscular;
    document.querySelector('#formEjercicio [name="dificultad"]').value = ejercicioEditando.dificultad;
    document.querySelector('#formEjercicio [name="descripcion"]').value = ejercicioEditando.descripcion;
    document.querySelector('#formEjercicio [name="instruccion"]').value = ejercicioEditando.instruccion;
    document.querySelector('#formEjercicio [name="demostracion"]').value = ejercicioEditando.demostracion;

    document.getElementById('modalTitle').textContent = 'Editar Ejercicio';
    const modal = new bootstrap.Modal(document.getElementById('modalEjercicio'));
    modal.show();
}

/**
 * Llama a la API para guardar (Crear o Actualizar) un ejercicio.
 */
async function guardarEjercicio(e) {
    e.preventDefault();
    const formData = new FormData(e.target);
    const id = formData.get('id');

    const data = {
        nombre: formData.get('nombre'),
        descripcion: formData.get('descripcion'),
        categoria: formData.get('categoria'),
        grupo_muscular: formData.get('grupo_muscular'),
        dificultad: formData.get('dificultad'),
        demostracion: formData.get('demostracion'),
        instruccion: formData.get('instruccion')
    };

    let url = `${BASE_API_URL}/admin/ejercicios`;
    let method = 'POST';

    if (id) { // Si hay ID, es una actualización (PUT)
        url = `${BASE_API_URL}/admin/ejercicios/${id}`;
        method = 'PUT';
    }

    try {
        const response = await fetch(url, {
            method: method,
            headers: getApiHeaders(), // Envía el token de Admin
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            const err = await response.json();
            throw new Error(err.error || `Error al ${id ? 'actualizar' : 'crear'} ejercicio`);
        }
        
        // Si todo OK, cierra modal y recarga lista
        bootstrap.Modal.getInstance(document.getElementById('modalEjercicio')).hide();
        await cargarEjercicios(); 
        
    } catch (error) {
        console.error('Error guardando ejercicio:', error);
        alert('Error: ' + error.message);
    }
}

/**
 * Llama a la API para eliminar un ejercicio.
 */
async function eliminarEjercicio(id) {
    if (!confirm('¿Estás seguro de que quieres eliminar este ejercicio?')) return;
    
    try {
        const response = await fetch(`${BASE_API_URL}/admin/ejercicios/${id}`, {
            method: 'DELETE',
            headers: getApiHeaders() // Envía el token de Admin
        });

        if (!response.ok) {
            const err = await response.json();
            throw new Error(err.error || 'Error al eliminar el ejercicio');
        }
        
        await cargarEjercicios(); // Recarga la lista
        
    } catch (error) {
        console.error('Error eliminando ejercicio:', error);
        alert('Error: ' + error.message);
    }
}

/**
 * Limpia el formulario del modal cuando se cierra.
 */
function limpiarModal() {
    document.getElementById('formEjercicio').reset();
    document.getElementById('ejercicioId').value = '';
    document.getElementById('modalTitle').textContent = 'Nuevo Ejercicio';
    ejercicioEditando = null;
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