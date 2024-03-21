const fs = require('fs');

// Function to read and parse Cypher queries from a .cypher file
function parseCypherQueries(filePath) {
    const cypherQueries = [];
    const fileContents = fs.readFileSync(filePath, 'utf-8');
    const queries = fileContents.split('\n\n'); // Assuming queries are separated by double newline

    queries.forEach(query => {
        // Remove comments and trim whitespace
        const cleanedQuery = query.replace(/(\/\/.*)/g, '').trim();
        if (cleanedQuery.length > 0) {
            cypherQueries.push(cleanedQuery);
        }
    });

    return cypherQueries;
}



const neo4j = require('neo4j-driver');

// Neo4j connection details
const uri = 'neo4j+s://30fb144b.databases.neo4j.io';
const user = 'neo4j';
const password = 'RbFNFJaJ0d5pPtLJ1GQkAZJhJgCl-ukDZRoquy1TFf4';

// Create a Neo4j driver instance
const driver = neo4j.driver(uri, neo4j.auth.basic(user, password));

// Function to execute Cypher query
async function executeCypherQuery(query, params = {}) {
    const session = driver.session();
    try {
        const result = await session.run(query, params);
        return result.records.map(record => record.toObject());
    } finally {
        await session.close();
    }
}


// Example usage
async function main() {
    try {
        // Execute a sample Cypher query
        const query = parseCypherQueries('query.cypher')[0];
	console.log(query);
        const result = await executeCypherQuery(query);
        for (elt of result){
	    console.log(elt['b.Title']); };
    } catch (error) {
        console.error('Error executing Cypher query:', error);
    } finally {
        await driver.close();
    }
}

// Call the main function
main();
