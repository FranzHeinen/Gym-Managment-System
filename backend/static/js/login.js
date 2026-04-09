document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');
    const messageDiv = document.getElementById('message');
    const submitBtn = document.getElementById('submitBtn');

    if (isAuthenticated()) {
        window.location.href = '/home';
        return;
    }

    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        if (!email || !password) {
            showMessage('Por favor, completa todos los campos', 'warning');
            return;
        }

        submitBtn.disabled = true;
        submitBtn.innerHTML = '<span class="spinner-border spinner-border-sm me-2"></span>Iniciando sesión...';

        try {
            const response = await fetch('/api/auth/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password })
            });
            const data = await response.json();

            if (response.ok) {
                localStorage.setItem('token', data.token);
                // Guardamos el rol que nos devuelve la API
                localStorage.setItem('userRole', data.user.rol);
                localStorage.setItem('userName', data.user.nombre);
                
                showMessage('¡Inicio de sesión exitoso! Redirigiendo...', 'success');
                setTimeout(() => {
                    window.location.href = '/home'; 
                }, 1500);
            } else {
                showMessage(data.error || 'Error al iniciar sesión', 'danger');
            }
        } catch (error) {
            console.error('Error:', error);
            showMessage('Error de conexión. Inténtalo de nuevo.', 'danger');
        } finally {
            submitBtn.disabled = false;
            submitBtn.innerHTML = '<i class="bi bi-box-arrow-in-right me-2"></i>Iniciar Sesión';
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