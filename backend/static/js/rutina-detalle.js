// URL base de la API real
const BASE_API_URL = "/api";

let rutinaActual = null;
let ejerciciosDisponibles = [];
let ejercicioEditando = null;
let rutinaId = null;

/**
 * Se ejecuta cuando el DOM está completamente cargado.
 */
document.addEventListener('DOMContentLoaded', async function() {
    if (!requireAuth()) return;

    // 1. Obtener el ID de la rutina desde la URL
    const params = new URLSearchParams(window.location.search);
    rutinaId = params.get('id');
    if (!rutinaId) {
        mostrarError('No se especificó un ID de rutina.');
        return;
    }

    // 2. Cargar todos los datos necesarios
    await cargarDatos();
    configurarEventos();
});

/**
 * Carga los datos de la rutina y la lista de todos los ejercicios.
 */
async function cargarDatos() {
    try {
        // Carga la rutina específica y el catálogo de ejercicios en paralelo
        const [rutinaRes, ejerciciosRes] = await Promise.all([
            fetch(`${BASE_API_URL}/rutinas/${rutinaId}`, { headers: getApiHeaders() }),
            fetch(`${BASE_API_URL}/ejercicios`, { headers: getApiHeaders() }) // Llama a la ruta de ejercicios refactorizada
        ]);

        if (rutinaRes.status === 401 || ejerciciosRes.status === 401) return logout();
        if (!rutinaRes.ok) throw new Error('No se pudo cargar la rutina.');
        if (!ejerciciosRes.ok) throw new Error('No se pudo cargar el catálogo de ejercicios.');

        rutinaActual = await rutinaRes.json();
        ejerciciosDisponibles = await ejerciciosRes.json() || [];

        actualizarUI();
        llenarSelectEjercicios();
        renderizarEjercicios();
    } catch (error) {
        console.error('Error cargando datos:', error);
        mostrarError('Error al cargar los datos de la rutina. ' + error.message);
    }
}

/**
 * Rellena el <select> del modal con todos los ejercicios disponibles.
 */
function llenarSelectEjercicios() {
    const select = document.getElementById('selectEjercicio');
    select.innerHTML = '<option value="">Seleccionar ejercicio...</option>';
    
    ejerciciosDisponibles.forEach(ej => {
        const option = document.createElement('option');
        option.value = ej.id; // Usa el ID del ejercicio
        option.textContent = ej.nombre;
        option.setAttribute('data-descripcion', ej.descripcion);
        select.appendChild(option);
    });
}

/**
 * Actualiza el título y la descripción de la página.
 */
function actualizarUI() {
    document.getElementById('breadcrumbRutina').textContent = rutinaActual.nombre;
    document.getElementById('rutinaNombreInput').value = rutinaActual.nombre;
    document.getElementById('rutinaDescripcionInput').value = rutinaActual.descripcion || '';   
}

/**
 * Dibuja la lista de ejercicios de la rutina.
 */
