document.addEventListener('DOMContentLoaded', function() {
    
   
    // L isAuthenticated() viene de common.js y revisa si existe un token.
    if (isAuthenticated()) {
        // Si existe, significa que el usuario ya inició sesión.
        // Lo redirige a su pag principal de rutinas.
        window.location.href = '/home';
        return; 
    }

    // Si no hay, muestra el menú de Iniciar Sesión y Registrarse.
    updateAuthSection();
    
    // Esto escucha si inicias sesión en otra pestaña, para actualizar el menú.
    window.addEventListener('storage', function(e) {
        if (e.key === 'token') {
            updateAuthSection();
        }
    });
});