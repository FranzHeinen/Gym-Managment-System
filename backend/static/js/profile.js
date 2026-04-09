const BASE_API_URL = "/api";

document.addEventListener('DOMContentLoaded', function() {
    if (!isAuthenticated()) {
        window.location.href = '/login';
        return;
    }

    const profileForm = document.getElementById('profileForm');
    const passwordForm = document.getElementById('passwordForm');

    loadProfile(); 

    profileForm.addEventListener('submit', handleProfileUpdate);
    passwordForm.addEventListener('submit', handlePasswordChange);
});

async function loadProfile() {
    try {
        const response = await fetch(`${BASE_API_URL}/profile`, {
            headers: getApiHeaders(false) // getApiHeaders() añade el token de autenticación
        });

        if (!response.ok) {
            throw new Error('Error al cargar el perfil desde la API.');
        }
        
        const data = await response.json();

        document.getElementById('nombre').value = data.nombre || '';
        document.getElementById('email').value = data.email || '';
        if (data.fecha_nacimiento) {
            // El formato de fecha de la BD (ISO 8601) se corta para el input type="date"
            document.getElementById('fecha_nacimiento').value = data.fecha_nacimiento.split('T')[0];
        }
        document.getElementById('peso').value = data.peso || '';
        document.getElementById('altura').value = data.altura || '';
        document.getElementById('nivel_experiencia').value = data.nivel_experiencia || 'Principiante';

        const objetivos = data.objetivos || [];
        document.querySelectorAll('#objetivos .form-check-input').forEach(checkbox => {
            checkbox.checked = objetivos.includes(checkbox.value);
        });

    } catch (error) {
        console.error('Error en loadProfile:', error);
        showMessage('No se pudo cargar tu perfil. Intenta recargar la página.', 'danger', 'message-profile');
    }
}


async function handleProfileUpdate(e) {
    e.preventDefault();
    const submitBtn = document.getElementById('submitProfileBtn');
    
    const objetivos = Array.from(document.querySelectorAll('#objetivos .form-check-input:checked')).map(cb => cb.value);

    const profileData = {
        nombre: document.getElementById('nombre').value,
        fecha_nacimiento: new Date(document.getElementById('fecha_nacimiento').value),
        peso: parseFloat(document.getElementById('peso').value) || 0,
        altura: parseInt(document.getElementById('altura').value) || 0,
        nivel_experiencia: document.getElementById('nivel_experiencia').value,
        objetivos: objetivos
    };
    
    submitBtn.disabled = true;
    submitBtn.innerHTML = '<span class="spinner-border spinner-border-sm me-2"></span>Guardando...';

    try {
        const response = await fetch(`${BASE_API_URL}/profile`, {
            method: 'PUT', 
            headers: getApiHeaders(),
            body: JSON.stringify(profileData)
        });
        const data = await response.json();
        
        showMessage(data.message || 'Error al actualizar', response.ok ? 'success' : 'danger', 'message-profile');

    } catch (error) {
        console.error('Error en handleProfileUpdate:', error);
        showMessage('Error de conexión al guardar el perfil.', 'danger', 'message-profile');
    } finally {
        submitBtn.disabled = false;
        submitBtn.innerHTML = '<i class="bi bi-save me-2"></i>Guardar Cambios';
    }
}

async function handlePasswordChange(e) {
    e.preventDefault();
    const submitBtn = document.getElementById('submitPasswordBtn');
    
    const passwordData = {
        old_password: document.getElementById('old_password').value,
        new_password: document.getElementById('new_password').value
    };

    submitBtn.disabled = true;
    submitBtn.innerHTML = '<span class="spinner-border spinner-border-sm me-2"></span>Actualizando...';

    try {
        const response = await fetch(`${BASE_API_URL}/profile/change-password`, {
            method: 'POST',
            headers: getApiHeaders(),
            body: JSON.stringify(passwordData)
        });
        const data = await response.json();
        
        if (response.ok) {
            showMessage(data.message, 'success', 'message-password');
            e.target.reset(); 
        } else {
            showMessage(data.error || 'Ocurrió un error', 'danger', 'message-password');
        }
    } catch (error) {
        console.error('Error en handlePasswordChange:', error);
        showMessage('Error de conexión al cambiar la contraseña.', 'danger', 'message-password');
    } finally {
        submitBtn.disabled = false;
        submitBtn.innerHTML = '<i class="bi bi-key me-2"></i>Actualizar Contraseña';
    }
}

function showMessage(message, type, elementId) {
    const messageDiv = document.getElementById(elementId);
    messageDiv.innerHTML = `
        <div class="alert alert-${type} alert-dismissible fade show" role="alert">
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        </div>
    `;
    setTimeout(() => {
        const alert = messageDiv.querySelector('.alert');
        if (alert) new bootstrap.Alert(alert).close();
    }, 4000);
}

function isAuthenticated() {
    return !!localStorage.getItem('token');
}

function getApiHeaders() {
    const token = localStorage.getItem('token');
    return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
    };
}