function renderizarEjercicios() {
    const container = document.getElementById('listaEjercicios');
    
    if (!rutinaActual.ejercicios || rutinaActual.ejercicios.length === 0) {
        container.innerHTML = `
            <div class="text-center py-5">
                <i class="bi bi-inbox text-muted" style="font-size: 3rem;"></i>
                <h3 class="mt-3">No hay ejercicios</h3>
                <p class="text-muted">Agrega tu primer ejercicio a esta rutina.</p>
            </div>
        `;
        return;
    }

    const ejerciciosOrdenados = [...rutinaActual.ejercicios].sort((a, b) => a.orden - b.orden);
    
    container.innerHTML = ejerciciosOrdenados.map((ejercicio) => {
        
        console.log('Renderizando ejercicio:', ejercicio);

        // Busca la info completa del ejercicio en nuestro catálogo
        const ejercicioInfo = ejerciciosDisponibles.find(e => e.id === ejercicio.ejercicio_id) || {};
        
        // Usar dataset attributes en lugar de onclick con parámetros
        return `
            <div class="card mb-3 ejercicio-card" data-ejercicio-id="${ejercicio.id}">
                <div class="card-header">
                    <div class="d-flex justify-content-between align-items-center">
                        <div>
                            <h5 class="mb-0">
                                <span class="badge bg-secondary me-2">${ejercicio.orden}</span>
                                ${ejercicioInfo.nombre || 'Ejercicio no encontrado'}
                            </h5>
                        </div>
                        <div class="d-flex gap-2">
                            <button class="btn btn-sm btn-outline-primary btn-editar" data-ejercicio-id="${ejercicio.id}">
                                <i class="bi bi-pencil"></i>
                            </button>
                            <button class="btn btn-sm btn-outline-danger btn-eliminar" data-ejercicio-id="${ejercicio.id}">
                                <i class="bi bi-trash"></i>
                            </button>
                        </div>
                    </div>
                </div>
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-center mb-3">
                        <small class="text-muted">
                            <i class="bi bi-clock me-1"></i>Descanso: ${ejercicio.tiempo_descanso || 60}s
                        </small>
                        <small class="text-muted">
                            ${ejercicio.series?.length || 0} series
                        </small>
                    </div>
                    
                    <div class="table-responsive">
                        <table class="table table-sm">
                            <thead><tr><th>Serie</th><th>Reps</th><th>Peso (kg)</th></tr></thead>
                            <tbody>
                                ${(ejercicio.series || []).map((serie) => `
                                    <tr>
                                        <td>${serie.numero_serie}</td>
                                        <td>${serie.repeticiones}</td>
                                        <td>${serie.peso}</td>
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        `;
    }).join('');

    // Agregar event listeners después de renderizar
    agregarEventListenersEjercicios();
}

/**
 * Agrega event listeners a los botones de editar y eliminar
 */
function agregarEventListenersEjercicios() {
    // Botones de eliminar
    document.querySelectorAll('.btn-eliminar').forEach(btn => {
        btn.addEventListener('click', function() {
            const ejercicioId = this.getAttribute('data-ejercicio-id');
            console.log('Botón eliminar clickeado. ID leído del atributo data-:', ejercicioId);
            eliminarEjercicio(ejercicioId);
        });
    });

    // Botones de editar
    document.querySelectorAll('.btn-editar').forEach(btn => {
        btn.addEventListener('click', function() {
            const ejercicioId = this.getAttribute('data-ejercicio-id');
            console.log('Editar ejercicio ID:', ejercicioId); // Para debug
            editarEjercicio(ejercicioId);
        });
    });
}

/**
 * Configura los event listeners para los formularios.
 */
function configurarEventos() {
    document.getElementById('formEjercicio').addEventListener('submit', guardarEjercicio);
    document.getElementById('formDetallesRutina').addEventListener('submit', guardarDetallesRutina);
}

async function guardarDetallesRutina(e) {
    e.preventDefault();
    const btn = document.getElementById('btnGuardarDetallesRutina');
    const nombre = document.getElementById('rutinaNombreInput').value;
    const descripcion = document.getElementById('rutinaDescripcionInput').value;

    if (!nombre) {
        mostrarAlertaDetalles('El nombre no puede estar vacío', 'danger');
        return;
    }

    const originalText = btn.innerHTML;
    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm" role="status"></span>';

    try {
        const response = await fetch(`${BASE_API_URL}/rutinas/${rutinaId}`, {
            method: 'PUT',
            headers: getApiHeaders(),
            body: JSON.stringify({
                nombre: nombre,
                descripcion: descripcion
            })
        });

        if (!response.ok) {
            const err = await response.json();
            throw new Error(err.error || 'Error al guardar');
        }

        // Actualiza el estado local y la UI
        rutinaActual.nombre = nombre;
        rutinaActual.descripcion = descripcion;
        actualizarUI(); // Esto actualiza el breadcrumb y los inputs
        mostrarAlertaDetalles('¡Rutina actualizada!', 'success');

    } catch (error) {
        console.error('Error guardando detalles rutina:', error);
        mostrarAlertaDetalles(error.message, 'danger');
    } finally {
        btn.disabled = false;
        btn.innerHTML = originalText;
    }
}

/**
 * Abre el modal para agregar un nuevo ejercicio.
 */
function agregarEjercicio() {
    ejercicioEditando = null;
    document.getElementById('modalEjercicioTitle').textContent = 'Agregar Ejercicio';
    document.getElementById('formEjercicio').reset();
    
// Limpia las series y deja una por defecto
    document.getElementById('configSeries').innerHTML = `
        <div class="series-item row g-2 mb-2">
            <div class="col-3">
                <input type="number" class="form-control" placeholder="Reps" name="repeticiones[]" min="1" value="10" required>
            </div>
            <div class="col-3">
                <input type="number" class="form-control" placeholder="Peso (kg)" name="peso[]" min="0" step="0.5" value="0" required>
            </div>
            <div class="col-4">
                <div class="form-check pt-2">
                  <input class="form-check-input" type="checkbox" name="completada[]" >
                  <label class="form-check-label small">Completada</label>
                </div>
            </div>
            <div class="col-2">
                <button type="button" class="btn btn-outline-danger btn-sm" onclick="eliminarSerie(this)" disabled><i class="bi bi-trash"></i></button>
            </div>
        </div>
    `;
    
    const modal = new bootstrap.Modal(document.getElementById('modalEjercicio'));
    modal.show();
}

/**
 * Abre el modal para editar un ejercicio existente.
 */
function editarEjercicio(ejercicioRutinaID) {
    const ejercicio = rutinaActual.ejercicios.find(e => e.id === ejercicioRutinaID);
    if (!ejercicio) return;

    ejercicioEditando = ejercicio;
    document.getElementById('modalEjercicioTitle').textContent = 'Editar Ejercicio';
    
    document.getElementById('selectEjercicio').value = ejercicio.ejercicio_id;
    document.querySelector('input[name="orden"]').value = ejercicio.orden;
    document.querySelector('input[name="tiempo_descanso"]').value = ejercicio.tiempo_descanso || 60;
    
    const seriesContainer = document.getElementById('configSeries');
    seriesContainer.innerHTML = ''; // Limpia series anteriores
    
    (ejercicio.series || []).forEach((serie, index) => {
        const serieHTML = `
            <div class="series-item row g-2 mb-2">
                <div class="col-3">
                    <input type="number" class="form-control" placeholder="Reps" name="repeticiones[]" min="1" value="${serie.repeticiones}" required>
                </div>
                <div class="col-3">
                    <input type="number" class="form-control" placeholder="Peso (kg)" name="peso[]" min="0" step="0.5" value="${serie.peso}" required>
                </div>
                <div class="col-4">
                    <div class="form-check pt-2">
                      <input class="form-check-input" type="checkbox" name="completada[]" ${serie.completada ? 'checked' : ''} >
                      <label class="form-check-label small">Completada</label>
                    </div>
                </div>
                <div class="col-2">
                    <button type="button" class="btn btn-outline-danger btn-sm" onclick="eliminarSerie(this)"><i class="bi bi-trash"></i></button>
                </div>
            </div>
        `;
        seriesContainer.innerHTML += serieHTML;
    });
    
    const modal = new bootstrap.Modal(document.getElementById('modalEjercicio'));
    modal.show();
}

/**
 * Llama a la API para guardar (Crear o Actualizar) un ejercicio en la rutina.
 */
async function guardarEjercicio(e) {
    e.preventDefault();
    const formData = new FormData(e.target);

    const seriesItems = document.querySelectorAll('#configSeries .series-item');

    const series = Array.from(seriesItems).map((item, index) => {
    const repeticiones = item.querySelector('input[name="repeticiones[]"]').value;
    const peso = item.querySelector('input[name="peso[]"]').value;
    const completada = item.querySelector('input[name="completada[]"]').checked; 

    return {
        numero_serie: index + 1,
        repeticiones: parseInt(repeticiones) || 0,
        peso: parseFloat(peso) || 0,
        completada: completada 
    };
    });

    const ejercicioData = {
    ejercicio_id: formData.get('ejercicio_id'),
    orden: parseInt(formData.get('orden')),
    tiempo_descanso: parseInt(formData.get('tiempo_descanso')) || 60,
    series: series 
    };

    let url = `${BASE_API_URL}/rutinas/${rutinaId}/ejercicios`;
    let method = 'POST';

    if (ejercicioEditando) {
        console.log('Editando ejercicio. ID usado en URL:', ejercicioEditando.id);
        url = `${BASE_API_URL}/rutinas/${rutinaId}/ejercicios/${ejercicioEditando.id}`;
        method = 'PUT';
    }

    try {
        const response = await fetch(url, {
            method: method,
            headers: getApiHeaders(),
            body: JSON.stringify(ejercicioData)
        });

        if (!response.ok) {
            const err = await response.json();
            throw new Error(err.error || 'Error al guardar el ejercicio');
        }
        
        // Todo salió bien, cerramos modal y recargamos los datos
        bootstrap.Modal.getInstance(document.getElementById('modalEjercicio')).hide();
        await cargarDatos(); // Recarga toda la rutina desde la API
        
    } catch (error) {
        console.error('Error guardando ejercicio:', error);
        alert('Error: ' + error.message);
    }
}

/**
 * Llama a la API para eliminar un ejercicio de la rutina.
 */
async function eliminarEjercicio(ejercicioRutinaID) {
    if (!confirm('¿Estás seguro de que quieres eliminar este ejercicio de la rutina?')) return;
    
    try {
        const response = await fetch(`${BASE_API_URL}/rutinas/${rutinaId}/ejercicios/${ejercicioRutinaID}`, {
            method: 'DELETE',
            headers: getApiHeaders()
        });

        if (!response.ok) {
            const err = await response.json();
            throw new Error(err.error || 'Error al eliminar el ejercicio');
        }
        
        await cargarDatos(); // Recarga toda la rutina
        
    } catch (error) {
        console.error('Error eliminando ejercicio:', error);
        alert('Error: ' + error.message);
    }
}

// --- Funciones de Series (dentro del modal) ---

function agregarSerie() {
    const seriesContainer = document.getElementById('configSeries');
    const nuevaSerie = document.createElement('div');
    nuevaSerie.className = 'series-item row g-2 mb-2';
    nuevaSerie.innerHTML = `
        <div class="col-3">
            <input type="number" class="form-control" placeholder="Reps" name="repeticiones[]" min="1" value="10" required>
        </div>
        <div class="col-3">
            <input type="number" class="form-control" placeholder="Peso (kg)" name="peso[]" min="0" step="0.5" value="0" required>
        </div>
        <div class="col-4">
            <div class="form-check pt-2">
              <input class="form-check-input" type="checkbox" name="completada[]" >
              <label class="form-check-label small">Completada</label>
            </div>
        </div>
        <div class="col-2">
            <button type="button" class="btn btn-outline-danger btn-sm" onclick="eliminarSerie(this)"><i class="bi bi-trash"></i></button>
        </div>
    `;
    seriesContainer.appendChild(nuevaSerie);
}

function eliminarSerie(boton) {
    const serieItem = boton.closest('.series-item');
    // No permite eliminar la última serie
    if (document.querySelectorAll('.series-item').length > 1) {
        serieItem.remove();
    }
}

// --- Funciones Auxiliares ---

function volverARutinas() {
    window.location.href = '/rutinas';
}

function mostrarError(mensaje) {
    const container = document.getElementById('listaEjercicios');
    container.innerHTML = `
        <div class="alert alert-danger" role="alert">
            <i class="bi bi-exclamation-triangle me-2"></i>
            ${mensaje}
            <button class="btn btn-sm btn-outline-danger ms-3" onclick="volverARutinas()">
                Volver a Rutinas
            </button>
        </div>
    `;
}

function escapeHtml(text) {
    if (typeof text !== 'string') return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}