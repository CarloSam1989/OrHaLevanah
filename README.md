# Luna Bíblica App

Aplicación web en Go para mostrar:

- fase lunar actual
- inicio de mes bíblico por observación
- comparación entre ubicación actual y Jerusalén
- fiestas de Levítico 23
- calendario lunar bíblico

## Enfoque
Este proyecto usará una lógica de raíces hebreas, tomando Jerusalén como referencia para el inicio observacional del mes por primer creciente visible.

La fecha civil de «Hoy» y la marca de hoy en el calendario siguen la zona horaria local del dispositivo. Los cálculos del calendario bíblico, las fiestas y el atardecer mantienen a Jerusalén como referencia.

## Tecnologías
- Go
- HTML
- CSS
- JavaScript
- PWA

## Despliegue en Vercel

Vercel publica las páginas estáticas generadas por `npm run build:vercel` en `dist`: la página de la comunidad queda en `/` y el calendario en `/app/`. Las peticiones `/api/*` se reescriben al backend Go de Railway (`orhalevanah-production.up.railway.app`); ese servicio debe estar desplegado y sus rutas API deben responder correctamente. Si las API no están disponibles, el calendario conserva la cuadrícula civil y explica qué datos bíblicos faltan, sin inventar inicios de mes o avistamientos.

En Vercel, configura el directorio raíz como la raíz del repositorio y el preset como **Other**. `vercel.json` establece `"framework": null` para impedir que la presencia de `go.mod` active el constructor Go: este despliegue compila un sitio estático, no un binario Go. Usa `npm run build:vercel` como comando de compilación y `dist` como directorio de salida. La instalación puede usar el comando detectado desde `package-lock.json`.

El build estático incluye una copia de solo lectura de `data/moon_sightings.json` en `dist/data/` para mostrar los reportes confirmados guardados cuando la API no está disponible. Esa copia puede quedar desactualizada y no reemplaza la actualización del backend. El despliegue Docker también incluye `data/`; conserva el original para que las APIs puedan leer y actualizar los reportes lunares.