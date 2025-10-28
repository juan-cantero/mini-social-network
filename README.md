# Mini Red Social con Neo4j

Una aplicación de red social basada en consola construida con Go y la base de datos de grafos Neo4j. Este proyecto demuestra conceptos de teoría de grafos y operaciones con bases de datos de grafos, incluyendo relaciones de amistad y sistemas de recomendación.

## Características

- **Gestión de Personas**: Agregar, listar, buscar y eliminar personas
- **Gestión de Amistades**: Crear, visualizar y eliminar amistades bidireccionales
- **Sistema de Recomendaciones**:
  - Recomendaciones basadas en ciudad (sugerir amigos de la misma ciudad)
  - Recomendaciones basadas en hobby (sugerir amigos con hobbies compartidos)
- **Estadísticas**: Ver el total de personas y amistades en la red
- **Integridad de Datos**: Nombres únicos, relaciones simétricas, sin amistades duplicadas

## Stack Tecnológico

- **Lenguaje**: Go
- **Base de Datos**: Neo4j 5 (Base de Datos de Grafos)
- **Driver**: neo4j-go-driver v5
- **Contenedorización**: Docker

## Prerequisitos

- [Go](https://golang.org/doc/install) 1.21 o superior
- [Docker](https://docs.docker.com/get-docker/) y Docker Compose
- Git

## Instalación y Configuración

### 1. Clonar el Repositorio

```bash
git clone https://github.com/juan-cantero/mini-social-network.git
cd mini-social-network
```

### 2. Iniciar la Base de Datos Neo4j

Usando Docker Compose (recomendado):

```bash
docker-compose up -d
```

O usando Docker directamente:

```bash
docker run --name neo4j-social -p 7474:7474 -p 7687:7687 -d -e NEO4J_AUTH=neo4j/password123 neo4j:5
```

### 3. Verificar que Neo4j está Corriendo

Accede al navegador de Neo4j en [http://localhost:7474](http://localhost:7474)

- **Usuario**: neo4j
- **Contraseña**: password123

### 4. Instalar Dependencias de Go

```bash
go mod download
```

### 5. Ejecutar la Aplicación

```bash
go run .
```

O compilar y ejecutar:

```bash
go build -o social-network
./social-network
```

## Uso

Una vez que la aplicación inicia, verás el menú principal:

```
=== MINI SOCIAL NETWORK ===
1. Add person
2. List all people
3. Search person
4. Create friendship
5. View friends of a person
6. Delete friendship
7. City-based recommendations
8. Recommendations by hobby
9. Statistics
0. Exit
```

### Flujo de Trabajo de Ejemplo

1. **Agregar algunas personas**:
   - Ana, Córdoba, Lectura
   - Juan, Córdoba, Videojuegos
   - Maria, Buenos Aires, Lectura
   - Carlos, Córdoba, Videojuegos

2. **Crear amistades**:
   - Ana ← → Juan

3. **Obtener recomendaciones**:
   - Basadas en ciudad para Ana: Sugerirá a Carlos (misma ciudad, no son amigos)
   - Basadas en hobby para Ana: Sugerirá a Maria (mismo hobby, no son amigos)

4. **Ver estadísticas**: Ver el total de personas y amistades

## Modelo de Datos

### Estructura de Nodos

```cypher
(:Persona {
  nombre: "Nombre",
  ciudad: "Ciudad",
  hobby: "Hobby"
})
```

### Estructura de Relaciones

```cypher
(persona1)-[:AMIGO_DE]-(persona2)
```

Las amistades son bidireccionales, lo que significa que se crean ambas direcciones:
- `(persona1)-[:AMIGO_DE]->(persona2)`
- `(persona2)-[:AMIGO_DE]->(persona1)`

## Consultas Cypher Clave

### Crear una Persona

```cypher
CREATE (:Persona {nombre: "Ana", ciudad: "Córdoba", hobby: "Lectura"})
```

### Buscar una Persona

```cypher
MATCH (p:Persona {nombre: "Ana"}) RETURN p
```

### Crear Amistad

```cypher
MATCH (a:Persona {nombre: "Ana"}), (b:Persona {nombre: "Juan"})
CREATE (a)-[:AMIGO_DE]->(b), (b)-[:AMIGO_DE]->(a)
```

### Recomendaciones Basadas en Ciudad

```cypher
MATCH (p:Persona {nombre: "Ana"})
MATCH (candidato:Persona {ciudad: p.ciudad})
WHERE candidato <> p AND NOT (p)-[:AMIGO_DE]-(candidato)
RETURN candidato.nombre
```

### Recomendaciones Basadas en Hobby

```cypher
MATCH (p:Persona {nombre: "Ana"})
MATCH (candidato:Persona {hobby: p.hobby})
WHERE candidato <> p AND NOT (p)-[:AMIGO_DE]-(candidato)
RETURN candidato.nombre
```

### Estadísticas

```cypher
// Contar personas
MATCH (p:Persona) RETURN count(p) AS total

// Contar amistades
MATCH ()-[r:AMIGO_DE]->() RETURN count(r)/2 AS total
```

## Estructura del Proyecto

```
mini-social-network/
├── main.go              # Aplicación principal con menú de consola
├── database.go          # Operaciones de base de datos Neo4j
├── docker-compose.yml   # Configuración de Docker Compose
├── go.mod              # Definición del módulo Go
├── go.sum              # Checksum de dependencias Go
├── .gitignore          # Archivo de ignorados de Git
└── README.md           # Este archivo
```

## Aspectos Destacados del Código

### Gestión de Conexiones

La aplicación utiliza el driver de Go para Neo4j para establecer y gestionar conexiones a la base de datos:

```go
driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
```

### Operaciones Basadas en Sesiones

Cada operación de base de datos utiliza una sesión para ejecutar consultas Cypher:

```go
session := sn.driver.NewSession(sn.ctx, neo4j.SessionConfig{})
defer session.Close(sn.ctx)
```

### Relaciones Bidireccionales

Las amistades se crean simétricamente para asegurar la integridad del grafo:

```go
CREATE (a)-[:AMIGO_DE]->(b), (b)-[:AMIGO_DE]->(a)
```

## Pruebas

Para probar la aplicación manualmente:

1. Inicia Neo4j y la aplicación
2. Agrega varias personas con diferentes ciudades y hobbies
3. Crea amistades entre algunas personas
4. Prueba las recomendaciones para verificar que sugieren no-amigos con propiedades compartidas
5. Ve las estadísticas para verificar los conteos
6. Elimina amistades y verifica que se eliminaron

## Detener la Aplicación

- Presiona `0` para salir de la aplicación
- Detener Neo4j: `docker-compose down` o `docker stop neo4j-social`
- Eliminar el contenedor de Neo4j: `docker rm neo4j-social`

## Preguntas de Análisis

### Pregunta 1: Arquitectura de Base de Datos de Grafos

**¿Por qué Neo4j es más adecuado para esta red social que una base de datos relacional tradicional?**

Neo4j ofrece varias ventajas clave para aplicaciones de redes sociales:

1. **Representación Natural de Relaciones**: En Neo4j, las relaciones son ciudadanos de primera clase almacenados directamente en la base de datos, no como claves foráneas en tablas de unión. Esto hace que modelar redes de amigos sea intuitivo y eficiente.

2. **Rendimiento de Consultas**: Recorrer relaciones en Neo4j es una operación O(1) porque las relaciones se almacenan físicamente como punteros. En una base de datos relacional, encontrar amigos de amigos requiere múltiples operaciones JOIN que se vuelven exponencialmente más lentas a medida que la red crece.

3. **Simplicidad de Cypher**: Las consultas complejas de relaciones son simples en Cypher. Por ejemplo, encontrar recomendaciones de amigos-de-amigos en SQL requiere múltiples auto-joins y cláusulas WHERE complejas, mientras que en Cypher es:
   ```cypher
   MATCH (yo)-[:AMIGO_DE]->()-[:AMIGO_DE]->(recomendado)
   WHERE NOT (yo)-[:AMIGO_DE]-(recomendado) AND recomendado <> yo
   RETURN recomendado
   ```

4. **Flexibilidad de Esquema**: Agregar nuevas propiedades o tipos de relaciones no requiere migraciones de esquema o sentencias ALTER TABLE.

5. **Coincidencia de Patrones**: Neo4j sobresale en encontrar patrones en datos conectados, lo cual es fundamental para redes sociales (amigos mutuos, clusters, influencers, etc.).

### Pregunta 2: Escalabilidad del Sistema y Rendimiento

**Desafíos de escalabilidad y optimizaciones para millones de usuarios:**

**Desafíos:**

1. **Restricciones de Memoria**: Cargar todos los nodos y relaciones en memoria se vuelve impráctico
2. **Rendimiento de Consultas**: Sin optimización, las consultas de recomendación podrían escanear millones de nodos
3. **Throughput de Escritura**: Crear muchas amistades concurrentemente requiere gestión de transacciones
4. **Latencia de Red**: Un único servidor se convierte en un cuello de botella

**Optimizaciones:**

**1. Estrategia de Indexación:**

```cypher
// Crear índices en propiedades consultadas frecuentemente
CREATE INDEX person_name FOR (p:Persona) ON (p.nombre);
CREATE INDEX person_city FOR (p:Persona) ON (p.ciudad);
CREATE INDEX person_hobby FOR (p:Persona) ON (p.hobby);

// Crear índice compuesto para consultas de recomendación
CREATE INDEX person_city_hobby FOR (p:Persona) ON (p.ciudad, p.hobby);
```

Beneficios:
- Reduce el tiempo de búsqueda de O(n) a O(log n)
- Acelera las operaciones MATCH de 10 a 100 veces
- Esencial para el filtrado con cláusulas WHERE

**2. Optimización de Consultas con LIMIT y Paginación:**

```cypher
// En lugar de retornar todas las recomendaciones:
MATCH (p:Persona {nombre: $nombre})
MATCH (candidato:Persona {ciudad: p.ciudad})
WHERE candidato <> p AND NOT (p)-[:AMIGO_DE]-(candidato)
RETURN candidato
LIMIT 10  // Retornar solo las 10 mejores recomendaciones
```

```cypher
// Implementar paginación:
MATCH (p:Persona {nombre: $nombre})
MATCH (candidato:Persona {ciudad: p.ciudad})
WHERE candidato <> p AND NOT (p)-[:AMIGO_DE]-(candidato)
RETURN candidato
SKIP $offset LIMIT $pageSize
```

**3. Técnicas Adicionales de Escalabilidad:**

- **Caché**: Usar Redis o caché en memoria para datos accedidos frecuentemente (usuarios populares, conexiones trending)
- **Réplicas de Lectura**: Desplegar Neo4j Causal Cluster con réplicas de lectura para distribuir la carga de consultas
- **Connection Pooling**: Reutilizar conexiones a la base de datos en lugar de crear nuevas por solicitud
- **Operaciones por Lotes**: Insertar múltiples personas o amistades en una única transacción
- **Normalización de Grafos de Propiedades**: Almacenar propiedades pesadas externamente y referenciarlas por ID

**4. Algoritmos de Grafos Avanzados:**

Para redes a gran escala, usar la biblioteca Neo4j Graph Data Science:
- Detección de Comunidades (encontrar clusters de usuarios altamente conectados)
- PageRank (identificar usuarios influyentes)
- Algoritmos de Similitud (mejores recomendaciones basadas en la estructura de la red)

## Contribuciones

Este es un proyecto de aprendizaje. ¡Siéntete libre de hacer fork y experimentar!

## Licencia

Licencia MIT - siéntete libre de usar este proyecto con fines de aprendizaje.

## Autor

Creado como parte de un ejercicio de aprendizaje de bases de datos de grafos con Neo4j y Go.

## Recursos

- [Documentación de Neo4j](https://neo4j.com/docs/)
- [Driver de Neo4j para Go](https://github.com/neo4j/neo4j-go-driver)
- [Lenguaje de Consultas Cypher](https://neo4j.com/developer/cypher/)
- [Conceptos de Bases de Datos de Grafos](https://neo4j.com/developer/graph-database/)
