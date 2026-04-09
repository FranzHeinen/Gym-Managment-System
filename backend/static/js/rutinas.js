// URL base de la API real
const BASE_API_URL = "/api"; 

// Variables globales
let rutinas = [];
let workouts = []; // Almacenará el historial de workouts

/**
 * Muestra una alerta temporal en la parte superior de la lista de rutinas.
 */
function showAlert(message, type = 'success') {
    const container = document.getElementById('listaRutinas');
    
    const alertDiv = document.createElement('div');
    alertDiv.className = `alert alert-${type} alert-dismissible fade show`;
    alertDiv.role = 'alert';
    alertDiv.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
    `;
    
    container.prepend(alertDiv);
    
    setTimeout(() => {
        const bsAlert = new bootstrap.Alert(alertDiv);
        bsAlert.close();
    }, 4000);
}

/**
 * Se ejecuta cuando el DOM está completamente cargado.
 */
document.addEventListener("DOMContentLoaded", () => {
    if (!requireAuth()) return; // requireAuth() viene de common.js

    document.getElementById("formRutina").addEventListener("submit", async (e) => {
        e.preventDefault();
        const fd = new FormData(e.target);
        const nombre = fd.get("nombre");
        const descripcion = fd.get("descripcion");

        await crearRutina(nombre, descripcion);

        e.target.reset();
        const modal = bootstrap.Modal.getInstance(document.getElementById("modalRutina"));
        if (modal) modal.hide();
    });

    // Carga las rutinas Y el historial de workouts en paralelo
    Promise.all([
        fetchRutinas(),
        fetchWorkouts()
    ]).then(() => {
        renderRutinas(); // Renderizar solo después de tener AMBOS
    }).catch(err => {
        console.error("Error al cargar datos iniciales:", err);
        document.getElementById("listaRutinas").innerHTML = `<p class="text-danger text-center">No se pudieron cargar las rutinas.</p>`;
    });
});

/**
 * Obtiene el historial de workouts del usuario desde la API.
 */
async function fetchWorkouts() {
    try {
        const response = await fetch(`${BASE_API_URL}/workouts`, {
            headers: getApiHeaders(false) // GET, no content-type
        });
        if (!response.ok) {
            if (response.status === 401) return logout();
            throw new Error('Error al obtener workouts');
        }
        workouts = await response.json() || [];
    } catch (error) {
        console.error("Error cargando workouts:", error);
        workouts = []; // Asegurarse que sea un array
    }
}

/**
 * Obtiene todas las rutinas del usuario desde la API.
 */
async function fetchRutinas() {
    const container = document.getElementById("listaRutinas");
    container.innerHTML = '<div class="col-12 text-center"><span class="spinner-border"></span></div>'; 

    try {
        const response = await fetch(`${BASE_API_URL}/rutinas`, {
            headers: getApiHeaders(false) // GET, no content-type
        });

        if (response.status === 401) return logout(); 
        if (!response.ok) throw new Error('Error al obtener las rutinas');

        rutinas = await response.json() || [];
        // No llamamos a renderRutinas() aquí, esperamos a que los workouts también carguen
    } catch (error) {
        console.error(error);
        container.innerHTML = `<p class="text-danger text-center">No se pudieron cargar las rutinas.</p>`;
        throw error; // Lanzar el error para que Promise.all lo capture
    }
}

/**
 * Helper para verificar si una rutina ya se completó HOY.
 */
function isCompletadaHoy(rutinaId) {
    const hoy = new Date().toDateString(); // "Fri Oct 24 2025"

    return workouts.some(workout => {
        const fechaWorkout = new Date(workout.fecha).toDateString();
        return workout.rutina_id === rutinaId && fechaWorkout === hoy;
    });
}

/**
 * Dibuja las tarjetas (cards) de las rutinas en el HTML.
 */
function renderRutinas() {
    const cont = document.getElementById("listaRutinas");
    if (!rutinas || rutinas.length === 0) {
        cont.innerHTML = `<div class="col-12 text-center text-muted mt-5">
                            <i class="bi bi-journal-x" style="font-size: 3rem;"></i>
                            <h4 class="mt-3">Aún no tienes rutinas</h4>
                            <p>¡Crea tu primera rutina para empezar!</p>
                          </div>`;
        return;
    }

    cont.innerHTML = rutinas.map(r => {
        
        // Lógica para el botón de completar
        const completadaHoy = isCompletadaHoy(r.id);
        let botonCompletarHTML = '';
        if (completadaHoy) {
            botonCompletarHTML = `
                <button classD="btn btn-sm btn-success" disabled>
                  <i class="bi bi-check-circle-fill"></i> Completada
                </button>`;
        } else {
            botonCompletarHTML = `
                <button class="btn btn-sm btn-outline-success" onclick="registrarEntrenamiento('${r.id}', '${escapeHtml(r.nombre)}')">
                  <i class="bi bi-check-lg"></i> Completar
                </button>`;
        }

        return `
        <div class="col-md-6 col-lg-4">
          <div class="card h-100 shadow-sm">
            <div class="card-body d-flex flex-column">
              <div class="d-flex justify-content-between align-items-start">
                <h5 class="card-title mb-1">${escapeHtml(r.nombre)}</h5>
                <span class="badge text-bg-secondary">${(r.ejercicios?.length ?? 0)} ej.</span>
              </div>
              <p class="card-text text-muted small flex-grow-1">${escapeHtml(r.descripcion) || "Sin descripción."}</p>

              <div class="d-grid d-md-flex gap-2 justify-content-md-end">
                <button class="btn btn-sm btn-outline-secondary" onclick="duplicar('${r.id}')">
                  <i class="bi bi-copy"></i> Duplicar
                </button>
                <button class="btn btn-sm btn-outline-danger" onclick="eliminar('${r.id}')">
                  <i class="bi bi-trash"></i> Eliminar
                </button>
                
                ${botonCompletarHTML}

                <a href="/rutina-detalle?id=${r.id}" class="btn btn-sm btn-primary">
                  Ver/Editar <i class="bi bi-arrow-right"></i>
                </a>
              </div>
            </div>
          </div>
        </div>
        `;
    }).join("");
}

/**
 * Llama a la API para crear una nueva rutina.
 */
async function crearRutina(nombre, descripcion) {
    try {
        const res = await fetch(`${BASE_API_URL}/rutinas`, {
            method: "POST",
            headers: getApiHeaders(true),
            body: JSON.stringify({ 
                nombre: nombre, 
                descripcion: descripcion,
                ejercicios: [] 
            })
        });

        if (res.ok) {
            await fetchRutinas(); // Recarga la lista de rutinas
            renderRutinas(); // Vuelve a dibujar
        } else {
            const err = await res.json();
            alert("Error al crear la rutina: " + (err.error || "Error desconocido"));
        }
    } catch (error) {
        console.error("Error en crearRutina:", error);
    }
}

/**
 * Llama a la API para eliminar una rutina.
 */
async function eliminar(id) {
    if (!confirm("¿Estás seguro de que quieres eliminar esta rutina?")) return;
    try {
        const res = await fetch(`${BASE_API_URL}/rutinas/${id}`, {
            method: "DELETE",
            headers: getApiHeaders(false) // No hay body
        });

        if (res.ok) {
            await fetchRutinas(); // Recarga
            renderRutinas(); // Re-dibuja
        } else {
             alert("Error al eliminar la rutina.");
        }
    } catch (error) {
        console.error("Error en eliminar:", error);
    }
}

/**
 * Llama a la API para duplicar una rutina.
 */
async function duplicar(id) {
    if (!confirm("¿Duplicar esta rutina?")) return;
    try {
        const res = await fetch(`${BASE_API_URL}/rutinas/${id}/duplicate`, {
            method: "POST",
            headers: getApiHeaders(false) // No hay body
        });

        if (res.ok) {
            await fetchRutinas(); // Recarga
            renderRutinas(); // Re-dibuja
        } else {
            alert("Error al duplicar la rutina.");
        }
    } catch (error) {
        console.error("Error en duplicar:", error);
    }
}

/**
 * Llama a la API para registrar un entrenamiento (workout) completado.
 */
async function registrarEntrenamiento(id, nombre) {
    if (!confirm(`¿Confirmas que has completado la rutina "${nombre}"?`)) return;

    try {
        const res = await fetch(`${BASE_API_URL}/workouts`, {
            method: "POST",
            headers: getApiHeaders(true), // 'true' para Content-Type
            body: JSON.stringify({ 
                rutina_id: id
            })
        });

        if (res.ok) {
            showAlert(`¡Bien hecho! Rutina "${nombre}" registrada.`, 'success');
            // Actualizamos el historial local y volvemos a renderizar
            await fetchWorkouts();
            renderRutinas();
        } else {
            const err = await res.json();
            showAlert(`Error al registrar: ${err.error}`, 'danger');
        }
    } catch (error) {
        console.error("Error en registrarEntrenamiento:", error);
        showAlert('Error de conexión al registrar el entrenamiento.', 'danger');
    }
}

/**
 * Helper para escapar HTML y prevenir ataques XSS.
 */
function escapeHtml(text) {
    if (typeof text !== 'string') return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}