document.addEventListener('DOMContentLoaded', function() {
    const registerForm = document.getElementById('registerForm');
    const messageDiv = document.getElementById('message');
    const submitBtn = document.getElementById('submitBtn');

    if (isAuthenticated()) {
        window.location.href = '/home';
        return;
    }

    registerForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        // Obtenemos el valor del nuevo campo "nombre"
        const nombre = document.getElementById('nombre').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        const confirmPassword = document.getElementById('confirmPassword').value;

        // Validamos también el campo "nombre"
        if (!nombre || !email || !password || !confirmPassword) {
            showMessage('Por favor, completa todos los campos', 'warning');
            return;
        }
        if (password.length < 6) {
            showMessage('La contraseña debe tener al menos 6 caracteres', 'warning');
            return;
        }
        if (password !== confirmPassword) {
            showMessage('Las contraseñas no coinciden', 'warning');
            return;
        }

        submitBtn.disabled = true;
        submitBtn.innerHTML = '<span class="spinner-border spinner-border-sm me-2"></span>Registrando...';

        try {
            const response = await fetch('/api/auth/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                // Añadimos el "nombre" al cuerpo de la petición
                body: JSON.stringify({ nombre, email, password })
            });
            const data = await response.json();

            if (response.ok) {
                localStorage.setItem('token', data.token);
                localStorage.setItem('userRole', data.user.rol);
                localStorage.setItem('userName', data.user.nombre);
                
                showMessage('¡Registro exitoso! Redirigiendo...', 'success');
                setTimeout(() => {
                    window.location.href = '/home';
                }, 1500);
            } else {
                showMessage(data.error || 'Error al registrarse', 'danger');
            }
        } catch (error) {
            console.error('Error:', error);
            showMessage('Error de conexión. Inténtalo de nuevo.', 'danger');
        } finally {
            submitBtn.disabled = false;
            submitBtn.innerHTML = '<i class="bi bi-person-plus me-2"></i>Crear Cuenta';
        }
    });

    function showMessage(message, type) {
        messageDiv.innerHTML = `<div class="alert alert-${type} alert-dismissible fade show" role="alert">${message}<button type="button" class="btn-close" data-bs-dismiss="alert"></button></div>`;
        setTimeout(() => {
            const alert = messageDiv.querySelector('.alert');
            if (alert) new bootstrap.Alert(alert).close();
        }, 5000);
    }
});