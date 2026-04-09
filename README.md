# Gym Management System 🏋️‍♂️

Potente sistema de gestión de gimnasios desarrollado en **Go**, diseñado para ser eficiente, escalable y fácil de desplegar mediante contenedores.

## 🚀 Guía de Uso (Manual de Ejecución)
Para probar la herramienta en tu entorno local, sigue estos pasos:

* **Requisitos**: Tener instalado Docker Desktop.

* **Clonar el repositorio**: `git clone https://github.com/FranzHeinen/Gym-Managment-System.git`

* **Desplegar con Docker**: Abrir una terminal en la raíz del proyecto y ejecutar: `docker-compose up --build.`

* **Acceso al Sistema**: Abrir el navegador en http://localhost:8080 para ver el Dashboard o en /swagger/index.html para la API.

* **Base de Datos**: El sistema levantará automáticamente una instancia de MongoDB configurada.

## 🛠️ Tecnologías Core
* **Lenguaje:** Go (Golang)
* **Base de Datos:** MongoDB (NoSQL)
* **Infraestructura:** Docker & Docker Compose
* **Seguridad:** Autenticación con JWT (JSON Web Tokens)
* **Frontend:** Templates HTML dinámicos con JavaScript vanilla

## 🏗️ Arquitectura del Proyecto
El sistema sigue el patrón de diseño por capas para una mejor mantenibilidad:
* **Handlers:** Gestión de peticiones HTTP y rutas.
* **Services:** Lógica de negocio centralizada.
* **Repositories:** Interacción directa con MongoDB.
* **Models:** Definición de esquemas de datos (Usuarios, Rutinas, Ejercicios, Workouts).

## 🐳 Despliegue con Docker
Este proyecto está preparado para funcionar en entornos contenedores. Utiliza Docker Desktop para levantar la base de datos MongoDB y el backend de forma orquestada, asegurando que el entorno de desarrollo sea idéntico al de producción.

## 🌟 Funcionalidades
* **Gestión de Usuarios:** Registro, login y perfiles con roles (Admin/Usuario).
* **Administración de Entrenamientos:** Creación de rutinas personalizadas, series y seguimiento de ejercicios.
* **Panel de Control:** Dashboard para usuarios y vista de administración para gestión de logs y estadísticas.

* ## 👥 Equipo de Desarrollo
* [Franz Heinen](https://github.com/FranzHeinen)
* [Tadeo Isaac](https://github.com/isaactadeo)
* [Felipe Williner](https://github.com/felipewilliner)
