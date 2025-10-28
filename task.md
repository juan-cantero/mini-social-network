Programming Challenge: Mini Social Network with Neo4j
1. Exercise Objective

Create a console application that simulates a mini social network using Neo4j as a graph database. Students will apply graph theory concepts while learning to work with graph databases, managing people and their friendship relationships, and implementing a basic friend recommendation system based on shared properties.

2. Technical Specifications
Data Input:

The application must be able to process and manage:

A collection of people with the following properties: name (required), city (required), and one additional custom property of your choice (e.g., age, profession, hobby).

Bidirectional friendship relationships between people.

Search queries by name and recommendation requests based on shared properties.

Technology Stack (Choose one):

Languages: Python (neo4j driver), JavaScript/Node.js (neo4j-driver), Java (neo4j-java-driver), C#/.NET (Neo4j.Driver), Go (neo4j-go-driver), or PHP (laudis/neo4j-php-client).

Database: Neo4j running on Docker.

Expected Output:

The system must provide a complete console menu with the following functionalities:

```
=== MINI SOCIAL NETWORK ===
1. Add person
2. List all people  
3. Search person
4. Create friendship
5. View friends of a person
6. Delete friendship
7. City-based recommendations
8. Recommendations by [your custom property]
9. Statistics
0. Exit
```

3. Rules and Constraints

A valid implementation must meet ALL of the following conditions:

People Management: Add new people with required properties (name, city, custom property), list all registered people, search people by name, and delete people (including their relationships).

Friendship Management: Create bidirectional friendship relationships, list friends of a specific person, delete friendship relationships, and show basic statistics (total people, total friendships).

Recommendation System: Implement TWO types of friend recommendations - by shared city (suggest people from the same city who are not friends) and by shared custom property (suggest people who share the additional property you chose).

Data Integrity: Each person must have a unique name, friendship relationships must be symmetrical, and no duplicate friendships can exist.

4. Example Scenario

Setup Instructions:

Step 1 - Launch Neo4j with Docker:
```bash
docker run --name neo4j-social -p 7474:7474 -p 7687:7687 -d -e NEO4J_AUTH=neo4j/password123 neo4j:5
```

Step 2 - Web Interface Access:
- URL: http://localhost:7474
- User: neo4j
- Password: password123

Data Structure:

Person Node:
```
(:Persona {
  nombre: "Juan Pérez",
  ciudad: "Córdoba",
  [your_property]: "value"
})
```

Friendship Relationship:
```
(person1)-[:AMIGO_DE]-(person2)
```

Sample Cypher Queries:

Create person:
```cypher
CREATE (:Persona {nombre: "Ana", ciudad: "Córdoba"})
```

Search person:
```cypher
MATCH (p:Persona {nombre: "Ana"}) RETURN p
```

Create friendship:
```cypher
MATCH (a:Persona {nombre: "Ana"}), (b:Persona {nombre: "Juan"})
CREATE (a)-[:AMIGO_DE]->(b), (b)-[:AMIGO_DE]->(a)
```

City-based recommendations:
```cypher
MATCH (p:Persona {nombre: "Ana"})
MATCH (candidato:Persona {ciudad: p.ciudad})
WHERE candidato <> p AND NOT (p)-[:AMIGO_DE]-(candidato)
RETURN candidato.nombre
```

5. Post-Implementation Analysis

Once you have a functional solution, answer the following questions to demonstrate your understanding of the problem and the implemented solution.

Question 1: Graph Database Architecture Why is Neo4j (a graph database) more suitable for this social network application compared to a traditional relational database? Describe the specific advantages the graph model provides for relationship queries and recommendations, and explain how Cypher queries simplify complex relationship traversals.

Question 2: System Scalability and Performance Analyze your current implementation's scalability. What challenges would arise if your social network grew to millions of users? Describe at least two specific optimizations you would implement to handle large-scale data while maintaining good performance for friend recommendations. Consider indexing strategies and query optimization techniques.

6. Delivery Format

Once the implementation is complete and the analysis questions are answered, you must submit a GitHub repository. The repository must contain:

The complete source code of your application, properly organized and documented.

A comprehensive README.md file with clear installation and usage instructions.

Documented commit history showing development progress.

Code comments explaining key functionality and Cypher queries.